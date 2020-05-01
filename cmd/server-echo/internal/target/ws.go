package target

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type WebSocketTarget struct {
	WSClients *sync.Map
}

func (w WebSocketTarget) Push(msg []string) {
	client, ok := w.WSClients.Load(msg[0])
	if ok {
		log.Printf("Forwarding to ws: %v", msg)
		client.(*websocket.Conn).WriteMessage(websocket.TextMessage, []byte(msg[1]))
	}
}
