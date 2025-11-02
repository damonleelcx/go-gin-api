package repository

import (
	"errors"
	"time"

	"github.com/damonleelcx/go-gin-api/entity"
	"gorm.io/gorm"
)

// SessionRepository session repository interface
type SessionRepository interface {
	// FindByToken find session by token
	FindByToken(token string) (*entity.Session, error)
	// FindByID find session by ID
	FindByID(id uint) (*entity.Session, error)
	// FindByUserID find all sessions by user ID
	FindByUserID(userID uint) ([]*entity.Session, error)
	// Create create session
	Create(session *entity.Session) error
	// Update update session
	Update(session *entity.Session) error
	// UpdateLastUsedAt update last used time
	UpdateLastUsedAt(sessionID uint, lastUsedAt time.Time) error
	// UpdateStatusByUserID update status of all active sessions for specified user
	UpdateStatusByUserID(userID uint, status string) error
	// Delete delete session
	Delete(id uint) error
}

// sessionRepository session repository implementation
type sessionRepository struct {
	db *gorm.DB
}

// NewSessionRepository creates a new session repository instance
func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{
		db: db,
	}
}

// FindByToken find session by token
func (r *sessionRepository) FindByToken(token string) (*entity.Session, error) {
	var session entity.Session
	if err := r.db.Where("token = ?", token).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("session does not exist")
		}
		return nil, err
	}
	return &session, nil
}

// FindByID find session by ID
func (r *sessionRepository) FindByID(id uint) (*entity.Session, error) {
	var session entity.Session
	if err := r.db.First(&session, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("session does not exist")
		}
		return nil, err
	}
	return &session, nil
}

// FindByUserID find all sessions by user ID
func (r *sessionRepository) FindByUserID(userID uint) ([]*entity.Session, error) {
	var sessions []*entity.Session
	if err := r.db.Where("user_id = ?", userID).Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

// Create create session
func (r *sessionRepository) Create(session *entity.Session) error {
	return r.db.Create(session).Error
}

// Update update session
func (r *sessionRepository) Update(session *entity.Session) error {
	return r.db.Save(session).Error
}

// UpdateLastUsedAt update last used time
func (r *sessionRepository) UpdateLastUsedAt(sessionID uint, lastUsedAt time.Time) error {
	return r.db.Model(&entity.Session{}).
		Where("id = ?", sessionID).
		Update("last_used_at", &lastUsedAt).Error
}

// UpdateStatusByUserID update status of all active sessions for specified user
func (r *sessionRepository) UpdateStatusByUserID(userID uint, status string) error {
	result := r.db.Model(&entity.Session{}).
		Where("user_id = ? AND status = ?", userID, "active").
		Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Delete delete session
func (r *sessionRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Session{}, id).Error
}

