package models

import "time"

type Post struct {
	ID       int       `json:"id,omitempty"`
	Author   string    `json:"author"`
	Created  time.Time `json:"created"`
	Forum    string    `json:"forum,omitempty"`
	IsEdited bool      `json:"isEdited"`
	Message  string    `json:"message"`
	Parent   int       `json:"parent,omitempty"`
	Thread   int       `json:"thread,omitempty"`
}

type Posts []*Post

type PostDetails struct {
	Author *User   `json:"author,omitempty"`
	Forum  *Forum  `json:"forum,omitempty"`
	Post   *Post   `json:"post,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
}
