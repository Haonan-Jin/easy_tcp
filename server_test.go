package goland

import (
	"fmt"
	"log"
	"net"
	"sync"
	"testing"
)

// decode your data in this func
func StringDecoder(b []byte) (interface{}, error) {
	fmt.Println(string(b))
	return string(b), nil
}

// encode your data in this func
func StringEncoder(msg interface{}) []byte {
	return []byte(msg.(string))
}

type StringHandler struct {
	mutex sync.Mutex
	times int
}

// process decoded message
func (t *StringHandler) HandleMsg(ctx Context, msg interface{}) {
	t.mutex.Lock()
	t.times++
	t.mutex.Unlock()
}

func (t *StringHandler) HandleErr(ctx Context, err error) {
	log.Println("a connect closed, error: ", err)
	ctx.Close()
}

func TestServer(t *testing.T) {
	listener, e := net.Listen("tcp", ":3333")
	if e != nil {
		panic(e)
	}

	for {
		conn, e := listener.Accept()
		if e != nil {
			continue
		}
		handler := NewConnectionHandler(conn)
		handler.AddEncoder(StringEncoder)
		handler.AddDecoder(StringDecoder)
		handler.AddHandler(new(StringHandler))
		handler.Serve()
	}
}
