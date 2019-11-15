package main

import (
	"fmt"
	"github.com/drags/pinboard"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Failed to parse form", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	token := r.FormValue("token")

	if len(username) < 1 {
		http.Error(w, "login requires username", http.StatusUnauthorized)
		return
	}

	if len(password) < 1 && len(token) < 1 {
		http.Error(w, "login requires either password or token", http.StatusUnauthorized)
		return
	}

	p := &pinboard.Pinboard{User: username}
	if len(token) > 0 {
		p.Token = token
	}

	if len(token) < 1 && len(password) > 0 {
		p.Password = password
	}

	lu, err := p.PostsUpdated()
	if err != nil {
		msg := fmt.Sprintf("login failed %v", err.Error())
		http.Error(w, msg, http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "Login OK, posts last updated %s", lu)

}

func main() {
	log.Println("tagger backend")

	mux := http.NewServeMux()
	mux.HandleFunc("/login", loginHandler)

	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":4567", handler)
}
