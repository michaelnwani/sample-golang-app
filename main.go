package main

import (
  "github.com/aws/aws-sdk-go/service/sqs"
  log "github.com/sirupsen/logrus"
  // "consumer"
  "time"
  "encoding/json"
  "os"
)

var (
  processing = make(chan bool, 1)
)

func init() {
  file, err := os.OpenFile("/var/log/anon/1.log", os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666)
  if err != nil {
    panic(err)
  }
  log.SetOutput(file)
  log.SetFormatter(&log.JSONFormatter{})
}

func main() {
  queue := NewSimpleQueue()
  // shutdown := shutdown.NewShutdown(time.Duration(60))
  //
  // workerPool := workerpool.NewWorkerPool(shutdown)
  // go workerPool.Start(queue.Consume)
  //
  // shutdown.WaitForSignal()
  for {
    payload := queue.PollMessages()
    if len(payload) > 0 {
      log.Info("processing messages")
      processMessages(payload)
      log.Info("processing messages: done. blocking")
      <-processing
      log.Info("processing messages: done. blocking: done.")
    } else {
      log.Info("no messages to process. sleeping for 5 seconds")
      time.Sleep(5 * time.Second)
    }
  }
}

func processMessages(payload []*sqs.Message) {
  log.Info("do stuff...")
  time.Sleep(3 * time.Second)
  log.Info("add to processing channel")
  for _, message := range payload {
    var raw map[string]interface{}
    json.Unmarshal([]byte(*message.Body), &raw)
    log.Info("json decoded:",raw)
    elem, ok := raw["body"]
    if ok {
      log.Info("BODY:", elem)
    }
  }
  processing <- false
}
