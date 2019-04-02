package models

type User struct {
	About    string `json:"about,omitempty"`
	Email    string `json:"email" validate:"email,required"`
	Fullname string `json:"fullname" validate:"required"`
	Nickname string `json:"nickname"`
}

type Users []*User
