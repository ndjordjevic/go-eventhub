package meventsource

type EventSource interface {
	Listen()
}

type MultipleEventSource struct {
	listeners []EventSource
}

func New() *MultipleEventSource {
	return &MultipleEventSource{
		listeners: make([]EventSource, 0),
	}
}

func (ml *MultipleEventSource) Add(el EventSource) {
	ml.listeners = append(ml.listeners, el)
}

func (ml *MultipleEventSource) Start() {
	for _, l := range ml.listeners {
		go l.Listen()
	}
}
