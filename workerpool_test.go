package workerpool_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ramadani/workerpool"
	"github.com/stretchr/testify/assert"
)

func TestWorkerPool(t *testing.T) {
	const numJobs = 5
	const numWorkers = 2

	callback := func(workerID int, job workerpool.Job) workerpool.Result {
		resultData := fmt.Sprintf("data %d", job.ID)
		return workerpool.Result{JobID: job.ID, Data: resultData, Error: nil}
	}

	wp := workerpool.New(numWorkers, numJobs, callback)
	wp.Start()

	expectedResults := make(map[int]string)
	for j := 1; j <= numJobs; j++ {
		expectedResults[j] = fmt.Sprintf("data %d", j)
		wp.AddJob(workerpool.Job{ID: j, Data: fmt.Sprintf("data %d", j)})
	}

	wp.Wait()

	results := make(map[int]workerpool.Result)
	for result := range wp.Results() {
		results[result.JobID] = result
	}

	assert.Equal(t, numJobs, len(results))

	for jobID, result := range results {
		expectedResult := expectedResults[jobID]
		actualResult := result.Data.(string)
		assert.Equal(t, expectedResult, actualResult)
		assert.Nil(t, result.Error)
	}
}

func TestWorkerPoolWithError(t *testing.T) {
	const numJobs = 5
	const numWorkers = 2

	callback := func(workerID int, job workerpool.Job) workerpool.Result {
		if job.ID%2 == 0 {
			return workerpool.Result{JobID: job.ID, Data: nil, Error: fmt.Errorf("error processing job %d", job.ID)}
		}
		resultData := fmt.Sprintf("data %d", job.ID)
		return workerpool.Result{JobID: job.ID, Data: resultData, Error: nil}
	}

	wp := workerpool.New(numWorkers, numJobs, callback)
	wp.Start()

	for j := 1; j <= numJobs; j++ {
		wp.AddJob(workerpool.Job{ID: j, Data: fmt.Sprintf("data %d", j)})
	}

	wp.Wait()

	results := make(map[int]workerpool.Result)
	for result := range wp.Results() {
		results[result.JobID] = result
	}

	assert.Equal(t, numJobs, len(results))

	for jobID, result := range results {
		if jobID%2 == 0 {
			assert.NotNil(t, result.Error)
		} else {
			assert.Nil(t, result.Error)
			expectedResult := fmt.Sprintf("data %d", jobID)
			actualResult := result.Data.(string)
			assert.Equal(t, expectedResult, actualResult)
			assert.Regexp(t, expectedResult, actualResult)
		}
	}
}

func TestWorkerPoolConcurrency(t *testing.T) {
	const numJobs = 10
	const numWorkers = 2

	callback := func(workerID int, job workerpool.Job) workerpool.Result {
		time.Sleep(100 * time.Millisecond)
		resultData := fmt.Sprintf("data %d", job.ID)
		return workerpool.Result{JobID: job.ID, Data: resultData, Error: nil}
	}

	wp := workerpool.New(numWorkers, numJobs, callback)
	wp.Start()

	for j := 1; j <= numJobs; j++ {
		wp.AddJob(workerpool.Job{ID: j, Data: fmt.Sprintf("data %d", j)})
	}

	wp.Wait()

	results := make(map[int]workerpool.Result)
	for result := range wp.Results() {
		results[result.JobID] = result
	}

	assert.Equal(t, numJobs, len(results))

	for jobID, result := range results {
		expectedResult := fmt.Sprintf("data %d", jobID)
		actualResult := result.Data.(string)
		assert.Equal(t, expectedResult, actualResult)
		assert.Nil(t, result.Error)
	}
}
