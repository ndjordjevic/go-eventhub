package server

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"sync"
)

var (
	syncMap  sync.Map
	socketMu sync.Mutex
)

func Run() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Login route
	e.POST("/login", login)

	// Web client
	e.Static("/", "./public")

	// Restricted group
	r := e.Group("/restricted")
	r.Use(middleware.JWT([]byte("secret")))
	r.GET("", restricted)
	r.GET("/ws", sbevents)

	go natsListener()

	e.Logger.Fatal(e.Start(":8080"))
}
