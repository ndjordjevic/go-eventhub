package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var user = flag.String("user", "ndjordjevic", "username to login to ws with")

func main() {
	flag.Parse()
	log.SetFlags(0)

	formData := url.Values{
		"username": {*user},
		"password": {"test"},
	}

	r, err := http.PostForm("http://"+*addr+"/login", formData)
	if err != nil {
		log.Fatal(err)
	}
	var result map[string]interface{}
	_ = json.NewDecoder(r.Body).Decode(&result)

	token, ok := result["token"]
	if ok == false {
		log.Fatal("User is not authorized to connect to ws")
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/restricted/ws"}
	log.Printf("connecting to %s", u.String())

	var reqH http.Header
	reqH = make(map[string][]string)
	reqH.Add("Authorization", "Bearer "+token.(string))

	c, _, err := websocket.DefaultDialer.Dial(u.String(), reqH)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
			log.Printf("Echoing the same msg back to server: %s", message)
			if err := c.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}()

	//ticker := time.NewTicker(time.Second)
	//defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		//case t := <-ticker.C:
		//	err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
		//	if err != nil {
		//		log.Println("write:", err)
		//		return
		//	}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
