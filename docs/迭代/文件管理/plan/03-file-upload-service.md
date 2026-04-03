# 文件管理 - 文件上传服务

## 任务描述

实现基于 SFTP 的文件上传执行服务。

## 详细说明

### 1. 新增 Service 文件

创建 `backend/internal/service/fileupload.go`：

### 2. FileUploadService 结构体

```go
type FileUploadService struct {
    repo    *repository.Repository
    sshPool *SSHPool
    baseDir string  // ./data/files/
}
```

### 3. UploadFile 方法

```go
func (s *FileUploadService) UploadFile(serverID int64, localFilePath, remotePath string) (string, error)
```

- 使用 SSH Pool 获取服务器连接
- 通过 SFTP 上传文件
- 构造完整的远程路径：`remotePath + "/" + filename`
- 返回执行结果或错误信息

### 4. ExecuteUpload 方法

```go
func (s *FileUploadService) ExecuteUpload(fileUploadID int64) error
```

- 查询 file_upload 配置信息
- 获取关联的服务器列表
- 遍历服务器列表，并行执行上传（参考 task.go 的并行执行模式）
- 更新每条 file_upload_server 的状态和结果
- 整体更新 file_upload 状态

### 5. SFTP 文件传输实现

参考 golang.org/x/crypto/ssh 包中的 SFTPClient 使用方式：

```go
// 建立 SFTP Session
sftpClient, err := sftp.NewClient(sshClient)
if err != nil {
    return "", err
}
defer sftpClient.Close()

// 打开本地文件
srcFile, err := os.Open(localPath)
if err != nil {
    return "", err
}
defer srcFile.Close()

// 创建远程文件
dstFile, err := sftpClient.Create(remoteFullPath)
if err != nil {
    return "", err
}
defer dstFile.Close()

// 复制内容
_, err = io.Copy(dstFile, srcFile)
```

### 6. 确保目录存在

上传前需确保远程目录存在，必要时创建：

```go
// 创建远程目录（如果不存在）
err = sftpClient.MkdirAll(filepath.Dir(remoteFullPath))
```

## 输入

- 需求文档 `06-新增文件管理.md`
- 现有 `backend/internal/service/ssh.go` (SSH Pool)
- 现有 `backend/internal/service/task.go` (并行执行模式)

## 输出

- `backend/internal/service/fileupload.go` - 文件上传服务

## 依赖

- 02-repository-methods.md

## 验收标准

- [ ] SFTP 文件上传功能正常
- [ ] 远程目录自动创建
- [ ] 并行上传多台服务器
- [ ] 错误信息准确记录
- [ ] 符合现有项目代码风格
