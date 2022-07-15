package core

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DbSession *gorm.DB

func initDb() {
	var err error

	DbSession, err = gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}
