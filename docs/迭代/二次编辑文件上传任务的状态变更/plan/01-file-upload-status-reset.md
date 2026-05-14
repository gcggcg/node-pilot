# 修改 UpdateFileUpload Handler 实现状态重置

## 任务描述

修改文件上传任务的 `UpdateFileUpload` 方法，实现编辑保存时状态重置为 "pending" 的功能。

## 详细说明

### 1. 添加 Running 状态校验

在 `UpdateFileUpload` 方法中，获取到 fileUpload 记录后，首先检查其状态：

```go
// 如果任务正在执行中，不允许编辑
if fu.Status == "running" {
    c.JSON(http.StatusBadRequest, gin.H{"error": "任务正在执行中，无法编辑"})
    return
}
```

### 2. 强制重置状态为 Pending

在更新记录前，强制将 status 设置为 "pending"：

```go
// 强制将状态重置为待执行
fu.Status = "pending"
```

然后调用 `h.repo.UpdateFileUpload(fu)` 更新记录。

### 3. 关联记录处理

当前代码在编辑时已经：
- 删除旧的 `file_upload_servers` 关联（line 193）
- 创建新的关联并设置 status = "pending"（line 196-208）

这块逻辑已经正确，无需修改。

## 输入

- 需求文档 `./需求.md`
- 现有代码 `backend/internal/handler/fileupload.go`

## 输出

- 修改后的 `backend/internal/handler/fileupload.go`

## 依赖

- 无

## 验收标准

- [ ] 编辑 success 状态任务后，保存后状态变为 pending
- [ ] 编辑 failed 状态任务后，保存后状态变为 pending
- [ ] 编辑 running 状态任务后，返回错误 "任务正在执行中，无法编辑"
- [ ] 编辑 pending 状态任务后，保存后状态变为 pending
- [ ] 关联的 file_upload_servers 记录状态同步重置为 pending
