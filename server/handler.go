package server

// A handler to handle decoded message
type Handler interface {
	Handle(ctx *ContextHandler, msg interface{})
}
