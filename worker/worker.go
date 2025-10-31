package worker

import (
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/louishuyng/orchestrator/task"
)

type Worker struct {
	ID        uuid.UUID
	Queue     queue.Queue
	Db        map[uuid.UUID]*task.Task
	TaskCount int
}
