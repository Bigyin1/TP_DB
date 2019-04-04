package api

import (
	"encoding/json"
	"fmt"
	"gohw/internal/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func (h *Handler) PostDetails(rw http.ResponseWriter, r *http.Request) {
	var (
		details models.PostDetails
		post    models.Post
		user    models.User
		forum   models.Forum
		thread  models.Thread
		err     error
	)
	related := strings.Split(r.FormValue("related"), ",")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	if post, err = h.db.GetPostByID(id); err != nil {
		message := models.Message{Message: "Can't find post"}
		fmt.Printf("PostDetails err: %s\n", err.Error())
		response(rw, http.StatusNotFound, message)
		return
	}
	details.Post = &post

	for _, rel := range related {
		if rel == "user" {
			if user, err = h.db.GetUserByName(details.Post.Author); err != nil {
				message := models.Message{Message: "Can't find post author"}
				response(rw, http.StatusNotFound, message)
				return
			}
			details.Author = &user
		}
		if rel == "thread" {
			if thread, err = h.db.GetThreadByID(details.Post.Thread); err != nil {
				message := models.Message{Message: "Can't find post thread"}
				response(rw, http.StatusNotFound, message)
				return
			}
			details.Thread = &thread
		}
		if rel == "forum" {
			if forum, err = h.db.GetForumBySlug(details.Post.Forum); err != nil {
				message := models.Message{Message: "Can't find post forum"}
				response(rw, http.StatusNotFound, message)
				return
			}
			details.Forum = &forum
		}
	}
	response(rw, http.StatusOK, details)

	return
}

func (h *Handler) PostUpdate(rw http.ResponseWriter, r *http.Request) {

	var (
		post     models.Post
		message  models.Message
		isEdited bool
		err      error
	)

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	if post, err = h.db.GetPostByID(id); err != nil {
		message = models.Message{Message: "Can't find post"}
		response(rw, http.StatusNotFound, message)
	}

	if err = json.NewDecoder(r.Body).Decode(&message); err != nil {
		fmt.Printf("PostChange error: %s at %s\n", err.Error(), r.URL)
		response(rw, http.StatusBadRequest, nil)
		return
	}
	if post.Message != message.Message {
		isEdited = true
	}
	post.Message = message.Message
	post.IsEdited = isEdited
	if err = h.db.UpdatePost(&post); err != nil {
		fmt.Printf("%s\n", err.Error())
		message = models.Message{Message: "Crash"}
		response(rw, 500, message)
		return
	}
	response(rw, http.StatusOK, post)
}
