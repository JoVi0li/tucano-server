package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/coder/websocket"
)

type Client struct {
	Conn    *websocket.Conn
	Offer   map[string]string `json:"offer"`
	Answer  map[string]string `json:"answer"`
	Partner *Client
}

var clients = make([]*Client, 0)

func removeClient(index int) {
	clients = append(clients[:index], clients[index+1:]...)
}

func main() {
	http.HandleFunc("/call", func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			OriginPatterns: []string{"*"},
		})
		if err != nil {
			log.Fatal(err)
		}

		ctx, cancel := context.WithCancel(r.Context())

		var client = &Client{Conn: c}
		var index int
		clients = append(clients, client)
		index = len(clients) - 1

		defer func() {
			client.Partner.Partner = nil
			removeClient(index)
			c.CloseNow()
			cancel()
		}()

		for {
			var data map[string]interface{}
			_, d, err := c.Read(ctx)
			if err != nil {
				if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
					websocket.CloseStatus(err) == websocket.StatusGoingAway {
					break
				}
				log.Printf("reading message error: %v", err)
				break
			}

			err = json.Unmarshal(d, &data)
			if err != nil {
				log.Println(string(d))
				log.Println(err)
				continue
			}

			switch data["type"] {
			case "offer":
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				rIndex := r.Intn(len(clients))
				var random = clients[rIndex]
				if random.Conn != c {
					log.Println(random.Conn)
					if err := random.Conn.Write(ctx, websocket.MessageText, d); err != nil {
						log.Println(err)
					}
					client.Partner = random
					random.Partner = client
				}
			case "answer":
				if client.Partner != nil {
					if err := client.Partner.Conn.Write(ctx, websocket.MessageText, d); err != nil {
						log.Println(err)
					}
				}
			case "candidate":
				if client.Partner != nil {
					if err := client.Partner.Conn.Write(ctx, websocket.MessageText, d); err != nil {
						log.Println(err)
					}
				}
			}
		}
	})

	certFile := "/app/cert/fullchain.pem"
	keyFile := "/app/cert/privkey.pem"

	err := http.ListenAndServeTLS(":443", certFile, keyFile, nil)
	if err != nil {
		log.Fatal(err)
	}

}
