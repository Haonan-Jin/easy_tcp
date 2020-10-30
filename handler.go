package goland

// A handler to handle decoded message
type Handler interface {
	HandleMsg(ctx ConnectionHandler, msg interface{})
	HandleErr(ctx ConnectionHandler, err error)
}

type ConnectionHandler interface {
	Write(interface{})
	Close()
}
