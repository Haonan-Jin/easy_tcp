package client

import (
	"bytes"
	"encoding/binary"
	"github.com/Haonan-Jin/tcp_server/codec"
	"github.com/Haonan-Jin/tcp_server/handler"
	"net"
)

type TcpClient struct {
	handler.ContextHandler
	decoder codec.Decoder
	encoder codec.Encoder
	handler handler.Handler
	conn    *net.TCPConn
	buffer  *bytes.Buffer
}

func NewTcpClient(localAddr, targetAddr *net.TCPAddr) (*TcpClient, error) {
	conn, e := net.DialTCP("tcp", localAddr, targetAddr)
	if e != nil {
		return nil, e
	}

	client := new(TcpClient)
	client.conn = conn
	client.buffer = bytes.NewBuffer(nil)

	return client, nil
}

func (tc *TcpClient) AddDecoder(decoder codec.Decoder) {
	tc.decoder = decoder
}

func (tc *TcpClient) AddEncoder(encoder codec.Encoder) {
	tc.encoder = encoder
}

func (tc *TcpClient) AddHandler(handler handler.Handler) {
	tc.handler = handler
}

func (tc *TcpClient) Dial() {
	go func() {
		buffer := make([]byte, 1024)
		for {
			i, e := tc.conn.Read(buffer)
			if e != nil {
				return
			}

			tc.buffer.Write(buffer[:i])
			tc.parseReadBytes()
		}
	}()
}

func (tc *TcpClient) parseReadBytes() {
	for {
		msg, e := handler.LengthFixedUnpack(tc.buffer)
		if e != nil {
			if msg != nil {
				tc.buffer = bytes.NewBuffer(msg)
			}
			break
		}

		decoded, e := tc.decoder(msg)
		if e != nil {
			// failed to decode drop this msg.
			break
		}

		tc.handler.Handle(tc, decoded)
	}
}

func (tc *TcpClient) Write(msg interface{}) {
	encoded := tc.encoder(msg)

	msgLen := make([]byte, 4)
	binary.BigEndian.PutUint32(msgLen, uint32(len(encoded)))

	buffer := bytes.NewBuffer(msgLen)
	buffer.Write(encoded)

	data := buffer.Bytes()
	_, _ = tc.conn.Write(data)
}

func (tc *TcpClient) Close() {
	_ = tc.conn.Close()
}
