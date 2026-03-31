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
                <textarea v-model="form.content" rows="12" required placeholder="#!/bin/bash&#10;echo 'Hello World'" class="code"></textarea>
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
</style>
