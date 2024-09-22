package internal

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/coder/websocket"
)

func HandleSignalling(w http.ResponseWriter, r *http.Request, clients *clientList) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{OriginPatterns: []string{"*"}})
	if err != nil {
		log.Fatal(err)
	}

	ctx := r.Context()

	client := newClient(conn)
	cIndex := clients.appendClient(client)

	for {
		var data map[string]interface{}

		_, result, err := conn.Read(ctx)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
				websocket.CloseStatus(err) == websocket.StatusGoingAway {
				log.Printf("connection closed: %v", err)
				client.Partner.Partner = nil
				clients.removeClient(cIndex)
				conn.CloseNow()
				break
			}
			log.Printf("reading message error: %v", err)
			break
		}

		err = json.Unmarshal(result, &data)
		if err != nil {
			log.Printf("json unmarshal error: %v", err)
			continue
		}

		switch data["type"] {
		case "offer":
			for len(clients.clients) >= 2 {
				random := clients.drawClient()
				if random.Conn != conn {
					if err := random.Conn.Write(ctx, websocket.MessageText, result); err != nil {
						log.Printf("send offer to random client error: %v", err)
					}
					client.Partner = random
					random.Partner = client
					break
				} else {
					log.Printf("draw client itself, do it again")
					continue
				}
			}
		case "answer":
			if client.Partner != nil {
				if err := client.Partner.Conn.Write(ctx, websocket.MessageText, result); err != nil {
					log.Printf("send answer to partner error: %v", err)
				}
			}
		case "candidate":
			log.Printf("candidate")
			if client.Partner != nil {
				if err := client.Partner.Conn.Write(ctx, websocket.MessageText, result); err != nil {
					log.Printf("send ICE candidate to partner error: %v", err)
				}
			}
		}
	}
}
