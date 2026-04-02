<template>
    <div class="page">
        <h1>个人中心</h1>

        <div class="profile-card">
            <div class="section">
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
                    <div class="info-item">
                        <label>邮箱</label>
                        <span>{{ authStore.user?.email || '-' }}</span>
                    </div>
                    <div class="info-item">
                        <label>电话</label>
                        <span>{{ authStore.user?.phone || '-' }}</span>
                    </div>
                </div>
            </div>

            <div class="section">
                <h3>修改个人信息</h3>
                <form @submit.prevent="handleUpdateProfile" class="form">
                    <div class="form-group">
                        <label>邮箱</label>
                        <input v-model="profileForm.email" type="email" placeholder="请输入邮箱" />
                    </div>
                    <div class="form-group">
                        <label>电话</label>
                        <input v-model="profileForm.phone" placeholder="请输入电话" />
                    </div>
                    <div class="form-actions">
                        <button type="submit" class="btn btn-primary">保存</button>
                    </div>
                </form>
            </div>

            <div class="section">
                <h3>修改密码</h3>
                <form @submit.prevent="handleChangePassword" class="form">
                    <div class="form-group">
                        <label>旧密码</label>
                        <input v-model="passwordForm.old_password" type="password" required placeholder="请输入旧密码" />
                    </div>
                    <div class="form-group">
                        <label>新密码</label>
                        <input v-model="passwordForm.new_password" type="password" required minlength="6" placeholder="请输入新密码（至少6位）" />
                    </div>
                    <div class="form-actions">
                        <button type="submit" class="btn btn-primary">修改密码</button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useAuthStore } from '@/stores/auth';

const authStore = useAuthStore();

const profileForm = ref({
    email: '',
    phone: ''
});

const passwordForm = ref({
    old_password: '',
    new_password: ''
});

onMounted(() => {
    if (authStore.user) {
        profileForm.value.email = authStore.user.email || '';
        profileForm.value.phone = authStore.user.phone || '';
    }
});

async function handleUpdateProfile() {
    try {
        await authStore.updateProfile(profileForm.value);
        alert('个人信息已更新');
    } catch (e: any) {
        alert(e.response?.data?.error || '更新失败');
    }
}

async function handleChangePassword() {
    if (passwordForm.value.new_password.length < 6) {
        alert('新密码至少6位');
        return;
    }
    try {
        await authStore.changePassword(passwordForm.value.old_password, passwordForm.value.new_password);
        alert('密码已修改');
        passwordForm.value = { old_password: '', new_password: '' };
    } catch (e: any) {
        alert(e.response?.data?.error || '密码修改失败');
    }
}
</script>

<style scoped>
.page {
    padding: 24px;
    max-width: 800px;
}

h1 {
    font-size: 1.5rem;
    color: #333;
    margin-bottom: 24px;
}

.profile-card {
    background: white;
    border-radius: 12px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    overflow: hidden;
}

.section {
    padding: 24px;
    border-bottom: 1px solid #eee;
}

.section:last-child {
    border-bottom: none;
}

.section h3 {
    margin: 0 0 20px 0;
    font-size: 1rem;
    color: #333;
    font-weight: 600;
}

.info-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
}

.info-item {
    display: flex;
    flex-direction: column;
    gap: 4px;
}

.info-item label {
    font-size: 13px;
    color: #888;
}

.info-item span {
    font-size: 14px;
    color: #333;
}

.role-badge {
    display: inline-block;
    padding: 4px 10px;
    border-radius: 12px;
    font-size: 12px;
    font-weight: 500;
    width: fit-content;
}

.role-badge.admin {
    background: #d4edda;
    color: #155724;
}

.role-badge.user {
    background: #e9ecef;
    color: #6c757d;
}

.form {
    max-width: 400px;
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

.form-group input {
    width: 100%;
    padding: 10px 12px;
    border: 1px solid #ddd;
    border-radius: 6px;
    font-size: 14px;
    box-sizing: border-box;
}

.form-group input:focus {
    outline: none;
    border-color: #667eea;
}

.form-actions {
    display: flex;
    gap: 12px;
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
</style>
