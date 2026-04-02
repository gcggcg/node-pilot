# 创建通用分页组件

## 任务描述

创建可复用的 Vue 分页组件，支持统一的分页交互体验。

## 详细说明

### 1. 创建 Pagination.vue 组件

创建 `frontend/src/components/Pagination.vue`：

```vue
<template>
    <div class="pagination">
        <div class="pagination-info">
            共 {{ total }} 条记录，当前 {{ currentPage }}/{{ totalPages }} 页
        </div>
        <div class="pagination-controls">
            <select v-model="pageSize" @change="onPageSizeChange" class="page-size-select">
                <option :value="10">10条/页</option>
                <option :value="20">20条/页</option>
                <option :value="50">50条/页</option>
                <option :value="100">100条/页</option>
            </select>
            <button @click="goToFirst" :disabled="currentPage === 1" class="page-btn">首页</button>
            <button @click="goToPrev" :disabled="currentPage === 1" class="page-btn">上一页</button>
            <span class="page-numbers">
                <button 
                    v-for="page in visiblePages" 
                    :key="page"
                    @click="goToPage(page)"
                    :class="['page-num', { active: page === currentPage }]"
                >
                    {{ page }}
                </button>
            </span>
            <button @click="goToNext" :disabled="currentPage === totalPages" class="page-btn">下一页</button>
            <button @click="goToLast" :disabled="currentPage === totalPages" class="page-btn">末页</button>
            <input 
                type="number" 
                v-model.number="jumpPage" 
                @keyup.enter="handleJump"
                class="jump-input"
                min="1"
                :max="totalPages"
                placeholder="跳转"
            />
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';

const props = defineProps<{
    currentPage: number;
    pageSize: number;
    total: number;
}>();

const emit = defineEmits<{
    (e: 'update:currentPage', page: number): void;
    (e: 'update:pageSize', size: number): void;
    (e: 'change', payload: { page: number; pageSize: number }): void;
}>();

const jumpPage = ref(props.currentPage);

const totalPages = computed(() => Math.ceil(props.total / props.pageSize) || 1);

const visiblePages = computed(() => {
    const pages: number[] = [];
    const maxVisible = 5;
    let start = Math.max(1, props.currentPage - Math.floor(maxVisible / 2));
    let end = Math.min(totalPages.value, start + maxVisible - 1);
    if (end - start < maxVisible - 1) {
        start = Math.max(1, end - maxVisible + 1);
    }
    for (let i = start; i <= end; i++) {
        pages.push(i);
    }
    return pages;
});

function goToPage(page: number) {
    if (page >= 1 && page <= totalPages.value) {
        emit('update:currentPage', page);
        emit('change', { page, pageSize: props.pageSize });
    }
}

function goToFirst() { goToPage(1); }
function goToPrev() { goToPage(props.currentPage - 1); }
function goToNext() { goToPage(props.currentPage + 1); }
function goToLast() { goToPage(totalPages.value); }

function handleJump() {
    const page = jumpPage.value;
    if (page >= 1 && page <= totalPages.value) {
        goToPage(page);
    }
    jumpPage.value = props.currentPage;
}

function onPageSizeChange() {
    emit('update:pageSize', props.pageSize);
    emit('change', { page: 1, pageSize: props.pageSize });
}

watch(() => props.currentPage, (val) => {
    jumpPage.value = val;
});
</script>

<style scoped>
.pagination {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px;
    background: white;
    border-top: 1px solid #eee;
}

.pagination-info {
    color: #666;
    font-size: 14px;
}

.pagination-controls {
    display: flex;
    align-items: center;
    gap: 8px;
}

.page-size-select {
    padding: 6px 10px;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 14px;
}

.page-btn {
    padding: 6px 12px;
    border: 1px solid #ddd;
    border-radius: 4px;
    background: white;
    cursor: pointer;
    font-size: 14px;
}

.page-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}

.page-numbers {
    display: flex;
    gap: 4px;
}

.page-num {
    min-width: 32px;
    height: 32px;
    border: 1px solid #ddd;
    border-radius: 4px;
    background: white;
    cursor: pointer;
    font-size: 14px;
}

.page-num.active {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    border-color: #667eea;
}

.jump-input {
    width: 60px;
    padding: 6px 10px;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 14px;
}
</style>
```

## 输入

- 需求文档 `04-支持分页.md`
- 当前 `frontend/src/views/*.vue` 列表组件

## 输出

- `frontend/src/components/Pagination.vue` - 通用分页组件

## 依赖

- 无

## 验收标准

- [ ] 分页组件支持首页、末页、上一页、下一页
- [ ] 分页组件支持跳页输入框
- [ ] 分页组件支持每页条数选择器（10/20/50/100）
- [ ] 分页组件显示总记录数和当前页/总页数
- [ ] 组件通过 emit 事件传递分页参数变化
