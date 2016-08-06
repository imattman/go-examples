package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/imattman/go-examples/stats"
)

func main() {
	var file string
	flag.StringVar(&file, "f", "", "data file")
	flag.Parse()

	var in io.Reader = os.Stdin
	if file != "" {
		f, err := os.Open(file)
		if err != nil {
			log.Fatalf("Error opening data file %q - %s", file, err)
		}
		defer f.Close()
		in = f
	}

	accum := stats.NewCompactAccumulator(1000)

	scan := bufio.NewScanner(bufio.NewReader(in))
	scan.Split(bufio.ScanWords)

	for i := 1; scan.Scan(); i++ {
		valStr := scan.Text()
		val, err := strconv.Atoi(strings.TrimSpace(valStr))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ignoring bad value %q", valStr)
			continue
		}

		accum.Add(val)
		if i > 3 {
			fmt.Printf("%d: %0.2f\n", i, accum.Mean())
		}
	}
	if scan.Err() != nil {
		log.Fatalf("Error scanning input - %s", scan.Err())
	}

	fmt.Println("summary:", accum)
}
