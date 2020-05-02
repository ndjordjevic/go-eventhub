package listeners

type EventListener interface {
	Listen()
}

type MultipleEventListener struct {
	listeners []EventListener
}

func NewMultiple() *MultipleEventListener {
	return &MultipleEventListener{
		listeners: make([]EventListener, 0),
	}
}

func (ml *MultipleEventListener) Add(el EventListener) {
	ml.listeners = append(ml.listeners, el)
}

func (ml *MultipleEventListener) Start() {
	for _, l := range ml.listeners {
		go l.Listen()
	}
}
