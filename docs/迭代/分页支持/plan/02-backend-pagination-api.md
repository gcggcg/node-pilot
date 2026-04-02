# 后端分页接口改造

## 任务描述

修改后端 Go Handler，支持 page、pageSize 查询参数，返回 { data: [], total: number } 格式。

## 详细说明

### 1. 修改 Repository 层

在 `backend/internal/repository/db.go` 中为 ListServers、ListScripts、ListTasks 添加分页支持：

```go
// ListServersWithPagination 获取服务器列表（分页）
func (r *Repository) ListServersWithPagination(page, pageSize int) ([]*model.Server, int64, error) {
    offset := (page - 1) * pageSize
    
    var total int64
    err := r.db.QueryRow("SELECT COUNT(*) FROM servers").Scan(&total)
    if err != nil {
        return nil, 0, err
    }
    
    rows, err := r.db.Query(
        "SELECT id, name, host, port, username, password_encrypted, connection_status, created_at, updated_at FROM servers ORDER BY id DESC LIMIT ? OFFSET ?",
        pageSize, offset,
    )
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()
    
    var servers []*model.Server
    for rows.Next() {
        s := &model.Server{}
        err := rows.Scan(&s.ID, &s.Name, &s.Host, &s.Port, &s.Username, &s.PasswordEncrypted, &s.ConnectionStatus, &s.CreatedAt, &s.UpdatedAt)
        if err != nil {
            return nil, 0, err
        }
        servers = append(servers, s)
    }
    return servers, total, nil
}
```

类似实现 `ListScriptsWithPagination` 和 `ListTasksWithPagination`。

### 2. 修改 Handler 层

在 `backend/internal/handler/handler.go` 中：

```go
func (h *Handler) ListServers(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
    
    if page < 1 {
        page = 1
    }
    if pageSize < 1 || pageSize > 100 {
        pageSize = 10
    }
    
    servers, total, err := h.repo.ListServersWithPagination(page, pageSize)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    if servers == nil {
        servers = []*model.Server{}
    }
    c.JSON(200, gin.H{
        "data": servers,
        "total": total,
        "page": page,
        "pageSize": pageSize,
    })
}
```

类似修改 ListScripts、ListTasks。

## 输入

- `backend/internal/repository/db.go`
- `backend/internal/handler/handler.go`

## 输出

- 修改后的 Repository 和 Handler，支持分页参数

## 依赖

- 01

## 验收标准

- [ ] GET /api/servers 支持 page、pageSize 参数
- [ ] GET /api/scripts 支持 page、pageSize 参数
- [ ] GET /api/tasks 支持 page、pageSize 参数
- [ ] 返回格式包含 data 数组和 total 总数
