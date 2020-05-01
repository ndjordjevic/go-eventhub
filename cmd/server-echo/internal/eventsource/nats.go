package eventsource

import (
	"github.com/nats-io/nats.go"
	"go-eventhub/cmd/server-echo/internal/pushers"
	"log"
	"strings"
	"time"
)

type NATS struct {
	Targets []pushers.EventPusher
}

func (n *NATS) Listen() {
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

		parsed := strings.Split(string(msg.Data), ",")

		for _, t := range n.Targets {
			t.Push(parsed)
		}
	}
}
