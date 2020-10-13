package main

import (
	"fmt"
	"net"
	"sync"
	"tcp_handler/server"
)

type decoder struct {
}

func (d *decoder) Decode(data []byte) interface{} {
	return string(data)
}

type stringHandler struct {
	mutex sync.Mutex
	times int
}

func (h *stringHandler) Handle(conn net.Conn, msg interface{}) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.times++
	fmt.Println(h.times)
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

	tcpServer.AddDecoder(server.NewLengthFixedDecoder(4))
	tcpServer.AddHandler(new(stringHandler))
	tcpServer.Start()
}
