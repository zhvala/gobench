package main

import (
	"time"
)

// Commands 接收命令行参数
// Commands store cmd args
type Commands struct {
	Time     time.Time
	Proxy    string
	Clients  int
	Protocol int
	Method   string
	Help     string
	Version  string
}

// ParserCommandsArgs 从命令行读取参数
// ParserCommandsArgs paser args from cmd
func ParserCommandsArgs() (cmds Commands) {
	return
}
