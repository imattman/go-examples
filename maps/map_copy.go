package main

import "fmt"

// demonstration that maps do not make a deep copy with
// assignment
func main() {
	m1 := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	m2 := m1 // new ref header points to same data

	fmt.Printf("m1: %v\n", m1)
	fmt.Printf("m2: %v\n", m2)

	m1["four"] = 4
	fmt.Printf("\nAfter editing m1:\n")
	fmt.Printf("m1: %v\n", m1)
	fmt.Printf("m2: %v\n", m2)

	m2["ten"] = 10
	fmt.Printf("\nAfter editing m2:\n")
	fmt.Printf("m1: %v\n", m1)
	fmt.Printf("m2: %v\n", m2)
}
