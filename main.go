package main

import (
	"fmt"
	"github.com/Haonan-Jin/tcp_server/server"
	"net"
	"sync"
)

// decode you data in this func
func Decode(b []byte) (interface{}, error) {
	return string(b), nil
}

func Encode(msg interface{}) []byte {
	return []byte(msg.(string))
}

type StringHandler struct {
	mutex sync.Mutex
	times int
}

// process decoded message
func (t *StringHandler) Handle(ctx *server.ContextHandler, msg interface{}) {
	t.mutex.Lock()
	t.times++
	fmt.Println(t.times)
	t.mutex.Unlock()
}

func main() {
	addr := net.TCPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 3333,
	}

	tcpServer, e := server.NewTcpServer(&addr)
	if e != nil {
		panic(e)
	}

	tcpServer.AddEncoder(Encode)
	tcpServer.AddDecoder(Decode)
	tcpServer.AddHandler(new(StringHandler))
	tcpServer.Start()
}
