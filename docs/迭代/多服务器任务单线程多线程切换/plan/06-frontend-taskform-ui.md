# 前端任务表单 UI 修改

## 任务描述

在任务创建/编辑表单中添加执行模式选择器，让用户可以选择并发执行或单线程执行。

## 详细说明

### 1. 修改 TaskForm.vue 组件 (`frontend/src/views/TaskForm.vue`)

在表单中添加执行模式单选框组：

```vue
<template>
  <!-- 现有表单项... -->
  
  <div class="form-item">
    <label>执行模式</label>
    <div class="radio-group">
      <label class="radio-label">
        <input 
          type="radio" 
          v-model="form.execution_mode" 
          value="concurrent"
        />
        <span>并发执行</span>
        <span class="hint">多个服务器同时执行脚本</span>
      </label>
      <label class="radio-label">
        <input 
          type="radio" 
          v-model="form.execution_mode" 
          value="sequential"
        />
        <span>单线程执行</span>
        <span class="hint">按顺序执行，失败时终止后续服务器</span>
      </label>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue';

const form = reactive({
  script_ids: '',
  name: '',
  server_ids: [],
  execution_mode: 'concurrent'  // 默认并发
});
</script>
```

### 2. 表单数据提交

确保 `execution_mode` 随其他表单数据一起提交到后端 API。

### 3. 样式要求

- 单选框组布局清晰
- 提供执行模式说明文字
- 默认选中并发执行

### 4. 编辑模式支持

如果是在编辑模式下加载表单，需要从任务数据中读取 `execution_mode` 并设置到表单中。

## 输入

- 需求文档：`10-多服务器任务支持单线程和多线程切换.md`
- 现有表单组件：`frontend/src/views/TaskForm.vue`

## 输出

- 修改后的 `frontend/src/views/TaskForm.vue`

## 依赖

- 04, 05

## 验收标准

- [ ] 表单包含执行模式单选框组
- [ ] 并发执行和单线程执行两个选项可用
- [ ] 默认选中并发执行
- [ ] 每个选项有对应的说明文字
- [ ] 表单提交时包含 execution_mode 字段
- [ ] 编辑模式下正确加载任务的 execution_mode
