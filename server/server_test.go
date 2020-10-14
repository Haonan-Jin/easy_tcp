package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"
)

func TestWG(t *testing.T) {
}

func TestMultiSelect(t *testing.T) {
	ints1 := make(chan int, 1)
	timer := time.NewTimer(time.Second * 2)

	go func() {
		ints1 <- 1
		timer.Stop()
	}()

	select {
	case <-ints1:
		fmt.Println(1)
	case <-timer.C:
		fmt.Println(2)
	}
}

func TestSelect(t *testing.T) {
	ints := make(chan int, 1)

	go func() {
		ints <- 1
	}()

	select {
	case s := <-ints:
		fmt.Println(s)
	}
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
