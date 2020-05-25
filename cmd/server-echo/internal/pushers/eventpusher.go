package pushers

import "github.com/ndjordjevic/go-eventhub/internal/protogen/api"

type EventPusher interface {
	Push(*api.Instrument)
}
