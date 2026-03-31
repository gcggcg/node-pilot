import { createRouter, createWebHistory } from 'vue-router';

const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/',
            redirect: '/servers'
        },
        {
            path: '/servers',
            name: 'servers',
            component: () => import('@/views/ServerList.vue')
        },
        {
            path: '/servers/new',
            name: 'server-new',
            component: () => import('@/views/ServerForm.vue')
        },
        {
            path: '/servers/:id/edit',
            name: 'server-edit',
            component: () => import('@/views/ServerForm.vue')
        },
        {
            path: '/scripts',
            name: 'scripts',
            component: () => import('@/views/ScriptList.vue')
        },
        {
            path: '/scripts/new',
            name: 'script-new',
            component: () => import('@/views/ScriptForm.vue')
        },
        {
            path: '/scripts/:id/edit',
            name: 'script-edit',
            component: () => import('@/views/ScriptForm.vue')
        },
        {
            path: '/tasks',
            name: 'tasks',
            component: () => import('@/views/TaskList.vue')
        },
        {
            path: '/tasks/:id/output',
            name: 'task-output',
            component: () => import('@/views/TaskOutput.vue')
        },
    ]
});

export default router;
