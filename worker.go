package main
// package main

import (
  "github.com/google/uuid"
)

// func main() {}

type WorkerID struct {
  uuid.UUID
}

type Worker struct {
  WorkerID          WorkerID
  ProcessingChannel chan<- bool
  FinishedChannel   chan<- WorkerID
}

func NewWorker(processingChannel chan<- bool, finishedChannel chan<- WorkerID) Worker {
  return Worker{
    WorkerID:           generateID(),
    ProcessingChannel:  processingChannel,
    FinishedChannel:    finishedChannel,
  }
}

func generateID() WorkerID {
  return WorkerID{
    UUID: uuid.New(),
  }
}
