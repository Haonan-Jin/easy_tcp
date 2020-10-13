package server

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
)

type tooShortBytesError struct {
}

func (t *tooShortBytesError) Error() string {
	return "received bytes too short to decode"
}

var tooShort = new(tooShortBytesError)

type bytesHolder struct {
	coder      Decoder
	buffer     *bytes.Buffer
	workStatus bool
	stopChan   chan int
}

func newByteBuffer(decoder Decoder) *bytesHolder {
	holder := new(bytesHolder)
	holder.buffer = bytes.NewBuffer(nil)
	holder.coder = decoder
	return holder
}

func (b *bytesHolder) write(data []byte, conn net.Conn, handler Handler) error {
	_, err := b.buffer.Write(data)
	if err != nil {
		return err
	}

	b.startUnpack(conn, handler)

	return nil
}

func (b *bytesHolder) startUnpack(conn net.Conn, handler Handler) {
	for {
		decode, e := b.coder.Unpack(b)
		if e != nil {
			if e != tooShort {
				log.Println(e)
			}
			log.Println(e)
			break
		}
		handler.Handle(conn, b.coder.Decode(decode))
	}
}

func (b *bytesHolder) decode() ([]byte, error) {
	if b.buffer.Len() < 4 {
		return nil, tooShort
	}

	header := make([]byte, 4)
	_, err := b.buffer.Read(header)
	if err != nil {
		return nil, err
	}

	bodyLen := int(binary.BigEndian.Uint32(header))

	if b.buffer.Len() < bodyLen {
		recent := bytes.NewBuffer(nil)
		recent.Write(header)
		recent.Write(b.buffer.Bytes())
		b.buffer = recent
		return nil, tooShort
	}

	body := make([]byte, bodyLen)
	n, err := b.buffer.Read(body)
	if err != nil {
		return nil, err
	}

	return body[:n], nil
}
