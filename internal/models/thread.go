package models

type Thread struct {
	ID      int    `json:"id"`
	Author  string `json:"author"`
	Created string `json:"created"`
	Forum   string `json:"forum"`
	Message string `json:"message,omitempty"`
	Slug    string `json:"slug,omitempty"`
	Title   string `json:"title"`
	Votes   int    `json:"votes"`
}

type Threads []*Thread
