<template>
    <div class="page">
        <div class="header">
            <h1>任务详情</h1>
            <router-link to="/tasks" class="btn btn-secondary">返回任务列表</router-link>
        </div>

        <div v-if="loading" class="loading">加载中...</div>
        <div v-else-if="error" class="error">{{ error }}</div>
        
        <div v-else-if="task" class="content">
            <div class="card">
                <h2>基本信息</h2>
                <div class="info-grid">
                    <div class="info-item">
                        <label>任务名称</label>
                        <span>{{ task.name }}</span>
                    </div>
                    <div class="info-item">
                        <label>状态</label>
                        <span :class="['status-badge', task.status]">{{ statusText(task.status) }}</span>
                    </div>
                    <div class="info-item">
                        <label>执行模式</label>
                        <span :class="['mode-badge', task.execution_mode]">{{ executionModeText(task.execution_mode) }}</span>
                    </div>
                    <div class="info-item">
                        <label>创建时间</label>
                        <span>{{ formatDate(task.created_at) }}</span>
                    </div>
                    <div class="info-item" v-if="task.started_at">
                        <label>开始时间</label>
                        <span>{{ formatDate(task.started_at) }}</span>
                    </div>
                    <div class="info-item" v-if="task.finished_at">
                        <label>结束时间</label>
                        <span>{{ formatDate(task.finished_at) }}</span>
                    </div>
                </div>
            </div>

            <div class="card">
                <h2>执行的脚本</h2>
                <div class="script-list">
                    <div v-if="task.script_ids" class="script-item" v-for="(scriptId, index) in scriptIdList" :key="scriptId">
                        <span class="script-index">{{ index + 1 }}</span>
                        <span class="script-id">ID: {{ scriptId }}</span>
                        <span v-if="getScriptName(scriptId)" class="script-name">{{ getScriptName(scriptId) }}</span>
                        <span v-else class="script-loading">加载中...</span>
                    </div>
                    <div v-else-if="task.script_id" class="script-item">
                        <span class="script-index">1</span>
                        <span class="script-id">ID: {{ task.script_id }}</span>
                        <span v-if="getScriptName(task.script_id)" class="script-name">{{ getScriptName(task.script_id) }}</span>
                    </div>
                    <div v-else class="empty-message">暂无脚本信息</div>
                </div>
            </div>

            <div class="card">
                <h2>目标服务器</h2>
                <div class="server-list">
                    <div v-if="servers.length > 0" class="server-item" v-for="server in servers" :key="server.server_id">
                        <span class="server-name">{{ server.server_name || `服务器 ${server.server_id}` }}</span>
                        <span :class="['server-status', server.status]">{{ serverStatusText(server.status) }}</span>
                        <span v-if="server.error" class="server-error">{{ server.error }}</span>
                    </div>
                    <div v-else class="empty-message">暂无服务器信息</div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { taskApi, scriptApi } from '@/api';
import type { Script } from '@/types';

const route = useRoute();

const task = ref<any>(null);
const servers = ref<any[]>([]);
const scripts = ref<Script[]>([]);
const loading = ref(true);
const error = ref<string | null>(null);

const taskId = computed(() => Number(route.params.id));

const scriptIdList = computed(() => {
    if (!task.value?.script_ids) return [];
    return task.value.script_ids.split(',').map((id: string) => parseInt(id.trim())).filter((id: number) => !isNaN(id));
});

function getScriptName(id: number): string | null {
    const script = scripts.value.find(s => s.id === id);
    return script ? script.name : null;
}

onMounted(async () => {
    try {
        const [taskRes, scriptRes] = await Promise.all([
            taskApi.get(taskId.value),
            scriptApi.list({ page: 1, pageSize: 100 })
        ]);
        task.value = taskRes.task;
        servers.value = taskRes.servers || [];
        scripts.value = scriptRes.data || [];
    } catch (e: any) {
        error.value = e.message || '加载任务失败';
    } finally {
        loading.value = false;
    }
});

function formatDate(dateStr: string): string {
    if (!dateStr) return '-';
    return new Date(dateStr).toLocaleString('zh-CN');
}

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

function executionModeText(mode: string): string {
    const modeMap: Record<string, string> = {
        'concurrent': '并发执行',
        'sequential': '单线程执行'
    };
    return modeMap[mode] || mode;
}

function serverStatusText(status: string): string {
    const statusMap: Record<string, string> = {
        'pending': '等待中',
        'running': '运行中',
        'success': '成功',
        'failed': '失败'
    };
    return statusMap[status] || status;
}
</script>

<style scoped>
.page {
    padding: 24px;
    max-width: 1000px;
}

.header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
}

h1 {
    font-size: 1.5rem;
    color: #333;
}

.content {
    display: flex;
    flex-direction: column;
    gap: 20px;
}

.card {
    background: white;
    padding: 24px;
    border-radius: 8px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.card h2 {
    font-size: 1.1rem;
    color: #333;
    margin-bottom: 16px;
    padding-bottom: 8px;
    border-bottom: 1px solid #eee;
}

.info-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 16px;
}

.info-item {
    display: flex;
    flex-direction: column;
    gap: 4px;
}

.info-item label {
    font-size: 12px;
    color: #999;
}

.info-item span {
    font-size: 14px;
    color: #333;
}

.status-badge {
    display: inline-block;
    padding: 4px 8px;
    border-radius: 4px;
    font-size: 12px;
    width: fit-content;
}

.status-badge.pending { background: #e9ecef; color: #495057; }
.status-badge.running { background: #cce5ff; color: #004085; }
.status-badge.completed { background: #d4edda; color: #155724; }
.status-badge.cancelled { background: #fff3cd; color: #856404; }
.status-badge.failed { background: #f8d7da; color: #721c24; }

.mode-badge {
    display: inline-block;
    padding: 4px 8px;
    border-radius: 4px;
    font-size: 12px;
    width: fit-content;
}

.mode-badge.concurrent { background: #e3f2fd; color: #1565c0; }
.mode-badge.sequential { background: #fff3e0; color: #e65100; }

.script-list, .server-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.script-item, .server-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px;
    background: #f8f9fa;
    border-radius: 6px;
}

.script-index {
    background: #667eea;
    color: white;
    padding: 2px 8px;
    border-radius: 12px;
    font-size: 12px;
    font-weight: 500;
}

.script-id {
    font-family: monospace;
    color: #666;
}

.script-name {
    color: #333;
    font-weight: 500;
}

.script-loading {
    color: #999;
    font-style: italic;
}

.server-name {
    flex: 1;
    font-weight: 500;
}

.server-status {
    padding: 2px 8px;
    border-radius: 4px;
    font-size: 12px;
}

.server-status.pending { background: #e9ecef; color: #495057; }
.server-status.running { background: #cce5ff; color: #004085; }
.server-status.success { background: #d4edda; color: #155724; }
.server-status.failed { background: #f8d7da; color: #721c24; }

.server-error {
    color: #dc3545;
    font-size: 12px;
    max-width: 300px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.empty-message {
    color: #999;
    font-style: italic;
    padding: 12px;
    text-align: center;
}

.loading, .error {
    text-align: center;
    padding: 48px;
}

.error {
    color: #dc3545;
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
