package test

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Task represents a unit of work
type Task struct {
	ID      int
	Payload string
}

// PerformTask simulates task processing and returns an error if something goes wrong
func PerformTask1(task Task) error {
	// Simulate task processing time
	time.Sleep(2 * time.Second)
	fmt.Println("task 1 done")
	// Simulate an error for demonstration
	if task.Payload == "error" {
		return fmt.Errorf("simulated error processing task %d", task.ID)
	}

	return nil
}
func PerformTask2(task Task) error {
	// Simulate task processing time
	time.Sleep(1 * time.Second)
	fmt.Println("task 2 done")

	// Simulate an error for demonstration
	if task.Payload == "error" {
		return fmt.Errorf("simulated error processing task %d", task.ID)
	}

	return nil
}

// HTTP handler to accept task requests and wait for both tasks to complete
func taskHandler(w http.ResponseWriter, r *http.Request) {
	payload1 := r.URL.Query().Get("payload1")
	payload2 := r.URL.Query().Get("payload2")
	if payload1 == "" || payload2 == "" {
		http.Error(w, "Missing payload(s)", http.StatusBadRequest)
		return
	}

	task1 := Task{ID: 1, Payload: payload1}
	task2 := Task{ID: 2, Payload: payload2}

	var wg sync.WaitGroup
	var err1, err2 error

	wg.Add(2) // We have two tasks to wait for

	go func() {
		defer wg.Done()
		err1 = PerformTask1(task1)
	}()

	go func() {
		defer wg.Done()
		err2 = PerformTask2(task2)
	}()

	// Wait for both tasks to complete
	wg.Wait()

	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
		return
	}
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Both tasks completed successfully\n")
}
func SyncWaitGroupUsageExample() {
	http.HandleFunc("/enqueue", taskHandler)

}
