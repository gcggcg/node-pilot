<template>
    <div class="page">
        <div class="header">
            <h1>{{ isEdit ? '编辑任务' : '新建任务' }}</h1>
        </div>

        <form @submit.prevent="handleSubmit" class="form">
            <div class="form-group">
                <label>任务名称</label>
                <input v-model="form.name" type="text" required placeholder="例如: 批量部署-1" />
            </div>

            <div class="form-group">
                <label>脚本 (支持批量)</label>
                <ScriptSelector 
                    v-model="selectedScriptIds"
                    multiple
                    placeholder="选择要执行的脚本（可多选）"
                />
            </div>

            <div class="form-group">
                <label>服务器</label>
                <div class="server-select">
                    <div class="selected-servers" v-if="form.server_ids.length > 0">
                        <span v-for="sid in form.server_ids" :key="sid" class="server-tag">
                            {{ getServerName(sid) }}
                            <button type="button" @click="removeServer(sid)">&times;</button>
                        </span>
                    </div>
                    <div v-if="form.server_ids.length === 0" class="no-servers">
                        请选择目标服务器
                    </div>
                    <div class="server-dropdown">
                        <div class="server-search">
                            <input 
                                v-model="serverSearch" 
                                type="text" 
                                placeholder="搜索服务器..." 
                                @focus="showServerDropdown = true"
                            />
                        </div>
                        <div v-if="showServerDropdown" class="server-list-dropdown">
                            <div 
                                v-for="server in filteredServers" 
                                :key="server.id" 
                                class="server-option"
                                :class="{ selected: form.server_ids.includes(server.id) }"
                                @click="toggleServer(server.id)"
                            >
                                <span class="server-name">{{ server.name }}</span>
                                <span class="server-info">{{ server.host }}:{{ server.port }}</span>
                            </div>
                            <div v-if="filteredServers.length === 0" class="no-results">
                                未找到服务器
                            </div>
                        </div>
                    </div>
                </div>
                <div class="hint-text">最多选择 10 台服务器</div>
            </div>

            <div class="actions">
                <button 
                    type="submit" 
                    class="btn btn-primary" 
                    :disabled="loading || !isFormValid"
                >
                    {{ loading ? '保存中...' : '保存' }}
                </button>
                <router-link to="/tasks" class="btn btn-secondary">取消</router-link>
            </div>
        </form>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useTaskStore } from '@/stores/task';
import { useScriptStore } from '@/stores/script';
import { useServerStore } from '@/stores/server';
import ScriptSelector from '@/components/ScriptSelector.vue';
import type { Server } from '@/types';

const router = useRouter();
const route = useRoute();
const store = useTaskStore();
const scriptStore = useScriptStore();
const serverStore = useServerStore();

const form = ref({
    name: '',
    script_id: '' as number | '',
    script_ids: '',
    server_ids: [] as number[]
});

const selectedScriptIds = ref<number[]>([]);
const loading = ref(false);
const isEdit = computed(() => !!route.params.id);
const serverSearch = ref('');
const showServerDropdown = ref(false);
const servers = ref<Server[]>([]);

const isFormValid = computed(() => {
    const hasScripts = selectedScriptIds.value.length > 0 || form.value.script_id;
    return form.value.name.trim() && 
           hasScripts && 
           form.value.server_ids.length > 0;
});

const scriptIdsString = computed(() => {
    return selectedScriptIds.value.join(',');
});

const filteredServers = computed(() => {
    if (!serverSearch.value) return servers.value;
    const search = serverSearch.value.toLowerCase();
    return servers.value.filter(s => 
        s.name.toLowerCase().includes(search) || 
        s.host.toLowerCase().includes(search)
    );
});

const MAX_SERVERS = 10;

function getServerName(id: number): string {
    const server = servers.value.find(s => s.id === id);
    return server ? `${server.name} (${server.host}:${server.port})` : `服务器 ${id}`;
}

function toggleServer(id: number) {
    const index = form.value.server_ids.indexOf(id);
    if (index >= 0) {
        form.value.server_ids.splice(index, 1);
    } else {
        if (form.value.server_ids.length >= MAX_SERVERS) {
            alert(`最多只能选择 ${MAX_SERVERS} 台服务器`);
            return;
        }
        form.value.server_ids.push(id);
    }
}

function removeServer(id: number) {
    const index = form.value.server_ids.indexOf(id);
    if (index >= 0) {
        form.value.server_ids.splice(index, 1);
    }
}

async function loadServers() {
    try {
        await serverStore.fetchServers(1, 100);
        servers.value = serverStore.servers;
    } catch (e) {
        console.error('Failed to load servers:', e);
    }
}

onMounted(async () => {
    await Promise.all([
        scriptStore.fetchScripts(),
        loadServers()
    ]);
    
    if (isEdit.value) {
        const id = Number(route.params.id);
        try {
            const res = await store.fetchTaskDetail(id);
            form.value.name = res.task.name;
            form.value.script_id = res.task.script_id;
            form.value.script_ids = res.task.script_ids || '';
            if (res.task.script_ids) {
                selectedScriptIds.value = res.task.script_ids.split(',').map((id: string) => parseInt(id.trim()));
            } else if (res.task.script_id) {
                selectedScriptIds.value = [res.task.script_id];
            }
            form.value.server_ids = res.servers?.map((s: any) => s.server_id) || [];
        } catch (e) {
            alert('加载任务失败');
            router.push('/tasks');
        }
    }
    
    document.addEventListener('click', (e) => {
        const target = e.target as HTMLElement;
        if (!target.closest('.server-select')) {
            showServerDropdown.value = false;
        }
    });
});

async function handleSubmit() {
    if (!isFormValid.value) return;
    
    loading.value = true;
    try {
        const payload: any = {
            name: form.value.name,
            server_ids: form.value.server_ids
        };
        
        if (selectedScriptIds.value.length > 0) {
            payload.script_ids = scriptIdsString.value;
        } else if (form.value.script_id) {
            payload.script_id = Number(form.value.script_id);
        }
        
        if (isEdit.value) {
            await store.updateTask(Number(route.params.id), payload);
        } else {
            await store.createTask(payload);
        }
        router.push('/tasks');
    } catch (e: any) {
        alert(e.message || '保存失败');
    } finally {
        loading.value = false;
    }
}
</script>

<style scoped>
.page {
    padding: 24px;
    max-width: 800px;
}

.header {
    margin-bottom: 24px;
}

h1 {
    font-size: 1.5rem;
    color: #333;
}

.form {
    background: white;
    padding: 24px;
    border-radius: 8px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.form-group {
    margin-bottom: 20px;
}

.form-group label {
    display: block;
    margin-bottom: 6px;
    font-weight: 500;
    color: #555;
}

.form-group input[type="text"],
.form-group select {
    width: 100%;
    padding: 10px 12px;
    border: 1px solid #ddd;
    border-radius: 6px;
    font-size: 14px;
}

.form-group input[type="text"]:focus,
.form-group select:focus {
    outline: none;
    border-color: #667eea;
}

.hint-text {
    color: #999;
    font-size: 12px;
    margin-top: 4px;
}

.server-select {
    position: relative;
}

.selected-servers {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-bottom: 8px;
}

.server-tag {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 4px 10px;
    background: #e9ecef;
    border-radius: 16px;
    font-size: 13px;
    color: #495057;
}

.server-tag button {
    background: none;
    border: none;
    font-size: 14px;
    color: #999;
    cursor: pointer;
    padding: 0;
    line-height: 1;
}

.server-tag button:hover {
    color: #dc3545;
}

.no-servers {
    padding: 12px;
    background: #f8f9fa;
    border-radius: 6px;
    color: #999;
    font-size: 14px;
    text-align: center;
}

.server-dropdown {
    position: relative;
}

.server-search input {
    width: 100%;
    padding: 8px 12px;
    border: 1px solid #ddd;
    border-radius: 6px;
    font-size: 14px;
}

.server-search input:focus {
    outline: none;
    border-color: #667eea;
}

.server-list-dropdown {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    background: white;
    border: 1px solid #ddd;
    border-radius: 6px;
    margin-top: 4px;
    max-height: 240px;
    overflow-y: auto;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    z-index: 100;
}

.server-option {
    padding: 10px 12px;
    cursor: pointer;
    display: flex;
    justify-content: space-between;
    border-bottom: 1px solid #eee;
}

.server-option:last-child {
    border-bottom: none;
}

.server-option:hover {
    background: #f8f9fa;
}

.server-option.selected {
    background: #e9ecef;
}

.server-name {
    font-weight: 500;
    color: #333;
}

.server-info {
    font-size: 12px;
    color: #999;
}

.no-results {
    padding: 16px;
    text-align: center;
    color: #999;
}

.actions {
    display: flex;
    gap: 12px;
    margin-top: 24px;
}

.btn {
    padding: 10px 20px;
    border: none;
    border-radius: 6px;
    cursor: pointer;
    font-size: 14px;
}

.btn-primary {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
}

.btn-primary:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.btn-secondary {
    background: #e9ecef;
    color: #495057;
}
</style>
