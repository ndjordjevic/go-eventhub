package server

import (
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
	"log"
	"strings"
	"sync"
	"time"
)

type NATSListener struct {
	wsClients *sync.Map
}

func (n *NATSListener) Listen() {
	nc, err := nats.Connect("localhost")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Subscribe
	sub, err := nc.SubscribeSync("sb-events")
	if err != nil {
		log.Fatal(err)
	}

	for {
		// Wait for a message
		msg, err := sub.NextMsg(10 * time.Second)
		if err != nil {
			continue
		}

		s := strings.Split(string(msg.Data), ",")

		client, ok := n.wsClients.Load(s[0])
		if ok {
			log.Printf("Forwarding to ws: %v", s)
			client.(*websocket.Conn).WriteMessage(websocket.TextMessage, []byte(s[1]))
		}
	}
}
