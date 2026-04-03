# 风险评估与运行效果报告 - node-pilot

## 执行摘要

本次代码审查覆盖从 `1a6dcf6` 到 `HEAD` 的3个提交，主要涉及任务管理改造功能。代码整体结构良好，**之前的重复执行bug已修复**
。但仍存在**1个严重安全问题**和若干代码质量问题需要关注。

**审查范围:**

- 提交: `e61795a`, `877af9d`, `b4a2de4`
- 主要变更: 任务管理改造（创建/执行分离，编辑功能）

---

## 🔧 自动修复汇总

### ✅ 已自动修复 (Auto-Fixed)

| 问题          | 位置              | 修复方式                                      |
|-------------|-----------------|-------------------------------------------|
| 任务执行后刷新显示失败 | task.go:261-280 | 修复 executeOnServer 重复创建 task_servers 记录问题 |
| 输出显示重复      | task.go:261-280 | 使用 UpdateTaskServerByIDs 更新已有记录而非创建新记录    |

---

## 🔴 严重风险 (Critical)

### 1. Task/Server/Script 路由缺少JWT认证

**位置:** `backend/cmd/server/main.go:136-146`

```go
tasks := api.Group("/tasks")
{
tasks.GET("", h.ListTasks)
tasks.GET("/:id", h.GetTask)
tasks.POST("", h.CreateTask)
tasks.PUT("/:id", h.UpdateTask)
tasks.POST("/:id/execute", h.ExecuteTask)
tasks.DELETE("/:id", h.CancelTask)
tasks.POST("/batch-delete", h.DeleteTasks)
tasks.GET("/:id/output", h.GetTaskOutput)
}
```

**问题:**

- `/tasks`, `/servers`, `/scripts` 路由**均未使用 JWT 认证中间件**
- 对比 `/v1/file-uploads` 正确使用了 `middleware.JWTAuth(jwtSecret)`

**影响:**

- **任何人无需登录即可操作用户任务、服务器、脚本**
- 可执行任意任务、查看输出、删除数据

**修复建议:**

```go
tasks := api.Group("/tasks")
tasks.Use(middleware.JWTAuth(jwtSecret)) // 添加认证
{
// ...
}
```

**严重程度:** 🔴 Critical  
**是否可自动修复:** 否，需修改路由配置

---

## 🟠 高风险 (High)

### 2. 重复的密码解密函数

**位置:**

- `backend/internal/handler/handler.go:635-662` (`decrypt` 方法)
- `backend/internal/service/task.go:105-130` (`decryptPassword` 方法)

**问题:** 两处代码几乎完全相同，违反DRY原则

```go
// handler.go
func (h *Handler) decrypt(encrypted string) (string, error) {
key := []byte("12345678901234567890123456789012")
// ... AES-256-GCM 解密逻辑
}

// task.go  
func (e *TaskExecutor) decryptPassword(encrypted string) (string, error) {
key := []byte("12345678901234567890123456789012") // 同样的硬编码密钥
// ... 完全相同的 AES-256-GCM 解密逻辑
}
```

**影响:**

- 代码维护困难
- 密钥分散在多处

**修复建议:** 抽取为公共函数，放入 `common/` 或 `crypto/` 包

**严重程度:** 🟠 High  
**是否可自动修复:** ⚠️ 可合并，但需验证不影响现有功能

---

## 🟡 中风险 (Medium)

### 3. 硬编码AES密钥

**位置:**

- `backend/internal/service/task.go:110`
- `backend/internal/handler/handler.go:43`

```go
key := []byte("12345678901234567890123456789012") // 32 bytes for AES-256
```

**问题:** 生产环境密钥不应硬编码

**严重程度:** 🟡 Medium  
**建议:** 使用环境变量或配置文件

---

### 4. 硬编码JWT密钥

**位置:** `backend/cmd/server/main.go:77`

```go
jwtSecret := "node-pilot-jwt-secret-key-32bytes!" // 32 bytes for HS256
```

**严重程度:** 🟡 Medium  
**建议:** 使用环境变量 `JWT_SECRET`

---

### 5. Task name 缺乏验证

**位置:** `backend/internal/handler/handler.go:444`

```go
type UpdateTaskInput struct {
ScriptID  int64   `json:"script_id"`
Name      string  `json:"name"` // 无长度/内容验证
ServerIDs []int64 `json:"server_ids"`
}
```

**问题:** 任务名称无验证，可能导致:

- 数据库存储过长字符串
- XSS风险（如果显示在Web界面）

**严重程度:** 🟡 Medium  
**建议:** 添加长度限制和HTML转义

---

## 🟢 低风险 (Low)

### 6. 注释中提及不应发生的fallback逻辑

**位置:** `backend/internal/service/task.go:277-284`

```go
} else {
// 如果没有找到已有记录（理论上不应该发生），创建新记录
ts := &model.TaskServer{
TaskID:   task.ID,
ServerID: srv.ID,
Status:   "running",
}
tsID, _ = e.repo.CreateTaskServer(ts)
}
```

**问题:** 虽然注释说"不应该发生"，但代码仍然创建记录，可能导致bug被掩盖

**严重程度:** 🟢 Low  
**建议:** 改为返回错误并记录日志

---

### 7. 前端使用 alert() 而非统一错误处理

**位置:** `frontend/src/views/TaskForm.vue:398`

```typescript
alert('加载任务失败');
```

**严重程度:** 🟢 Low  
**建议:** 使用统一的错误提示组件

---

## 已识别风险汇总

| #   | 风险类型         | 严重程度        | 位置                  | 状态     |
|-----|--------------|-------------|---------------------|--------|
| 1   | 路由缺少JWT认证    | 🔴 Critical | main.go:136-146     | ⚠️ 待修复 |
| 2   | 重复解密函数       | 🟠 High     | task.go, handler.go | ⚠️ 待修复 |
| 3   | 硬编码AES密钥     | 🟡 Medium   | task.go:110         | ⚠️ 待修复 |
| 4   | 硬编码JWT密钥     | 🟡 Medium   | main.go:77          | ⚠️ 待修复 |
| 5   | 任务名称无验证      | 🟡 Medium   | handler.go:444      | ⚠️ 待修复 |
| 6   | Fallback创建记录 | 🟢 Low      | task.go:277         | ⚠️ 待修复 |
| 7   | 使用alert()    | 🟢 Low      | TaskForm.vue        | ⚠️ 待修复 |

---

## 代码质量评分

| 维度   | 评分     | 说明           |
|------|--------|--------------|
| 可维护性 | 75/100 | 重复代码需重构      |
| 安全性  | 60/100 | JWT认证缺失是严重问题 |
| 性能   | 90/100 | 无明显性能问题      |

---

## 运行效果验证

```bash
# Backend构建
cd backend && go build -o /dev/null ./cmd/server
✅ 构建成功

# Frontend构建
cd frontend && npm run build
✅ 构建成功
```

---

## 修复优先级建议

1. **立即修复** - 添加JWT认证到 /tasks, /servers, /scripts 路由
2. **尽快修复** - 合并重复的 decryptPassword 函数
3. **计划修复** - 将密钥移至环境变量
4. **后续优化** - 添加输入验证和错误处理改进

---

## 结论

本次变更 (`b4a2de4`) 成功修复了任务重复执行的bug，任务管理改造功能实现正确。但**安全认证缺失**
是需要立即处理的关键问题，建议参考 `fileUploads` 路由的认证模式进行修复。

**审查完成时间:** 2026-04-02
