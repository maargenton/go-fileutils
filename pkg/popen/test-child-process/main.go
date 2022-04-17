package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var shutdownPeriod time.Duration
	if len(os.Args) > 1 {
		shutdownPeriod, _ = time.ParseDuration(os.Args[1])

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		done := make(chan bool, 1)

		go func() {
			select {
			case sig := <-sigs:
				fmt.Printf("signal: %v\n", sig)
				done <- true
			}
		}()

		fmt.Printf("awaiting signal\n")
		<-done

		fmt.Printf("shuting down ...\n")
		<-time.After(shutdownPeriod)
	}

	fmt.Printf("exiting\n")
}
