# 创建用户模型与数据库迁移

## 任务描述

在 Node-Pilot 后端添加用户认证模块，创建 User 模型和 users 表数据库迁移。

## 详细说明

### 1. 添加 User 模型

在 `backend/internal/model/model.go` 添加 User 结构体：

```go
type User struct {
    ID           int64     `json:"id"`
    Username     string    `json:"username"`
    PasswordHash string    `json:"-"`           // bcrypt hash, never expose
    Email        string    `json:"email"`
    Phone        string    `json:"phone"`
    Role         string    `json:"role"`        // ROLE_ADMIN or ROLE_USER
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
```

### 2. 添加 users 表迁移

在 `backend/internal/repository/db.go` 的 `initSchema` 函数中添加：

```go
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    email TEXT DEFAULT '',
    phone TEXT DEFAULT '',
    role TEXT DEFAULT 'ROLE_USER',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### 3. 添加 User Repository 方法

在 `backend/internal/repository/db.go` 添加：

- `CreateUser(user *User) (int64, error)` - 创建用户
- `GetUserByUsername(username string) (*User, error)` - 根据用户名查询
- `GetUserByID(id int64) (*User, error)` - 根据ID查询
- `UpdateUser(user *User) error` - 更新用户
- `ListUsers(page, pageSize int, keyword string) ([]*User, int64, error)` - 分页查询用户
- `DeleteUsers(ids []int64) error` - 批量删除用户

### 4. 默认 root 用户初始化

在 `NewDB` 函数或 `initSchema` 末尾添加 root 用户创建逻辑：

```go
// 检查 root 用户是否存在，不存在则创建
var count int64
r.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = 'root'").Scan(&count)
if count == 0 {
    // password: root, 使用 bcrypt hash
    hash, _ := bcrypt.GenerateFromPassword([]byte("root"), bcrypt.DefaultCost)
    r.db.Exec("INSERT INTO users (username, password_hash, role) VALUES ('root', ?, 'ROLE_ADMIN')", string(hash))
}
```

### 5. 依赖项

确保 `go.mod` 包含：

```
golang.org/x/crypto/bcrypt
github.com/golang-jwt/jwt/v5
```

## 输入

- 需求文档 `05-添加授权登录.md`
- 现有 `backend/internal/model/model.go`
- 现有 `backend/internal/repository/db.go`

## 输出

- 修改 `backend/internal/model/model.go` - 添加 User 结构体
- 修改 `backend/internal/repository/db.go` - 添加 users 表和 Repository 方法

## 依赖

- 无

## 验收标准

- [ ] User 结构体包含 id, username, password_hash, email, phone, role, created_at, updated_at
- [ ] users 表正确创建，包含唯一约束和默认角色
- [ ] CreateUser 方法使用 bcrypt 加密密码
- [ ] GetUserByUsername 方法用于登录验证
- [ ] ListUsers 支持分页和关键字搜索
- [ ] 系统启动时自动创建 root/root 管理员用户
