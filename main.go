package main

import (
	"fmt"
	"time"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/louishuyng/orchestrator/task"
	"github.com/louishuyng/orchestrator/worker"
)

func main() {
	db := make(map[uuid.UUID]*task.Task)

	w := worker.Worker{
		Db:    db,
		Queue: *queue.New(),
	}

	t := task.Task{
		ID:    uuid.New(),
		Name:  "test-container-1",
		State: task.Scheduled,
		Image: "strm/helloworld-http",
	}

	fmt.Println("Starting task...")
	w.AddTask(t)

	result := w.RunTask()
	if result.Error != nil {
		panic(result.Error)
	}

	t.ContainerID = result.ContainerID
	fmt.Println("Task running with Container ID:", t.ContainerID)

	time.Sleep(10 * time.Second)

	fmt.Println("Stopping task...")
	t.State = task.Completed
	w.AddTask(t)

	result = w.RunTask()
	if result.Error != nil {
		panic(result.Error)
	}
}
