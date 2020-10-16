package goland

// A handler to handle decoded message
type Handler interface {
	Handle(ctx ConnectionHandler, msg interface{})
}

type ConnectionHandler interface {
	Write(interface{}) (int, error)
	Close()
}
