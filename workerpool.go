package main
// package main

import (
  // "worker"
  // "shutdown"
  "time"
  "log"
)

const (
  stateInitial = iota
  stateMain
  stateExit
  stateLaunch
  stateProcessing
  stateFinish
  stateSleep
  stateWait
  stateTimeout
  stateQuit
)

type WorkerPool struct {
  maxWorkers        int
  interval          int
  timeout           int
  shutdown          *Shutdown
  processingChannel chan bool
  finishedChannel   chan WorkerID
  needsToQuit       bool
  activeWorkers     map[WorkerID]Worker
  currentTransition Transition
}

// updated every cycle
type Transition struct {
  state   int
  payload Payload
}

type Payload struct {
  workerID      WorkerID
  isProcessing  bool
}

// func main(){}
func NewWorkerPool(shutdown *Shutdown) WorkerPool {
  maxWorkers := 1
  initialTransition := Transition {
    state: stateInitial,
  }

  return WorkerPool{
    maxWorkers:         maxWorkers,
    interval:           30,
    timeout:            60,
    shutdown:           shutdown,
    processingChannel:  make(chan bool, maxWorkers),
    finishedChannel:    make(chan WorkerID, maxWorkers),
    activeWorkers:      map[WorkerID]Worker{},
    currentTransition:  initialTransition,
  }
}

// actionHandler func will be called by the individual workers
func (workerPool *WorkerPool) Start(actionHandler func(worker Worker)) {
  log.Print("workerPool.Start")
  for workerPool.currentTransition.state != stateExit {
    switch workerPool.currentTransition.state {
    case stateInitial:
      workerPool.processInitialState()
    case stateMain:
      workerPool.processMainState()
    case stateWait:
      workerPool.processWaitState()
    case stateSleep:
      workerPool.processSleepState()
    case stateLaunch:
      workerPool.processLaunchState(actionHandler)
    case stateQuit:
      workerPool.processQuitState()
    case stateTimeout:
      workerPool.processTimeoutState()
    case stateFinish:
      workerPool.processFinishState()
    case stateProcessing:
      log.Print("workerPool.Start calling processing state")
      workerPool.processProcessingState()
    default:
      panic("invalid state transition")
    }
  }

  workerPool.shutdown.DoneChannel <- true
}

func (workerPool *WorkerPool) processInitialState() {
  log.Print("workerPool.processInitialState")
  workerPool.goToState(stateLaunch, nil)
}

func (workerPool *WorkerPool) processMainState() {
  log.Print("workerPool.processMainState")
  // for _, worker := range workerPool.activeWorkers {
  //
  // }
}

func (workerPool *WorkerPool) processWaitState() {
  log.Print("workerPool.processWaitState")
  select {
  case finishedWorkerID := <-workerPool.finishedChannel:
    log.Print("workerPool finished channel entry")
    payload := Payload{
      workerID: finishedWorkerID,
    }
    workerPool.goToState(stateFinish, &payload)
  case processing := <-workerPool.processingChannel:
    log.Print("workerPool processing channel entry")
    payload := Payload{
      isProcessing: processing,
    }
    workerPool.goToState(stateProcessing, &payload)
  case <- workerPool.shutdown.InitiateChannel:
    log.Print("workerPool shutdown initiate channel entry")
    workerPool.goToState(stateQuit, nil)
  case <- time.After(time.Second * time.Duration(workerPool.timeout)):
    log.Print("workerPool timeout channel entry")
    workerPool.goToState(stateTimeout, nil)
  // default:
  //   log.Print("workerPool defaulting to main state from wait state")
  //   workerPool.goToState(stateMain, nil)
  }
}

func (workerPool *WorkerPool) processSleepState() {}

// in charge of creating a new worker and spawning a new goroutine with the actionHandler function
func (workerPool *WorkerPool) processLaunchState(actionHandler func(worker Worker)) {
  log.Print("workerPool.processLaunchState")
  newWorker := NewWorker(workerPool.processingChannel, workerPool.finishedChannel)

  workerPool.activeWorkers[newWorker.WorkerID] = newWorker
  go actionHandler(newWorker)

  workerPool.goToState(stateWait, nil)
}

func (workerPool *WorkerPool) processQuitState() {
  log.Print("workerPool.processQuitState")
}

func (workerPool *WorkerPool) processTimeoutState() {
  log.Print("workerPool.processTimeoutState")
}

func (workerPool *WorkerPool) processFinishState() {
  log.Print("workerPool.processFinishState")
}

func (workerPool *WorkerPool) processProcessingState() {
  log.Print("workerPool.processProcessingState")
  // workerPool.goToState(stateWait, nil)
  return
}

func (workerPool *WorkerPool) goToState(state int, payload *Payload) {
  transition := Transition{
    state: state,
  }

  if payload != nil {
    transition.payload = *payload
  }

  workerPool.currentTransition = transition
}
