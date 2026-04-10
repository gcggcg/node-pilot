# 后端 Handler 输入结构体修改

## 任务描述

在 CreateTaskInput 和 UpdateTaskInput 结构体中添加 `execution_mode` 字段，支持创建和更新任务时指定执行模式。

## 详细说明

### 1. 修改 CreateTaskInput 结构体 (`backend/internal/handler/handler.go`)

在现有结构体中添加 `ExecutionMode` 字段：
```go
type CreateTaskInput struct {
    ScriptID     int64   `json:"script_id"`
    ScriptIDs    string  `json:"script_ids"`
    Name         string  `json:"name"`
    ServerIDs    []int64 `json:"server_ids"`
    ExecutionMode string  `json:"execution_mode"` // 新增
}
```

### 2. 修改 UpdateTaskInput 结构体 (`backend/internal/handler/handler.go`)

在现有结构体中添加 `ExecutionMode` 字段：
```go
type UpdateTaskInput struct {
    ScriptID     int64   `json:"script_id"`
    ScriptIDs    string  `json:"script_ids"`
    Name         string  `json:"name"`
    ServerIDs    []int64 `json:"server_ids"`
    ExecutionMode string  `json:"execution_mode"` // 新增
}
```

### 3. 修改 CreateTask 函数 (`backend/internal/handler/handler.go`)

在创建任务时，将 `ExecutionMode` 保存到数据库：
- 如果 `input.ExecutionMode` 为空，默认为 `"concurrent"`
- 只接受 `"concurrent"` 或 `"sequential"` 两个值

### 4. 修改 UpdateTask 函数 (`backend/internal/handler/handler.go`)

在更新任务时，支持更新 `ExecutionMode`：
- 如果 `input.ExecutionMode` 非空，则更新任务执行模式

### 5. 修改 Repository UpdateTask 方法 (`backend/internal/repository/db.go`)

可能需要新增或修改 `UpdateTask` 方法签名以支持 `execution_mode` 更新，或新增 `UpdateTaskExecutionMode` 方法。

## 输入

- 需求文档：`10-多服务器任务支持单线程和多线程切换.md`
- 现有 Handler 文件：`backend/internal/handler/handler.go`
- 现有 Repository 文件：`backend/internal/repository/db.go`
- Task 01 输出：已添加 execution_mode 字段的 model

## 输出

- 修改后的 `backend/internal/handler/handler.go`
- 修改后的 `backend/internal/repository/db.go`（如需要）

## 依赖

- 01

## 验收标准

- [ ] CreateTaskInput 包含 `ExecutionMode string` 字段
- [ ] UpdateTaskInput 包含 `ExecutionMode string` 字段
- [ ] 创建任务时正确保存 execution_mode（默认 "concurrent"）
- [ ] 更新任务时正确更新 execution_mode
- [ ] 无效的 execution_mode 值被拒绝（只接受 "concurrent" 或 "sequential"）
