package repository

import (
	"log"
	"zetian-personal-website-hertz/biz/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitPostgres() {
    dsn := config.GetSpecificConfig().DB_DSN

    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Printf("failed to connect to postgres: %v", err)
    }

    log.Println("PostgreSQL connected.")
}

