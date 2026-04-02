# 任务管理改造 - Frontend Types & API

## 任务描述

更新前端 Types 和 API 定义，添加任务执行和更新相关类型与方法。

## 详细说明

### 1. 更新 TaskForm 类型

在 `frontend/src/types/index.ts` 中修改 TaskForm：

```typescript
// TaskForm 用于创建和更新任务
export interface TaskForm {
    script_id: number;
    name: string;
    server_ids: number[];
}
```

### 2. 更新 taskApi

在 `frontend/src/api/index.ts` 中添加：

```typescript
export const taskApi = {
    list: (params?: PaginationParams) => {
        const query = params 
            ? `?page=${params.page || 1}&pageSize=${params.pageSize || 10}`
            : '';
        return api.get<any, PaginatedResponse<any>>(`/tasks${query}`);
    },
    get: (id: number) => api.get<any, any>(`/tasks/${id}`),
    create: (data: TaskForm) => api.post('/tasks', {
        script_id: data.script_id,
        name: data.name,
        server_ids: data.server_ids
    }),
    update: (id: number, data: TaskForm) => api.put(`/tasks/${id}`, {
        script_id: data.script_id,
        name: data.name,
        server_ids: data.server_ids
    }),
    execute: (id: number) => api.post(`/tasks/${id}/execute`),
    cancel: (id: number) => api.delete(`/tasks/${id}`),
    deleteMany: (ids: number[]) => api.post('/tasks/batch-delete', { ids }),
};
```

## 输入

- `frontend/src/types/index.ts`
- `frontend/src/api/index.ts`

## 输出

- 修改后的 `frontend/src/types/index.ts`
- 修改后的 `frontend/src/api/index.ts`

## 依赖

- 03-backend-handlers-routes

## 验收标准

- [ ] TaskForm 类型正确
- [ ] taskApi.execute 方法正确调用 POST /tasks/:id/execute
- [ ] taskApi.update 方法正确调用 PUT /tasks/:id
