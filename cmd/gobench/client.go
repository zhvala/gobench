package main

import (
	"sync"
)

//ClientPool contains a amount of clients, each client runs in a goroutine
type ClientPool struct {
	tasks chan *Task
	size  int
	wg    sync.WaitGroup
}

//CreateClientPool create a client pool, retuen its pointer
func CreateClientPool(clientNum int) *ClientPool {
	if clientNum <= 0 {
		panic("invalid client num")
	}
	pool := ClientPool{
		tasks: make(chan *Task),
		size:  clientNum,
	}
	for i := 0; i < clientNum; i++ {
		pool.wg.Add(1)
		go func() {
			defer pool.wg.Done()
			client := &Client{}
			for task := range pool.tasks {
				client.Process(task)
			}
		}()
	}
	return &pool
}

// Run put a task into tasks queue
func (pool *ClientPool) Run(task *Task) {
	pool.tasks <- task
}

// Close close tasks queue
func (pool *ClientPool) Close() {
	close(pool.tasks)
	pool.wg.Wait()
}

// Client http client
type Client struct {
}

// Process do http request
func (client *Client) Process(task *Task) {

}
