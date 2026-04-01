# NodePilot 后端技术规格书

## 1. 项目概述

- **项目名称**: NodePilot 后端服务
- **技术栈**: Go + Gin + SQLite
- **功能定位**: 批量服务器管理平台的REST API和WebSocket服务

## 2. 项目结构

```
backend/
├── cmd/server/
│   └── main.go              # 程序入口
├── internal/
│   ├── config/
│   │   └── config.go        # 配置结构体
│   ├── handler/
│   │   └── handler.go       # HTTP处理器(16个API端点)
│   ├── logger/
│   │   └── logger.go        # 日志系统(DEBUG/INFO/WARN/ERROR)
│   ├── model/
│   │   └── model.go         # 数据模型(5个结构体)
│   ├── repository/
│   │   └── db.go           # 数据库访问层(CRUD+批量删除+迁移)
│   ├── service/
│   │   ├── ssh.go          # SSH连接池管理
│   │   └── task.go         # 任务执行器(异步+取消支持)
│   └── websocket/
│       └── hub.go          # WebSocket Hub(TaskID过滤广播)
├── data/                    # SQLite数据目录
│   └── servers.db           # SQLite数据库文件
└── web/                     # 前端构建产物(嵌入式)
```

## 3. 启动参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--db` | `./data/servers.db` | SQLite数据库路径 |
| `--listen` | `:8080` | 监听地址 |
| `--debug` | `false` | 启用Debug日志 |
| `--log` | `""` | 日志文件路径 |

**启动示例**:
```bash
./node-pilot --db ../data/servers.db --listen :8080 --debug --log ../data/app.log
```

## 4. 数据库表结构

### servers (服务器表)
| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER PRIMARY KEY | 自增ID |
| name | TEXT NOT NULL | 服务器名称 |
| host | TEXT NOT NULL | IP地址 |
| port | INTEGER DEFAULT 22 | SSH端口 |
| username | TEXT NOT NULL | 用户名 |
| password_encrypted | TEXT | AES-256-GCM加密密码 |
| connection_status | TEXT DEFAULT 'unknown' | 连接状态 (online/offline/unknown) |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

### scripts (脚本表)
| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER PRIMARY KEY | 自增ID |
| name | TEXT NOT NULL | 脚本名称 |
| description | TEXT | 脚本描述 |
| content | TEXT NOT NULL | 脚本内容 |
| target_path | TEXT NOT NULL | 远程目标路径 |
| created_at | DATETIME | 创建时间 |
| updated_at | DATETIME | 更新时间 |

### tasks (任务表)
| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER PRIMARY KEY | 自增ID |
| script_id | INTEGER | 关联脚本ID |
| name | TEXT NOT NULL | 任务名称 |
| status | TEXT DEFAULT 'pending' | 状态 |
| created_at | DATETIME | 创建时间 |
| started_at | DATETIME | 开始时间 |
| finished_at | DATETIME | 结束时间 |

**状态枚举**: `pending` | `running` | `completed` | `cancelled` | `failed`

### task_servers (任务服务器关联表)
| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER PRIMARY KEY | 自增ID |
| task_id | INTEGER | 关联任务ID |
| server_id | INTEGER | 关联服务器ID |
| status | TEXT | 状态 |
| output | TEXT | 输出(限制最新50行) |
| error | TEXT | 错误信息 |
| started_at | DATETIME | 开始时间 |
| finished_at | DATETIME | 结束时间 |

**状态枚举**: `pending` | `running` | `success` | `failed`

## 5. API接口

### 5.1 服务器管理 `/api/servers`

| 方法 | 路径 | 说明 | 请求体 |
|------|------|------|--------|
| GET | `/servers` | 获取服务器列表 | - |
| GET | `/servers/:id` | 获取单个服务器 | - |
| POST | `/servers` | 创建服务器 | `{name, host, port, username, password}` |
| PUT | `/servers/:id` | 更新服务器 | `{name, host, port, username, password}` |
| DELETE | `/servers/:id` | 删除单个服务器 | - |
| POST | `/servers/batch-delete` | 批量删除服务器 | `{ids: [1,2,3]}` |
| POST | `/servers/:id/test` | 测试SSH连接并更新状态 | - |

**连接状态更新逻辑**:
- 测试成功: `connection_status` → `online`
- 测试失败: `connection_status` → `offline`
- 初始创建: `connection_status` → `unknown`

### 5.2 脚本管理 `/api/scripts`

| 方法 | 路径 | 说明 | 请求体 |
|------|------|------|--------|
| GET | `/scripts` | 获取脚本列表 | - |
| GET | `/scripts/:id` | 获取单个脚本 | - |
| POST | `/scripts` | 创建脚本 | `{name, description, content, target_path}` |
| PUT | `/scripts/:id` | 更新脚本 | `{name, description, content, target_path}` |
| DELETE | `/scripts/:id` | 删除单个脚本 | - |
| POST | `/scripts/batch-delete` | 批量删除脚本 | `{ids: [1,2,3]}` |

### 5.3 任务管理 `/api/tasks`

| 方法 | 路径 | 说明 | 请求体 |
|------|------|------|--------|
| GET | `/tasks` | 获取任务列表 | - |
| GET | `/tasks/:id` | 获取任务详情(含关联服务器) | - |
| POST | `/tasks` | 创建并执行任务 | `{script_id, name, server_ids}` |
| DELETE | `/tasks/:id` | 取消任务 | - |
| POST | `/tasks/batch-delete` | 批量删除任务(级联删除task_servers) | `{ids: [1,2,3]}` |
| GET | `/tasks/:id/output` | SSE实时输出流 | - |

### 5.4 文件操作

| 方法 | 路径 | 说明 | 请求体 |
|------|------|------|--------|
| POST | `/api/upload` | 上传文件到本地 | `multipart/form-data` |
| POST | `/api/deploy` | 部署文件到远程服务器 | `{server_ids, local_path, remote_path}` |

### 5.5 WebSocket

| 路径 | 说明 | 参数 |
|------|------|------|
| `/ws?task_id=N` | 任务输出实时推送 | task_id: 任务ID |

## 6. WebSocket消息格式

```go
type WSMessage struct {
    Type       string    // output|server_start|server_done|task_start|task_done
    TaskID     int64     // 任务ID
    ServerID   int64     // 服务器ID(可选)
    ServerName string    // 服务器名称(可选)
    Content    string    // 输出内容
    Status     string    // 状态
    ExitCode   int       // 退出码
    Timestamp  time.Time // 时间戳
    Total      int       // 总数(任务完成时)
    Success    int       // 成功数(任务完成时)
    Failed     int       // 失败数(任务完成时)
}
```

**消息类型**:
- `task_start`: 任务开始执行
- `server_start`: 开始在服务器上执行
- `output`: 实时输出
- `server_done`: 单服务器执行完成
- `task_done`: 任务全部完成

## 7. 核心功能

### 7.1 SSH连接池
- 连接池管理，避免频繁建立连接
- 支持密码认证
- 支持SFTP文件传输

### 7.2 任务执行流程
1. 创建任务记录
2. SSH连接目标服务器
3. SFTP上传脚本到目标路径
4. 设置执行权限 `chmod +x`
5. 执行脚本并捕获输出
6. 更新任务状态

### 7.3 任务取消机制
- 使用 `sync.Map` 存储已取消的任务ID: `cancelledTasks`
- `CancelTask` 将任务ID存入 `cancelledTasks`
- `ExecuteScript` 每个批次执行前检查:
  ```go
  if _, cancelled := cancelledTasks.Load(task.ID); cancelled {
      // 停止执行新批次
      break
  }
  ```
- 取消后保留 `cancelled` 状态，不会被覆盖为 `completed`
- 任务完成后自动从 `cancelledTasks` 中移除

### 7.4 WebSocket广播
- Client结构存储TaskID
- `BroadcastToTask` 只向观看指定任务的客户端广播
- 防止旧任务消息污染新任务

### 7.5 输出限制
- 数据库只存储最新50行输出
- 使用 `limitLines()` 函数截断

### 7.6 批量删除级联
- `DeleteTasks` 使用事务同时删除:
  - `task_servers` 表中关联记录
  - `tasks` 表中任务记录
- 删除顺序: 先删 `task_servers`，再删 `tasks`

## 8. 数据库迁移

启动时自动执行迁移:

```go
func migrateSchema(db *sql.DB) error {
    // servers表: 添加 connection_status 字段
    _, err := db.Exec(`
        ALTER TABLE servers ADD COLUMN connection_status TEXT DEFAULT 'unknown'
    `)
    if err != nil && !strings.Contains(err.Error(), "duplicate column") {
        return err
    }
    return nil
}
```

**注意**: 迁移逻辑使用 `ALTER TABLE`，SQLite需支持 FTS5 等扩展。

## 9. 安全特性

- **密码加密**: AES-256-GCM加密存储
- **密码不返回**: JSON序列化时忽略 `password_encrypted` 字段 (使用 `json:"-"`)
- **CORS**: 允许所有来源访问(开发环境)

## 10. 日志系统

```go
logger.Debug(format, args...)  // Debug级别
logger.Info(format, args...)   // Info级别
logger.Warn(format, args...)  // Warn级别
logger.Error(format, args...)  // Error级别
```

**日志格式**:
```
2026-03-31 17:30:19.024 [INFO]  NodePilot starting on :8080
```

## 11. 依赖库

| 库 | 版本 | 用途 |
|----|------|------|
| github.com/gin-gonic/gin | latest | HTTP框架 |
| github.com/mattn/go-sqlite3 | latest | SQLite驱动 |
| github.com/pkg/sftp | latest | SFTP客户端 |
| golang.org/x/crypto/ssh | latest | SSH客户端 |
| github.com/gorilla/websocket | latest | WebSocket |

## 12. 构建部署

```bash
# 构建
cd backend
go build -o node-pilot ./cmd/server

# 运行
./node-pilot --db ../data/servers.db --listen :8080

# 或使用启动脚本
./scripts/start.sh [IP] [PORT] [debug]
```

## 13. 数据模型

```go
type Server struct {
    ID                int64
    Name              string
    Host              string
    Port              int
    Username          string
    PasswordEncrypted string    `json:"-"`
    ConnectionStatus  string    // online|offline|unknown
    CreatedAt         time.Time
    UpdatedAt         time.Time
}

type Script struct {
    ID          int64
    Name        string
    Description string
    Content     string
    TargetPath  string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type Task struct {
    ID         int64
    ScriptID   int64
    Name       string
    Status     string  // pending|running|completed|cancelled|failed
    CreatedAt  time.Time
    StartedAt  *time.Time
    FinishedAt *time.Time
}

type TaskServer struct {
    ID         int64
    TaskID     int64
    ServerID   int64
    ServerName string
    Status     string  // pending|running|success|failed
    Output     string
    Error      string
    StartedAt  *time.Time
    FinishedAt *time.Time
}

type WSMessage struct {
    Type       string
    TaskID     int64
    ServerID   int64
    ServerName string
    Content    string
    Status     string
    ExitCode   int
    Timestamp  time.Time
    Total      int
    Success    int
    Failed     int
}
```
