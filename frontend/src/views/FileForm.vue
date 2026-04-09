<template>
    <div class="page">
        <div class="header">
            <h1>{{ isEdit ? '编辑上传配置' : '创建上传配置' }}</h1>
        </div>

        <form @submit.prevent="handleSubmit" class="form">
            <div class="form-group">
                <label>配置名称</label>
                <input v-model="form.name" type="text" required placeholder="例如: 部署配置文件" />
            </div>

            <!-- 本地文件选择 (仅新增模式) -->
            <div v-if="!isEdit" class="form-group">
                <label>选择文件</label>
                <input 
                    type="file" 
                    ref="fileInput" 
                    multiple
                    accept=".tar,.tar.gz,.tgz,.sh,.zip,.conf,.txt,.json,.yml,.xml,.log,.sql"
                    @change="onFileSelect"
                    style="display: none"
                />
                
                <div class="upload-area" 
                    @dragover.prevent="onDragOver"
                    @dragleave.prevent="onDragLeave"
                    @drop.prevent="onDrop"
                    :class="{ 'drag-over': isDragging }"
                >
                    <div v-if="form.files.length === 0" class="upload-placeholder" @click="triggerFileSelect">
                        <span class="upload-icon">📁</span>
                        <p>点击选择文件或拖拽文件到此处</p>
                        <p class="hint">支持 .tar, .tar.gz, .tgz, .sh, .zip, .conf, .txt, .json, .yml, .xml, .log, .sql 文件，单文件不超过 500MB，最多 20 个文件</p>
                    </div>
                    
                    <div v-else class="file-list">
                        <div v-for="(file, index) in form.files" :key="index" class="file-item">
                            <span class="file-name">{{ file.name }}</span>
                            <span class="file-size">{{ formatSize(file.size) }}</span>
                            <button type="button" class="btn-remove" @click="removeFile(index)">&times;</button>
                        </div>
                        <div class="file-actions">
                            <button type="button" class="btn-add" @click="triggerFileSelect">+ 添加更多文件</button>
                        </div>
                    </div>
                </div>
                
                <div v-if="fileError" class="error-text">{{ fileError }}</div>
            </div>

            <!-- 目标服务器选择 -->
            <div class="form-group">
                <label>目标服务器</label>
                <div class="server-select">
                    <div class="selected-servers" v-if="form.serverIds.length > 0">
                        <span v-for="sid in form.serverIds" :key="sid" class="server-tag">
                            {{ getServerName(sid) }}
                            <button type="button" @click="removeServer(sid)">&times;</button>
                        </span>
                    </div>
                    <div v-if="form.serverIds.length === 0" class="no-servers">
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
                                :class="{ selected: form.serverIds.includes(server.id) }"
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

            <!-- 远程路径配置 -->
            <div class="form-group">
                <label>远程路径</label>
                <input 
                    v-model="form.remotePath" 
                    type="text" 
                    required
                    placeholder="例如: /opt/configs/"
                    :class="{ 'invalid': form.remotePath && !isPathValid }"
                />
                <div v-if="form.remotePath && !isPathValid" class="error-text">
                    远程路径必须以 / 开头
                </div>
            </div>

            <div class="actions">
                <button 
                    type="submit" 
                    class="btn btn-primary" 
                    :disabled="loading || !isFormValid"
                >
                    {{ loading ? '保存中...' : '保存' }}
                </button>
                <router-link to="/files" class="btn btn-secondary">取消</router-link>
            </div>
        </form>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useFileUploadStore } from '@/stores/fileupload';
import { serverApi, fileUploadApi } from '@/api';
import type { Server } from '@/types';

const router = useRouter();
const route = useRoute();
const store = useFileUploadStore();

const form = ref({
    name: '',
    files: [] as File[],
    serverIds: [] as number[],
    remotePath: ''
});

const loading = ref(false);
const isEdit = computed(() => !!route.params.id);
const fileInput = ref<HTMLInputElement | null>(null);
const fileError = ref<string | null>(null);
const isDragging = ref(false);
const serverSearch = ref('');
const showServerDropdown = ref(false);
const servers = ref<Server[]>([]);

const isPathValid = computed(() => /^\//.test(form.value.remotePath));

const isFormValid = computed(() => {
    if (!form.value.name.trim()) return false;
    if (!isPathValid.value) return false;
    if (form.value.serverIds.length === 0) return false;
    if (!isEdit.value && form.value.files.length === 0) return false;
    return true;
});

const filteredServers = computed(() => {
    if (!serverSearch.value) return servers.value;
    const search = serverSearch.value.toLowerCase();
    return servers.value.filter(s => 
        s.name.toLowerCase().includes(search) || 
        s.host.toLowerCase().includes(search)
    );
});

const ALLOWED_EXTENSIONS = ['.tar', '.tar.gz', '.tgz', '.sh', '.zip', '.conf', '.txt', '.json', '.yml', '.xml', '.log', '.sql'];
const MAX_FILE_SIZE = 500 * 1024 * 1024; // 500MB
const MAX_FILES = 20;
const MAX_SERVERS = 10;

function getServerName(id: number): string {
    const server = servers.value.find(s => s.id === id);
    return server ? `${server.name} (${server.host}:${server.port})` : `服务器 ${id}`;
}

function formatSize(bytes: number): string {
    if (bytes < 1024) return bytes + ' B';
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
}

function triggerFileSelect() {
    fileInput.value?.click();
}

function onDragOver(_event: DragEvent) {
    isDragging.value = true;
}

function onDragLeave(_event: DragEvent) {
    isDragging.value = false;
}

function onDrop(event: DragEvent) {
    isDragging.value = false;
    const files = Array.from(event.dataTransfer?.files || []);
    handleFiles(files);
}

function onFileSelect(event: Event) {
    const target = event.target as HTMLInputElement;
    const files = Array.from(target.files || []);
    handleFiles(files);
}

function validateFile(file: File): string | null {
    const lowerName = file.name.toLowerCase();
    // 检查复合扩展名 (.tar.gz, .tar.bz2 等)
    const compoundExtensions = ['.tar.gz', '.tar.bz2', '.tar.xz', '.tar.zst'];
    for (const compExt of compoundExtensions) {
        if (lowerName.endsWith(compExt)) {
            if (file.size > MAX_FILE_SIZE) {
                return `文件超过500MB限制: ${file.name}`;
            }
            return null;
        }
    }
    // 单扩展名
    const ext = '.' + file.name.split('.').pop()?.toLowerCase();
    if (!ALLOWED_EXTENSIONS.includes(ext)) {
        return `不支持的文件格式: ${file.name}`;
    }
    if (file.size > MAX_FILE_SIZE) {
        return `文件超过500MB限制: ${file.name}`;
    }
    return null;
}

function handleFiles(files: File[]) {
    fileError.value = null;
    
    const remainingSlots = MAX_FILES - form.value.files.length;
    if (files.length > remainingSlots) {
        fileError.value = `最多只能选择 ${MAX_FILES} 个文件`;
        files = files.slice(0, remainingSlots);
    }
    
    for (const file of files) {
        const error = validateFile(file);
        if (error) {
            fileError.value = error;
            return;
        }
    }
    
    form.value.files.push(...files);
}

function removeFile(index: number) {
    form.value.files.splice(index, 1);
    fileError.value = null;
}

function toggleServer(id: number) {
    const index = form.value.serverIds.indexOf(id);
    if (index >= 0) {
        form.value.serverIds.splice(index, 1);
    } else {
        if (form.value.serverIds.length >= MAX_SERVERS) {
            alert(`最多只能选择 ${MAX_SERVERS} 台服务器`);
            return;
        }
        form.value.serverIds.push(id);
    }
}

function removeServer(id: number) {
    const index = form.value.serverIds.indexOf(id);
    if (index >= 0) {
        form.value.serverIds.splice(index, 1);
    }
}

async function loadServers() {
    try {
        const res = await serverApi.list({ page: 1, pageSize: 100 });
        servers.value = res.data;
    } catch (e) {
        console.error('Failed to load servers:', e);
    }
}

onMounted(async () => {
    await loadServers();
    
    if (isEdit.value) {
        const id = Number(route.params.id);
        try {
            await store.fetchResults(id);
            const upload = store.currentUpload;
            if (upload) {
                form.value.name = upload.name;
                form.value.remotePath = upload.remote_path;
                form.value.serverIds = upload.servers?.map(s => s.id) || [];
            }
        } catch (e) {
            alert('加载配置失败');
            router.push('/files');
        }
    }
    
    // 点击外部关闭服务器下拉框
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
        if (isEdit.value) {
            await store.updateUpload(Number(route.params.id), {
                name: form.value.name,
                localPath: '',
                remotePath: form.value.remotePath,
                serverIds: form.value.serverIds
            });
        } else {
            // 先上传文件获取路径
            const formData = new FormData();
            formData.append('name', form.value.name);
            formData.append('remotePath', form.value.remotePath);
            formData.append('serverIds', JSON.stringify(form.value.serverIds));
            
            // 上传文件
            if (form.value.files.length > 0) {
                const fileRes = await fileUploadApi.uploadFile(form.value.files[0]) as any;
                
                await store.createUpload({
                    name: form.value.name,
                    localPath: fileRes.path,
                    remotePath: form.value.remotePath,
                    serverIds: form.value.serverIds
                });
            }
        }
        router.push('/files');
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

.form-group input[type="text"] {
    width: 100%;
    padding: 10px 12px;
    border: 1px solid #ddd;
    border-radius: 6px;
    font-size: 14px;
}

.form-group input[type="text"]:focus {
    outline: none;
    border-color: #667eea;
}

.form-group input[type="text"].invalid {
    border-color: #dc3545;
}

.error-text {
    color: #dc3545;
    font-size: 12px;
    margin-top: 4px;
}

.hint-text {
    color: #999;
    font-size: 12px;
    margin-top: 4px;
}

.upload-area {
    border: 2px dashed #ddd;
    border-radius: 8px;
    padding: 24px;
    text-align: center;
    transition: all 0.2s;
}

.upload-area.drag-over {
    border-color: #667eea;
    background: #f8f7ff;
}

.upload-placeholder {
    cursor: pointer;
}

.upload-icon {
    font-size: 48px;
    display: block;
    margin-bottom: 12px;
}

.upload-placeholder p {
    margin: 8px 0;
    color: #666;
}

.upload-placeholder .hint {
    font-size: 12px;
    color: #999;
}

.file-list {
    text-align: left;
}

.file-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 8px 12px;
    background: #f8f9fa;
    border-radius: 6px;
    margin-bottom: 8px;
}

.file-name {
    flex: 1;
    font-size: 14px;
    color: #333;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.file-size {
    font-size: 12px;
    color: #999;
}

.btn-remove {
    background: none;
    border: none;
    font-size: 18px;
    color: #999;
    cursor: pointer;
    padding: 0 4px;
}

.btn-remove:hover {
    color: #dc3545;
}

.file-actions {
    margin-top: 8px;
}

.btn-add {
    background: none;
    border: 1px dashed #ddd;
    border-radius: 6px;
    padding: 8px 16px;
    color: #667eea;
    cursor: pointer;
    font-size: 14px;
    width: 100%;
}

.btn-add:hover {
    background: #f8f7ff;
    border-color: #667eea;
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
