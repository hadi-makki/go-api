package main

import (
	"fmt"
	"go-api/test"
	"log"
	"net/http"
	"time"
)

// Job represents a unit of work
type Job struct {
	ID      int
	Payload string
}

// PerformJob simulates job processing and returns an error if something goes wrong
func PerformJob(job Job) error {
	// Simulate job processing time
	time.Sleep(2 * time.Second)

	// Simulate an error for demonstration
	if job.Payload == "error" {
		return fmt.Errorf("simulated error processing job %d", job.ID)
	}

	return nil
}

// HTTP handler to accept job requests and wait for job completion
func jobHandler(w http.ResponseWriter, r *http.Request) {
	payload := r.URL.Query().Get("payload")
	if payload == "" {
		http.Error(w, "Missing payload", http.StatusBadRequest)
		return
	}

	jobID := time.Now().UnixNano() // Generate a unique job ID
	job := Job{ID: int(jobID), Payload: payload}

	// Channel to signal job completion
	done := make(chan error)

	// Start a goroutine to perform the job
	go func() {
		err := PerformJob(job)
		done <- err
	}()

	// Wait for the job to complete
	err := <-done
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Job %d completed successfully\n", job.ID)
}

func main() {
	// http.HandleFunc("/enqueue", jobHandler)
	test.SyncWaitGroupUsageExample()
	server := &http.Server{Addr: ":8080"}

	// Start the server
	fmt.Println("Starting server on port 8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %v\n", err)
	}
}
