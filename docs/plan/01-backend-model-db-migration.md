# 后端模型与数据库迁移

## 任务描述

在 Task 模型中添加 `execution_mode` 字段，并添加数据库迁移脚本。

## 详细说明

### 1. 修改 Task 模型 (`backend/internal/model/model.go`)

在 `Task` 结构体中添加 `ExecutionMode` 字段：
- 字段名：`ExecutionMode`
- 类型：`string`
- JSON tag：`json:"execution_mode"`
- 可选值：`"concurrent"` (并发，默认) | `"sequential"` (单线程顺序)

### 2. 添加数据库迁移 (`backend/internal/repository/db.go`)

在 `migrateSchema` 函数中添加迁移语句：
```sql
ALTER TABLE tasks ADD COLUMN execution_mode TEXT NOT NULL DEFAULT 'concurrent';
```

**注意**：
- 使用 `ALTER TABLE` 添加新列
- 默认值为 `'concurrent'` 保证向后兼容
- 历史数据默认使用并发模式

## 输入

- 需求文档：`10-多服务器任务支持单线程和多线程切换.md`
- 现有模型文件：`backend/internal/model/model.go`
- 现有数据库迁移文件：`backend/internal/repository/db.go`

## 输出

- 修改后的 `backend/internal/model/model.go`
- 修改后的 `backend/internal/repository/db.go`

## 依赖

- 无

## 验收标准

- [ ] Task 结构体包含 `ExecutionMode string` 字段，json tag 为 `execution_mode`
- [ ] `migrateSchema` 函数包含 `execution_mode` 列的 ALTER TABLE 迁移
- [ ] 默认值为 `concurrent`
- [ ] 数据库查询能够正确读取和写入 `execution_mode` 字段
