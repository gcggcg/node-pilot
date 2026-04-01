# 实现文件读取和内容填充逻辑

## 任务描述

实现客户端文件读取功能，在文件上传成功后读取文件内容并填充到脚本内容文本框。

## 详细说明

### 1. 添加响应式状态

在 `<script setup>` 中添加：

```typescript
const fileInput = ref<HTMLInputElement | null>(null)
const isUploading = ref(false)
const uploadError = ref<string | null>(null)

// 保留两个模式的内容
const manualContent = ref(form.value.content)
const uploadedContent = ref('')

// 当前输入模式
const inputMode = ref<'manual' | 'upload'>('manual')
```

### 2. 文件选择处理函数

```typescript
function triggerFileSelect() {
    fileInput.value?.click()
}

function onFileSelect(event: Event) {
    const target = event.target as HTMLInputElement
    const file = target.files?.[0]
    if (file) {
        handleFile(file)
    }
}

function onDragOver(event: DragEvent) {
    (event.currentTarget as HTMLElement).classList.add('drag-over')
}

function onDragLeave(event: DragEvent) {
    (event.currentTarget as HTMLElement).classList.remove('drag-over')
}

function onDrop(event: DragEvent) {
    (event.currentTarget as HTMLElement).classList.remove('drag-over')
    const file = event.dataTransfer?.files?.[0]
    if (file) {
        handleFile(file)
    }
}
```

### 3. 文件验证和读取函数

```typescript
const ALLOWED_EXTENSIONS = ['.txt', '.sh', '.py', '.js', '.sql']
const MAX_FILE_SIZE = 5 * 1024 * 1024 // 5MB

function validateFile(file: File): string | null {
    const ext = '.' + file.name.split('.').pop()?.toLowerCase()
    if (!ALLOWED_EXTENSIONS.includes(ext)) {
        return `不支持的文件格式，请上传 ${ALLOWED_EXTENSIONS.join(', ')} 文件`
    }
    if (file.size > MAX_FILE_SIZE) {
        return '文件大小超过5MB限制'
    }
    return null
}

function handleFile(file: File) {
    uploadError.value = null
    
    // 验证文件
    const error = validateFile(file)
    if (error) {
        uploadError.value = error
        return
    }
    
    // 读取文件内容
    isUploading.value = true
    const reader = new FileReader()
    
    reader.onload = (e) => {
        const content = e.target?.result as string
        uploadedContent.value = content
        form.value.content = content // 直接填充到表单
        isUploading.value = false
    }
    
    reader.onerror = () => {
        uploadError.value = '文件读取失败，请重试'
        isUploading.value = false
    }
    
    reader.readAsText(file)
}
```

### 4. 更新提交逻辑

修改 `handleSubmit`，确保使用当前模式的内容：

```typescript
// 确保提交前同步内容
if (inputMode.value === 'manual') {
    form.value.content = manualContent.value
} else {
    form.value.content = uploadedContent.value
}
```

## 输入

- 修改后的 `frontend/src/views/ScriptForm.vue`（包含模式切换UI）

## 输出

- 更新 `frontend/src/views/ScriptForm.vue`，实现文件读取和内容填充

## 依赖

- 01

## 验收标准

- [ ] 选择文件后自动验证格式和大小
- [ ] 验证失败显示错误提示
- [ ] 验证通过后读取文件内容
- [ ] 文件内容成功填充到脚本内容文本框
- [ ] 可以手动编辑填充后的内容
