package goland

import (
	"bytes"
	"net"
	"sync"
)

type ClientContext struct {
	Context

	closedMutex sync.RWMutex
	closed      bool

	conn   net.Conn
	buffer *bytes.Buffer

	decode Decoder
	encode Encoder
	handle Handler
	unPack UnPacker
}

func NewConnectionHandler(conn net.Conn) *ClientContext {
	context := new(ClientContext)
	context.conn = conn
	context.buffer = bytes.NewBuffer(nil)
	context.DefaultUnPacker()
	return context
}

func (ch *ClientContext) DefaultUnPacker() {
	ch.unPack = LengthFixedUnpack
}

func (ch *ClientContext) AddUnPacker(packer UnPacker) {
	ch.unPack = packer
}

func (ch *ClientContext) AddDecoder(decoder Decoder) {
	ch.decode = decoder
}

func (ch *ClientContext) AddEncoder(encoder Encoder) {
	ch.encode = encoder
}

func (ch *ClientContext) AddHandler(handler Handler) {
	ch.handle = handler
}

func (ch *ClientContext) Serve() {
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := ch.conn.Read(buffer)
			if err != nil {
				if ch.isOpen() {
					ch.handle.HandleErr(ch, err)
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
		msg, e := ch.unPack(ch.buffer)
		if e != nil {
			if msg != nil {
				ch.buffer = bytes.NewBuffer(msg)
			}
			break
		}

		decoded, e := ch.decode(msg)
		if e != nil {
			// failed to decode, drop this msg.
			if ch.isOpen() {
				ch.handle.HandleErr(ch, e)
			}
			break
		}

		go ch.handle.HandleMsg(ch, decoded)
	}
}

func (ch *ClientContext) Write(msg interface{}) {
	encoded := ch.encode(msg)
	_, _ = ch.conn.Write(encoded)
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
