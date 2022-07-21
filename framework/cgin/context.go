package cgin

import (
	"context"
	"encoding/json"
	"net/http"
)

type Context struct {
	w   http.ResponseWriter
	r   *http.Request
	ctx context.Context
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		w:   w,
		r:   r,
		ctx: r.Context(),
	}
}

func (ctx *Context) Json(code int, obj any) error {

	ctx.w.Header().Set("Content-Type", "application/json")

	ctx.w.WriteHeader(code)

	return json.NewEncoder(ctx.w).Encode(obj)

}

type CiginHandler func(c *Context)
