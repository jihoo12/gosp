package main

import (
	"log"
	"net/http"
)

//go:generate go run ../../cmd/gosp -src . -out pages_gen.go -pkg main
func main() {
	mux := http.NewServeMux()
	RegisterPages(mux)
	log.Println("GOSP example: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
