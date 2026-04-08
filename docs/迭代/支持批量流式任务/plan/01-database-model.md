# Database Model - Add Batch Script Field

## 任务描述

修改 Task 数据模型，添加批量脚本支持字段，使任务能够关联多个脚本。

## 详细说明

### 1. 修改 Task 模型
在 `backend/internal/model/model.go` 中：
- 将 `ScriptID int64` 字段修改为 `ScriptIDs string` (TEXT 类型，存储逗号分隔的脚本ID列表)
- 添加注释说明该字段的存储格式

### 2. 数据库迁移
由于项目使用 SQLite，需要在数据库层面添加新字段：
- Task 表添加 `script_ids TEXT` 字段（存储逗号分隔的脚本ID列表，如 "1,2,3"）
- 保持向后兼容：原有的 `script_id` 字段保留，但标记为废弃

### 3. 代码修改示例

**model.go 修改:**
```go
type Task struct {
    ID         int64      `json:"id"`
    ScriptID   int64      `json:"script_id"`   // 废弃字段，保留兼容性
    ScriptIDs  string     `json:"script_ids"`  // 新字段：逗号分隔的脚本ID列表
    Name       string     `json:"name"`
    Status     string     `json:"status"`
    CreatedAt  time.Time  `json:"created_at"`
    StartedAt  *time.Time `json:"started_at,omitempty"`
    FinishedAt *time.Time `json:"finished_at,omitempty"`
}
```

### 4. 兼容性处理
- 读取任务时，优先读取 `ScriptIDs` 字段
- 如果 `ScriptIDs` 为空，则回退到 `ScriptID` 字段（单脚本兼容）
- 创建/更新任务时，优先使用 `ScriptIDs` 字段

## 输入

- 现有 Task 模型定义 (backend/internal/model/model.go)
- Repository 数据库访问层 (backend/internal/repository/db.go)

## 输出

- 修改后的 model.go
- db.go 中添加必要的数据库字段更新逻辑

## 依赖

- 无

## 验收标准

- [ ] Task 结构添加了 ScriptIDs 字段 (TEXT 类型)
- [ ] 原有 ScriptID 字段保留（向后兼容）
- [ ] 数据库迁移逻辑正确（添加 script_ids 字段）
- [ ] 代码能正确处理空 ScriptIDs 的情况（回退到 ScriptID）
