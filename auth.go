package main

import (
	"fmt"
	"github.com/drags/pinboard"
	"net/http"
	"strings"
)

func authedPinboard(w http.ResponseWriter, r *http.Request) *pinboard.Pinboard {
	auth := r.Header.Get("Pinboard-Auth")
	authInfo := strings.Split(auth, ":")
	if len(authInfo) != 2 {
		msg := fmt.Errorf("Pinboard-Auth header should be in the form <username>:<token>")
		http.Error(w, msg.Error(), http.StatusUnauthorized)
		return nil
	}
	username := authInfo[0]
	token := authInfo[1]
	if len(username) < 1 || len(token) < 1 {
		msg := fmt.Errorf("Pinboard-Auth header should be in the form <username>:<token>")
		http.Error(w, msg.Error(), http.StatusUnauthorized)
		return nil
	}

	p := &pinboard.Pinboard{User: authInfo[0], Token: authInfo[1]}
	return p
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	p := authedPinboard(w, r)

	lu, err := p.PostsUpdated()
	if err != nil {
		msg := fmt.Sprintf("login failed %v", err.Error())
		http.Error(w, msg, http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "Login OK, posts last updated %s", lu)
}
