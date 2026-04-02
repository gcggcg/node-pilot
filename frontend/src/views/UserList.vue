<template>
    <div class="page">
        <div class="header">
            <h1>用户管理</h1>
            <div class="header-actions">
                <button 
                    v-if="selectedIds.length > 0" 
                    @click="deleteSelected" 
                    class="btn btn-danger"
                >
                    删除已选 ({{ selectedIds.length }})
                </button>
                <button @click="showAddModal = true" class="btn btn-primary">+ 添加用户</button>
            </div>
        </div>

        <div class="search-bar">
            <input 
                v-model="keyword" 
                type="text" 
                placeholder="搜索用户名..." 
                class="search-input"
                @keyup.enter="handleSearch"
            />
            <button @click="handleSearch" class="btn btn-secondary">搜索</button>
        </div>

        <div v-if="loading" class="loading">加载中...</div>
        <div v-else-if="error" class="error">{{ error }}</div>

        <table v-else class="table">
            <thead>
                <tr>
                    <th class="checkbox-col">
                        <label class="checkbox-wrapper">
                            <input type="checkbox" v-model="selectAll" @change="toggleSelectAll" />
                            <span class="checkmark"></span>
                        </label>
                    </th>
                    <th>ID</th>
                    <th>用户名</th>
                    <th>邮箱</th>
                    <th>电话</th>
                    <th>角色</th>
                    <th>创建时间</th>
                    <th>操作</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="user in users" :key="user.id">
                    <td class="checkbox-col">
                        <label class="checkbox-wrapper">
                            <input type="checkbox" :value="user.id" v-model="selectedIds" />
                            <span class="checkmark"></span>
                        </label>
                    </td>
                    <td>{{ user.id }}</td>
                    <td>{{ user.username }}</td>
                    <td>{{ user.email || '-' }}</td>
                    <td>{{ user.phone || '-' }}</td>
                    <td>
                        <span :class="['role-badge', user.role === 'ROLE_ADMIN' ? 'admin' : 'user']">
                            {{ user.role === 'ROLE_ADMIN' ? '管理员' : '用户' }}
                        </span>
                    </td>
                    <td>{{ formatDate(user.created_at) }}</td>
                    <td class="actions">
                        <button 
                            v-if="user.username !== 'root'" 
                            @click="deleteUser(user.id)" 
                            class="btn btn-small btn-danger"
                        >
                            删除
                        </button>
                    </td>
                </tr>
                <tr v-if="users.length === 0">
                    <td colspan="8" class="empty">暂无用户</td>
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

        <!-- 添加用户弹窗 -->
        <div v-if="showAddModal" class="modal-overlay" @click.self="showAddModal = false">
            <div class="modal">
                <h3>添加用户</h3>
                <form @submit.prevent="handleAdd">
                    <div class="form-group">
                        <label>用户名</label>
                        <input v-model="newUser.username" required />
                    </div>
                    <div class="form-group">
                        <label>密码</label>
                        <input v-model="newUser.password" type="password" required minlength="6" />
                    </div>
                    <div class="form-group">
                        <label>邮箱</label>
                        <input v-model="newUser.email" type="email" />
                    </div>
                    <div class="form-group">
                        <label>电话</label>
                        <input v-model="newUser.phone" />
                    </div>
                    <div class="form-group">
                        <label>角色</label>
                        <select v-model="newUser.role" required>
                            <option value="ROLE_USER">普通用户</option>
                            <option value="ROLE_ADMIN">管理员</option>
                        </select>
                    </div>
                    <div class="modal-actions">
                        <button type="button" @click="showAddModal = false" class="btn btn-secondary">取消</button>
                        <button type="submit" class="btn btn-primary">添加</button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { userApi } from '@/api';
import type { User } from '@/types';
import Pagination from '@/components/Pagination.vue';

const router = useRouter();
const route = useRoute();

const users = ref<User[]>([]);
const selectedIds = ref<number[]>([]);
const keyword = ref('');
const showAddModal = ref(false);
const loading = ref(false);
const error = ref('');

const newUser = ref({
    username: '',
    password: '',
    email: '',
    phone: '',
    role: 'ROLE_USER' as 'ROLE_USER' | 'ROLE_ADMIN'
});

const pagination = ref({
    page: Number(route.query.page) || 1,
    pageSize: Number(route.query.pageSize) || 10,
    total: 0
});

const selectAll = computed({
    get: () => users.value.length > 0 && selectedIds.value.length === users.value.length,
    set: (val: boolean) => {
        selectedIds.value = val ? users.value.map(u => u.id) : [];
    }
});

async function loadUsers() {
    loading.value = true;
    error.value = '';
    try {
        const res = await userApi.list({
            page: pagination.value.page,
            pageSize: pagination.value.pageSize,
            keyword: keyword.value || undefined
        });
        users.value = res.data;
        pagination.value.total = res.total;
    } catch (e: any) {
        error.value = e.response?.data?.error || '加载失败';
    } finally {
        loading.value = false;
    }
}

function handleSearch() {
    pagination.value.page = 1;
    router.replace({ 
        query: { page: 1, pageSize: pagination.value.pageSize, keyword: keyword.value || undefined } 
    });
    loadUsers();
}

function handlePageChange(payload: { page: number; pageSize: number }) {
    pagination.value.page = payload.page;
    pagination.value.pageSize = payload.pageSize;
    router.replace({ 
        query: { page: payload.page, pageSize: payload.pageSize, keyword: keyword.value || undefined } 
    });
    loadUsers();
}

function toggleSelectAll() {
    if (selectAll.value) {
        selectedIds.value = users.value.map(u => u.id);
    } else {
        selectedIds.value = [];
    }
}

async function handleAdd() {
    try {
        await userApi.create(newUser.value);
        showAddModal.value = false;
        loadUsers();
        // 重置表单
        newUser.value = { username: '', password: '', email: '', phone: '', role: 'ROLE_USER' };
    } catch (e: any) {
        alert(e.response?.data?.error || '添加失败');
    }
}

async function deleteUser(id: number) {
    if (!confirm('确认删除此用户？')) return;
    try {
        await userApi.delete(id);
        loadUsers();
    } catch (e: any) {
        alert(e.response?.data?.error || '删除失败');
    }
}

async function deleteSelected() {
    if (!confirm(`确认删除选中的 ${selectedIds.value.length} 个用户？`)) return;
    try {
        await userApi.deleteMany(selectedIds.value);
        selectedIds.value = [];
        loadUsers();
    } catch (e: any) {
        alert(e.response?.data?.error || '删除失败');
    }
}

function formatDate(dateStr: string) {
    return new Date(dateStr).toLocaleDateString('zh-CN');
}

onMounted(loadUsers);
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

.search-bar {
    display: flex;
    gap: 12px;
    margin-bottom: 20px;
}

.search-input {
    flex: 1;
    max-width: 300px;
    padding: 8px 12px;
    border: 1px solid #ddd;
    border-radius: 6px;
    font-size: 14px;
}

.search-input:focus {
    outline: none;
    border-color: #667eea;
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

.btn-secondary {
    background: #e9ecef;
    color: #495057;
}

.btn-secondary:hover {
    background: #dee2e6;
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

.role-badge {
    display: inline-block;
    padding: 4px 10px;
    border-radius: 12px;
    font-size: 12px;
    font-weight: 500;
}

.role-badge.admin {
    background: #d4edda;
    color: #155724;
}

.role-badge.user {
    background: #e9ecef;
    color: #6c757d;
}

/* Modal Styles */
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
    padding: 24px;
    width: 90%;
    max-width: 480px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.modal h3 {
    margin: 0 0 20px 0;
    font-size: 1.25rem;
    color: #333;
}

.form-group {
    margin-bottom: 16px;
}

.form-group label {
    display: block;
    margin-bottom: 6px;
    font-size: 14px;
    color: #555;
    font-weight: 500;
}

.form-group input,
.form-group select {
    width: 100%;
    padding: 10px 12px;
    border: 1px solid #ddd;
    border-radius: 6px;
    font-size: 14px;
    box-sizing: border-box;
}

.form-group input:focus,
.form-group select:focus {
    outline: none;
    border-color: #667eea;
}

.modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    margin-top: 24px;
}
</style>
