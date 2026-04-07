# Stream Execution Engine - Batch Script Sequential Execution

## 任务描述

实现批量脚本的流式顺序执行引擎，支持错误终止和详细的错误信息返回。

## 详细说明

### 1. 修改 ExecuteScript 方法
在 `backend/internal/service/task.go` 中，修改 `ExecuteScript` 方法支持批量脚本：

**核心逻辑：**
```go
func (e *TaskExecutor) ExecuteScript(task *model.Task, script *model.Script, servers []*model.Server, password string, serverTaskServerMap map[int64]int64) {
    // ... 现有逻辑 ...
    
    // 解析脚本ID列表
    scriptIDs := ParseScriptIDs(task.ScriptIDs)
    
    // 单脚本兼容：如果没有 ScriptIDs，回退到 ScriptID
    if len(scriptIDs) == 0 && script != nil {
        scriptIDs = []int64{script.ID}
    }
    
    // 遍历服务器执行批量脚本
    for _, srv := range servers {
        e.executeScriptsOnServer(task, srv, scriptIDs, password, serverTaskServerMap)
    }
}
```

### 2. 新增 executeScriptsOnServer 方法
创建新方法处理单服务器上的批量脚本执行：

**核心逻辑：**
```go
func (e *TaskExecutor) executeScriptsOnServer(task *model.Task, srv *model.Server, scriptIDs []int64, password string, serverTaskServerMap map[int64]int64) {
    // 获取所有脚本
    scripts, err := e.repo.GetScripts(scriptIDs)
    if err != nil {
        // 错误处理
        return
    }
    
    // 顺序执行每个脚本
    for i, scr := range scripts {
        // 执行单个脚本
        success := e.executeSingleScript(task, srv, scr, password, serverTaskServerMap, i, len(scripts))
        
        // 如果失败，终止后续脚本
        if !success {
            // 记录错误并返回
            e.recordBatchError(task, srv, scr, i, len(scripts), "脚本执行失败")
            return
        }
    }
}
```

### 3. 修改 executeSingleScript 方法（原有 executeOnServer 的抽取）

将原来在 `executeOnServer` 中的单个脚本执行逻辑抽取为独立方法：

```go
func (e *TaskExecutor) executeSingleScript(task *model.Task, srv *model.Server, script *model.Script, password string, serverTaskServerMap map[int64]int64, scriptIndex, totalScripts int) bool {
    // 构建执行命令（使用目录切换）
    targetDir := filepath.Dir(script.TargetPath)
    targetFile := script.TargetPath
    
    // 目录切换执行命令
    execCmd := fmt.Sprintf("cd %s && /bin/bash %s", targetDir, targetFile)
    
    // 执行并捕获输出
    // ... 执行逻辑 ...
    
    // 返回是否成功
    return err == nil
}
```

### 4. 修改命令执行格式（目录切换）

**原格式（需废弃）：**
```go
execCmd := fmt.Sprintf("/bin/bash %s", targetFile)
```

**新格式：**
```go
execCmd := fmt.Sprintf("cd %s && /bin/bash %s", targetDir, targetFile)
```

### 5. 增强错误信息

**错误结构：**
```go
type BatchScriptError struct {
    CurrentIndex int    // 当前执行到第几个（从1开始）
    TotalCount   int    // 总脚本数量
    ScriptPath   string // 失败脚本的完整路径
    ScriptName   string // 失败脚本的名称
    ErrorMsg     string // 具体错误信息
}
```

**错误消息格式：**
```
脚本执行失败 (2/5)
脚本路径: /path/to/script2.sh
错误: exit status 1
命令输出: ...
```

### 6. WebSocket 消息增强

在执行过程中，通过 WebSocket 发送每个脚本的执行状态：

```go
// 脚本开始执行
WSMessage{
    Type: "script_start",
    TaskID: task.ID,
    ServerID: srv.ID,
    ScriptIndex: i + 1,    // 从1开始
    ScriptPath: scr.TargetPath,
    ScriptName: scr.Name,
    TotalScripts: len(scripts),
}

// 脚本执行完成
WSMessage{
    Type: "script_done",
    TaskID: task.ID,
    ServerID: srv.ID,
    ScriptIndex: i + 1,
    ScriptPath: scr.TargetPath,
    Status: "success" | "failed",
    Output: output,
}
```

## 输入

- backend/internal/service/task.go
- 批量脚本 ID 列表
- Repository.GetScripts 方法

## 输出

- 修改后的 task.go，包含流式执行引擎

## 依赖

- 01（Task.ScriptIDs 字段）
- 02（GetScripts, ParseScriptIDs 方法）
- 03（Handler 已支持批量脚本）

## 验收标准

- [ ] 单脚本模式仍然正常工作（向后兼容）
- [ ] 批量脚本按顺序执行
- [ ] 错误发生时立即终止后续脚本
- [ ] 错误信息包含脚本索引（从1开始）和总数量
- [ ] 错误信息包含失败脚本的完整路径
- [ ] WebSocket 消息正确传递脚本执行状态
- [ ] 使用目录切换执行命令格式 `cd <dir> && /bin/bash <script>`
- [ ] 所有日志包含时间戳和脚本路径
