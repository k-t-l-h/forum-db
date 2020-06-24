package main

import (
	"forum-db/internal/pkg/database"
	"forum-db/internal/pkg/forum"
	"forum-db/internal/pkg/post"
	"forum-db/internal/pkg/service"
	"forum-db/internal/pkg/thread"
	"forum-db/internal/pkg/user"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"log"
)



type fastService struct {
	Port   string
	Router *routing.Router
}

func main() {

	if err := database.Open(); err != nil {
		//panic(err)
	}
	defer database.Close()

	router := routing.New()
	router.To("POST", "/api/forum/create", forum.FCreateForum)
	router.To("POST", "/api/forum/<slug>/create", forum.FCreateSlug)
	router.To("GET", "/api/forum/<slug>/details", forum.FSlugDetails)
	router.To("GET", "/api/forum/<slug>/threads", forum.FSlugThreads)
	router.To("GET", "/api/forum/<slug>/users", forum.FSlugUsers)

	router.To("POST", "/api/post/<id>/details", post.FPostUpdateDetails)
	router.To("GET", "/api/post/<id>/details", post.FPostDetails)


	router.To("POST","/api/service/status", service.FStatus)
	router.To("GET", "/api/service/clear", service.FClear)


	router.To("POST", "/api/thread/<slug_or_id>/create", thread.FCreate)
	router.To("POST", "/api/thread/<slug_or_id>/details", thread.FUpdate)
	router.To("POST", "/api/thread/<slug_or_id>/vote", thread.FVote)

	router.To("GET", "/api/thread/<slug_or_id>/details", thread.FDetails)
	router.To("GET", "/api/thread/<slug_or_id>/posts", thread.FPosts)

	router.To("POST", "/api/user/<nickname>/create", user.FCreate)
	router.To("POST", "/api/user/<nickname>/profile", user.FUpdate)
	router.To("GET", "/api/user/<nickname>/profile", user.FDetails)

	s := fastService{
		Port:   ":5000",
		Router: router,
	}
	log.Printf("Server running at %v\n", s.Port)
	fasthttp.ListenAndServe(s.Port, s.Router.HandleRequest)



}
