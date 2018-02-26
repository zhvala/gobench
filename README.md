# gobench

`gobench` is a simple web benchmark wrote by `golang`.

## Installation

```bash
go get -u github.com/zhvala/gobench/cmd/gobench
# PS: adding $GOPATH/bin into your $PATH may be needed
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