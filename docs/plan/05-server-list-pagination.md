# 服务器管理列表分页

## 任务描述

为 ServerList.vue 添加分页组件和数据加载逻辑。

## 详细说明

### 1. 引入分页组件和 useRouter

```typescript
import { ref, onMounted, watch } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useServerStore } from '@/stores/server';
import Pagination from '@/components/Pagination.vue';

const router = useRouter();
const route = useRoute();
const store = useServerStore();

const pagination = ref({
    page: Number(route.query.page) || 1,
    pageSize: Number(route.query.pageSize) || 10,
    total: 0
});

onMounted(() => {
    loadData();
});

function loadData() {
    store.fetchServers(pagination.value.page, pagination.value.pageSize);
    pagination.value.total = store.pagination.total;
}

function handlePageChange(payload: { page: number; pageSize: number }) {
    pagination.value.page = payload.page;
    pagination.value.pageSize = payload.pageSize;
    router.replace({ 
        query: { page: payload.page, pageSize: payload.pageSize } 
    });
    store.fetchServers(payload.page, payload.pageSize);
}
```

### 2. 模板添加分页组件

```vue
<table v-else class="table">
    <!-- ... 表格内容 ... -->
</table>

<Pagination 
    :current-page="pagination.page"
    :page-size="pagination.pageSize"
    :total="pagination.total"
    @change="handlePageChange"
/>
```

## 输入

- `frontend/src/views/ServerList.vue`
- `frontend/src/components/Pagination.vue`

## 输出

- ServerList.vue 支持 URL 参数同步分页状态

## 依赖

- 04

## 验收标准

- [ ] 分页参数通过 URL query 同步
- [ ] 切换页码/页大小时更新 URL
- [ ] 页面初始化时从 URL 读取分页参数
- [ ] 显示总记录数和分页信息
