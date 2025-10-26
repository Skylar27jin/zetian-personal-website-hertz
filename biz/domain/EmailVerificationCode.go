package domain

type EmailVerificationCode struct {
    Email     string `gorm:"primaryKey;size:255" json:"email"`
    Purpose   string `gorm:"size:50;default:signup" json:"purpose"`
    Code      string `gorm:"size:10" json:"code"`
    CreatedAt int64  `json:"created_at"`
    ExpireAt  int64  `json:"expire_at"`
}
