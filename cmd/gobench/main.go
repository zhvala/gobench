// Copyright 2017-2018 zhvala
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func globalPanicHandle() {
	if err := recover(); err != nil {
		fmt.Fprintln(os.Stderr, "program terminate, error: ", err)
	}
}

func showCopyRight() {
	fmt.Fprintln(os.Stdout, "gobench - simple web benchmark - version ", AppVersion)
	fmt.Fprintln(os.Stdout, Copyright)
	fmt.Fprintf(os.Stdout, "\n")
}

func listenSysSignal() chan os.Signal {
	osSignal := make(chan os.Signal, 1)
	sysSignalListen := []os.Signal{
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
	}
	signal.Notify(osSignal, sysSignalListen...)
	return osSignal
}

func main() {
	defer globalPanicHandle()
	// show copy right
	showCopyRight()

	// get cmd args
	cmdArgs := ParseCmdArgs()
	fmt.Fprintln(os.Stdout, "Bench start:")
	fmt.Fprintln(os.Stdout, cmdArgs)

	// create bencher
	bencher := NewBencher(cmdArgs)

	/* listen sys signal */
	osSignal := listenSysSignal()

	ctx := bencher.Run()

	select {
	case signal := <-osSignal:
		fmt.Fprintln(os.Stdout, "Bench interrupted by signal: ", signal)
		bencher.Terminate()
	case <-ctx.Done():
		bencher.Close()
	}

	// show task result here
	fmt.Fprintln(os.Stdout, "Bench result:")
	fmt.Fprintln(os.Stdout, StatusFmt(bencher.Status()))
}
