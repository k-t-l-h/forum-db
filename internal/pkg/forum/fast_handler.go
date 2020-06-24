package forum

import (
	"encoding/json"
	"forum-db/internal/models"
	"forum-db/internal/pkg/database"
	"github.com/valyala/fasthttp"
	routing "github.com/qiangxue/fasthttp-routing"
	"strings"
)

func FCreateForum(ctx *routing.Context) error {

	var f models.Forum
	json.Unmarshal(ctx.PostBody(), &f)
	forums, err := database.CreateForum(f)

	switch err {
	case database.OK:
		ctx.SetStatusCode(fasthttp.StatusCreated)
		data, _ := json.Marshal(forums[0])
		ctx.SetBody(data)

	case database.UserNotFound:
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		data, _ := json.Marshal(models.Error{Message: "User not found"})
		ctx.SetBody(data)

	case database.ForumConflict:
		ctx.SetStatusCode(fasthttp.StatusConflict)
		data, _ := json.Marshal(forums[0])
		ctx.SetBody(data)

	}
	return nil
}

func FCreateSlug(ctx *routing.Context) error{

	slug := ctx.Param("slug")

	var t models.Thread
	t.Forum = slug
	json.Unmarshal(ctx.PostBody(), &t)

	th, err := database.CreateSlug(t)

	switch err {
	case database.OK:
		ctx.SetStatusCode(fasthttp.StatusCreated)
		data, _ := json.Marshal(th[0])
		ctx.SetBody(data)


	case database.UserNotFound:
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		data, _ := json.Marshal(models.Error{Message: "User not found"})
		ctx.SetBody(data)

	case database.ForumConflict:
		ctx.SetStatusCode(fasthttp.StatusConflict)
		data, _ := json.Marshal(th[0])
		ctx.SetBody(data)
	}
	return nil
}

func FSlugDetails(ctx *routing.Context) error {
	slug := ctx.Param("slug")
	f := models.Forum{}
	f.Slug = slug

	f, err := database.GetForumBySlag(f)

	switch err {
	case database.OK:
		ctx.SetStatusCode(fasthttp.StatusOK)
		data, _ := json.Marshal(f)
		ctx.SetBody(data)


	case database.NotFound:
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		data, _ := json.Marshal(models.Error{Message: "User not found"})
		ctx.SetBody(data)
	}
	return nil
}

func FSlugThreads(ctx *routing.Context)error {

	slug := ctx.Param("slug")

	limits := strings.Split(string(ctx.QueryArgs().Peek("limit")), ",)")
	sinces := strings.Split(string(ctx.QueryArgs().Peek("since")), ",)")
	descs := strings.Split(string(ctx.QueryArgs().Peek("desc")), ",)")

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
		ctx.SetStatusCode(fasthttp.StatusOK)
		data, _ := json.Marshal(t)
		ctx.SetBody(data)
	case database.NotFound:
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		data, _ := json.Marshal(models.Error{Message: "User not found"})
		ctx.SetBody(data)
	}
	return nil
}

func FSlugUsers(ctx *routing.Context) error{
	slug := ctx.Param("slug")

	limits := strings.Split(string(ctx.QueryArgs().Peek("limit")), ",)")
	sinces := strings.Split(string(ctx.QueryArgs().Peek("since")), ",)")
	descs := strings.Split(string(ctx.QueryArgs().Peek("desc")), ",)")

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
		ctx.SetStatusCode(fasthttp.StatusOK)
		data, _ := json.Marshal(u)
		ctx.SetBody(data)
	case database.NotFound:
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		data, _ := json.Marshal(models.Error{Message: "User not found"})
		ctx.SetBody(data)
	}
	return nil

}
