# Frontend Script Selector - Multi-Select Component

## 任务描述

创建前端多选脚本组件，支持在任务表单中选择多个脚本。

## 详细说明

### 1. 分析现有 ScriptList 组件
查看 `frontend/src/views/ScriptList.vue` 和 `frontend/src/stores/script.ts` 了解现有脚本列表实现。

### 2. 创建 ScriptSelector 组件
在 `frontend/src/components/` 下创建 `ScriptSelector.vue` 组件：

**功能需求：**
- 显示所有可用脚本列表
- 支持单选和多选模式
- 显示已选脚本数量
- 支持搜索/过滤脚本

**组件接口：**
```typescript
interface Props {
  modelValue: number[]        // 选中的脚本ID数组
  multiple?: boolean = true   // 是否多选
  placeholder?: string
}

interface Emits {
  (e: 'update:modelValue', value: number[]): void
}
```

**UI 设计：**
- 下拉选择框形式
- 选中的脚本显示为标签（tag）
- 显示已选数量，如 "已选择 3 个脚本"
- 清空按钮

### 3. 修改 ScriptStore（如需要）
在 `frontend/src/stores/script.ts` 中：

```typescript
// 添加获取脚本列表的 action（如果还没有）
async function fetchScripts() {
    // ... 现有逻辑
}

// 解析脚本ID列表为逗号分隔字符串
function parseScriptIDs(ids: number[]): string {
    return ids.join(',')
}
```

### 4. 组件代码示例

```vue
<template>
  <div class="script-selector">
    <el-select
      v-model="selectedIds"
      multiple
      filterable
      placeholder="选择脚本"
      @change="handleChange"
    >
      <el-option
        v-for="script in scripts"
        :key="script.id"
        :label="script.name"
        :value="script.id"
      >
        <span>{{ script.name }}</span>
        <span class="script-path">{{ script.target_path }}</span>
      </el-option>
    </el-select>
    <div class="selected-info">
      已选择 {{ selectedIds.length }} 个脚本
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useScriptStore } from '@/stores/script'

const props = defineProps<{
  modelValue: number[]
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: number[]): void
}>()

const scriptStore = useScriptStore()
const selectedIds = ref<number[]>([])

// 初始化
if (props.modelValue.length > 0) {
  selectedIds.value = [...props.modelValue]
}

function handleChange(value: number[]) {
  emit('update:modelValue', value)
}
</script>
```

## 输入

- 现有 Vue 组件结构
- script store (frontend/src/stores/script.ts)
- Element Plus UI 库

## 输出

- 新组件：frontend/src/components/ScriptSelector.vue
- 可选：修改 script store

## 依赖

- 前端开发环境正常运行
- Element Plus 已安装

## 验收标准

- [ ] ScriptSelector 组件支持多选
- [ ] 正确显示所有可用脚本
- [ ] 支持 filterable（可搜索）
- [ ] 正确 emit update:modelValue 事件
- [ ] 初始值正确显示已选项
- [ ] UI 友好，显示已选数量
