package pushers

type EventPusher interface {
	Push([]string)
}
