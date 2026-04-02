# 风险评估与运行效果报告 - node-pilot

## 执行摘要

本报告审核了从提交 `10cb3c81` (验收授权登录) 到 HEAD 的所有代码变更，主要涉及文件上传管理功能的实现。代码构建成功，发现以下风险等级的问题。

## 代码变更概览

- **提交数**: 11
- **新增文件**: 6 (backend: 2, frontend: 4)
- **修改文件**: 8
- **删除文件**: 0

### 提交列表

| 提交 | 描述 |
|------|------|
| `a36fcf1` | feat: 添加--files参数支持可配置文件上传目录 |
| `f5b41e4` | fix: 文件管理SSH密码解密失败+编辑模式数据删除bug |
| `39673c6` | AI自动化开发+08-file-form-page.md |
| `21bcfb9` | AI自动化开发+07-file-list-page.md |
| `3efe00d` | AI自动化开发+06-frontend-store.md |
| `1b04267` | AI自动化开发+05-frontend-types-api.md |
| `dc9fceb` | AI自动化开发+04-backend-handlers-routes.md |
| `99fa242` | AI自动化开发+03-file-upload-service.md |
| `60397fd` | AI自动化开发+02-repository-methods.md |
| `f7fa7a7` | AI自动化开发+01-database-model-and-migration.md |

---

## 🔧 自动修复汇总

### ✅ 已自动修复 (Auto-Fixed)

| 问题 | 位置 | 修复方式 |
|------|------|----------|
| SSH密码为空导致连接失败 | `fileupload.go:155` | 添加AES解密方法，解密server.PasswordEncrypted |
| 编辑模式删除主记录bug | `fileupload.go:196` | 调用DeleteFileUploadServers(id)替代DeleteFileUploads([]int64{id}) |
| 缺少DeleteFileUploadServers方法 | `db.go:780` | 新增方法仅删除file_upload_servers关联表 |

---

## 已识别风险汇总

### 🔴 严重风险 (Critical)

#### 1. 硬编码AES加密密钥
**位置**: `backend/internal/service/fileupload.go:24`
```go
encKey: []byte("12345678901234567890123456789012"),
```
**问题**: 密钥硬编码在代码中，且与handler.go中相同
**影响**: 如果代码泄露，攻击者可解密所有服务器密码
**修复建议**: 从环境变量或配置文件读取密钥

#### 2. 删除操作使用string拼接SQL
**位置**: `backend/internal/repository/db.go:729-736`
```go
query := `DELETE FROM file_uploads WHERE id IN (` + strings.Join(placeholders, ",") + `)`
_, err := r.db.Exec(query, args...)
```
**问题**: 虽然当前ID为int64类型安全，但使用字符串拼接SQL是危险模式
**修复建议**: 使用 `fmt.Sprintf` 或 `strings.Builder` 预编译SQL语句

---

### 🟠 高风险 (High)

#### 3. 文件名未做安全校验
**位置**: `backend/internal/handler/fileupload.go:328-329`
```go
filename := strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + filepath.Base(header.Filename)
localPath := filepath.Join(h.baseDir, filename)
```
**问题**: `filepath.Base` 可阻止目录遍历攻击，但未验证文件扩展名
**影响**: 可能上传恶意文件
**修复建议**: 后端也应校验文件扩展名

#### 4. 上传文件大小无限制
**位置**: `backend/internal/handler/fileupload.go:305-335`
```go
func (h *FileUploadHandler) UploadFileToStorage(c *gin.Context) {
    file, header, err := c.Request.FormFile("file")
    // 无大小校验
    io.Copy(out, file)
}
```
**问题**: 未限制上传文件大小
**影响**: 可能导致磁盘空间耗尽 (DoS)
**修复建议**: 添加文件大小校验 (`header.Size > maxSize`)

#### 5. 缺少授权检查
**位置**: `backend/internal/handler/fileupload.go` (整个文件)
**问题**: 文件上传API未检查用户是否有权限访问指定的服务器
**影响**: 用户可能上传文件到无权访问的服务器
**修复建议**: 在执行前验证用户对目标服务器的权限

---

### 🟡 中风险 (Medium)

#### 6. 前端密码类型断言
**位置**: `frontend/src/views/FileForm.vue:323`
```typescript
const fileRes = await fileUploadApi.uploadFile(form.value.files[0]) as any;
```
**问题**: 使用 `as any` 绕过类型检查
**影响**: 如果API返回格式变化，可能导致运行时错误
**修复建议**: 定义明确的返回类型接口

#### 7. JWT密钥硬编码
**位置**: `backend/cmd/server/main.go:74`
```go
jwtSecret := "node-pilot-jwt-secret-key-32bytes!"
```
**问题**: JWT密钥硬编码在代码中
**影响**: 生产环境中密钥应从环境变量读取
**修复建议**: 从环境变量 `JWT_SECRET` 读取

#### 8. CORS允许所有来源
**位置**: `backend/internal/handler/handler.go:34`
```go
CheckOrigin: func(r *http.Request) bool {
    return true
}
```
**问题**: 允许所有来源跨域请求
**影响**: 在生产环境中存在安全风险
**修复建议**: 配置允许的特定域名列表

---

### 🟢 低风险 (Low)

#### 9. 重复的解密方法
**问题**: `FileUploadService.decrypt` 与 `Handler.decrypt` 实现相同
**位置**: `fileupload.go:47-70` vs `handler.go:554-597`
**修复建议**: 抽取为共享工具函数

#### 10. 未使用的变量
**位置**: `frontend/src/views/FileForm.vue:316`
```typescript
const formData = new FormData();
formData.append('name', form.value.name);
formData.append('remotePath', form.value.remotePath);
formData.append('serverIds', JSON.stringify(form.value.serverIds));
```
**问题**: `formData` 创建后未使用
**影响**: 轻微性能浪费
**修复建议**: 删除未使用的代码

#### 11. 错误处理使用alert
**位置**: `frontend/src/views/FileForm.vue:335`
```typescript
alert(e.message || '保存失败');
```
**问题**: 使用alert用户体验不佳
**修复建议**: 使用更友好的提示组件

---

## 冗余函数分析

| 函数名 | 文件位置 | 调用次数 | 状态 |
|--------|----------|----------|------|
| 无 | - | - | - |

---

## 安全漏洞分析

| 漏洞类型 | 位置 | 严重程度 | 修复建议 |
|----------|------|----------|----------|
| 硬编码密钥 | `fileupload.go:24`, `main.go:74` | Critical | 使用环境变量 |
| 路径遍历潜在风险 | `fileupload.go:328` | High | 添加文件名白名单校验 |
| 文件大小无限制 | `fileupload.go:305` | High | 添加maxSize校验 |
| 授权缺失 | `fileupload.go` 整体 | High | 添加用户-服务器权限验证 |
| CORS全开 | `handler.go:34` | Medium | 配置允许域名 |

---

## 代码质量评分

| 指标 | 评分 | 说明 |
|------|------|------|
| 可维护性 | 75/100 | 代码结构清晰，但有重复代码 |
| 安全性 | 60/100 | 存在硬编码密钥和授权缺失问题 |
| 性能 | 85/100 | 并发处理正确，无明显性能问题 |

---

## 运行效果验证

```
[Backend Build]
BUILD_SUCCESS

[Frontend Build]
✓ built in 2.51s
```

---

## 修复优先级建议

### P0 - 必须修复 (影响安全)
1. **硬编码AES密钥** → 移至环境变量或配置文件
2. **JWT密钥硬编码** → 移至环境变量
3. **文件大小无限制** → 添加 maxSize 校验

### P1 - 强烈建议
4. **授权检查缺失** → 添加用户-服务器权限验证
5. **文件名安全校验** → 后端也应校验文件扩展名

### P2 - 建议优化
6. **CORS配置** → 生产环境配置允许域名
7. **重复解密方法** → 抽取共享函数
8. **前端类型断言** → 定义明确接口

---

## 总结

本次审核的代码变更实现了文件上传管理功能，代码结构清晰，主要问题是**安全配置不当**（硬编码密钥、CORS全开、授权缺失）和**防御性编程不足**（文件大小无限制、文件名未校验）。

建议优先修复P0级别的安全问题后再上线生产环境。
