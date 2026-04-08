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
            path: '/tasks/new',
            name: 'task-new',
            component: () => import('@/views/TaskForm.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/tasks/:id/edit',
            name: 'task-edit',
            component: () => import('@/views/TaskForm.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/tasks/:id/output',
            name: 'task-output',
            component: () => import('@/views/TaskOutput.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/tasks/:id/detail',
            name: 'task-detail',
            component: () => import('@/views/TaskDetail.vue'),
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
        {
            path: '/files',
            name: 'files',
            component: () => import('@/views/FileList.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/files/new',
            name: 'file-new',
            component: () => import('@/views/FileForm.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/files/:id/edit',
            name: 'file-edit',
            component: () => import('@/views/FileForm.vue'),
            meta: { requiresAuth: true }
        },
    ]
});

setupAuthGuard(router);

export default router;
