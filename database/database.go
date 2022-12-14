package database

import (
	//	"fmt"
	"math"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Note struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	Text      string
}

func InitDatabase() error {
	db, err := gorm.Open(sqlite.Open("mynotes.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	db.AutoMigrate(&Note{})

	return nil
}

func GetAllNotes() ([]Note, error) {
	var notes []Note

	db, err := gorm.Open(sqlite.Open("mynotes.db"), &gorm.Config{})
	if err != nil {
		return notes, err
	}

	db.Find(&notes) // get all data
	//db.Limit(5).Offset(0).Find(&notes)
	return notes, nil
}

func GetNotes(start int) ([]Note, int, error) {
	var notes []Note
	offset := 5
	cstart := start - 1
	cstart = cstart * offset
	var pageCount float64

	db, err := gorm.Open(sqlite.Open("mynotes.db"), &gorm.Config{})
	if err != nil {
		return notes, 0, err
	}

	getCount := db.Find(&notes)
	pageCount = float64(getCount.RowsAffected) / float64(offset) // returns count of records found
	pageCount = math.Ceil(pageCount)

	db.Limit(offset).Offset(cstart).Find(&notes)
	//res := db.First(&notes)
	//result.RowsAffected // returns count of records found
	//fmt.Printf("%v\n", res.RowsAffected)
	//db.Limit(5).Offset(5).Find(&notes)
	return notes, int(pageCount), nil
}

func CreateNote(text string) error {
	//var newNotes = Note{Text: text}

	db, err := gorm.Open(sqlite.Open("mynotes.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	db.Create(&Note{Text: text})

	return nil
}
