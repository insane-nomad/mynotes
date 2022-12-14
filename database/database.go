package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Note struct {
	gorm.Model
	Text string
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

	db.Find(&notes)

	return notes, nil
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
