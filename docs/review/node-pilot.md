# 风险评估与运行效果报告 - NodePilot

## 执行摘要

本次代码审查覆盖了从提交 `c89e37761a34c0e7aa020d157f9ecb85a9e22254` 之后的变更，主要涉及**批量脚本流式执行功能**的实现。审查发现**1个高危安全风险**（硬编码AES密钥）和**1个冗余代码问题**（已自动修复）。

| 类型 | 数量 |
|------|------|
| 🔴 严重风险 | 0 |
| 🟠 高风险 | 1 |
| 🟡 中风险 | 2 |
| 🟢 低风险 | 0 |
| ✅ 已自动修复 | 1 |

---

## 代码变更概览

- **提交数**: 4
- **新增文件**: 14 (主要是frontend/dist构建产物)
- **修改文件**: 9
- **删除文件**: 0

### 提交历史
```
ac08221 feat: 批量脚本流式执行功能完成
75f8dd2 AI自动化开发+03-backend-handlers
0e182b5 AI自动化开发+02-repository-methods
24df3b8 AI自动化开发+01-database-model
```

---

## 🔧 自动修复汇总

### ✅ 已自动修复 (Auto-Fixed)

| 问题 | 位置 | 修复方式 |
|------|------|----------|
| 冗余函数 `executeOnServer` | backend/internal/service/task.go:678 | 已删除 - 该函数被新实现的 `executeScriptsOnServer` 替代，属于死代码 |

---

## 已识别风险汇总

### 🟠 高风险 (High)

#### 1. 硬编码AES加密密钥

**严重程度**: High

**位置**:
- `backend/internal/handler/handler.go:43`
- `backend/internal/service/task.go:139`
- `backend/internal/service/fileupload.go:31`

**问题描述**:
项目中使用了硬编码的AES-256密钥 `12345678901234567890123456789012`，该密钥在三个不同的文件中重复出现。

```go
// handler.go:43
key := []byte("12345678901234567890123456789012") // exactly 32 bytes for AES-256

// task.go:139
key := []byte("12345678901234567890123456789012") // 32 bytes for AES-256

// fileupload.go:31
encKey:  []byte("12345678901234567890123456789012"),
```

**安全风险**:
1. **密钥泄露风险**: 密钥被硬编码在源代码中，如果代码仓库被泄露（如GitHub public repo），所有使用该密钥加密的数据都可以被解密
2. **密钥复用风险**: 同一密钥用于多个加密场景（密码加密、文件上传加密），增加了密钥暴露的攻击面
3. **无法轮换密钥**: 硬编码密钥无法在不更新代码的情况下进行轮换

**修复建议**:
1. 将密钥移至配置文件或环境变量
2. 使用 `os.Getenv("NODE_PILOT_KEY")` 或配置文件读取
3. 启动时生成随机密钥并持久化（首次运行）
4. 确保不同环境使用不同密钥

**示例修复**:
```go
func getEncryptionKey() []byte {
    key := os.Getenv("NODE_PILOT_KEY")
    if key == "" {
        // 生成随机密钥（仅首次启动）
        key = generateRandomKey()
        // 可选：保存到配置文件
    }
    return []byte(key) // 确保是32字节
}
```

---

### 🟡 中风险 (Medium)

#### 2. SQL注入风险 - GetScripts方法

**严重程度**: Medium

**位置**: `backend/internal/repository/db.go:332-358`

**问题描述**:
`GetScripts` 方法使用字符串拼接构建SQL查询：

```go
func (r *Repository) GetScripts(ids []int64) ([]*model.Script, error) {
    // ...
    placeholders := make([]string, len(ids))
    args := make([]interface{}, len(ids))
    for i, id := range ids {
        placeholders[i] = "?"
        args[i] = id
    }
    query := `SELECT id, name, description, content, target_path, created_at, updated_at FROM scripts WHERE id IN (` + strings.Join(placeholders, ",") + `)`
    rows, err := r.db.Query(query, args...)
    // ...
}
```

**风险分析**:
当前实现使用 `?` 占位符和参数化查询，**风险较低**。但代码构造方式存在潜在风险，如果未来有人修改为直接拼接 `id` 而非使用占位符，可能引入SQL注入。

**当前状态**: ✅ 使用参数化查询，风险可控

**建议**:
1. 添加输入验证，确保 `ids` 数组中的每个元素都是有效的正整数
2. 添加最大ID数量限制，防止大量ID导致的DoS风险

---

#### 3. 路径遍历风险 - 脚本执行

**严重程度**: Medium

**位置**: `backend/internal/service/task.go:395-676`

**问题描述**:
`executeSingleScript` 函数使用 `filepath.Dir()` 获取目标目录后执行shell命令：

```go
targetDir := filepath.Dir(script.TargetPath)
targetFile := script.TargetPath
// ...
execCmd := fmt.Sprintf("cd %s && /bin/bash %s", targetDir, targetFile)
session.CombinedOutput(execCmd)
```

**风险分析**:
如果攻击者能够控制 `script.TargetPath`（如通过SQL注入或修改数据库），可能通过构造路径如 `/tmp/../../../etc/cron.d/malicious` 进行路径遍历攻击。

**当前缓解措施**:
1. `filepath.Dir()` 会解析 `..` 但不会阻止 `cd` 到任意目录
2. SFTP上传限制了目标文件路径

**修复建议**:
1. 验证 `targetDir` 不包含 `..`
2. 使用 `filepath.Clean()` 并检查结果
3. 添加允许目录的白名单机制

```go
// 添加路径安全验证
targetDir := filepath.Dir(script.TargetPath)
cleanDir := filepath.Clean(targetDir)
if strings.Contains(cleanDir, "..") {
    return false, fmt.Errorf("path traversal detected: %s", targetDir)
}
// 验证目录在允许的基础目录下
allowedBase := "/opt/scripts"
if !strings.HasPrefix(cleanDir, allowedBase) {
    return false, fmt.Errorf("invalid script directory: %s", targetDir)
}
```

---

## 代码质量分析

### 冗余函数分析

| 函数名 | 文件位置 | 调用次数 | 状态 |
|--------|----------|----------|------|
| `executeOnServer` | backend/internal/service/task.go:678 | 0 | ✅ 已删除 |

**说明**: `executeOnServer` 是旧的单脚本执行函数，被新的 `executeScriptsOnServer`（批量脚本执行）替代后成为死代码。

---

## 安全漏洞分析

| 漏洞类型 | 位置 | 严重程度 | 修复建议 |
|----------|------|----------|----------|
| 硬编码密钥 | handler.go:43, task.go:139, fileupload.go:31 | High | 移至环境变量或配置文件 |
| 路径遍历 | task.go:executeSingleScript | Medium | 添加路径验证和白名单 |

---

## 运行效果验证

### 构建测试

```bash
cd backend && go build -o /tmp/node-pilot-review ./cmd/server
# ✅ 构建成功，无错误
```

### 代码变更验证

```
修改文件:
- backend/internal/handler/handler.go    (ScriptIDs字段支持)
- backend/internal/model/model.go         (WSMessage批量字段)
- backend/internal/repository/db.go      (GetScripts方法, script_ids列)
- backend/internal/service/task.go       (批量执行引擎, executeOnServer已删除)
```

---

## 修复优先级建议

### P0 - 必须立即修复 (高危)
1. **硬编码AES密钥**: 立即将密钥移至环境变量或配置中心

### P1 - 尽快修复 (中危)
2. **路径遍历防护**: 在生产环境部署前添加路径验证
3. **SQL查询限制**: 添加ID数量上限和输入验证

---

## 总结

本次审查的批量脚本流式执行功能在业务逻辑实现上较为完善，具备：
- ✅ 顺序执行与错误终止
- ✅ WebSocket实时进度推送
- ✅ 批量脚本执行状态追踪

但存在**高危安全风险**（硬编码密钥）需要立即修复。建议：
1. 优先处理密钥配置问题
2. 后续迭代中添加路径遍历防护
3. 建议使用环境变量或安全的配置中心管理敏感信息

---

*报告生成时间: 2026-04-07*
*审查版本: c89e37761a34c0e7aa020d157f9ecb85a9e22254..HEAD*
