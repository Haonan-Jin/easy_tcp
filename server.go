package goland

import (
	"log"
	"net"
)

type NetWorkServer interface {
	Start()
}

type TcpServer struct {
	NetWorkServer
	listener *net.TCPListener
	decoder  Decoder
	encoder  Encoder
	handler  Handler
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

func (ts *TcpServer) AddDecoder(decoder Decoder) {
	ts.decoder = decoder
}

func (ts *TcpServer) AddEncoder(encoder Encoder) {
	ts.encoder = encoder
}

func (ts *TcpServer) AddHandler(handler Handler) {
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
