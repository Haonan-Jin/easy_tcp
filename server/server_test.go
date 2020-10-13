package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"testing"
)

func TestClient(t *testing.T) {
	content := []byte("啊")
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
