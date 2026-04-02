# 任务管理改造 - TaskList + TaskForm + Router

## 任务描述

改造前端任务管理页面：
1. 修改 TaskList.vue - 添加执行、编辑按钮
2. 新建 TaskForm.vue - 任务表单页面
3. 更新路由 - 添加 /tasks/new, /tasks/:id/edit

## 详细说明

### 1. TaskList.vue 改造

修改 `frontend/src/views/TaskList.vue`：

**添加操作按钮：**
```vue
<td class="actions">
    <!-- 手动执行按钮（仅 pending 状态显示） -->
    <button
        v-if="task.status === 'pending'"
        @click="executeTask(task.id)"
        class="btn btn-small btn-primary"
    >执行</button>
    
    <!-- 编辑按钮（仅 pending 状态显示） -->
    <router-link 
        v-if="task.status === 'pending'"
        :to="`/tasks/${task.id}/edit`" 
        class="btn btn-small"
    >编辑</router-link>
    
    <!-- 输出按钮 -->
    <router-link :to="`/tasks/${task.id}/output`" class="btn btn-small">输出</router-link>
    
    <!-- 取消按钮（仅 running 状态显示） -->
    <button
        v-if="task.status === 'running'"
        @click="cancelTask(task.id)"
        class="btn btn-small btn-danger"
    >取消</button>
</td>
```

**修改新建按钮文字：**
```vue
<button @click="router.push('/tasks/new')" class="btn btn-primary">+ 新建任务</button>
```

**添加执行方法：**
```typescript
async function executeTask(id: number) {
    if (confirm('确定要执行此任务吗？')) {
        try {
            await store.executeTask(id);
        } catch (e: any) {
            alert(e.message || '执行任务失败');
        }
    }
}
```

### 2. 新建 TaskForm.vue

创建 `frontend/src/views/TaskForm.vue`：

参考 `FileForm.vue` 的结构，实现：
- 任务名称输入
- 脚本选择下拉框
- 服务器选择（支持搜索和多选）
- 保存按钮
- 取消按钮

**关键差异：**
- 任务不需要选择文件
- 任务创建后直接是 pending 状态
- 执行需要另外点击执行按钮

```vue
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useTaskStore } from '@/stores/task';
import { useScriptStore } from '@/stores/script';
import { useServerStore } from '@/stores/server';

const router = useRouter();
const route = useRoute();
const store = useTaskStore();
const scriptStore = useScriptStore();
const serverStore = useServerStore();

const form = ref({
    name: '',
    script_id: '' as number | '',
    server_ids: [] as number[]
});

const loading = ref(false);
const isEdit = computed(() => !!route.params.id);

const isFormValid = computed(() => {
    return form.value.name.trim() && 
           form.value.script_id && 
           form.value.server_ids.length > 0;
});

onMounted(async () => {
    await Promise.all([
        scriptStore.fetchScripts(),
        serverStore.fetchServers()
    ]);
    
    if (isEdit.value) {
        const id = Number(route.params.id);
        try {
            const res = await store.fetchTaskDetail(id);
            form.value.name = res.task.name;
            form.value.script_id = res.task.script_id;
            form.value.server_ids = res.servers?.map((s: any) => s.server_id) || [];
        } catch (e) {
            alert('加载任务失败');
            router.push('/tasks');
        }
    }
});

async function handleSubmit() {
    if (!isFormValid.value) return;
    
    loading.value = true;
    try {
        if (isEdit.value) {
            await store.updateTask(Number(route.params.id), {
                name: form.value.name,
                script_id: Number(form.value.script_id),
                server_ids: form.value.server_ids
            });
        } else {
            await store.createTask({
                name: form.value.name,
                script_id: Number(form.value.script_id),
                server_ids: form.value.server_ids
            });
        }
        router.push('/tasks');
    } catch (e: any) {
        alert(e.message || '保存失败');
    } finally {
        loading.value = false;
    }
}
</script>
```

### 3. 更新 Router

修改 `frontend/src/router/index.ts`：

```typescript
{
    path: '/tasks',
    name: 'tasks',
    component: () => import('@/views/TaskList.vue'),
    meta: { requiresAuth: true }
},
{
    path: '/tasks/new',
    name: 'task-new',
    component: () => import('@/views/TaskForm.vue'),
    meta: { requiresAuth: true }
},
{
    path: '/tasks/:id/edit',
    name: 'task-edit',
    component: () => import('@/views/TaskForm.vue'),
    meta: { requiresAuth: true }
},
{
    path: '/tasks/:id/output',
    name: 'task-output',
    component: () => import('@/views/TaskOutput.vue'),
    meta: { requiresAuth: true }
},
```

### 4. Store 添加 fetchTaskDetail

在 `frontend/src/stores/task.ts` 中添加：

```typescript
async function fetchTaskDetail(id: number) {
    loading.value = true;
    error.value = null;
    try {
        const res = await taskApi.get(id);
        currentTask.value = res.task;
        return res;
    } catch (e: any) {
        error.value = e.message;
        throw e;
    } finally {
        loading.value = false;
    }
}
```

## 输入

- `frontend/src/views/TaskList.vue`
- `frontend/src/views/TaskForm.vue` (新建)
- `frontend/src/router/index.ts`
- `frontend/src/stores/task.ts`

## 输出

- 修改后的 `TaskList.vue`
- 新建的 `TaskForm.vue`
- 修改后的 `router/index.ts`
- 修改后的 `stores/task.ts`

## 依赖

- 04-frontend-types-api
- 05-frontend-store

## 验收标准

- [ ] TaskList 显示执行按钮（仅 pending 任务）
- [ ] TaskList 显示编辑按钮（仅 pending 任务）
- [ ] TaskList 正确调用 executeTask
- [ ] TaskForm 正确创建和更新任务
- [ ] 路由正确配置
- [ ] 编辑模式正确加载任务数据
