package main

import (
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {
	log.Println("tagger backend")

	mux := http.NewServeMux()
	mux.HandleFunc("/login", loginHandler)

	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":4567", handler)
}
