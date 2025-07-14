package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDb() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("storage.db"))
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&User{}, &Book{})
	log.Println("db connected!!")

	return db, nil
}