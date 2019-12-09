package main

import (
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {
	log.Println("tagger backend")

	mux := http.NewServeMux()
	mux.HandleFunc("/auth/test", requireAuth(authTestHandler))
	mux.HandleFunc("/posts/dates", requireAuth(postDatesHandler))
	mux.HandleFunc("/posts/get", requireAuth(postsGetHandler))
	mux.HandleFunc("/tags/get", requireAuth(tagsGetHandler))
	mux.HandleFunc("/posts/deleteTag", requireAuth(requireParams(postDeleteTagHandler, "url", "tag")))
	mux.HandleFunc("/posts/addTag", requireAuth(requireParams(postAddTagHandler, "url", "tag")))

	c := cors.New(cors.Options{
		AllowedHeaders: []string{"Pinboard-Auth"},
	})

	handler := c.Handler(mux)
	http.ListenAndServe(":4567", handler)
}
