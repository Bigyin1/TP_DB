package models

type Forum struct {
	ID      int    `json:"-"`
	Posts   int    `json:"posts"`
	Slug    string `json:"slug"`
	Threads int    `json:"threads"`
	Title   string `json:"title"`
	Author  string `json:"user"`
}
