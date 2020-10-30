package goland

import (
	"bytes"
	"encoding/binary"
	"net"
	"sync"
)

type ContextHandler struct {
	ConnectionHandler
	mutex sync.Mutex

	conn net.Conn

	buffer  *bytes.Buffer
	decoder Decoder

	dataChan chan int

	handler Handler
	encoder Encoder
}

func handleConnection(conn net.Conn, encoder Encoder, decoder Decoder, handler Handler) {
	contextHandler := new(ContextHandler)

	contextHandler.dataChan = make(chan int, 100)
	contextHandler.buffer = bytes.NewBuffer(nil)
	contextHandler.handler = handler
	contextHandler.conn = conn
	contextHandler.decoder = decoder
	contextHandler.encoder = encoder

	contextHandler.start()
}

func (ch *ContextHandler) start() {
	terminateChan := make(chan error, 1)
	go ch.read(terminateChan)

	select {
	case <-terminateChan:
		return
	}
}

func (ch *ContextHandler) read(terminateChan chan<- error) {

	buffer := make([]byte, 1024)

	for {

		n, err := ch.conn.Read(buffer)
		if err != nil {
			_ = ch.conn.Close()
			terminateChan <- err
			return
		}

		// read bytes to buffer only

		ch.mutex.Lock()
		ch.buffer.Write(buffer[:n])
		ch.mutex.Unlock()

		ch.parseReadBytes()
	}

}

func (ch *ContextHandler) parseReadBytes() {
	for {
		ch.mutex.Lock()
		msg, e := LengthFixedUnpack(ch.buffer)

		if e != nil {
			if msg != nil {
				ch.buffer = bytes.NewBuffer(msg)
			}
			ch.mutex.Unlock()
			break
		}

		decoded, e := ch.decoder(msg)
		if e != nil {
			// failed to decode drop this msg.
			ch.mutex.Unlock()
			break
		}

		ch.mutex.Unlock()

		ch.handler.HandleMsg(ch, decoded)
	}
}

func (ch *ContextHandler) Write(msg interface{}) {
	encoded := ch.encoder(msg)

	msgLen := make([]byte, 4)
	binary.BigEndian.PutUint32(msgLen, uint32(len(encoded)))

	buffer := bytes.NewBuffer(msgLen)
	buffer.Write(encoded)

	ch.conn.Write(buffer.Bytes())
}

func (ch *ContextHandler) Close() {
	_ = ch.conn.Close()
}
