package user

import (
	"mynotes/database"
	mycookie "mynotes/internal/cookie"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func SetCookie(c *fiber.Ctx, login *string) {
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(mycookie.LifeTime)

	cookie := new(fiber.Cookie)
	cookie.Name = mycookie.Name
	cookie.Value = sessionToken
	cookie.Expires = expiresAt
	c.Cookie(cookie)

	database.CreateSession(sessionToken, *login, expiresAt)

}

func ClearCookie(c *fiber.Ctx) {
	sessionToken := c.Cookies(mycookie.Name)

	cookie := new(fiber.Cookie)
	cookie.Name = mycookie.Name
	cookie.Value = ""
	cookie.Expires = time.Now()
	c.Cookie(cookie)

	database.DelSession(sessionToken)
	database.DelOldSession()
}

func IsLogged(c *fiber.Ctx) (string, bool) {
	sessionToken := c.Cookies(mycookie.Name)
	_, user, expire := database.GetSessionData(sessionToken)
	return user, !expire.Before(time.Now())
}
