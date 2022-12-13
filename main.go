package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

func main() {
	// Create a new engine
	engine := html.New("./ui/html", ".html")

	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views:                engine,
		CompressedFileSuffix: ".compressed.gz",
	})

	//app.Static("ui/html/static", "ui/html/static")   //рабочая строка
	app.Static("ui/html/static", "ui/html/static", fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		CacheDuration: 1000 * time.Hour,
		MaxAge:        3600,
	})

	app.Get("/", MainPageHandler)

	app.Get("/layout", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		return c.Render("index2", fiber.Map{
			"Title": "Hello, World!",
		}, "layouts/main")
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).Render("index", fiber.Map{
			"Error": "Error 404. Not found!",
			"Title": "Error 404. Not found!",
		}, "errors/404")
	})

	log.Fatal(app.Listen(":3000"))
}
func MainPageHandler(c *fiber.Ctx) error {
	// Render index
	return c.Render("index", fiber.Map{
		"Title": "Hello, World!",
	})
}
