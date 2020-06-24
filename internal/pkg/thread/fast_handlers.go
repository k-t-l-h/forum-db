package thread

import (
	"encoding/json"
	"forum-db/internal/models"
	"forum-db/internal/pkg/database"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
)

//POST // /thread/{slug_or_id}/create
func FCreate(ctx *routing.Context) error{

	slug := ctx.Param("slug_or_id")
	id, err := strconv.Atoi(slug)

	posts := []models.Post{}
	status := 0
	json.Unmarshal(ctx.PostBody(), &posts)

	if err == nil {
		posts, status = database.CreateThreadPostID(id, posts)
	} else {
		posts, status = database.CreateThreadPost(slug, posts)
	}


	switch status {
	case database.OK:
		ctx.SetStatusCode(fasthttp.StatusOK)
		data, _ := json.Marshal(posts)
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

func FUpdate(ctx *routing.Context) error{
	slug := ctx.Param("slug_or_id")
	id, err := strconv.Atoi(slug)
	status := 0


	var t, tr models.Thread
	json.Unmarshal(ctx.PostBody(), &t)

	if err == nil {
		tr, status = database.ThreadUpdateID(id, t)
	}else {
	tr, status = database.ThreadUpdate(slug, t)
	}

	switch status {
	case database.OK:
		ctx.SetStatusCode(fasthttp.StatusOK)
		data, _ := json.Marshal(tr)
		ctx.Write(data)
	case database.NotFound:
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		data, _ := json.Marshal(models.Error{Message: "User not found"})
		ctx.Write(data)
	}
	return nil

}

///thread/{slug_or_id}/vote
func FVote(ctx *routing.Context) error{
	slug := ctx.Param("slug_or_id")
	id, err := strconv.Atoi(slug)
	status := 0



	vote := models.Vote{}
	thread := models.Thread{}
	json.Unmarshal(ctx.PostBody(), &vote)
	if err == nil {
		thread, status = database.ThreadVoteID(id, vote)
	} else {
		thread, status = database.ThreadVote(slug, vote)
	}


	switch status {
	case database.OK:
		ctx.SetStatusCode(fasthttp.StatusOK)
		data, _ := json.Marshal(thread)
		ctx.Write(data)
	case database.NotFound:
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		data, _ := json.Marshal(models.Error{Message: "User not found"})
		ctx.Write(data)

	}
	return nil
}

///thread/{slug_or_id}/vote

func FPosts(ctx *routing.Context) error{

	Ps := []models.Post{}
	slug := ctx.Param("slug_or_id")
	id, err := strconv.Atoi(slug)
	status := 0

	limits := strings.Split(string(ctx.QueryArgs().Peek("limit")), ",)")
	sinces := strings.Split(string(ctx.QueryArgs().Peek("since")), ",)")
	descs := strings.Split(string(ctx.QueryArgs().Peek("desc")), ",)")
	sorts := strings.Split(string(ctx.QueryArgs().Peek("sort")), ",)")

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


	if err == nil {
		Ps, status = database.GetThreadsPostsID(limit, since, desc, sort, id)
	} else {
		Ps, status = database.GetThreadsPosts(limit, since, desc, sort, slug)
	}

	switch status {
	case database.OK:
		ctx.SetStatusCode(fasthttp.StatusOK)
		data, _ := json.Marshal(Ps)
		ctx.Write(data)
	case database.NotFound:
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		data, _ := json.Marshal(models.Error{Message: "User not found"})
		ctx.Write(data)
	}

	return nil
}

///thread/{slug_or_id}/details
func FDetails(ctx *routing.Context) error{
	slug := ctx.Param("slug_or_id")
	id, err := strconv.Atoi(slug)
	status := 0

	var t models.Thread
	json.Unmarshal(ctx.PostBody(), &t)
	t.Slug = slug

	if err == nil {
		t, status = database.GetThreadByID(id, t)
	} else {
		t, status = database.GetThreadBySlug(slug, t)
	}

	switch status {
	case database.OK:
		ctx.SetStatusCode(fasthttp.StatusOK)
		data, _ := json.Marshal(t)
		ctx.Write(data)
	case database.NotFound:
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		data, _ := json.Marshal(models.Error{Message: "User not found"})
		ctx.Write(data)
	}

	return nil
}

