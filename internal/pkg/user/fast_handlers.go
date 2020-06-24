package user

import (
	"encoding/json"
	"forum-db/internal/models"
	"forum-db/internal/pkg/database"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

func FCreate(ctx *routing.Context) error{

	name := ctx.Param("nickname")

	var u models.User
	json.Unmarshal(ctx.PostBody(), &u)
	u.NickName = name

	user, status := database.CreateUser(u)

	switch status {
	case database.OK:
		ctx.SetStatusCode(fasthttp.StatusCreated)
		data, _ := json.Marshal(user[0])
		ctx.Write(data)

	case database.ForumConflict:
		ctx.SetStatusCode(fasthttp.StatusConflict)
		data, _ := json.Marshal(user)
		ctx.Write(data)
	}

	return nil
}

func FUpdate(ctx *routing.Context) error{
	name := ctx.Param("nickname")

	var p models.User
	json.Unmarshal(ctx.PostBody(), &p)
	p.NickName = name

	u, status := database.UpdateUser(p)

	switch status {
	case database.OK:
		ctx.SetStatusCode(fasthttp.StatusOK)
		data, _ := json.Marshal(u)
		ctx.Write(data)
	case database.NotFound:
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		data, _ := json.Marshal(models.Error{Message: "User not found"})
		ctx.Write(data)
	case database.ForumConflict:
		ctx.SetStatusCode(fasthttp.StatusConflict)
		data, _ := json.Marshal(models.Error{Message: "StatusConflict"})
		ctx.Write(data)
	}
	return nil
}

//GET /user/{nickname}/profile
func FDetails(ctx *routing.Context)error {
	name := ctx.Param("nickname")

	us := models.User{}
	us.NickName = name

	u, status := database.GetUser(us)

	switch status {
	case database.OK:
		ctx.SetStatusCode(fasthttp.StatusOK)
		data, _ := json.Marshal(u)
		ctx.Write(data)
	case database.NotFound:
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		data, _ := json.Marshal(models.Error{Message: "User not found"})
		ctx.Write(data)
	case database.ForumConflict:
		ctx.SetStatusCode(fasthttp.StatusConflict)
		data, _ := json.Marshal(models.Error{Message: "StatusConflict"})
		ctx.Write(data)
	}
	return nil
}

