package main

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

var clients []*Client

type Client struct {
	conn *websocket.Conn
	send chan *Message

	roomId string
	user *User
}

const messageBufferSize = 256

func newClient(conn *websocket.Conn, roomId string, u *User) {
	c := &Client{
		conn: conn,
		send: make(chan *Message, messageBufferSize),
		roomId: roomId,
		user : u,
	}

	clients = append(clients, c)

	go c.readLoop()
	go c.writeLoop()
}

func (c *Client) Close() {
	for i, client := range clients {
		if client == c {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
	close(c.send)

	c.conn.Close()
	log.Printf("close connectiuon. addr : %s", c.conn.RemoteAddr())
}

func (c *Client) readLoop() {
	for {
		m, err := c.read()
		if err != nil {
			log.Println("read message error : ", err )
			break
		}
		m.create()
		broadcast(m)
	}
	c.Close()
}

func (c *Client) writeLoop()  {
	for msg := range c.send {
		if c.roomId == msg.RoomId.Hex() {
			c.write(msg)
		}
	}
}

func broadcast(m *Message)  {
	for _, client := range clients {
		client.send <- m
	}
}

func (c *Client) read() (*Message, error) {
	var msg *Message

	if err := c.conn.ReadJSON(&msg); err != nil {
		return nil, err
	}

	msg.CreatedAt =  time.Now()
	msg.User = c.user

	log.Println("read from websocket : ", msg)

	return msg, nil
}

func (c *Client) write(m *Message) error {
	log.Println("write to websoket:", m)

	return c.conn.WriteJSON(m)
}