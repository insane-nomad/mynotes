package main

import (
	"html/template"
	"log"
	"mynotes/database"
	"mynotes/internal/routes"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/template/html"
)

func main() {
	// Create a new database
	dbErr := (&database.SqlHandler{}).InitDatabase()
	if dbErr != nil {
		panic(dbErr)
	}

	// Create a new engine
	engine := html.New("./ui/html", ".html")
	engine.AddFunc(
		// add unescape function
		"unescape", func(s string) template.HTML {
			return template.HTML(s)
		},
	)

	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views:                engine,
		CompressedFileSuffix: ".compressed.gz",
	})

	app.Static("ui/html/static", "ui/html/static", fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		CacheDuration: 1000 * time.Hour,
		MaxAge:        3600,
	})

	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format:       "(${time}) - [${ip}]:${port} ${status} - ${method} ${referer} -> ${path}\n",
		EnableColors: true,
	}))

	// Routes
	app.Post("/addnote", routes.AddnoteHandler)
	app.Get("/", routes.MainPageHandler)
	app.Get("/layout", routes.LayoutHandler)
	app.Get("/add", routes.AddHandler)
	app.Get("/pages/:id<range(1,1000)>", routes.PaginationHandler)
	app.Get("/delnote/:id<range(1,10000)>", routes.DelNoteHandler)

	app.Use(routes.Return404Handler)

	// start server on 127.0.0.1:3000
	log.Fatal(app.Listen(":3000"))

}
