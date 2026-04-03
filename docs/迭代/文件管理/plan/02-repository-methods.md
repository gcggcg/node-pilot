# 文件管理 - Repository 层实现

## 任务描述

实现文件上传管理的数据库访问层方法。

## 详细说明

在 `backend/internal/repository/db.go` 中添加以下方法：

### 1. CreateFileUpload

```go
func (r *Repository) CreateFileUpload(fu *model.FileUpload) (int64, error)
```

- 插入 file_uploads 记录
- 返回新记录 ID

### 2. CreateFileUploadServers

```go
func (r *Repository) CreateFileUploadServers(records []*model.FileUploadServer) error
```

- 批量插入 file_upload_servers 记录

### 3. GetFileUploadByID

```go
func (r *Repository) GetFileUploadByID(id int64) (*model.FileUpload, error)
```

- 根据 ID 查询单条记录

### 4. ListFileUploads

```go
func (r *Repository) ListFileUploads(page, pageSize int, status, keyword string, startTime, endTime *time.Time) ([]*model.FileUpload, int64, error)
```

- 分页查询，支持状态筛选、关键字搜索、日期范围
- 返回记录列表和总数

### 5. UpdateFileUpload

```go
func (r *Repository) UpdateFileUpload(fu *model.FileUpload) error
```

- 更新 file_uploads 记录

### 6. DeleteFileUploads

```go
func (r *Repository) DeleteFileUploads(ids []int64) error
```

- 批量删除 file_uploads 记录（级联删除关联记录）

### 7. GetFileUploadServers

```go
func (r *Repository) GetFileUploadServers(fileUploadID int64) ([]*model.FileUploadServer, error)
```

- 根据 file_upload_id 查询关联的服务器记录

### 8. UpdateFileUploadServerStatus

```go
func (r *Repository) UpdateFileUploadServerStatus(id int64, status, errorMsg string) error
```

- 更新单条执行结果状态

### 9. GetLocalFiles

```go
func (r *Repository) GetLocalFiles(dir string) ([]string, error)
```

- 获取 ./data/files/ 目录下的文件列表

## 输入

- 需求文档 `06-新增文件管理.md`
- 现有 `backend/internal/repository/db.go` 模式

## 输出

- `backend/internal/repository/db.go` - 添加所有 FileUpload 相关方法

## 依赖

- 01-database-model-and-migration.md

## 验收标准

- [ ] 所有 Repository 方法实现完整
- [ ] 分页逻辑正确
- [ ] 支持筛选和搜索
- [ ] 符合现有项目代码风格
