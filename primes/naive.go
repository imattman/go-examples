// Naive implementation of prime number finder

package primes

import (
	"fmt"
	"io"
)

func Naive(out io.Writer, limit int) {
	for i := 2; i <= limit; i++ {
		iroot := IntSqrt(i)
		isPrime := true
		for j := 2; j <= iroot; j++ {
			if i%j == 0 {
				isPrime = false
				break
			}
		}
		if isPrime {
			fmt.Fprintln(out, i)
		}
	}
}
