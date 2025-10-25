package repository

import (
    "log"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

func InitPostgres() {
    dsn := "host=zetian-personal-website-postgre.c1uyeekq4253.us-east-2.rds.amazonaws.com user=skylar27jin password=zetian-personal-website-postgre dbname=postgres port=5432 sslmode=require"

    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Printf("failed to connect to postgres: %v", err)
    }

    log.Println("PostgreSQL connected.")
}

