import { defineStore } from 'pinia';
import { ref } from 'vue';
import { serverApi } from '@/api';
import type { Server, ServerForm } from '@/types';

export const useServerStore = defineStore('server', () => {
    const servers = ref<Server[]>([]);
    const loading = ref(false);
    const error = ref<string | null>(null);
    const pagination = ref({
        page: 1,
        pageSize: 10,
        total: 0
    });

    async function fetchServers(page = 1, pageSize = 10) {
        loading.value = true;
        error.value = null;
        try {
            const res = await serverApi.list({ page, pageSize });
            servers.value = res.data;
            pagination.value = {
                page: res.page,
                pageSize: res.pageSize,
                total: res.total
            };
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
            await fetchServers(pagination.value.page, pagination.value.pageSize);
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
            await fetchServers(pagination.value.page, pagination.value.pageSize);
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
            await fetchServers(pagination.value.page, pagination.value.pageSize);
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
            await fetchServers(pagination.value.page, pagination.value.pageSize);
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
        pagination,
        fetchServers,
        createServer,
        updateServer,
        deleteServer,
        deleteServers,
        testConnection,
    };
});
