// Package manager implements a task management system that handles task distribution
// and tracking across multiple workers.
package manager

import (
	"fmt"
	"github.com/christinavaneyssen/cube/task"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

type Manager struct {
	// Pending contains tasks that are waiting to be assigned to workers
	Pending queue.Queue

	// TaskDb stores tasks indexed by a string key
	TaskDb map[string][]*task.Task

	// TaskEvent stores task events indexed by a string key
	EventDb map[string][]*task.TaskEvent

	// Workers contains a list of available workers
	Workers []string

	// WorkerTaskMap maintains the relationship between workers and their assigned tasks
	// Key: worker name, Value: slice of task UUIDs assigned to the worker
	WorkerTaskMap map[string][]uuid.UUID

	// TaskWorkerMap maintains the reverse mapping of tasks to workers
	// Key: task UUID, Value: name of the worker the task is assigned to
	TaskWorkerMap map[uuid.UUID]string // k = task UUID, v = name of worker
}

// SelectWorker chooses an appropriate worker from the available pool
// based on current workload and capacity
func (m *Manager) SelectWorker() {
	fmt.Println("I will select an appropriate worker")
}

// UpdateTasks maintains the current state of all tasks in the system.
func (m *Manager) UpdateTasks() {
	fmt.Println("I keep track of tasks, their states and the machines they run on")
}

// SendWork dispatches a task to its assigned worker for execution.
func (m *Manager) SendWork() {
	fmt.Println("I send the task to the worker")
}
