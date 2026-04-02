# 创建认证接口处理器

## 任务描述

实现登录、获取当前用户、刷新 Token 等认证相关接口。

## 详细说明

### 1. 创建 Auth 处理器

在 `backend/internal/handler/auth.go`：

```go
package handler

import (
    "net/http"
    "node-pilot/internal/auth"
    "node-pilot/internal/repository"
    "github.com/gin-gonic/gin"
)

type AuthHandler struct {
    repo     *repository.Repository
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
    ExpiresIn    int    `json:"expires_in"` // 秒
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
        ExpiresIn:    86400, // 24小时
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
```

### 2. 注册路由

在 `backend/cmd/server/main.go` 添加：

```go
// 导入新模块
"node-pilot/internal/auth"
"node-pilot/internal/middleware"

// 在 main 函数中初始化
jwtSecret := "your-32-byte-secret-key-here-123" // 生产环境从配置读取
authHandler := handler.NewAuthHandler(repo, jwtSecret)

// 在路由配置中添加
auth := api.Group("/v1/auth")
{
    auth.POST("/login", authHandler.Login)
    auth.GET("/me", middleware.JWTAuth(jwtSecret), authHandler.Me)
    auth.POST("/refresh", authHandler.RefreshToken)
    auth.PUT("/profile", middleware.JWTAuth(jwtSecret), authHandler.UpdateProfile)
    auth.PUT("/password", middleware.JWTAuth(jwtSecret), authHandler.ChangePassword)
}
```

## 输入

- 需求文档 `05-添加授权登录.md`
- 现有的 `backend/internal/handler/handler.go`
- 现有的 `backend/cmd/server/main.go`

## 输出

- 新建 `backend/internal/handler/auth.go` - 认证处理器
- 修改 `backend/cmd/server/main.go` - 添加认证路由

## 依赖

- 01-user-model-and-migration.md
- 02-jwt-auth-middleware.md

## 验收标准

- [ ] POST /api/v1/auth/login 正确验证用户名密码并返回 token
- [ ] GET /api/v1/auth/me 返回当前登录用户信息
- [ ] POST /api/v1/auth/refresh 使用 refresh_token 获取新 access_token
- [ ] PUT /api/v1/auth/profile 更新用户邮箱和电话
- [ ] PUT /api/v1/auth/password 验证旧密码后更新新密码
- [ ] 所有需要认证的接口正确使用 JWTAuth 中间件
