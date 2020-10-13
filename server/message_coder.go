package server

import (
	"bytes"
	"encoding/binary"
)

type Coder interface {
}

type Decoder interface {
	Coder
	// unpack bytes read from socket
	Unpack(data *bytesHolder) ([]byte, error)
	// decode unpacked bytes
	Decode(unpacked []byte) interface{}
}

type LengthFixedDecoder struct {
	bodyLen int
}

func NewLengthFixedDecoder(bodyLen int) *LengthFixedDecoder {
	decoder := new(LengthFixedDecoder)
	decoder.bodyLen = bodyLen
	return decoder
}

func (lf LengthFixedDecoder) Unpack(b *bytesHolder) ([]byte, error) {
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

func (lf LengthFixedDecoder) Decode(unpacked []byte) interface{} {
	return string(unpacked)
}
