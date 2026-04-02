import axios from 'axios';
import type { ServerForm, ScriptForm, TaskForm, PaginatedResponse, PaginationParams, LoginResponse, User, RefreshRequest, UpdateProfileRequest, ChangePasswordRequest, CreateUserRequest, UserListResponse } from '@/types';

const api = axios.create({
    baseURL: '/api',
    timeout: 30000,
    headers: {
        'Content-Type': 'application/json'
    }
});

api.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('access_token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => Promise.reject(error)
);

api.interceptors.response.use(
    response => response.data,
    async (error) => {
        const originalRequest = error.config;
        
        if (error.response?.status === 401 && !originalRequest._retry) {
            originalRequest._retry = true;
            
            const refreshToken = localStorage.getItem('refresh_token');
            if (refreshToken) {
                try {
                    const response = await api.post<any, LoginResponse>('/v1/auth/refresh', { refresh_token: refreshToken } as RefreshRequest);
                    localStorage.setItem('access_token', response.access_token);
                    localStorage.setItem('refresh_token', response.refresh_token);
                    
                    originalRequest.headers.Authorization = `Bearer ${response.access_token}`;
                    return api(originalRequest);
                } catch (refreshError) {
                    localStorage.removeItem('access_token');
                    localStorage.removeItem('refresh_token');
                    window.location.href = '/login';
                    return Promise.reject(refreshError);
                }
            }
        }
        
        const message = error.response?.data?.error || error.message || 'Request failed';
        console.error('API Error:', message);
        return Promise.reject(error);
    }
);

export const authApi = {
    login: (data: { username: string; password: string }) => api.post<any, LoginResponse>('/v1/auth/login', data),
    me: () => api.get<any, User>('/v1/auth/me'),
    refresh: (data: RefreshRequest) => api.post<any, LoginResponse>('/v1/auth/refresh', data),
    updateProfile: (data: UpdateProfileRequest) => api.put('/v1/auth/profile', data),
    changePassword: (data: ChangePasswordRequest) => api.put('/v1/auth/password', data),
};

export const userApi = {
    list: (params?: { page?: number; pageSize?: number; keyword?: string }) => {
        const query = params
            ? `?page=${params.page || 1}&pageSize=${params.pageSize || 10}${params.keyword ? `&keyword=${params.keyword}` : ''}`
            : '';
        return api.get<any, UserListResponse>(`/v1/admin/users${query}`);
    },
    create: (data: CreateUserRequest) => api.post('/v1/admin/users', data),
    delete: (id: number) => api.delete(`/v1/admin/users/${id}`),
    deleteMany: (ids: number[]) => api.post('/v1/admin/users/batch-delete', { ids }),
};

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
