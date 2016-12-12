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

// CaptureStdOut takes a zero-argument function and invokes the function while
// capturing all data written to stdout.
func CaptureStdOut(f func()) ([]byte, error) {
	oldOut := os.Stdout
	defer func() {
		os.Stdout = oldOut
	}()

	// Go doesn't use io.Writer for Stdout/Stderr
	// so we have to jump through some hoops using an
	// in-memory pipe to accomplish the capture
	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	var buf []byte
	var wg sync.WaitGroup
	wg.Add(1)

	// fire off the go-routine before the function is invoked and has any
	// chance to write to stdout
	go func() {
		buf, err = ioutil.ReadAll(r)
		if err != nil {
			log.Fatalf("error reading os.Pipe %q", err)
		}

		wg.Done() // signal that capture is done
	}()

	os.Stdout = w
	f()
	w.Close()

	wg.Wait() // allow the capture routine to finish

	return buf, nil
}

func nPrintMsg(msg string, n int) {
	for i := 0; i < n; i++ {
		fmt.Println(msg)
	}
}
