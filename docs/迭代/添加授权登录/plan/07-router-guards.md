# 创建前端路由守卫

## 任务描述

实现前端路由守卫，保护需要认证的页面，并处理管理员专属路由。

## 详细说明

### 1. 创建路由守卫模块

创建 `frontend/src/router/auth-guard.ts`：

```typescript
import type { Router } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

export function setupAuthGuard(router: Router) {
    router.beforeEach(async (to, from, next) => {
        const authStore = useAuthStore();
        
        // 需要认证的路由
        const requiresAuth = to.meta.requiresAuth !== false;
        // 需要管理员权限的路由
        const requiresAdmin = to.meta.requiresAdmin === true;

        // 如果已登录，尝试获取用户信息
        if (authStore.isLoggedIn && !authStore.user) {
            await authStore.fetchUser();
        }

        // 检查是否需要登录
        if (requiresAuth && !authStore.isLoggedIn) {
            next({ name: 'login', query: { redirect: to.fullPath } });
            return;
        }

        // 检查是否已登录却访问登录页
        if (to.name === 'login' && authStore.isLoggedIn) {
            next({ name: 'home' });
            return;
        }

        // 检查管理员权限
        if (requiresAdmin && !authStore.isAdmin) {
            next({ name: 'home' });
            return;
        }

        next();
    });
}
```

### 2. 更新路由配置

修改 `frontend/src/router/index.ts`：

```typescript
import { createRouter, createWebHistory } from 'vue-router';
import { setupAuthGuard } from './auth-guard';

const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/',
            redirect: '/servers'
        },
        {
            path: '/login',
            name: 'login',
            component: () => import('@/views/Login.vue'),
            meta: { requiresAuth: false }
        },
        {
            path: '/servers',
            name: 'servers',
            component: () => import('@/views/ServerList.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/servers/new',
            name: 'server-new',
            component: () => import('@/views/ServerForm.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/servers/:id/edit',
            name: 'server-edit',
            component: () => import('@/views/ServerForm.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/scripts',
            name: 'scripts',
            component: () => import('@/views/ScriptList.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/scripts/new',
            name: 'script-new',
            component: () => import('@/views/ScriptForm.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/scripts/:id/edit',
            name: 'script-edit',
            component: () => import('@/views/ScriptForm.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/tasks',
            name: 'tasks',
            component: () => import('@/views/TaskList.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/tasks/:id/output',
            name: 'task-output',
            component: () => import('@/views/TaskOutput.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/users',
            name: 'users',
            component: () => import('@/views/UserList.vue'),
            meta: { requiresAuth: true, requiresAdmin: true }
        },
        {
            path: '/profile',
            name: 'profile',
            component: () => import('@/views/Profile.vue'),
            meta: { requiresAuth: true }
        },
    ]
});

// 设置路由守卫
setupAuthGuard(router);

export default router;
```

### 3. 更新 NavBar 组件

修改 `frontend/src/components/NavBar.vue`，根据登录状态显示不同内容：

```vue
<template>
    <div class="navbar">
        <div class="logo" @click="router.push('/')">NodePilot</div>
        <nav class="nav-links" v-if="authStore.isLoggedIn">
            <router-link to="/servers" class="nav-link">服务器</router-link>
            <router-link to="/scripts" class="nav-link">脚本</router-link>
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
/* ... existing styles ... */
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
```

## 输入

- 需求文档 `05-添加授权登录.md`
- 现有的 `frontend/src/router/index.ts`
- 现有的 `frontend/src/components/NavBar.vue`

## 输出

- 新建 `frontend/src/router/auth-guard.ts` - 路由守卫模块
- 修改 `frontend/src/router/index.ts` - 添加 meta 和守卫
- 修改 `frontend/src/components/NavBar.vue` - 根据登录状态显示

## 依赖

- 05-frontend-auth-api-store.md (Auth Store 实现)

## 验收标准

- [ ] 未登录访问需要认证的页面时跳转到 /login
- [ ] 登录后访问 /login 自动跳转到首页
- [ ] 非管理员访问 /users 跳转到首页
- [ ] 导航栏根据登录状态显示不同内容
- [ ] 登录后显示用户名和退出按钮
- [ ] 退出登录后清除 token 并跳转登录页
- [ ] token 过期自动刷新或跳转登录
