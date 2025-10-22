package db

import (
    "log"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

func InitPostgres() {
    dsn := "host=localhost user=postgres password=yourpassword dbname=yourdb port=5432 sslmode=disable"
    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Printf("failed to connect to postgres: %v", err)
    }

    log.Println("PostgreSQL connected.")
}
