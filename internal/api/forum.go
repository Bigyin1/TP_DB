package api

import (
	"encoding/json"
	"fmt"
	"gohw/internal/models"
	rerrors "gohw/internal/return_errors"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) CreateForum(rw http.ResponseWriter, r *http.Request) {

	var (
		forum models.Forum
		err   error
	)

	if err = json.NewDecoder(r.Body).Decode(&forum); err != nil {
		fmt.Printf("CreateForum error: %s at %s\n", err.Error(), r.URL)
		response(rw, http.StatusBadRequest, nil)
		return
	}
	r.Body.Close()

	switch err = h.db.CreateForum(&forum); err.Error() {
	case rerrors.UserNotFound:
		message := models.Message{Message: "Can't find user with nickname: " + forum.Author}
		response(rw, http.StatusNotFound, message)
	case rerrors.AlreadyExist:
		response(rw, http.StatusConflict, forum)
	default:
		response(rw, http.StatusCreated, forum)
	}

}

func (h *Handler) ForumDetails(rw http.ResponseWriter, r *http.Request) {

	var (
		forum models.Forum
		err   error
	)

	if forum, err = h.db.GetForumBySlug(mux.Vars(r)["slug"]); err != nil {
		message := models.Message{Message: "Can't find forum with slug: " + mux.Vars(r)["slug"]}
		response(rw, http.StatusBadRequest, message)
		return
	}
	response(rw, http.StatusOK, forum)

}
