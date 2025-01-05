package manager

import (
	"cube/task"
	"fmt"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

type Manager struct {
	Pending       queue.Queue
	TaskDb        map[string][]*task.Task
	EventDb       map[string][]*task.TaskEvent
	Workers       []string
	WorkerTaskMap map[string][]uuid.UUID // k = name of worker, v = slice of task UUIDs
	TaskWorkerMap map[uuid.UUID]string   // k = task UUID, v = name of worker
}

func (m *Manager) SelectWorker() {
	fmt.Println("I will select an appropriate worker")
}

func (m *Manager) UpdateTasks() {
	fmt.Println("I keep track of tasks, their states and the machines they run on")
}

func (m *Manager) SendTask() {
	fmt.Println("I send the task to the worker")
}
