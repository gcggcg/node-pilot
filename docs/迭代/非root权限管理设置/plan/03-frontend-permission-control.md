# 任务三：前端权限控制与交互限制

## 任务描述

完成前端权限控制：对 NavBar 菜单、路由守卫、ScriptSelector 组件和 TaskForm 进行 admin 专属限制，确保普通用户无法看到、访问或使用脚本相关功能。

## 详细说明

### 1. NavBar 菜单权限过滤（已部分完成，需确认）

文件：`frontend/src/components/NavBar.vue`

确认第 6-9 行已包含正确的 `v-if` 条件：
```vue
<router-link to="/files" class="nav-link" v-if="authStore.isAdmin">文件</router-link>
<router-link to="/scripts" class="nav-link" v-if="authStore.isAdmin">脚本</router-link>
```

如果未添加，需添加。

### 2. 路由守卫确认（已存在）

文件：`frontend/src/router/auth-guard.ts`

确认守卫逻辑中已有对 `requiresAdmin` 的检查（当前代码第 25-28 行已有）。无需修改。

### 3. 路由 meta 添加 requiresAdmin（已有，需确认）

文件：`frontend/src/router/index.ts`

确认以下路由的 `meta` 中已包含 `requiresAdmin: true`：
- `/scripts` (line ~36-40)
- `/scripts/new` (line ~42-46)
- `/scripts/:id/edit` (line ~48-52)
- `/files` (line ~96-100)
- `/files/new` (line ~102-106)
- `/files/:id/edit` (line ~108-112)

如果未添加，需在每个路由的 `meta` 中添加：
```typescript
meta: { requiresAuth: true, requiresAdmin: true }
```

### 4. ScriptSelector 组件权限限制

文件：`frontend/src/components/ScriptSelector.vue`

在组件根元素上添加 `v-if="authStore.isAdmin"`：
```vue
<div class="script-selector" v-if="authStore.isAdmin">
```

并在组件底部添加一个提示给非 admin 用户：
```vue
<div v-else class="no-permission">
    <p>您没有权限管理脚本</p>
</div>
```

添加样式：
```css
.no-permission {
    padding: 12px;
    background: #f8f9fa;
    border-radius: 6px;
    color: #999;
    font-size: 14px;
    text-align: center;
    border: 1px dashed #ddd;
}
```

同时在 `<script setup>` 中导入 authStore：
```typescript
import { useAuthStore } from '@/stores/auth';
const authStore = useAuthStore();
```

### 5. TaskForm 脚本选择区域权限限制

文件：`frontend/src/views/TaskForm.vue`

在"脚本"选择区域的容器上添加 `v-if="authStore.isAdmin"`：
```vue
<div class="form-group" v-if="authStore.isAdmin">
    <label>脚本 (支持批量)</label>
    <ScriptSelector 
        v-model="selectedScriptIds"
        multiple
        placeholder="选择要执行的脚本（可多选）"
    />
</div>
```

在"脚本"区域下方添加提示（v-else）：
```vue
<div class="form-group" v-else>
    <label>脚本</label>
    <div class="readonly-hint">普通用户无法选择脚本，请联系管理员</div>
</div>
```

导入 authStore（如果没有）：
```typescript
import { useAuthStore } from '@/stores/auth';
const authStore = useAuthStore();
```

添加提示样式：
```css
.readonly-hint {
    padding: 10px 12px;
    background: #fff3cd;
    border: 1px solid #ffc107;
    border-radius: 6px;
    color: #856404;
    font-size: 14px;
}
```

### 6. TaskForm 表单验证调整

因为非 admin 用户无法选择脚本，需要调整 `isFormValid` 逻辑：

修改 `isFormValid`：
```typescript
const isFormValid = computed(() => {
    if (authStore.isAdmin) {
        // Admin 用户必须选择脚本
        const hasScripts = selectedScriptIds.value.length > 0 || form.value.script_id;
        return form.value.name.trim() && hasScripts && form.value.server_ids.length > 0;
    } else {
        // 普通用户不选脚本，只需名称和服务器
        return form.value.name.trim() && form.value.server_ids.length > 0;
    }
});
```

### 7. 任务创建逻辑调整（非 admin 用户）

由于普通用户无法选择脚本，`handleSubmit` 中的 `payload` 构造需要处理：

```typescript
// 如果是 admin 且选择了脚本，才传 script_ids
if (authStore.isAdmin && selectedScriptIds.value.length > 0) {
    payload.script_ids = scriptIdsString.value;
} else if (authStore.isAdmin && form.value.script_id) {
    payload.script_id = Number(form.value.script_id);
}
// 非 admin 不传 script_ids（后端会创建不含脚本的任务）
```

> 注：后端 CreateTask 接口本身不强制要求有脚本 ID，所以普通用户可以创建不含脚本的任务。如果后端需要脚本才能执行，则普通用户创建的任务会无法执行（这是预期行为）。

## 输入

- `frontend/src/components/NavBar.vue`
- `frontend/src/router/index.ts`
- `frontend/src/router/auth-guard.ts`
- `frontend/src/components/ScriptSelector.vue`
- `frontend/src/views/TaskForm.vue`
- `frontend/src/stores/auth.ts`（确认 isAdmin 已存在）

## 输出

- 修改后的 NavBar.vue
- 修改后的 ScriptSelector.vue
- 修改后的 TaskForm.vue
- 确认/修改后的 router/index.ts

## 依赖

- 无（可与任务一、任务二并行开发，但最终一起验收）

## 验收标准

- [ ] NavBar 中脚本和文件菜单有 `v-if="authStore.isAdmin"`
- [ ] 脚本/文件路由的 `meta` 中有 `requiresAdmin: true`
- [ ] ScriptSelector 组件对非 admin 用户显示"没有权限"提示
- [ ] TaskForm 中非 admin 用户看到"无法选择脚本"提示
- [ ] 非 admin 用户的 `isFormValid` 不强制要求选择脚本
- [ ] 前端 TypeScript 编译通过（`cd frontend && npm run build` 或类型检查）
