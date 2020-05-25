package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

func main() {
	// use a closure so that the invoked function conforms
	// to no-arg signature required by CaptureStdOut(f)
	noArgFunc := func() {
		nPrintMsg("important message", 3)
	}

	cap, err := CaptureStdOut(noArgFunc)
	if err != nil {
		log.Fatalf("Unexpected error: %q", err)
	}

	fmt.Printf("Captured out:\n%q\n", string(cap))
}

// CaptureStdOut takes and invokes a zero-argument function.
// Any data the function writes to stdout is captured and returned.
func CaptureStdOut(f func()) ([]byte, error) {

	// Go doesn't use io.Writer for stdout/stderr
	// so we have to jump through some hoops using an
	// in-memory pipe to accomplish the capture
	pr, pw, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	var buf []byte
	var wg sync.WaitGroup
	wg.Add(1)

	// The two ends of the os.Pipe need to execute in separate goroutines.
	// Launch a separate goroutine for the reader while the writer is used in this.
	go func() {
		buf, err = ioutil.ReadAll(pr)
		if err != nil {
			log.Fatalf("error reading os.Pipe %q", err)
		}

		wg.Done() // capture is done
	}()

	oldOut := os.Stdout
	defer func() { os.Stdout = oldOut }() // replace regular stdout when done here
	os.Stdout = pw

	f()
	pw.Close()
	wg.Wait() // allow the capture routine to finish

	return buf, nil
}

func nPrintMsg(msg string, n int) {
	for i := 0; i < n; i++ {
		fmt.Println(msg)
	}
}
