# Backend Handlers - Batch Script Support

## 任务描述

修改 Handler 层，支持批量脚本的创建和更新功能。

## 详细说明

### 1. 修改 CreateTaskInput 结构
在 `backend/internal/handler/handler.go` 中：

**原有结构：**
```go
type CreateTaskInput struct {
    ScriptID  int64   `json:"script_id"`
    Name      string  `json:"name"`
    ServerIDs []int64 `json:"server_ids"`
}
```

**修改为：**
```go
type CreateTaskInput struct {
    ScriptID  int64   `json:"script_id"`   // 保留兼容
    ScriptIDs string  `json:"script_ids"`  // 新字段：逗号分隔的脚本ID列表
    Name      string  `json:"name"`
    ServerIDs []int64 `json:"server_ids"`
}
```

### 2. 修改 CreateTask Handler
在 `CreateTask` 方法中：
- 优先使用 ScriptIDs 字段
- 如果 ScriptIDs 为空，则回退使用 ScriptID（单脚本兼容）
- 将脚本ID列表存储到 ScriptIDs 字段

**伪代码：**
```go
func (h *Handler) CreateTask(c *gin.Context) {
    var input CreateTaskInput
    // ... 解析代码 ...
    
    task := &model.Task{
        Name:     input.Name,
        Status:   "pending",
    }
    
    // 处理脚本ID
    if input.ScriptIDs != "" {
        task.ScriptIDs = input.ScriptIDs
    } else if input.ScriptID > 0 {
        task.ScriptIDs = strconv.FormatInt(input.ScriptID, 10)
    }
    
    // ... 创建任务逻辑 ...
}
```

### 3. 修改 UpdateTaskInput 结构
同样修改 UpdateTaskInput 结构体，添加 ScriptIDs 字段。

### 4. 修改 UpdateTask Handler
在 `UpdateTask` 方法中：
- 优先使用 ScriptIDs 字段更新
- 如果 ScriptIDs 变化，需要重新验证脚本存在性

### 5. 修改 GetTask Response
确保 API 返回时包含 ScriptIDs：

```go
func (h *Handler) GetTask(c *gin.Context) {
    // ...
    c.JSON(200, gin.H{
        "task":      task,
        "servers":   taskServers,
        "script_ids": task.ScriptIDs,  // 添加返回
    })
}
```

## 输入

- 现有 Handler 定义 (backend/internal/handler/handler.go)
- CreateTaskInput, UpdateTaskInput 结构体

## 输出

- 修改后的 handler.go，包含批量脚本支持的输入结构

## 依赖

- 01（需要 Task 模型支持 ScriptIDs 字段）

## 验收标准

- [ ] CreateTaskInput 包含 ScriptIDs 字段
- [ ] UpdateTaskInput 包含 ScriptIDs 字段
- [ ] CreateTask 优先使用 ScriptIDs，回退到 ScriptID
- [ ] UpdateTask 正确处理 ScriptIDs 更新
- [ ] GetTask API 返回包含 script_ids 字段
- [ ] 向后兼容：单脚本模式（使用 ScriptID）仍然正常工作
