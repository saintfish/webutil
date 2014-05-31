package webutil

import (
	"encoding/json"
	"github.com/hoisie/web"
)

func Error(ctx *web.Context, err error) {
	ctx.ContentType("text/plain")
	ctx.ResponseWriter.WriteHeader(500)
	ctx.WriteString(err.Error())
}

func Json(ctx *web.Context, obj interface{}) {
	content, err := json.Marshal(obj)
	if err != nil {
		Error(ctx, err)
		return
	}
	ctx.ContentType("application/json")
	ctx.Write(content)
}

func ReadJson(ctx *web.Context, data interface{}) error {
	decoder := json.NewDecoder(ctx.Request.Body)
	return decoder.Decode(data)
}
