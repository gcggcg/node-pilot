<template>
    <div class="page">
        <div class="header">
            <h1>{{ isEdit ? '编辑服务器' : '添加服务器' }}</h1>
        </div>

        <form @submit.prevent="handleSubmit" class="form">
            <div class="form-group">
                <label>名称</label>
                <input v-model="form.name" type="text" required placeholder="例如: web-01" />
            </div>

            <div class="form-group">
                <label>IP地址</label>
                <input v-model="form.host" type="text" required placeholder="例如: 192.168.1.10" />
            </div>

            <div class="form-group">
                <label>端口</label>
                <input v-model.number="form.port" type="number" required min="1" max="65535" />
            </div>

            <div class="form-group">
                <label>用户名</label>
                <input v-model="form.username" type="text" required placeholder="SSH用户名" />
            </div>

            <div class="form-group">
                <label>密码</label>
                <input v-model="form.password" type="password" required placeholder="SSH密码" />
            </div>

            <div class="actions">
                <button type="submit" class="btn btn-primary" :disabled="loading">
                    {{ loading ? '保存中...' : '保存' }}
                </button>
                <router-link to="/servers" class="btn btn-secondary">取消</router-link>
            </div>
        </form>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useServerStore } from '@/stores/server';
import { serverApi } from '@/api';
import type { ServerForm } from '@/types';

const router = useRouter();
const route = useRoute();
const store = useServerStore();

const form = ref<ServerForm>({
    name: '',
    host: '',
    port: 22,
    username: '',
    password: ''
});

const loading = ref(false);
const isEdit = computed(() => !!route.params.id);

onMounted(async () => {
    if (isEdit.value) {
        const id = Number(route.params.id);
        try {
            const server = await serverApi.get(id);
            form.value = {
                name: server.name,
                host: server.host,
                port: server.port,
                username: server.username,
                password: ''
            };
        } catch (e) {
            alert('加载服务器失败');
            router.push('/servers');
        }
    }
});

async function handleSubmit() {
    loading.value = true;
    try {
        if (isEdit.value) {
            await store.updateServer(Number(route.params.id), form.value);
        } else {
            await store.createServer(form.value);
        }
        router.push('/servers');
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
    max-width: 600px;
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

.form-group input {
    width: 100%;
    padding: 10px 12px;
    border: 1px solid #ddd;
    border-radius: 6px;
    font-size: 14px;
}

.form-group input:focus {
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
