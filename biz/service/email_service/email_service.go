package email_service

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"zetian-personal-website-hertz/biz/pkg/SES_email"
)

func GenerateVeriCode(length int, hasNumber bool, hasLowerEngChar bool, hasUpperEngChar bool, hasSpecialChar bool) (string, error) {
	if length <= 0 {
		return "", errors.New("length must be greater than 0")
	}

	var charset string
	if hasNumber {
		charset += "0123456789"
	}
	if hasLowerEngChar {
		charset += "abcdefghijklmnopqrstuvwxyz"
	}
	if hasUpperEngChar {
		charset += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	if hasSpecialChar {
		charset += "!@#$%^&*()-_=+[]{}<>?"
	}

	if len(charset) == 0 {
		return "", errors.New("no character types selected")
	}

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[index.Int64()]
	}

	return string(result), nil
}

func SendVeriCodeEmailTo(ctx context.Context, to, code, purpose string) error {
	err := SES_email.SendEmail(
		ctx,
		to,
		getDefaultVeriCodeEmailSubject(),
		getDefaultVeriCodeEmailBody(to, code, purpose),
	)
	if err != nil {
		return err
	}
	return nil
}

func getDefaultVeriCodeEmailSubject() string {
	return "Hello! Verification Code From skylar27.com~"
}

func getDefaultVeriCodeEmailBody(to, code, purpose string) string {
	body := ""
	body += "Hi " + to + ",\n\n"
	body += "Your verification code is: " + code + "\n\n"
	body += "This code seems to be used to: " + purpose + "\n\n"
	body += "⚠️ Please do NOT share this code with anyone. and don't reply to this email\n\n"
	body += "Thank you for using our service!\n"
	body += "— skylar27.com"
	return body
}
