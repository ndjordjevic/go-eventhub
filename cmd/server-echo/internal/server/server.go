package server

import (
	"context"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/ndjordjevic/go-eventhub/cmd/server-echo/internal/listeners"
	"github.com/ndjordjevic/go-eventhub/cmd/server-echo/internal/pushers"
	"go.etcd.io/etcd/clientv3"
	"log"
	"net/http"
	"sync"
	"time"
)

type CustomContext struct {
	echo.Context
	wsClients *sync.Map
}

func Run() {
	etcd, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 2 * time.Second,
	})

	if err == context.DeadlineExceeded {
		log.Fatal(err)
	}

	e := echo.New()

	var wsClients sync.Map
	setupEchoServer(e, &wsClients)

	startMultiEventListener(&wsClients, etcd)

	e.Logger.Fatal(e.Start(":8080"))
}

func setupEchoServer(e *echo.Echo, wsClients *sync.Map) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{c, wsClients}
			return next(cc)
		}
	})

	e.GET("/debug/pprof/*", echo.WrapHandler(http.DefaultServeMux))

	// Login route
	e.POST("/login", login)

	// Web client
	e.Static("/", "./public")

	// Restricted group
	r := e.Group("/restricted")
	r.Use(middleware.JWT([]byte("secret")))
	r.GET("/ws", wsEndpoint)
}

func startMultiEventListener(wsClients *sync.Map, etcd *clientv3.Client) {
	ml := listeners.NewMultiple()

	ml.Add(&listeners.NATS{
		Pushers:  []pushers.EventPusher{&pushers.WebSocket{WSClients: wsClients}, &pushers.Kafka{}},
		Subjects: []string{"instrument", "order"},
		Config:   etcd,
	})

	ml.Start()
}
