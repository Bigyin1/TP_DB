package api

import (
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
				message := models.Message{Message: "Can't find post author"}
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
