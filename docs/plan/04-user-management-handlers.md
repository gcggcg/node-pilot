# 创建用户管理接口

## 任务描述

实现管理员操作用户的相关接口，包括用户列表、添加用户、删除用户。

## 详细说明

### 1. 创建 User 处理器

在 `backend/internal/handler/user.go`：

```go
package handler

import (
    "net/http"
    "strconv"
    "node-pilot/internal/auth"
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

type ListUsersRequest struct {
    Page    int    `form:"page"`
    PageSize int   `form:"pageSize"`
    Keyword string `form:"keyword"`
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
        users = []*repository.User{}
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

    // 检查用户名是否已存在
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

    user := &repository.User{
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

    // 禁止删除 root 用户
    for _, id := range req.IDs {
        user, err := h.repo.GetUserByID(id)
        if err == nil && user.Username == "root" {
            c.JSON(http.StatusForbidden, gin.H{"error": "cannot delete root user"})
            return
        }
    }

    if err := h.repo.DeleteUsers(req.IDs); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
```

### 2. 注册路由

在 `backend/cmd/server/main.go` 添加：

```go
userHandler := handler.NewUserHandler(repo, jwtSecret)

admin := api.Group("/v1/admin")
admin.Use(middleware.JWTAuth(jwtSecret))
admin.Use(middleware.RequireRole("ROLE_ADMIN"))
{
    admin.GET("/users", userHandler.ListUsers)
    admin.POST("/users", userHandler.CreateUser)
    admin.DELETE("/users/:id", userHandler.DeleteUsers)
    admin.POST("/users/batch-delete", userHandler.DeleteUsers)
}
```

### 3. 敏感操作日志

在 `backend/internal/handler/user.go` 的 DeleteUsers 中添加日志记录：

```go
// 记录删除操作日志
logger.Info("[AUDIT] User %s deleted users: %v", username, req.IDs)
```

## 输入

- 需求文档 `05-添加授权登录.md`
- 现有的 `backend/internal/handler/auth.go`
- 现有的 `backend/cmd/server/main.go`

## 输出

- 新建 `backend/internal/handler/user.go` - 用户管理处理器
- 修改 `backend/cmd/server/main.go` - 添加用户管理路由

## 依赖

- 01-user-model-and-migration.md
- 02-jwt-auth-middleware.md
- 03-auth-handlers.md

## 验收标准

- [ ] GET /api/v1/admin/users 返回分页用户列表
- [ ] GET /api/v1/admin/users 支持 keyword 模糊搜索
- [ ] POST /api/v1/admin/users 正确创建用户
- [ ] POST /api/v1/admin/users 检查用户名唯一性
- [ ] DELETE /api/v1/admin/users/:id 删除单个用户
- [ ] POST /api/v1/admin/users/batch-delete 批量删除用户
- [ ] 禁止删除 root 用户
- [ ] 所有管理接口需要 ROLE_ADMIN 权限
- [ ] 敏感操作记录日志
