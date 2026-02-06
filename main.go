package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/m-arshad-dev/job-queue/jobs"
	"github.com/m-arshad-dev/job-queue/workers"
)

func main() {
	store := jobs.NewJobStore()
	workerPool := workers.NewWorkerPool(store, 100)
	workerPool.Start(3)

	// POST /jobs
	http.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Type    string `json:"type"`
			Payload string `json:"payload"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		job := store.AddJob(req.Type, req.Payload)
		workerPool.JobQueue <- job

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(job)
	})

	// GET /jobs/{id}
	http.HandleFunc("/jobs/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/jobs/")

		job, ok := store.GetJob(id)
		if !ok {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(job)
	})

	log.Println("Job queue running on :8080")
	http.ListenAndServe(":8080", nil)
}
