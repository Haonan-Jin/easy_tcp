package goland

import (
	"bytes"
	"encoding/binary"
	"net"
)

type TcpClient struct {
	ConnectionHandler
	decoder    Decoder
	encoder    Encoder
	handler    Handler
	conn       *net.TCPConn
	buffer     *bytes.Buffer
	localAddr  *net.TCPAddr
	targetAddr *net.TCPAddr
}

func NewTcpClient(localAddr, targetAddr *net.TCPAddr) (*TcpClient, error) {
	conn, e := net.DialTCP("tcp", localAddr, targetAddr)
	if e != nil {
		return nil, e
	}

	client := new(TcpClient)
	client.localAddr = localAddr
	client.targetAddr = targetAddr
	client.conn = conn
	client.buffer = bytes.NewBuffer(nil)

	return client, nil
}

func (tc *TcpClient) AddDecoder(decoder Decoder) {
	tc.decoder = decoder
}

func (tc *TcpClient) AddEncoder(encoder Encoder) {
	tc.encoder = encoder
}

func (tc *TcpClient) AddHandler(handler Handler) {
	tc.handler = handler
}

func (tc *TcpClient) ReConn() {
	_ = tc.conn.Close()
	conn, err := net.DialTCP("tcp", tc.localAddr, tc.targetAddr)
	if err != nil {
		tc.handler.HandleErr(tc, err)
		return
	}

	tc.conn = conn
	tc.Dial()
}

func (tc *TcpClient) Dial() {
	go func() {
		buffer := make([]byte, 1024)
		for {
			i, e := tc.conn.Read(buffer)
			if e != nil {
				tc.handler.HandleErr(tc, e)
				return
			}

			tc.buffer.Write(buffer[:i])
			tc.parseReadBytes()
		}
	}()
}

func (tc *TcpClient) parseReadBytes() {
	for {
		msg, e := LengthFixedUnpack(tc.buffer)
		if e != nil {
			if msg != nil {
				tc.buffer = bytes.NewBuffer(msg)
			}
			break
		}

		decoded, e := tc.decoder(msg)
		if e != nil {
			// failed to decode drop this msg.
			tc.handler.HandleErr(tc, e)
			break
		}

		tc.handler.HandleMsg(tc, decoded)
	}
}

func (tc *TcpClient) Write(msg interface{}) {
	encoded := tc.encoder(msg)

	msgLen := make([]byte, 4)
	binary.BigEndian.PutUint32(msgLen, uint32(len(encoded)))

	buffer := bytes.NewBuffer(msgLen)
	buffer.Write(encoded)

	data := buffer.Bytes()

	_, e := tc.conn.Write(data)
	if e != nil {
		tc.handler.HandleErr(tc, e)
		return
	}

}

func (tc *TcpClient) Close() {
	_ = tc.conn.Close()
}
