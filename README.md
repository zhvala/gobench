# Gobench

![Scrutinizer Build](https://img.shields.io/scrutinizer/build/g/filp/whoops.svg)![Go Report Card](https://goreportcard.com/badge/github.com/zhvala/gobench)![License](https://img.shields.io/badge/license-Apache%202.0-green.svg)

`Gobench` is a simple web benchmark written in `golang`. It use `goroutine` to simulate concurrent HTTP requests. 

## Feature

- Supports GET, POST, HEAD, OPTION, TRACE
- Support both HTTP/HTTPS and SOCKS5 proxy
- Support HTTP/1.1 and HTTP/2
- High concurrency

## Installation

```bash
go get -u github.com/zhvala/gobench/cmd/gobench
# adding $GOPATH/bin into your $PATH may be needed
export PATH=$GOPATH/bin:$PATH
```

## Usage

```bash
gobench [option]... URL:

  -data string
    	Post data, json format supports only.
  -duration int
    	Run gobench for <sec> seconds. (default 60)
  -get
    	Use GET(default) request method.
  -head
    	Use HEAD request method.
  -http2
    	Use HTTP2 protocol.
  -interval int
    	Interval between each request of every client <millisecond>.
  -option
    	Use OPTIONS request method.
  -post
    	Use POST request method.
  -proxy string
    	Use http/https proxy server for request <host:port>.
  -reload
    	Send reload request - Pragma: no-cache.
  -socks5 string
    	Use socks5 proxy server for request <host:port>.
  -thread int
    	Run <n> threads at once. (default 1)
  -timeout int
    	Request timeout <millisecond>.
  -trace
    	Use TRACE request method.
  -version
    	Display program version.
```

## Author
zhvala(zhvala@foxmail.com)