package helpers

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Job struct {
	conn    *websocket.Conn
	Message []byte
}

type WorkerPool struct {
	JobQueue   chan Job
	WorkerSize int
	wg         sync.WaitGroup
}

func NewWorkerPool(workerSize int, queueSize int) *WorkerPool {
	return &WorkerPool{
		JobQueue:   make(chan Job, queueSize),
		WorkerSize: workerSize,
	}
}

func (wp *WorkerPool) Start() {

}
