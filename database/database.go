package database

import (
	"math"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

type notes []Note

var Db *gorm.DB

func InitDatabase() {

	database, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	err = database.AutoMigrate(&Note{}, &User{})
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

func CreateUser(username, password string) error {
	err := Db.Create(&User{Username: username,
		Password: password,
	}).Error

	if err != nil {
		return err
	}
	return nil
}

func DelNotes(id int) {
	var notes notes
	Db.Delete(&notes, id)
}
