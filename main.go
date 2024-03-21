package main

import (
	"bufio"
	"fmt"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/google/uuid"
)

// var allTemplates *template.Template

var SeedPages = []string{"https://www.wikipedia.org/"}

type Subscriber struct {
	Id            string
	LogChannel    chan (string)
	ResultChannel chan (string)
}

func NewSubscriber() *Subscriber {
	return &Subscriber{
		Id:            uuid.NewString(),
		LogChannel:    make(chan string),
		ResultChannel: make(chan string),
	}
}

var subscribersList []*Subscriber

func startSending(s *Subscriber, w *bufio.Writer) {
	for {
		msg := fmt.Sprintf("Scraping %s", "wikipedia.org")
		fmt.Fprintf(w, "event:addLog\ndata:%s<br>%s<br>\n\n", msg, s.Id)
		fmt.Fprintf(w, "event:searchResult\ndata:<div><a href='%s' class='w-full'>%s</a></div>\n\n", "www.wikipedia.com", "wikipedia.com")

		err := w.Flush()
		if err != nil {
			// Refreshing page in web browser will establish a new
			// SSE connection, but only (the last) one is alive, so
			// dead connections must be closed here.
			fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)
			break
		}
		time.Sleep(time.Second)
	}
}

func main() {
	engine := html.New("./templates", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// app.Get("/sse", func(c *fiber.Ctx) error {
	// 	c.Set("Content-Type", "text/event-stream")
	// 	c.Set("Cache-Control", "no-cache")
	// 	c.Set("Connection", "keep-alive")
	// 	c.Set("Transfer-Encoding", "chunked")

	// 	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
	// 		sub := NewSubscriber()
	// 		subscribersList = append(subscribersList, sub)
	// 		startSending(sub, w)
	// 	}))

	// 	return nil
	// })

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			c.Next()
			return nil
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		log.Println(c.Locals("allowed"))
		log.Println(c.Params("id"))
		log.Println(c.Query("v"))
		log.Println(c.Cookies("session"))
		for {
			// run separate goroutines with channels for read/write
			if mt, msg, err := c.ReadMessage(); err == nil {
				spew.Print("Received", mt, string(msg))
				if string(msg) == "marco" {
					err := c.WriteMessage(websocket.TextMessage, []byte("polo"))
					if err != nil {
						spew.Printf("Failed to send ws message")
					}
				}
			}
		}
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("base", nil)
	})

	log.Fatal(app.Listen("localhost:3000"))
}
