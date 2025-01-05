package task

import (
	"github.com/docker/go-connections/nat"
	_ "github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"time"
)

type State int

const (
	Pending   State = iota // Enqueued, waiting for the Scheduler to work out the best place for the task to be run
	Scheduled              // Scheduler has worked out where to run the task
	Running                // Worker successfully starts the task
	Completed              // Task completed successfully OR stopped by user
	Failed                 // Task has crashed or has stopped working as expected
)

type Task struct {
	ID            uuid.UUID
	Name          string
	State         State
	Image         string
	Memory        int
	Disk          int
	ExposedPorts  nat.PortMap
	PortBindings  map[string]string
	RestartPolicy string
	StartTime     time.Time
	FinishTime    time.Time
}

type TaskEvent struct {
	ID        uuid.UUID
	State     State
	Timestamp time.Time
	Task      Task
}
