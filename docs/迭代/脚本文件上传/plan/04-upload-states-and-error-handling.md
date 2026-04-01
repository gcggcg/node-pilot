# 添加上传状态和错误处理

## 任务描述

完善上传功能的状态反馈，包括上传进度、成功/失败状态显示，以及错误重试机制。

## 详细说明

### 1. 添加状态相关的响应式变量

```typescript
const uploadStatus = ref<'idle' | 'uploading' | 'success' | 'error'>('idle')
const selectedFileName = ref<string>('')

// 重置上传状态
function resetUpload() {
    uploadStatus.value = 'idle'
    uploadError.value = null
    selectedFileName.value = ''
    uploadedContent.value = ''
    form.value.content = manualContent.value || '#!/bin/bash\n'
}
```

### 2. 更新 handleFile 函数的状态管理

```typescript
function handleFile(file: File) {
    uploadError.value = null
    selectedFileName.value = file.name
    
    // 验证文件
    const error = validateFile(file)
    if (error) {
        uploadError.value = error
        uploadStatus.value = 'error'
        return
    }
    
    // 读取文件内容
    uploadStatus.value = 'uploading'
    const reader = new FileReader()
    
    reader.onload = (e) => {
        const content = e.target?.result as string
        uploadedContent.value = content
        form.value.content = content
        uploadStatus.value = 'success'
        isUploading.value = false
    }
    
    reader.onerror = () => {
        uploadError.value = '文件读取失败，请重试'
        uploadStatus.value = 'error'
        isUploading.value = false
    }
    
    reader.readAsText(file)
}
```

### 3. 更新模板添加状态显示

```vue
<!-- 上传区域 -->
<div v-if="inputMode === 'upload'" class="upload-area">
    <input 
        type="file" 
        ref="fileInput" 
        accept=".txt,.sh,.py,.js,.sql"
        @change="onFileSelect"
        style="display: none"
    />
    
    <!-- 空闲/错误状态：显示上传入口 -->
    <div v-if="uploadStatus === 'idle' || uploadStatus === 'error'" 
         class="upload-placeholder" 
         @click="triggerFileSelect">
        <span class="upload-icon">📁</span>
        <p>拖拽文件到此处，或 <span class="link">点击选择</span></p>
        <p class="hint">支持 .txt, .sh, .py, .js, .sql 文件，不超过 5MB</p>
    </div>
    
    <!-- 上传中状态 -->
    <div v-if="uploadStatus === 'uploading'" class="upload-progress">
        <div class="spinner"></div>
        <p>正在读取文件...</p>
    </div>
    
    <!-- 成功状态 -->
    <div v-if="uploadStatus === 'success'" class="upload-success">
        <div class="success-header">
            <span class="success-icon">✓</span>
            <span>文件读取成功</span>
        </div>
        <p class="file-name">{{ selectedFileName }}</p>
        <button type="button" class="btn-reupload" @click="triggerFileSelect">
            重新上传
        </button>
    </div>
    
    <!-- 错误提示 -->
    <div v-if="uploadError" class="upload-error">
        <span class="error-icon">⚠</span>
        <span>{{ uploadError }}</span>
        <button type="button" class="btn-retry" @click="resetUpload">
            重新上传
        </button>
    </div>
</div>
```

### 4. 添加样式

```css
.upload-progress {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 40px;
}

.spinner {
    width: 32px;
    height: 32px;
    border: 3px solid #f3f3f3;
    border-top: 3px solid #667eea;
    border-radius: 50%;
    animation: spin 1s linear infinite;
}

@keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
}

.upload-success,
.upload-error {
    padding: 16px;
    border-radius: 6px;
    margin-top: 12px;
}

.upload-success {
    background: #d4edda;
    border: 1px solid #c3e6cb;
}

.success-header {
    display: flex;
    align-items: center;
    gap: 8px;
    color: #155724;
    font-weight: 500;
}

.success-icon {
    color: #28a745;
    font-size: 18px;
}

.file-name {
    margin: 8px 0;
    color: #666;
    font-size: 13px;
}

.upload-error {
    background: #f8d7da;
    border: 1px solid #f5c6cb;
    display: flex;
    align-items: center;
    gap: 8px;
    color: #721c24;
}

.error-icon {
    font-size: 18px;
}

.btn-reupload,
.btn-retry {
    margin-left: auto;
    padding: 4px 12px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 12px;
}

.btn-reupload {
    background: #667eea;
    color: white;
}

.btn-retry {
    background: #dc3545;
    color: white;
}

.drag-over {
    border-color: #667eea !important;
    background: #f8f7ff !important;
}
```

## 输入

- 更新后的 `frontend/src/views/ScriptForm.vue`（包含模式切换和数据保留）

## 输出

- 更新 `frontend/src/views/ScriptForm.vue`，完善所有状态反馈和错误处理

## 依赖

- 03

## 验收标准

- [ ] 上传过程中显示加载动画
- [ ] 上传成功后显示绿色成功提示和文件名
- [ ] 上传失败后显示红色错误提示
- [ ] 成功和错误状态都提供重新上传按钮
- [ ] 拖拽文件悬停时显示高亮效果
- [ ] 网络异常等情况下能显示错误并重试
