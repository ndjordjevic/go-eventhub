package server

import (
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
	"log"
	"strings"
	"time"
)

func natsListener() {
	nc, err := nats.Connect("localhost")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	for {
		// Subscribe
		sub, err := nc.SubscribeSync("sb-events")
		if err != nil {
			log.Fatal(err)
		}

		// Wait for a message
		msg, err := sub.NextMsg(10 * time.Second)
		if err != nil {
			continue
		}

		s := strings.Split(string(msg.Data), ",")

		// Use the response
		log.Printf("Forwarding to ws: %v", s)

		client, _ := syncMap.Load(s[0])
		socketMu.Lock()
		client.(*websocket.Conn).WriteMessage(websocket.TextMessage, []byte(s[1]))
		socketMu.Unlock()
	}
}
