package api

import (
	"encoding/json"
	"fmt"
	"gohw/internal/models"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) CreateUser(rw http.ResponseWriter, r *http.Request) {

	var (
		user models.User
		err  error
	)
	fmt.Println("CreateUser start")
	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		fmt.Printf("CreateUser error: %s at %s\n", err.Error(), r.URL)
		response(rw, http.StatusBadRequest, nil)
		return
	}
	r.Body.Close()
	user.Nickname = mux.Vars(r)["name"]

	fmt.Println(user.Nickname)
	if err = h.db.InsertNewUser(&user); err != nil {
		fmt.Printf("CreateUser error: %s at %s\n", err.Error(), r.URL)
		users, err := h.db.GetAllCollisionUsers(&user)
		if err != nil {
			response(rw, http.StatusInternalServerError, nil)
			return
		}
		response(rw, http.StatusConflict, users)
		return
	}
	response(rw, http.StatusCreated, user)
	return
}

func (h *Handler) ProfileUser(rw http.ResponseWriter, r *http.Request) {

	fmt.Println("ProfileUser start")
	var (
		user models.User
		err  error
	)

	nickname := mux.Vars(r)["name"]
	if user, err = h.db.GetUserByName(nickname); err != nil {
		fmt.Printf("CreateUser error: %s at %s\n", err.Error(), r.URL)
		message := models.Message{Message: "Can't find user by nickname: " + nickname}
		response(rw, http.StatusNotFound, message)
		return
	}
	response(rw, http.StatusOK, user)
}

func (h *Handler) UpdateUser(rw http.ResponseWriter, r *http.Request) {

	fmt.Println("UpdateUser start")
	var (
		user models.User
		err  error
	)

	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		fmt.Printf("UpdateUser error: %s at %s\n", err.Error(), r.URL)
		response(rw, http.StatusBadRequest, nil)
		return
	}
	r.Body.Close()
	user.Nickname = mux.Vars(r)["name"]

	if _, err = h.db.GetUserByName(user.Nickname); err != nil {
		fmt.Printf("UpdateUser error: %s at %s\n", err.Error(), r.URL)
		message := models.Message{Message: "Can't find user by nickname: " + user.Nickname}
		response(rw, http.StatusNotFound, message)
		return
	}

	if err = h.db.UpdateProfile(&user); err != nil {
		fmt.Printf("UpdateUser error: %s at %s\n", err.Error(), r.URL)
		message := models.Message{Message: "This email already in use"}
		response(rw, http.StatusConflict, message)
		return
	}
	response(rw, http.StatusOK, user)

}

func (h *Handler) ForumUsers(rw http.ResponseWriter, r *http.Request) {

	fmt.Println("ForumUsers start")
	var (
		users models.Users
		err   error
		query models.URLQuery
	)
	query.Init(r)

	if _, err = h.db.GetForumBySlug(query.Slug); err != nil {
		message := models.Message{Message: "Forum not found"}
		response(rw, http.StatusNotFound, message)
		return
	}
	if users, err = h.db.GetUsersByForum(query); err != nil {
		fmt.Printf("ForumUsers error: %s\n", err.Error())
		return
	}
	response(rw, http.StatusOK, users)

}
