package server

type Handler interface {
	Handle(conn *ContextHandler, msg interface{})
}
