package goland

// Encode your msg to bytes by your protocol
type Encoder func(message interface{}) []byte

// Decode unpacked bytes to any type you want.
type Decoder func(unpacked []byte) (interface{}, error)
