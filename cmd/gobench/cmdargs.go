package main

import (
	"flag"
	"time"
)

// supports GET|HEAD|OPTION|TRACE
var (
	// UsageInfo cmd args usage info
	UsageInfo = map[string][]interface{}{
		"help":    {"help", false, "This information."},
		"h":       {"h", false, "This information."},
		"?":       {"?", false, "This information."},
		"version": {"version", false, "Display program version."},
		"V":       {"V", false, "Display program version."},
		"client":  {"client", 1, "Run <n> HTTP clients at once. Default one."},
		"C":       {"C", 1, "Run <n> HTTP clients at once. Default one."},
		"time":    {"time", 60, "Run gobench for <sec> seconds. Default 60."},
		"http2":   {"http2", false, "Use HTTP/2.0 protocol."},
		"H2":      {"H2", false, "Use HTTP/2.0 protocol."},
		"proxy":   {"proxy", "", "Use proxy server for request. <host:port>."},
		"P":       {"P", "", "Use proxy server for request. <host:port>."},
		"get":     {"get", false, "Use GET request method."},
		"head":    {"head", false, "Use HEAD request method."},
		"option":  {"option", false, "Use OPTIONS request method."},
		"trace":   {"trace", false, "Use TRACE request method."},
		"force":   {"force", false, "Don't wait for reply from server."},
		"F":       {"F", false, "Don't wait for reply from server."},
		"reload":  {"reload", false, "Send reload request - Pragma: no-cache."},
		"R":       {"R", false, "Send reload request - Pragma: no-cache."},
	}
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
	// Help, show help info 是否显示帮助信息
	Help bool
	// AppVersion, show app version 显示软件版本
	AppVersion string
}

// ParseCmdArgs 从命令行读取参数
// ParseCmdArgs paser args from cmd
func ParseCmdArgs() (cmdArgs CmdArgs) {
	var url string
	url = flag.String("url", "", "url address")
	flag.StringVar(&url, "u")
	return
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
