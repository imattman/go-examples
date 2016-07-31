// A concurrent prime sieve

package primes

import (
	"fmt"
	"io"
)

// empty struct used as signal on 'done' channel
type sig struct{}

// Send the sequence 2, 3, 4, ... to out channel
func generate(out chan<- int, done chan sig) {
	for i := 2; ; i++ {
		select {
		case out <- i:
		case <-done:
			// fmt.Println("generate done")
			return
		}
	}
}

// Copy the values from channel 'in' to channel 'out',
// removing those divisible by 'prime'.
func filter(prime int, in <-chan int, out chan<- int, done chan sig) {
	for {
		select {
		case nextVal := <-in:
			if nextVal%prime != 0 {
				out <- nextVal
			}
		case <-done:
			// fmt.Printf("filter(%d) done\n", prime)
			return
		}
	}
}

func Concurrent(w io.Writer, limit int) {
	done := make(chan sig) // signal for shutting down
	defer close(done)

	ch := make(chan int)
	go generate(ch, done)

	for i := 0; ; i++ {
		prime := <-ch
		if prime > limit {
			break
		}

		fmt.Fprintln(w, prime)

		// set up next stage, using last 'out' channel as new source
		out := make(chan int)
		go filter(prime, ch, out, done)
		ch = out
	}
}
