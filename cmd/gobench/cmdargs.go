// Copyright 2017 zhvala
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
	"flag"
	"fmt"
	"os"
	"time"
)

// supports GET|HEAD|OPTION|TRACE

// CmdArgs 接收命令行参数
// CmdArgs store cmd args
type CmdArgs struct {
	// URL target url
	URL string
	// Time duration 运行时间
	Time time.Duration
	// Proxy use proxy 使用代理
	Proxy string
	// Clients, concurrent clients 并发数
	Clients int
	// HTTP Version, HTTP协议版本 HTTP1.0 HTTP1.1 HTTP2.0
	HTTPVersion int
	// HTTP Method, HTTP方法 GET HEAD OPTION TRACE
	HTTPMethod string
	// Reload, sent reload request 发生重新加载请求
	Reload bool
}

func (cmdArgs CmdArgs) String() (str string) {
	str = fmt.Sprintf("%s %s, currency %d, run %s", cmdArgs.HTTPMethod, cmdArgs.URL, cmdArgs.Clients, cmdArgs.Time)
	if cmdArgs.HTTPVersion == HTTP2 {
		str += fmt.Sprintf(", HTTP/2.0")
	}
	if cmdArgs.Reload {
		str += fmt.Sprintf(", disable cache")
	}
	if cmdArgs.Proxy != "" {
		str += fmt.Sprintf(", proxy: %s", cmdArgs.Proxy)
	}
	return
}

func checkURL(url string) bool {
	// Todo
	return true
}

// ParseCmdArgs 从命令行读取参数
// ParseCmdArgs paser args from cmd
func ParseCmdArgs() (cmdArgs CmdArgs) {
	argc := len(os.Args)
	if argc <= 1 {
		panic("gobench need at least one parameter")
	}
	url := os.Args[argc-1]
	if !checkURL(url) {
		panic("")
	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "gobench [option]... URL:\n\n")
		flag.PrintDefaults()
	}

	// AppVersion, show app version 显示软件版本
	var appVersion bool
	flag.BoolVar(&appVersion, "version", false, "Display program version.")

	var clients int
	flag.IntVar(&clients, "client", 1, "Run <n> HTTP clients at once.")

	var runTime int
	flag.IntVar(&runTime, "time", 60, "Run gobench for <sec> seconds.")

	var http2 bool
	flag.BoolVar(&http2, "http2", false, "Use HTTP/2.0 protocol.")

	var proxy string
	flag.StringVar(&proxy, "proxy", "", "Use proxy server for request. <host:port>.")

	var getMethod, headMethod, optionMethod, traceMethod bool
	flag.BoolVar(&getMethod, "get", false, "Use GET(default) request method.")
	flag.BoolVar(&headMethod, "head", false, "Use HEAD request method.")
	flag.BoolVar(&optionMethod, "option", false, "Use OPTIONS request method.")
	flag.BoolVar(&traceMethod, "trace", false, "Use TRACE request method.")

	var reload bool
	flag.BoolVar(&reload, "reload", false, "Send reload request - Pragma: no-cache.")

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
	} else if headMethod {
		httpMethod = HEAD
	} else if optionMethod {
		httpMethod = OPTION
	} else if traceMethod {
		httpMethod = TRACE
	}

	cmdArgs = CmdArgs{
		URL:         url,
		Time:        time.Second * time.Duration(runTime),
		Proxy:       proxy,
		Clients:     clients,
		HTTPVersion: httpVersion,
		HTTPMethod:  httpMethod,
		Reload:      reload,
	}
	return
}

// Task struct
type Task struct {
	URL         string
	HTTPVersion int
	HTTPMethod  string
	HTTPHeader  map[string]string // Todo, support more parameters
	Proxy       string
	Timeout     time.Duration
}

// CreateTask create a task from given cmd args
func CreateTask(cmdArgs CmdArgs) (task *Task) {
	timeout := cDefaultTimeout
	// if !cmdArgs.Timeout != cTimeMax {
	// 	timeout = cmdArgs.Timeout
	// }
	header := make(map[string]string)
	if cmdArgs.Reload {
		header["Pragma"] = "no-cache"
	}
	task = &Task{
		URL:         cmdArgs.URL,
		HTTPVersion: cmdArgs.HTTPVersion,
		HTTPMethod:  cmdArgs.HTTPMethod,
		HTTPHeader:  header,
		Proxy:       cmdArgs.Proxy,
		Timeout:     timeout,
	}
	return
}
