# 创建 JWT 认证中间件

## 任务描述

实现 JWT Token 生成、验证和刷新功能的 Gin 中间件。

## 详细说明

### 1. 添加 JWT 配置

在 `backend/internal/config/config.go` 添加：

```go
type Config struct {
    // ... existing fields
    JWTSecret          string  // JWT 密钥 (32字节)
    JWTExpireHours     int     // access_token 过期时间 (默认24)
    RefreshExpireDays   int     // refresh_token 过期时间 (默认7)
}
```

### 2. 创建 JWT 工具模块

创建 `backend/internal/auth/jwt.go`：

```go
package auth

import (
    "errors"
    "time"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

var (
    ErrInvalidToken = errors.New("invalid token")
    ErrExpiredToken = errors.New("token expired")
)

type Claims struct {
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

// GenerateAccessToken 生成 access_token (24小时)
func GenerateAccessToken(userID int64, username, role, secret string) (string, error) {
    claims := Claims{
        UserID:   userID,
        Username: username,
        Role:     role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

// GenerateRefreshToken 生成 refresh_token (7天)
func GenerateRefreshToken(userID int64, secret string) (string, error) {
    claims := jwt.RegisteredClaims{
        Subject:   fmt.Sprintf("%d", userID),
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
        IssuedAt:  jwt.NewNumericDate(time.Now()),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

// ValidateToken 验证 token 并返回 Claims
func ValidateToken(tokenString, secret string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, ErrInvalidToken
        }
        return []byte(secret), nil
    })
    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, ErrExpiredToken
        }
        return nil, ErrInvalidToken
    }
    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, ErrInvalidToken
    }
    return claims, nil
}

// ValidateRefreshToken 验证 refresh_token
func ValidateRefreshToken(tokenString, secret string) (int64, error) {
    claims, err := ValidateToken(tokenString, secret)
    if err != nil {
        return 0, err
    }
    userID, _ := strconv.ParseInt(claims.Subject, 10, 64)
    return userID, nil
}

// HashPassword 使用 bcrypt 哈希密码
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// CheckPassword 验证密码
func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

### 3. 创建 Gin 中间件

创建 `backend/internal/middleware/auth.go`：

```go
package middleware

import (
    "net/http"
    "strings"
    "node-pilot/internal/auth"
    "github.com/gin-gonic/gin"
)

func JWTAuth(secret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
            c.Abort()
            return
        }

        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
            c.Abort()
            return
        }

        claims, err := auth.ValidateToken(parts[1], secret)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
            c.Abort()
            return
        }

        // 将用户信息存入 Context
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("role", claims.Role)
        c.Next()
    }
}

func RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("role")
        if !exists {
            c.JSON(http.StatusForbidden, gin.H{"error": "role not found"})
            c.Abort()
            return
        }

        role := userRole.(string)
        for _, r := range roles {
            if role == r {
                c.Next()
                return
            }
        }

        c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
        c.Abort()
    }
}
```

### 4. 更新 Handler 依赖

修改 `backend/internal/handler/handler.go` 的 NewHandler 函数，添加 JWT 密钥参数。

## 输入

- 需求文档 `05-添加授权登录.md`
- JWT 标准库文档

## 输出

- 新建 `backend/internal/auth/jwt.go` - JWT 工具函数
- 新建 `backend/internal/middleware/auth.go` - Gin 中间件
- 修改 `backend/internal/config/config.go` - 添加 JWT 配置

## 依赖

- 01-user-model-and-migration.md

## 验收标准

- [ ] GenerateAccessToken 生成 24 小时有效的 JWT
- [ ] GenerateRefreshToken 生成 7 天有效的 refresh token
- [ ] ValidateToken 正确验证 token 并返回 Claims
- [ ] HashPassword 使用 bcrypt 生成密码哈希
- [ ] CheckPassword 正确验证密码
- [ ] JWTAuth 中间件正确提取和验证 Authorization header
- [ ] RequireRole 中间件正确检查用户角色
- [ ] token 失效时返回 401 状态码
