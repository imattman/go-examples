package main

import "fmt"

func main() {
	ws := []string{
		"the",
		"quick",
		"brown",
		"fox",
		"jumped",
		"over",
		"the",
		"lazy",
		"dog",
	}

	fmt.Println(ws)
        reverse(ws)
	fmt.Println(ws)
}

func reverse(ws []string) {
	for i, j := 0, len(ws)-1; i < j; i, j = i+1, j-1 {
		ws[i], ws[j] = ws[j], ws[i]
	}
}
