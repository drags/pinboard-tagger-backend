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

	p := &pinboard.Pinboard{User: authInfo[0], Token: authInfo[1]}
	return p
}

func authTestHandler(w http.ResponseWriter, r *http.Request) {
	p := authedPinboard(w, r)

	lu, err := p.PostsUpdated()
	if err != nil {
		msg := fmt.Sprintf("Authentication failed %v", err.Error())
		http.Error(w, msg, http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "Authentication OK, posts for user %s last updated %s", p.User, lu)
}
