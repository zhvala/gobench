package main

import (
	"time"
)

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

// CmdArgs 接收命令行参数
// CmdArgs store cmd args
type CmdArgs struct {
	// URL target url
	URL string
	// Time duration 运行时间
	Time time.Duration
	// Proxy 使用代理
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
	return
}
