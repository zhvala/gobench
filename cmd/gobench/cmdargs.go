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
	"flag"
	"fmt"
	"net/url"
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
	// Time
	// duration
	// 运行时间
	Time time.Duration
	// Proxy
	// http/https proxy address
	// http/https代理地址
	Proxy string
	// SOCKS5
	// socks5 proxy address
	// socks5 代理地址
	SOCKS5 string
	// Clients
	// concurrent clients
	// 并发数
	Clients int
	// HTTPVersion
	// HTTP version, supports HTTP1.1 HTTP2
	// HTTP协议版本, 支持 HTTP1.1 HTTP2
	HTTPVersion int
	// HTTPMethod
	// HTTP method, supports GET HEAD OPTION TRACE
	// HTTP方法, 支持 GET HEAD OPTION TRACE
	HTTPMethod string
	// Data
	// post data
	// post数据
	Data string
	// Reload
	// sent reload request, no-cache
	// 发生重新加载请求, 禁用缓存
	Reload bool
	// Interval
	// interval between each request of every client, millisecond, default no interval
	// 客户端发送请求的间隔，单位毫秒, 默认没有间隔
	Interval int
	// Force
	// cancel requests, don't wait response from server after force time
	// 超过时间不等待服务器回复，强制取消请求
	Force int
	// Timeout
	// timeout of request, millisecond
	// 请求超时时间, 单位毫秒
	Timeout int
}

func (cmdArgs CmdArgs) String() (str string) {
	str = fmt.Sprintf("%s %s, currency %d, run %s", cmdArgs.HTTPMethod, cmdArgs.URL, cmdArgs.Clients, cmdArgs.Time)

	if cmdArgs.Interval != 0 {
		str += fmt.Sprintf(", request interval %dms", cmdArgs.Interval)
	} else {
		str += fmt.Sprintf(", no request interval")
	}

	if cmdArgs.Timeout != 0 {
		str += fmt.Sprintf(", request timeout %dms", cmdArgs.Timeout)
	}

	if cmdArgs.HTTPVersion == HTTP2 {
		str += fmt.Sprintf(", HTTP2")
	}

	if cmdArgs.Reload {
		str += fmt.Sprintf(", disable cache")
	}

	if cmdArgs.Proxy != "" {
		str += fmt.Sprintf(", proxy: %s", cmdArgs.Proxy)
	}

	if cmdArgs.SOCKS5 != "" {
		str += fmt.Sprintf(", socks5 proxy: %s", cmdArgs.SOCKS5)
	}
	return
}

func checkTargetURL(task *Task) {
	if !strings.HasPrefix(task.URL, HTTPPrefix) && !strings.HasPrefix(task.URL, HTTPSPrefix) {
		if HTTP2 == task.HTTPVersion {
			task.URL += HTTPSPrefix
		} else {
			task.URL += HTTPPrefix
		}
	} else if strings.HasPrefix(task.URL, HTTPPrefix) {
		if HTTP2 == task.HTTPVersion {
			panic("http2 only support https")
		}
	}
	if !checkURL(task.URL) {
		panic("invalid target url")
	}
}

func checkURL(str string) bool {
	if _, err := url.ParseRequestURI(str); err != nil {
		return false
	}
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
	flag.StringVar(&data, "data", "", "Send data only if the method is post.")

	var reload bool
	flag.BoolVar(&reload, "reload", false, "Send reload request - Pragma: no-cache.")

	var interval int
	flag.IntVar(&interval, "interval", 0, "Interval between each request of every client <millisecond>.")

	var force int
	flag.IntVar(&force, "force", 100, "Client will cancel request and not wait response from server after a given time duration <millisecond>.")

	var timeout int
	flag.IntVar(&timeout, "timeout", 1000, "Request timeout <millisecond>.")

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

	if proxy != "" {
		if !checkURL(proxy) {
			panic("invalid http proxy url")
		}
	}

	if socks5 != "" {
		if !checkURL(socks5) {
			panic("invalid socks5 proxy url")
		}
	}

	if timeout == 0 {
		timeout = cDefaultTimeout
	}

	cmdArgs = CmdArgs{
		URL:         url,
		Time:        time.Second * time.Duration(runTime),
		Proxy:       proxy,
		SOCKS5:      socks5,
		Clients:     clients,
		HTTPVersion: httpVersion,
		HTTPMethod:  httpMethod,
		Data:        data,
		Reload:      reload,
		Force:       force,
		Interval:    interval,
		Timeout:     timeout,
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
	SOCKS5      string
	Data        string
	Timeout     time.Duration
	Force       time.Duration
}

// CreateTask create a task from given cmd args
func CreateTask(cmdArgs CmdArgs) (task *Task) {
	force := time.Duration(cmdArgs.Force) * time.Millisecond
	timeout := force
	if cmdArgs.Force <= cDefaultTimeout {
		timeout = time.Duration(cmdArgs.Timeout) * time.Millisecond
	}

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
		Force:       force,
		SOCKS5:      cmdArgs.SOCKS5,
		Data:        cmdArgs.Data,
	}
	return
}
