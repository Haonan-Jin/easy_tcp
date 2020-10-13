package server

import "bytes"

type Encoder interface {
	Encode(msg interface{}) []byte
}

type Decoder interface {
	// parseReadBytes unpacked bytes
	Decode(unpacked *bytes.Buffer) (interface{}, error)
}
