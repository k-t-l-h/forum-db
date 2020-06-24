package user

import (
	"encoding/json"
	"forum-db/internal/models"
	"forum-db/internal/pkg/database"
	"forum-db/internal/pkg/response"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"net/http"
)

//POST
func Create(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	name := vars["nickname"]

	var u models.User
	err := json.NewDecoder(request.Body).Decode(&u)
	u.NickName = name

	if err != nil {
		panic(err)
	}
	user, status := database.CreateUser(u)

	switch status {
	case database.OK:
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		easyjson.MarshalToHTTPResponseWriter(user[0], writer)
	case database.ForumConflict:
		response.Respond(writer, http.StatusConflict, user)

	}

}

func Update(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	name := vars["nickname"]

	var p models.User
	easyjson.UnmarshalFromReader(request.Body, &p)

	p.NickName = name

	u, status := database.UpdateUser(p)

	switch status {
	case database.OK:
		response.Respond(writer, http.StatusOK, u)
	case database.NotFound:
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)
		easyjson.MarshalToHTTPResponseWriter(models.Error{Message: "User not found"}, writer)
	case database.ForumConflict:
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusConflict)
		easyjson.MarshalToHTTPResponseWriter(models.Error{Message: "User not found"}, writer)
	}
}

//GET /user/{nickname}/profile
func Details(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	name := vars["nickname"]

	us := models.User{}
	us.NickName = name

	u, status := database.GetUser(us)

	switch status {
	case database.OK:
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		easyjson.MarshalToHTTPResponseWriter(u, writer)
	case database.NotFound:
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)
		easyjson.MarshalToHTTPResponseWriter(models.Error{Message: "User not found"}, writer)
	}
}
