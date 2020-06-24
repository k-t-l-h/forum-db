package post

import (
	"forum-db/internal/models"
	"forum-db/internal/pkg/database"
	"forum-db/internal/pkg/response"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"net/http"
	"strconv"
	"strings"
)

//GET /post/{id}/details
func PostDetails(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	ids := vars["id"]
	query := request.URL.Query()
	relateds := query["related"]
	related := []string{}

	if len(relateds) > 0 {
		related = strings.Split(relateds[0], ",")
	}
	id, _ := strconv.Atoi(ids)

	pu := models.PostFull{}
	easyjson.UnmarshalFromReader(request.Body, &pu)

	pu.Post.Id = id
	res, status := database.GetPost(pu, related)

	switch status {
	case database.OK:
		response.Respond(writer, http.StatusOK, res)
	case database.NotFound:
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)
		easyjson.MarshalToHTTPResponseWriter(models.Error{Message: "User not found"}, writer)
	}
}

//POST /post/{id}/details
func PostUpdateDetails(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	ids := vars["id"]

	pu := models.PostUpdate{}
	easyjson.UnmarshalFromReader(request.Body, &pu)
	id, err := strconv.Atoi(ids)

	if err == nil {
		pu.Id = id
	}

	up, status := database.UpdatePost(pu)
	switch status {
	case database.OK:
		response.Respond(writer, http.StatusOK, up)
	default:
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)
		easyjson.MarshalToHTTPResponseWriter(models.Error{Message: "User not found"}, writer)

	}

}