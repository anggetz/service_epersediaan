package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello bos")
}

func main() {
	http.HandleFunc("/v1/public-api/hello", hello)

	fmt.Println("server listening")

	http.ListenAndServe(":8090", nil)
}
