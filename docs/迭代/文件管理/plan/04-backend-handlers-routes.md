# 文件管理 - 后端 Handlers 和路由

## 任务描述

实现文件上传管理的 REST API 处理器和路由注册。

## 详细说明

### 1. 新增 Handler 文件

创建 `backend/internal/handler/fileupload.go`：

### 2. FileUploadHandler 结构体

```go
type FileUploadHandler struct {
    repo     *repository.Repository
    uploadSvc *service.FileUploadService
}
```

### 3. API 接口实现

#### ListFileUploads - GET /api/v1/file-uploads

- 参数：page, pageSize, status, fileName, startTime, endTime
- 调用 `repo.ListFileUploads()`
- 返回分页数据

#### CreateFileUpload - POST /api/v1/file-uploads

- 参数：name, localPath, remotePath, serverIds
- 文件存储：接收文件列表，保存到 ./data/files/
- 创建 file_upload 记录
- 创建 file_upload_server 关联记录

#### UpdateFileUpload - PUT /api/v1/file-uploads/:id

- 参数：serverIds, remotePath
- 更新 file_upload 记录
- 不可修改 localPath

#### ExecuteFileUpload - POST /api/v1/file-uploads/:id/execute

- 调用 `uploadSvc.ExecuteUpload()`
- 返回执行状态

#### GetFileUploadResults - GET /api/v1/file-uploads/:id/results

- 查询 file_upload_server 列表
- 返回每台服务器的执行结果

#### DeleteFileUploads - DELETE /api/v1/file-uploads

- 参数：ids 数组
- 删除记录和关联文件
- 调用 `os.Remove()` 删除 ./data/files/ 下的文件

### 4. 路由注册

在 `backend/cmd/server/main.go` 中添加：

```go
fileUploadHandler := handler.NewFileUploadHandler(repo, uploadSvc)

fileUploads := api.Group("/v1/file-uploads")
fileUploads.Use(middleware.JWTAuth(jwtSecret))
{
    fileUploads.GET("", fileUploadHandler.ListFileUploads)
    fileUploads.POST("", fileUploadHandler.CreateFileUpload)
    fileUploads.PUT("/:id", fileUploadHandler.UpdateFileUpload)
    fileUploads.DELETE("", fileUploadHandler.DeleteFileUploads)
    fileUploads.POST("/:id/execute", fileUploadHandler.ExecuteFileUpload)
    fileUploads.GET("/:id/results", fileUploadHandler.GetFileUploadResults)
}
```

### 5. 文件上传处理（multipart/form-data）

CreateFileUpload 接口需要处理文件上传：

- 接收 multipart/form-data
- 使用 `r.FormFile()` 获取文件
- 保存到 `./data/files/` 目录
- 返回保存后的相对路径

```go
func (h *FileUploadHandler) CreateFileUpload(c *gin.Context) {
    // 解析 multipart form
    form, err := c.MultipartForm()
    if err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 获取文件
    files := form.File["files"]
    // ... 保存文件
}
```

## 输入

- 需求文档 `06-新增文件管理.md`
- 现有 `backend/internal/handler/auth.go` 模式
- 现有 `backend/cmd/server/main.go` 路由注册

## 输出

- `backend/internal/handler/fileupload.go` - 文件上传处理器
- `backend/cmd/server/main.go` - 添加路由注册

## 依赖

- 03-file-upload-service.md

## 验收标准

- [ ] 所有 API 接口正确实现
- [ ] 文件上传到 ./data/files/ 目录
- [ ] 路由正确注册
- [ ] 符合现有项目代码风格
