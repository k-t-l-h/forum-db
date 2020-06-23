package thread

import (
	"encoding/json"
	"forum-db/internal/models"
	"forum-db/internal/pkg/database"
	"forum-db/internal/pkg/response"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"net/http"
	"strconv"
)

//POST // /thread/{slug_or_id}/create
func Create(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	slug := vars["slug_or_id"]

	posts := []models.Post{}
	decoder := json.NewDecoder(request.Body)
	decoder.Decode(&posts)

	posts, status := database.CreateThreadPost(slug, posts)

	switch status {
	case database.OK:
		response.Respond(writer, http.StatusCreated, posts)
	case database.NotFound:
		response.Respond(writer, http.StatusNotFound, models.Error{Message: "Forum not found"})
	case database.ForumConflict:
		response.Respond(writer, http.StatusConflict, models.Error{Message: "Thread not found"})
	}
}

func Update(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	slug := vars["slug_or_id"]

	var t models.Thread
	easyjson.UnmarshalFromReader(request.Body, &t)

	tr, status := database.ThreadUpdate(slug, t)

	switch status {
	case database.OK:
		//успешно
		response.Respond(writer, http.StatusOK, tr)
	case database.NotFound:
		//нет ветки
		response.Respond(writer, http.StatusNotFound, models.Error{Message: "Thread not found"})
	}

}

func UpdateID(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	slug := vars["id"]
	id, _ := strconv.Atoi(slug)

	var t models.Thread
	easyjson.UnmarshalFromReader(request.Body, &t)

	tr, status := database.ThreadUpdateID(id, t)

	switch status {
	case database.OK:
		//успешно
		response.Respond(writer, http.StatusOK, tr)
	case database.NotFound:
		//нет ветки
		response.Respond(writer, http.StatusNotFound, models.Error{Message: "Thread not found"})
	}

}

///thread/{slug_or_id}/vote
func Vote(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	slug := vars["slug_or_id"]
	vote := models.Vote{}
	easyjson.UnmarshalFromReader(request.Body, &vote)
	thread, status := database.ThreadVote(slug, vote)

	switch status {
	case database.OK:
		response.Respond(writer, http.StatusOK, thread)
	case database.NotFound:
		response.Respond(writer, http.StatusNotFound, models.Error{Message: "Thread not found"})

	}
}

///thread/{slug_or_id}/vote
func VoteID(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	slug := vars["id"]
	vote := models.Vote{}

	id, _ := strconv.Atoi(slug)

	easyjson.UnmarshalFromReader(request.Body, &vote)
	thread, status := database.ThreadVoteID(id, vote)

	switch status {
	case database.OK:
		response.Respond(writer, http.StatusOK, thread)
	case database.NotFound:
		response.Respond(writer, http.StatusNotFound, models.Error{Message: "Thread not found"})

	}
}

//GET

func Posts(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	slug := vars["slug_or_id"]

	query := request.URL.Query()

	limits := query["limit"]
	sinces := query["since"]
	descs := query["desc"]
	sorts := query["sort"]

	limit := ""
	since := ""
	desc := ""
	sort := ""

	if len(limits) > 0 {
		limit = limits[0]
	}

	if len(sinces) > 0 {
		since = sinces[0]
	}

	if len(descs) > 0 {
		desc = descs[0]
	}
	if len(sorts) > 0 {
		sort = sorts[0]
	}

	Ps, status := database.GetThreadsPosts(limit, since, desc, sort, slug)
	switch status {
	case database.OK:
		//успешно
		response.Respond(writer, http.StatusOK, Ps)
	case database.NotFound:
		//нет ветки
		response.Respond(writer, http.StatusNotFound, models.Error{Message: "Thread not found"})
	}

}

///thread/{slug_or_id}/details
func Details(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	slug := vars["slug_or_id"]

	var t models.Thread
	easyjson.UnmarshalFromReader(request.Body, &t)
	t.Slug = slug

	thread, status := database.GetThreadBySlug(slug, t)

	switch status {
	case database.OK:
		//успешно
		response.Respond(writer, http.StatusOK, thread)
	case database.NotFound:
		response.Respond(writer, http.StatusNotFound, models.Error{Message: "Thread not found"})
	}

}

///thread/{slug_or_id}/details
func DetailsID(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	slug := vars["id"]
	id, _ := strconv.Atoi(slug)

	var t models.Thread
	easyjson.UnmarshalFromReader(request.Body, &t)
	t.Id = id

	thread, status := database.GetThreadByID(id, t)

	switch status {
	case database.OK:
		//успешно
		response.Respond(writer, http.StatusOK, thread)
	case database.NotFound:
		response.Respond(writer, http.StatusNotFound, models.Error{Message: "Thread not found"})
	}

}

func CreateID(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	slug := vars["id"]
	id, _ := strconv.Atoi(slug)

	posts := []models.Post{}
	decoder := json.NewDecoder(request.Body)
	decoder.Decode(&posts)

	posts, status := database.CreateThreadPostID(id, posts)

	switch status {
	case database.OK:
		response.Respond(writer, http.StatusCreated, posts)
	case database.NotFound:
		response.Respond(writer, http.StatusNotFound, models.Error{Message: "Forum not found"})
	case database.ForumConflict:
		response.Respond(writer, http.StatusConflict, models.Error{Message: "Thread not found"})
	}
}


func PostsID(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	slug := vars["id"]
	id, _ := strconv.Atoi(slug)

	query := request.URL.Query()

	limits := query["limit"]
	sinces := query["since"]
	descs := query["desc"]
	sorts := query["sort"]

	limit := ""
	since := ""
	desc := ""
	sort := ""

	if len(limits) > 0 {
		limit = limits[0]
	}

	if len(sinces) > 0 {
		since = sinces[0]
	}

	if len(descs) > 0 {
		desc = descs[0]
	}
	if len(sorts) > 0 {
		sort = sorts[0]
	}

	Ps, status := database.GetThreadsPostsID(limit, since, desc, sort, id)
	switch status {
	case database.OK:
		//успешно
		response.Respond(writer, http.StatusOK, Ps)
	case database.NotFound:
		//нет ветки
		response.Respond(writer, http.StatusNotFound, models.Error{Message: "Thread not found"})
	}

}