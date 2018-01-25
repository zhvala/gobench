package pool

import (
	"sync"
)

type Job interface {
}

type Worker interface {
	Task(*Job)
}

type WorkerPool struct {
	job chan *Job
	wg  sync.WaitGroup
}

func CreateWorkerPool(poolSize, queueSize int) *WorkerPool {
	// pool := WorkerPool{
	// 	jobs: make(chan *ScanJob, queueSize),
	// }
	// for i := 0; i < poolSize; i++ {
	// 	go func() {
	// 		worker := &ScanWorker{}
	// 		for job := range pool.jobs {
	// 			worker.Process(job)
	// 		}
	// 	}()
	// }
	// return &pool
}
