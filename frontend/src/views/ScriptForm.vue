<template>
    <div class="page">
        <div class="header">
            <h1>{{ isEdit ? '编辑脚本' : '创建脚本' }}</h1>
        </div>

        <form @submit.prevent="handleSubmit" class="form">
            <div class="form-group">
                <label>名称</label>
                <input v-model="form.name" type="text" required placeholder="例如: deploy-app.sh" />
            </div>

            <div class="form-group">
                <label>描述</label>
                <textarea v-model="form.description" rows="2" placeholder="可选描述"></textarea>
            </div>

            <div class="form-group">
                <label>目标路径</label>
                <input v-model="form.target_path" type="text" required placeholder="例如: /opt/scripts/deploy.sh" />
            </div>

            <div class="form-group">
                <label>脚本内容</label>
                
                <div class="mode-toggle">
                    <button 
                        type="button" 
                        :class="['mode-btn', { active: inputMode === 'manual' }]"
                        @click="inputMode = 'manual'"
                    >
                        手动输入
                    </button>
                    <button 
                        type="button" 
                        :class="['mode-btn', { active: inputMode === 'upload' }]"
                        @click="inputMode = 'upload'"
                    >
                        文件上传
                    </button>
                </div>

                <div v-if="inputMode === 'upload'" class="upload-area" 
                    @dragover.prevent="onDragOver"
                    @dragleave.prevent="onDragLeave"
                    @drop.prevent="onDrop"
                >
                    <input 
                        type="file" 
                        ref="fileInput" 
                        accept=".txt,.sh,.py,.js,.sql"
                        @change="onFileSelect"
                        style="display: none"
                    />
                    <div class="upload-placeholder" @click="triggerFileSelect">
                        <span class="upload-icon">📁</span>
                        <p>拖拽文件到此处，或 <span class="link">点击选择</span></p>
                        <p class="hint">支持 .txt, .sh, .py, .js, .sql 文件，不超过 5MB</p>
                    </div>
                </div>

                <textarea v-else v-model="form.content" rows="12" required placeholder="#!/bin/bash&#10;echo 'Hello World'" class="code"></textarea>
            </div>

            <div class="actions">
                <button type="submit" class="btn btn-primary" :disabled="loading">
                    {{ loading ? '保存中...' : '保存' }}
                </button>
                <router-link to="/scripts" class="btn btn-secondary">取消</router-link>
            </div>
        </form>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useScriptStore } from '@/stores/script';
import { scriptApi } from '@/api';
import type { ScriptForm } from '@/types';

const router = useRouter();
const route = useRoute();
const store = useScriptStore();

const form = ref<ScriptForm>({
    name: '',
    description: '',
    content: '#!/bin/bash\n',
    target_path: ''
});

const loading = ref(false);
const isEdit = computed(() => !!route.params.id);
const inputMode = ref<'manual' | 'upload'>('manual');
const fileInput = ref<HTMLInputElement | null>(null);
const isUploading = ref(false);
const uploadError = ref<string | null>(null);
const manualContent = ref(form.value.content);
const uploadedContent = ref('');

onMounted(async () => {
    if (isEdit.value) {
        const id = Number(route.params.id);
        try {
            const script = await scriptApi.get(id);
            form.value = {
                name: script.name,
                description: script.description || '',
                content: script.content,
                target_path: script.target_path
            };
        } catch (e) {
            alert('加载脚本失败');
            router.push('/scripts');
        }
    }
});

function triggerFileSelect() {
    fileInput.value?.click();
}

function onDragOver(event: DragEvent) {
    (event.currentTarget as HTMLElement).classList.add('drag-over');
}

function onDragLeave(event: DragEvent) {
    (event.currentTarget as HTMLElement).classList.remove('drag-over');
}

function onDrop(event: DragEvent) {
    (event.currentTarget as HTMLElement).classList.remove('drag-over');
    const file = event.dataTransfer?.files?.[0];
    if (file) {
        handleFile(file);
    }
}

function onFileSelect(event: Event) {
    const target = event.target as HTMLInputElement;
    const file = target.files?.[0];
    if (file) {
        handleFile(file);
    }
}

const ALLOWED_EXTENSIONS = ['.txt', '.sh', '.py', '.js', '.sql'];
const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB

function validateFile(file: File): string | null {
    const ext = '.' + file.name.split('.').pop()?.toLowerCase();
    if (!ALLOWED_EXTENSIONS.includes(ext)) {
        return `不支持的文件格式，请上传 ${ALLOWED_EXTENSIONS.join(', ')} 文件`;
    }
    if (file.size > MAX_FILE_SIZE) {
        return '文件大小超过5MB限制';
    }
    return null;
}

function handleFile(file: File) {
    uploadError.value = null;

    // 验证文件
    const error = validateFile(file);
    if (error) {
        uploadError.value = error;
        return;
    }

    // 读取文件内容
    isUploading.value = true;
    const reader = new FileReader();

    reader.onload = (e) => {
        const content = e.target?.result as string;
        uploadedContent.value = content;
        form.value.content = content; // 直接填充到表单
        isUploading.value = false;
    };

    reader.onerror = () => {
        uploadError.value = '文件读取失败，请重试';
        isUploading.value = false;
    };

    reader.readAsText(file);
}

async function handleSubmit() {
    loading.value = true;
    try {
        if (isEdit.value) {
            await store.updateScript(Number(route.params.id), form.value);
        } else {
            await store.createScript(form.value);
        }
        router.push('/scripts');
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
    margin-bottom: 16px;
}

.form-group label {
    display: block;
    margin-bottom: 6px;
    font-weight: 500;
    color: #555;
}

.form-group input,
.form-group textarea {
    width: 100%;
    padding: 10px 12px;
    border: 1px solid #ddd;
    border-radius: 6px;
    font-size: 14px;
    font-family: inherit;
}

.form-group textarea.code {
    font-family: 'Courier New', monospace;
    resize: vertical;
}

.form-group input:focus,
.form-group textarea:focus {
    outline: none;
    border-color: #667eea;
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
}

.btn-secondary {
    background: #e9ecef;
    color: #495057;
}

.mode-toggle {
    display: flex;
    gap: 8px;
    margin-bottom: 12px;
}

.mode-btn {
    padding: 8px 16px;
    border: 1px solid #ddd;
    border-radius: 6px;
    background: white;
    cursor: pointer;
    font-size: 14px;
    color: #666;
    transition: all 0.2s;
}

.mode-btn.active {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    border-color: #667eea;
}

.upload-area {
    border: 2px dashed #ddd;
    border-radius: 8px;
    padding: 32px;
    text-align: center;
    transition: all 0.2s;
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

.upload-placeholder .link {
    color: #667eea;
    text-decoration: underline;
}

.upload-placeholder .hint {
    font-size: 12px;
    color: #999;
}

.drag-over {
    border-color: #667eea !important;
    background: #f8f7ff;
}
</style>
