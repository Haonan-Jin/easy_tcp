package goland

// A handle to handle decoded message
type Handler interface {
	// Implementation will be called when get a request or response.
	// Msg's type depend on func Decoder's return type.
	HandleMsg(ctx Context, msg interface{})

	// Implementation will be called when an error occurred.
	// Handle errors in the implementation.
	HandleErr(ctx Context, err error)
}

type Context interface {
	Write(interface{})
	ReConn() error
	Close()
}
