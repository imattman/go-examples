package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func hello(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/hello/")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello, %s!", name)
}

func main() {
	listen := flag.String("listen", ":8080", "Listen address and port")
	flag.Parse()

	http.HandleFunc("/hello/", hello)
	log.Fatal(http.ListenAndServe(*listen, nil))
}
