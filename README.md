# Gobench

![Scrutinizer Build](https://img.shields.io/scrutinizer/build/g/filp/whoops.svg)![Go Report Card](https://goreportcard.com/badge/github.com/zhvala/gobench)![License](https://img.shields.io/badge/license-Apache%202.0-green.svg)

`Gobench` is a simple web benchmark written in `golang`.

## Installation

```bash
go get -u github.com/zhvala/gobench/cmd/gobench
# adding $GOPATH/bin into your $PATH may be needed
export PATH=$GOPATH/bin:$PATH
```

## Usage

```bash
gobench [option]... URL:

  -client int
    	Run <n> HTTP clients at once. (default 1)
  -data string
    	Send data only if the method is post.
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
  -time int
    	Run gobench for <sec> seconds. (default 60)
  -timeout int
    	Request timeout <millisecond>. (default 1000)
  -trace
    	Use TRACE request method.
  -version
    	Display program version.
```

## Performance test
To do