# 文件管理 - 前端 Pinia Store

## 任务描述

创建文件上传管理的 Pinia 状态管理。

## 详细说明

### 1. 创建 Store 文件

创建 `frontend/src/stores/fileupload.ts`：

### 2. useFileUploadStore 定义

```typescript
import { defineStore } from 'pinia';
import { ref } from 'vue';
import { fileUploadApi } from '@/api';
import type { FileUpload, FileUploadServer } from '@/types';

export const useFileUploadStore = defineStore('fileupload', () => {
    const uploads = ref<FileUpload[]>([]);
    const results = ref<FileUploadServer[]>([]);
    const loading = ref(false);
    const pagination = ref({ page: 1, pageSize: 10, total: 0 });

    async function fetchUploads(page = 1, pageSize = 10, filters = {}) {
        loading.value = true;
        try {
            const res = await fileUploadApi.list({ page, pageSize, ...filters });
            uploads.value = res.data;
            pagination.value = { page: res.page, pageSize: res.pageSize, total: res.total };
        } finally {
            loading.value = false;
        }
    }

    async function createUpload(formData: FormData) {
        await fileUploadApi.create(formData);
        await fetchUploads();
    }

    async function updateUpload(id: number, data: { serverIds: number[]; remotePath: string }) {
        await fileUploadApi.update(id, data);
        await fetchUploads(pagination.value.page, pagination.value.pageSize);
    }

    async function executeUpload(id: number) {
        await fileUploadApi.execute(id);
    }

    async function fetchResults(id: number) {
        const res = await fileUploadApi.getResults(id);
        results.value = res.data;
    }

    async function deleteUploads(ids: number[]) {
        await fileUploadApi.delete(ids);
        await fetchUploads(pagination.value.page, pagination.value.pageSize);
    }

    return {
        uploads,
        results,
        loading,
        pagination,
        fetchUploads,
        createUpload,
        updateUpload,
        executeUpload,
        fetchResults,
        deleteUploads,
    };
});
```

## 输入

- 需求文档 `06-新增文件管理.md`
- 现有 `frontend/src/stores/auth.ts` 模式
- `05-frontend-types-api.md`

## 输出

- `frontend/src/stores/fileupload.ts` - 文件上传状态管理

## 依赖

- 05-frontend-types-api.md

## 验收标准

- [ ] Store 实现完整
- [ ] 状态管理正确
- [ ] 方法完整
- [ ] 符合现有项目代码风格
