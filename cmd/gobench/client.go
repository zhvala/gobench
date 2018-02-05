package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"sync"

	"golang.org/x/net/http2"
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
	if task.HTTPVersion != HTTP && task.HTTPVersion != HTTP2 {
		return
	}
	var httpCli *http.Client
	if task.HTTPVersion == HTTP2 {
		transport := &http2.Transport{
			AllowHTTP: false, // Allow unsafe connection 允许非加密的链接
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				NextProtos:         []string{"h2"},
			},
			DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) {
				return tls.Dial(netw, addr, cfg)
			},
		}
		httpCli = &http.Client{
			Timeout:   task.Timeout,
			Transport: transport,
		}
	} else {
		httpCli = &http.Client{
			Timeout: task.Timeout,
		}
	}
	req, err := http.NewRequest(task.HTTPMethod, task.URL, nil)
	if err != nil {
		return
	}
	rep, err := httpCli.Do(req)
	if err != nil {
		return
	}
	fmt.Println(rep.Status)
}
