package server

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	usrConns = make(map[string]*websocket.Conn)
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
	r.GET("/ws", hello)

	go natsListener()

	e.Logger.Fatal(e.Start(":8080"))
}
