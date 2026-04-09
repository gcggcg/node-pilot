# 后端 Service 执行模式切换实现

## 任务描述

修改 `TaskExecutor.ExecuteScript` 方法，根据 `execution_mode` 选择并发或顺序执行，并实现 `executeSequential` 方法。

## 详细说明

### 1. 修改 ExecuteScript 方法 (`backend/internal/service/task.go`)

在 `ExecuteScript` 方法中：
- 读取 `task.ExecutionMode`
- 如果为空，默认为 `"concurrent"`
- 根据 mode 调用不同的执行方法：
  - `"sequential"` → `executeSequential`
  - `"concurrent"` → `executeConcurrent`（原 `executeBatch` 重命名）

### 2. 新增 executeSequential 方法

实现单线程顺序执行逻辑：
```go
func (e *TaskExecutor) executeSequential(task *model.Task, scripts []*model.Script, servers []*model.Server, password string, serverTaskServerMap map[int64]int64) {
    for i, srv := range servers {
        // 检查任务是否被取消
        if e.IsTaskCancelled(task.ID) {
            break
        }
        
        // 在单个服务器上执行所有脚本
        success := e.executeScriptsOnServer(task, srv, scripts, password, serverTaskServerMap)
        
        if !success {
            // 单线程模式：服务器执行失败，终止后续服务器
            logger.Debug("[TASK-%d] 单线程模式：服务器 %s 执行失败，终止后续服务器", task.ID, srv.Name)
            
            // 标记剩余服务器为 skipped
            for j := i + 1; j < len(servers); j++ {
                remainingSrv := servers[j]
                if tsID, ok := serverTaskServerMap[remainingSrv.ID]; ok {
                    finished := time.Now()
                    e.repo.UpdateTaskServerStatus(tsID, "skipped", "", "前置服务器执行失败，跳过", nil, &finished)
                }
            }
            break
        }
    }
    
    // 更新任务最终状态
    e.finalizeTask(task, servers)
}
```

### 3. 重命名 executeBatch 为 executeConcurrent

将现有的 `executeBatch` 方法重命名为 `executeConcurrent`（或保持原名，只需在调用处修改逻辑）。

### 4. 新增 finalizeTask 辅助方法（如不存在）

将任务状态 finalization 逻辑提取为独立方法，供两种执行模式复用。

### 5. 注意事项

- 单线程模式**不使用** `sync.WaitGroup`，使用普通 for 循环
- 单线程模式需要在 `executeScriptsOnServer` 返回 false 时立即终止
- 即使在单线程模式，也要检查 `IsTaskCancelled()` 允许用户取消
- 单线程模式添加明确日志说明当前执行到第几个服务器
- 失败后剩余服务器标记为 `skipped` 状态而非 `cancelled`

## 输入

- 需求文档：`10-多服务器任务支持单线程和多线程切换.md`
- 现有 Service 文件：`backend/internal/service/task.go`
- Task 01 输出：已添加 execution_mode 字段的 model

## 输出

- 修改后的 `backend/internal/service/task.go`

## 依赖

- 01, 02

## 验收标准

- [ ] ExecuteScript 方法正确读取 task.ExecutionMode
- [ ] execution_mode 为空时默认为 "concurrent"
- [ ] 新增 executeSequential 方法实现顺序执行逻辑
- [ ] executeSequential 在服务器失败时正确终止后续服务器
- [ ] executeSequential 在服务器失败时将剩余服务器标记为 skipped
- [ ] 单线程模式正确检查任务取消状态
- [ ] 日志清晰说明当前执行进度
- [ ] 并发模式行为与修改前一致
