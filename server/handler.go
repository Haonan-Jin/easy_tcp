package server

type Handler interface {
	Handle(ctx *ContextHandler, msg interface{})
}
