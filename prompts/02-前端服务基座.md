# NodePilot 前端技术规格书

## 1. 项目概述

- **项目名称**: NodePilot 前端
- **技术栈**: Vue 3 + TypeScript + Vite + Pinia
- **功能定位**: 批量服务器管理平台的Web管理界面

## 2. 项目结构

```
frontend/
├── src/
│   ├── api/
│   │   └── index.ts          # API调用封装(Axios)
│   ├── components/
│   │   ├── NavBar.vue         # 顶部导航栏
│   │   └── OutputPanel.vue    # 输出面板组件
│   ├── router/
│   │   └── index.ts          # Vue Router配置(7个路由)
│   ├── stores/
│   │   ├── server.ts         # 服务器状态管理
│   │   ├── script.ts         # 脚本状态管理
│   │   └── task.ts           # 任务状态管理+WebSocket
│   ├── types/
│   │   └── index.ts          # TypeScript类型定义
│   ├── views/
│   │   ├── ServerList.vue    # 服务器列表(含批量删除)
│   │   ├── ServerForm.vue    # 服务器表单(新建/编辑)
│   │   ├── ScriptList.vue    # 脚本列表(含批量删除)
│   │   ├── ScriptForm.vue    # 脚本表单(新建/编辑)
│   │   ├── TaskList.vue      # 任务列表(含批量删除)
│   │   └── TaskOutput.vue    # 任务输出(实时+历史)
│   ├── App.vue               # 根组件
│   └── main.ts               # 入口文件
├── package.json
└── vite.config.ts
```

## 3. 技术栈

| 技术 | 版本 | 用途 |
|------|------|------|
| Vue | 3.4+ | 渐进式JavaScript框架 |
| TypeScript | 5.x | 类型安全 |
| Vite | 5.x | 构建工具 |
| Pinia | 2.x | 状态管理 |
| Vue Router | 4.x | 路由管理 |
| Axios | 1.x | HTTP客户端 |

## 4. 路由配置

| 路径 | 组件 | 说明 |
|------|------|------|
| `/` | redirect to `/servers` | 首页重定向 |
| `/servers` | ServerList.vue | 服务器列表 |
| `/servers/new` | ServerForm.vue | 新建服务器 |
| `/servers/:id/edit` | ServerForm.vue | 编辑服务器 |
| `/scripts` | ScriptList.vue | 脚本列表 |
| `/scripts/new` | ScriptForm.vue | 新建脚本 |
| `/scripts/:id/edit` | ScriptForm.vue | 编辑脚本 |
| `/tasks` | TaskList.vue | 任务列表 |
| `/tasks/:id/output` | TaskOutput.vue | 任务输出 |

## 5. 数据类型

### Server
```typescript
interface Server {
    id: number;
    name: string;
    host: string;
    port: number;
    username: string;
    password_encrypted?: string;
    created_at: string;
    updated_at: string;
}
```

### Script
```typescript
interface Script {
    id: number;
    name: string;
    description: string;
    content: string;
    target_path: string;
    created_at: string;
    updated_at: string;
}
```

### Task
```typescript
interface Task {
    id: number;
    script_id: number;
    name: string;
    status: 'pending' | 'running' | 'completed' | 'cancelled' | 'failed';
    created_at: string;
    started_at?: string;
    finished_at?: string;
}
```

### TaskServer
```typescript
interface TaskServer {
    id: number;
    task_id: number;
    server_id: number;
    server_name?: string;
    status: 'pending' | 'running' | 'success' | 'failed';
    output: string;
    error: string;
    started_at?: string;
    finished_at?: string;
}
```

### WSMessage
```typescript
interface WSMessage {
    type: 'output' | 'server_start' | 'server_done' | 'task_start' | 'task_done';
    task_id: number;
    server_id?: number;
    server_name?: string;
    content?: string;
    status?: string;
    exit_code?: number;
    timestamp: string;
    total?: number;
    success?: number;
    failed?: number;
}
```

### 表单类型
```typescript
interface ServerForm {
    name: string;
    host: string;
    port: number;
    username: string;
    password: string;
}

interface ScriptForm {
    name: string;
    description: string;
    content: string;
    target_path: string;
}

interface TaskForm {
    script_id: number;
    name: string;
    server_ids: number[];
}
```

## 6. API封装

### serverApi
```typescript
export const serverApi = {
    list: () => api.get('/servers'),
    get: (id: number) => api.get(`/servers/${id}`),
    create: (data: ServerForm) => api.post('/servers', data),
    update: (id: number, data: ServerForm) => api.put(`/servers/${id}`, data),
    delete: (id: number) => api.delete(`/servers/${id}`),
    deleteMany: (ids: number[]) => api.post('/servers/batch-delete', { ids }),
    test: (id: number) => api.post(`/servers/${id}/test`),
};
```

### scriptApi
```typescript
export const scriptApi = {
    list: () => api.get('/scripts'),
    get: (id: number) => api.get(`/scripts/${id}`),
    create: (data: ScriptForm) => api.post('/scripts', data),
    update: (id: number, data: ScriptForm) => api.put(`/scripts/${id}`, data),
    delete: (id: number) => api.delete(`/scripts/${id}`),
    deleteMany: (ids: number[]) => api.post('/scripts/batch-delete', { ids }),
};
```

### taskApi
```typescript
export const taskApi = {
    list: () => api.get('/tasks'),
    get: (id: number) => api.get(`/tasks/${id}`),
    create: (data: TaskForm) => api.post('/tasks', data),
    cancel: (id: number) => api.delete(`/tasks/${id}`),
    deleteMany: (ids: number[]) => api.post('/tasks/batch-delete', { ids }),
};
```

## 7. 状态管理 (Pinia Stores)

### serverStore
```typescript
const serverStore = defineStore('server', () => {
    const servers = ref<Server[]>([]);
    const loading = ref(false);
    const error = ref<string | null>(null);
    
    // 方法
    fetchServers(), createServer(), updateServer(), deleteServer(), 
    deleteServers(), testConnection()
    
    return { servers, loading, error, ... }
});
```

### scriptStore
```typescript
const scriptStore = defineStore('script', () => {
    const scripts = ref<Script[]>([]);
    const loading = ref(false);
    const error = ref<string | null>(null);
    
    // 方法
    fetchScripts(), createScript(), updateScript(), deleteScript(), deleteScripts()
    
    return { scripts, loading, error, ... }
});
```

### taskStore
```typescript
const taskStore = defineStore('task', () => {
    const tasks = ref<Task[]>([]);
    const currentTask = ref<Task | null>(null);
    const currentTaskId = ref<number | null>(null);
    const outputs = ref<Map<number, string>>(new Map()); // server_id -> output
    const ws = ref<WebSocket | null>(null);
    
    // 方法
    fetchTasks(), createTask(), cancelTask(), deleteTasks(),
    connectWebSocket(), disconnectWebSocket(), clearOutputs()
    
    return { tasks, currentTask, currentTaskId, outputs, ws, ... }
});
```

## 8. 组件功能

### 8.1 NavBar.vue
- 顶部紫色渐变导航栏
- 三个导航链接: 服务器 / 脚本 / 任务
- 当前页面高亮显示

### 8.2 ServerList.vue
- 服务器列表表格展示
- **批量删除**: 勾选框 + 全选 + 删除已选按钮
- 自定义紫色checkbox样式
- 操作按钮: 测试 / 编辑 / 删除
- 空状态提示

### 8.3 ServerForm.vue
- 新建/编辑服务器表单
- 字段: 名称 / IP地址 / 端口 / 用户名 / 密码
- 表单验证
- 提交时自动AES加密密码(后端处理)

### 8.4 ScriptList.vue
- 卡片式脚本列表
- **批量删除**: 卡片右上角勾选 + 删除已选按钮
- 卡片选中高亮效果
- 脚本名称 / 描述 / 目标路径
- 操作按钮: 编辑 / 删除

### 8.5 ScriptForm.vue
- 新建/编辑脚本表单
- 字段: 名称 / 描述 / 内容 / 目标路径
- 代码编辑器样式(等宽字体)

### 8.6 TaskList.vue
- 任务列表表格展示
- **批量删除**: 勾选框 + 全选 + 删除已选按钮
- 状态标签: 待执行(灰) / 运行中(蓝) / 已完成(绿) / 已取消(黄) / 失败(红)
- 创建任务弹窗: 选择脚本 + 多选服务器
- 服务器选择列表: 卡片式布局 + 自定义checkbox
- 操作: 查看输出 / 取消(仅运行中任务)

### 8.7 TaskOutput.vue
- 任务输出面板
- **历史输出**: 从API加载已完成任务的历史输出
- **实时输出**: WebSocket接收实时日志
- **TaskID过滤**: 只显示当前任务的输出
- 运行中任务自动清空旧数据
- 已完成任务显示历史数据

### 8.8 OutputPanel.vue
- 通用输出面板组件
- 按服务器分栏显示输出
- 实时追加输出内容
- 错误信息显示

## 9. UI特性

### 9.1 样式系统
- CSS变量管理主题色
- 紫色渐变按钮(#667eea → #764ba2)
- 卡片式布局
- 圆角和阴影效果

### 9.2 自定义Checkbox
```css
.checkbox-wrapper input {
    position: absolute;
    opacity: 0;
}
.checkmark {
    height: 20px;
    width: 20px;
    border: 2px solid #ddd;
    border-radius: 4px;
}
.checkbox-wrapper input:checked ~ .checkmark {
    background-color: #667eea;
    border-color: #667eea;
}
.checkmark:after {
    /* 勾选标记 */
}
```

### 9.3 服务器选择列表
- 卡片式服务器选择器
- 每行: checkbox + 服务器名称 + IP地址
- hover高亮效果

### 9.4 状态颜色
| 状态 | 背景色 | 文字色 |
|------|--------|--------|
| pending | #e9ecef | #495057 |
| running | #cce5ff | #004085 |
| completed | #d4edda | #155724 |
| cancelled | #fff3cd | #856404 |
| failed | #f8d7da | #721c24 |

## 10. WebSocket处理

### 连接建立
```typescript
function connectWebSocket(taskId: number) {
    currentTaskId.value = taskId;
    ws.value = new WebSocket(`${protocol}//${host}/ws?task_id=${taskId}`);
}
```

### 消息过滤
```typescript
function handleWSMessage(data: WSMessage) {
    // 忽略不属于当前任务的消息
    if (currentTaskId.value !== null && data.task_id !== currentTaskId.value) {
        return;
    }
    // 处理消息...
}
```

### 生命周期
- `onMounted`: 连接WebSocket
- `onUnmounted`: 断开WebSocket连接

## 11. 批量删除流程

1. 用户勾选多个项目
2. 显示"删除已选 (N)"按钮
3. 用户点击删除按钮
4. 弹出确认对话框
5. 确认后调用 `deleteMany([...ids])`
6. 刷新列表
7. 清空选中状态

## 12. 任务创建流程

1. 点击"运行任务"按钮
2. 填写任务名称
3. 选择脚本
4. 多选目标服务器(卡片式勾选)
5. 点击"运行"按钮
6. 跳转到任务列表
7. 任务自动异步执行

## 13. 输出显示逻辑

### TaskOutput.vue
```
onMounted:
    1. 获取任务详情
    2. 根据任务状态决定:
        - running: 清空旧数据, 连接WebSocket
        - completed/cancelled/failed: 加载历史输出, 连接WebSocket(仅接收新消息)
    3. 实时追加WebSocket消息到对应服务器的输出
```

## 14. 构建部署

```bash
# 安装依赖
npm install

# 开发模式
npm run dev

# 生产构建
npm run build

# 产物复制到后端
cp -r dist/* ../backend/web/
```

## 15. 启动脚本

```bash
# 一键启动(自动构建)
./scripts/start.sh [IP] [PORT] [debug]

# 示例
./scripts/start.sh 127.0.0.1 8080 debug
```
