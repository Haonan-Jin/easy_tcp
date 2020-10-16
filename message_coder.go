package goland

// encode your data in this func
// you need only encode the message part
// this frame work will add the header
// witch describes your message bytes len
type Encoder func(message interface{}) []byte

// decode unpacked bytes to a interface{}
type Decoder func(unpacked []byte) (interface{}, error)
