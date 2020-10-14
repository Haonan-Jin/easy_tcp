package main

import (
	"fmt"
	"net"
	"sync"
	"tcp_server/server"
)

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
