package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/valyala/fasthttp"
)

// var allTemplates *template.Template

var SeedPages = []string{"https://www.wikipedia.org/"}

type Subscriber struct {
	stream *bufio.Writer
}

var subscribersList []*Subscriber

func startSending(s *Subscriber) error {
	fmt.Println(s)
	w := s.stream
	if w == nil {
		return errors.New("write stream missing")
	}
	for {
		msg := fmt.Sprintf("Scraping %s", "wikipedia.org")
		// fmt.Fprintf(w, "data: <h2>%s</h2>\n\n", msg)
		fmt.Fprintf(w, "event:addLog\ndata:%s<br>\n\n", msg)
		fmt.Fprintf(w, "event:searchResult\ndata:<div><a href='%s' class='w-full'>%s</a></div>\n\n", "www.wikipedia.com", "wikipedia.com")

		err := w.Flush()
		if w == nil {
			fmt.Println("write stream missing")
		}
		if err != nil {
			// Refreshing page in web browser will establish a new
			// SSE connection, but only (the last) one is alive, so
			// dead connections must be closed here.
			fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)

			break
		}
		time.Sleep(time.Second)
	}
	return nil
}

func main() {
	engine := html.New("./templates", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// allTemplates, _ = ParseAllTemplates("templates")
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("base", nil)
	})

	app.Get("/sse", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		var sub Subscriber
		c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
			sub = Subscriber{stream: w}
			subscribersList = append(subscribersList, &sub)
		}))
		_ = startSending(&sub)

		return nil
	})
	// app.Get("/search", func(c *fiber.Ctx) error {
	// 	term := c.Query("term")

	// })

	log.Fatal(app.Listen("localhost:3000"))
}
