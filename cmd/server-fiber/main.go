package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	jwtware "github.com/gofiber/jwt"
	"github.com/gofiber/websocket"
	"log"
	"time"
)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
	}))

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("secret"),
	}))

	// Login route
	app.Post("/login", login)

	// Unauthenticated route
	app.Get("/", accessible)

	// Restricted Routes
	app.Get("/restricted", restricted)

	// Optional middleware
	app.Use("/ws", func(c *fiber.Ctx) {
		if c.Get("host") == "localhost:3000" {
			c.Locals("Host", "Localhost:3000")
			c.Next()
			return
		}
		c.Status(403).Send("Request origin not allowed")
	})

	// Upgraded websocket request
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		fmt.Println(c.Locals("Host")) // "Localhost:3000"
		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)
			err = c.WriteMessage(mt, msg)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}))

	// ws://localhost:3000/ws
	app.Listen(3000)
}

func accessible(c *fiber.Ctx) {
	c.Send("Accessible")
}

func restricted(c *fiber.Ctx) {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	c.Send("Welcome " + name)
}

func login(c *fiber.Ctx) {
	user := c.FormValue("user")
	pass := c.FormValue("pass")

	// Throws Unauthorized error
	if user != "john" || pass != "doe" {
		c.SendStatus(fiber.StatusUnauthorized)
		return
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "John Doe"
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		c.SendStatus(fiber.StatusInternalServerError)
		return
	}

	c.JSON(fiber.Map{"token": t})
}
