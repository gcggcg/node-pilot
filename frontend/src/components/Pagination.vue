<template>
    <div class="pagination">
        <div class="pagination-info">
            共 {{ total }} 条记录，当前 {{ currentPage }}/{{ totalPages }} 页
        </div>
        <div class="pagination-controls">
            <select v-model="localPageSize" @change="onPageSizeChange($event)" class="page-size-select">
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

// Use computed wrapper for pageSize to allow v-model binding (props are read-only)
const localPageSize = computed({
    get: () => props.pageSize,
    set: (val) => emit('update:pageSize', val)
});

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
        emit('change', { page, pageSize: localPageSize.value });
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

function onPageSizeChange(e: Event) {
    const size = Number((e.target as HTMLSelectElement).value);
    emit('update:pageSize', size);
    emit('change', { page: 1, pageSize: size });
}

watch(() => props.currentPage, (val) => {
    jumpPage.value = val;
});

watch(() => props.pageSize, () => {
    if (localPageSize.value !== props.pageSize) {
        localPageSize.value = props.pageSize;
    }
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
