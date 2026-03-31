import { defineStore } from 'pinia';
import { ref } from 'vue';
import { scriptApi } from '@/api';
import type { Script, ScriptForm } from '@/types';

export const useScriptStore = defineStore('script', () => {
    const scripts = ref<Script[]>([]);
    const loading = ref(false);
    const error = ref<string | null>(null);

    async function fetchScripts() {
        loading.value = true;
        error.value = null;
        try {
            scripts.value = await scriptApi.list();
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
            await fetchScripts();
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
            await fetchScripts();
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
            await fetchScripts();
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
            await fetchScripts();
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
        fetchScripts,
        createScript,
        updateScript,
        deleteScript,
        deleteScripts,
    };
});
