package entity

import (
	"time"
)

// Session represents user session entity
type Session struct {
	ID        uint      `json:"id" gorm:"primaryKey"`                                  // Session ID
	UserID    uint      `json:"user_id" gorm:"not null;index"`                        // User ID, foreign key to User table
	Token     string    `json:"token" gorm:"uniqueIndex;not null;type:varchar(255)"`   // Session token, unique index
	IPAddress string    `json:"ip_address" gorm:"type:varchar(45)"`                    // IP address (supports IPv6)
	UserAgent string    `json:"user_agent" gorm:"type:varchar(500)"`                   // User agent information
	Device    string    `json:"device" gorm:"type:varchar(50)"`                        // Device type: web, mobile, tablet
	Platform  string    `json:"platform" gorm:"type:varchar(50)"`                      // Platform: windows, macos, linux, ios, android
	Status    string    `json:"status" gorm:"type:varchar(20);default:'active'"`       // Status: active, expired, revoked
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;index"`                      // Expiration time, indexed
	LastUsedAt *time.Time `json:"last_used_at" gorm:"index"`                          // Last used time, indexed
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`                      // Created at
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`                      // Updated at
}

// TableName specifies table name
func (Session) TableName() string {
	return "sessions"
}

// IsExpired checks if session has expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// IsActive checks if session is active
func (s *Session) IsActive() bool {
	return s.Status == "active" && !s.IsExpired()
}

// User associated user (if using GORM association feature)
// func (s *Session) User() User {
// 	var user User
// 	// Can query through GORM association here
// 	return user
// }

