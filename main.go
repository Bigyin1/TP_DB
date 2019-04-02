package main

import (
	"fmt"
	"gohw/internal/api"

	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	API, err := api.GetHandler()
	if err != nil {
		fmt.Println("Some error happened with configuration file or database" + err.Error())
		return
	}
	r := mux.NewRouter()
	var v = r.PathPrefix("/api").Subrouter()
	//API.Ok()
	v.HandleFunc("/user/{name}/create", API.CreateUser).Methods("POST")
	v.HandleFunc("/user/{name}/profile", API.ProfileUser).Methods("GET")
	v.HandleFunc("/user/{name}/profile", API.UpdateUser).Methods("POST")
	v.HandleFunc("/forum/create", API.CreateForum).Methods("POST")

	// v.HandleFunc("forum/{slug}/details", forum.InfoForum).Methods("GET")
	// v.HandleFunc("forum/{slug}/users", forum.UsersForum).Methods("GET")

	// v.HandleFunc("forum/{slug}/threads", thread.Threads).Methods("GET")
	// v.HandleFunc("thread/{slug_or_id}/create", thread.CreatePosts).Methods("POST")
	// v.HandleFunc("thread/{slug_or_id}/vote", thread.Vote).Methods("POST")
	// v.HandleFunc("thread/{slug_or_id}/details", thread.ThreadInfo).Methods("GET")
	// v.HandleFunc("thread/{slug_or_id}/posts", thread.SortPosts).Methods("GET")
	// v.HandleFunc("thread/{slug_or_id}/details", thread.UpdateThread).Methods("POST")

	// v.HandleFunc("post/{id}/details", post.InfoPost).Methods("GET")
	// v.HandleFunc("post/{id}/details", post.UpdatePost).Methods("POST")

	// v.HandleFunc("service/clear", service.Clear).Methods("POST")
	// v.HandleFunc("service/status", service.Status).Methods("GET")

	fmt.Println("Running on port 3000")
	if err = http.ListenAndServe(":3000", r); err != nil {
		fmt.Println("oh, this is error:" + err.Error())
	}
}
