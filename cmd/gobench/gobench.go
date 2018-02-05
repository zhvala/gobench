package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cmdArgs := ParseCmdArgs()
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

LOOP:
	for {
		select {
		case signal := <-osSignal:
			fmt.Fprintf(os.Stderr, "gobench interrupted by signal: %s\n", signal)
			break LOOP
		case <-stopSignal:
			fmt.Fprintf(os.Stderr, "gobench task finish.\n")
			break LOOP
		default:
			pool.Run(task)
		}
	}

	// Show task result here
	fmt.Println("test")
}
