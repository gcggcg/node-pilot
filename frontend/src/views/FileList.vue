<template>
    <div class="page">
        <div class="header">
            <h1>文件管理</h1>
            <div class="header-actions">
                <button 
                    v-if="selectedIds.length > 0" 
                    @click="deleteSelected" 
                    class="btn btn-danger"
                >
                    删除已选 ({{ selectedIds.length }})
                </button>
                <router-link to="/files/new" class="btn btn-primary">+ 新增上传配置</router-link>
            </div>
        </div>

        <!-- 筛选区域 -->
        <div class="filters">
            <select v-model="filters.status" class="filter-select">
                <option value="">全部状态</option>
                <option value="pending">待执行</option>
                <option value="success">成功</option>
                <option value="failed">失败</option>
            </select>
            <input 
                v-model="filters.keyword" 
                type="text" 
                placeholder="搜索配置名称" 
                class="filter-input"
            />
            <button @click="handleSearch" class="btn btn-small">查询</button>
            <button @click="handleReset" class="btn btn-small btn-secondary">重置</button>
        </div>

        <div v-if="store.loading" class="loading">加载中...</div>

        <table v-if="!store.loading" class="table">
            <thead>
                <tr>
                    <th class="checkbox-col">
                        <label class="checkbox-wrapper">
                            <input type="checkbox" v-model="selectAll" @change="toggleSelectAll" />
                            <span class="checkmark"></span>
                        </label>
                    </th>
                    <th>序号</th>
                    <th>配置名称</th>
                    <th>本地路径</th>
                    <th>目标服务器</th>
                    <th>远程路径</th>
                    <th>状态</th>
                    <th>创建时间</th>
                    <th>操作</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="(upload, index) in store.uploads" :key="upload.id">
                    <td class="checkbox-col">
                        <label class="checkbox-wrapper">
                            <input type="checkbox" :value="upload.id" v-model="selectedIds" />
                            <span class="checkmark"></span>
                        </label>
                    </td>
                    <td>{{ (pagination.page - 1) * pagination.pageSize + index + 1 }}</td>
                    <td>{{ upload.name }}</td>
                    <td class="path-cell">{{ upload.local_path }}</td>
                    <td>
                        <span v-if="upload.servers && upload.servers.length > 0" class="server-tags">
                            <span v-for="server in upload.servers.slice(0, 3)" :key="server.id" class="server-tag">
                                {{ server.name }}
                            </span>
                            <span v-if="upload.servers.length > 3" class="server-tag more">
                                +{{ upload.servers.length - 3 }}
                            </span>
                        </span>
                        <span v-else class="no-servers">未选择服务器</span>
                    </td>
                    <td class="path-cell">{{ upload.remote_path }}</td>
                    <td>
                        <span :class="['status', upload.status]">
                            {{ statusText(upload.status) }}
                        </span>
                    </td>
                    <td>{{ formatDate(upload.created_at) }}</td>
                    <td class="actions">
                        <router-link :to="`/files/${upload.id}/edit`" class="btn btn-small">编辑</router-link>
                        <button @click="handleExecute(upload.id)" class="btn btn-small btn-primary">执行</button>
                        <button @click="handleViewResults(upload.id)" class="btn btn-small">结果</button>
                        <button @click="deleteUpload(upload.id)" class="btn btn-small btn-danger">删除</button>
                    </td>
                </tr>
                <tr v-if="store.uploads.length === 0">
                    <td colspan="9" class="empty">暂无上传配置，点击"+ 新增上传配置"创建</td>
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

        <!-- 查看结果弹窗 -->
        <div v-if="showResultsModal" class="modal-overlay" @click.self="closeResultsModal">
            <div class="modal">
                <div class="modal-header">
                    <h2>{{ store.currentUpload?.name }} - 执行结果</h2>
                    <button @click="closeResultsModal" class="modal-close">&times;</button>
                </div>
                <div class="modal-body">
                    <table class="results-table">
                        <thead>
                            <tr>
                                <th>服务器</th>
                                <th>状态</th>
                                <th>远程路径</th>
                                <th>错误信息</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="result in store.results" :key="result.id">
                                <td>{{ result.server_name || `服务器 ${result.server_id}` }}</td>
                                <td>
                                    <span :class="['status', result.status]">
                                        {{ statusText(result.status) }}
                                    </span>
                                </td>
                                <td class="path-cell">{{ result.remote_full_path || '-' }}</td>
                                <td class="error-cell">{{ result.error_message || '-' }}</td>
                            </tr>
                            <tr v-if="store.results.length === 0">
                                <td colspan="4" class="empty">暂无执行结果</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useFileUploadStore } from '@/stores/fileupload';
import Pagination from '@/components/Pagination.vue';

const router = useRouter();
const route = useRoute();
const store = useFileUploadStore();
const selectedIds = ref<number[]>([]);
const showResultsModal = ref(false);

const filters = ref({
    status: '',
    keyword: '',
    startTime: '',
    endTime: ''
});

const pagination = ref({
    page: Number(route.query.page) || 1,
    pageSize: Number(route.query.pageSize) || 10,
    total: 0
});

const selectAll = computed({
    get: () => store.uploads.length > 0 && selectedIds.value.length === store.uploads.length,
    set: (val: boolean) => {
        selectedIds.value = val ? store.uploads.map(u => u.id) : [];
    }
});

async function loadData() {
    await store.fetchUploads(pagination.value.page, pagination.value.pageSize, filters.value);
    pagination.value.total = store.pagination.total;
}

async function handleSearch() {
    pagination.value.page = 1;
    await loadData();
}

async function handleReset() {
    filters.value = {
        status: '',
        keyword: '',
        startTime: '',
        endTime: ''
    };
    pagination.value.page = 1;
    await loadData();
}

function handlePageChange(payload: { page: number; pageSize: number }) {
    pagination.value.page = payload.page;
    pagination.value.pageSize = payload.pageSize;
    router.replace({ 
        query: { page: payload.page, pageSize: payload.pageSize } 
    });
    store.fetchUploads(payload.page, payload.pageSize, filters.value);
}

onMounted(() => {
    loadData();
});

watch(() => route.query, (query) => {
    pagination.value.page = Number(query.page) || 1;
    pagination.value.pageSize = Number(route.query.pageSize) || 10;
    loadData();
});

function toggleSelectAll() {
    if (selectAll.value) {
        selectedIds.value = store.uploads.map(u => u.id);
    } else {
        selectedIds.value = [];
    }
}

async function handleExecute(id: number) {
    if (!confirm('确定要执行上传吗？')) return;
    try {
        await store.executeUpload(id);
        await loadData();
    } catch (e: any) {
        alert('执行失败: ' + (e.message || '未知错误'));
    }
}

async function handleViewResults(id: number) {
    await store.fetchResults(id);
    showResultsModal.value = true;
}

function closeResultsModal() {
    showResultsModal.value = false;
    store.resetResults();
}

function statusText(status: string): string {
    const statusMap: Record<string, string> = {
        'pending': '待执行',
        'success': '成功',
        'failed': '失败',
        'running': '执行中'
    };
    return statusMap[status] || status;
}

function formatDate(dateStr: string): string {
    if (!dateStr) return '-';
    const date = new Date(dateStr);
    return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
    });
}

async function deleteUpload(id: number) {
    if (confirm('确定要删除此上传配置吗？')) {
        await store.deleteUploads([id]);
        selectedIds.value = selectedIds.value.filter(i => i !== id);
        pagination.value.total = store.pagination.total;
    }
}

async function deleteSelected() {
    if (confirm(`确定要删除已选的 ${selectedIds.value.length} 个上传配置吗？`)) {
        await store.deleteUploads(selectedIds.value);
        selectedIds.value = [];
        pagination.value.total = store.pagination.total;
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

.filters {
    display: flex;
    gap: 12px;
    margin-bottom: 20px;
    flex-wrap: wrap;
}

.filter-select,
.filter-input {
    padding: 8px 12px;
    border: 1px solid #ddd;
    border-radius: 6px;
    font-size: 14px;
}

.filter-input {
    min-width: 200px;
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

.table tr:last-child td {
    border-bottom: none;
}

.table tbody tr:hover {
    background-color: #f8f9fa;
}

.path-cell {
    font-family: monospace;
    font-size: 12px;
    max-width: 200px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.server-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
}

.server-tag {
    display: inline-block;
    padding: 2px 8px;
    background: #e9ecef;
    border-radius: 10px;
    font-size: 12px;
    color: #495057;
}

.server-tag.more {
    background: #dee2e6;
    color: #6c757d;
}

.no-servers {
    color: #999;
    font-size: 12px;
}

.actions {
    display: flex;
    gap: 6px;
}

.btn {
    padding: 8px 16px;
    border: none;
    border-radius: 6px;
    cursor: pointer;
    font-size: 14px;
    transition: all 0.2s;
}

.btn-primary {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
}

.btn-primary:hover {
    opacity: 0.9;
}

.btn-secondary {
    background: #e9ecef;
    color: #495057;
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

.btn-danger:hover {
    background: #c82333;
}

.loading, .error, .empty {
    text-align: center;
    padding: 48px;
    color: #666;
}

.error {
    color: #dc3545;
}

.status {
    display: inline-block;
    padding: 4px 10px;
    border-radius: 12px;
    font-size: 12px;
    font-weight: 500;
}

.status.pending {
    background: #fff3cd;
    color: #856404;
}

.status.success {
    background: #d4edda;
    color: #155724;
}

.status.failed {
    background: #f8d7da;
    color: #721c24;
}

.status.running {
    background: #cce5ff;
    color: #004085;
}

/* Modal styles */
.modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
}

.modal {
    background: white;
    border-radius: 12px;
    width: 90%;
    max-width: 800px;
    max-height: 80vh;
    overflow: hidden;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 20px;
    border-bottom: 1px solid #eee;
}

.modal-header h2 {
    font-size: 1.2rem;
    margin: 0;
}

.modal-close {
    background: none;
    border: none;
    font-size: 24px;
    cursor: pointer;
    color: #999;
    line-height: 1;
}

.modal-close:hover {
    color: #333;
}

.modal-body {
    padding: 20px;
    overflow-y: auto;
    max-height: calc(80vh - 70px);
}

.results-table {
    width: 100%;
    border-collapse: collapse;
}

.results-table th,
.results-table td {
    padding: 10px 12px;
    text-align: left;
    border-bottom: 1px solid #eee;
}

.results-table th {
    background: #f8f9fa;
    font-weight: 600;
    color: #555;
}

.error-cell {
    color: #dc3545;
    font-size: 13px;
}
</style>
