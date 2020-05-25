package pushers

import (
	"github.com/gorilla/websocket"
	"github.com/ndjordjevic/go-eventhub/internal/protogen/api"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"sync"
)

type WebSocket struct {
	WSClients *sync.Map
}

func (w WebSocket) Push(i *api.Instrument) {
	client, ok := w.WSClients.Load(i.User)
	if ok {
		json, err := protojson.Marshal(i)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Forwarding to ws: %s", json)
		if err := client.(*websocket.Conn).WriteMessage(websocket.TextMessage, json); err != nil {
			log.Fatal(err)
		}
	}
}
