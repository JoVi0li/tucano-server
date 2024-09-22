package internal

import (
	"math/rand"
	"time"

	"github.com/coder/websocket"
)

type client struct {
	Conn    *websocket.Conn
	Offer   map[string]string `json:"offer"`
	Answer  map[string]string `json:"answer"`
	Partner *client
}

type clientList struct {
	clients []*client
}

func newClient(c *websocket.Conn) *client {
	return &client{Conn: c}
}

func NewClientList() *clientList {
	return &clientList{}
}

func (c* clientList) removeClient(index int) {
	if len(c.clients) == 1 {
		c.clients = make([]*client, 0)
	}
	c.clients = append(c.clients[:index], c.clients[index+1:]...)
}

func (c* clientList) appendClient(client *client) int {
	c.clients = append(c.clients, client)
	return len(c.clients) - 1
}

func (c* clientList) drawClient() *client {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := r.Intn(len(c.clients))
	random := c.clients[index]
	return random
}
