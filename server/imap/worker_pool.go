package imap

import (
	"time"
)

type Pool interface {
	Get() *Worker
	Put(worker *Worker) bool
	Close()
}

type Worker struct {
	LastTime time.Time
}

type WorkerPool struct {
	client  chan *Worker
	timeout time.Duration
}

func (pool *WorkerPool) Get() *Worker {
	select {
	case worker := <-pool.client:
		{
			if time.Now().Sub(worker.LastTime) < pool.timeout {
				return worker
			} else {
				worker.LastTime = time.Now()
				return worker
			}
		}
	}
}

func (pool *WorkerPool) Put(worker *Worker) bool {
	if worker != nil {
		worker.LastTime = time.Now()
		pool.client <- worker
	} else {
		return false
	}
	return true
}

func (pool *WorkerPool) Close() {

}

func NewPool(workNum int, timout time.Duration) (pool WorkerPool, err error) {
	mailWorkers := make(chan *Worker, workNum)
	for i := 0; i < cap(mailWorkers); i++ {
		mailWorkers <- &Worker{LastTime: time.Now()}
	}
	pool.client = mailWorkers
	pool.timeout = timout
	return
}
