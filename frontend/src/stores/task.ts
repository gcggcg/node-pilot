import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { taskApi } from '@/api';
import type { Task, TaskForm, WSMessage } from '@/types';

export const useTaskStore = defineStore('task', () => {
    const tasks = ref<Task[]>([]);
    const currentTask = ref<Task | null>(null);
    const currentTaskId = ref<number | null>(null);
    const loading = ref(false);
    const error = ref<string | null>(null);
    const outputs = ref<Map<number, string>>(new Map());
    const ws = ref<WebSocket | null>(null);

    const taskStats = computed(() => {
        const running = tasks.value.filter(t => t.status === 'running').length;
        const completed = tasks.value.filter(t => t.status === 'completed').length;
        const cancelled = tasks.value.filter(t => t.status === 'cancelled').length;
        return { running, completed, cancelled, total: tasks.value.length };
    });

    async function fetchTasks() {
        loading.value = true;
        error.value = null;
        try {
            tasks.value = await taskApi.list();
        } catch (e: any) {
            error.value = e.message;
        } finally {
            loading.value = false;
        }
    }

    async function createTask(data: TaskForm) {
        loading.value = true;
        error.value = null;
        try {
            const result = await taskApi.create(data);
            await fetchTasks();
            return result;
        } catch (e: any) {
            error.value = e.message;
            throw e;
        } finally {
            loading.value = false;
        }
    }

    async function cancelTask(id: number) {
        try {
            await taskApi.cancel(id);
            await fetchTasks();
        } catch (e: any) {
            error.value = e.message;
            throw e;
        }
    }

    async function deleteTasks(ids: number[]) {
        loading.value = true;
        error.value = null;
        try {
            await taskApi.deleteMany(ids);
            await fetchTasks();
        } catch (e: any) {
            error.value = e.message;
            throw e;
        } finally {
            loading.value = false;
        }
    }

    function connectWebSocket(taskId: number) {
        currentTaskId.value = taskId;
        
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const host = window.location.host;
        const wsUrl = `${protocol}//${host}/ws?task_id=${taskId}`;
        
        console.log('[WS] Connecting to:', wsUrl);
        
        ws.value = new WebSocket(wsUrl);

        ws.value.onopen = () => {
            console.log('[WS] Connected successfully');
        };

        ws.value.onmessage = (event) => {
            console.log('[WS] Received message:', event.data);
            const messages = event.data.split('\n');
            for (const msg of messages) {
                if (!msg.trim()) continue;
                try {
                    const data: WSMessage = JSON.parse(msg);
                    handleWSMessage(data);
                } catch (e) {
                    console.error('[WS] Failed to parse message:', e);
                }
            }
        };

        ws.value.onerror = (error) => {
            console.error('[WS] Error:', error);
        };

        ws.value.onclose = (event) => {
            console.log('[WS] Closed:', event.code, event.reason);
        };
    }

    function handleWSMessage(data: WSMessage) {
        if (currentTaskId.value !== null && data.task_id !== currentTaskId.value) {
            return;
        }
        if (data.type === 'output' && data.server_id) {
            const current = outputs.value.get(data.server_id) || '';
            outputs.value.set(data.server_id, current + (data.content || ''));
        } else if (data.type === 'task_done' || data.type === 'server_done') {
            fetchTasks();
        }
    }

    function disconnectWebSocket() {
        if (ws.value) {
            ws.value.close();
            ws.value = null;
        }
        currentTaskId.value = null;
    }

    function clearOutputs() {
        outputs.value.clear();
    }

    return {
        tasks,
        currentTask,
        loading,
        error,
        outputs,
        taskStats,
        fetchTasks,
        createTask,
        cancelTask,
        deleteTasks,
        connectWebSocket,
        disconnectWebSocket,
        clearOutputs,
    };
});
