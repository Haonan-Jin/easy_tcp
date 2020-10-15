package handler

// A handler to handle decoded message
type Handler interface {
	Handle(ctx ContextHandler, msg interface{})
}

type ContextHandler interface {
	Write(interface{})
	Close()
}
