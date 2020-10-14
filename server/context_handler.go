package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
)

type ContextHandler struct {
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
		msg, e := unpack(ch.buffer)

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

		ch.handler.Handle(ch, decoded)
	}
}

func (ch *ContextHandler) Write(msg interface{}) {
	encode := ch.encoder(msg)
	_, err := ch.conn.Write(encode)
	if err != nil {
		_ = ch.conn.Close()
		fmt.Println(err)
	}
}

func unpack(b *bytes.Buffer) ([]byte, error) {
	if b.Len() < 4 {
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
		return recent.Bytes(), io.ErrShortBuffer
	}

	body := make([]byte, bodyLen)
	n, err := b.Read(body)
	if err != nil {
		return nil, err
	}

	return body[:n], nil
}
