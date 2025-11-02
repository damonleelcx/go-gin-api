package entity

import (
	"time"
)

// PasswordResetToken password reset token entity
type PasswordResetToken struct {
	ID        uint      `json:"id" gorm:"primaryKey"`                                    // Primary key ID
	UserID    uint      `json:"user_id" gorm:"not null;index"`                         // User ID, foreign key to User table
	Token     string    `json:"token" gorm:"uniqueIndex;not null;type:varchar(255)"`    // Reset token, unique index
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;index"`                       // Expiration time, indexed
	Used      bool      `json:"used" gorm:"default:false"`                              // Whether it has been used
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`                       // Created at
}

// TableName specifies table name
func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}

// IsExpired checks if token has expired
func (p *PasswordResetToken) IsExpired() bool {
	return time.Now().After(p.ExpiresAt)
}

// IsValid checks if token is valid (not used and not expired)
func (p *PasswordResetToken) IsValid() bool {
	return !p.Used && !p.IsExpired()
}

