package server

import (
	"bytes"
	"encoding/binary"
)

type tooShortBytesError struct {
}

func (t *tooShortBytesError) Error() string {
	return "received bytes too short to parseReadBytes"
}

var tooShort = new(tooShortBytesError)

type BytesHolder struct {
	decoder    Decoder
	Buffer     *bytes.Buffer
	workStatus bool
	stopChan   chan int
}

func newByteBuffer(decoder Decoder) *BytesHolder {
	holder := new(BytesHolder)
	holder.Buffer = bytes.NewBuffer(nil)
	holder.decoder = decoder
	return holder
}

func (b *BytesHolder) decode() ([]byte, error) {
	if b.Buffer.Len() < 4 {
		return nil, tooShort
	}

	header := make([]byte, 4)
	_, err := b.Buffer.Read(header)
	if err != nil {
		return nil, err
	}

	bodyLen := int(binary.BigEndian.Uint32(header))

	if b.Buffer.Len() < bodyLen {
		recent := bytes.NewBuffer(nil)
		recent.Write(header)
		recent.Write(b.Buffer.Bytes())
		b.Buffer = recent
		return nil, tooShort
	}

	body := make([]byte, bodyLen)
	n, err := b.Buffer.Read(body)
	if err != nil {
		return nil, err
	}

	return body[:n], nil
}

//
//func (b *BytesHolder) write(data []byte, conn *ContextHandler, handler Handler) error {
//	_, err := b.Buffer.Write(data)
//	if err != nil {
//		return err
//	}
//
//	b.startUnpack(conn, handler)
//
//	return nil
//}
//
//func (b *BytesHolder) startUnpack(conn *ContextHandler, handler Handler) {
//	for {
//		decode, e := b.decoder.Unpack(b)
//		if e != nil {
//			if e != tooShort {
//				log.Println(e)
//			}
//			log.Println(e)
//			break
//		}
//		handler.Handle(conn, b.decoder.Decode(decode))
//	}
//}
