<template>
    <div class="output-panel">
        <div class="tabs">
            <button
                v-for="server in servers"
                :key="server.id"
                :class="['tab', { active: activeTab === server.id, [server.status]: true }]"
                @click="activeTab = server.id"
            >
                {{ server.name }}
                <span v-if="server.error" class="error-indicator">!</span>
            </button>
        </div>
        <div class="terminal" ref="terminalRef">
            <pre v-if="currentOutput">{{ currentOutput }}</pre>
            <div v-else-if="currentError" class="error-output">{{ currentError }}</div>
            <div v-else class="empty">等待输出...</div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue';

interface ServerStatus {
    id: number;
    name: string;
    status: string;
    error?: string;
}

const props = defineProps<{
    outputs: Map<number, string>;
    servers: ServerStatus[];
}>();

const activeTab = ref<number>(0);
const terminalRef = ref<HTMLElement | null>(null);

const currentOutput = computed(() => {
    if (activeTab.value && props.outputs.has(activeTab.value)) {
        return props.outputs.get(activeTab.value);
    }
    return '';
});

const currentError = computed(() => {
    const server = props.servers.find(s => s.id === activeTab.value);
    if (server && server.status === 'failed' && server.error) {
        return `❌ 执行失败:\n${server.error}`;
    }
    return '';
});

watch(() => props.outputs.size, () => {
    nextTick(() => {
        if (terminalRef.value) {
            terminalRef.value.scrollTop = terminalRef.value.scrollHeight;
        }
    });
});

watch(() => props.servers.length, (len) => {
    if (len > 0 && !activeTab.value) {
        activeTab.value = props.servers[0].id;
    }
});
</script>

<style scoped>
.output-panel {
    background: #1e1e1e;
    border-radius: 8px;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    height: 100%;
}

.tabs {
    display: flex;
    background: #2d2d2d;
    padding: 8px 8px 0;
    gap: 4px;
    flex-wrap: wrap;
}

.tab {
    padding: 8px 16px;
    background: transparent;
    border: none;
    color: #888;
    cursor: pointer;
    border-radius: 4px 4px 0 0;
    font-size: 13px;
    transition: all 0.2s;
}

.tab:hover {
    color: #fff;
    background: #3d3d3d;
}

.tab.active {
    color: #fff;
    background: #1e1e1e;
}

.tab.success {
    border-left: 3px solid #4caf50;
}

.tab.failed {
    border-left: 3px solid #f44336;
}

.tab.running {
    border-left: 3px solid #2196f3;
}

.error-indicator {
    background: #f44336;
    color: white;
    border-radius: 50%;
    width: 16px;
    height: 16px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    font-size: 11px;
    font-weight: bold;
    margin-left: 4px;
}

.error-output {
    color: #f44336;
    white-space: pre-wrap;
    word-break: break-all;
}

.terminal {
    flex: 1;
    padding: 16px;
    overflow-y: auto;
    font-family: 'Courier New', monospace;
    font-size: 13px;
    line-height: 1.5;
    color: #ddd;
    min-height: 300px;
}

.terminal pre {
    margin: 0;
    white-space: pre-wrap;
    word-break: break-all;
}

.empty {
    color: #666;
    font-style: italic;
}
</style>
