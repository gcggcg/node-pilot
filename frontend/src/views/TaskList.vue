<template>
    <div class="page">
        <div class="header">
            <h1>任务管理</h1>
            <div class="header-actions">
                <button 
                    v-if="selectedIds.length > 0" 
                    @click="deleteSelected" 
                    class="btn btn-danger"
                >
                    删除已选 ({{ selectedIds.length }})
                </button>
                <button 
                    @click="toggleAutoRefresh" 
                    :class="['btn', 'btn-small', autoRefresh ? 'btn-primary' : 'btn-secondary']"
                >
                    {{ autoRefresh ? '⏸ 暂停刷新' : '▶ 启用刷新' }}
                </button>
                <button @click="router.push('/tasks/new')" class="btn btn-primary">+ 新建任务</button>
            </div>
        </div>

        <div v-if="store.loading" class="loading">加载中...</div>
        <div v-else-if="store.error" class="error">{{ store.error }}</div>

        <table v-else class="table">
            <thead>
                <tr>
                    <th class="checkbox-col">
                        <label class="checkbox-wrapper">
                            <input type="checkbox" v-model="selectAll" @change="toggleSelectAll" />
                            <span class="checkmark"></span>
                        </label>
                    </th>
                    <th>名称</th>
                    <th>脚本</th>
                    <th>状态</th>
                    <th>创建时间</th>
                    <th>操作</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="task in store.tasks" :key="task.id">
                    <td class="checkbox-col">
                        <label class="checkbox-wrapper">
                            <input type="checkbox" :value="task.id" v-model="selectedIds" />
                            <span class="checkmark"></span>
                        </label>
                    </td>
                    <td>{{ task.name }}</td>
                    <td>{{ formatScriptInfo(task) }}</td>
                    <td>
                        <span :class="['status', task.status]">{{ statusText(task.status) }}</span>
                    </td>
                    <td>{{ formatDate(task.created_at) }}</td>
                    <td class="actions">
                        <!-- 手动执行按钮（仅 pending 状态显示） -->
                        <button
                            v-if="task.status === 'pending'"
                            @click="executeTask(task.id)"
                            class="btn btn-small btn-primary"
                        >执行</button>
                        
                        <!-- 编辑按钮（仅 pending 状态显示） -->
                        <router-link 
                            v-if="task.status === 'pending'"
                            :to="`/tasks/${task.id}/edit`" 
                            class="btn btn-small"
                        >编辑</router-link>
                        
                        <!-- 输出按钮 -->
                        <router-link :to="`/tasks/${task.id}/output`" class="btn btn-small">输出</router-link>
                        
                        <!-- 取消按钮（仅 running 状态显示） -->
                        <button
                            v-if="task.status === 'running'"
                            @click="cancelTask(task.id)"
                            class="btn btn-small btn-danger"
                        >取消</button>
                    </td>
                </tr>
                <tr v-if="store.tasks.length === 0">
                    <td colspan="6" class="empty">暂无任务，点击"+ 新建任务"创建</td>
                </tr>
            </tbody>
        </table>

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
import { ref, computed, onMounted, watch, onUnmounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useTaskStore } from '@/stores/task';
import Pagination from '@/components/Pagination.vue';

const router = useRouter();
const route = useRoute();
const store = useTaskStore();

const selectedIds = ref<number[]>([]);
const autoRefresh = ref(false);
let autoRefreshTimer: ReturnType<typeof setInterval> | null = null;

const pagination = ref({
    page: Number(route.query.page) || 1,
    pageSize: Number(route.query.pageSize) || 10,
    total: 0
});

const selectAll = computed({
    get: () => store.tasks.length > 0 && selectedIds.value.length === store.tasks.length,
    set: (val: boolean) => {
        selectedIds.value = val ? store.tasks.map(t => t.id) : [];
    }
});

async function loadData() {
    await store.fetchTasks(pagination.value.page, pagination.value.pageSize);
    pagination.value.total = store.pagination.total;
}

function handlePageChange(payload: { page: number; pageSize: number }) {
    pagination.value.page = payload.page;
    pagination.value.pageSize = payload.pageSize;
    router.replace({ 
        query: { page: payload.page, pageSize: payload.pageSize } 
    });
    store.fetchTasks(payload.page, payload.pageSize);
}

onMounted(() => {
    loadData();
    startAutoRefresh();
});

onUnmounted(() => {
    stopAutoRefresh();
});

function startAutoRefresh() {
    stopAutoRefresh();
    if (autoRefresh.value) {
        autoRefreshTimer = setInterval(() => {
            loadData();
        }, 3000);
    }
}

function stopAutoRefresh() {
    if (autoRefreshTimer) {
        clearInterval(autoRefreshTimer);
        autoRefreshTimer = null;
    }
}

function toggleAutoRefresh() {
    autoRefresh.value = !autoRefresh.value;
    if (autoRefresh.value) {
        startAutoRefresh();
    } else {
        stopAutoRefresh();
    }
}

watch(() => route.query, (query) => {
    pagination.value.page = Number(query.page) || 1;
    pagination.value.pageSize = Number(query.pageSize) || 10;
    loadData();
});

function toggleSelectAll() {
    if (selectAll.value) {
        selectedIds.value = store.tasks.map(t => t.id);
    } else {
        selectedIds.value = [];
    }
}

async function cancelTask(id: number) {
    if (confirm('确定要取消此任务吗？')) {
        await store.cancelTask(id);
        pagination.value.total = store.pagination.total;
    }
}

async function executeTask(id: number) {
    if (confirm('确定要执行此任务吗？')) {
        try {
            await store.executeTask(id);
        } catch (e: any) {
            alert(e.message || '执行任务失败');
        }
    }
}

async function deleteSelected() {
    if (confirm(`确定要删除已选的 ${selectedIds.value.length} 个任务吗？`)) {
        await store.deleteTasks(selectedIds.value);
        selectedIds.value = [];
        pagination.value.total = store.pagination.total;
    }
}

function formatDate(dateStr: string) {
    return new Date(dateStr).toLocaleString();
}

function statusText(status: string): string {
    const statusMap: Record<string, string> = {
        'pending': '等待中',
        'running': '运行中',
        'completed': '已完成',
        'cancelled': '已取消',
        'failed': '失败'
    };
    return statusMap[status] || status;
}

function formatScriptInfo(task: any): string {
    if (task.script_ids) {
        const ids = task.script_ids.split(',').map((id: string) => (id as string).trim()).filter((id: string) => id);
        if (ids.length > 0) {
            return `脚本: ${ids.join(', ')}`;
        }
    }
    if (task.script_id && task.script_id > 0) {
        return `ID: ${task.script_id}`;
    }
    return '-';
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

.checkbox-col {
    width: 48px;
    text-align: center;
}

.checkbox-wrapper {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    position: relative;
}

.checkbox-wrapper input {
    position: absolute;
    opacity: 0;
    cursor: pointer;
    height: 0;
    width: 0;
}

.checkmark {
    height: 20px;
    width: 20px;
    background-color: #fff;
    border: 2px solid #ddd;
    border-radius: 4px;
    transition: all 0.2s ease;
}

.checkbox-wrapper:hover .checkmark {
    border-color: #667eea;
}

.checkbox-wrapper input:checked ~ .checkmark {
    background-color: #667eea;
    border-color: #667eea;
}

.checkmark:after {
    content: "";
    position: absolute;
    display: none;
    left: 7px;
    top: 3px;
    width: 5px;
    height: 10px;
    border: solid white;
    border-width: 0 2px 2px 0;
    transform: rotate(45deg);
}

.checkbox-wrapper input:checked ~ .checkmark:after {
    display: block;
}

.table {
    width: 100%;
    border-collapse: collapse;
    background: white;
    border-radius: 8px;
    overflow: hidden;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.table th,
.table td {
    padding: 12px 16px;
    text-align: left;
    border-bottom: 1px solid #eee;
}

.table th {
    background: #f8f9fa;
    font-weight: 600;
    color: #555;
}

.status {
    padding: 4px 8px;
    border-radius: 4px;
    font-size: 12px;
    font-weight: 500;
}

.status.pending { background: #e9ecef; color: #495057; }
.status.running { background: #cce5ff; color: #004085; }
.status.completed { background: #d4edda; color: #155724; }
.status.cancelled { background: #fff3cd; color: #856404; }
.status.failed { background: #f8d7da; color: #721c24; }

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

.btn-secondary {
    background: #e9ecef;
    color: #495057;
}

.loading, .error, .empty {
    text-align: center;
    padding: 48px;
    color: #666;
}

.dialog-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 200;
}

.dialog {
    background: white;
    padding: 24px;
    border-radius: 8px;
    width: 500px;
    max-height: 80vh;
    overflow-y: auto;
}

.dialog h2 {
    margin: 0 0 16px;
}

.form-group {
    margin-bottom: 16px;
}

.form-group label {
    display: block;
    margin-bottom: 6px;
    font-weight: 500;
    color: #555;
}

.form-group input,
.form-group select {
    width: 100%;
    padding: 10px 12px;
    border: 1px solid #ddd;
    border-radius: 6px;
    font-size: 14px;
}

.server-list {
    max-height: 240px;
    overflow-y: auto;
    border: 1px solid #e0e0e0;
    border-radius: 8px;
    background: #fafafa;
}

.server-item {
    display: flex;
    align-items: center;
    padding: 12px 16px;
    border-bottom: 1px solid #eee;
    cursor: pointer;
    transition: background 0.2s;
}

.server-item:last-child {
    border-bottom: none;
}

.server-item:hover {
    background: #f0f0f0;
}

.server-item .checkbox-wrapper {
    margin-right: 12px;
}

.server-info {
    display: flex;
    flex-direction: column;
}

.server-name {
    font-weight: 500;
    color: #333;
}

.server-host {
    font-size: 12px;
    color: #888;
}

.no-servers {
    padding: 24px;
    text-align: center;
    color: #888;
}

.dialog-actions {
    display: flex;
    gap: 12px;
    margin-top: 24px;
}

.checkbox-wrapper {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    position: relative;
}

.checkbox-wrapper input {
    position: absolute;
    opacity: 0;
    cursor: pointer;
    height: 0;
    width: 0;
}

.checkmark {
    height: 20px;
    width: 20px;
    background-color: #fff;
    border: 2px solid #ddd;
    border-radius: 4px;
    transition: all 0.2s ease;
}

.checkbox-wrapper:hover .checkmark {
    border-color: #667eea;
}

.checkbox-wrapper input:checked ~ .checkmark {
    background-color: #667eea;
    border-color: #667eea;
}

.checkmark:after {
    content: "";
    position: absolute;
    display: none;
    left: 7px;
    top: 3px;
    width: 5px;
    height: 10px;
    border: solid white;
    border-width: 0 2px 2px 0;
    transform: rotate(45deg);
}

.checkbox-wrapper input:checked ~ .checkmark:after {
    display: block;
}
</style>
