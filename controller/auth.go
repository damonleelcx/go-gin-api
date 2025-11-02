package controller

import (
	"net/http"

	"github.com/damonleelcx/go-gin-api/service"
	"github.com/gin-gonic/gin"
)

// AuthController authentication controller
type AuthController struct {
	authService *service.AuthService
}

// NewAuthController creates a new authentication controller instance
func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Signup user registration
// @Summary User registration
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body service.SignupRequest true "Registration information"
// @Success 200 {object} service.SignupResponse
// @Failure 400 {object} map[string]string
// @Router /auth/signup [post]
func (ac *AuthController) Signup(c *gin.Context) {
	var req service.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request parameters: " + err.Error(),
		})
		return
	}

	// Get client IP and User-Agent
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// Call service layer
	response, err := ac.authService.Signup(&req, ipAddress, userAgent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Signin user login
// @Summary User login
// @Description User login and get session token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body service.SigninRequest true "Login information"
// @Success 200 {object} service.SigninResponse
// @Failure 400 {object} map[string]string
// @Router /auth/signin [post]
func (ac *AuthController) Signin(c *gin.Context) {
	var req service.SigninRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request parameters: " + err.Error(),
		})
		return
	}

	// Get client IP and User-Agent
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// Call service layer
	response, err := ac.authService.Signin(&req, ipAddress, userAgent)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Logout user logout
// @Summary User logout
// @Description Invalidate current session
// @Tags auth
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/logout [post]
func (ac *AuthController) Logout(c *gin.Context) {
	// Get token from request header
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Missing authentication token",
		})
		return
	}

	// Remove "Bearer " prefix if exists
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// Call service layer
	if err := ac.authService.Logout(token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}

// LogoutAll logout all sessions
// @Summary Logout all sessions
// @Description Invalidate all sessions of the current user
// @Tags auth
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/logout-all [post]
func (ac *AuthController) LogoutAll(c *gin.Context) {
	// Get token from request header and validate
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Missing authentication token",
		})
		return
	}

	// Remove "Bearer " prefix if exists
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// Validate token and get user information
	_, user, err := ac.authService.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token: " + err.Error(),
		})
		return
	}

	// Call service layer to logout all sessions
	if err := ac.authService.LogoutAll(user.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All sessions logged out",
	})
}

// ForgotPassword forgot password
// @Summary Forgot password
// @Description Send password reset link to user email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body service.ForgotPasswordRequest true "Email information"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /auth/forgot-password [post]
func (ac *AuthController) ForgotPassword(c *gin.Context) {
	var req service.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request parameters: " + err.Error(),
		})
		return
	}

	// Call service layer
	message, err := ac.authService.ForgotPassword(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}

// ResetPassword reset password
// @Summary Reset password
// @Description Set new password using reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body service.ResetPasswordRequest true "Reset information"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /auth/reset-password [post]
func (ac *AuthController) ResetPassword(c *gin.Context) {
	var req service.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request parameters: " + err.Error(),
		})
		return
	}

	// Call service layer
	if err := ac.authService.ResetPassword(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset successful",
	})
}

// ValidateToken validate token
// @Summary Validate token
// @Description Validate user token validity and return user information
// @Tags auth
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /auth/validate [get]
func (ac *AuthController) ValidateToken(c *gin.Context) {
	// Get token from request header
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Missing authentication token",
		})
		return
	}

	// Remove "Bearer " prefix if exists
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// Call service layer
	session, user, err := ac.authService.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":    user,
		"session": session,
		"valid":   true,
	})
}

// RegisterRoutes register routes
// @Description Register authentication-related routes to Gin router
func (ac *AuthController) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/signup", ac.Signup)
		auth.POST("/signin", ac.Signin)
		auth.POST("/logout", ac.Logout)
		auth.POST("/logout-all", ac.LogoutAll)
		auth.POST("/forgot-password", ac.ForgotPassword)
		auth.POST("/reset-password", ac.ResetPassword)
		auth.GET("/validate", ac.ValidateToken)
	}
}

