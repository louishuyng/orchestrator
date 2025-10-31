package main

import (
	"fmt"
	"time"

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
}
