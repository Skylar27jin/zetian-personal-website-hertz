package domain

import (
	"time"

	"github.com/lib/pq"
)

type Category struct {
	ID          int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string         `json:"name" gorm:"type:varchar(64);not null"`         // 展示名
	Key         string         `json:"key"  gorm:"type:varchar(32);not null;unique"` // 稳定英文 key
	Aliases     pq.StringArray `json:"aliases" gorm:"type:text[]"`                   // 搜索别名
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
