package email_verification_code_repo

import (
	"context"
	"fmt"
	"zetian-personal-website-hertz/biz/domain"
	DB "zetian-personal-website-hertz/biz/repository"
)

// CreateOrUpdateCode inserts or updates an email verification code record.
//
// This function ensures that each email (and purpose) only has one active verification code.
// If a record with the same email and purpose already exists, it will be overwritten
// with the new code and refreshed expiration time.
//
// Parameters:
//   - ctx:            Context for database operations (recommended for request-level timeout/cancel).
//   - email:          The target user's email address (serves as unique identifier).
//   - purpose:        The purpose of the verification (e.g., "signup", "reset_password").
//   - code:           The verification code to be stored (typically a short alphanumeric string).
//   - now:            The current timestamp in Unix seconds. If set to -1, current system time will be used.
//   - validDuration:  The validity period in seconds. If set to -1, defaults to 15 minutes.
//
// Behavior:
//   - Inserts a new record if none exists for (email, purpose).
//   - If an existing record is found, updates code, created_at, expire_at,
//     and resets is_used = false (meaning the code becomes active again).
//
// Returns:
//   - error: Any database execution error during insert/update.
func CreateOrUpdateCode(ctx context.Context, email, purpose, code string, now, validDuration int64) error {
	if email == "" || code == "" {
		return fmt.Errorf("email or code cannot be nil")
	}
	exp := now + validDuration

	return DB.DB.WithContext(ctx).Exec(`
		INSERT INTO email_verification_code (email, purpose, code, created_at, expire_at, is_used)
		VALUES (?, ?, ?, ?, ?, false)
		ON DUPLICATE KEY UPDATE
			code = VALUES(code),
			created_at = VALUES(created_at),
			expire_at = VALUES(expire_at),
			is_used = false;
	`, email, purpose, code, now, exp).Error
}

// GetCodeByEmail fetches the verification code record by email and purpose.
func GetCodeByEmail(ctx context.Context, email string) (domain.EmailVerificationCode, error) { 

	var record domain.EmailVerificationCode
	err := DB.DB.WithContext(ctx).
		Where("email = ?", email).
		First(&record).Error
	if err != nil {
		return domain.EmailVerificationCode{}, err
	}
	return record, nil
}

// MarkCodeAsUsed sets the code record's is_used to true.
func MarkCodeAsUsed(ctx context.Context, email string) error {
	return DB.DB.WithContext(ctx).
		Model(&domain.EmailVerificationCode{}).
		Where("email = ?", email).
		Update("is_used", true).Error
}

// DeleteCode removes a code record (e.g., after successful verification or cleanup).
func DeleteCode(ctx context.Context, email string) error {
	return DB.DB.WithContext(ctx).
		Where("email = ?", email).
		Delete(&domain.EmailVerificationCode{}).Error
}

// GetTableName returns the DB table name.
func GetTableName() string {
	return "email_verification_code"
}
