package models

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Forum struct {
	ID      int    `json:"-"`
	Posts   int    `json:"posts"`
	Slug    string `json:"slug"`
	Threads int    `json:"threads"`
	Title   string `json:"title"`
	Author  string `json:"user"`
}

type ForumUsersQuery struct {
	Slug  string
	Limit int
	Since string
	Desc  bool
}

func (f *ForumUsersQuery) Init(r *http.Request) {

	f.Slug = mux.Vars(r)["slug"]
	f.Limit, _ = strconv.Atoi(r.FormValue("limit"))
	f.Since = r.FormValue("since")
	if desc := r.FormValue("desc"); desc == "true" {
		f.Desc = true
	} else {
		f.Desc = false
	}
}
