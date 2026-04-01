# 添加模式切换组件和文件上传UI

## 任务描述

在脚本表单页面（ScriptForm.vue）中添加「手动输入」和「文件上传」两种模式的切换按钮组，以及文件上传区域的基础UI结构。

## 详细说明

### 1. 添加模式切换按钮组

在"脚本内容"标签下方添加切换按钮组件：

```vue
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
```

### 2. 添加文件上传区域UI

当 `inputMode === 'upload'` 时显示上传区域：

```vue
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
```

### 3. 添加样式

在 `<style scoped>` 中添加：
- `.mode-toggle` - 按钮组容器
- `.mode-btn` - 切换按钮样式（选中态用紫色渐变背景）
- `.upload-area` - 上传区域容器
- `.upload-placeholder` - 上传提示区域
- `.drag-over` - 拖拽悬停状态

## 输入

- 现有 `frontend/src/views/ScriptForm.vue` 源文件
- 需求文档 `03-脚本文件上传.md`

## 输出

- 修改 `frontend/src/views/ScriptForm.vue`，添加模式切换按钮和上传区域UI骨架

## 依赖

- 无

## 验收标准

- [ ] 页面显示"手动输入"和"文件上传"两个切换按钮
- [ ] 默认选中"手动输入"模式
- [ ] 点击"文件上传"时显示上传区域
- [ ] 上传区域包含拖拽提示和点击选择入口
- [ ] 显示支持的文件格式提示
