package forum

import (
	"forum-db/internal/models"
	"forum-db/internal/pkg/database"
	"forum-db/internal/pkg/response"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"net/http"
)

func CreateForum(writer http.ResponseWriter, request *http.Request) {
	var f models.Forum

	jsonerr := easyjson.UnmarshalFromReader(request.Body, &f)

	if jsonerr != nil {
		panic(jsonerr)
	}

	forums, err := database.CreateForum(f)

	switch err {
	case database.OK:
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		easyjson.MarshalToHTTPResponseWriter(forums[0], writer)

	case database.UserNotFound:
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)
		easyjson.MarshalToHTTPResponseWriter(models.Error{Message: "User not found"}, writer)

	case database.ForumConflict:
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusConflict)
		easyjson.MarshalToHTTPResponseWriter(forums[0], writer)
	}
}

func CreateSlug(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	slug := vars["slug"]
	var t models.Thread
	t.Forum = slug

	jsonerr := easyjson.UnmarshalFromReader(request.Body, &t)

	if jsonerr != nil {
		panic(jsonerr)
	}

	th, err := database.CreateSlug(t)

	switch err {
	case database.OK:
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		easyjson.MarshalToHTTPResponseWriter(th[0], writer)


	case database.UserNotFound:
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)
		easyjson.MarshalToHTTPResponseWriter(models.Error{Message: "User not found"}, writer)

	case database.ForumConflict:
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusConflict)
		easyjson.MarshalToHTTPResponseWriter(th[0], writer)

	default:
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusTeapot)
		easyjson.MarshalToHTTPResponseWriter(th[0], nil)
	}

}

func SlugDetails(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	slug := vars["slug"]

	f := models.Forum{}
	f.Slug = slug

	f, err := database.GetForumBySlag(f)

	switch err {
	case database.OK:
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		easyjson.MarshalToHTTPResponseWriter(f, writer)

	case database.NotFound:
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)
		easyjson.MarshalToHTTPResponseWriter(models.Error{Message: "User not found"}, writer)
	}

}

func SlugThreads(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	slug := vars["slug"]

	query := request.URL.Query()
	limits := query["limit"]
	sinces := query["since"]
	descs := query["desc"]

	limit := ""
	since := ""
	desc := ""

	var f models.Thread
	f.Forum = slug

	if len(limits) > 0 {
		limit = limits[0]
	}

	if len(sinces) > 0 {
		since = sinces[0]
	}

	if len(descs) > 0 {
		desc = descs[0]
	}

	t, status := database.GetForumThreads(f, limit, since, desc)

	switch status {
	case database.OK:
		response.Respond(writer, http.StatusOK, t)
	case database.NotFound:
		response.Respond(writer, http.StatusNotFound, models.Error{Message: "Forum not found"})
	}
}

func SlugUsers(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	slug := vars["slug"]

	query := request.URL.Query()
	limits := query["limit"]
	sinces := query["since"]
	descs := query["desc"]

	limit := ""
	since := ""
	desc := ""

	if len(limits) > 0 {
		limit = limits[0]
	}

	if len(sinces) > 0 {
		since = sinces[0]
	}

	if len(descs) > 0 {
		desc = descs[0]
	}

	var f models.Forum
	f.Slug = slug

	u, status := database.GetForumUsers(f, limit, since, desc)

	switch status {
	case database.OK:
		response.Respond(writer, http.StatusOK, u)
	case database.NotFound:
		response.Respond(writer, http.StatusNotFound, models.Error{Message: "Forum not found"})
	}

}
