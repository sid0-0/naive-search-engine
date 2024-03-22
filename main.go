package main

import (
	"bufio"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

// var allTemplates *template.Template

var SeedPages = []string{"https://www.wikipedia.org/"}

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
			c.Locals("expiration", 0)
			c.Next()
			return nil
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		// log.Println(c.Locals("allowed"))
		// log.Println(c.Params("id"))
		// log.Println(c.Query("v"))
		// log.Println(c.Cookies("session"))
		log.Println(c.Cookies("expiration"))

		sub := NewSubscriber(c)
		subscribersList = append(subscribersList, sub)
		sub.RunMessageChannels()
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("base", nil)
	})

	log.Fatal(app.Listen("localhost:3000"))
}
