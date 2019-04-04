package models

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type URLQuery struct {
	Slug  string
	Limit int
	Since string
	Desc  bool
}

func (f *URLQuery) Init(r *http.Request) {

	f.Slug = mux.Vars(r)["slug"]
	f.Limit, _ = strconv.Atoi(r.FormValue("limit"))
	f.Since = r.FormValue("since")
	if desc := r.FormValue("desc"); desc == "true" {
		f.Desc = true
	} else {
		f.Desc = false
	}
}
