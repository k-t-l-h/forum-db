package service

import (
	"encoding/json"
	"forum-db/internal/pkg/database"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

func FClear(ctx *routing.Context) error{

	database.Clear()
	ctx.SetStatusCode(fasthttp.StatusOK)
	return nil
}

// /service/status
func FStatus(ctx *routing.Context) error{
	info := database.Info()
	ctx.SetStatusCode(fasthttp.StatusOK)
	data, _ := json.Marshal(info)
	ctx.Write(data)
	return nil
}

func FIndex(ctx *routing.Context) {
	ctx.SetStatusCode(fasthttp.StatusOK)
}

