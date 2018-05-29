package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx := signalHandler(3, 5*time.Second)

	// wait for shutdown notification via context
	<-ctx.Done()
	if err := ctx.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Reported error: %v\n", err)
		os.Exit(1)
	}
}

func signalHandler(repeat int, timeout time.Duration) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	go func() {
		defer func() {
			// deregister channel from signal notification
			signal.Stop(sigCh)
			cancel()
		}()

		n := 0
		for {
			select {
			case <-ctx.Done():
			case <-sigCh:
				n++
				fmt.Printf("Signal count: %d\n", n)
				if n >= repeat {
					fmt.Printf("Repeat count reached.  Exiting...\n")
					return
				}
			}
		}
	}()

	return ctx
}
