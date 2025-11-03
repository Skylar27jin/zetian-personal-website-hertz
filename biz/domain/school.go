package domain

import "time"

type School struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"type:varchar(128);not null"`
	ShortName   string    `json:"short_name" gorm:"type:varchar(32);not null"`
	Aliases     []string  `json:"aliases" gorm:"type:text[]"` // PostgreSQL 数组
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
