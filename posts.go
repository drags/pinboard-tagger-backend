package main

import (
	"encoding/json"
	"errors"
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

func postByUrl(w http.ResponseWriter, r *http.Request) (pinboard.Post, error) {
	p := authedPinboard(w, r)
	u := r.FormValue("url")
	t := r.FormValue("tag")
	if len(u) < 1 || len(t) < 1 {
		http.Error(w, "Both `url` and `tag` parameters are required", http.StatusBadRequest)
	}
	pf := pinboard.PostsFilter{Url: u}
	tmp, err := p.PostsGet(pf)
	if err != nil {
		msg := fmt.Sprintf("Failed to retrieve post: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return pinboard.Post{}, errors.New(msg)
	}
	return tmp[0], nil
}

func updatePost(w http.ResponseWriter, r *http.Request, po pinboard.Post) error {
	p := authedPinboard(w, r)
	err := p.PostsAdd(po, false, false)
	if err != nil {
		msg := fmt.Sprintf("Failed to update post: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return errors.New(msg)
	}
	return nil
}

func postDeleteTag(w http.ResponseWriter, r *http.Request) {
	po, err := postByUrl(w, r)
	if err != nil {
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

	err = updatePost(w, r, po)
	if err != nil {
		return
	}

	fmt.Fprintf(w, "Removed tag: %s from URL: %s", t, r.FormValue("url"))
}

func postAddTag(w http.ResponseWriter, r *http.Request) {
	u := r.FormValue("url")
	t := r.FormValue("tag")
	po, err := postByUrl(w, r)
	if err != nil {
		return
	}

	po.Tags = append(po.Tags, t)

	err = updatePost(w, r, po)
	if err != nil {
		return
	}

	fmt.Fprintf(w, "Added tag: %s from URL: %s", t, u)
}
