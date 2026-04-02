# 任务管理改造 - Backend Repository

## 任务描述

改造任务管理功能，实现"创建任务不自动执行，需要手动点击执行"的工作流。Backend层需要：
1. 新增批量创建任务服务器关联方法 `CreateTaskServers`
2. 新增删除任务服务器关联方法 `DeleteTaskServers`
3. 修改现有 `CreateTask` 不再自动执行

## 详细说明

### 1. 新增 Repository 方法

在 `backend/internal/repository/db.go` 中添加：

```go
// CreateTaskServers 批量创建任务-服务器关联
func (r *Repository) CreateTaskServers(taskID int64, serverIDs []int64) error {
    if len(serverIDs) == 0 {
        return nil
    }
    tx, err := r.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    stmt, err := tx.Prepare(`INSERT INTO task_servers (task_id, server_id, status) VALUES (?, ?, 'pending')`)
    if err != nil {
        return err
    }
    defer stmt.Close()
    
    for _, serverID := range serverIDs {
        _, err := stmt.Exec(taskID, serverID)
        if err != nil {
            return err
        }
    }
    
    return tx.Commit()
}

// DeleteTaskServers 删除任务的所有服务器关联
func (r *Repository) DeleteTaskServers(taskID int64) error {
    _, err := r.db.Exec(`DELETE FROM task_servers WHERE task_id = ?`, taskID)
    return err
}
```

### 2. 修改 CreateTask 不自动执行

`CreateTask` 方法需要：
- 创建任务记录（status = "pending"）
- 创建任务-服务器关联记录（使用新的 CreateTaskServers）
- **不再启动 goroutine 自动执行**

## 输入

- `backend/internal/repository/db.go` 文件

## 输出

- 修改后的 `backend/internal/repository/db.go`
- 新增 `CreateTaskServers` 方法
- 新增 `DeleteTaskServers` 方法

## 依赖

- 无

## 验收标准

- [ ] `CreateTaskServers` 方法正确批量创建任务-服务器关联
- [ ] `DeleteTaskServers` 方法正确删除任务-服务器关联
- [ ] `CreateTask` 不再自动执行任务（移除 goroutine 代码）
