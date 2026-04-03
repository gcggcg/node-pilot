# Bug修复：任务服务器关联缺失

## 任务描述

修复任务管理模块中的两个关联Bug：
1. 任务创建时不创建 task_servers 关联，导致编辑页无法回显服务器列表
2. 任务执行时无法获取关联服务器，导致500错误

## 详细说明

### 问题分析

**Bug 1: 任务编辑页服务器列表数据丢失**
- 位置: `backend/internal/handler/handler.go` CreateTask 函数
- 原因: 第397-409行代码注释说明"不再在此创建 task_servers 关联"
- 影响: `CreateTask` 创建任务时不创建 `task_servers` 记录
- 后果: `GetTaskServers` 返回空数组，编辑页无法回显

**Bug 2: 任务执行接口返回500错误**
- 位置: `backend/internal/service/task.go` ExecuteTask 函数
- 原因: 第65-67行检查 `taskServers` 是否为空
- 影响: 由于创建时没有关联，执行时返回 "任务没有关联服务器"
- 后果: 点击执行按钮返回500错误

### 修复方案

1. 修改 `CreateTask` (handler.go): 恢复创建 `task_servers` 关联的逻辑
2. 修改 `ExecuteTask` (task.go): 兼容处理没有关联服务器的情况

## 输入

- 需求文档: `./任务处理Bug修复.md`
- 依赖库文档: `./docs/modules/node-pilot.md`

## 输出

- 修复后端 `CreateTask` 函数: `backend/internal/handler/handler.go`
- 修复后端 `ExecuteTask` 函数: `backend/internal/service/task.go`
- 运行测试验证修复

## 依赖

无（独立Bug修复任务）

## 验收标准

- [ ] 修复 `CreateTask` 函数，恢复创建 task_servers 关联
- [ ] 修复 `ExecuteTask` 函数，正确处理无关联服务器情况
- [ ] 重新构建后端: `cd backend && go build ./cmd/server`
- [ ] 测试创建任务后编辑页面能正确回显服务器列表
- [ ] 测试任务执行接口能正常工作（不返回500）
