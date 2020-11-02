package goland

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"testing"
)

func Decode(data []byte) (interface{}, error) {
	return string(data), nil
}
func Encode(msg interface{}) []byte {
	return []byte(msg.(string))
}

type ClientHandler struct {
	times int
}

func (h *ClientHandler) HandleMsg(ctx Context, msg interface{}) {
	h.times++
	fmt.Println(msg)
	ctx.Write("asd")
}

func (h *ClientHandler) HandleErr(ctx Context, err error) {
	log.Println("disconnected to server, error: ", err)
	ctx.ReConn()
	ctx.Write("reconn")
	ctx.Close()
}

var randomSentences = []string{"Bad days will pass", "Your dream is not dre", "the manner in which someone behaves toward or deals with someone or something.", "是啊是啊", "不是不是"}
var serverAddr = net.TCPAddr{
	IP:   net.ParseIP("0.0.0.0"),
	Port: 3333,
}

func TestMultiClient(t *testing.T) {
	maxClient := 1000

	for i := 0; i < maxClient; i++ {
		go func() {
			client, e := NewTcpClient(nil, &serverAddr)
			if e != nil {
				panic(e)
			}

			client.AddEncoder(Encode)
			client.AddDecoder(Decode)
			client.AddHandler(new(ClientHandler))
			client.Dial()

			for j := 0; j < 1000; j++ {
				client.Write(randomSentences[rand.Intn(len(randomSentences))])
			}
		}()
	}
}

func TestDial(t *testing.T) {

	client, e := NewTcpClient(nil, &serverAddr)
	if e != nil {
		panic(e)
	}

	client.AddEncoder(Encode)
	client.AddDecoder(Decode)
	client.AddHandler(new(ClientHandler))

	client.Dial()

	client.Write("asd")

	select {}
}
