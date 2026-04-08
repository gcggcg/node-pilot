# Repository Methods - Batch Script Operations

## 任务描述

在 Repository 层添加获取多个脚本和批量脚本操作的方法，支持 TaskExecutor 的批量脚本执行需求。

## 详细说明

### 1. 添加 GetScripts 方法
在 `backend/internal/repository/db.go` 中添加：

```go
// GetScripts 根据ID列表获取多个脚本
func (r *Repository) GetScripts(ids []int64) ([]*model.Script, error)
```

**实现逻辑：**
- 将 ids 转换为逗号分隔的字符串
- 执行 SQL: `SELECT * FROM scripts WHERE id IN (?,?,?)`
- 返回脚本列表

### 2. 添加 ParseScriptIDs 辅助函数
在 service 层添加解析函数（用于 TaskExecutor）：

```go
// ParseScriptIDs 解析逗号分隔的脚本ID字符串
func ParseScriptIDs(scriptIDs string) ([]int64, error)
```

**实现逻辑：**
- 如果为空字符串，返回空数组（向后兼容单脚本模式）
- 按逗号分割字符串
- 转换为 int64 数组
- 返回解析结果

### 3. 添加 UpdateTaskScripts 方法（如需要）
如果任务创建时需要单独保存脚本关联，可以添加：

```go
// UpdateTaskScripts 更新任务的脚本列表
func (r *Repository) UpdateTaskScripts(taskID int64, scriptIDs []int64) error
```

### 4. 修改 GetTask 方法
确保 GetTask 返回时包含 ScriptIDs 字段：

```go
// GetTask 获取任务（已包含 ScriptIDs 字段）
func (r *Repository) GetTask(id int64) (*model.Task, error)
```

## 输入

- 现有 Repository 定义 (backend/internal/repository/db.go)
- model.Script 结构定义

## 输出

- db.go 中新增的 GetScripts 方法
- service 层新增的 ParseScriptIDs 函数

## 依赖

- 01（需要 Task 模型已添加 ScriptIDs 字段）

## 验收标准

- [ ] GetScripts(ids []int64) 方法正确返回多个脚本
- [ ] ParseScriptIDs 函数正确解析逗号分隔的ID字符串
- [ ] 解析空字符串时返回空数组（不报错）
- [ ] GetTask 方法正确返回 ScriptIDs 字段
