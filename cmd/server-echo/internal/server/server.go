package server

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go-eventhub/cmd/server-echo/internal/listeners"
	"go-eventhub/cmd/server-echo/internal/mlistener"
	"go-eventhub/cmd/server-echo/internal/target"
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

	startMListener(&wsClients)

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

	// Login route
	e.POST("/login", login)

	// Web client
	e.Static("/", "./public")

	// Restricted group
	r := e.Group("/restricted")
	r.Use(middleware.JWT([]byte("secret")))
	r.GET("", restricted)
	r.GET("/ws", wsEndpoint)
}

func startMListener(wsClients *sync.Map) {
	ml := mlistener.New()

	ml.Add(&listeners.NATSListener{
		Targets: []target.Target{&target.WebSocketTarget{WSClients: wsClients}, &target.KafkaTarget{}},
	})

	ml.Start()
}
