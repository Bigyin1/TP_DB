package api

import (
	"encoding/json"
	"fmt"
	"gohw/internal/models"
	rerrors "gohw/internal/return_errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *Handler) CreateThread(rw http.ResponseWriter, r *http.Request) {

	var (
		thread models.Thread
		err    error
	)

	if err = json.NewDecoder(r.Body).Decode(&thread); err != nil {
		fmt.Printf("CreateThread error: %s at %s\n", err.Error(), r.URL)
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

func (h *Handler) ThreadDetails(rw http.ResponseWriter, r *http.Request) {

	var (
		thread models.Thread
	)

	if id, err := strconv.Atoi(mux.Vars(r)["slug_or_id"]); err != nil {
		if thread, err = h.db.GetThreadBySlug(mux.Vars(r)["slug_or_id"]); err != nil {
			message := models.Message{Message: "Can't find thread"}
			response(rw, http.StatusNotFound, message)
			return
		}
	} else {
		if thread, err = h.db.GetThreadByID(id); err != nil {
			message := models.Message{Message: "Can't find thread"}
			response(rw, http.StatusNotFound, message)
			return
		}
	}

	response(rw, http.StatusOK, thread)
}

func (h *Handler) ThreadUpdate(rw http.ResponseWriter, r *http.Request) {

	var (
		thread models.Thread
		upd    models.Thread
		err    error
	)

	if id, err := strconv.Atoi(mux.Vars(r)["slug_or_id"]); err != nil {
		if thread, err = h.db.GetThreadBySlug(mux.Vars(r)["slug_or_id"]); err != nil {
			message := models.Message{Message: "Can't find thread"}
			response(rw, http.StatusNotFound, message)
			return
		}
	} else {
		if thread, err = h.db.GetThreadByID(id); err != nil {
			message := models.Message{Message: "Can't find thread"}
			response(rw, http.StatusNotFound, message)
			return
		}
	}

	if err = json.NewDecoder(r.Body).Decode(&upd); err != nil {
		fmt.Printf("ThreadUpdate error: %s at %s\n", err.Error(), r.URL)
		response(rw, http.StatusBadRequest, nil)
		return
	}
	r.Body.Close()

	if err = h.db.UpdateThread(&thread, &upd); err != nil {
		fmt.Printf("ThreadUpdate error: %s", err.Error())
		return
	}

	response(rw, http.StatusOK, thread)
}
