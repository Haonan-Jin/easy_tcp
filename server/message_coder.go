package server

type Encoder func(msg interface{}) []byte

type Decoder func(unpacked []byte) (interface{}, error)
