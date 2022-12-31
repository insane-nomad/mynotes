package user

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var Sessions = map[string]session{}

// each session contains the username of the user and the time at which it expires
type session struct {
	username string
	expiry   time.Time
}

func (s session) IsExpired() bool {
	return s.expiry.Before(time.Now())
}

func (s session) IsLogged() bool {
	return !s.expiry.Before(time.Now())
}

func SetCookie(c *fiber.Ctx, login *string) {
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(30 * time.Second)

	cookie := new(fiber.Cookie)
	cookie.Name = "mycookie"
	cookie.Value = sessionToken
	cookie.Expires = expiresAt
	c.Cookie(cookie)

	Sessions[sessionToken] = session{
		username: *login,
		expiry:   expiresAt,
	}
}

func IsLogged(c *fiber.Ctx) bool {
	sessionToken := c.Cookies("mycookie")
	userSession, exists := Sessions[sessionToken]

	if !exists {
		// пользователь не авторизован
		return false
	}

	if userSession.IsExpired() {
		delete(Sessions, sessionToken)
		return false
	}

	fmt.Println(userSession.IsExpired())
	fmt.Println("-------------")
	fmt.Println(sessionToken)
	fmt.Println("-------------")
	fmt.Println(Sessions[c.Cookies("mycookie")])
	fmt.Println("-------------")
	//fmt.Println(IsLogged(c))

	return true
}
