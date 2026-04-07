export interface Server {
    id: number;
    name: string;
    host: string;
    port: number;
    username: string;
    password_encrypted?: string;
    connection_status: 'online' | 'offline' | 'unknown';
    created_at: string;
    updated_at: string;
}

export interface Script {
    id: number;
    name: string;
    description: string;
    content: string;
    target_path: string;
    created_at: string;
    updated_at: string;
}

export interface Task {
    id: number;
    script_id: number;
    script_ids: string;
    name: string;
    status: 'pending' | 'running' | 'completed' | 'cancelled' | 'failed';
    created_at: string;
    started_at?: string;
    finished_at?: string;
}

export interface TaskServer {
    id: number;
    task_id: number;
    server_id: number;
    server_name?: string;
    status: 'pending' | 'running' | 'success' | 'failed';
    output: string;
    error: string;
    started_at?: string;
    finished_at?: string;
}

export interface WSMessage {
    type: 'output' | 'server_start' | 'server_done' | 'task_start' | 'task_done';
    task_id: number;
    server_id?: number;
    server_name?: string;
    content?: string;
    status?: string;
    exit_code?: number;
    timestamp: string;
    total?: number;
    success?: number;
    failed?: number;
}

export interface ServerForm {
    name: string;
    host: string;
    port: number;
    username: string;
    password: string;
}

export interface ScriptForm {
    name: string;
    description: string;
    content: string;
    target_path: string;
}

export interface TaskForm {
    script_id?: number;
    script_ids?: string;
    name: string;
    server_ids: number[];
}

export interface PaginatedResponse<T> {
    data: T[];
    total: number;
    page: number;
    pageSize: number;
}

export interface PaginationParams {
    page?: number;
    pageSize?: number;
}

export interface User {
    id: number;
    username: string;
    email: string;
    phone: string;
    role: 'ROLE_ADMIN' | 'ROLE_USER';
    created_at: string;
    updated_at: string;
}

export interface LoginRequest {
    username: string;
    password: string;
}

export interface LoginResponse {
    access_token: string;
    refresh_token: string;
    token_type: string;
    expires_in: number;
}

export interface RefreshRequest {
    refresh_token: string;
}

export interface UpdateProfileRequest {
    email?: string;
    phone?: string;
}

export interface ChangePasswordRequest {
    old_password: string;
    new_password: string;
}

export interface CreateUserRequest {
    username: string;
    password: string;
    email?: string;
    phone?: string;
    role: 'ROLE_ADMIN' | 'ROLE_USER';
}

export interface UserListResponse {
    data: User[];
    total: number;
    page: number;
    pageSize: number;
}

export interface FileUpload {
    id: number;
    name: string;
    local_path: string;
    remote_path: string;
    status: 'pending' | 'success' | 'failed';
    created_at: string;
    updated_at: string;
    servers?: Server[];
}

export interface FileUploadServer {
    id: number;
    file_upload_id: number;
    server_id: number;
    server_name?: string;
    status: 'pending' | 'success' | 'failed';
    error_message?: string;
    file_name?: string;
    remote_full_path?: string;
    created_at: string;
}

export interface FileUploadForm {
    name: string;
    localPath: string;
    remotePath: string;
    serverIds: number[];
}

export interface FileUploadListResponse {
    data: FileUpload[];
    total: number;
    page: number;
    pageSize: number;
}

export interface FileUploadResultResponse {
    file_upload: FileUpload;
    results: FileUploadServer[];
}
