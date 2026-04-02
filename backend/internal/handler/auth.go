package handler

import (
	"net/http"

	"node-pilot/internal/auth"
	"node-pilot/internal/repository"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	repo      *repository.Repository
	jwtSecret string
}

func NewAuthHandler(repo *repository.Repository, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.repo.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	accessToken, err := auth.GenerateAccessToken(user.ID, user.Username, user.Role, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    86400,
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	user, err := h.repo.GetUserByID(userID.(int64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"phone":    user.Phone,
		"role":     user.Role,
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := auth.ValidateRefreshToken(req.RefreshToken, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	user, err := h.repo.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	accessToken, err := auth.GenerateAccessToken(user.ID, user.Username, user.Role, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    86400,
	})
}

type UpdateProfileRequest struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.repo.GetUserByID(userID.(int64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	user.Email = req.Email
	user.Phone = req.Phone

	if err := h.repo.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated"})
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.repo.GetUserByID(userID.(int64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if !auth.CheckPassword(req.OldPassword, user.PasswordHash) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect old password"})
		return
	}

	newHash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user.PasswordHash = newHash
	if err := h.repo.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed"})
}
