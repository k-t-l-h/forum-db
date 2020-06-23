package service

import (
	"forum-db/internal/pkg/database"
	"forum-db/internal/pkg/response"
	"net/http"
)

// /service/clear
func Clear(writer http.ResponseWriter, request *http.Request) {
	database.Clear()
	response.Respond(writer, http.StatusOK, nil)
}

// /service/status
func Status(writer http.ResponseWriter, request *http.Request) {
	response.Respond(writer, http.StatusOK, database.Info())
}

func Index(writer http.ResponseWriter, request *http.Request) {
	response.Respond(writer, http.StatusOK, nil)
}
