// Package task provides types and constants for managing containerized workloads
// in an orchestration system. It defines the core abstraction for
// tasks, their lifecycle states, events, and configuration.
package task

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"io"
	"log"
	"math"
	"os"
	"time"
)

// State represents the current lifecycle stage of a task in the orchestration
type State int

const (
	// Pending indicates the task is queued and awaiting scheduler placement
	// iota automatically generates sequential constant values starting from 0
	Pending State = iota

	// Scheduled indicates the task has been assigned to a specific worker by scheduler
	Scheduled

	// Running indicates the task is actively executing on worker
	Running

	// Completed indicates the task has finished successfully or was gracefully terminated
	Completed

	// Failed indicates the task terminated abnormally due to error or crash
	Failed
)

// Task represents a containerized workload with its configuration and runtime state.
// It encapsulates all necessary information to schedule, run, and monitor a task and container.
type Task struct {
	// Id uniquely identifies the task
	ID uuid.UUID

	// Name is a human-readable identifier for the task
	Name string

	// State indicates the current lifecycle stage of the task
	State State

	// Image specifies the container image to be used
	Image string

	// Memory specifies the amount of memory in MB to allocate to the container
	Memory int

	// Disk specifies the amount of disk space in MB to allocate to the container
	Disk int

	// ExposedPorts defines which ports are exposed by the container
	ExposedPorts nat.PortMap

	// PortBindings maps container ports to host ports
	// Format: "containerPort/protocol": "hostPort"
	PortBindings map[string]string

	// RestartPolicy defines how the container should be restarted on exit
	// Valid values are:
	// 	- "" (empty string): no restart
	// 	- "always": restart the container any time it stops
	// - "unless-stopped": restart the container unless explicitly stopped
	// 	- "on-failure": restart the container only on non-zero exit code
	RestartPolicy string

	// StartTime records when the task began execution
	StartTime time.Time

	// FinishTime records when the task completed execution
	FinishTime time.Time
}

// TaskEvent represents a point-in-time state change of a task in the orchestration.
// It captures the transition details including when it occurred and the task's full state.
type TaskEvent struct {
	// ID uniquely identifies this event
	ID uuid.UUID

	// State indicates the new state that the task transitioned to
	State State

	// Timestamp records when this state transition occurred
	Timestamp time.Time

	// Task contains the complete task information at the time of the event
	Task Task
}

// Config defines the configuration parameters for an orchestration task.
type Config struct {
	// Name specifies both the task name and container name
	Name string

	// AttachStdin indicates whether to attach to the container's standard input
	AttachStdin bool

	// AttachStdout indicates whether to attach to the container's standard output
	AttachStdout bool

	// AttachStderr indicates whether to attach to the container's standard error
	AttachStderr bool

	// ExposedPorts defines the network ports to expose from the container
	ExposedPorts nat.PortSet

	// Cmd specifies the command to run in the container
	Cmd []string

	// Image represents the name of the container image to run
	Image string

	// Cpu defines the amount of CPU resources to allocate to the container in CPU shares
	Cpu float64

	// Memory specifies the memory limit in bytes for the container
	// The scheduler uses this value to find a suitable node in the cluster
	Memory int64

	// Disk specifies the disk space limit in bytes for the container
	// The scheduler uses this value to find a suitable node in the cluster
	Disk int64

	// Env specifies environment variables to pass to the container
	Env []string

	// RestartPolicy defines the container's restart behaviour on exit
	RestartPolicy container.RestartPolicyMode
}

// Docker provides an interface to interact with the Docker daemon through the Docker API.
type Docker struct {
	// Client is the Docker client used to communicate with the Docker daemon
	Client *client.Client

	// Config holds both the initial task configuration and runtime information
	// such as ContainerID once the task is running
	Config Config

	ContainerID string
}

// DockerResult encapsulates the outcome of Docker operations
// such as starting or stopping containers.
type DockerResult struct {
	// Error holds any error that occurred during the operation
	Error error

	// Action describes the operation performed (eg. "start" or "stop")
	Action string

	// ContainerID uniquely identifies the target container
	ContainerID string

	// Result contains additional operation-specific output
	Result string
}

func (d *Docker) Run() DockerResult {
	log.Printf("Attempting to start container")
	ctx := context.Background()

	reader, err := d.Client.ImagePull(
		ctx, d.Config.Image, image.PullOptions{})
	if err != nil {
		log.Printf("Failed to pull image %s: %v\n", d.Config.Image, err)
		return DockerResult{Error: err}
	}
	io.Copy(os.Stdout, reader)
	rp := container.RestartPolicy{
		Name: d.Config.RestartPolicy,
	}

	r := container.Resources{
		Memory:   d.Config.Memory,
		NanoCPUs: int64(d.Config.Cpu * math.Pow(10, 9)),
	}

	cc := container.Config{
		Image:        d.Config.Image,
		Tty:          false,
		Env:          d.Config.Env,
		ExposedPorts: d.Config.ExposedPorts,
	}

	hc := container.HostConfig{
		RestartPolicy:   rp,
		Resources:       r,
		PublishAllPorts: true,
	}

	resp, err := d.Client.ContainerCreate(ctx, &cc, &hc, nil, nil, d.Config.Name)
	if err != nil {
		log.Printf("Error creating container using image: %s: %v\n", d.Config.Image, err)
		return DockerResult{Error: err}
	}

	err = d.Client.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		log.Printf("Error starting container: %s: %v\n", resp.ID, err)
		return DockerResult{Error: err}
	}

	d.ContainerID = resp.ID

	out, err := d.Client.ContainerLogs(
		ctx,
		resp.ID,
		container.LogsOptions{ShowStdout: true, ShowStderr: true},
	)
	if err != nil {
		log.Printf("Error getting logs for container %s: %v\n", resp.ID, err)
		return DockerResult{Error: err}
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	return DockerResult{
		Error:       nil,
		Action:      "start",
		ContainerID: resp.ID,
		Result:      "success",
	}
}
