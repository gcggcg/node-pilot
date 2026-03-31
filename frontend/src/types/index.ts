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
    script_id: number;
    name: string;
    server_ids: number[];
}
