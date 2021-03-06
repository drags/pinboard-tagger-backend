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
		msg := fmt.Sprintf("Failed to encode Posts as JSON: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
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
		msg := fmt.Sprintf("Failed to encode Posts as JSON: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
}

func postDeleteTagHandler(w http.ResponseWriter, r *http.Request) {
	p := authedPinboard(w, r)
	po, err := postByUrl(p, r.FormValue("url"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	t := r.FormValue("tag")
	newTags := make([]string, 0)
	for _, tag := range po.Tags {
		if tag != t {
			newTags = append(newTags, tag)
		}
	}
	po.Tags = newTags

	err = updatePost(p, po)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Removed tag: %s from URL: %s", t, r.FormValue("url"))
}

func postAddTagHandler(w http.ResponseWriter, r *http.Request) {
	p := authedPinboard(w, r)
	po, err := postByUrl(p, r.FormValue("url"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	t := r.FormValue("tag")
	po.Tags = append(po.Tags, t)

	err = updatePost(p, po)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Added tag: %s from URL: %s", t, r.FormValue("url"))
}

func postByUrl(p *pinboard.Pinboard, u string) (pinboard.Post, error) {
	pf := pinboard.PostsFilter{Url: u}
	tmp, err := p.PostsGet(pf)
	if err != nil {
		return pinboard.Post{}, fmt.Errorf("Failed to retrieve post: %v", err)
	}
	return tmp[0], nil
}

func updatePost(p *pinboard.Pinboard, po pinboard.Post) error {
	err := p.PostsAdd(po, false, false)
	if err != nil {
		return fmt.Errorf("Failed to update post: %v", err)
	}
	return nil
}
