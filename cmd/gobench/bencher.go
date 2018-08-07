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
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"sync"
	"sync/atomic"
	"time"

	socks5 "github.com/zhvala/gosocks5"
	"golang.org/x/net/http2"
)

// Status contains:
// reqCounter, num of requests
// sucCounter, success num of requests, status code 200
// failCounter, failed num of requests, status code not 200
// errCounter, err num of requests, other error
// costCounter, total time cost of all requests
// sendCounter, total input bytes
// recvCounter, total output bytes
type Status struct {
	reqCounter  int64
	sucCounter  int64
	failCounter int64
	costCounter int64
	sendCounter int64
	recvCounter int64
}

// StatusFmt convert Status to  readable string
func StatusFmt(duration time.Duration, status Status) (str string) {
	seconds := int64(duration.Seconds())

	str += fmt.Sprintln("Total requests: ", status.reqCounter)
	str += fmt.Sprintln("Success requests: ", status.sucCounter)
	str += fmt.Sprintln("Failed requests: ", status.failCounter)
	str += fmt.Sprintln("Requests per second: ", status.reqCounter/seconds)
	str += fmt.Sprintln("Average response time: ", time.Duration(status.costCounter/status.reqCounter))

	str += fmt.Sprintln("Total send: ", humanize(status.sendCounter))
	str += fmt.Sprintln("Total recv: ", humanize(status.recvCounter))
	str += fmt.Sprintln("Send per second: ", humanize(status.sendCounter/seconds))
	str += fmt.Sprintln("Recv per second: ", humanize(status.recvCounter/seconds))

	return str
}

func humanize(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

//Bencher contains a amount of clients, each client runs in a goroutine
type Bencher struct {
	args     *CmdArgs
	httpTran *http.Transport
	wg       *sync.WaitGroup
	cancel   context.CancelFunc
	status   Status
}

func createTransport(args *CmdArgs) *http.Transport {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	if args.Proxy != "" {
		transport.Proxy = func(*http.Request) (*url.URL, error) {
			return url.Parse(args.Proxy)
		}
	} else if args.SOCKS5 != "" {
		transport.Dial = (&socks5.Client{
			Addr: args.SOCKS5,
		}).Dial
	}

	if args.Version == HTTP2 {
		http2.ConfigureTransport(transport)
	}

	return transport
}

//NewBencher create a client bencher, return its pointer
func NewBencher(args *CmdArgs) *Bencher {
	if args == nil {
		panic("empty args")
	}

	bencher := Bencher{
		args:     args,
		httpTran: createTransport(args),
	}

	return &bencher
}

func (bencher *Bencher) process(ctx context.Context) {
	defer bencher.wg.Done()

	httpCli := &http.Client{
		Timeout:   bencher.args.Timeout,
		Transport: bencher.httpTran,
	}

	var body io.Reader
	if bencher.args.Data != "" {
		body = bytes.NewBuffer([]byte(bencher.args.Data))
	}

	req, err := http.NewRequest(bencher.args.Method, bencher.args.URL, body)
	if err != nil {
		return
	}

	reqDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		return
	}
	sendSize := int64(len(string(reqDump)))

	if body != nil {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
	}

	if bencher.args.Reload {
		req.Header.Add("Pragma", "no-cache")
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			atomic.AddInt64(&bencher.status.reqCounter, 1)
			start := time.Now()
			rep, err := httpCli.Do(req)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				atomic.AddInt64(&bencher.status.failCounter, 1)
				continue
			}
			defer rep.Body.Close()

			atomic.AddInt64(&bencher.status.sendCounter, sendSize)

			if rep.StatusCode == 200 {
				atomic.AddInt64(&bencher.status.sucCounter, 1)
			} else {
				atomic.AddInt64(&bencher.status.failCounter, 1)
			}

			cost := time.Now().Sub(start)

			repDump, err := httputil.DumpResponse(rep, true)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				atomic.AddInt64(&bencher.status.failCounter, 1)
				continue
			}
			atomic.AddInt64(&bencher.status.recvCounter, int64(len(string(repDump))))
			atomic.AddInt64(&bencher.status.costCounter, int64(cost))
		}
	}
}

// Run and stop after deadline
func (bencher *Bencher) Run() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), bencher.args.Duration)
	bencher.cancel = cancel
	bencher.wg = &sync.WaitGroup{}

	bencher.wg.Add(bencher.args.Thread)
	for i := 0; i < bencher.args.Thread; i++ {
		go bencher.process(ctx)
	}

	return ctx
}

// Close all routine and exec cancel func
func (bencher *Bencher) Close() {
	bencher.wg.Wait()
	bencher.cancel()
}

// Terminate exec cancel func before timeout
func (bencher *Bencher) Terminate() {
	bencher.cancel()
	bencher.wg.Wait()
}

// Status return bencher.status
func (bencher *Bencher) Status() Status {
	return bencher.status
}
