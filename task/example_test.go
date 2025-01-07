package task_test

import (
	"fmt"
	"github.com/christinavaneyssen/cube/task"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"time"
)

// Example_newTask demonstrates creating a basic task with minimal configuration
func Example_newTask() {
	task := task.Task{
		ID:     uuid.New(),
		Name:   "nginx-server",
		State:  task.Pending,
		Image:  "nginx:latest",
		Memory: 512,  // 512MB
		Disk:   1024, // 1GB
	}

	fmt.Printf("Created task: %s with state: %v\n\n", task.Name, task.State)
	// Output: Created task: nginx-server with state: 0
}

// Example_configureTaskPorts shows how to configure port mappings for a task
func Example_configureTaskPorts() {
	//ports, _ := nat.ParsePortSpec("80:8080")
	portMap := nat.PortMap{
		"80/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "8080"}},
	}

	task := task.Task{
		ID:           uuid.New(),
		Name:         "web-server",
		ExposedPorts: portMap,
		PortBindings: map[string]string{"80/tcp": "8080"},
	}

	fmt.Printf("Port 80 is mapped to host port: %s\n", task.PortBindings["80/tcp"])
	// Output: Port 80 is mapped to host port: 8080
}

// Example_taskEventLifecycle demonstrates tracking task state changes
func Example_taskEventLifecycle() {
	t := task.Task{
		ID:    uuid.New(),
		Name:  "background-job",
		State: task.Pending,
	}

	event := task.TaskEvent{
		ID:        uuid.New(),
		State:     task.Scheduled,
		Timestamp: time.Now(),
		Task:      t,
	}

	fmt.Printf("Task %s transitioned to state: %v\n", event.Task.Name, event.State)
	// Output: Task background-job transitioned to state: 1
}

// Example_fullTaskConfig shows how to create a complete task configuration
func Example_fullTaskConfig() {
	config := task.Config{
		Name:         "redis-cache",
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		ExposedPorts: nat.PortMap{
			"6379/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "6379"}},
		},
		Cmd:           []string{"redis-server"},
		Image:         "redis:latest",
		Cpu:           1.0,
		Memory:        1024 * 1024 * 1024,      // 1GB
		Disk:          10 * 1024 * 1024 * 1024, // 10GB
		Env:           []string{"REDIS_PASSWORD=secret"},
		RestartPolicy: "always",
	}

	fmt.Printf("Created Redis configuration with memory: %d bytes\n", config.Memory)
	// Output: Created Redis configuration with memory: 1073741824 bytes
}
