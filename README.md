# Gobench

![Scrutinizer Build](https://img.shields.io/badge/build-pass-green.svg)![Go Report Card](https://goreportcard.com/badge/github.com/zhvala/gobench)![License](https://img.shields.io/badge/license-Apache%202.0-green.svg)

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

## Example

```shell
# GET
gobench -thread=50 -duration=120 http://127.0.0.1:8082/hello
# POST
gobench -thread=50 -duration=120 -post -data='{"data": "hello,world"}' http://127.0.0.1:8082/echo
```

## Performance

### Environment

| Hardware | Parameters                                                   |
| -------- | ------------------------------------------------------------ |
| CPU      | Intel(R) Xeon(R) CPU E5-2620 v3 @  2.40GHz (4 x 6 = 24 cores) |
| MEM      | 64GB                                                         |

### Echo-Server-Test

#### GET

```shell
> gobench -thread=50  http://172.31.0.193:8082/hello
gobench - simple web benchmark - version  1.00
Copyright (c) zhvala 2017-2018, Apache 2.0 Open Source Software.

Bench start:
URL:  http://172.31.0.193:8082/hello
HTTP method:  GET
HTTP version: HTTP/1.1
Thread num:  50
Duration:  1m0s
Request interval: none
Request timeout: none

Bench finish.
Total requests:  4982972
Success requests:  4982972
Failed requests:  0
Requests per second:  83049
Average response time:  385.216µs
Total send:  228.1 MiB
Total recv:  613.0 MiB
Send per second:  3.8 MiB
Recv per second:  10.2 MiB
```

#### POST

```shell
> gobench -thread=50 -post -data='{"data": "hello,world"}' http://172.31.0.193:8082/echo
gobench - simple web benchmark - version  1.00
Copyright (c) zhvala 2017-2018, Apache 2.0 Open Source Software.

Bench start:
URL:  http://172.31.0.193:8082/echo
HTTP method:  POST
HTTP version: HTTP/1.1
Thread num:  50
Duration:  1m0s
Request interval: none
Request timeout: none

Bench finish.
Total requests:  3909748
Success requests:  3909748
Failed requests:  0
Requests per second:  65162
Average response time:  535.24µs
Total send:  264.7 MiB
Total recv:  522.0 MiB
Send per second:  4.4 MiB
Recv per second:  8.7 MiB
```

## Author

zhvala(zhvala@foxmail.com)