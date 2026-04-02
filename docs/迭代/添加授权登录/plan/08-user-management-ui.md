# 创建用户管理页面和个人中心

## 任务描述

创建管理员用户管理页面和个人信息修改页面。

## 详细说明

### 1. 创建用户管理页面

创建 `frontend/src/views/UserList.vue`：

```vue
<template>
    <div class="user-list">
        <div class="page-header">
            <h2>用户管理</h2>
            <button @click="showAddModal = true" class="btn-primary">添加用户</button>
        </div>

        <div class="search-bar">
            <input 
                v-model="keyword" 
                type="text" 
                placeholder="搜索用户名..."
                @keyup.enter="handleSearch"
            />
            <button @click="handleSearch">搜索</button>
        </div>

        <table class="data-table">
            <thead>
                <tr>
                    <th><input type="checkbox" v-model="selectAll" @change="toggleSelectAll" /></th>
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
                    <td><input type="checkbox" :value="user.id" v-model="selectedIds" /></td>
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
                    <td>
                        <button 
                            v-if="user.username !== 'root'" 
                            @click="handleDelete(user.id)" 
                            class="btn-danger"
                        >
                            删除
                        </button>
                    </td>
                </tr>
            </tbody>
        </table>

        <div class="table-footer" v-if="selectedIds.length > 0">
            <span>已选择 {{ selectedIds.length }} 项</span>
            <button @click="handleBatchDelete" class="btn-danger">批量删除</button>
        </div>

        <Pagination 
            v-if="pagination.total > 0"
            :current-page="pagination.page"
            :page-size="pagination.pageSize"
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
                        <button type="button" @click="showAddModal = false" class="btn-secondary">取消</button>
                        <button type="submit" class="btn-primary">添加</button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { userApi } from '@/api';
import type { User } from '@/types';
import Pagination from '@/components/Pagination.vue';

const users = ref<User[]>([]);
const selectedIds = ref<number[]>([]);
const selectAll = ref(false);
const keyword = ref('');
const showAddModal = ref(false);

const newUser = ref({
    username: '',
    password: '',
    email: '',
    phone: '',
    role: 'ROLE_USER'
});

const pagination = ref({
    page: 1,
    pageSize: 10,
    total: 0
});

async function loadUsers() {
    const res = await userApi.list({
        page: pagination.value.page,
        pageSize: pagination.value.pageSize,
        keyword: keyword.value
    });
    users.value = res.data;
    pagination.value.total = res.total;
}

function handleSearch() {
    pagination.value.page = 1;
    loadUsers();
}

function handlePageChange({ page, pageSize }: { page: number; pageSize: number }) {
    pagination.value.page = page;
    pagination.value.pageSize = pageSize;
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
    await userApi.create(newUser.value);
    showAddModal.value = false;
    loadUsers();
    // 重置表单
    newUser.value = { username: '', password: '', email: '', phone: '', role: 'ROLE_USER' };
}

async function handleDelete(id: number) {
    if (!confirm('确认删除此用户？')) return;
    await userApi.delete(id);
    loadUsers();
}

async function handleBatchDelete() {
    if (!confirm(`确认删除选中的 ${selectedIds.value.length} 个用户？`)) return;
    await userApi.deleteMany(selectedIds.value);
    selectedIds.value = [];
    selectAll.value = false;
    loadUsers();
}

function formatDate(dateStr: string) {
    return new Date(dateStr).toLocaleDateString('zh-CN');
}

onMounted(loadUsers);
</script>

<style scoped>
/* 参考现有列表页样式 */
</style>
```

### 2. 创建个人中心页面

创建 `frontend/src/views/Profile.vue`：

```vue
<template>
    <div class="profile-page">
        <h2>个人中心</h2>
        
        <div class="profile-section">
            <h3>基本信息</h3>
            <div class="info-grid">
                <div class="info-item">
                    <label>用户名</label>
                    <span>{{ authStore.user?.username }}</span>
                </div>
                <div class="info-item">
                    <label>角色</label>
                    <span :class="['role-badge', authStore.user?.role === 'ROLE_ADMIN' ? 'admin' : 'user']">
                        {{ authStore.user?.role === 'ROLE_ADMIN' ? '管理员' : '普通用户' }}
                    </span>
                </div>
            </div>
        </div>

        <div class="profile-section">
            <h3>修改个人信息</h3>
            <form @submit.prevent="handleUpdateProfile">
                <div class="form-group">
                    <label>邮箱</label>
                    <input v-model="profileForm.email" type="email" />
                </div>
                <div class="form-group">
                    <label>电话</label>
                    <input v-model="profileForm.phone" />
                </div>
                <button type="submit" class="btn-primary">保存</button>
            </form>
        </div>

        <div class="profile-section">
            <h3>修改密码</h3>
            <form @submit.prevent="handleChangePassword">
                <div class="form-group">
                    <label>旧密码</label>
                    <input v-model="passwordForm.old_password" type="password" required />
                </div>
                <div class="form-group">
                    <label>新密码</label>
                    <input v-model="passwordForm.new_password" type="password" required minlength="6" />
                </div>
                <button type="submit" class="btn-primary">修改密码</button>
            </form>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useAuthStore } from '@/stores/auth';

const authStore = useAuthStore();

const profileForm = ref({
    email: authStore.user?.email || '',
    phone: authStore.user?.phone || ''
});

const passwordForm = ref({
    old_password: '',
    new_password: ''
});

async function handleUpdateProfile() {
    await authStore.updateProfile(profileForm.value);
    alert('个人信息已更新');
}

async function handleChangePassword() {
    await authStore.changePassword(passwordForm.value.old_password, passwordForm.value.new_password);
    alert('密码已修改');
    passwordForm.value = { old_password: '', new_password: '' };
}
</script>
```

## 输入

- 需求文档 `05-添加授权登录.md`
- 现有的 `frontend/src/views/*.vue` 页面组件

## 输出

- 新建 `frontend/src/views/UserList.vue` - 用户管理页面
- 新建 `frontend/src/views/Profile.vue` - 个人中心页面

## 依赖

- 05-frontend-auth-api-store.md
- 06-login-page.md
- 07-router-guards.md

## 验收标准

- [ ] UserList 页面仅管理员可见
- [ ] 用户列表支持分页
- [ ] 用户列表支持关键字搜索
- [ ] 支持单选和批量选择删除
- [ ] 添加用户弹窗包含所有必填字段
- [ ] 禁止删除 root 用户
- [ ] Profile 页面显示当前用户信息
- [ ] 修改个人信息成功
- [ ] 修改密码验证旧密码
