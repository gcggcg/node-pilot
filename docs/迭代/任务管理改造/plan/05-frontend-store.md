# 任务管理改造 - Frontend Store

## 任务描述

在任务管理 Store 中添加 `executeTask` 和 `updateTask` 方法。

## 详细说明

### 1. 更新 task store

在 `frontend/src/stores/task.ts` 中添加：

```typescript
async function executeTask(id: number) {
    loading.value = true;
    error.value = null;
    try {
        await taskApi.execute(id);
        await fetchTasks(pagination.value.page, pagination.value.pageSize);
    } catch (e: any) {
        error.value = e.message;
        throw e;
    } finally {
        loading.value = false;
    }
}

async function updateTask(id: number, data: TaskForm) {
    loading.value = true;
    error.value = null;
    try {
        await taskApi.update(id, data);
        await fetchTasks(pagination.value.page, pagination.value.pageSize);
    } catch (e: any) {
        error.value = e.message;
        throw e;
    } finally {
        loading.value = false;
    }
}
```

### 2. 更新 return 语句

在 store 的 return 中添加新方法：

```typescript
return {
    // ... existing exports
    executeTask,
    updateTask,
}
```

## 输入

- `frontend/src/stores/task.ts`

## 输出

- 修改后的 `frontend/src/stores/task.ts`

## 依赖

- 04-frontend-types-api

## 验收标准

- [ ] executeTask 方法正确调用 API
- [ ] updateTask 方法正确调用 API
- [ ] 方法在 store 返回中正确导出
