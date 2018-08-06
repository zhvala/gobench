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
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	socks5 "github.com/zhvala/gosocks5"
	"golang.org/x/net/http2"
)

//ClientPool contains a amount of clients, each client runs in a goroutine
type ClientPool struct {
	args     *CmdArgs
	httpTran *http.Transport
	wg       *sync.WaitGroup
	cancel   context.CancelFunc
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

//NewClientPool create a client pool, return its pointer
func NewClientPool(args *CmdArgs) *ClientPool {
	if args == nil {
		panic("empty args")
	}

	pool := ClientPool{
		args:     args,
		httpTran: createTransport(args),
	}

	return &pool
}

func (pool *ClientPool) process(ctx context.Context) {
	defer pool.wg.Done()

	httpCli := &http.Client{
		Timeout:   pool.args.Timeout,
		Transport: pool.httpTran,
	}

	var body io.Reader
	if pool.args.Data != "" {
		body = bytes.NewBuffer([]byte(pool.args.Data))
	}

	req, err := http.NewRequest(pool.args.Method, pool.args.URL, body)
	if err != nil {
		return
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
	}

	if pool.args.Reload {
		req.Header.Add("Pragma", "no-cache")
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			rep, err := httpCli.Do(req)
			if err != nil {
				return
			}
			defer rep.Body.Close()

			_, err = ioutil.ReadAll(rep.Body)
			if err != nil {
				return
			}
		}
	}
}

// Run and stop after deadline
func (pool *ClientPool) Run() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), pool.args.Duration)
	pool.cancel = cancel
	pool.wg = &sync.WaitGroup{}

	for i := 0; i < pool.args.Thread; i++ {
		pool.wg.Add(1)
		go pool.process(ctx)
	}

	return ctx
}

// Close all routine and exec cancel func
func (pool *ClientPool) Close() {
	pool.wg.Wait()
	pool.cancel()
}

// Terminate exec cancel func before timeout
func (pool *ClientPool) Terminate() {
	pool.cancel()
	pool.wg.Wait()
}
