package main

import (
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/client"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/louishuyng/orchestrator/manager"
	"github.com/louishuyng/orchestrator/node"
	"github.com/louishuyng/orchestrator/task"
	"github.com/louishuyng/orchestrator/worker"
)

func main() {
	t := task.Task{
		ID:     uuid.New(),
		Name:   "Task-1",
		State:  task.Pending,
		Image:  "Image-1",
		Memory: 1024,
		Disk:   1,
	}

	te := task.TaskEvent{
		ID:        uuid.New(),
		State:     task.Pending,
		TimeStamp: time.Now(),
		Task:      t,
	}

	fmt.Printf("Task: %+v\n", t)
	fmt.Printf("Task Event: %+v\n", te)

	w := worker.Worker{
		Name:  "Worker-1",
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}

	fmt.Printf("Worker: %+v\n", w)
	w.CollectStats()
	w.RunTask()
	w.StartTask()
	w.StopTask()

	m := manager.Manager{
		Pending: *queue.New(),
		TaskDb:  make(map[string][]*task.Task),
		EventDb: make(map[string][]*task.TaskEvent),
		Workers: []string{w.Name},
	}

	fmt.Printf("Manager: %+v\n", m)
	m.SelectWorker()
	m.UpdateTask()
	m.SendWork()

	n := node.Node{
		Name:   "Node-1",
		Ip:     "168.192.1.1",
		Cores:  4,
		Memory: 1024,
		Disk:   25,
		Role:   "worker",
	}

	fmt.Printf("Node: %+v\n", n)

	fmt.Println("Creating Docker container...")
	dockerTask, createResult := createContainer()
	if createResult.Error != nil {
		fmt.Printf("Failed to create Docker container: %v\n", createResult.Error)
		os.Exit(1)
	}

	time.Sleep(5 * time.Second)

	fmt.Println("Stopping Docker container...")
	_ = stopContainer(dockerTask, createResult.ContainerID)

}

func createContainer() (*task.Docker, *task.DockerResult) {
	c := task.Config{
		Name:  "test-container-1",
		Image: "postgres:13",
		Env: []string{
			"POSTGRES_USER=testuser",
			"POSTGRES_PASSWORD=testpassword",
		},
	}

	dc, _ := client.NewClientWithOpts(client.FromEnv)
	d := task.Docker{
		Client: dc,
		Config: c,
	}

	result := d.Run()
	if result.Error != nil {
		fmt.Printf("Error running container: %v\n", result.Error)
		return nil, &result
	}

	fmt.Printf("Container %s is runnning with container config: %+v\n", result.ContainerID, c)
	return &d, &result
}

func stopContainer(d *task.Docker, containerID string) *task.DockerResult {
	result := d.Stop(containerID)
	if result.Error != nil {
		fmt.Printf("Error stopping container: %v\n", result.Error)
		return nil
	}

	fmt.Printf("Container %s is stopped and removed\n", result.ContainerID)
	return &result
}
