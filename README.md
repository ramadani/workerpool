# workerpool
The workerpool package provides a simple and efficient worker pool implementation in Go, enabling concurrent processing of tasks with a fixed number of workers. This package is ideal for handling workloads where you need to manage and balance multiple goroutines performing tasks in parallel.

## Features
- **Concurrency Management**: Easily manage a pool of workers to process tasks concurrently.
- **Task Scheduling**: Submit tasks to the worker pool and have them automatically distributed to available workers.
- **Callback Support**: Define and execute callback functions for each worker before they start processing a task.
- **Result Collection**: Collect results from each task processed by the workers.

## Installation
To install the package, run:

```bash
go get github.com/ramadani/workerpool
```

## Usage
Here is an example of how to use the worker pool:

```go
package main

import (
    "fmt"
    "time"
    "github.com/ramadani/workerpool"
)

func main() {
	const numJobs = 10
	const numWorkers = 2

	callback := func(workerID int, job workerpool.Job) workerpool.Result {
		fmt.Printf("Worker %d processing job %d with data: %v\n", workerID, job.ID, job.Data)
		// Simulate processing time
		time.Sleep(100 * time.Millisecond)
		return workerpool.Result{JobID: job.ID, Data: job.Data, Error: nil}
	}

	wp := workerpool.New(numWorkers, numJobs, callback)

	// Start the worker pool
	wp.Start()

	// Send jobs to the jobs channel
	for j := 1; j <= numJobs; j++ {
		wp.AddJob(workerpool.Job{ID: j, Data: fmt.Sprintf("data %d", j)})
	}

	// Wait for all workers to finish
	wp.Wait()

	// Collect results (if needed)
	for result := range wp.Results() {
		if result.Error != nil {
			fmt.Printf("Job %d encountered error: %v\n", result.JobID, result.Error)
		} else {
			fmt.Printf("Job %d processed with result: %v\n", result.JobID, result.Data)
		}
	}
}
```

## Worker Pool API

### New
Creates a new worker pool.

```go
func New(numWorkers, numJobs int, callbackFunc CallbackFunc) *WorkerPool
```

- `numWorkers`: Number of workers in the pool.
- `numJobs`: Capacity of the job queue.
- `callbackFunc`: Function to be executed by each worker for each job.

### Start
Starts the worker pool.

```go
func (wp *WorkerPool) Start()
```

### AddJob
Adds a job to the worker pool.

```go
func (wp *WorkerPool) AddJob(job Job)
```

- `job`: Job to be added to the pool.

### Wait
Waits for all workers to finish processing jobs.

```go
func (wp *WorkerPool) Wait()
```

### Results
Returns a channel to collect results.

```go
func (wp *WorkerPool) Results() <-chan Result
```

## License
This project is licensed under the MIT License - see the `LICENSE` file for details.