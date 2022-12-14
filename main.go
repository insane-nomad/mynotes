package main

import (
	"fmt"
	"html/template"
	"log"
	"mynotes/database"
	"strconv"

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

	app.Get("/pages/:id<range(1,1000)>", PaginationHandler)
	app.Use(Return404Handler)

	// start server on 127.0.0.1:3000
	log.Fatal(app.Listen(":3000"))

}

func PaginationHandler(c *fiber.Ctx) error {
	var table []string
	var pages []string

	id, err := c.ParamsInt("id") // int 123 and no error
	if err != nil {
		return err
	}
	result, pageCounter, _ := database.GetNotes(id)

	for _, value := range result {
		table = append(table, fmt.Sprintf("<tr><th scope=\"row\">%v</th><td>%v</td><td>%v</td></tr>",
			value.ID, value.CreatedAt.Format("2006/01/02 15:04"), value.Text))
	}
	joinedLines := strings.Join(table, "")
	if id > 1 {
		pages = append(pages, fmt.Sprintf("<li class=\"page-item\"><a class=\"page-link\" href=\"/pages/%v\">Previous</a></li>", id-1))
	}

	for i := 1; i <= pageCounter; i++ {
		pages = append(pages, fmt.Sprintf("<li class=\"page-item\"><a class=\"page-link\" href=\"/pages/%v\">%[1]v</a></li>", i))
	}
	if id < pageCounter {
		pages = append(pages, fmt.Sprintf("<li class=\"page-item\"><a class=\"page-link\" href=\"/pages/%v\">Next</a></li>", id+1))
	}

	if id > pageCounter {
		return c.Status(fiber.StatusNotFound).Render("index", fiber.Map{
			"Error": "Error 404. Not found!",
			"Title": "Error 404. Not found!",
		}, "errors/404")
	}
	pagination := strings.Join(pages, "")

	return c.Render("index", fiber.Map{
		"Title":      "Заметки. Страница " + strconv.Itoa(id),
		"Table":      joinedLines,
		"Pagination": pagination,
	})

}

func MainPageHandler(c *fiber.Ctx) error {
	var table []string
	var pages []string
	result, pageCounter, _ := database.GetNotes(0)

	for _, value := range result {
		table = append(table, fmt.Sprintf("<tr><th scope=\"row\">%v</th><td>%v</td><td>%v</td></tr>",
			value.ID, value.CreatedAt.Format("2006/01/02 15:04"), value.Text))
	}
	joinedLines := strings.Join(table, "")

	for i := 1; i <= pageCounter; i++ {
		pages = append(pages, fmt.Sprintf("<li class=\"page-item\"><a class=\"page-link\" href=\"/pages/%v\">%[1]v</a></li>", i))
	}
	pages = append(pages, "<li class=\"page-item\"><a class=\"page-link\" href=\"/pages/2\">Next</a></li>")

	pagination := strings.Join(pages, "")

	return c.Render("index", fiber.Map{
		"Title":      "Заметки",
		"Table":      joinedLines,
		"Pagination": pagination,
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

	/*убираем весь колхоз и получаем данные правильно
	confirmationText, err := url.ParseQuery(string(c.Body()))
	if err != nil {
		return err
	}

	fmt.Fprintf(c, "%v\n", c.FormValue("confirmationText"))
	result := strings.Join(confirmationText["confirmationText"], " ")
	*/
	err := database.CreateNote(c.FormValue("confirmationText"))
	if err != nil {
		return err
	}

	return c.Render("success", fiber.Map{
		"Title": "Add note",
	})

}
