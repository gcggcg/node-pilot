# 前端Store分页状态管理

## 任务描述

修改前端 Pinia Store，支持分页状态管理。

## 详细说明

### 1. 修改 server.ts store

```typescript
import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { serverApi } from '@/api';
import type { Server, ServerForm, PaginatedResponse } from '@/types';

export const useServerStore = defineStore('server', () => {
    const servers = ref<Server[]>([]);
    const loading = ref(false);
    const error = ref<string | null>(null);
    const pagination = ref({
        page: 1,
        pageSize: 10,
        total: 0
    });

    async function fetchServers(page = 1, pageSize = 10) {
        loading.value = true;
        error.value = null;
        try {
            const res = await serverApi.list({ page, pageSize });
            servers.value = res.data;
            pagination.value = {
                page: res.page,
                pageSize: res.pageSize,
                total: res.total
            };
        } catch (e: any) {
            error.value = e.message;
        } finally {
            loading.value = false;
        }
    }

    // ... 其他方法保持不变，只需将 fetchServers() 调用改为带分页参数
});
```

### 2. 类似修改 script.ts 和 task.ts stores

## 输入

- `frontend/src/stores/server.ts`
- `frontend/src/stores/script.ts`
- `frontend/src/stores/task.ts`

## 输出

- 修改后的 Store，支持分页状态（page, pageSize, total）

## 依赖

- 03

## 验收标准

- [ ] Store 包含 pagination 状态（page, pageSize, total）
- [ ] fetchServers/scripts/tasks 支持分页参数
- [ ] 分页状态在数据获取后更新
