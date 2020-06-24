package post

import (
	"encoding/json"
	"forum-db/internal/models"
	"forum-db/internal/pkg/database"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
)

//GET /post/{id}/details
func FPostDetails(ctx *routing.Context) error{
	ids := ctx.Param("id")
	id, _ := strconv.Atoi(ids)

	relateds := strings.Split(string(ctx.QueryArgs().Peek("related")), ",)")

	related := []string{}

	if len(relateds) > 0 {
		related = strings.Split(relateds[0], ",")
	}



	pu := models.PostFull{}

	json.Unmarshal(ctx.PostBody(), &pu)

	pu.Post.Id = id
	res, status := database.GetPost(pu, related)

	switch status {
	case database.OK:
		ctx.SetStatusCode(fasthttp.StatusOK)
		data, _ := json.Marshal(res)
		ctx.Write(data)

	case database.NotFound:
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		data, _ := json.Marshal(models.Error{Message: "Thread not in forum"})
		ctx.Write(data)
	}
	return nil
}

//POST /post/{id}/details
func FPostUpdateDetails(ctx *routing.Context) error {

	ids := ctx.Param("id")


	pu := models.PostUpdate{}
	json.Unmarshal(ctx.PostBody(), &pu)
	id, err := strconv.Atoi(ids)

	if err == nil {
		pu.Id = id
	}

	up, status := database.UpdatePost(pu)
	switch status {
	case database.OK:
		ctx.SetStatusCode(fasthttp.StatusOK)
		data, _ := json.Marshal(up)
		ctx.Write(data)
	default:
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		data, _ := json.Marshal(models.Error{Message: "Something went wrong"})
		ctx.Write(data)

	}

	return nil
}
