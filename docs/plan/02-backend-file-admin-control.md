# 任务二：后端文件上传 Admin 权限控制

## 任务描述

在后端对文件上传相关的所有 API 路由添加 `RequireRole("ROLE_ADMIN")` 中间件，确保非 admin 用户调用文件上传/部署接口时返回 403 Forbidden。

## 详细说明

### 1. 修改 file-uploads 路由组

文件：`backend/cmd/server/main.go`

找到第 154-163 行的 file-uploads 路由组，将：
```go
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

改为：
```go
fileUploads := api.Group("/v1/file-uploads")
fileUploads.Use(middleware.JWTAuth(jwtSecret))
fileUploads.Use(middleware.RequireRole("ROLE_ADMIN"))
{
    fileUploads.GET("", fileUploadHandler.ListFileUploads)
    fileUploads.POST("", fileUploadHandler.CreateFileUpload)
    fileUploads.PUT("/:id", fileUploadHandler.UpdateFileUpload)
    fileUploads.DELETE("", fileUploadHandler.DeleteFileUploads)
    fileUploads.POST("/:id/execute", fileUploadHandler.ExecuteFileUpload)
    fileUploads.GET("/:id/results", fileUploadHandler.GetFileUploadResults)
}
```

### 2. 修改独立的文件上传路由

文件：`backend/cmd/server/main.go`

找到第 165 行：
```go
api.POST("/v1/file-uploads/upload-file", fileUploadHandler.UploadFileToStorage)
```

改为（增加 JWT 和 Admin 双重中间件）：
```go
api.POST("/v1/file-uploads/upload-file", 
    middleware.JWTAuth(jwtSecret),
    middleware.RequireRole("ROLE_ADMIN"),
    fileUploadHandler.UploadFileToStorage)
```

### 3. （可选）修改通用 upload/deploy 路由

如果 `/api/upload` 和 `/api/deploy` 也需要保护，找到第 151-152 行：
```go
api.POST("/upload", middleware.JWTAuth(jwtSecret), h.UploadFile)
api.POST("/deploy", middleware.JWTAuth(jwtSecret), h.DeployFile)
```

改为：
```go
api.POST("/upload", middleware.JWTAuth(jwtSecret), middleware.RequireRole("ROLE_ADMIN"), h.UploadFile)
api.POST("/deploy", middleware.JWTAuth(jwtSecret), middleware.RequireRole("ROLE_ADMIN"), h.DeployFile)
```

> 注：如果 `/api/upload` 和 `/api/deploy` 有其他合法用途（如普通用户也需要），跳过此步，仅保护 file-uploads 路由组。

## 输入

- `backend/cmd/server/main.go` 第 151-152 行（upload/deploy）
- `backend/cmd/server/main.go` 第 154-165 行（file-uploads 路由组）

## 输出

- 修改后的 `backend/cmd/server/main.go`

## 依赖

- 01（任务一完成后进行）

## 验收标准

- [ ] file-uploads 路由组中增加了 `.Use(middleware.RequireRole("ROLE_ADMIN"))`
- [ ] `/api/v1/file-uploads/upload-file` 路由增加了 Admin 中间件
- [ ] Go 编译通过（`cd backend && go build ./...`）
