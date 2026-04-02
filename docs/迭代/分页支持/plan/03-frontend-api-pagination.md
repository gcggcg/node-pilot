# 前端API层分页支持

## 任务描述

修改前端 API 层，支持分页参数并处理新的返回格式。

## 详细说明

### 1. 修改 API 类型定义

在 `frontend/src/types/index.ts` 添加分页响应类型：

```typescript
export interface PaginatedResponse<T> {
    data: T[];
    total: number;
    page: number;
    pageSize: number;
}

export interface PaginationParams {
    page?: number;
    pageSize?: number;
}
```

### 2. 修改 api/index.ts

```typescript
export const serverApi = {
    list: (params?: PaginationParams) => {
        const query = params 
            ? `?page=${params.page || 1}&pageSize=${params.pageSize || 10}`
            : '';
        return api.get<any, PaginatedResponse<any>>(`/servers${query}`);
    },
    // ... 其他方法
};

export const scriptApi = {
    list: (params?: PaginationParams) => {
        const query = params 
            ? `?page=${params.page || 1}&pageSize=${params.pageSize || 10}`
            : '';
        return api.get<any, PaginatedResponse<any>>(`/scripts${query}`);
    },
    // ... 其他方法
};

export const taskApi = {
    list: (params?: PaginationParams) => {
        const query = params 
            ? `?page=${params.page || 1}&pageSize=${params.pageSize || 10}`
            : '';
        return api.get<any, PaginatedResponse<any>>(`/tasks${query}`);
    },
    // ... 其他方法
};
```

## 输入

- `frontend/src/api/index.ts`
- `frontend/src/types/index.ts`

## 输出

- 修改后的 API 层，支持分页参数和 PaginatedResponse 处理

## 依赖

- 02

## 验收标准

- [ ] API 方法支持传入分页参数
- [ ] 返回类型包含 PaginatedResponse 泛型
- [ ] list 接口返回 { data, total, page, pageSize }
