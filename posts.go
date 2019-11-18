package main

import (
	"encoding/json"
	"fmt"
	"github.com/drags/pinboard"
	"net/http"
	"time"
)

func postsGetHandler(w http.ResponseWriter, r *http.Request) {
	p := authedPinboard(w, r)

	pf := pinboard.PostsFilter{}

	if len(r.FormValue("date")) > 0 {
		dt, err := time.Parse("2006-01-02", r.FormValue("date"))
		if err != nil {
			msg := fmt.Sprintf("Failed to parse date: %v", err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		pf.Date = dt
	}

	if len(r.FormValue("tag")) > 0 {
		pf.Tags = append(pf.Tags, r.FormValue("tag"))
	}

	po, err := p.PostsGet(pf)
	if err != nil {
		msg := fmt.Sprintf("Failed to retrieve Posts: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(po)
	if err != nil {
		http.Error(w, "Failed to encode Posts as JSON", http.StatusInternalServerError)
		return
	}
}

func postDatesHandler(w http.ResponseWriter, r *http.Request) {
	p := authedPinboard(w, r)

	pd, err := p.PostsDates(r.FormValue("tag"))
	if err != nil {
		msg := fmt.Sprintf("Failed to retrieve PostsDates: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(pd)
	if err != nil {
		http.Error(w, "Failed to encode PostDates as JSON", http.StatusInternalServerError)
		return
	}
}
