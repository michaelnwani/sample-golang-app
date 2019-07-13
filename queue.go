package main

// import (
//   "fmt"
//   "sync"
//   "github.com/aws/aws-sdk-go/aws"
//   "github.com/aws/aws-sdk-go/aws/awserr"
//   "github.com/aws/aws-sdk-go/aws/session"
//   "github.com/aws/aws-sdk-go/service/sqs"
//   "reflect"
// )
//
// const (
//   region = "us-east-1"
//   queueName = "sqs_create_q.fifo"
// )
//
// var (
//   // A buffered channel that we can send work requests on
//   // has a capacity of 100
//   JobQueue = make(chan Job, 100)
//   wg sync.WaitGroup
// )
//
// type Job struct {
//   Type  string
//
// }
//
// func Init() {
//   // make a request with a capacity of 100
//   // jobChan := make(chan Job, 100)
//
//   // increment the WaitGroup before starting the 'worker' (goroutine)
//   wg.Add(1)
//
//   // start the worker
//   go worker(jobChan)
// }
//
// func worker(jobChan <-chan Job) {
//   defer wg.Done()
//
//   for job := range jobChan {
//     process(job)
//   }
// }
//
// func process(job Job) {
//   fmt.Println("job ", job)
// }
//
// func fetchMessages() {
//   sess, err := session.NewSession(&aws.Config{
//     Region: aws.String(region)},
//   )
//
//   // Create a SQS service client.
//   svc := sqs.New(sess)
//
//   resultUrl, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
//     QueueName: aws.String(queueName),
//   })
//
//   if err != nil {
//     if aerr, ok := err.(awserr.Error); ok && aerr.Code() == sqs.ErrCodeQueueDoesNotExist {
//       exitErrorf("Unable to find queue %q.", name)
//     }
//     exitErrorf("Unable to find queue %q, %v", name, err)
//   }
// }
//
// // TryEnqueue tries to enqueue a job to the given job channel. Returns true if
// // the operation was successful, and false if enqueuing would not have been
// // possible without blocking. Job is not enqueued in the latter case.
// func TryEnqueue(job Job, jobChan <-chan Job) bool {
//   select {
//   case jobChan <- job:
//     return true
//   default:
//     return false
//   }
// }
//
// // WaitTimeout does a Wait on a sync.WaitGroup object but with a specified
// // timeout. Returns true if the wait completed without timing out, false otherwise.
// func WaitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
//   ch := make(chan struct{})
//   go func() {
//     wg.Wait()
//     close(ch)
//   }()
//   select {
//   case <-ch:
//     return true
//   case <- time.After(timeout):
//     return false
//   }
// }
