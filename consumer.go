package main

import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/sqs"
  // "worker"
  log "github.com/sirupsen/logrus"
  // "fmt"
  "os"
)

const (
  region = "us-east-1"
  queueName = "sqs_create_q.fifo"
)

type SimpleQueue struct {
  client  *sqs.SQS
  url     *sqs.GetQueueUrlOutput
}

// func main() {
//   q := NewSimpleQueue()
//   fmt.Println("q: %v", q)
//   q.pollMessages()
// }

func init() {
  file, err := os.OpenFile("/var/log/anon-penny/1.log", os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666)
  if err != nil {
    panic(err)
  }
  log.SetOutput(file)
  log.SetFormatter(&log.JSONFormatter{})
}

func NewSimpleQueue() SimpleQueue {
  sess, err := session.NewSession(&aws.Config{
    Region: aws.String(region)},
  )

  // Create a SQS service client.
  svc := sqs.New(sess)

  resultUrl, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
    QueueName: aws.String(queueName),
  })

  if err != nil {
    log.Error(err)
    panic(err)
  }

  return SimpleQueue{
    client: svc,
    url: resultUrl,
  }
}

func (queue *SimpleQueue) PollMessages() []*sqs.Message {
  result, err := queue.client.ReceiveMessage(&sqs.ReceiveMessageInput{
    QueueUrl: queue.url.QueueUrl,
    AttributeNames: aws.StringSlice([]string{
      "SentTimestamp",
    }),
    MaxNumberOfMessages: aws.Int64(10),
    MessageAttributeNames: aws.StringSlice([]string{
      "All",
    }),
    WaitTimeSeconds: aws.Int64(20),
  })

  if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
      log.Error(aerr)
    }
    log.Error("Unable to receive message from queue %q, %v", queueName, err)
  }

  log.Info("Received %d messages.\n", len(result.Messages))
  if len(result.Messages) > 0 {
    log.Info(result.Messages)
  }

  return result.Messages
}


// TODO: each worker will be passed this method as its action handler
// (starting from 'main' in myapp; e.g. 'sqs.Consume'),
// where it polls the queue asking for work to do; handles the message, updates the pool
// about its state along the way.
func (queue *SimpleQueue) Consume(worker Worker) {
  log.Info("queue.Consume")
  // signal the pool that the worker has finished after the method is executed
  defer func() {
    worker.FinishedChannel <- worker.WorkerID
  }()

  // poll the queue
  messages := queue.PollMessages()
  // if err != nil {
  //   // signal the pool that there's no message being processed
  //   worker.ProcessingChannel <- false
  //
  //   return err
  // }

  if len(messages) == 0 {
    // signal the pool that there's no message being processed
    worker.ProcessingChannel <- false
    // return nil
  } else {
    // signal the pool that a message will be processed
    worker.ProcessingChannel <- true

    // handle the message...
    log.Info("handle the message...")
  }

  // return nil
}
