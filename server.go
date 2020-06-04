package main

import (
	"flag"
	"github.com/gorilla/handlers"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

func main() {
	listenAddress := flag.String("listen", ":4567", "address:port for server to accept connections")
	flag.Parse()

	log.Println("Starting tagger-backend on", *listenAddress)

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

	handler := handlers.CombinedLoggingHandler(os.Stdout, c.Handler(mux))
	http.ListenAndServe(*listenAddress, handler)
}
