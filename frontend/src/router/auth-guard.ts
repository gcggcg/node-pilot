import type { Router } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

export function setupAuthGuard(router: Router) {
    router.beforeEach(async (to, _from, next) => {
        const authStore = useAuthStore();
        
        const requiresAuth = to.meta.requiresAuth !== false;
        const requiresAdmin = to.meta.requiresAdmin === true;

        if (authStore.isLoggedIn && !authStore.user) {
            await authStore.fetchUser();
        }

        if (requiresAuth && !authStore.isLoggedIn) {
            next({ name: 'login', query: { redirect: to.fullPath } });
            return;
        }

        if (to.name === 'login' && authStore.isLoggedIn) {
            next({ name: 'home' });
            return;
        }

        if (requiresAdmin && !authStore.isAdmin) {
            next({ name: 'home' });
            return;
        }

        next();
    });
}
