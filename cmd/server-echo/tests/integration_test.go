package tests

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"testing"
)

const (
	addr string = `http://localhost:8080`
	user string = "ndjordjevic"
)

func TestIntegration(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	go wsListener(t, &wg)

	sendToNats()

	wg.Wait()
}

func sendToNats() {
	opts := []nats.Option{nats.Name("NATS Test Publisher")}

	// Connect to NATS
	nc, err := nats.Connect("localhost", opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	subj, msg := "sb-events", []byte("ndjordjevic,new_order")

	_ = nc.Publish(subj, msg)

	_ = nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Published [%s] : '%s'\n", subj, msg)
	}
}

func wsListener(t *testing.T, wg *sync.WaitGroup) {
	defer wg.Done()

	formData := url.Values{
		"username": {user},
		"password": {"test"},
	}

	r, err := http.PostForm(addr+"/login", formData)
	var result map[string]interface{}
	_ = json.NewDecoder(r.Body).Decode(&result)

	token, ok := result["token"]
	if ok == false {
		log.Fatal("User is not authorized to connect to ws")
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/restricted/ws"}
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
			fmt.Println("Listen ws")
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
			if string(message) == "new_order" {
				t.Log("Successful")
				return
			}
		}
	}()

	for {
		select {
		case <-done:
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
			}
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			//select {
			//case <-done:
			//case <-time.After(time.Second):
			//}
			wg.Done()
			return
		}
	}
}
