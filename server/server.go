package server

import (
	"github.com/Haonan-Jin/tcp_server/codec"
	"github.com/Haonan-Jin/tcp_server/handler"
	"log"
	"net"
)

type NetWorkServer interface {
	Start()
}

type TcpServer struct {
	NetWorkServer
	listener *net.TCPListener
	decoder  codec.Decoder
	encoder  codec.Encoder
	handler  handler.Handler
}

func NewTcpServer(addr *net.TCPAddr) (*TcpServer, error) {
	listener, e := listen(addr)
	if e != nil {
		return nil, e
	}
	server := new(TcpServer)
	server.listener = listener

	return server, nil
}

func (ts *TcpServer) AddDecoder(decoder codec.Decoder) {
	ts.decoder = decoder
}

func (ts *TcpServer) AddEncoder(encoder codec.Encoder) {
	ts.encoder = encoder
}

func (ts *TcpServer) AddHandler(handler handler.Handler) {
	ts.handler = handler
}

func listen(addr *net.TCPAddr) (*net.TCPListener, error) {
	listener, e := net.ListenTCP("tcp", addr)
	if e != nil {
		return nil, e
	}

	return listener, nil
}

func (ts *TcpServer) Start() {
	for {
		conn, e := ts.listener.Accept()
		if e != nil {
			log.Println(e)
			continue
		}

		go handleConnection(conn, ts.encoder, ts.decoder, ts.handler)
	}
}
