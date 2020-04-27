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
		msg, err := sub.NextMsg(100 * time.Second)
		if err != nil {
			log.Fatal(err)
		}

		s := strings.Split(string(msg.Data), ",")

		// Use the response
		log.Printf("Forwarding to ws: %v", s)

		usrConns[s[0]].WriteMessage(websocket.TextMessage, []byte(s[1]))
	}
}
