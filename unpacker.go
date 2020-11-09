package goland

import (
	"bytes"
	"encoding/binary"
	"errors"
)

var ReadLengthError = errors.New("not read completely")

// Define a packer that try to solve bytes have read.
// Always return your service or unsolved bytes.
// Return an error such as ReadLengthError when you are returning unsolved bytes.
type UnPacker func(b *bytes.Buffer) ([]byte, error)

// Try to unpack read bytes.
// Returned bytes without header bytes.
// The message that be sent to server must be encoded like:
// +------+-----------------------------------+
// | 0-3  | header describes data length      |
// | 4... | data body                         |
// +------+-----------------------------------+
// for example:
// the data part []byte{4,2}, you have to add the header
// []byte{0,0,0,2} to describes the length of data part.
// If received []byte{0,0,0,2,4,3} ,
// Will returns []byte{4,3}. the header will be dropped.
func LengthFixedUnpack(b *bytes.Buffer) ([]byte, error) {
	if b.Len() < 4 {
		return nil, ReadLengthError
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
		return recent.Bytes(), ReadLengthError
	}

	body := make([]byte, bodyLen)
	n, err := b.Read(body)
	if err != nil {
		return nil, err
	}

	return body[:n], nil
}
