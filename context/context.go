package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	total = 1 * time.Second
	delta = 100 * time.Millisecond
	iter  = 10
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), total)

	// start tasks so that approx half are able to complete before timeout
	start := time.Now()
	timeForTasks := iter * delta
	taskStart := start.Add(timeForTasks / 2)
	fmt.Println(taskStart.Sub(time.Now()))

	var wg sync.WaitGroup
	wg.Add(iter)
	for i := 1; i <= iter; i++ {
		// note: `i` must be converted to a Duration before multiply with `delta`
		sleep := time.Until(taskStart.Add(time.Duration(i) * delta))

		go func(c context.Context, i int, sleep time.Duration) {
			select {
			case <-time.After(sleep):
				fmt.Printf("[%d] done. (%s sleep,)\n", i, sleep)
			case <-ctx.Done():
				fmt.Printf("[%d]: terminated early: %s\n", i, ctx.Err())
			}
			wg.Done()
		}(ctx, i, sleep)
	}

	select {
	case <-time.After(3 * time.Second):
		fmt.Printf("Forcing cancel (%s)...\n", time.Since(start))
	case <-ctx.Done():
		fmt.Println("Timeout triggered by Context.")
	}
	cancel()
	wg.Wait()

	fmt.Printf("\nTime elapsed: %s\n", time.Since(start))
}
