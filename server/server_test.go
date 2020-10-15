package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/Haonan-Jin/tcp_server/handler"
	"math/rand"
	"net"
	"sync"
	"testing"
)

// decode your data in this func
func Decode(b []byte) (interface{}, error) {
	return string(b), nil
}

// encode your data in this func
func Encode(msg interface{}) []byte {
	return []byte(msg.(string))
}

type StringHandler struct {
	mutex sync.Mutex
	times int
}

// process decoded message
func (t *StringHandler) Handle(ctx handler.ContextHandler, msg interface{}) {
	t.mutex.Lock()
	t.times++
	t.mutex.Unlock()

	ctx.Write("1")
	fmt.Println(t.times)
	fmt.Println("read from client: ", msg)
}

func TestServer(t *testing.T) {
	addr := net.TCPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 3333,
	}

	tcpServer, e := NewTcpServer(&addr)
	if e != nil {
		panic(e)
	}

	tcpServer.AddEncoder(Encode)
	tcpServer.AddDecoder(Decode)
	tcpServer.AddHandler(new(StringHandler))
	tcpServer.Start()
}

func TestClient(t *testing.T) {
	content := []byte("a")
	counter := 0

	for j := 0; j < 1000; j++ {
		go func() {
			buffer := bytes.NewBuffer(nil)
			conn, e := net.DialTCP("tcp", nil, &net.TCPAddr{
				IP:   net.ParseIP("0.0.0.0"),
				Port: 3333,
			})
			if e != nil {
				panic(e)
			}

			go func() {
				bff := make([]byte, 1024)
				for {
					i, err := conn.Read(bff)
					if err != nil {
						continue
					}
					fmt.Println(bff[:i])
				}
			}()

			for i := 0; i < 1000; i++ {
				header := make([]byte, 4)
				intn := rand.Intn(9)
				if intn == 0 {
					intn = 1
				}
				binary.BigEndian.PutUint32(header, uint32(len(content)*intn))
				buffer.Write(header)
				for i := 0; i < intn; i++ {
					buffer.Write(content)
				}

				content := buffer.Bytes()
				_, e = conn.Write(content)
				if e != nil {
					fmt.Println(e)
				}
				counter++
				buffer.Reset()
			}
			conn.Close()
		}()
	}

	select {}
}
