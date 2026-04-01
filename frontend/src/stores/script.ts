import { defineStore } from 'pinia';
import { ref } from 'vue';
import { scriptApi } from '@/api';
import type { Script, ScriptForm } from '@/types';

export const useScriptStore = defineStore('script', () => {
    const scripts = ref<Script[]>([]);
    const loading = ref(false);
    const error = ref<string | null>(null);
    const pagination = ref({
        page: 1,
        pageSize: 10,
        total: 0
    });

    async function fetchScripts(page = 1, pageSize = 10) {
        loading.value = true;
        error.value = null;
        try {
            const res = await scriptApi.list({ page, pageSize });
            scripts.value = res.data;
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

    async function createScript(data: ScriptForm) {
        loading.value = true;
        error.value = null;
        try {
            const result = await scriptApi.create(data);
            await fetchScripts(pagination.value.page, pagination.value.pageSize);
            return result;
        } catch (e: any) {
            error.value = e.message;
            throw e;
        } finally {
            loading.value = false;
        }
    }

    async function updateScript(id: number, data: ScriptForm) {
        loading.value = true;
        error.value = null;
        try {
            await scriptApi.update(id, data);
            await fetchScripts(pagination.value.page, pagination.value.pageSize);
        } catch (e: any) {
            error.value = e.message;
            throw e;
        } finally {
            loading.value = false;
        }
    }

    async function deleteScript(id: number) {
        loading.value = true;
        error.value = null;
        try {
            await scriptApi.delete(id);
            await fetchScripts(pagination.value.page, pagination.value.pageSize);
        } catch (e: any) {
            error.value = e.message;
            throw e;
        } finally {
            loading.value = false;
        }
    }

    async function deleteScripts(ids: number[]) {
        loading.value = true;
        error.value = null;
        try {
            await scriptApi.deleteMany(ids);
            await fetchScripts(pagination.value.page, pagination.value.pageSize);
        } catch (e: any) {
            error.value = e.message;
            throw e;
        } finally {
            loading.value = false;
        }
    }

    return {
        scripts,
        loading,
        error,
        pagination,
        fetchScripts,
        createScript,
        updateScript,
        deleteScript,
        deleteScripts,
    };
});
