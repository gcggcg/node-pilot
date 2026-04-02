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
