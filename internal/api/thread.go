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

	fmt.Println("CreateThreadstart")
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

	fmt.Println("ForumThreadList start")
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

	fmt.Println("ThreadDetails start")
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

	fmt.Println("ThreadUpdate start")
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
		fmt.Printf("ThreadUpdate error: %s\n", err.Error())
		return
	}

	response(rw, http.StatusOK, thread)
}

func (h *Handler) ThreadPosts(rw http.ResponseWriter, r *http.Request) {

	fmt.Println("ThreadPosts start")
	var (
		posts  models.Posts
		query  models.URLQuery
		thread models.Thread
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
	query.Init(r)
	query.ID = thread.ID

	if err = h.db.GetThreadPosts(&posts, query); err != nil {
		fmt.Printf("ThreadPosts error: %s\n", err.Error())
	}
	response(rw, http.StatusOK, posts)
}

func (h *Handler) ThreadVote(rw http.ResponseWriter, r *http.Request) {

	fmt.Println("ThreadVote start")
	var (
		voice  models.Vote
		thread models.Thread
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
	if err = json.NewDecoder(r.Body).Decode(&voice); err != nil {
		fmt.Printf("ThreadVote error: %s at %s\n", err.Error(), r.URL)
		response(rw, http.StatusBadRequest, nil)
		return
	}
	r.Body.Close()

	if _, err = h.db.GetUserByName(voice.Nickname); err != nil {
		message := models.Message{Message: "Can't find user"}
		response(rw, http.StatusNotFound, message)
		return
	}

	if err = h.db.InsertOrUpdateVote(voice, &thread); err != nil {
		fmt.Printf("ThreadVote error: %s at %s\n", err.Error(), r.URL)
		return
	}

	response(rw, http.StatusOK, thread)

}
