# 前端 TypeScript 类型定义

## 任务描述

在 TypeScript 接口中添加 `execution_mode` 字段，支持前端类型检查和表单数据绑定。

## 详细说明

### 1. 修改 Task 接口 (`frontend/src/types/index.ts`)

```typescript
export interface Task {
    id: number;
    script_id: number;
    script_ids: string;
    name: string;
    status: 'pending' | 'running' | 'completed' | 'cancelled' | 'failed';
    execution_mode: 'concurrent' | 'sequential';  // 新增
    created_at: string;
    started_at?: string;
    finished_at?: string;
}
```

### 2. 修改 TaskForm 接口 (`frontend/src/types/index.ts`)

```typescript
export interface TaskForm {
    script_id?: number;
    script_ids?: string;
    name: string;
    server_ids: number[];
    execution_mode: 'concurrent' | 'sequential';  // 新增
}
```

## 输入

- 需求文档：`10-多服务器任务支持单线程和多线程切换.md`
- 现有类型文件：`frontend/src/types/index.ts`

## 输出

- 修改后的 `frontend/src/types/index.ts`

## 依赖

- 无（前端独立修改）

## 验收标准

- [ ] Task 接口包含 `execution_mode: 'concurrent' | 'sequential'`
- [ ] TaskForm 接口包含 `execution_mode: 'concurrent' | 'sequential'`
- [ ] 类型定义与后端模型一致
