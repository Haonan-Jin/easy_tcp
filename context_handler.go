package goland

import (
	"bytes"
	"encoding/binary"
	"net"
	"sync"
)

type ClientContext struct {
	Context

	closedMutex sync.RWMutex
	closed      bool

	conn    net.Conn
	buffer  *bytes.Buffer
	decoder Decoder
	encoder Encoder
	handler Handler
}

func handleConnection(conn net.Conn, encoder Encoder, decoder Decoder, handler Handler) {
	context := &ClientContext{
		conn:    conn,
		buffer:  bytes.NewBuffer(nil),
		decoder: decoder,
		encoder: encoder,
		handler: handler,
	}

	context.serve()
}

func (ch *ClientContext) serve() {
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := ch.conn.Read(buffer)
			if err != nil {
				if ch.isOpen() {
					ch.handler.HandleErr(ch, err)
				}
				return
			}

			// read bytes to buffer only
			ch.buffer.Write(buffer[:n])
			ch.parseReadBytes()
		}
	}()
}

func (ch *ClientContext) parseReadBytes() {
	for {
		msg, e := LengthFixedUnpack(ch.buffer)
		if e != nil {
			if msg != nil {
				ch.buffer = bytes.NewBuffer(msg)
			}
			break
		}

		decoded, e := ch.decoder(msg)
		if e != nil {
			// failed to decode drop this msg.
			if ch.isOpen() {
				ch.handler.HandleErr(ch, e)
			}
			break
		}

		ch.handler.HandleMsg(ch, decoded)
	}
}

func (ch *ClientContext) Write(msg interface{}) {
	encoded := ch.encoder(msg)

	msgLen := make([]byte, 4)
	binary.BigEndian.PutUint32(msgLen, uint32(len(encoded)))

	buffer := bytes.NewBuffer(msgLen)
	buffer.Write(encoded)

	ch.conn.Write(buffer.Bytes())
}

func (ch *ClientContext) isOpen() bool {
	ch.closedMutex.RLock()
	defer ch.closedMutex.RUnlock()
	return !ch.closed
}

func (ch *ClientContext) ReConn() error { return nil }

func (ch *ClientContext) Close() {
	ch.closedMutex.Lock()
	ch.closed = true
	_ = ch.conn.Close()
	ch.closedMutex.Unlock()
}
