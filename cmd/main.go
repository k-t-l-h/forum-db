package main

import (
	"forum-db/internal/pkg/database"
	"forum-db/internal/pkg/forum"
	"forum-db/internal/pkg/post"
	"forum-db/internal/pkg/service"
	"forum-db/internal/pkg/thread"
	"forum-db/internal/pkg/user"
	"github.com/gorilla/mux"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"log"
	"net/http"
)

func main() {

	if err := database.Open(); err != nil {
		//panic(err)
	}
	defer database.Close()

	muxRouter := mux.NewRouter()

	muxRouter.HandleFunc("/api/forum/create", forum.CreateForum).Methods("POST")
	muxRouter.HandleFunc("/api/forum/{slug}/create", forum.CreateSlug).Methods("POST")
	muxRouter.HandleFunc("/api/forum/{slug}/details", forum.SlugDetails).Methods("GET")
	muxRouter.HandleFunc("/api/forum/{slug}/threads", forum.SlugThreads).Methods("GET")
	muxRouter.HandleFunc("/api/forum/{slug}/users", forum.SlugUsers).Methods("GET")

	muxRouter.HandleFunc("/api/post/{id}/details", post.PostDetails).Methods("GET")
	muxRouter.HandleFunc("/api/post/{id}/details", post.PostUpdateDetails).Methods("POST")

	muxRouter.HandleFunc("/api/service/status", service.Status).Methods("GET")
	muxRouter.HandleFunc("/api/service/clear", service.Clear).Methods("POST")

	muxRouter.HandleFunc("/api/thread/{id:[0-9]+}/create", thread.CreateID).Methods("POST")
	muxRouter.HandleFunc("/api/thread/{slug_or_id}/create", thread.Create).Methods("POST")

	muxRouter.HandleFunc("/api/thread/{id:[0-9]+}/details", thread.UpdateID).Methods("POST")
	muxRouter.HandleFunc("/api/thread/{slug_or_id}/details", thread.Update).Methods("POST")

	muxRouter.HandleFunc("/api/thread/{id:[0-9]+}/vote", thread.VoteID).Methods("POST")
	muxRouter.HandleFunc("/api/thread/{slug_or_id}/vote", thread.Vote).Methods("POST")

	muxRouter.HandleFunc("/api/thread/{id:[0-9]+}/details", thread.DetailsID).Methods("GET")
	muxRouter.HandleFunc("/api/thread/{slug_or_id}/details", thread.Details).Methods("GET")

	muxRouter.HandleFunc("/api/thread/{id:[0-9]+}/posts", thread.PostsID).Methods("GET")
	muxRouter.HandleFunc("/api/thread/{slug_or_id}/posts", thread.Posts).Methods("GET")

	muxRouter.HandleFunc("/api/user/{nickname}/create", user.Create).Methods("POST")
	muxRouter.HandleFunc("/api/user/{nickname}/profile", user.Update).Methods("POST")
	muxRouter.HandleFunc("/api/user/{nickname}/profile", user.Details).Methods("GET")

	muxRouter.HandleFunc("/api", service.Index).Methods("GET")
	http.Handle("/", muxRouter)
	log.Print(fasthttp.ListenAndServe(":5000", fasthttpadaptor.NewFastHTTPHandler(muxRouter)))

}
