# 前端 API 层修改

## 任务描述

在 API 请求函数中添加 `execution_mode` 参数，支持创建和更新任务时传递执行模式。

## 详细说明

### 1. 修改创建任务 API (`frontend/src/api/index.ts` 或相关 API 文件)

找到 `createTask` 函数，添加 `execution_mode` 到请求体：
```typescript
export const createTask = async (data: {
    script_ids: string;
    name: string;
    server_ids: number[];
    execution_mode: 'concurrent' | 'sequential';
}) => {
    const response = await fetch('/api/tasks', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
    });
    return response.json();
};
```

### 2. 修改更新任务 API (`frontend/src/api/index.ts` 或相关 API 文件)

找到 `updateTask` 函数，添加 `execution_mode` 到请求体：
```typescript
export const updateTask = async (id: number, data: {
    script_ids: string;
    name: string;
    server_ids: number[];
    execution_mode: 'concurrent' | 'sequential';
}) => {
    const response = await fetch(`/api/tasks/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
    });
    return response.json();
};
```

### 3. 确保任务列表 API 返回包含 execution_mode

检查 `fetchTasks` 或类似函数是否正确处理返回的 `execution_mode` 字段。

## 输入

- 需求文档：`10-多服务器任务支持单线程和多线程切换.md`
- 现有 API 文件：`frontend/src/api/index.ts` 或相关文件

## 输出

- 修改后的 API 文件

## 依赖

- 04

## 验收标准

- [ ] createTask 请求体包含 execution_mode
- [ ] updateTask 请求体包含 execution_mode
- [ ] API 调用时 execution_mode 默认为 'concurrent'
- [ ] 任务列表 API 正确返回 execution_mode 字段
