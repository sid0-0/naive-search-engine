package main

import (
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type Subscriber struct {
	Id           string
	WriteChannel chan (string)
	ReadChannel  chan (string)
	Connection   *websocket.Conn
}

func NewSubscriber(c *websocket.Conn) *Subscriber {
	return &Subscriber{
		Id:           uuid.NewString(),
		WriteChannel: make(chan string),
		ReadChannel:  make(chan string),
		Connection:   c,
	}
}

func (sub *Subscriber) RunMessageChannels() {
	go sub.StartWriting()
	// run an infinite reading loop
	// so that the connection does not close
	// idk if that's ideal but whatever
	sub.StartReadLoop()
}

func (sub *Subscriber) StartReadLoop() {
	for {
		mt, msg, err := sub.Connection.ReadMessage()
		if err != nil {
			spew.Print("Connection dropped because read failed")
			return
		}
		spew.Print("Received", mt, string(msg))
		switch strings.TrimSpace(string(msg)) {
		case "marco":
			sub.WriteChannel <- "polo"
		default:
			sub.WriteChannel <- "did not understand"
		}
	}
}

func (sub *Subscriber) StartWriting() {
	for {
		data := <-sub.WriteChannel
		fmt.Println("message: ", data)
		err := sub.Connection.WriteMessage(websocket.TextMessage, []byte(data))
		if err != nil {
			spew.Print("Failed to send ws message", err)
		}
	}
}
