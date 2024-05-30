package websocket

import "context"

type Handler interface {
	Register(r *Router)
}

type HandleFunc func(ctx *Context)

type Context struct {
	context.Context

	req *Request
	rf  *ResponseFactory
}

func (ctx *Context) Payload() interface{} {
	return ctx.req.Payload
}

func (ctx *Context) CurrentAction() string {
	return ctx.req.Action
}

func (ctx *Context) ResponseWriter() *ResponseWriter {
	return NewResponseWriter(ctx.rf, ctx.req.Nonce)
}

type Request struct {
	Action  string      `json:"action"`
	Nonce   interface{} `json:"nonce"`
	Payload interface{} `json:"payload"`
}
