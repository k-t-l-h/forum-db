package user

import (
	"encoding/json"
	"forum-db/internal/models"
	"forum-db/internal/pkg/database"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"log"
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
		m := models.User{user[0].About, user[0].Email, user[0].FullName, user[0].NickName}
		data, err := json.Marshal(m)

		ctx.SetBody(data)
		log.Print(err)

	case database.ForumConflict:
		ctx.SetStatusCode(fasthttp.StatusConflict)
		data, _ := json.Marshal(user)
		ctx.SetBody(data)

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
		ctx.SetBody(data)
	case database.NotFound:
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		data, _ := json.Marshal(models.Error{Message: "User not found"})
		ctx.SetBody(data)
	case database.ForumConflict:
		ctx.SetStatusCode(fasthttp.StatusConflict)
		data, _ := json.Marshal(models.Error{Message: "StatusConflict"})
		ctx.SetBody(data)
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
		ctx.SetBody(data)
	case database.NotFound:
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		data, _ := json.Marshal(models.Error{Message: "User not found"})
		ctx.SetBody(data)
	case database.ForumConflict:
		ctx.SetStatusCode(fasthttp.StatusConflict)
		data, _ := json.Marshal(models.Error{Message: "StatusConflict"})
		ctx.SetBody(data)
	}
	return nil
}

