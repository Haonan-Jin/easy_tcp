package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"tcp_handler/server"
)

type LengthFixedDecoder struct {
	headerLen int
}

func NewLengthFixedDecoder(headerLen int) *LengthFixedDecoder {
	decoder := new(LengthFixedDecoder)
	decoder.headerLen = headerLen
	return decoder
}

func (lf LengthFixedDecoder) Decode(b *bytes.Buffer) (interface{}, error) {
	if b.Len() < lf.headerLen {
		return nil, io.ErrShortBuffer
	}

	header := make([]byte, 4)
	_, err := b.Read(header)
	if err != nil {
		return nil, err
	}

	bodyLen := int(binary.BigEndian.Uint32(header))

	if b.Len() < bodyLen {
		recent := bytes.NewBuffer(nil)
		recent.Write(header)
		recent.Write(b.Bytes())
		return recent, io.ErrShortBuffer
	}

	body := make([]byte, bodyLen)
	n, err := b.Read(body)
	if err != nil {
		return nil, err
	}

	return string(body[:n]), nil
}

type StringEncoder struct {
}

func (se StringEncoder) Encode(msg interface{}) []byte {
	return []byte(msg.(string))
}

func Handle(ctx *server.ContextHandler, msg interface{}) {
	fmt.Println(msg)
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

	tcpServer.AddEncoder(new(StringEncoder))
	tcpServer.AddDecoder(NewLengthFixedDecoder(4))
	tcpServer.AddHandler(Handle)
	tcpServer.Start()
}
