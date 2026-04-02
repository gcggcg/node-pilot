# 任务管理改造 - Backend Handlers & Routes

## 任务描述

修改 Backend Handlers 和 Routes，实现：
1. 修改 `CreateTask` 不自动执行
2. 新增 `ExecuteTask` handler
3. 新增 `UpdateTask` handler
4. 新增相关路由

## 详细说明

### 1. 修改 CreateTask Handler

在 `handler.go` 中修改 `CreateTask`：

```go
func (h *Handler) CreateTask(c *gin.Context) {
    var input CreateTaskInput
    if err := c.BindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    task := &model.Task{
        ScriptID: input.ScriptID,
        Name:     input.Name,
        Status:   "pending",
    }

    id, err := h.repo.CreateTask(task)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    // 创建任务-服务器关联（批量）
    if err := h.repo.CreateTaskServers(id, input.ServerIDs); err != nil {
        h.repo.DeleteTasks([]int64{id})
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    // 移除自动执行逻辑 - 现在需要手动调用 ExecuteTask
    // go func() { ... }()

    c.JSON(201, gin.H{"id": id})
}
```

### 2. 新增 ExecuteTask Handler

```go
func (h *Handler) ExecuteTask(c *gin.Context) {
    idStr := c.Param("id")
    id, _ := strconv.ParseInt(idStr, 10, 64)
    
    task, err := h.repo.GetTask(id)
    if err != nil {
        c.JSON(404, gin.H{"error": "task not found"})
        return
    }
    
    if task.Status != "pending" {
        c.JSON(400, gin.H{"error": "only pending tasks can be executed"})
        return
    }
    
    if err := h.taskSvc.ExecuteTask(id); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{"message": "task execution started"})
}
```

### 3. 新增 UpdateTask Handler

```go
type UpdateTaskInput struct {
    ScriptID  int64   `json:"script_id"`
    Name      string  `json:"name"`
    ServerIDs []int64 `json:"server_ids"`
}

func (h *Handler) UpdateTask(c *gin.Context) {
    idStr := c.Param("id")
    id, _ := strconv.ParseInt(idStr, 10, 64)
    
    task, err := h.repo.GetTask(id)
    if err != nil {
        c.JSON(404, gin.H{"error": "task not found"})
        return
    }
    
    // 只能修改 pending 状态的任务
    if task.Status != "pending" {
        c.JSON(400, gin.H{"error": "only pending tasks can be modified"})
        return
    }
    
    var input UpdateTaskInput
    if err := c.BindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 更新任务基本信息
    task.Name = input.Name
    task.ScriptID = input.ScriptID
    
    // 删除旧的服务器关联
    if err := h.repo.DeleteTaskServers(id); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    // 创建新的服务器关联
    if err := h.repo.CreateTaskServers(id, input.ServerIDs); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{"message": "task updated"})
}
```

### 4. 新增路由

在 `main.go` 中添加：

```go
tasks := api.Group("/tasks")
{
    tasks.GET("", h.ListTasks)
    tasks.GET("/:id", h.GetTask)
    tasks.POST("", h.CreateTask)
    tasks.PUT("/:id", h.UpdateTask)           // 新增
    tasks.POST("/:id/execute", h.ExecuteTask) // 新增
    tasks.DELETE("/:id", h.CancelTask)
    tasks.POST("/batch-delete", h.DeleteTasks)
    tasks.GET("/:id/output", h.GetTaskOutput)
}
```

## 输入

- `backend/internal/handler/handler.go`
- `backend/cmd/server/main.go`

## 输出

- 修改后的 `handler.go`（CreateTask 移除自动执行，新增 ExecuteTask, UpdateTask）
- 修改后的 `main.go`（新增路由）

## 依赖

- 01-backend-repository
- 02-task-executor-service

## 验收标准

- [ ] CreateTask 不再自动执行
- [ ] ExecuteTask handler 正确触发任务执行
- [ ] UpdateTask handler 只能修改 pending 状态任务
- [ ] 路由正确注册
