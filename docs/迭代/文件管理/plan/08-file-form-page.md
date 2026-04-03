# 文件管理 - 上传配置表单页面

## 任务描述

创建文件上传配置的新增和编辑表单页面。

## 详细说明

### 1. 创建页面文件

创建 `frontend/src/views/FileForm.vue`：

### 2. 页面结构

#### 新增模式 (`/file-uploads/new`)

包含以下表单字段：

##### 本地文件选择
- 使用 `<input type="file" multiple>` 触发系统文件选择器
- 文件格式限制：.tar, .sh, .zip, .conf, .txt, .json, .yml, .xml, .log, .sql
- 数量限制：最多 20 个文件
- 文件大小限制：单文件不超过 500MB
- 选择后显示已选文件列表（可移除）

##### 目标服务器选择
- 多选下拉框
- 数据源：`serverApi.list()` 获取服务器列表
- 选中展示：IP:PORT 格式
- 选择上限：最多 10 台

##### 远程路径配置
- 文本输入框
- 格式校验：必须以 `/` 开头
- 路径为空或格式错误时禁用保存按钮

##### 配置名称
- 文本输入框
- 用于标识此上传配置

##### 保存行为
- 点击保存后，使用 FormData 上传文件到后端
- 后端将文件保存到 `./data/files/` 目录
- 创建数据库记录
- 上传状态默认为"待执行"
- 成功后关闭弹窗/跳回列表

#### 编辑模式 (`/file-uploads/:id/edit`)

- 加载已有配置信息
- 可修改：目标服务器、远程路径
- 不可修改：本地文件列表
- 保存后更新数据库

### 3. 主要功能实现

```typescript
const form = ref({
    name: '',
    files: [] as File[],
    serverIds: [] as number[],
    remotePath: ''
});

// 文件选择处理
function handleFileSelect(event: Event) {
    const input = event.target as HTMLInputElement;
    const files = Array.from(input.files || []);
    
    // 格式校验
    const allowedExts = ['.tar', '.sh', '.zip', '.conf', '.txt', '.json', '.yml', '.xml', '.log', '.sql'];
    // 大小校验（500MB）
    
    form.value.files = files;
}

// 移除文件
function removeFile(index: number) {
    form.value.files.splice(index, 1);
}

// 远程路径校验
const isPathValid = computed(() => /^\//.test(form.value.remotePath));

// 提交表单
async function handleSubmit() {
    const formData = new FormData();
    formData.append('name', form.value.name);
    formData.append('remotePath', form.value.remotePath);
    formData.append('serverIds', JSON.stringify(form.value.serverIds));
    
    form.value.files.forEach(file => {
        formData.append('files', file);
    });
    
    await store.createUpload(formData);
    router.push('/file-uploads');
}
```

### 4. 路由配置

在 `frontend/src/router/index.ts` 中添加：

```typescript
{
    path: '/file-uploads',
    name: 'file-uploads',
    component: () => import('@/views/FileList.vue'),
    meta: { requiresAuth: true }
},
{
    path: '/file-uploads/new',
    name: 'file-upload-new',
    component: () => import('@/views/FileForm.vue'),
    meta: { requiresAuth: true }
},
{
    path: '/file-uploads/:id/edit',
    name: 'file-upload-edit',
    component: () => import('@/views/FileForm.vue'),
    meta: { requiresAuth: true }
}
```

### 5. NavBar 更新

在 `frontend/src/components/NavBar.vue` 中添加入口：

```typescript
// 如果用户有脚本管理权限，显示文件管理入口
{ path: '/file-uploads', name: '文件管理', icon: '📁' }
```

## 输入

- 需求文档 `06-新增文件管理.md`
- 现有 `frontend/src/views/ScriptForm.vue` 模式
- `07-file-list-page.md`

## 输出

- `frontend/src/views/FileForm.vue` - 上传配置表单
- `frontend/src/router/index.ts` - 添加路由
- `frontend/src/components/NavBar.vue` - 添加导航入口

## 依赖

- 07-file-list-page.md

## 验收标准

- [ ] 文件选择器正常
- [ ] 格式和大小校验正常
- [ ] 远程路径校验正常
- [ ] 服务器多选正常
- [ ] FormData 上传正常
- [ ] 编辑模式正常
- [ ] 路由正确注册
- [ ] NavBar 入口添加
