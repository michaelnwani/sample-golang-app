package main
// package main

import (
  "time"
  "os"
  "os/signal"
  log "github.com/sirupsen/logrus"
)

type Shutdown struct {
  InitiateChannel chan bool
  DoneChannel     chan bool
  Timeout         time.Duration
}

// func main(){}
func NewShutdown(timeout time.Duration) *Shutdown {
  return &Shutdown{
    InitiateChannel:  make(chan bool, 1),
    DoneChannel:      make(chan bool, 1),
    Timeout:          timeout,
  }
}

// makes the main thread wait for the interrupt signal from the system
// (via Ctrl+C) that will ping the pool to initiate its shutdown
func (shutdown *Shutdown) WaitForSignal() {
  signalChannel := make(chan os.Signal, 1)
  signal.Notify(signalChannel, os.Interrupt)
  <-signalChannel

  log.Info("received interrupt signal")
  shutdown.InitiateChannel <- true
  select {
  case <-signalChannel:
    log.Fatal("forcing shutdown")
    // os.Exit(1) // automatically called by log.Fatal
  case <-shutdown.DoneChannel:
    log.Info("cleanup successful, exiting")
  case <-time.After(time.Second * shutdown.Timeout):
    log.Info("cleanup timed out, exiting")
  }
}
