// Copyright 2017-2018 zhvala
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
	"flag"
	"fmt"
	uri "net/url"
	"os"
	"strings"
	"time"
)

// CmdArgs store cmd args
// CmdArgs 命令行参数
type CmdArgs struct {
	// URL
	// target url
	// 目标地址
	URL string
	// Duration
	// duration
	// 运行时间
	Duration time.Duration
	// Proxy
	// http/https proxy address
	// http/https代理地址
	Proxy string
	// SOCKS5
	// socks5 proxy address
	// socks5 代理地址
	SOCKS5 string
	// Thread
	// concurrent thread
	// 并发线程数
	Thread int
	// Version
	// HTTP version, supports HTTP1.1 HTTP2
	// HTTP协议版本, 支持 HTTP1.1 HTTP2
	Version int
	// Method
	// HTTP method, supports GET HEAD OPTION TRACE
	// HTTP方法, 支持 GET HEAD OPTION TRACE
	Method string
	// Data
	// post data
	// post数据
	Data string
	// Reload
	// sent reload request, no-cache
	// 发生重新加载请求, 禁用缓存
	Reload bool
	// Interval
	// interval between each request of every thread, millisecond, default no interval
	// 客户端发送请求的间隔，单位毫秒, 默认没有间隔
	Interval time.Duration
	// Timeout
	// timeout of request, millisecond
	// 请求超时时间, 单位毫秒
	Timeout time.Duration
}

func (cmdArgs CmdArgs) String() (str string) {
	str += fmt.Sprintln("URL: ", cmdArgs.URL)
	str += fmt.Sprintln("HTTP method: ", cmdArgs.Method)

	if cmdArgs.Version == HTTP2 {
		str += fmt.Sprintln("HTTP version: HTTP2")
	} else {
		str += fmt.Sprintln("HTTP version: HTTP/1.1")
	}

	str += fmt.Sprintln("Thread num: ", cmdArgs.Thread)
	str += fmt.Sprintln("Duration: ", cmdArgs.Duration)

	if cmdArgs.Interval != 0 {
		str += fmt.Sprintln("Request interval: ", cmdArgs.Interval)
	} else {
		str += fmt.Sprintln("Request interval: none")
	}

	if cmdArgs.Timeout != 0 {
		str += fmt.Sprintln("Request timeout: ", cmdArgs.Timeout)
	} else {
		str += fmt.Sprintln("Request timeout: none")
	}

	if cmdArgs.Reload {
		str += fmt.Sprintln("Disable cache: true")
	}

	if cmdArgs.Proxy != "" {
		str += fmt.Sprintln("Proxy: ", cmdArgs.Proxy)
	}

	if cmdArgs.SOCKS5 != "" {
		str += fmt.Sprintln("Socks5 proxy: ", cmdArgs.SOCKS5)
	}
	return
}

// ParseCmdArgs 从命令行读取参数
// ParseCmdArgs paser args from cmd
func ParseCmdArgs() (cmdArgs *CmdArgs) {
	argc := len(os.Args)
	if argc <= 1 {
		panic("gobench need at least one parameter")
	}

	url := os.Args[argc-1]
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "gobench [option]... URL:\n\n")
		flag.PrintDefaults()
	}

	// AppVersion, show app version 显示软件版本
	var appVersion bool
	flag.BoolVar(&appVersion, "version", false, "Display program version.")

	var thread int
	flag.IntVar(&thread, "thread", 1, "Run <n> threads at once.")

	var duration int
	flag.IntVar(&duration, "duration", 60, "Run gobench for <sec> seconds.")

	var http2 bool
	flag.BoolVar(&http2, "http2", false, "Use HTTP2 protocol.")

	var proxy string
	flag.StringVar(&proxy, "proxy", "", "Use http/https proxy server for request <host:port>.")

	var socks5 string
	flag.StringVar(&socks5, "socks5", "", "Use socks5 proxy server for request <host:port>.")

	var getMethod, postMethod, headMethod, optionMethod, traceMethod bool
	flag.BoolVar(&getMethod, "get", false, "Use GET(default) request method.")
	flag.BoolVar(&postMethod, "post", false, "Use POST request method.")
	flag.BoolVar(&headMethod, "head", false, "Use HEAD request method.")
	flag.BoolVar(&optionMethod, "option", false, "Use OPTIONS request method.")
	flag.BoolVar(&traceMethod, "trace", false, "Use TRACE request method.")

	var data string
	flag.StringVar(&data, "data", "", "Post data, json format supports only.")

	var reload bool
	flag.BoolVar(&reload, "reload", false, "Send reload request - Pragma: no-cache.")

	var interval int
	flag.IntVar(&interval, "interval", 0, "Interval between each request of every client <millisecond>.")

	var timeout int
	flag.IntVar(&timeout, "timeout", 0, "Request timeout <millisecond>.")

	flag.Parse()

	if appVersion {
		fmt.Printf("gobench version %s\n", AppVersion)
		os.Exit(0)
	}

	var httpVersion = HTTP
	if http2 {
		httpVersion = HTTP2
	}

	var httpMethod = GET
	if getMethod {
		httpMethod = GET
	} else if postMethod {
		httpMethod = POST
	} else if headMethod {
		httpMethod = HEAD
	} else if optionMethod {
		httpMethod = OPTION
	} else if traceMethod {
		httpMethod = TRACE
	}

	if !strings.HasPrefix(url, HTTPPrefix) &&
		!strings.HasPrefix(url, HTTPSPrefix) {
		if http2 {
			url = HTTPSPrefix + url
		} else {
			url = HTTPPrefix + url
		}
	} else if strings.HasPrefix(url, HTTPPrefix) {
		if http2 {
			panic("http2 only support https")
		}
	}

	if _, err := uri.ParseRequestURI(url); err != nil {
		panic("invalid target url")
	}

	if proxy != "" {
		if _, err := uri.ParseRequestURI(proxy); err != nil {
			panic("invalid http proxy url")
		}
	}

	if socks5 != "" {
		if _, err := uri.ParseRequestURI(socks5); err != nil {
			panic("invalid socks5 proxy url")
		}
	}

	cmdArgs = &CmdArgs{
		URL:      url,
		Duration: time.Second * time.Duration(duration),
		Proxy:    proxy,
		SOCKS5:   socks5,
		Thread:   thread,
		Version:  httpVersion,
		Method:   httpMethod,
		Data:     data,
		Reload:   reload,
		Interval: time.Millisecond * time.Duration(interval),
		Timeout:  time.Millisecond * time.Duration(timeout),
	}
	return
}
