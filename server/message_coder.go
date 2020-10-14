package server

type Encoder func(msg interface{}) []byte

// decode unpacked bytes to a interface{}
// the message send to server must be encoded like:
// +------+----------------+
// | 0-3 | data length     |
// | 4.. | data body       |
// +------+----------------+
// for example:
// you send a string "hello"
// you should encode the message "hello" to
// 0 0 0 5 []byte("hello")
type Decoder func(unpacked []byte) (interface{}, error)
