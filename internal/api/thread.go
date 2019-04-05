package api

import (
	"encoding/json"
	"fmt"
	"gohw/internal/models"
	rerrors "gohw/internal/return_errors"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) CreateThread(rw http.ResponseWriter, r *http.Request) {

	var (
		thread models.Thread
		err    error
	)

	if err = json.NewDecoder(r.Body).Decode(&thread); err != nil {
		fmt.Printf("CreateForum error: %s at %s\n", err.Error(), r.URL)
		response(rw, http.StatusBadRequest, nil)
		return
	}
	r.Body.Close()

	thread.Forum = mux.Vars(r)["slug"]
	switch err = h.db.CreateThread(&thread); err.Error() {
	case rerrors.UserNotFound:
		message := models.Message{Message: "Can't find user or forum"}
		response(rw, http.StatusNotFound, message)
	case rerrors.AlreadyExist:
		response(rw, http.StatusConflict, thread)
	default:
		response(rw, http.StatusCreated, thread)
	}

}

func (h *Handler) ForumThreadList(rw http.ResponseWriter, r *http.Request) {

	var (
		threads models.Threads
		err     error
		query   models.URLQuery
	)
	query.Init(r)

	if _, err = h.db.GetForumBySlug(query.Slug); err != nil {
		message := models.Message{Message: "Forum not found"}
		response(rw, http.StatusNotFound, message)
		return
	}

	if threads, err = h.db.GetForumThreads(query); err != nil {
		fmt.Printf("ForumUsers error: ", err.Error())
		return
	}
	response(rw, http.StatusOK, threads)

}
