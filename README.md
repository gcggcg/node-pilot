# NodePilot

> 🚀 批量服务器管理平台 - 一键部署、批量操作、实时监控

[![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue)](https://golang.org/)
[![Vue Version](https://img.shields.io/badge/Vue-3.4%2B-green)](https://vuejs.org/)
[![License](https://img.shields.io/badge/license-MIT-yellow)](LICENSE)
[![Stars](https://img.shields.io/github/stars/gcggcg/node-pilot?style=social)](https://github.com/gcggcg/node-pilot)

[English](./README.md) | 中文

<!-- TOC -->

- [特性](#特性)
- [快速开始](#快速开始)
- [架构设计](#架构设计)
- [适用场景](#适用场景)
- [API 参考](#api-参考)
- [开发指南](#开发指南)
- [安全说明](#安全说明)
- [贡献指南](#贡献指南)
- [许可证](#许可证)

<!-- /TOC -->

## ✨ 特性

### 🎯 核心能力

- **批量服务器管理** - 集中管理多台服务器的连接信息
- **批量脚本执行** - 同时在多台服务器上执行相同的脚本
- **实时输出回显** - 类似 `kubectl logs -f` 的实时日志查看
- **连接状态监控** - 一键测试服务器连接状态（在线/离线/未知）

### 🔒 安全特性

- **AES-256-GCM 加密** - 密码安全存储，不在客户端暴露
- **会话隔离** - 每个任务独立执行，互不干扰

### 🚀 性能优势

- **单端口部署** - 前端嵌入二进制，无需 Nginx/Apache
- **批量并发执行** - 支持 10 台服务器同时执行（可配置批次大小）
- **WebSocket 实时推送** - 无需轮询，即时获取执行结果

### 💻 用户体验

- **Web UI** - 简洁易用的管理界面
- **一键启动** - 下载即运行，无需复杂配置
- **跨平台** - 支持 Linux/macOS/Windows
- **分页支持** - 服务器/脚本/任务列表支持分页浏览

## 📦 快速开始

### 二进制下载

```bash
# Linux/macOS
curl -fsSL https://github.com/gcggcg/node-pilot/releases/latest/download/node-pilot-linux-amd64 -o node-pilot
chmod +x node-pilot

# Windows
# 从 https://github.com/gcggcg/node-pilot/releases 下载

# 运行
./node-pilot --db ./data/servers.db --listen :8080
```

访问 `http://localhost:8080`

### 源码构建

**环境要求:**

- Go 1.21+
- Node.js 18+ (仅前端开发需要)

```bash
# 克隆仓库
git clone https://github.com/gcggcg/node-pilot.git
cd node-pilot

# 构建后端
cd backend
go build -o node-pilot ./cmd/server

# 构建前端
cd ../frontend
npm install
npm run build

# 复制前端到后端 web 目录
cp -r dist/* ../backend/web/

# 运行
cd ../backend
./node-pilot --db ../data/servers.db --listen :8080
```

### 一键启动脚本

```bash
cd node-pilot
chmod +x scripts/start.sh
./scripts/start.sh
```

### Docker 部署

```bash
# 拉取镜像
docker pull gcggcg/node-pilot:latest

# 运行
docker run -d -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  node-pilot:latest
```

## 🏗 架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                        浏览器端                              │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      NodePilot 服务端                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │
│  │  REST API   │  │  WebSocket  │  │   静态资源           │ │
│  │  (Gin)      │  │  (Hub)      │  │   (嵌入式前端)       │ │
│  └─────────────┘  └─────────────┘  └─────────────────────┘ │
│         │                │                    │            │
│         ▼                ▼                    │            │
│  ┌─────────────────────────────────────────┐ │            │
│  │              业务逻辑层                   │ │            │
│  │  ┌─────────────┐  ┌─────────────────┐   │ │            │
│  │  │ SSH 连接池  │  │ 任务执行器       │   │ │            │
│  │  │ (连接管理)  │  │ (异步批量执行)   │   │ │            │
│  │  └─────────────┘  └─────────────────┘   │ │            │
│  └─────────────────────────────────────────┘ │            │
│         │                                      │            │
│         ▼                                      │            │
│  ┌─────────────────────────────────────────┐ │            │
│  │              数据访问层                    │ │            │
│  │              (SQLite)                    │ │            │
│  └─────────────────────────────────────────┘ │            │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      目标服务器集群                          │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐                  │
│  │ 服务器 1 │  │ 服务器 2 │  │ 服务器 N │  ...             │
│  │ (SSH+SFTP)│  │ (SSH+SFTP)│  │ (SSH+SFTP)│                  │
│  └──────────┘  └──────────┘  └──────────┘                  │
└─────────────────────────────────────────────────────────────┘
```

### 技术栈

| 层级        | 技术                        |
|-----------|---------------------------|
| 后端        | Go + Gin                  |
| 前端        | Vue 3 + TypeScript + Vite |
| 数据库       | SQLite (嵌入式)              |
| SSH       | golang.org/x/crypto/ssh   |
| WebSocket | gorilla/websocket         |
| 密码加密      | AES-256-GCM               |

## 💡 适用场景

### 1. 微服务组件集群一键搭建

```bash
# 创建微服务部署脚本列表
#!/bin/bash
cd /opt/myapp
git pull origin main
docker-compose up -d --build

# 在 NodePilot UI 中选择 10 台服务器
# 一键在所有服务器上执行该脚本
```

### 2. OpenClaw 小龙虾集群一键搭建

```bash
#!/bin/bash
# 一键初始化集群
apt-get update && apt-get install -y docker.io
docker pull openclust/node:latest
docker run -d --name master openclust/node --role=master
```

### 3. AI 编程环境集群一键搭建

```bash
#!/bin/bash
# 安装 OpenCode AI 开发环境
curl -fsSL https://setup.opencode.com | bash
opencode config set api-key your-key
```

### 4. 系统批量维护

```bash
#!/bin/bash
# 批量系统更新
apt-get update && apt-get upgrade -y
docker system prune -f
journalctl --vacuum-time=7d
```

### 5. 定制化批量服务处理

## 🌐 API 参考

### 服务器管理

| 方法     | 端点                          | 说明      | 分页支持 |
|--------|-----------------------------|---------|------|
| GET    | `/api/servers`              | 获取服务器列表 | ✅    |
| GET    | `/api/servers/:id`          | 获取服务器详情 |      |
| POST   | `/api/servers`              | 创建服务器   |      |
| PUT    | `/api/servers/:id`          | 更新服务器   |      |
| DELETE | `/api/servers/:id`          | 删除服务器   |      |
| POST   | `/api/servers/:id/test`     | 测试连接    |      |
| POST   | `/api/servers/batch-delete` | 批量删除    |      |

### 脚本管理

| 方法     | 端点                          | 说明     | 分页支持 |
|--------|-----------------------------|--------|------|
| GET    | `/api/scripts`              | 获取脚本列表 | ✅    |
| GET    | `/api/scripts/:id`          | 获取脚本详情 |      |
| POST   | `/api/scripts`              | 创建脚本   |      |
| PUT    | `/api/scripts/:id`          | 更新脚本   |      |
| DELETE | `/api/scripts/:id`          | 删除脚本   |      |
| POST   | `/api/scripts/batch-delete` | 批量删除   |      |

### 任务管理

| 方法     | 端点                        | 说明        | 分页支持 |
|--------|---------------------------|-----------|------|
| GET    | `/api/tasks`              | 获取任务列表    | ✅    |
| GET    | `/api/tasks/:id`          | 获取任务详情    |      |
| POST   | `/api/tasks`              | 创建并执行任务   |      |
| DELETE | `/api/tasks/:id`          | 取消任务      |      |
| POST   | `/api/tasks/batch-delete` | 批量删除      |      |
| GET    | `/api/tasks/:id/output`   | SSE 实时输出流 |      |

### 分页参数

列表接口支持分页查询，参数通过 URL query 传递：

| 参数     | 类型    | 默认值 | 说明    |
|---------|-------|------|-------|
| `page`    | int   | 1    | 页码    |
| `pageSize` | int   | 10   | 每页条数 |

**分页响应格式：**

```json
{
  "data": [...],
  "total": 100,
  "page": 1,
  "pageSize": 10
}
```

### WebSocket

```
ws://localhost:8080/ws?task_id=123
```

**消息类型:**

- `task_start` - 任务开始执行
- `server_start` - 开始在服务器上执行
- `output` - 实时输出
- `server_done` - 单服务器执行完成
- `task_done` - 全部服务器执行完成

## 🔧 开发指南

### 项目结构

```
node-pilot/
├── backend/
│   ├── cmd/server/
│   │   └── main.go           # 程序入口
│   ├── internal/
│   │   ├── config/          # 配置管理
│   │   ├── handler/          # HTTP 处理器
│   │   ├── model/            # 数据模型
│   │   ├── repository/       # 数据库操作
│   │   ├── service/          # 业务逻辑
│   │   │   ├── ssh.go       # SSH 连接池
│   │   │   └── task.go      # 任务执行器
│   │   └── websocket/        # WebSocket Hub
│   ├── web/                   # 嵌入式前端
│   └── data/                  # SQLite 数据库
├── frontend/
│   └── src/
│       ├── api/              # API 调用封装
│       ├── components/       # Vue 组件 (含 Pagination.vue)
│       ├── views/            # 页面视图
│       ├── stores/           # Pinia 状态管理
│       └── types/            # TypeScript 类型
├── docs/
│   ├── plan/                 # 任务计划文档
│   └── review/               # 代码审查报告
└── scripts/
    └── start.sh             # 启动脚本
```

### 环境变量

| 变量                  | 默认值                 | 说明              |
|---------------------|---------------------|-----------------|
| `NODE_PILOT_DB`     | `./data/servers.db` | 数据库路径           |
| `NODE_PILOT_LISTEN` | `:8080`             | 监听地址            |
| `NODE_PILOT_KEY`    | (自动生成)              | AES 加密密钥 (32字节) |

### 运行测试

```bash
# 后端测试
cd backend
go test ./...

# 前端测试
cd frontend
npm test
```

## 🔒 安全说明

### 密码存储

所有密码使用 AES-256-GCM 加密后存储：

```go
// 32字节密钥用于 AES-256
key := []byte("your-32-byte-encryption-key")
```

### 最佳实践

1. **数据库权限** - 设置数据库文件权限为 `600`
2. **网络隔离** - 仅在内网环境使用
3. **SSH 密钥** - 生产环境建议使用 SSH 密钥认证替代密码
4. **加密密钥** - 生产环境使用随机生成的强密钥

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！请在提交前阅读 [贡献指南](./CONTRIBUTING.md)。

### 开发环境搭建

```bash
# Fork 并克隆仓库
git clone https://github.com/gcggcg/node-pilot.git

# 创建功能分支
git checkout -b feature/your-feature-name

# 修改代码并测试
go test ./...
npm test

# 提交并推送
git commit -m "feat: 添加新功能"
git push origin feature/your-feature-name

# 创建 Pull Request
```

### 代码规范

- Go: 遵循 [Go 编码规范](https://go.dev/doc/effective_go)
- JavaScript/TypeScript: 遵循项目 ESLint 配置

## 📄 许可证

MIT 许可证 - 详见 [LICENSE](./LICENSE) 文件。

## 🙏 致谢

- [gin-gonic/gin](https://github.com/gin-gonic/gin) - 高性能 HTTP Web 框架
- [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto) - SSH 和加密库
- [vuejs/core](https://github.com/vuejs/core) - 渐进式 JavaScript 框架

---

<p align="center">
  使用 ❤️ 开发 | <a href="https://github.com/gcggcg">光之翼</a>
</p>
