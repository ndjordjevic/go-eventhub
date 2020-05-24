package main

import (
	"github.com/ndjordjevic/go-eventhub/cmd/server-echo/internal/server"
	_ "net/http/pprof"
)

func main() {
	server.Run()
}
