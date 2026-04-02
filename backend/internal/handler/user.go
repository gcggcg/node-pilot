package handler

import (
	"net/http"
	"strconv"

	"node-pilot/internal/auth"
	"node-pilot/internal/logger"
	"node-pilot/internal/model"
	"node-pilot/internal/repository"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	repo      *repository.Repository
	jwtSecret string
}

func NewUserHandler(repo *repository.Repository, jwtSecret string) *UserHandler {
	return &UserHandler{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, total, err := h.repo.ListUsers(page, pageSize, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if users == nil {
		users = []*model.User{}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     users,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role" binding:"required,oneof=ROLE_USER ROLE_ADMIN"`
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existing, _ := h.repo.GetUserByUsername(req.Username)
	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user := &model.User{
		Username:     req.Username,
		PasswordHash: hash,
		Email:        req.Email,
		Phone:        req.Phone,
		Role:         req.Role,
	}

	id, err := h.repo.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

type DeleteUsersRequest struct {
	IDs []int64 `json:"ids" binding:"required"`
}

func (h *UserHandler) DeleteUsers(c *gin.Context) {
	var req DeleteUsersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, id := range req.IDs {
		user, err := h.repo.GetUserByID(id)
		if err == nil && user.Username == "root" {
			c.JSON(http.StatusForbidden, gin.H{"error": "cannot delete root user"})
			return
		}
	}

	username, _ := c.Get("username")
	logger.Info("[AUDIT] User %s deleted users: %v", username, req.IDs)

	if err := h.repo.DeleteUsers(req.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
