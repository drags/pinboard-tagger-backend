package main

import (
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {
	log.Println("tagger backend")

	mux := http.NewServeMux()
	mux.Handle("/auth/test", requireAuth(http.HandlerFunc(authTestHandler)))
	mux.Handle("/posts/dates", requireAuth(http.HandlerFunc(postDatesHandler)))

	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":4567", handler)
}
