# 创建前端认证 API 和状态管理

## 任务描述

在 Vue 前端添加认证相关的 API 调用封装和 Pinia 状态管理。

## 详细说明

### 1. 添加 User 类型

在 `frontend/src/types/index.ts` 添加：

```typescript
export interface User {
    id: number;
    username: string;
    email: string;
    phone: string;
    role: 'ROLE_ADMIN' | 'ROLE_USER';
    created_at: string;
    updated_at: string;
}

export interface LoginRequest {
    username: string;
    password: string;
}

export interface LoginResponse {
    access_token: string;
    refresh_token: string;
    token_type: string;
    expires_in: number;
}

export interface RefreshRequest {
    refresh_token: string;
}

export interface UpdateProfileRequest {
    email?: string;
    phone?: string;
}

export interface ChangePasswordRequest {
    old_password: string;
    new_password: string;
}

export interface CreateUserRequest {
    username: string;
    password: string;
    email?: string;
    phone?: string;
    role: 'ROLE_ADMIN' | 'ROLE_USER';
}

export interface UserListResponse {
    data: User[];
    total: number;
    page: number;
    pageSize: number;
}
```

### 2. 添加认证 API

在 `frontend/src/api/index.ts` 添加：

```typescript
export const authApi = {
    login: (data: LoginRequest) => api.post<any, LoginResponse>('/v1/auth/login', data),
    me: () => api.get<any, User>('/v1/auth/me'),
    refresh: (data: RefreshRequest) => api.post<any, LoginResponse>('/v1/auth/refresh', data),
    updateProfile: (data: UpdateProfileRequest) => api.put('/v1/auth/profile', data),
    changePassword: (data: ChangePasswordRequest) => api.put('/v1/auth/password', data),
};

export const userApi = {
    list: (params?: { page?: number; pageSize?: number; keyword?: string }) => {
        const query = params
            ? `?page=${params.page || 1}&pageSize=${params.pageSize || 10}${params.keyword ? `&keyword=${params.keyword}` : ''}`
            : '';
        return api.get<any, UserListResponse>(`/v1/admin/users${query}`);
    },
    create: (data: CreateUserRequest) => api.post('/v1/admin/users', data),
    delete: (id: number) => api.delete(`/v1/admin/users/${id}`),
    deleteMany: (ids: number[]) => api.post('/v1/admin/users/batch-delete', { ids }),
};
```

### 3. 添加 Auth Store

创建 `frontend/src/stores/auth.ts`：

```typescript
import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { authApi } from '@/api';
import type { User, LoginRequest } from '@/types';

const TOKEN_KEY = 'access_token';
const REFRESH_TOKEN_KEY = 'refresh_token';

export const useAuthStore = defineStore('auth', () => {
    const token = ref<string | null>(localStorage.getItem(TOKEN_KEY));
    const refreshToken = ref<string | null>(localStorage.getItem(REFRESH_TOKEN_KEY));
    const user = ref<User | null>(null);
    const loading = ref(false);

    const isLoggedIn = computed(() => !!token.value);
    const isAdmin = computed(() => user.value?.role === 'ROLE_ADMIN');

    async function login(credentials: LoginRequest) {
        loading.value = true;
        try {
            const response = await authApi.login(credentials);
            token.value = response.access_token;
            refreshToken.value = response.refresh_token;
            localStorage.setItem(TOKEN_KEY, response.access_token);
            localStorage.setItem(REFRESH_TOKEN_KEY, response.refresh_token);
            await fetchUser();
            return true;
        } catch (error) {
            console.error('Login failed:', error);
            return false;
        } finally {
            loading.value = false;
        }
    }

    async function fetchUser() {
        if (!token.value) return;
        try {
            user.value = await authApi.me();
        } catch (error) {
            logout();
        }
    }

    async function refreshAccessToken() {
        if (!refreshToken.value) return false;
        try {
            const response = await authApi.refresh({ refresh_token: refreshToken.value });
            token.value = response.access_token;
            refreshToken.value = response.refresh_token;
            localStorage.setItem(TOKEN_KEY, response.access_token);
            localStorage.setItem(REFRESH_TOKEN_KEY, response.refresh_token);
            return true;
        } catch (error) {
            logout();
            return false;
        }
    }

    async function updateProfile(data: { email?: string; phone?: string }) {
        await authApi.updateProfile(data);
        await fetchUser();
    }

    async function changePassword(oldPassword: string, newPassword: string) {
        await authApi.changePassword({
            old_password: oldPassword,
            new_password: newPassword,
        });
    }

    function logout() {
        token.value = null;
        refreshToken.value = null;
        user.value = null;
        localStorage.removeItem(TOKEN_KEY);
        localStorage.removeItem(REFRESH_TOKEN_KEY);
    }

    // 初始化时获取用户信息
    if (token.value) {
        fetchUser();
    }

    return {
        token,
        refreshToken,
        user,
        loading,
        isLoggedIn,
        isAdmin,
        login,
        fetchUser,
        refreshAccessToken,
        updateProfile,
        changePassword,
        logout,
    };
});
```

### 4. 更新 Axios 拦截器

修改 `frontend/src/api/index.ts`，添加请求拦截器自动附加 Token：

```typescript
api.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('access_token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => Promise.reject(error)
);

api.interceptors.response.use(
    response => response.data,
    async (error) => {
        const originalRequest = error.config;
        
        // 如果是 401 错误且没有重试过，尝试刷新 token
        if (error.response?.status === 401 && !originalRequest._retry) {
            originalRequest._retry = true;
            
            const refreshToken = localStorage.getItem('refresh_token');
            if (refreshToken) {
                try {
                    const response = await api.post('/v1/auth/refresh', { refresh_token: refreshToken });
                    localStorage.setItem('access_token', response.access_token);
                    localStorage.setItem('refresh_token', response.refresh_token);
                    
                    // 重试原始请求
                    originalRequest.headers.Authorization = `Bearer ${response.access_token}`;
                    return api(originalRequest);
                } catch (refreshError) {
                    // 刷新失败，清除 token
                    localStorage.removeItem('access_token');
                    localStorage.removeItem('refresh_token');
                    window.location.href = '/login';
                    return Promise.reject(refreshError);
                }
            }
        }
        
        const message = error.response?.data?.error || error.message || 'Request failed';
        console.error('API Error:', message);
        return Promise.reject(error);
    }
);
```

## 输入

- 需求文档 `05-添加授权登录.md`
- 现有的 `frontend/src/api/index.ts`
- 现有的 `frontend/src/types/index.ts`
- 现有的 Pinia store 模式

## 输出

- 修改 `frontend/src/types/index.ts` - 添加 User 和认证相关类型
- 修改 `frontend/src/api/index.ts` - 添加认证 API 和 Token 拦截器
- 新建 `frontend/src/stores/auth.ts` - 认证状态管理

## 依赖

- 后端 API 实现完成

## 验收标准

- [ ] LoginRequest, LoginResponse, User 等类型正确定义
- [ ] authApi.login 正确调用登录接口
- [ ] authApi.me 正确获取当前用户
- [ ] Axios 请求拦截器自动附加 Authorization header
- [ ] 401 错误时自动尝试刷新 token
- [ ] 刷新失败时清除 token 并跳转登录
- [ ] Auth Store 正确管理登录状态
- [ ] logout 函数正确清除本地存储
