package main

import (
	"fmt"
	"html/template"
	"log"
	"mynotes/database"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

func main() {
	// Create a new database
	dbErr := database.InitDatabase()

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

	//app.Static("ui/html/static", "ui/html/static")   //рабочая строка
	app.Static("ui/html/static", "ui/html/static", fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		CacheDuration: 1000 * time.Hour,
		MaxAge:        3600,
	})

	// Routes
	app.Post("/addnote", AddnoteHandler)

	app.Get("/", MainPageHandler)
	app.Get("/layout", LayoutHandler)
	app.Get("/add", AddHandler)

	app.Use(Return404Handler)

	// start server on 127.0.0.1:3000
	log.Fatal(app.Listen(":3000"))

}

func MainPageHandler(c *fiber.Ctx) error {
	var table []string
	result, _ := database.GetAllNotes()
	//fmt.Fprintf(c, "%T\n", result)
	for _, value := range result {
		//fmt.Fprintf(c, "%v\n", value.Text)
		table = append(table, fmt.Sprintf("<tr><th scope=\"row\">%v</th><td>%v</td><td>%v</td></tr>",
			value.ID, value.CreatedAt.Format("2006/01/02 15:04"), value.Text))
	}
	joinedLines := strings.Join(table, "")
	fmt.Fprintf(c, "%v\n", joinedLines)

	return c.Render("index", fiber.Map{
		"Title": "Заметки",
		"Table": joinedLines,
	})
	//return nil
}

func AddHandler(c *fiber.Ctx) error {
	// Render index
	return c.Render("add", fiber.Map{
		"Title": "Add note",
	})
}

func LayoutHandler(c *fiber.Ctx) error {
	// Render index within layouts/main
	return c.Render("index2", fiber.Map{
		"Title": "Hello, World!",
	}, "layouts/main")
}

func Return404Handler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).Render("index", fiber.Map{
		"Error": "Error 404. Not found!",
		"Title": "Error 404. Not found!",
	}, "errors/404")
}

func AddnoteHandler(c *fiber.Ctx) error {
	confirmationText, err := url.ParseQuery(string(c.Body()))
	if err != nil {
		return err
	}
	result := strings.Join(confirmationText["confirmationText"], " ")

	err = database.CreateNote(result)
	if err != nil {
		return err
	}

	return c.Render("success", fiber.Map{
		"Title": "Add note",
	})
}
