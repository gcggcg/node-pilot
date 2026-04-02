# 任务管理改造 - TaskExecutor Service

## 任务描述

在 TaskExecutor Service 中添加 `ExecuteTask` 方法，供手动执行任务使用。

## 详细说明

### 1. 新增 ExecuteTask 方法

在 `backend/internal/service/task.go` 中添加：

```go
// ExecuteTask 手动执行指定任务
func (s *TaskExecutor) ExecuteTask(taskID int64) error {
    task, err := s.repo.GetTask(taskID)
    if err != nil {
        return err
    }
    
    // 只能执行 pending 状态的任务
    if task.Status != "pending" {
        return fmt.Errorf("任务不是pending状态，无法执行")
    }
    
    servers, err := s.getTaskServers(task.ID)
    if err != nil {
        return err
    }
    
    if len(servers) == 0 {
        return fmt.Errorf("任务没有关联服务器")
    }
    
    // 获取第一个服务器的密码（假设所有服务器密码相同）
    // 注意：实际应该从 task_servers 或单独存储获取
    password, err := s.decryptPassword(servers[0].PasswordEncrypted)
    if err != nil {
        return err
    }
    
    script, err := s.repo.GetScript(task.ScriptID)
    if err != nil {
        return err
    }
    
    // 在 goroutine 中执行
    go s.ExecuteScript(task, script, servers, password)
    
    return nil
}
```

### 2. 修改 CreateTask 移除自动执行

在 `CreateTask` 调用后，不再启动 goroutine 执行。任务创建后保持 pending 状态。

### 3. 添加 decryptPassword 辅助方法

参考 FileUploadService 的实现，添加密码解密方法。

## 输入

- `backend/internal/service/task.go` 文件

## 输出

- 修改后的 `backend/internal/service/task.go`
- 新增 `ExecuteTask` 方法
- 新增 `decryptPassword` 辅助方法

## 依赖

- 01-backend-repository

## 验收标准

- [ ] `ExecuteTask` 方法正确执行 pending 状态的任务
- [ ] 非 pending 状态任务返回错误
- [ ] 任务在 goroutine 中异步执行
- [ ] 密码正确解密
