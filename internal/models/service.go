package models

type Service struct {
	Post   int `json:"post"`
	Forum  int `json:"forum"`
	Thread int `json:"thread"`
	User   int `json:"user"`
}
