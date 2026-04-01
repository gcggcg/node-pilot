import axios from 'axios';
import type { ServerForm, ScriptForm, TaskForm, PaginatedResponse, PaginationParams } from '@/types';

const api = axios.create({
    baseURL: '/api',
    timeout: 30000,
    headers: {
        'Content-Type': 'application/json'
    }
});

api.interceptors.response.use(
    response => response.data,
    error => {
        const message = error.response?.data?.error || error.message || 'Request failed';
        console.error('API Error:', message);
        return Promise.reject(error);
    }
);

export const serverApi = {
    list: (params?: PaginationParams) => {
        const query = params 
            ? `?page=${params.page || 1}&pageSize=${params.pageSize || 10}`
            : '';
        return api.get<any, PaginatedResponse<any>>(`/servers${query}`);
    },
    get: (id: number) => api.get<any, any>(`/servers/${id}`),
    create: (data: ServerForm) => api.post('/servers', data),
    update: (id: number, data: ServerForm) => api.put(`/servers/${id}`, data),
    delete: (id: number) => api.delete(`/servers/${id}`),
    deleteMany: (ids: number[]) => api.post('/servers/batch-delete', { ids }),
    test: (id: number) => api.post(`/servers/${id}/test`),
};

export const scriptApi = {
    list: (params?: PaginationParams) => {
        const query = params 
            ? `?page=${params.page || 1}&pageSize=${params.pageSize || 10}`
            : '';
        return api.get<any, PaginatedResponse<any>>(`/scripts${query}`);
    },
    get: (id: number) => api.get<any, any>(`/scripts/${id}`),
    create: (data: ScriptForm) => api.post('/scripts', data),
    update: (id: number, data: ScriptForm) => api.put(`/scripts/${id}`, data),
    delete: (id: number) => api.delete(`/scripts/${id}`),
    deleteMany: (ids: number[]) => api.post('/scripts/batch-delete', { ids }),
};

export const taskApi = {
    list: (params?: PaginationParams) => {
        const query = params 
            ? `?page=${params.page || 1}&pageSize=${params.pageSize || 10}`
            : '';
        return api.get<any, PaginatedResponse<any>>(`/tasks${query}`);
    },
    get: (id: number) => api.get<any, any>(`/tasks/${id}`),
    create: (data: TaskForm) => api.post('/tasks', data),
    cancel: (id: number) => api.delete(`/tasks/${id}`),
    deleteMany: (ids: number[]) => api.post('/tasks/batch-delete', { ids }),
};

export default api;
