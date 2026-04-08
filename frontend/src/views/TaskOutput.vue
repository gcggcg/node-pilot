<template>
    <div class="page">
        <div class="header">
            <div class="header-left">
                <h1>{{ task?.name || '任务输出' }}</h1>
                <span v-if="task" :class="['status', task.status]">{{ statusText(task.status) }}</span>
            </div>
            <router-link to="/tasks" class="btn btn-secondary">返回任务列表</router-link>
        </div>

        <div class="content">
            <OutputPanel :outputs="taskStore.outputs" :servers="serverStatusesWithError" />
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useRoute } from 'vue-router';
import { useTaskStore } from '@/stores/task';
import { taskApi } from '@/api';
import OutputPanel from '@/components/OutputPanel.vue';

const route = useRoute();
const taskStore = useTaskStore();

const task = ref<any>(null);
const taskServers = ref<any[]>([]);

const taskId = computed(() => Number(route.params.id));

const serverStatusesWithError = computed(() => {
    return taskServers.value.map(s => ({
        id: s.server_id,
        name: s.server_name || `Server ${s.server_id}`,
        status: s.status,
        error: s.error
    }));
});

onMounted(async () => {
    try {
        const response = await taskApi.get(taskId.value);
        task.value = response.task;
        taskServers.value = response.servers || [];
        
        // For running tasks, clear old outputs to avoid stale data
        if (task.value?.status === 'running') {
            taskStore.outputs = {};
        } else {
            // Pre-populate outputs with stored output from DB for completed tasks
            for (const server of taskServers.value) {
                if (server.status === 'failed' && server.error) {
                    taskStore.outputs[server.server_id] = `❌ 执行失败: ${server.error}`;
                } else if (server.output) {
                    taskStore.outputs[server.server_id] = server.output;
                }
            }
        }
    } catch (e) {
        console.error('Failed to load task:', e);
    }

    taskStore.connectWebSocket(taskId.value);
});

onUnmounted(() => {
    taskStore.disconnectWebSocket();
});

function statusText(status: string): string {
    const statusMap: Record<string, string> = {
        'pending': '等待中',
        'running': '运行中',
        'completed': '已完成',
        'cancelled': '已取消',
        'failed': '失败'
    };
    return statusMap[status] || status;
}
</script>

<style scoped>
.page {
    padding: 24px;
    height: calc(100vh - 60px);
    display: flex;
    flex-direction: column;
}

.header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
}

.header-left {
    display: flex;
    align-items: center;
    gap: 16px;
}

h1 {
    font-size: 1.5rem;
    color: #333;
}

.status {
    padding: 4px 12px;
    border-radius: 4px;
    font-size: 14px;
    font-weight: 500;
}

.status.pending { background: #e9ecef; color: #495057; }
.status.running { background: #cce5ff; color: #004085; }
.status.completed { background: #d4edda; color: #155724; }
.status.cancelled { background: #fff3cd; color: #856404; }
.status.failed { background: #f8d7da; color: #721c24; }

.content {
    flex: 1;
    min-height: 0;
}

.btn {
    padding: 8px 16px;
    border: none;
    border-radius: 6px;
    cursor: pointer;
    font-size: 14px;
    text-decoration: none;
}

.btn-secondary {
    background: #e9ecef;
    color: #495057;
}
</style>
