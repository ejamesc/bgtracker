package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func testHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", testHandler)

	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
}
