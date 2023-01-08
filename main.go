package main

import (
	"log"
	"mynotes/database"
	"mynotes/internal/routes"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/template/html"
)

func main() {
	// Create a new database
	database.InitDatabase()

	// Create a new engine
	engine := html.New("./ui/html", ".html")

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
		Format: "(${time}) - [${ip}]:${port} ${status} - ${method} ${referer} -> ${path}\n",
	}))

	// Routes
	app.Post("/savenote", routes.SavenoteHandler)
	app.Post("/adduser", routes.AdduserHandler)
	app.Post("/login", routes.LoginHandler)
	app.Get("/register", routes.RegisterHandler)
	app.Get("/logout", routes.LogoutHandler)
	app.Get("/", routes.MainPageHandler)
	app.Get("/add", routes.AddnoteHandler)
	app.Get("/pages/:id<range(1,1000)>", routes.PaginationHandler)
	app.Get("/delnote/:id<range(1,10000)>", routes.DelNoteHandler)
	app.Get("/metrics", monitor.New(monitor.Config{
		Title:      "Fiber Monitor",
		Refresh:    1 * time.Second,
		APIOnly:    false,
		Next:       nil,
		CustomHead: "",
		FontURL:    "https://fonts.googleapis.com/css2?family=Roboto:wght@400;900&display=swap",
		ChartJsURL: "https://cdn.jsdelivr.net/npm/chart.js@2.9/dist/Chart.bundle.min.js",
	}))

	app.Use(routes.Return404Handler)

	// start server on 127.0.0.1:3000
	log.Fatal(app.Listen(":3000"))

}
