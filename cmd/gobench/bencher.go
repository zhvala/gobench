package main

import (
	"sync"
)

//BenchClient interface
type BenchClient interface {
	Process(*BenchJob)
}

//WorkerPool contains a amount of workers, each worker runs in a goroutine
type WorkerPool struct {
	job chan *BenchJob
	wg  sync.WaitGroup
}

func createWorker(workerType int) (worker BenchWorker) {
	switch workerType {
	default:
	}
	return
}

//CreateWorkerPool create a worker pool, retuen its pointer
func CreateWorkerPool(workerType, workerNum int) *WorkerPool {
	pool := WorkerPool{
		job: make(chan *BenchJob),
	}

	for i := 0; i < workerNum; i++ {
		go func() {
			worker := createWorker(workerType)
			for job := range pool.job {
				worker.Process(job)
			}
		}()
	}
	return &pool
}
