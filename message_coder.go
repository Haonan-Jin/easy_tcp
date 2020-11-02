package goland

// Encode your msg to bytes
// You need only encode the message part,
// this framework will add the header automatically
// witch describes your message bytes len.
type Encoder func(message interface{}) []byte

// Decode unpacked bytes to any type you want.
type Decoder func(unpacked []byte) (interface{}, error)
