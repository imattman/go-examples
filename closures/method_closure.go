package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type names []string

func (ns names) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(ns)
}

func index(w http.ResponseWriter, r *http.Request) {
	page := `<html>
	<body><ul>
	<li><a href="/old">old crew</a></li>
	<li><a href="/new">new crew</a></li>
	</ul></body></html>
	`
	fmt.Fprintf(w, page)
}

func main() {
	ns1 := names{
		"Kirk",
		"Spock",
		"McCoy",
		"Scotty",
		"Uhura",
		"Sulu",
		"Chekov",
	}

	ns2 := names{
		"Picard",
		"Riker",
		"Data",
		"LaForge",
		"Worf",
		"Troi",
		"Crusher",
	}

	http.HandleFunc("/", index)
	http.HandleFunc("/old", ns1.ServeHTTP)
	http.HandleFunc("/new", ns2.ServeHTTP)

	port := ":8080"
	fmt.Printf("Listening on %s...\n", port)
	http.ListenAndServe(port, nil)
}
