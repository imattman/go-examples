// Prime number algorithms implemented in Go.

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/imattman/go-examples/primes"
)

const defaultLimit = 100

func main() {
	var algorithm string
	flag.StringVar(&algorithm, "a", "sieve", "algorithm: {sieve|naive|filter|channel}")
	flag.Parse()

	var limit = defaultLimit

	if len(flag.Args()) > 0 {
		limStr := flag.Args()[0]
		limInt, err := strconv.Atoi(limStr)
		if err != nil || limInt < 1 {
			fmt.Fprintf(os.Stderr,
				"Illegal argument: %q must be an integer greater than zero\n", limStr)
			os.Exit(1)
		}

		limit = limInt
	}

	switch algorithm {
	case "sieve":
		primes.TraditionalSieve(os.Stdout, limit)
	case "naive":
		primes.Naive(os.Stdout, limit)
	case "filter":
		fmt.Fprintf(os.Stderr, "Not yet implemented\n")
		os.Exit(1)
	case "channel":
		fmt.Fprintf(os.Stderr, "Not yet implemented\n")
		os.Exit(1)
	}
}
