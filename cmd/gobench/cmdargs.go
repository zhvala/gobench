package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

// supports GET|HEAD|OPTION|TRACE
var (
	// UsageInfo cmd args usage info
	UsageInfo = map[string][]interface{}{}
)

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
	// Force, don't wait reply from server 不需要等待服务器的回复
	Force bool
	// Reload, sent reload request 发生重新加载请求
	Reload bool
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

	// Help, show help info 是否显示帮助信息
	var help bool
	flag.BoolVar(&help, "help", false, "This information.")
	flag.BoolVar(&help, "H", false, "This information.")
	flag.BoolVar(&help, "?", false, "This information.")
	// AppVersion, show app version 显示软件版本
	var appVersion bool
	flag.BoolVar(&appVersion, "version", false, "Display program version.")
	flag.BoolVar(&appVersion, "V", false, "Display program version.")

	var clients int
	flag.IntVar(&clients, "client", 1, "Run <n> HTTP clients at once. Default 1.")
	flag.IntVar(&clients, "C", 1, "Run <n> HTTP clients at once. Default 1.")

	var runTime int
	flag.IntVar(&runTime, "time", 60, "Run gobench for <sec> seconds. Default 60.")
	flag.IntVar(&runTime, "T", 60, "Run gobench for <sec> seconds. Default 60.")

	var http2 bool
	flag.BoolVar(&http2, "http2", false, "Use HTTP/2.0 protocol.")
	flag.BoolVar(&http2, "H2", false, "Use HTTP/2.0 protocol.")

	var proxy string
	flag.StringVar(&proxy, "proxy", "", "Use proxy server for request. <host:port>.")
	flag.StringVar(&proxy, "P", "", "Use proxy server for request. <host:port>.")

	var getMethod, headMethod, optionMethod, traceMethod bool
	flag.BoolVar(&getMethod, "get", false, "Use GET request method.")
	flag.BoolVar(&headMethod, "head", false, "Use HEAD request method.")
	flag.BoolVar(&optionMethod, "option", false, "Use OPTIONS request method.")
	flag.BoolVar(&traceMethod, "trace", false, "Use TRACE request method.")

	var force, reload bool
	flag.BoolVar(&force, "force", false, "Don't wait for reply from server.")
	flag.BoolVar(&force, "F", false, "Don't wait for reply from server.")
	flag.BoolVar(&reload, "reload", false, "Send reload request - Pragma: no-cache.")
	flag.BoolVar(&reload, "R", false, "Send reload request - Pragma: no-cache.")

	flag.Parse()

	if help {
		fmt.Println("help")
		os.Exit(0)
	}
	if appVersion {
		fmt.Printf("version: %s\n", AppVersion)
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
		Force:       force,
		Reload:      reload,
	}
	fmt.Println(cmdArgs)
	return
}

func checkURL(url string) bool {
	return true
}

const (
	// AppVersion gobench version
	AppVersion = "version 1.0"
)

// HTTP Method supported
const (
	// GET HTTP GET
	GET = "GET"
	// HEAD HTTP HEAD
	HEAD = "HEAD"
	// OPTION HTTP OPTION
	OPTION = "OPTION"
	// TRACE HTTP TRACE
	TRACE = "TRACE"
)

const (
	// HTTP http1.1 is used as a default version in golang http.client
	HTTP = iota
	// HTTP2 http2
	HTTP2
)

const (
	cDefaultTimeout = time.Second
)

// Task struct
type Task struct {
	URL         string
	HTTPVersion int
	HTTPMethod  string
	HTTPHeader  map[string]string
	Timeout     time.Duration
}

// CreateTask create a task from given cmd args
func CreateTask(cmdArgs CmdArgs) (task *Task) {
	timeout := time.Duration(0)
	if !cmdArgs.Force {
		timeout = cDefaultTimeout
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
		Timeout:     timeout,
	}
	return
}
