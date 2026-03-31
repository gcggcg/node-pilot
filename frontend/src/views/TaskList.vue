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
                <button @click="showCreateDialog = true" class="btn btn-primary">+ 运行任务</button>
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
                    <th>脚本ID</th>
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
                    <td>{{ task.script_id }}</td>
                    <td>
                        <span :class="['status', task.status]">{{ statusText(task.status) }}</span>
                    </td>
                    <td>{{ formatDate(task.created_at) }}</td>
                    <td class="actions">
                        <router-link :to="`/tasks/${task.id}/output`" class="btn btn-small">输出</router-link>
                        <button
                            v-if="task.status === 'running'"
                            @click="cancelTask(task.id)"
                            class="btn btn-small btn-danger"
                        >取消</button>
                    </td>
                </tr>
                <tr v-if="store.tasks.length === 0">
                    <td colspan="6" class="empty">暂无任务，点击"+ 运行任务"创建</td>
                </tr>
            </tbody>
        </table>

        <div v-if="showCreateDialog" class="dialog-overlay" @click.self="showCreateDialog = false">
            <div class="dialog">
                <h2>运行任务</h2>
                <div class="form-group">
                    <label>任务名称</label>
                    <input v-model="newTask.name" type="text" placeholder="例如: 批量部署-1" />
                </div>
                <div class="form-group">
                    <label>脚本</label>
                    <select v-model="newTask.script_id">
                        <option value="">选择脚本...</option>
                        <option v-for="s in scriptStore.scripts" :key="s.id" :value="s.id">{{ s.name }}</option>
                    </select>
                </div>
                <div class="form-group">
                    <label>服务器</label>
                    <div class="server-list">
                        <label v-for="server in serverStore.servers" :key="server.id" class="server-item">
                            <label class="checkbox-wrapper">
                                <input type="checkbox" :value="server.id" v-model="newTask.server_ids" />
                                <span class="checkmark"></span>
                            </label>
                            <span class="server-info">
                                <span class="server-name">{{ server.name }}</span>
                                <span class="server-host">{{ server.host }}</span>
                            </span>
                        </label>
                        <div v-if="serverStore.servers.length === 0" class="no-servers">
                            暂无服务器，请先添加服务器
                        </div>
                    </div>
                </div>
                <div class="dialog-actions">
                    <button @click="createTask" class="btn btn-primary" :disabled="loading">运行</button>
                    <button @click="showCreateDialog = false" class="btn btn-secondary">取消</button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useTaskStore } from '@/stores/task';
import { useScriptStore } from '@/stores/script';
import { useServerStore } from '@/stores/server';

const store = useTaskStore();
const scriptStore = useScriptStore();
const serverStore = useServerStore();

const showCreateDialog = ref(false);
const loading = ref(false);
const selectedIds = ref<number[]>([]);
const newTask = ref({
    name: '',
    script_id: '' as number | '',
    server_ids: [] as number[]
});

const selectAll = computed({
    get: () => store.tasks.length > 0 && selectedIds.value.length === store.tasks.length,
    set: (val: boolean) => {
        selectedIds.value = val ? store.tasks.map(t => t.id) : [];
    }
});

onMounted(() => {
    store.fetchTasks();
    scriptStore.fetchScripts();
    serverStore.fetchServers();
});

function toggleSelectAll() {
    if (selectAll.value) {
        selectedIds.value = store.tasks.map(t => t.id);
    } else {
        selectedIds.value = [];
    }
}

async function createTask() {
    if (!newTask.value.name || !newTask.value.script_id || newTask.value.server_ids.length === 0) {
        alert('请填写所有字段');
        return;
    }
    loading.value = true;
    try {
        await store.createTask({
            name: newTask.value.name,
            script_id: Number(newTask.value.script_id),
            server_ids: newTask.value.server_ids
        });
        showCreateDialog.value = false;
        newTask.value = { name: '', script_id: '', server_ids: [] };
    } catch (e: any) {
        alert(e.message || '创建任务失败');
    } finally {
        loading.value = false;
    }
}

async function cancelTask(id: number) {
    if (confirm('确定要取消此任务吗？')) {
        await store.cancelTask(id);
    }
}

async function deleteSelected() {
    if (confirm(`确定要删除已选的 ${selectedIds.value.length} 个任务吗？`)) {
        await store.deleteTasks(selectedIds.value);
        selectedIds.value = [];
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
