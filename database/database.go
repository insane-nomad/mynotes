package database

import (
	//"log"
	"math"
	//"os"
	"mynotes/internal/cookie"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	//"gorm.io/gorm/logger"
)

const dbName string = "mynotes.db"

type Note struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	Text      string
}

type User struct {
	ID       uint `gorm:"primarykey"`
	Username string
	Password string
}

type Session struct {
	ID       uint `gorm:"primarykey"`
	Session  string
	Username string
	Expiry   time.Time
}

type notes []Note

var Db *gorm.DB

func InitDatabase() {
	// newLogger := logger.New(
	// 	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	// 	logger.Config{
	// 		SlowThreshold:             time.Second, // Slow SQL threshold
	// 		LogLevel:                  logger.Info, // Log level
	// 		IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
	// 		Colorful:                  false,       // Disable color
	// 	},
	// )

	database, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{
		// Logger: newLogger,
	})

	if err != nil {
		panic("Failed to connect to database!")
	}

	err = database.AutoMigrate(&Note{}, &User{}, &Session{})
	if err != nil {
		return
	}
	//database.AutoMigrate(&User{})

	Db = database
}

func GetNotes(start int) (notes, int, error) {
	var notes notes
	offset := 10
	cstart := (start - 1) * offset
	var pageCount float64

	getCount := Db.Find(&notes)
	pageCount = float64(getCount.RowsAffected) / float64(offset)
	pageCount = math.Ceil(pageCount)

	Db.Limit(offset).Offset(cstart).Find(&notes)

	return notes, int(pageCount), nil
}

func CreateNote(text string) {
	Db.Create(&Note{Text: text})
}

func DelNotes(id int) {
	var notes notes
	Db.Delete(&notes, id)
}

func CreateUser(username, password string) error {
	err := Db.Create(&User{Username: username,
		Password: password,
	}).Error

	if err != nil {
		return err
	}
	return nil
}

func Login(username, password string) bool {
	var user User
	Db.Where("username = ?", username).First(&user)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return false
	}
	return true
}

func CreateSession(sessionToken, username string, expire time.Time) error {
	err := Db.Create(&Session{
		Session:  sessionToken,
		Username: username,
		Expiry:   expire,
	}).Error

	if err != nil {
		return err
	}
	return nil
}

func GetSessionData(sessionToken string) (string, string, time.Time) {
	var sess Session
	if sessionToken == "" {
		return "", "", time.Date(2000, 11, 14, 16, 45, 16, 36, time.UTC)
	}
	Db.Where("session = ?", sessionToken).First(&sess)
	return sess.Session, sess.Username, sess.Expiry

}

func DelSession(sessionToken string) {
	Db.Delete(&Session{}, "session LIKE ?", sessionToken)
}

func DelOldSession() {
	Db.Delete(&Session{}, "expiry < '?'", time.Now().Add(-cookie.LifeTime))
}
