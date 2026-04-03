# 文件管理 - 前端类型定义和 API

## 任务描述

为文件上传管理功能添加 TypeScript 类型定义和 API 调用封装。

## 详细说明

### 1. 更新 Types

在 `frontend/src/types/index.ts` 中添加：

```typescript
export interface FileUpload {
    id: number;
    name: string;
    local_path: string;
    remote_path: string;
    status: 'pending' | 'success' | 'failed';
    created_at: string;
    updated_at: string;
}

export interface FileUploadServer {
    id: number;
    file_upload_id: number;
    server_id: number;
    server_name: string;
    status: 'pending' | 'success' | 'failed';
    error_message: string;
    file_name: string;
    remote_full_path: string;
    created_at: string;
}

export interface FileUploadForm {
    name: string;
    files: File[];  // 上传的文件列表
    serverIds: number[];
    remotePath: string;
}

export interface FileUploadListResponse {
    data: FileUpload[];
    total: number;
    page: number;
    pageSize: number;
}

export interface FileUploadResultResponse {
    data: FileUploadServer[];
}
```

### 2. 更新 API

在 `frontend/src/api/index.ts` 中添加 fileUploadApi：

```typescript
export const fileUploadApi = {
    list: (params?: {
        page?: number;
        pageSize?: number;
        status?: string;
        fileName?: string;
        startTime?: string;
        endTime?: string;
    }) => {
        const query = params
            ? `?page=${params.page || 1}&pageSize=${params.pageSize || 10}${params.status ? `&status=${params.status}` : ''}${params.fileName ? `&fileName=${params.fileName}` : ''}${params.startTime ? `&startTime=${params.startTime}` : ''}${params.endTime ? `&endTime=${params.endTime}` : ''}`
            : '';
        return api.get<any, FileUploadListResponse>(`/v1/file-uploads${query}`);
    },
    
    create: (formData: FormData) => {
        return api.post('/v1/file-uploads', formData, {
            headers: {
                'Content-Type': 'multipart/form-data'
            }
        });
    },
    
    update: (id: number, data: { serverIds: number[]; remotePath: string }) => {
        return api.put(`/v1/file-uploads/${id}`, data);
    },
    
    execute: (id: number) => {
        return api.post(`/v1/file-uploads/${id}/execute`);
    },
    
    getResults: (id: number) => {
        return api.get<any, FileUploadResultResponse>(`/v1/file-uploads/${id}/results`);
    },
    
    delete: (ids: number[]) => {
        return api.delete('/v1/file-uploads', { data: { ids } });
    },
};
```

## 输入

- 需求文档 `06-新增文件管理.md`
- 现有 `frontend/src/types/index.ts`
- 现有 `frontend/src/api/index.ts` 模式

## 输出

- `frontend/src/types/index.ts` - 添加 FileUpload 相关类型
- `frontend/src/api/index.ts` - 添加 fileUploadApi

## 依赖

- 04-backend-handlers-routes.md

## 验收标准

- [ ] 所有类型定义完整
- [ ] API 方法正确实现
- [ ] 支持 multipart/form-data 上传
- [ ] 符合现有项目代码风格
