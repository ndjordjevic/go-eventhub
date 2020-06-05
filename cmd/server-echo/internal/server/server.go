package server

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/ndjordjevic/go-eventhub/cmd/server-echo/internal/listeners"
	"github.com/ndjordjevic/go-eventhub/cmd/server-echo/internal/pushers"
	"net/http"
	"sync"
)

type CustomContext struct {
	echo.Context
	wsClients *sync.Map
}

func Run() {
	e := echo.New()

	var wsClients sync.Map
	setupEchoServer(e, &wsClients)

	startMultiEventListener(&wsClients)

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

func startMultiEventListener(wsClients *sync.Map) {
	ml := listeners.NewMultiple()

	ml.Add(&listeners.NATS{
		Pushers:  []pushers.EventPusher{&pushers.WebSocket{WSClients: wsClients}, &pushers.Kafka{}},
		Subjects: []string{"instrument", "order"},
	})

	ml.Start()
}
