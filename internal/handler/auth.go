package handler

import (
	"net/http"
	"lighthouse/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type RegisterWithEmailRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterWithPhoneRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"` // Can be email or phone
	Password   string `json:"password" binding:"required"`
}

// RegisterWithEmail godoc
// @Summary      Register a new user with email
// @Description  Create a new user account using email
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      RegisterWithEmailRequest  true  "Registration data"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]string
// @Router       /auth/register/email [post]
func (h *AuthHandler) RegisterWithEmail(c *gin.Context) {
	var req RegisterWithEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.authService.RegisterWithEmail(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":  user,
		"token": token,
	})
}

// RegisterWithPhone godoc
// @Summary      Register a new user with phone
// @Description  Create a new user account using phone number
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      RegisterWithPhoneRequest  true  "Registration data"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]string
// @Router       /auth/register/phone [post]
func (h *AuthHandler) RegisterWithPhone(c *gin.Context) {
	var req RegisterWithPhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.authService.RegisterWithPhone(req.Phone, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":  user,
		"token": token,
	})
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate user with email or phone and get JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      LoginRequest  true  "Login credentials (identifier can be email or phone)"
// @Success      200      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]string
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.authService.Login(req.Identifier, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	// TODO: Implement refresh token logic
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// Me godoc
// @Summary      Get current user
// @Description  Get authenticated user information
// @Tags         auth
// @Security     BearerAuth
// @Produce      json
// @Success      200      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]string
// @Router       /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.authService.GetUserByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
