package goland

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"
	"testing"
)

// decode your data in this func
func StringDecoder(b []byte) (interface{}, error) {
	fmt.Println(string(b))
	return string(b), nil
}

// encode your data in this func
func StringEncoder(msg interface{}) []byte {
	msgBody := []byte(msg.(string))
	msgLen := make([]byte, 4)
	binary.BigEndian.PutUint32(msgLen, uint32(len(msgBody)))
	buffer := bytes.NewBuffer(msgLen)
	buffer.Write(msgBody)

	return buffer.Bytes()
}

func HttpEncoder(msg interface{}) []byte {
	return []byte(msg.(string))
}

type HttpHandler struct {
}

func (http *HttpHandler) HandleMsg(ctx Context, msg interface{}) {
	response := "HTTP/1.1 200 OK\n" +
		"Date: Tue, 10 Nov 2020 03:22:58 GMT\n" +
		"\n"
	ctx.Write(response)
	fmt.Println("write end")
	ctx.Close()
}

func (http *HttpHandler) HandleErr(ctx Context, err error) {
	log.Println("a connect closed, error: ", err)
}

type StringHandler struct {
	mutex sync.Mutex
	times int
}

// process decoded message
func (t *StringHandler) HandleMsg(ctx Context, msg interface{}) {
	t.mutex.Lock()
	t.times++
	t.mutex.Unlock()
	ctx.Write("reconn")
}

func (t *StringHandler) HandleErr(ctx Context, err error) {
	log.Println("a connect closed, error: ", err)
	ctx.Close()
}

func TestServer(t *testing.T) {
	listener, e := net.Listen("tcp", ":3333")
	if e != nil {
		panic(e)
	}

	for {
		conn, e := listener.Accept()
		if e != nil {
			continue
		}
		handler := NewConnectionHandler(conn)
		handler.AddEncoder(HttpEncoder)
		handler.AddDecoder(StringDecoder)
		handler.AddHandler(new(HttpHandler))
		handler.AddUnPacker(HttpUnPacker)
		handler.Serve()
	}
}

func HttpUnPacker(buffer *bytes.Buffer) ([]byte, error) {
	if buffer.Len() == 0 {
		return nil, ReadLengthError
	}
	i := make([]byte, 1024)
	buffer.Read(i)
	return i, nil
}
