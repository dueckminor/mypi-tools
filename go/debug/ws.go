package debug

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

type msg struct {
	Topic string `json:"topic"`
	Value any    `json:"value"`
}

type WS struct {
	io.Closer
	listeners chan []chan msg
}

func NewWS() *WS {
	result := &WS{}
	result.listeners = make(chan []chan msg, 1)
	result.listeners <- []chan msg{}
	return result
}

func (ws *WS) Publish(topic string, value any) {
	fmt.Printf("publish %v: %v\n", topic, value)
	listeners := <-ws.listeners
	for _, listener := range listeners {
		listener <- msg{topic, value}
	}
	ws.listeners <- listeners
}

func (ws *WS) disconnect(ch chan msg) {
	listeners := <-ws.listeners
	for i, e := range listeners {
		if e == ch {
			listeners = append(listeners[:i], listeners[i+1:]...)
			break
		}
	}
	ws.listeners <- listeners

}

func (ws *WS) handleConnection(conn *websocket.Conn) {
	ch := make(chan msg, 10)
	defer close(ch)

	ws.listeners <- append(<-ws.listeners, ch)

	for {
		m := <-ch
		fmt.Printf("publish-ws %v: %v\n", m.Topic, m.Value)
		err := websocket.JSON.Send(conn, m)
		if err != nil {
			break
		}
	}

	conn.Close()

	ws.disconnect(ch)
}

func (ws *WS) Run(r *gin.RouterGroup) {
	r.Any("ws", gin.WrapH(websocket.Handler(ws.handleConnection)))
	r.Any("ws/", gin.WrapH(websocket.Handler(ws.handleConnection)))
}
