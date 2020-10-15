package client

import (
	"fmt"
	"github.com/Haonan-Jin/tcp_server/handler"
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

type Handler struct {
	times int
}

func (h *Handler) Handle(ctx handler.ContextHandler, msg interface{}) {
	//fmt.Println("response from server: ", msg)
	h.times++
	fmt.Println(h.times)
	if h.times == 1000 {
		ctx.Close()
	}
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
			client.AddHandler(new(Handler))
			client.Dial()

			for j := 0; j < 1000; j++ {
				client.Write(randomSentences[rand.Intn(len(randomSentences))])
			}
		}()
	}

	select {}

}

func TestDial(t *testing.T) {

	client, e := NewTcpClient(nil, &serverAddr)
	if e != nil {
		panic(e)
	}

	client.AddEncoder(Encode)
	client.AddDecoder(Decode)
	client.AddHandler(new(Handler))

	client.Dial()

	client.Write("hello")

	select {}
}
