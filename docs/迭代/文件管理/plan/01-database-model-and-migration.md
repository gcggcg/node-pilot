# 文件管理 - 数据库模型与迁移

## 任务描述

为文件上传管理功能设计数据库模型和表结构。

## 详细说明

### 1. 新增 Model

在 `backend/internal/model/model.go` 中添加 FileUpload 和 FileUploadServer 结构体：

```go
type FileUpload struct {
    ID          int64     `json:"id"`
    Name        string    `json:"name"`        // 配置名称
    LocalPath   string    `json:"local_path"`   // ./data/files/ 下的相对路径
    RemotePath  string    `json:"remote_path"`  // 远程目标路径，必须以/开头
    Status      string    `json:"status"`      // pending|success|failed
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type FileUploadServer struct {
    ID             int64     `json:"id"`
    FileUploadID   int64     `json:"file_upload_id"`
    ServerID       int64     `json:"server_id"`
    ServerName     string    `json:"server_name,omitempty"`
    Status         string    `json:"status"`       // pending|success|failed
    ErrorMessage   string    `json:"error_message,omitempty"`
    FileName       string    `json:"file_name"`    // 文件名
    RemoteFullPath string    `json:"remote_full_path"` // 完整远程路径
    CreatedAt      time.Time `json:"created_at"`
}
```

### 2. 新增数据库表

在 `backend/internal/repository/db.go` 的 `initSchema` 函数中添加：

```sql
CREATE TABLE IF NOT EXISTS file_uploads (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    local_path TEXT NOT NULL,
    remote_path TEXT NOT NULL,
    status TEXT DEFAULT 'pending',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS file_upload_servers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    file_upload_id INTEGER NOT NULL,
    server_id INTEGER NOT NULL,
    server_name TEXT DEFAULT '',
    status TEXT DEFAULT 'pending',
    error_message TEXT DEFAULT '',
    file_name TEXT NOT NULL,
    remote_full_path TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (file_upload_id) REFERENCES file_uploads(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_file_upload_servers_upload ON file_upload_servers(file_upload_id);
CREATE INDEX IF NOT EXISTS idx_file_upload_servers_server ON file_upload_servers(server_id);
```

## 输入

- 需求文档 `06-新增文件管理.md`
- 现有 `backend/internal/model/model.go`
- 现有 `backend/internal/repository/db.go`

## 输出

- `backend/internal/model/model.go` - 添加 FileUpload 和 FileUploadServer 模型
- `backend/internal/repository/db.go` - 添加表结构定义

## 依赖

- 无

## 验收标准

- [ ] FileUpload 结构体包含所有必要字段
- [ ] FileUploadServer 结构体包含所有必要字段
- [ ] 数据库表结构定义完整，包含外键和索引
- [ ] 符合现有项目代码风格
