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
	_ = json.NewDecoder(r.Body).Decode(&user)
	r.Body.Close()
	user.Nickname = mux.Vars(r)["name"]

	if err = h.db.InsertNewUser(&user); err != nil {
		fmt.Printf("CreateUser error: %s\n", err.Error())
		users, err := h.db.GetAllCollisionUsers(&user)
		if err != nil {
			response(rw, 500, nil)
			return
		}
		response(rw, 409, users)
		return
	}
	response(rw, http.StatusCreated, nil)
	return
}
