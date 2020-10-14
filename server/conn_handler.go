package server

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type ContextHandler struct {
	mutex sync.Mutex

	conn net.Conn

	buffer  *bytes.Buffer
	decoder Decoder

	dataChan chan int
	msgChan  chan interface{}

	handler Handler
	encoder Encoder
}

func handleConnection(conn net.Conn, encoder Encoder, decoder Decoder, handler Handler) {
	contextHandler := new(ContextHandler)

	contextHandler.msgChan = make(chan interface{}, 100)

	contextHandler.dataChan = make(chan int, 100)
	contextHandler.buffer = bytes.NewBuffer(nil)
	contextHandler.handler = handler
	contextHandler.conn = conn
	contextHandler.decoder = decoder
	contextHandler.encoder = encoder

	// TODO open a goroutine to parseReadBytes read bytes

	contextHandler.start()
}

func (ch *ContextHandler) start() {
	terminateChan := make(chan error, 1)
	go ch.read(terminateChan)
	go ch.parseReadBytes()
	go ch.handleMsg()

	select {
	case err := <-terminateChan:
		log.Println(err)
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

		ch.dataChan <- 1
	}

}

func (ch *ContextHandler) parseReadBytes() {
	var timeoutChan *time.Timer
loop:
	for {
		timeoutChan = time.NewTimer(time.Second * 30)
		select {
		case <-ch.dataChan:
			for {
				ch.mutex.Lock()
				msg, e := ch.decoder.Decode(ch.buffer)

				if e != nil {
					if msg != nil {
						ch.buffer = msg.(*bytes.Buffer)
					}
					timeoutChan.Stop()
					ch.mutex.Unlock()
					break
				}

				ch.mutex.Unlock()
				ch.msgChan <- msg
			}
		case <-timeoutChan.C:
			ch.buffer.Reset()
			break loop
		}
	}
}

func (ch *ContextHandler) handleMsg() {
	var timeoutChan *time.Timer
loop:
	for {
		timeoutChan = time.NewTimer(time.Second * 30)
		select {
		case msg := <-ch.msgChan:
			timeoutChan.Stop()
			ch.handler(ch, msg)
		case <-timeoutChan.C:
			break loop
		}
	}
}

func (ch *ContextHandler) Write(msg interface{}) {
	encode := ch.encoder.Encode(msg)
	_, err := ch.conn.Write(encode)
	if err != nil {
		_ = ch.conn.Close()
		fmt.Println(err)
	}
}
