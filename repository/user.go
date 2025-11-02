package repository

import (
	"errors"

	"github.com/damonleelcx/go-gin-api/entity"
	"gorm.io/gorm"
)

// UserRepository user repository interface
type UserRepository interface {
	// FindByID find user by ID
	FindByID(id uint) (*entity.User, error)
	// FindByUsername find user by username
	FindByUsername(username string) (*entity.User, error)
	// FindByEmail find user by email
	FindByEmail(email string) (*entity.User, error)
	// FindByUsernameOrEmail find user by username or email
	FindByUsernameOrEmail(usernameOrEmail string) (*entity.User, error)
	// Exists check if username or email already exists
	Exists(username, email string) (bool, *entity.User, error)
	// Create create user
	Create(user *entity.User) error
	// Update update user
	Update(user *entity.User) error
}

// userRepository user repository implementation
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// FindByID find user by ID
func (r *userRepository) FindByID(id uint) (*entity.User, error) {
	var user entity.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user does not exist")
		}
		return nil, err
	}
	return &user, nil
}

// FindByUsername find user by username
func (r *userRepository) FindByUsername(username string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user does not exist")
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail find user by email
func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user does not exist")
		}
		return nil, err
	}
	return &user, nil
}

// FindByUsernameOrEmail find user by username or email
func (r *userRepository) FindByUsernameOrEmail(usernameOrEmail string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user does not exist")
		}
		return nil, err
	}
	return &user, nil
}

// Exists check if username or email already exists
func (r *userRepository) Exists(username, email string) (bool, *entity.User, error) {
	var user entity.User
	err := r.db.Where("username = ? OR email = ?", username, email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil, nil
		}
		return false, nil, err
	}
	return true, &user, nil
}

// Create create user
func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

// Update update user
func (r *userRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

