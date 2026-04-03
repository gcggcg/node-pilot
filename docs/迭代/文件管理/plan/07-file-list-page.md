# 文件管理 - 文件列表页面

## 任务描述

创建文件上传管理的列表页面。

## 详细说明

### 1. 创建页面文件

创建 `frontend/src/views/FileList.vue`：

### 2. 页面结构

参考 `ScriptList.vue` 和 `ServerList.vue` 的卡片布局或表格布局：

#### 顶部区域
- 标题：文件管理
- 右侧按钮：新增上传配置

#### 筛选区域
- 状态下拉框：全部/待执行/成功/失败
- 文件名搜索框
- 日期范围选择器（可选）
- 查询按钮、重置按钮

#### 列表区域
- 表格布局，显示字段：
  - 序号
  - 配置名称
  - 本地路径（./data/files/ 下相对路径）
  - 目标服务器列表
  - 远程路径
  - 上传状态（待执行/成功/失败）
  - 创建时间
  - 操作列

#### 操作列
- 编辑按钮
- 执行按钮
- 查看结果按钮
- 删除按钮

#### 底部区域
- 分页组件
- 批量删除按钮（选中时显示）

### 3. 查看结果弹窗

点击"查看结果"按钮时弹出模态框：
- 标题：配置名称 + 执行时间
- 列表显示：
  - 服务器信息（IP:PORT）
  - 文件名
  - 状态（成功图标/失败图标）
  - 失败原因（仅失败项显示）

### 4. 主要功能实现

```typescript
const filters = ref({
    status: '',
    fileName: '',
    startTime: '',
    endTime: ''
});

async function handleSearch() {
    pagination.value.page = 1;
    await store.fetchUploads(1, pagination.value.pageSize, filters.value);
}

async function handleExecute(id: number) {
    if (!confirm('确定要执行上传吗？')) return;
    await store.executeUpload(id);
    await store.fetchUploads(pagination.value.page, pagination.value.pageSize);
}

async function handleViewResults(id: number) {
    await store.fetchResults(id);
    showResultsModal.value = true;
}
```

### 5. 样式

复用 `ScriptList.vue` 和 `ServerList.vue` 的样式：
- 卡片式布局或表格布局
- 按钮样式复用
- 分页组件复用

## 输入

- 需求文档 `06-新增文件管理.md`
- 现有 `frontend/src/views/ScriptList.vue`
- 现有 `frontend/src/views/ServerList.vue`
- `06-frontend-store.md`

## 输出

- `frontend/src/views/FileList.vue` - 文件列表页面

## 依赖

- 06-frontend-store.md

## 验收标准

- [ ] 列表页面完整实现
- [ ] 筛选功能正常
- [ ] 分页功能正常
- [ ] 查看结果弹窗正常
- [ ] 批量删除功能正常
- [ ] 样式与项目一致
