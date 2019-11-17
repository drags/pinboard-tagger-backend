package main

import (
	"fmt"
	"github.com/drags/pinboard"
	"net/http"
	"strings"
)

// Ensure's an appropriately formed Pinboard-Auth header is present. Cannot assert
// on format of API keys since the format has changed over the years without expiring
// previous format keys. Auth against the Pinboard API is _explicitly_ not done here
// to avoid adding an extra API call to every request. Auth issues are handled by
// pinboard.get
func requireAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Pinboard-Auth")
		authInfo := strings.Split(auth, ":")
		if len(authInfo) != 2 {
			msg := fmt.Errorf("Pinboard-Auth header should be in the form <username>:<token>")
			http.Error(w, msg.Error(), http.StatusUnauthorized)
			return
		}
		username := authInfo[0]
		token := authInfo[1]
		if len(username) < 1 || len(token) < 1 {
			msg := fmt.Errorf("Pinboard-Auth header should be in the form <username>:<token>")
			http.Error(w, msg.Error(), http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func authedPinboard(w http.ResponseWriter, r *http.Request) *pinboard.Pinboard {
	auth := r.Header.Get("Pinboard-Auth")
	authInfo := strings.Split(auth, ":")

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
