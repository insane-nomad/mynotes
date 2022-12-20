package database

import (
	//	"fmt"
	// "log"
	"math"

	// "os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const dbName string = "mynotes.db"

type Note struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	Text      string
}

type notes []Note

type SqlHandler struct {
	db *gorm.DB
}

func (s *SqlHandler) InitDatabase() error {
	var err error
	s.db, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		return err
	}

	s.db.AutoMigrate(&Note{})
	return nil
}

func (s *SqlHandler) GetNotes(start int) (notes, int, error) {
	// newLogger := logger.New(
	// 	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	// 	logger.Config{
	// 		SlowThreshold:             time.Second, // Slow SQL threshold
	// 		LogLevel:                  logger.Info, // Log level
	// 		IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
	// 		Colorful:                  false,       // Disable color
	// 	},
	// )
	var notes notes
	var err error
	offset := 10
	cstart := (start - 1) * offset
	var pageCount float64

	s.db, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{
		//		Logger: newLogger,
	})
	if err != nil {
		return notes, 0, err
	}

	getCount := s.db.Find(&notes)
	pageCount = float64(getCount.RowsAffected) / float64(offset)
	pageCount = math.Ceil(pageCount)

	s.db.Limit(offset).Offset(cstart).Find(&notes)
	return notes, int(pageCount), nil
}

func CreateNote(text string) error {
	//var newNotes = Note{Text: text}

	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}
	db.Create(&Note{Text: text})

	return nil
}

func DelNotes(id int) error {
	var notes notes

	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{
		//	Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	db.Delete(&notes, id)
	return nil
}
