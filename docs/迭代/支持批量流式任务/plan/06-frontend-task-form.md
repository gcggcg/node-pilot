# Frontend Task Form - Batch Script Integration

## 任务描述

修改任务表单（TaskForm.vue），集成多选脚本组件，支持创建和编辑批量任务。

## 详细说明

### 1. 分析现有 TaskForm.vue
查看 `frontend/src/views/TaskForm.vue` 了解现有表单结构。

### 2. 修改表单结构

**原表单字段：**
- 任务名称 (name)
- 脚本选择 (script_id，单选)
- 服务器选择 (server_ids，多选)

**新表单字段：**
- 任务名称 (name)
- 脚本选择 (script_ids，多选) - **使用新的 ScriptSelector 组件**
- 服务器选择 (server_ids，多选)

### 3. 修改表单数据结构

```typescript
// 任务表单数据
interface TaskFormData {
  name: string
  script_ids: string   // 逗号分隔的脚本ID列表
  server_ids: number[]
}

// 初始化默认值
const formData = reactive({
  name: '',
  script_ids: '',
  server_ids: []
})
```

### 4. 集成 ScriptSelector 组件

```vue
<template>
  <el-form :model="formData" label-width="120px">
    <!-- 任务名称 -->
    <el-form-item label="任务名称" required>
      <el-input v-model="formData.name" placeholder="输入任务名称" />
    </el-form-item>

    <!-- 批量脚本选择 -->
    <el-form-item label="选择脚本" required>
      <ScriptSelector
        v-model="selectedScriptIds"
        multiple
        placeholder="选择要执行的脚本（可多选）"
      />
    </el-form-item>

    <!-- 服务器选择 -->
    <el-form-item label="选择服务器" required>
      <ServerSelector
        v-model="formData.server_ids"
        multiple
        placeholder="选择目标服务器"
      />
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import ScriptSelector from '@/components/ScriptSelector.vue'
import ServerSelector from '@/components/ServerSelector.vue'

const selectedScriptIds = ref<number[]>([])

// 将数组转换为逗号分隔字符串
const scriptIdsString = computed(() => {
  return selectedScriptIds.value.join(',')
})

// 提交时使用
function handleSubmit() {
  const payload = {
    name: formData.name,
    script_ids: scriptIdsString.value,
    server_ids: formData.server_ids
  }
  // 调用 API
}
</script>
```

### 5. 处理编辑模式
如果支持编辑已有任务：

```typescript
// 加载任务数据
async function loadTask(taskId: number) {
  const response = await taskApi.get(taskId)
  formData.name = response.task.name
  
  // 解析 script_ids 字段
  if (response.task.script_ids) {
    selectedScriptIds.value = response.task.script_ids.split(',').map(id => parseInt(id))
  }
}
```

### 6. 任务列表显示修改
修改 `TaskList.vue` 显示批量脚本信息：

- 显示脚本数量：如 "3 个脚本"
- 或者显示脚本名称列表

### 7. API 调用修改

**CreateTask API 请求：**
```typescript
interface CreateTaskRequest {
  name: string
  script_ids: string  // "1,2,3" 格式
  server_ids: number[]
}
```

## 输入

- 现有 TaskForm.vue (frontend/src/views/TaskForm.vue)
- ScriptSelector 组件 (05-frontend-script-selector.md 产出)
- Task API (frontend/src/api/index.ts)

## 输出

- 修改后的 TaskForm.vue，支持批量脚本选择
- 可选：修改 TaskList.vue 显示脚本信息

## 依赖

- 05（ScriptSelector 组件已完成）

## 验收标准

- [ ] TaskForm 支持多选脚本
- [ ] 使用 ScriptSelector 组件
- [ ] 正确处理单脚本向后兼容（script_id 字段）
- [ ] 提交时正确构造 script_ids 字符串
- [ ] 编辑模式正确加载批量脚本数据
- [ ] 表单验证正确（至少选择一个脚本）
- [ ] 任务列表正确显示脚本信息
