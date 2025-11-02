package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/damonleelcx/go-gin-api/entity"
	"github.com/damonleelcx/go-gin-api/repository"
	"golang.org/x/crypto/bcrypt"
)

// AuthService authentication service
type AuthService struct {
	userRepo                repository.UserRepository
	sessionRepo             repository.SessionRepository
	passwordResetTokenRepo  repository.PasswordResetTokenRepository
}

// NewAuthService creates a new authentication service instance
func NewAuthService(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	passwordResetTokenRepo repository.PasswordResetTokenRepository,
) *AuthService {
	return &AuthService{
		userRepo:               userRepo,
		sessionRepo:            sessionRepo,
		passwordResetTokenRepo: passwordResetTokenRepo,
	}
}

// SignupRequest registration request
type SignupRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

// SignupResponse registration response
type SignupResponse struct {
	User    *entity.User `json:"user"`
	Message string       `json:"message"`
}

// SigninRequest login request
type SigninRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// SigninResponse login response
type SigninResponse struct {
	User    *entity.User   `json:"user"`
	Session *entity.Session `json:"session"`
	Token   string         `json:"token"`
	Message string         `json:"message"`
}

// ForgotPasswordRequest forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest reset password request
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// Signup user registration
func (s *AuthService) Signup(req *SignupRequest, ipAddress, userAgent string) (*SignupResponse, error) {
	// Check if username or email already exists
	exists, existingUser, err := s.userRepo.Exists(req.Username, req.Email)
	if err != nil {
		return nil, errors.New("failed to query user: " + err.Error())
	}
	if exists {
		if existingUser.Username == req.Username {
			return nil, errors.New("username already exists")
		}
		if existingUser.Email == req.Email {
			return nil, errors.New("email already registered")
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("password encryption failed: " + err.Error())
	}

	// Create user
	user := &entity.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Status:    "active",
		Role:      "user",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user: " + err.Error())
	}

	return &SignupResponse{
		User:    user,
		Message: "Registration successful",
	}, nil
}

// Signin user login
func (s *AuthService) Signin(req *SigninRequest, ipAddress, userAgent string) (*SigninResponse, error) {
	// Find user (supports username or email login)
	user, err := s.userRepo.FindByUsernameOrEmail(req.Username)
	if err != nil {
		return nil, errors.New("username or password incorrect")
	}

	// Check user status
	if user.Status != "active" {
		return nil, errors.New("account has been disabled")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("username or password incorrect")
	}

	// Generate session token
	token, err := generateToken()
	if err != nil {
		return nil, errors.New("failed to generate token: " + err.Error())
	}

	// Create session
	now := time.Now()
	session := &entity.Session{
		UserID:    user.ID,
		Token:     token,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Status:    "active",
		ExpiresAt: time.Now().Add(24 * 7 * time.Hour), // 7 days expiry
		LastUsedAt: &now,
	}

	if err := s.sessionRepo.Create(session); err != nil {
		return nil, errors.New("failed to create session: " + err.Error())
	}

	// Clear password field
	user.Password = ""

	return &SigninResponse{
		User:    user,
		Session: session,
		Token:   token,
		Message: "Login successful",
	}, nil
}

// Logout user logout
func (s *AuthService) Logout(token string) error {
	// Find session
	session, err := s.sessionRepo.FindByToken(token)
	if err != nil {
		return err
	}

	// Update session status to logged out
	session.Status = "logout"
	if err := s.sessionRepo.Update(session); err != nil {
		return errors.New("failed to update session status: " + err.Error())
	}

	return nil
}

// LogoutAll logout all sessions of the user
func (s *AuthService) LogoutAll(userID uint) error {
	// Update status of all active sessions for this user
	if err := s.sessionRepo.UpdateStatusByUserID(userID, "logout"); err != nil {
		return errors.New("failed to logout all sessions: " + err.Error())
	}

	return nil
}

// ForgotPassword forgot password
func (s *AuthService) ForgotPassword(req *ForgotPasswordRequest) (string, error) {
	// Find user
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		// For security, return success message even if email doesn't exist
		return "If the email exists, reset link has been sent", nil
	}

	// Generate reset token
	token, err := generateToken()
	if err != nil {
		return "", errors.New("failed to generate reset token: " + err.Error())
	}

	// Create password reset token (valid for 1 hour)
	resetToken := &entity.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
	}

	if err := s.passwordResetTokenRepo.Create(resetToken); err != nil {
		return "", errors.New("failed to create reset token: " + err.Error())
	}

	// In actual application, email should be sent here
	// Now only return token (production should send via email)
	return token, nil
}

// ResetPassword reset password
func (s *AuthService) ResetPassword(req *ResetPasswordRequest) error {
	// Find reset token
	resetToken, err := s.passwordResetTokenRepo.FindByToken(req.Token)
	if err != nil {
		return err
	}

	// Check if token has been used
	if resetToken.Used {
		return errors.New("reset token has been used")
	}

	// Check if token has expired
	if time.Now().After(resetToken.ExpiresAt) {
		return errors.New("reset token has expired")
	}

	// Find user
	user, err := s.userRepo.FindByID(resetToken.UserID)
	if err != nil {
		return errors.New("user does not exist")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("password encryption failed: " + err.Error())
	}

	// Update user password
	user.Password = string(hashedPassword)
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("failed to update password: " + err.Error())
	}

	// Mark token as used
	resetToken.Used = true
	if err := s.passwordResetTokenRepo.Update(resetToken); err != nil {
		return errors.New("failed to update token status: " + err.Error())
	}

	// Invalidate all sessions for this user (security consideration)
	if err := s.LogoutAll(user.ID); err != nil {
		// Log error but don't affect password reset flow
		// log.Printf("failed to logout all sessions: %v", err)
	}

	return nil
}

// ValidateToken validate session token
func (s *AuthService) ValidateToken(token string) (*entity.Session, *entity.User, error) {
	// Find session
	session, err := s.sessionRepo.FindByToken(token)
	if err != nil {
		return nil, nil, err
	}

	// Check if session is valid
	if !session.IsActive() {
		return nil, nil, errors.New("session has expired")
	}

	// Update last used time
	now := time.Now()
	if err := s.sessionRepo.UpdateLastUsedAt(session.ID, now); err != nil {
		// Log error but don't affect validation flow
	}

	// Find user
	user, err := s.userRepo.FindByID(session.UserID)
	if err != nil {
		return nil, nil, errors.New("user does not exist")
	}

	// Check user status
	if user.Status != "active" {
		return nil, nil, errors.New("user has been disabled")
	}

	// Clear password field
	user.Password = ""

	return session, user, nil
}

// generateToken generate random token
func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
