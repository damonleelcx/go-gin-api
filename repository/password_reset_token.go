package repository

import (
	"errors"

	"github.com/damonleelcx/go-gin-api/entity"
	"gorm.io/gorm"
)

// PasswordResetTokenRepository password reset token repository interface
type PasswordResetTokenRepository interface {
	// FindByToken find password reset token by token
	FindByToken(token string) (*entity.PasswordResetToken, error)
	// Create create password reset token
	Create(token *entity.PasswordResetToken) error
	// Update update password reset token
	Update(token *entity.PasswordResetToken) error
}

// passwordResetTokenRepository password reset token repository implementation
type passwordResetTokenRepository struct {
	db *gorm.DB
}

// NewPasswordResetTokenRepository creates a new password reset token repository instance
func NewPasswordResetTokenRepository(db *gorm.DB) PasswordResetTokenRepository {
	return &passwordResetTokenRepository{
		db: db,
	}
}

// FindByToken find password reset token by token
func (r *passwordResetTokenRepository) FindByToken(token string) (*entity.PasswordResetToken, error) {
	var resetToken entity.PasswordResetToken
	if err := r.db.Where("token = ?", token).First(&resetToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("reset token invalid")
		}
		return nil, err
	}
	return &resetToken, nil
}

// Create create password reset token
func (r *passwordResetTokenRepository) Create(token *entity.PasswordResetToken) error {
	return r.db.Create(token).Error
}

// Update update password reset token
func (r *passwordResetTokenRepository) Update(token *entity.PasswordResetToken) error {
	return r.db.Save(token).Error
}

