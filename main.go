package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!\n")
	})
	fmt.Println("Listen on 0.0.0.0:8000, see http://127.0.0.1:8000")
	http.ListenAndServe(":8000", nil)
}