<template>
    <div class="page">
        <div class="header">
            <h1>脚本管理</h1>
            <div class="header-actions">
                <button 
                    v-if="selectedIds.length > 0" 
                    @click="deleteSelected" 
                    class="btn btn-danger"
                >
                    删除已选 ({{ selectedIds.length }})
                </button>
                <router-link to="/scripts/new" class="btn btn-primary">+ 创建脚本</router-link>
            </div>
        </div>

        <div v-if="store.loading" class="loading">加载中...</div>
        <div v-else-if="store.error" class="error">{{ store.error }}</div>

        <div v-else class="grid">
            <div v-for="script in store.scripts" :key="script.id" class="card" :class="{ selected: selectedIds.includes(script.id) }">
                <div class="card-checkbox">
                    <input type="checkbox" :value="script.id" v-model="selectedIds" />
                </div>
                <h3>{{ script.name }}</h3>
                <p class="description">{{ script.description || '暂无描述' }}</p>
                <p class="path">{{ script.target_path }}</p>
                <div class="actions">
                    <router-link :to="`/scripts/${script.id}/edit`" class="btn btn-small">编辑</router-link>
                    <button @click="deleteScript(script.id)" class="btn btn-small btn-danger">删除</button>
                </div>
            </div>
            <div v-if="store.scripts.length === 0" class="empty">
                暂无脚本，点击"+ 创建脚本"创建
            </div>
        </div>

        <Pagination 
            v-if="pagination.total > 0"
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :total="pagination.total"
            @change="handlePageChange"
        />
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useScriptStore } from '@/stores/script';
import Pagination from '@/components/Pagination.vue';

const router = useRouter();
const route = useRoute();
const store = useScriptStore();
const selectedIds = ref<number[]>([]);

const pagination = ref({
    page: Number(route.query.page) || 1,
    pageSize: Number(route.query.pageSize) || 10,
    total: 0
});

function loadData() {
    store.fetchScripts(pagination.value.page, pagination.value.pageSize);
    pagination.value.total = store.pagination.total;
}

function handlePageChange(payload: { page: number; pageSize: number }) {
    pagination.value.page = payload.page;
    pagination.value.pageSize = payload.pageSize;
    router.replace({ 
        query: { page: payload.page, pageSize: payload.pageSize } 
    });
    store.fetchScripts(payload.page, payload.pageSize);
}

onMounted(() => {
    loadData();
});

watch(() => route.query, (query) => {
    pagination.value.page = Number(query.page) || 1;
    pagination.value.pageSize = Number(query.pageSize) || 10;
    loadData();
});

async function deleteScript(id: number) {
    if (confirm('确定要删除此脚本吗？')) {
        await store.deleteScript(id);
        selectedIds.value = selectedIds.value.filter(i => i !== id);
    }
}

async function deleteSelected() {
    if (confirm(`确定要删除已选的 ${selectedIds.value.length} 个脚本吗？`)) {
        await store.deleteScripts(selectedIds.value);
        selectedIds.value = [];
    }
}
</script>

<style scoped>
.page {
    padding: 24px;
}

.header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
}

.header-actions {
    display: flex;
    gap: 12px;
}

h1 {
    font-size: 1.5rem;
    color: #333;
}

.grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 16px;
}

.card {
    background: white;
    padding: 20px;
    border-radius: 8px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    position: relative;
    transition: all 0.2s;
}

.card.selected {
    box-shadow: 0 0 0 2px #dc3545;
}

.card-checkbox {
    position: absolute;
    top: 12px;
    right: 12px;
}

.card-checkbox input {
    cursor: pointer;
    width: 18px;
    height: 18px;
}

.card h3 {
    margin: 0 0 8px;
    color: #333;
}

.description {
    color: #666;
    font-size: 14px;
    margin: 0 0 8px;
}

.path {
    font-family: monospace;
    font-size: 13px;
    color: #888;
    margin: 0 0 16px;
}

.actions {
    display: flex;
    gap: 8px;
}

.btn {
    padding: 8px 16px;
    border: none;
    border-radius: 6px;
    cursor: pointer;
    font-size: 14px;
}

.btn-primary {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
}

.btn-small {
    padding: 4px 12px;
    font-size: 13px;
    background: #e9ecef;
    color: #495057;
}

.btn-danger {
    background: #dc3545;
    color: white;
}

.loading, .error, .empty {
    text-align: center;
    padding: 48px;
    color: #666;
}

.error {
    color: #dc3545;
}
</style>
