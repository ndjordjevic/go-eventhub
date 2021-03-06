package tests

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
	"net/url"
	"sync"
	"testing"
)

const (
	addr string = `http://localhost:8080`
	user string = "ndjordjevic"
)

var ready = make(chan interface{})

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

	<-ready
	_ = nc.Publish(subj, msg)

	_ = nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Published [%s] : '%s'\n", subj, msg)
	}
}

func wsListener(t *testing.T, wg *sync.WaitGroup) {
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
			ready <- nil
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
			if string(message) != "new_order" {
				t.Error("No right msg")
			}
			wg.Done()
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
		}
	}
}
