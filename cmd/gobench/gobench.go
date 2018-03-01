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
	"time"
)

func globalPanicHandle() {
	if err := recover(); err != nil {
		fmt.Fprintf(os.Stderr, "program terminate, error: %s.\n", err)
	}
}

func main() {
	defer globalPanicHandle()
	fmt.Fprintf(os.Stderr, "Gobench - simple web benchmark - version %s\n", AppVersion)
	fmt.Fprintln(os.Stderr, Copyright)
	fmt.Fprintf(os.Stderr, "\n")
	// get cmd args
	cmdArgs := ParseCmdArgs()
	fmt.Fprintln(os.Stderr, cmdArgs)
	task := CreateTask(cmdArgs)
	pool := CreateClientPool(cmdArgs.Clients)

	/* listen sys signal */
	osSignal := make(chan os.Signal, 1)
	sysSignalListen := []os.Signal{
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		// syscall.SIGTSTP,
	}
	signal.Notify(osSignal, sysSignalListen...)
	/* listen stop signal */
	stopSignal := time.After(cmdArgs.Time)

	counter := 0
LOOP:
	for {
		select {
		case signal := <-osSignal:
			fmt.Fprintf(os.Stderr, "gobench interrupted by signal: %s\n", signal)
			break LOOP
		case <-stopSignal:
			break LOOP
		default:
			pool.Run(task)
			counter++
			if counter >= cmdArgs.Clients {
				counter = 0
				if cmdArgs.Interval >= 0 {
					interval := time.Duration(cmdArgs.Interval) * time.Millisecond
					time.Sleep(interval)
				}
			}
		}
	}

	// Show task result here
	pool.ShowResult()
}
