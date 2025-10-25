package crypto

import (

	"golang.org/x/crypto/bcrypt"
)

// HashPassword takes a plain-text password and returns the bcrypt-hashed password
func HashPassword(password string) (string, error) {
    // bcrypt.DefaultCost = 10（推荐值），值越高越安全，但加密越慢
    hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedBytes), nil
}

func CheckPassword(rawPassword, hashedPassword string) (bool) {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(rawPassword))
    return err == nil
}

