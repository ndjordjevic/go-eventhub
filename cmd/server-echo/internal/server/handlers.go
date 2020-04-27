package server

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Throws unauthorized error
	if (username != "ndjordjevic" && password != "test") || (username != "vpopovic" && password != "test") {
		return echo.ErrUnauthorized
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = username
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func sbevents(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	c.Logger().Print("User logged in: ", name)

	syncMap.Store(name, ws)

	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		// Write
		//err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		//if err != nil {
		//	c.Logger().Error(err)
		//}

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("Client disconnected")
			c.Logger().Error(err)

			syncMap.Delete(name)
			return nil
		}
		fmt.Printf("%s\n", msg)
	}
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.String(http.StatusOK, "Welcome "+name+"!")
}
