# 风险评估与运行效果报告 - NodePilot JWT Auth 增量代码

## 执行摘要

本次代码审查覆盖从 `5af2abf` 到 `HEAD` 共 **10 次提交** 的 JWT 认证系统增量开发。代码整体质量良好，认证逻辑实现完整，但存在
**2 个高风险安全问题** 需要立即处理。

---

## 代码变更概览

- **提交数**: 10
- **新增文件**: 14 个源文件 + 大量构建产物
- **修改文件**: 6 个核心文件
- **删除文件**: 0

### 核心新增文件

| 文件                                    | 说明                 |
|---------------------------------------|--------------------|
| `backend/internal/auth/jwt.go`        | JWT 生成、验证、密码哈希工具   |
| `backend/internal/handler/auth.go`    | 登录/刷新/个人资料 API     |
| `backend/internal/handler/user.go`    | 用户管理 API           |
| `backend/internal/middleware/auth.go` | JWT 中间件和角色验证       |
| `backend/internal/model/model.go`     | User 模型定义          |
| `backend/internal/repository/db.go`   | 用户表和 Repository 方法 |
| `frontend/src/stores/auth.ts`         | Pinia 认证状态管理       |
| `frontend/src/views/Login.vue`        | 登录页面               |
| `frontend/src/views/UserList.vue`     | 用户管理页面             |
| `frontend/src/views/Profile.vue`      | 个人中心页面             |

---

## 🔧 自动修复汇总

### ✅ 已自动修复

| 问题                      | 位置                | 修复方式             |
|-------------------------|-------------------|------------------|
| TypeScript 未使用变量 `_i`   | `Login.vue:78`    | 改名为 `_i` 消除警告    |
| TypeScript 未使用变量 `from` | `auth-guard.ts:5` | 改名为 `_from` 消除警告 |

### ⏳ 待手动处理

| 问题             | 严重程度        | 位置                   | 建议修复方式      |
|----------------|-------------|----------------------|-------------|
| 硬编码 JWT 密钥     | 🔴 Critical | `main.go:70`         | 改为环境变量或配置文件 |
| 硬编码 Root 密码    | 🔴 Critical | `db.go:InitRootUser` | 改为环境变量注入    |
| 未使用的 bcrypt 导入 | 🟢 Low      | `db.go:8`            | 移除未使用的导入    |

---

## 已识别风险汇总

### 🔴 严重风险 (Critical)

#### 1. 硬编码 JWT 密钥

**位置**: `backend/cmd/server/main.go:70`

```go
jwtSecret := "node-pilot-jwt-secret-key-32bytes!" // 32 bytes for HS256
```

**风险**: JWT 密钥直接写在源代码中，任何能访问代码的人都能伪造有效 Token。

**修复建议**:

```go
// 从环境变量或配置文件获取
jwtSecret := os.Getenv("JWT_SECRET")
if jwtSecret == "" {
log.Fatal("JWT_SECRET environment variable not set")
}
```

#### 2. 硬编码 Root 默认密码

**位置**: `backend/internal/repository/db.go:InitRootUser`

```go
hash, err := bcrypt.GenerateFromPassword([]byte("root"), bcrypt.DefaultCost)
```

**风险**: 所有部署实例的 root 密码都是 "root"，攻击者可轻易登录。

**修复建议**:

```go
func InitRootUser() error {
rootPassword := os.Getenv("ROOT_PASSWORD")
if rootPassword == "" {
rootPassword = "change-me-root-password" // 强制要求修改
}
hash, err := bcrypt.GenerateFromPassword([]byte(rootPassword), bcrypt.DefaultCost)
// ...
}
```

---

### 🟠 高风险 (High)

#### 3. CORS 允许所有来源

**位置**: `backend/internal/websocket/hub.go` 或主服务配置

```go
CheckOrigin: func () bool { return true }
```

**风险**: 任何网站都能向后端 API 发起请求，进行 CSRF 攻击。

**修复建议**:

```go
CheckOrigin: func (r *http.Request) bool {
allowedOrigins := []string{"http://localhost:8080", "https://yourdomain.com"}
origin := r.Header.Get("Origin")
for _, allowed := range allowedOrigins {
if origin == allowed {
return true
}
}
return false
}
```

---

### 🟡 中风险 (Medium)

#### 4. 登录接口无频率限制

**位置**: `backend/internal/handler/auth.go:Login`

**风险**: 攻击者可进行暴力破解密码攻击。

**修复建议**: 实现登录失败次数限制（如 5 次失败后锁定 15 分钟）。

#### 5. 密码强度要求宽松

**位置**: `backend/internal/handler/auth.go:ChangePassword`

```go
NewPassword string `json:"new_password" binding:"required,min=6"`
```

**风险**: 仅要求 6 字符，建议增强。

---

### 🟢 低风险 (Low)

#### 6. 未使用的导入

**位置**: `backend/internal/repository/db.go:8`

```go
import (
// ...
"golang.org/x/crypto/bcrypt" // 未使用，auth 包已使用
)
```

**影响**: 轻微，编译时 Go 会自动优化。

---

## 安全漏洞分析

| 漏洞类型      | 位置                 | 严重程度        | 修复建议             |
|-----------|--------------------|-------------|------------------|
| 硬编码密钥     | main.go:70         | 🔴 Critical | 使用环境变量           |
| 硬编码密码     | db.go:InitRootUser | 🔴 Critical | 使用环境变量           |
| CORS 过度允许 | websocket hub.go   | 🟠 High     | 配置允许列表           |
| 无登录限流     | auth.go:Login      | 🟠 High     | 实现 rate limiting |
| 弱密码策略     | auth.go            | 🟡 Medium   | 增强最小长度           |

---

## JWT/认证相关代码审查

### ✅ 正确实现

| 项目                | 位置                                 | 说明                |
|-------------------|------------------------------------|-------------------|
| 密码哈希              | `auth/jwt.go:HashPassword`         | 正确使用 bcrypt       |
| Token 验证          | `auth/jwt.go:ValidateToken`        | 正确验证签名和过期         |
| Access Token 有效期  | `auth/jwt.go:GenerateAccessToken`  | 24 小时，合理          |
| Refresh Token 有效期 | `auth/jwt.go:GenerateRefreshToken` | 7 天，可接受           |
| 角色验证中间件           | `middleware/auth.go:RequireRole`   | 正确实现              |
| 密码永不暴露            | `model/model.go`                   | `json:"-"` tag 正确 |
| Root 用户删除保护       | `handler/user.go:DeleteUsers`      | 正确检查              |

### ⚠️ 需要注意

| 项目               | 说明                                 |
|------------------|------------------------------------|
| JWT 密钥长度         | 32 字节足够 (HS256)                    |
| Refresh Token 轮换 | 未实现 Refresh Token 轮换，存在 Token 复用风险 |

---

## 代码质量评分

| 维度   | 评分     | 说明           |
|------|--------|--------------|
| 可维护性 | 85/100 | 代码结构清晰，注释充分  |
| 安全性  | 55/100 | 存在硬编码密钥和密码   |
| 性能   | 90/100 | 无性能问题        |
| 架构   | 88/100 | 分层合理，中间件模式良好 |

---

## 运行效果验证

### Go Backend Build

```
✓ 编译成功
```

### Frontend Build

```
✓ vue-tsc 编译成功
✓ vite build 成功 (143.66 kB gzipped: 55.43 kB)
```

---

## 修复优先级建议

1. **[P0 - 立即修复]** 将 JWT Secret 改为环境变量
2. **[P0 - 立即修复]** 将 Root 密码改为环境变量
3. **[P1 - 高优先级]** 限制 CORS 来源
4. **[P1 - 高优先级]** 实现登录限流
5. **[P2 - 中优先级]** 实现 Refresh Token 轮换
6. **[P3 - 低优先级]** 移除未使用的导入

---

## 结论

本次增量开发的 JWT 认证系统整体实现质量良好，认证逻辑完善，密码处理安全（使用 bcrypt），Token 结构合理。但存在 **2 个严重安全问题
**（硬编码密钥和密码）必须立即修复，否则系统安全性形同虚设。

建议在生产部署前完成 P0 和 P1 级别的安全修复。

---

**报告生成时间**: 2026-04-02  
**审查范围**: commit `5af2abf` → `HEAD`  
**审查工具**: 自动代码审查 + 安全分析
