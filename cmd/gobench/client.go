// Copyright 2017-2018 zhvala
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	socks5 "github.com/zhvala/gosocks5"
	"golang.org/x/net/http2"
)

//ClientPool contains a amount of clients, each client runs in a goroutine
type ClientPool struct {
	taskChan   chan *Task
	resultChan chan Result
	taskWg     sync.WaitGroup
	resultWg   sync.WaitGroup
	resultsSum []Result
	startTime  time.Time
}

//CreateClientPool create a client pool, retuen its pointer
func CreateClientPool(clientNum int) *ClientPool {
	if clientNum <= 0 {
		panic("invalid client num")
	}
	pool := ClientPool{
		taskChan:   make(chan *Task),
		resultChan: make(chan Result),
		startTime:  time.Now(),
	}

	go func() {
		for result := range pool.resultChan {
			defer pool.resultWg.Done()
			pool.resultsSum = append(pool.resultsSum, result)
		}
	}()

	for i := 0; i < clientNum; i++ {
		pool.taskWg.Add(1)
		go func() {
			defer pool.taskWg.Done()
			client := &Client{}
			for task := range pool.taskChan {
				pool.resultChan <- client.Process(task)
				pool.resultWg.Add(1)
			}
		}()
	}
	return &pool
}

// Run put a task into taskChan queue
func (pool *ClientPool) Run(task *Task) {
	pool.taskChan <- task
}

// Close close taskChan queue
func (pool *ClientPool) Close() {
	close(pool.taskChan)
	pool.taskWg.Wait()
	close(pool.resultChan)
	pool.resultWg.Wait()
}

// ShowResult show result of all taskChan
func (pool *ClientPool) ShowResult() {
	costSecs := time.Now().Sub(pool.startTime)
	var recvSizeSum int64
	var successNum int
	var totalCost time.Duration
	maxCost := cTimeMin
	minCost := cTimeMax
	statusMap := make(map[int]int)
	for _, result := range pool.resultsSum {
		if result.Success {
			successNum++
			if _, ok := statusMap[result.StatusCode]; ok {
				statusMap[result.StatusCode]++
			} else {
				statusMap[result.StatusCode] = 1
			}
			recvSizeSum += result.RecvSize
			totalCost += result.TimeCost
			if result.TimeCost > maxCost {
				maxCost = result.TimeCost
			} else if result.TimeCost < minCost {
				minCost = result.TimeCost
			}
		} else {
			if _, ok := statusMap[result.StatusCode]; ok {
				statusMap[result.StatusCode]++
			} else {
				statusMap[result.StatusCode] = 1
			}
		}
	}

	totalReq := len(pool.resultsSum)
	avarReq := totalReq / int(costSecs/time.Second)
	avarBytes := recvSizeSum / int64(costSecs/time.Second)
	fmt.Fprintf(os.Stderr, "Request %d times, total cost %s, avarage: %d request/second, ", totalReq, costSecs, avarReq)
	if avarBytes > 0 {
		fmt.Fprintf(os.Stderr, "%d bytes/second.\n", avarBytes)
	} else {
		fmt.Fprintf(os.Stderr, "\n")
	}
	fmt.Fprintf(os.Stderr, "Request success %d times, failed %d times, details:\n", successNum, totalReq-successNum)
	for status, num := range statusMap {
		fmt.Fprintf(os.Stderr, "*status code: %d, %d times\n", status, num)
	}
	avarCost := totalCost / time.Duration(totalReq)
	fmt.Fprintf(os.Stderr, "Response cost max: %s, mix: %s, avarage: %s.\n", maxCost, minCost, avarCost)
}

// Client http client
type Client struct {
}

// Process do http request
func (client *Client) Process(task *Task) (result Result) {
	if task == nil || (task.HTTPVersion != HTTP && task.HTTPVersion != HTTP2) {
		return
	}

	success := false
	start := time.Now()
	statusCode := -1
	recvSize := int64(0)

	defer func() {
		end := time.Now()
		timeCost := end.Sub(start)
		result = Result{
			Success:    success,
			StatusCode: statusCode,
			RecvSize:   recvSize,
			TimeCost:   timeCost,
		}
	}()

	var httpCli *http.Client

	proxyFunc := (func(*http.Request) (*url.URL, error))(nil)
	dialFunc := net.Dial
	if task.Proxy != "" {
		proxyFunc = func(*http.Request) (*url.URL, error) {
			return url.Parse(task.Proxy)
		}
	} else if task.SOCKS5 != "" {
		client := &socks5.Client{
			Network: "tcp",
			Addr:    task.SOCKS5,
		}
		dialFunc = client.Dial
	}

	transport := &http.Transport{
		Dial:  dialFunc,
		Proxy: proxyFunc,
	}

	if task.HTTPVersion == HTTP2 {
		http2.ConfigureTransport(transport)
	}

	httpCli = &http.Client{
		Timeout:   task.Timeout,
		Transport: transport,
	}

	var body io.Reader
	if task.Data != "" {
		body = strings.NewReader(task.Data)
	}

	req, err := http.NewRequest(task.HTTPMethod, task.URL, body)
	if err != nil {
		return
	}

	rep, err := httpCli.Do(req)
	if err != nil {
		return
	}
	defer rep.Body.Close()

	success = true
	statusCode = rep.StatusCode
	recvSize = rep.ContentLength
	// fmt.Println(rep.Proto)
	return
}

// Result store task result
type Result struct {
	Success    bool
	StatusCode int
	RecvSize   int64
	TimeCost   time.Duration
}
