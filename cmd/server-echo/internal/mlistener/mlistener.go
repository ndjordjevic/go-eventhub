package mlistener

type EventListener interface {
	Listen()
}

type MultipleListener struct {
	listeners []EventListener
}

func New() *MultipleListener {
	return &MultipleListener{
		listeners: make([]EventListener, 0),
	}
}

func (ml *MultipleListener) Add(el EventListener) {
	ml.listeners = append(ml.listeners, el)
}

func (ml *MultipleListener) Start() {
	for _, l := range ml.listeners {
		go l.Listen()
	}
}
