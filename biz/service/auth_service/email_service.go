package auth_service

import (
	"context"
	"time"
	"zetian-personal-website-hertz/biz/domain"
	"zetian-personal-website-hertz/biz/repository/email_verification_code_repo"
)

// CreateOrUpdateCode creates or updates an email verification code record.
//
// This service-level wrapper ensures default values and provides business logic encapsulation
// on top of the repository function.
//
// Parameters:
//   - ctx:            Request context.
//   - email:          Target user's email (unique identifier).
//   - purpose:        Purpose of verification (e.g., "signup", "reset_password").
//   - code:           Verification code string.
//   - now:            Current Unix timestamp (if -1, will use system time).
//   - validDuration:  Code validity in seconds (if -1, defaults to 15 minutes).
//
// Behavior:
//   - Delegates to email_verification_code_repo.CreateOrUpdateCode.
//   - Automatically applies defaults for time-related fields.
//   - Can be extended later for rate limiting, notification triggers, or audit logging.
func CreateOrUpdateCode(ctx context.Context, email, purpose, code string, now, validDuration int64) error {
	if now == -1 {
		now = time.Now().Unix()
	}
	if validDuration == -1 {
		validDuration = 15 * 60 // default: 15 minutes
	}
	return email_verification_code_repo.CreateOrUpdateCode(ctx, email, purpose, code, now, validDuration)
}

// GetCodeByEmail fetches the verification code record by email.
//
// Returns:
//   - *domain.EmailVerificationCode: The matched record.
//   - error: If record not found or DB error occurs.
func GetCodeByEmail(ctx context.Context, email string) (domain.EmailVerificationCode, error) {
	return email_verification_code_repo.GetCodeByEmail(ctx, email)
}

// DeleteCode removes a verification code record (e.g., after successful verification or expiration cleanup).
func DeleteCode(ctx context.Context, email string) error {
	return email_verification_code_repo.DeleteCode(ctx, email)
}
