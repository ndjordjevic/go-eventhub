package listeners

import (
	"github.com/nats-io/nats.go"
	"github.com/ndjordjevic/go-eventhub/cmd/server-echo/internal/pushers"
	"log"
	"os"
	"strings"
	"time"
)

type NATS struct {
	Targets []pushers.EventPusher
}

var natsAddr = os.Getenv("NATS_ADDR")

func (n *NATS) Listen() {
	log.Println("Connecting to NATS on:", natsAddr)
	nc, err := nats.Connect(natsAddr)
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

		log.Println("Msg received:", string(msg.Data))

		parsed := strings.Split(string(msg.Data), ",")

		for _, t := range n.Targets {
			t.Push(parsed)
		}
	}
}
