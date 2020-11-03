package goland

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"sync"
)

type TcpClient struct {
	Context
	decoder    Decoder
	encoder    Encoder
	handler    Handler
	conn       *net.TCPConn
	buffer     *bytes.Buffer
	localAddr  *net.TCPAddr
	targetAddr *net.TCPAddr

	closed      bool
	closedMutex sync.RWMutex
}

// If localAddr is nil, a local address is automatically chosen
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

func (tc *TcpClient) ReConn() error {
	_ = tc.conn.Close()
	conn, err := net.DialTCP("tcp", tc.localAddr, tc.targetAddr)
	if err != nil {
		tc.handler.HandleErr(tc, err)
		return err
	}

	tc.conn = conn
	tc.Dial()
	tc.closed = false
	tc.buffer.Reset()
	return nil
}

func (tc *TcpClient) Dial() {
	go func() {
		buffer := make([]byte, 1024)
		for {
			i, e := tc.conn.Read(buffer)
			if e != nil {
				if errors.Is(e, io.EOF) {
					tc.buffer.Write(buffer)
					tc.parseReadBytes()
					continue
				}
				if tc.isOpen() {
					tc.handler.HandleErr(tc, e)
				}
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
			if tc.isOpen() {
				tc.handler.HandleErr(tc, e)
			}
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
		if tc.isOpen() {
			tc.handler.HandleErr(tc, e)
		}
		return
	}

}

func (tc *TcpClient) isOpen() bool {
	tc.closedMutex.RLock()
	defer tc.closedMutex.RUnlock()
	return !tc.closed
}

func (tc *TcpClient) Close() {
	tc.closedMutex.Lock()
	tc.closed = true
	_ = tc.conn.Close()
	tc.closedMutex.Unlock()
}
