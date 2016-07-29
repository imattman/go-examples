// Sieve of Eratosthenes implemented in Go.

package primes

import (
	"fmt"
	"io"
	"math"
)

func IntSqrt(n int) int {
	return int(math.Sqrt(float64(n)))
}

func TraditionalSieve(out io.Writer, limit int) {
	// true values represent a composite/non-prime for the integer represented by the index
	nums := make([]bool, limit+1)

	// special cases 0 and 1 are premarked as non-prime
	nums[0] = true
	nums[1] = true

	// optimization: no need to check candidates beyond square root of limit
	limRoot := IntSqrt(limit)

	for i := 2; i <= limRoot; i++ {
		if !nums[i] {
			// i is the next prime. Disqualify all multiples of i starting from i^2
			// because all lower multiples were marked by prev iterations.
			for p := i * i; p <= limit; p += i {
				nums[p] = true
			}
		}
	}

	// note: printing could be performed in-line above with each prime found
	for i, disqualified := range nums {
		if !disqualified {
			fmt.Fprintln(out, i)
		}
	}
}
