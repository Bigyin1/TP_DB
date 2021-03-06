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

	v.HandleFunc("/forum/{slug}/details", API.ForumDetails).Methods("GET")
	v.HandleFunc("/forum/{slug}/create", API.CreateThread).Methods("POST")
	v.HandleFunc("/forum/{slug}/users", API.ForumUsers).Methods("GET")

	v.HandleFunc("/forum/{slug}/threads", API.ForumThreadList).Methods("GET")
	v.HandleFunc("/thread/{slug_or_id}/create", API.CreatePosts).Methods("POST")
	v.HandleFunc("/thread/{slug_or_id}/vote", API.ThreadVote).Methods("POST")
	v.HandleFunc("/thread/{slug_or_id}/details", API.ThreadDetails).Methods("GET")
	v.HandleFunc("/thread/{slug_or_id}/posts", API.ThreadPosts).Methods("GET")
	v.HandleFunc("/thread/{slug_or_id}/details", API.ThreadUpdate).Methods("POST")

	v.HandleFunc("/post/{id}/details", API.PostDetails).Methods("GET")
	v.HandleFunc("/post/{id}/details", API.PostUpdate).Methods("POST")

	v.HandleFunc("/service/clear", API.ServiceClear).Methods("POST")
	v.HandleFunc("/service/status", API.ServiceInfo).Methods("GET")

	fmt.Println("Running on port 5000")
	if err = http.ListenAndServe(":5000", r); err != nil {
		fmt.Println("oh, this is error:" + err.Error())
	}
}
