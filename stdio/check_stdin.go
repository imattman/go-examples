package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// Test if stdin is reading from a pipe
	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Fatalln(err)
	}

	if (stat.Mode() & os.ModeCharDevice) == 0 {
		fmt.Println("data is being piped to stdin")
	} else {
		fmt.Println("stdin is from a terminal")
	}
}
