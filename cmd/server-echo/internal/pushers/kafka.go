package pushers

import "github.com/ndjordjevic/go-eventhub/internal/protogen/api"

type Kafka struct {
}

func (k Kafka) Push(*api.Instrument) {
	//panic("implement me")
}
