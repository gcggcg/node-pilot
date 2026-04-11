<template>
    <div class="navbar">
        <div class="logo" @click="router.push('/')">NodePilot</div>
        <nav class="nav-links" v-if="authStore.isLoggedIn">
            <router-link to="/servers" class="nav-link">服务器</router-link>
            <router-link to="/files" class="nav-link" v-if="authStore.isAdmin">文件</router-link>
            <router-link to="/scripts" class="nav-link" v-if="authStore.isAdmin">脚本</router-link>
            <router-link to="/tasks" class="nav-link">任务</router-link>
            <router-link to="/users" class="nav-link" v-if="authStore.isAdmin">用户</router-link>
            <router-link to="/profile" class="nav-link">个人</router-link>
        </nav>
        <div class="nav-actions" v-if="authStore.isLoggedIn">
            <span class="username">{{ authStore.user?.username }}</span>
            <button @click="handleLogout" class="logout-btn">退出</button>
        </div>
    </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

const router = useRouter();
const authStore = useAuthStore();

async function handleLogout() {
    authStore.logout();
    router.push('/login');
}
</script>

<style scoped>
.navbar {
    height: 60px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    display: flex;
    align-items: center;
    padding: 0 24px;
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    z-index: 100;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.logo {
    font-size: 1.5rem;
    font-weight: bold;
    color: white;
    margin-right: 48px;
    cursor: pointer;
}

.nav-links {
    display: flex;
    gap: 24px;
}

.nav-link {
    color: rgba(255, 255, 255, 0.85);
    text-decoration: none;
    font-weight: 500;
    padding: 8px 16px;
    border-radius: 6px;
    transition: all 0.2s;
}

.nav-link:hover {
    color: white;
    background: rgba(255, 255, 255, 0.1);
}

.nav-link.router-link-active {
    color: white;
    background: rgba(255, 255, 255, 0.2);
}

.nav-actions {
    display: flex;
    align-items: center;
    gap: 16px;
    margin-left: auto;
}

.username {
    color: rgba(255, 255, 255, 0.7);
    font-size: 0.9rem;
}

.logout-btn {
    background: rgba(255, 255, 255, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.2);
    border-radius: 6px;
    padding: 6px 12px;
    color: rgba(255, 255, 255, 0.8);
    cursor: pointer;
    transition: all 0.2s;
}

.logout-btn:hover {
    background: rgba(255, 255, 255, 0.2);
    color: white;
}
</style>
