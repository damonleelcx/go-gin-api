package entity

import (
	"time"
)

// User represents user entity
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`                    // User ID
	Username  string    `json:"username" gorm:"uniqueIndex;not null"`    // Username, unique index
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`       // Email, unique index
	Password  string    `json:"-" gorm:"not null"`                       // Password (not serialized to JSON)
	FirstName string    `json:"first_name" gorm:"type:varchar(100)"`      // First name
	LastName  string    `json:"last_name" gorm:"type:varchar(100)"`      // Last name
	Phone     string    `json:"phone" gorm:"type:varchar(20)"`           // Phone number
	Avatar    string    `json:"avatar" gorm:"type:varchar(255)"`         // Avatar URL
	Status    string    `json:"status" gorm:"type:varchar(20);default:'active'"` // Status: active, inactive, banned
	Role      string    `json:"role" gorm:"type:varchar(20);default:'user'"`     // Role: user, admin, moderator
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`        // Created at
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`        // Updated at
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`     // Soft delete time
}

// TableName specifies table name
func (User) TableName() string {
	return "users"
}

// BeforeCreate hook function before creation (if needed)
// func (u *User) BeforeCreate(tx *gorm.DB) error {
// 	// Can add pre-creation logic here, such as password encryption, etc.
// 	return nil
// }

