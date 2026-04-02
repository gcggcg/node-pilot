import { defineStore } from 'pinia';
import { ref } from 'vue';
import { fileUploadApi } from '@/api';
import type { FileUpload, FileUploadServer, FileUploadForm } from '@/types';

export const useFileUploadStore = defineStore('fileupload', () => {
    const uploads = ref<FileUpload[]>([]);
    const currentUpload = ref<FileUpload | null>(null);
    const results = ref<FileUploadServer[]>([]);
    const loading = ref(false);
    const pagination = ref({ page: 1, pageSize: 10, total: 0 });

    async function fetchUploads(page = 1, pageSize = 10, filters: {
        status?: string;
        keyword?: string;
        startTime?: string;
        endTime?: string;
    } = {}) {
        loading.value = true;
        try {
            const res = await fileUploadApi.list({ page, pageSize, ...filters });
            uploads.value = res.data;
            pagination.value = { page: res.page, pageSize: res.pageSize, total: res.total };
        } finally {
            loading.value = false;
        }
    }

    async function createUpload(form: FileUploadForm) {
        loading.value = true;
        try {
            await fileUploadApi.create(form);
            await fetchUploads();
        } finally {
            loading.value = false;
        }
    }

    async function updateUpload(id: number, data: FileUploadForm) {
        loading.value = true;
        try {
            await fileUploadApi.update(id, data);
            await fetchUploads(pagination.value.page, pagination.value.pageSize);
        } finally {
            loading.value = false;
        }
    }

    async function executeUpload(id: number) {
        loading.value = true;
        try {
            await fileUploadApi.execute(id);
        } finally {
            loading.value = false;
        }
    }

    async function fetchResults(id: number) {
        loading.value = true;
        try {
            const res = await fileUploadApi.getResults(id);
            results.value = res.results;
            currentUpload.value = res.file_upload;
        } finally {
            loading.value = false;
        }
    }

    async function deleteUploads(ids: number[]) {
        loading.value = true;
        try {
            await fileUploadApi.delete(ids);
            await fetchUploads(pagination.value.page, pagination.value.pageSize);
        } finally {
            loading.value = false;
        }
    }

    function resetResults() {
        results.value = [];
        currentUpload.value = null;
    }

    return {
        uploads,
        currentUpload,
        results,
        loading,
        pagination,
        fetchUploads,
        createUpload,
        updateUpload,
        executeUpload,
        fetchResults,
        deleteUploads,
        resetResults,
    };
});
