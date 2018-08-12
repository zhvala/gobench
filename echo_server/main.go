package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

var (
	addr    = flag.String("addr", "localhost:8080", "http service address")
	counter int64
)

func hello(rw http.ResponseWriter, req *http.Request) {
	atomic.AddInt64(&counter, 1)
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("hello, world"))
}

func echo(rw http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	atomic.AddInt64(&counter, 1)
	rw.WriteHeader(http.StatusOK)
	io.Copy(rw, req.Body)
}

func main() {
	flag.Parse()
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				qps := atomic.LoadInt64(&counter)
				atomic.AddInt64(&counter, -qps)
				fmt.Println("qps:", qps)
			}
		}
	}()

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/echo", echo)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
