package queue

// Package queue provides a simple in-memory task queue with worker pool.
// It supports starting a fixed number of workers that consume tasks from a
// bounded channel, enqueuing new tasks, and gracefully stopping the queue.
//
// Example:
//
//	tq := queue.NewQueue(100)
//	tq.Start(5)
//	tq.Enqueue(func(ctx context.Context) error {
//		fmt.Println("Processing task")
//		return nil
//	})
//	tq.Stop()

import "context"

// Task represents a unit of work to be executed by the queue.
type Task func(context.Context) error

// Queue represents a simple in-memory task queue with worker pool.
type Queue struct {
	tasks chan Task
}

// NewQueue creates a new queue with the given size.
func NewQueue(size int) *Queue {
	return &Queue{tasks: make(chan Task, size)}
}

// Start starts the worker pool.
func (q *Queue) Start(workers int) {
	for i := 0; i < workers; i++ {
		go func() {
			for t := range q.tasks {
				t(context.Background())
			}
		}()
	}
}

// Enqueue adds a new task to the queue.
func (q *Queue) Enqueue(t Task) {
	q.tasks <- t
}

// Stop closes the queue.
func (q *Queue) Stop() {
	close(q.tasks)
}

// Size returns the number of tasks in the queue.
func (q *Queue) Size() int {
	return len(q.tasks)
}
