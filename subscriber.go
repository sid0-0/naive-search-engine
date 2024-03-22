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
	// find a way to keep connection from auto closing!!!
	// till then only this is th way
	go sub.StartReading()
	go sub.StartWriting()
	for {
	}
}

func (sub *Subscriber) StartReading() {
	for {
		if mt, msg, err := sub.Connection.ReadMessage(); err == nil {
			spew.Print("Received", mt, string(msg))
			switch strings.TrimSpace(string(msg)) {
			case "marco":
				sub.WriteChannel <- "polo"
			default:
				sub.WriteChannel <- "did not understand"
			}
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
