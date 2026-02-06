package workers

import (
	"log"
	"time"

	"github.com/m-arshad-dev/job-queue/jobs"
)

type WorkerPool struct {
	JobQueue chan *jobs.Job
	Store    *jobs.JobStore
}

func NewWorkerPool(store *jobs.JobStore, queueSize int) *WorkerPool {

	return &WorkerPool{
		JobQueue: make(chan *jobs.Job, queueSize),
		Store:    store,
	}
}

func (wp *WorkerPool) Start(numWoekers int) {
	for i := 1; i < numWoekers+1; i++ {
		go wp.worker(i)
	}
}
func (wp *WorkerPool) worker(id int) {
	log.Printf("Worker %d started\n", id)

	for job := range wp.JobQueue {
		wp.process(job)
	}
}

func (wp *WorkerPool) process(job *jobs.Job) {

	wp.Store.UpdateJobStatus(job.ID, "porcessing")

	time.Sleep(2 * time.Second)

	wp.Store.UpdateJobStatus(job.ID, "done")

}
