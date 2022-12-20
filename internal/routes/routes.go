package routes

import (
	"fmt"
	"html/template"
	"mynotes/database"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func DelNoteHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id") // int 123 and no error
	if err != nil {
		return err
	}

	err = database.DelNotes(id)
	if err != nil {
		return err
	}

	return c.Render("deleted", fiber.Map{
		"Title": "Add note",
		"ID":    id,
	})
	//fmt.Fprintf(c, "%v\n", id)
	//return nil
}

func PaginationHandler(c *fiber.Ctx) error {
	var table []string
	var pages []string

	id, err := c.ParamsInt("id") // int 123 and no error
	if err != nil {
		return err
	}
	//result, pageCounter, _ := database.GetNotes(id)
	result, pageCounter, _ := (&database.SqlHandler{}).GetNotes(id)

	for _, value := range result {
		table = append(table, fmt.Sprintf(`<tr><th scope="row">%v</th><td>%v</td><td>%v</td>
		<td><a class="btn btn-outline-warning" style="float: right;" href="/delnote/%[1]v" role="button">
		Delete</a></td></tr>`, value.ID, value.CreatedAt.Format("2006/01/02 15:04"), value.Text))
	}
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
		//return c.Redirect("/404")
		return c.Status(fiber.StatusNotFound).Render("errors/404", fiber.Map{
			"Error": "Error 404. Not found!",
			"Title": "Error 404. Not found!",
		})
	}

	joinedLines := strings.Join(table, "")
	pagination := strings.Join(pages, "")

	return c.Render("index", fiber.Map{
		"Title":      "Заметки. Страница " + strconv.Itoa(id),
		"Table":      joinedLines,
		"Pagination": pagination,
	})

}

func MainPageHandler(c *fiber.Ctx) error {
	//	var table []string
	var pages []string

	result, pageCounter, _ := (&database.SqlHandler{}).GetNotes(0)
	//date: = database.Note{}.CreatedAt.Format("2006/01/02 15:04")
	for i := 1; i <= pageCounter; i++ {
		pages = append(pages, fmt.Sprintf("<li class=\"page-item\"><a class=\"page-link\" href=\"/pages/%v\">%[1]v</a></li>", i))
	}
	pages = append(pages, "<li class=\"page-item\"><a class=\"page-link\" href=\"/pages/2\">Next</a></li>")

	pagination := strings.Join(pages, "")

	date := fmt.Sprint("%v", database.Note{}.CreatedAt.Format("2006/01/02 15:04"))

	return c.Render("index", fiber.Map{
		"Title":      "Заметки",
		"result":     result,
		"ID":         (&database.Note{}).ID,
		"Text":       (&database.Note{}).Text,
		"CreatedAt":  date,
		"Pagination": pagination,
	})
	//return nil
}

func AddHandler(c *fiber.Ctx) error {
	return c.Render("add", fiber.Map{
		"Title": "Add note",
	})
}

func LayoutHandler(c *fiber.Ctx) error {
	return c.Render("index2", fiber.Map{
		"Title": "Hello, World!",
	}, "layouts/main")
}

func Return404Handler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).Render("errors/404", fiber.Map{
		"Error": "Error 404. Not found!",
		"Title": "Error 404. Not found!",
	})
}

func AddnoteHandler(c *fiber.Ctx) error {
	err := database.CreateNote(template.HTMLEscaper(c.FormValue("confirmationText")))
	if err != nil {
		return err
	}

	return c.Render("success", fiber.Map{
		"Title": "Add note",
	})
}
