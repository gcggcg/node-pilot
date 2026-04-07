<template>
    <div class="script-selector">
        <div class="selected-scripts" v-if="selectedIds.length > 0">
            <span v-for="sid in selectedIds" :key="sid" class="script-tag">
                {{ getScriptName(sid) }}
                <button type="button" @click="removeScript(sid)">&times;</button>
            </span>
        </div>
        <div class="script-dropdown">
            <input 
                v-model="search" 
                type="text" 
                :placeholder="placeholder || '搜索脚本...'"
                @focus="showDropdown = true"
                class="script-search-input"
            />
            <div v-if="showDropdown" class="script-list-dropdown">
                <div 
                    v-for="script in filteredScripts" 
                    :key="script.id" 
                    class="script-option"
                    :class="{ selected: selectedIds.includes(script.id) }"
                    @click="toggleScript(script.id)"
                >
                    <div class="script-info">
                        <span class="script-name">{{ script.name }}</span>
                        <span class="script-path">{{ script.target_path }}</span>
                    </div>
                    <span v-if="selectedIds.includes(script.id)" class="check-mark">✓</span>
                </div>
                <div v-if="filteredScripts.length === 0" class="no-results">
                    未找到脚本
                </div>
            </div>
        </div>
        <div class="selected-info">
            <span v-if="selectedIds.length > 0">已选择 {{ selectedIds.length }} 个脚本</span>
            <span v-else class="no-selection">请选择脚本</span>
            <button 
                v-if="selectedIds.length > 0" 
                type="button" 
                class="clear-btn"
                @click="clearAll"
            >
                清空
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import { useScriptStore } from '@/stores/script';
import type { Script } from '@/types';

interface Props {
    modelValue: number[];
    multiple?: boolean;
    placeholder?: string;
}

interface Emits {
    (e: 'update:modelValue', value: number[]): void;
}

const props = withDefaults(defineProps<Props>(), {
    multiple: true
});

const emit = defineEmits<Emits>();

const scriptStore = useScriptStore();
const selectedIds = ref<number[]>([]);
const search = ref('');
const showDropdown = ref(false);
const scripts = ref<Script[]>([]);

const filteredScripts = computed(() => {
    if (!search.value) return scripts.value;
    const s = search.value.toLowerCase();
    return scripts.value.filter(sc => 
        sc.name.toLowerCase().includes(s) || 
        sc.target_path.toLowerCase().includes(s)
    );
});

watch(() => props.modelValue, (newVal) => {
    selectedIds.value = [...newVal];
}, { immediate: true });

function getScriptName(id: number): string {
    const script = scripts.value.find(s => s.id === id);
    return script ? script.name : `脚本 ${id}`;
}

function toggleScript(id: number) {
    const index = selectedIds.value.indexOf(id);
    if (index >= 0) {
        selectedIds.value.splice(index, 1);
    } else {
        if (props.multiple) {
            selectedIds.value.push(id);
        } else {
            selectedIds.value = [id];
            showDropdown.value = false;
        }
    }
    emit('update:modelValue', [...selectedIds.value]);
}

function removeScript(id: number) {
    const index = selectedIds.value.indexOf(id);
    if (index >= 0) {
        selectedIds.value.splice(index, 1);
        emit('update:modelValue', [...selectedIds.value]);
    }
}

function clearAll() {
    selectedIds.value = [];
    emit('update:modelValue', []);
}

async function loadScripts() {
    try {
        await scriptStore.fetchScripts(1, 100);
        scripts.value = scriptStore.scripts;
    } catch (e) {
        console.error('Failed to load scripts:', e);
    }
}

onMounted(() => {
    loadScripts();
    document.addEventListener('click', (e) => {
        const target = e.target as HTMLElement;
        if (!target.closest('.script-dropdown')) {
            showDropdown.value = false;
        }
    });
});
</script>

<style scoped>
.script-selector {
    width: 100%;
}

.selected-scripts {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-bottom: 8px;
}

.script-tag {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 4px 10px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 16px;
    font-size: 13px;
    color: white;
}

.script-tag button {
    background: none;
    border: none;
    font-size: 14px;
    color: rgba(255,255,255,0.8);
    cursor: pointer;
    padding: 0;
    line-height: 1;
}

.script-tag button:hover {
    color: white;
}

.script-dropdown {
    position: relative;
}

.script-search-input {
    width: 100%;
    padding: 10px 12px;
    border: 1px solid #ddd;
    border-radius: 6px;
    font-size: 14px;
}

.script-search-input:focus {
    outline: none;
    border-color: #667eea;
}

.script-list-dropdown {
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

.script-option {
    padding: 10px 12px;
    cursor: pointer;
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid #eee;
}

.script-option:last-child {
    border-bottom: none;
}

.script-option:hover {
    background: #f8f9fa;
}

.script-option.selected {
    background: #e9ecef;
}

.script-info {
    display: flex;
    flex-direction: column;
    gap: 2px;
}

.script-name {
    font-weight: 500;
    color: #333;
}

.script-path {
    font-size: 12px;
    color: #999;
    font-family: monospace;
}

.check-mark {
    color: #667eea;
    font-weight: bold;
}

.no-results {
    padding: 16px;
    text-align: center;
    color: #999;
}

.selected-info {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 6px;
    font-size: 12px;
    color: #666;
}

.no-selection {
    color: #999;
}

.clear-btn {
    background: none;
    border: none;
    color: #dc3545;
    cursor: pointer;
    font-size: 12px;
}

.clear-btn:hover {
    text-decoration: underline;
}
</style>
