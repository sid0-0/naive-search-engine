package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

// var allTemplates *template.Template

var SeedPages = []string{"https://www.wikipedia.org/"}

func main() {
	engine := html.New("./templates", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// allTemplates, _ = ParseAllTemplates("templates")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("base", nil)
	})

	// app.Get("/search", func(c *fiber.Ctx) error {
	// 	term := c.Query("term")

	// })

	log.Fatal(app.Listen("localhost:3000"))
}
