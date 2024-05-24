package workerpool

import (
	"sync"
)

// CallbackFunc defines the function signature for worker callbacks
type CallbackFunc func(workerID int, job Job) Result

// Job represents the job to be run
type Job struct {
	ID   int
	Data interface{}
}

// Result represents the result of a job
type Result struct {
	JobID int
	Data  interface{}
	Error error
}

// WorkerPool struct to manage workers and jobs
type WorkerPool struct {
	numWorkers int
	jobs       chan Job
	results    chan Result
	wg         sync.WaitGroup
	callback   CallbackFunc
}

// NewWorkerPool creates a new WorkerPool
func NewWorkerPool(numWorkers, numJobs int, callbackFunc CallbackFunc) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		jobs:       make(chan Job, numJobs),
		results:    make(chan Result, numJobs),
		callback:   callbackFunc,
	}
}

// Start initializes and runs the worker pool
func (wp *WorkerPool) Start() {
	for w := 1; w <= wp.numWorkers; w++ {
		wp.wg.Add(1)
		go wp.worker(w)
	}
}

// worker function processes jobs from the jobs channel and sends results to the results channel
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	for job := range wp.jobs {
		result := wp.callback(id, job)
		wp.results <- result
	}
}

// AddJob sends a job to the jobs channel
func (wp *WorkerPool) AddJob(job Job) {
	wp.jobs <- job
}

// Wait waits for all workers to finish
func (wp *WorkerPool) Wait() {
	close(wp.jobs)
	wp.wg.Wait()
	close(wp.results)
}

// Results returns the results channel
func (wp *WorkerPool) Results() <-chan Result {
	return wp.results
}
