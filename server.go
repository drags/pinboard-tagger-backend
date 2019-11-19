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
	mux.Handle("/posts/get", requireAuth(http.HandlerFunc(postsGetHandler)))
	mux.Handle("/posts/deleteTag", requireAuth(http.HandlerFunc(postDeleteTag)))
	mux.Handle("/posts/addTag", requireAuth(http.HandlerFunc(postAddTag)))

	c := cors.New(cors.Options{
		AllowedHeaders: []string{"Pinboard-Auth"},
	})

	handler := c.Handler(mux)
	http.ListenAndServe(":4567", handler)
}
