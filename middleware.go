package main

import (
	"fmt"
	"net/http"
	"strings"
)

// Ensure's an appropriately formed Pinboard-Auth header is present. Cannot assert
// on format of API keys since the format has changed over the years without expiring
// previous format keys. Auth against the Pinboard API is _explicitly_ not done here
// to avoid adding an extra API call to every request. Auth issues are handled by
// pinboard.get
func requireAuth(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		f(w, r)
	}
}

func requireParams(f http.HandlerFunc, params ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		missingParams := make([]string, 0)
		for _, p := range params {
			if len(r.FormValue(p)) < 1 {
				missingParams = append(missingParams, p)
			}
		}

		if len(missingParams) > 0 {
			msg := fmt.Sprintf("The following required parameters are missing: %v", missingParams)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		f(w, r)
	}
}
