package server

import (
	"log"
	"net"
)

type connHandler struct {
	conn    net.Conn
	holder  *bytesHolder
	handler Handler
}

func handleConnection(conn net.Conn, decoder Decoder, handler Handler) {
	connHandler := new(connHandler)

	holder := newByteBuffer(decoder)

	connHandler.handler = handler
	connHandler.holder = holder
	connHandler.conn = conn

	connHandler.start()
}

func (ch *connHandler) start() {
	terminateChan := make(chan error, 1)
	go ch.read(terminateChan)

	select {
	case err := <-terminateChan:
		ch.holder.buffer.Reset()
		log.Println(err)
		return
	}
}

func (ch *connHandler) read(terminateChan chan<- error) {

	buffer := make([]byte, 1024)

	for {

		n, err := ch.conn.Read(buffer)
		if err != nil {
			_ = ch.conn.Close()
			terminateChan <- err
			return
		}

		err = ch.holder.write(buffer[:n], ch.conn, ch.handler)
		if err != nil {
			log.Println(err)
			continue
		}

	}

}
