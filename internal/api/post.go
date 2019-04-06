package api

import (
	"encoding/json"
	"fmt"
	"gohw/internal/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func (h *Handler) PostDetails(rw http.ResponseWriter, r *http.Request) {

	fmt.Println("PostDetails start")
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

	fmt.Println("PostUpdate start")
	var (
		post    models.Post
		message models.Message
		err     error
	)

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	if post, err = h.db.GetPostByID(id); err != nil {
		message = models.Message{Message: "Can't find post"}
		response(rw, http.StatusNotFound, message)
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&message); err != nil {
		fmt.Printf("PostUpdate error: %s at %s\n", err.Error(), r.URL)
		response(rw, http.StatusBadRequest, nil)
		return
	}
	if message.Message != "" && post.Message != message.Message {
		post.Message = message.Message
		post.IsEdited = true
	} else {
		response(rw, http.StatusOK, post)
		return
	}
	if err = h.db.UpdatePost(&post); err != nil {
		fmt.Printf("UpdatePost error: %s\n", err.Error())
		message = models.Message{Message: "Crash"}
		response(rw, 500, message)
		return
	}
	response(rw, http.StatusOK, post)
}

func (h *Handler) CreatePosts(rw http.ResponseWriter, r *http.Request) {

	fmt.Println("CreatePosts start")
	var (
		posts  models.Posts
		thread models.Thread
		err    error
	)

	if id, err := strconv.Atoi(mux.Vars(r)["slug_or_id"]); err != nil {
		if thread, err = h.db.GetThreadBySlug(mux.Vars(r)["slug_or_id"]); err != nil {
			message := models.Message{Message: "Can't find thread for new posts"}
			response(rw, http.StatusNotFound, message)
			return
		}
	} else {
		if thread, err = h.db.GetThreadByID(id); err != nil {
			message := models.Message{Message: "Can't find thread for new posts"}
			response(rw, http.StatusNotFound, message)
			return
		}
	}

	if err = json.NewDecoder(r.Body).Decode(&posts); err != nil {
		fmt.Printf("PostChange error: %s at %s\n", err.Error(), r.URL)
		response(rw, http.StatusBadRequest, nil)
		return
	}

	created := time.Now()
	for _, post := range posts {
		post.Created = created
		post.Forum = thread.Forum
		post.Thread = thread.ID
		if post.Parent != 0 {
			parent, err := h.db.GetPostByID(post.Parent)
			if parent.Thread != post.Thread || err != nil {
				message := models.Message{Message: "Can't find parent post for new post: " + strconv.Itoa(post.ID)}
				response(rw, http.StatusConflict, message)
				return
			}
		}
		if _, err = h.db.GetUserByName(post.Author); err != nil {
			message := models.Message{Message: "Can't find author for new post: " + post.Author}
			response(rw, http.StatusNotFound, message)
			return
		}
		if err = h.db.CreatePost(post); err != nil {
			response(rw, http.StatusInternalServerError, nil)
			return
		}
	}
	response(rw, http.StatusCreated, posts)
}
