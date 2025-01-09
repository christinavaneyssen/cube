package manager_test

import (
	"github.com/christinavaneyssen/cube/manager"
	"github.com/christinavaneyssen/cube/task"
	"github.com/google/uuid"
	"testing"
	"time"
)

// ExampleManager_BasicUsage demonstrates the basic workflow for creating and using a Manager.
func ExampleManager_BasicUsage() {
	mgr := &manager.Manager{
		TaskDb:        make(map[string][]*task.Task),
		EventDb:       make(map[string][]*task.TaskEvent),
		Workers:       []string{"worker1", "worker2", "worker3"},
		WorkerTaskMap: make(map[string][]uuid.UUID),
		TaskWorkerMap: make(map[uuid.UUID]string),
	}

	newTask := &task.Task{
		ID:    uuid.New(),
		Name:  "data-processing",
		State: task.Pending,
	}

	mgr.Pending.Enqueue(newTask)

	// Process the pending task
	mgr.SelectWorker()
	mgr.SendWork()
	mgr.UpdateTasks()

	// Output:
	// I will select an appropriate worker
	// I send the task to the worker
	// I keep track of tasks, their states and the machines they run on
}

// TestManager_TaskDistribution tests the distribution of tasks across workers.
func TestManager_TaskDistribution(t *testing.T) {
	mgr := &manager.Manager{
		TaskDb:        make(map[string][]*task.Task),
		EventDb:       make(map[string][]*task.TaskEvent),
		Workers:       []string{"worker1", "worker2", "worker3"},
		WorkerTaskMap: make(map[string][]uuid.UUID),
		TaskWorkerMap: make(map[uuid.UUID]string),
	}

	tasks := []*task.Task{
		{
			ID:    uuid.New(),
			Name:  "task-1",
			State: task.Pending,
		},
		{
			ID:    uuid.New(),
			Name:  "task-2",
			State: task.Pending,
		},
		{
			ID:    uuid.New(),
			Name:  "task-3",
			State: task.Pending,
		},
	}

	for _, task := range tasks {
		mgr.Pending.Enqueue(task)
	}

	// Process all pending tasks
	//for mgr.Pending.Len() > 0 {
	//	mgr.SelectWorker()
	//	mgr.SendWork()
	//	mgr.UpdateTasks()
	//}
}

func TestManager_TaskStateTransitions(t *testing.T) {
	mgr := &manager.Manager{
		TaskDb:        make(map[string][]*task.Task),
		EventDb:       make(map[string][]*task.TaskEvent),
		Workers:       []string{"worker1", "worker2", "worker3"},
		WorkerTaskMap: make(map[string][]uuid.UUID),
		TaskWorkerMap: make(map[uuid.UUID]string),
	}

	taskID := uuid.New()
	newTask := task.Task{
		ID:    taskID,
		Name:  "state-transition-task",
		State: task.Pending,
	}

	mgr.Pending.Enqueue(newTask)
	mgr.SelectWorker()
	mgr.SendWork()

	mgr.EventDb[taskID.String()] = []*task.TaskEvent{
		{
			ID:        uuid.New(),
			State:     task.Pending,
			Timestamp: time.Now(),
			Task:      newTask,
		},
		{
			ID:        uuid.New(),
			State:     task.Scheduled,
			Timestamp: time.Now(),
			Task:      newTask,
		},
		{
			ID:        uuid.New(),
			State:     task.Running,
			Timestamp: time.Now(),
			Task:      newTask,
		},
		{
			ID:        uuid.New(),
			State:     task.Completed,
			Timestamp: time.Now(),
			Task:      newTask,
		},
	}
	mgr.UpdateTasks()
}
