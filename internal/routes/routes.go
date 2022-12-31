package routes

import (
	"fmt"
	"html/template"
	"mynotes/database"
	"mynotes/internal/user"
	"strconv"

	//"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
)

const error404 string = "Error 404. Not found!"

func DelNoteHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id") // int 123 and no error
	if err != nil {
		return err
	}

	database.DelNotes(id)

	return c.Render("deleted", fiber.Map{
		"Title": "Add note",
		"ID":    id,
	})
	//fmt.Fprintf(c, "%v\n", id)
	//return nil
}

func PaginationHandler(c *fiber.Ctx) error {
	var pages []int
	var NotFirstPage, NotLastPage bool
	var nextPage, prevPage int

	id, err := c.ParamsInt("id") // int 123 and no error
	if err != nil {
		return err
	}
	result, pageCounter, _ := database.GetNotes(id)

	nextPage = id + 1
	prevPage = id - 1

	if id > 1 {
		NotFirstPage = true
	}

	if id < pageCounter {
		NotLastPage = true
	}

	for i := 1; i <= pageCounter; i++ {
		pages = append(pages, i)
	}

	if id > pageCounter {
		return c.Status(fiber.StatusNotFound).Render("errors/404", fiber.Map{
			"Error": error404,
			"Title": error404,
		})
	}
	//fmt.Fprintf(os.Stdout, "%v\n", NotLastPage)

	return c.Render("index", fiber.Map{
		"Title":        "Заметки. Страница " + strconv.Itoa(id),
		"result":       result,
		"pages":        pages,
		"NotFirstPage": NotFirstPage,
		"NotLastPage":  NotLastPage,
		"nextPage":     nextPage,
		"prevPage":     prevPage,
	})

}

func MainPageHandler(c *fiber.Ctx) error {

	fmt.Println(user.IsLogged(c))

	var pages []int
	result, pageCounter, _ := database.GetNotes(0)

	for i := 1; i <= pageCounter; i++ {
		pages = append(pages, i)
	}

	return c.Render("index", fiber.Map{
		"Title":        "Заметки",
		"result":       result,
		"pages":        pages,
		"NotFirstPage": false,
		"NotLastPage":  true,
		"nextPage":     2,
	})
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
		"Error": error404,
		"Title": error404,
	})
}

func AddnoteHandler(c *fiber.Ctx) error {
	database.CreateNote(template.HTMLEscaper(c.FormValue("confirmationText")))
	return c.Render("success", fiber.Map{
		"Title": "Add note",
	})
}

func RegisterHandler(c *fiber.Ctx) error {
	return c.Render("user/register", fiber.Map{
		"Title": "Регистрация нового пользователя",
	})
}

func AdduserHandler(c *fiber.Ctx) error {
	password := template.HTMLEscaper(c.FormValue("password"))
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return err
	}
	err = database.CreateUser(template.HTMLEscaper(c.FormValue("login")), string(hashedPassword))
	if err != nil {
		return c.Render("user/fail", fiber.Map{
			"Title":   "Add user",
			"Message": "Что-то пошло не так",
		})
	}
	return c.Render("user/success", fiber.Map{
		"Title":   "Add user",
		"Message": "Вы успешно зарегистрировались",
	})
}

func LoginHandler(c *fiber.Ctx) error {
	login := template.HTMLEscaper(c.FormValue("login"))
	password := template.HTMLEscaper(c.FormValue("password"))

	isLogged := database.Login(login, password)

	if isLogged {

		user.SetCookie(c, &login)

		return c.Render("user/success", fiber.Map{
			"Title":   "Add user",
			"Message": "Вы успешно залогинились",
		})
	} else {
		return c.Render("user/fail", fiber.Map{
			"Title":   "Add user",
			"Message": "Неверный логин или пароль",
		})
	}
}
