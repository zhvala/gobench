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
  -get
    	Use GET(default) request method.
  -head
    	Use HEAD request method.
  -http2
    	Use HTTP/2.0 protocol.
  -option
    	Use OPTIONS request method.
  -proxy string
    	Use proxy server for request. <host:port>.
  -reload
    	Send reload request - Pragma: no-cache.
  -time int
    	Run gobench for <sec> seconds. (default 60)
  -trace
    	Use TRACE request method.
  -version
    	Display program version.
```

## To do

 Add proxy support (socks, etc).