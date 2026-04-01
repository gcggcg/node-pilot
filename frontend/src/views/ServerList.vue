<template>
    <div class="page">
        <div class="header">
            <h1>服务器管理</h1>
            <div class="header-actions">
                <button 
                    v-if="selectedIds.length > 0" 
                    @click="deleteSelected" 
                    class="btn btn-danger"
                >
                    删除已选 ({{ selectedIds.length }})
                </button>
                <router-link to="/servers/new" class="btn btn-primary">+ 添加服务器</router-link>
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
                    <th>IP地址</th>
                    <th>端口</th>
                    <th>用户名</th>
                    <th>状态</th>
                    <th>操作</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="server in store.servers" :key="server.id">
                    <td class="checkbox-col">
                        <label class="checkbox-wrapper">
                            <input type="checkbox" :value="server.id" v-model="selectedIds" />
                            <span class="checkmark"></span>
                        </label>
                    </td>
                    <td>{{ server.name }}</td>
                    <td>{{ server.host }}</td>
                    <td>{{ server.port }}</td>
                    <td>{{ server.username }}</td>
                    <td>
                        <span :class="['status', server.connection_status]">
                            {{ statusText(server.connection_status) }}
                        </span>
                    </td>
                    <td class="actions">
                        <button @click="testConnection(server.id)" class="btn btn-small">测试</button>
                        <router-link :to="`/servers/${server.id}/edit`" class="btn btn-small">编辑</router-link>
                        <button @click="deleteServer(server.id)" class="btn btn-small btn-danger">删除</button>
                    </td>
                </tr>
                <tr v-if="store.servers.length === 0">
                    <td colspan="7" class="empty">暂无服务器，点击"+ 添加服务器"创建</td>
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
import { ref, computed, onMounted, watch } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useServerStore } from '@/stores/server';
import Pagination from '@/components/Pagination.vue';

const router = useRouter();
const route = useRoute();
const store = useServerStore();
const selectedIds = ref<number[]>([]);

const pagination = ref({
    page: Number(route.query.page) || 1,
    pageSize: Number(route.query.pageSize) || 10,
    total: 0
});

const selectAll = computed({
    get: () => store.servers.length > 0 && selectedIds.value.length === store.servers.length,
    set: (val: boolean) => {
        selectedIds.value = val ? store.servers.map(s => s.id) : [];
    }
});

async function loadData() {
    await store.fetchServers(pagination.value.page, pagination.value.pageSize);
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

onMounted(() => {
    loadData();
});

watch(() => route.query, (query) => {
    pagination.value.page = Number(query.page) || 1;
    pagination.value.pageSize = Number(query.pageSize) || 10;
    loadData();
});

function toggleSelectAll() {
    if (selectAll.value) {
        selectedIds.value = store.servers.map(s => s.id);
    } else {
        selectedIds.value = [];
    }
}

async function testConnection(id: number) {
    if (confirm('确定要测试此服务器的SSH连接吗？')) {
        try {
            await store.testConnection(id);
            await store.fetchServers(pagination.value.page, pagination.value.pageSize);
        } catch (e: any) {
            alert('连接失败: ' + e.message);
            await store.fetchServers(pagination.value.page, pagination.value.pageSize);
        }
    }
}

function statusText(status: string): string {
    const statusMap: Record<string, string> = {
        'online': '在线',
        'offline': '异常',
        'unknown': '未知'
    };
    return statusMap[status] || status;
}

async function deleteServer(id: number) {
    if (confirm('确定要删除此服务器吗？')) {
        await store.deleteServer(id);
        selectedIds.value = selectedIds.value.filter(i => i !== id);
        pagination.value.total = store.pagination.total;
    }
}

async function deleteSelected() {
    if (confirm(`确定要删除已选的 ${selectedIds.value.length} 个服务器吗？`)) {
        await store.deleteServers(selectedIds.value);
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
    transition: all 0.2s;
}

.btn-primary {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
}

.btn-primary:hover {
    opacity: 0.9;
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

.status.online {
    background: #d4edda;
    color: #155724;
}

.status.offline {
    background: #f8d7da;
    color: #721c24;
}

.status.unknown {
    background: #e9ecef;
    color: #6c757d;
}
</style>
