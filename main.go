package main

import (
	"fmt"
	"github.com/christinavaneyssen/cube/manager"
	"github.com/christinavaneyssen/cube/node"
	"github.com/christinavaneyssen/cube/task"
	"github.com/christinavaneyssen/cube/worker"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"time"
)

func main() {
	t := task.Task{
		ID:     uuid.New(),
		Name:   "first-task",
		State:  task.Pending,
		Image:  "Image-1",
		Memory: 1024,
		Disk:   1,
	}

	te := task.TaskEvent{
		ID:        uuid.New(),
		State:     task.Pending,
		Timestamp: time.Now(),
		Task:      t,
	}

	fmt.Printf("task: %v\n", t)
	fmt.Printf("task event: %v\n", te)

	w := worker.Worker{
		Name:  "first-worker",
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}
	fmt.Printf("worker: %v\n", w)
	w.CollectStats()
	w.RunTask()
	w.StartTask()
	w.StopTask()

	m := manager.Manager{
		Pending:       *queue.New(),
		TaskDb:        map[string][]*task.Task{},
		EventDb:       map[string][]*task.TaskEvent{},
		Workers:       []string{w.Name},
		WorkerTaskMap: nil,
		TaskWorkerMap: nil,
	}
	fmt.Printf("manager: %v\n", m)
	m.SelectWorker()
	m.UpdateTasks()
	m.SendTask()

	n := node.Node{
		Name:   "first-node",
		Ip:     "192.168.1.1",
		Cores:  4,
		Memory: 1024,
		Disk:   25,
		Role:   "worker",
	}
	fmt.Printf("node: %v\n", n)
}
