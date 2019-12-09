package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func tagsGetHandler(w http.ResponseWriter, r *http.Request) {
	p := authedPinboard(w, r)

	t, err := p.TagsGet()
	if err != nil {
		msg := fmt.Sprintf("Failed to retrieve Tags: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(t)
	if err != nil {
		msg := fmt.Sprintf("Failed to encode Tags as JSON: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
}
