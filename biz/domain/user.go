package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"size:64;not null"`
	Password string `gorm:"not null"`
    Email    string `gorm:"uniqueIndex;size:255"`
	AvatarUrl string `gorm:"type:text"`
}
//note : gorm note only effect autoMigrate, it is not used to validate input