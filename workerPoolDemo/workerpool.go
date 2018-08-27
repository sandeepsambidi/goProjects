package main

import (
  "fmt"
  "math/rand"
  "sync"
  "time"
)

type Job struct {
  id int
  randomNum int
}

type Result struct {
  job Job
  res int
}

//channel to store input
var inputStore = make(chan Job, 10)
//channel to store result after processing each input
var results = make(chan Result, 10)

//work to be done on the input
func work(num int) int{
  sum := 0
  for num != 0 {
    digit := num%10
    sum += digit
    num = num/10;
  }
  time.Sleep(2*time.Second)
  return sum
}

//worker will results of the work and puts the output in results channel
func worker(wg *sync.WaitGroup) {
  for input := range inputStore {
    output := Result{input, work(input.randomNum)}
    results <- output
  }
  //If no jobs are availabe in the inputStore, worker will mark itself done
  wg.Done()
}

//Create the worker pool with the number of workers required
func createWorkerPool(numOfWorkers int) {
  var wg sync.WaitGroup
  for i:=0; i < numOfWorkers; i++ {
    wg.Add(1)
    //worker will need to be created in a seperate thread,
    //else it will be blocking worker creation
    go worker(&wg)
  }
  wg.Wait()
  //unless we close results "for range" on results resultswill never end
  close(results)
}

//simulation function to add inputs to the inputStore channel
func addInputsToStore(numOfInputs int) {
  for i:=0; i<numOfInputs; i++ {
    randomNum := rand.Intn(10000)
    input := Job{i, randomNum}
    inputStore <- input
  }
  //unless we close jobs "for range" on inputStore channel will never end
  close(inputStore)
}

//function to check results in the results channel
func readResult(done chan bool) {
  for result := range results {
        fmt.Printf("Job id %d, input random no %d , sum of digits %d\n",
          result.job.id, result.job.randomNum, result.res)
    }
    done <- true
}

//test func to trigger
func main() {
  startTime := time.Now()
  noOfInputs := 100
  go addInputsToStore(noOfInputs)
  done := make(chan bool)
  go readResult(done)
  noOfWorkers := 100
  createWorkerPool(noOfWorkers)
  <-done
  endTime := time.Now()
  diff := endTime.Sub(startTime)
  fmt.Println("total time taken ", diff.Seconds(), "seconds")
}
