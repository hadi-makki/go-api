package test

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Job represents a unit of work
type Job struct {
	ID      int
	Payload string
}

var jobQueue = make(chan Job, 10) // Buffered channel to hold incoming jobs
// var exit = make(chan struct{})

// Worker function to process jobs
func worker(id int, jobs <-chan Job) {
	for job := range jobs {
		fmt.Printf("Worker %d processing job %d with payload: %s\n", id, job.ID, job.Payload)

		// Simulate job processing time and handle potential errors
		if err := processJob(job); err != nil {
			log.Printf("Worker %d failed to process job %d: %v\n", id, job.ID, err)
			// Depending on the requirements, you might want to retry the job, move it to a dead-letter queue, etc.
			continue
		}

		fmt.Printf("Worker %d completed job %d\n", id, job.ID)
	}
}

// processJob simulates job processing and returns an error if something goes wrong
func processJob(job Job) error {
	// Simulate job processing time
	time.Sleep(2 * time.Second)

	// Simulate an error for demonstration
	if job.Payload == "error" {
		return fmt.Errorf("simulated error processing job %d", job.ID)
	}

	return nil
}

// HTTP handler to accept job requests
func jobHandler(w http.ResponseWriter, r *http.Request) {
	payload := r.URL.Query().Get("payload")
	if payload == "" {
		http.Error(w, "Missing payload", http.StatusBadRequest)
		return
	}

	jobID := time.Now().UnixNano() // Generate a unique job ID
	job := Job{ID: int(jobID), Payload: payload}

	select {
	case jobQueue <- job:
		fmt.Fprintf(w, "Job %d enqueued\n", job.ID)
	default:
		http.Error(w, "Job queue is full", http.StatusServiceUnavailable)
	}
}

func PerformBackgroundJobs() {
	// Start worker goroutines
	numWorkers := 3
	for i := 1; i <= numWorkers; i++ {
		go worker(i, jobQueue)
	}

	// Start the HTTP server
	http.HandleFunc("/enqueue", jobHandler)
	server := &http.Server{Addr: ":8080"}

	// Start the server in a goroutine
	go func() {
		fmt.Println("Starting server on port 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("ListenAndServe(): %v\n", err)
		}
	}()

	// Handle graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	// Graceful shutdown of the server
	fmt.Println("Shutting down server...")
	if err := server.Close(); err != nil {
		log.Printf("Error closing server: %v", err)
	}
	close(jobQueue) // Close job queue to signal workers to stop

	// Wait for a short duration to allow workers to complete
	time.Sleep(3 * time.Second)
	fmt.Println("Server exited gracefully")
}
