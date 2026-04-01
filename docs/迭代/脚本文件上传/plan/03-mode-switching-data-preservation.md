# 实现模式切换与数据保留

## 任务描述

实现模式切换时保留各模式下的输入内容，确保用户切换模式后可以恢复之前的内容。

## 详细说明

### 1. 修改 watch 监听 inputMode

```typescript
import { watch } from 'vue'

// 切换模式时同步内容
watch(inputMode, (newMode, oldMode) => {
    if (newMode === 'manual' && uploadedContent.value) {
        // 切换到手动模式，保存上传内容用于恢复
        manualContent.value = form.value.content
    } else if (newMode === 'upload' && manualContent.value) {
        // 切换到上传模式，保存手动输入内容用于恢复
        uploadedContent.value = form.value.content
    }
})
```

### 2. 修改 v-model 绑定

将原来的 `v-model="form.content"` 改为根据模式绑定：

```vue
<!-- 手动输入模式 -->
<textarea 
    v-if="inputMode === 'manual'"
    v-model="form.content" 
    rows="12" 
    required 
    placeholder="#!/bin/bash&#10;echo 'Hello World'" 
    class="code"
></textarea>

<!-- 上传模式：显示只读内容和上传状态 -->
<div v-else class="upload-preview">
    <div class="preview-header">
        <span>文件内容预览</span>
        <span class="file-info">{{ uploadedContent.length }} 字符</span>
    </div>
    <textarea 
        v-model="form.content" 
        rows="12" 
        class="code"
        readonly
    ></textarea>
</div>
```

### 3. 添加样式

```css
.upload-preview {
    border: 1px solid #ddd;
    border-radius: 6px;
    overflow: hidden;
}

.preview-header {
    display: flex;
    justify-content: space-between;
    padding: 8px 12px;
    background: #f8f9fa;
    border-bottom: 1px solid #ddd;
    font-size: 12px;
    color: #666;
}

.file-info {
    color: #667eea;
}
```

### 4. 修改 handleSubmit 确保提交正确内容

```typescript
async function handleSubmit() {
    // 切换到手动模式以确保内容同步
    if (inputMode.value === 'manual') {
        manualContent.value = form.value.content
    }
    form.value.content = inputMode.value === 'manual' 
        ? manualContent.value 
        : uploadedContent.value
    
    loading.value = true
    try {
        if (isEdit.value) {
            await store.updateScript(Number(route.params.id), form.value)
        } else {
            await store.createScript(form.value)
        }
        router.push('/scripts')
    } catch (e: any) {
        alert(e.message || '保存失败')
    } finally {
        loading.value = false
    }
}
```

## 输入

- 更新后的 `frontend/src/views/ScriptForm.vue`（包含文件读取功能）

## 输出

- 更新 `frontend/src/views/ScriptForm.vue`，实现模式切换和数据保留

## 依赖

- 02

## 验收标准

- [ ] 切换模式时保留各模式下的内容
- [ ] 切换回原模式时恢复之前的内容
- [ ] 上传模式下文本框显示只读的内容预览
- [ ] 提交时正确提交当前模式的内容
