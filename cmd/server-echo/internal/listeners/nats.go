package listeners

import (
	"github.com/nats-io/nats.go"
	"github.com/ndjordjevic/go-eventhub/cmd/server-echo/internal/pushers"
	"github.com/ndjordjevic/go-eventhub/internal/protogen/api"
	"google.golang.org/protobuf/proto"
	"log"
	"os"
	"time"
)

type NATS struct {
	Pushers  []pushers.EventPusher
	Subjects []string
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
	sub, err := nc.SubscribeSync("instrument")
	if err != nil {
		log.Fatal(err)
	}

	for {
		// Wait for a message
		msg, err := sub.NextMsg(10 * time.Second)
		if err != nil {
			continue
		}

		instrument := &api.Instrument{}
		if err := proto.Unmarshal(msg.Data, instrument); err != nil {
			log.Fatalln("Failed to parse instrument:", err)
		}

		log.Printf("Msg received from nats, user: %v", instrument)

		for _, p := range n.Pushers {
			p.Push(instrument)
		}
	}
}
