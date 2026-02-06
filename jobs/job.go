package jobs

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type Job struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Payload   string    `json:"payload"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type JobStore struct {
	mu   sync.RWMutex
	jobs map[string]*Job
}

func NewJobStore() *JobStore {
	return &JobStore{
		jobs: make(map[string]*Job),
	}
}

func (s *JobStore) AddJob(jobtype, payload string) *Job {
	s.mu.Lock()
	defer s.mu.Unlock()

	job := &Job{
		ID:        uuid.NewString(),
		Type:      jobtype,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	s.jobs[job.ID] = job
	return job
}

func (s *JobStore) GetJob(id string) (*Job, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	job, ok := s.jobs[id]
	return job, ok
}

func (s *JobStore) UpdateJobStatus(id, status string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if job, ok := s.jobs[id]; ok {

		job.Status = status
		job.UpdatedAt = time.Now()
	}

}
