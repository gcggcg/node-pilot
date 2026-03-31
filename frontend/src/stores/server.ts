import { defineStore } from 'pinia';
import { ref } from 'vue';
import { serverApi } from '@/api';
import type { Server, ServerForm } from '@/types';

export const useServerStore = defineStore('server', () => {
    const servers = ref<Server[]>([]);
    const loading = ref(false);
    const error = ref<string | null>(null);

    async function fetchServers() {
        loading.value = true;
        error.value = null;
        try {
            servers.value = await serverApi.list();
        } catch (e: any) {
            error.value = e.message;
        } finally {
            loading.value = false;
        }
    }

    async function createServer(data: ServerForm) {
        loading.value = true;
        error.value = null;
        try {
            const result = await serverApi.create(data);
            await fetchServers();
            return result;
        } catch (e: any) {
            error.value = e.message;
            throw e;
        } finally {
            loading.value = false;
        }
    }

    async function updateServer(id: number, data: ServerForm) {
        loading.value = true;
        error.value = null;
        try {
            await serverApi.update(id, data);
            await fetchServers();
        } catch (e: any) {
            error.value = e.message;
            throw e;
        } finally {
            loading.value = false;
        }
    }

    async function deleteServer(id: number) {
        loading.value = true;
        error.value = null;
        try {
            await serverApi.delete(id);
            await fetchServers();
        } catch (e: any) {
            error.value = e.message;
            throw e;
        } finally {
            loading.value = false;
        }
    }

    async function deleteServers(ids: number[]) {
        loading.value = true;
        error.value = null;
        try {
            await serverApi.deleteMany(ids);
            await fetchServers();
        } catch (e: any) {
            error.value = e.message;
            throw e;
        } finally {
            loading.value = false;
        }
    }

    async function testConnection(id: number) {
        try {
            return await serverApi.test(id);
        } catch (e: any) {
            error.value = e.message;
            throw e;
        }
    }

    return {
        servers,
        loading,
        error,
        fetchServers,
        createServer,
        updateServer,
        deleteServer,
        deleteServers,
        testConnection,
    };
});
