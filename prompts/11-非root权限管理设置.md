# 非 Root 权限管理设置

## 需求概述

基于当前系统的 `ROLE_ADMIN` / `ROLE_USER` 角色体系，新增细粒度权限控制：**仅 `ROLE_ADMIN` 用户可以访问"脚本管理"和"文件管理"页面，普通 `ROLE_USER` 用户不可见且无法访问这两个模块**。

## 详细需求

### 1. 权限矩阵

| 功能模块     | ROLE_ADMIN | ROLE_USER |
|-----------|------------|-----------|
| 服务器管理   | ✅ 可见/可操作 | ✅ 可见/可操作 |
| 任务管理     | ✅ 可见/可操作 | ✅ 可见/可操作 |
| 个人设置     | ✅ 可见/可操作 | ✅ 可见/可操作 |
| 用户管理     | ✅ 可见/可操作 | ❌ 不可见     |
| **脚本管理** | ✅ 可见/可操作 | ❌ 不可见     |
| **文件管理** | ✅ 可见/可操作 | ❌ 不可见     |

### 2. 前端权限控制

#### 2.1 导航栏菜单过滤

修改 `frontend/src/components/NavBar.vue`，在脚本和文件的 `router-link` 上添加 `v-if` 条件：

```vue
<!-- NavBar.vue -->
<router-link to="/files" class="nav-link" v-if="authStore.isAdmin">文件</router-link>
<router-link to="/scripts" class="nav-link" v-if="authStore.isAdmin">脚本</router-link>
```

#### 2.2 路由守卫增强

修改 `frontend/src/router/auth-guard.ts`，对脚本和文件相关路由添加 `requiresAdmin` 检查：

```typescript
// router/auth-guard.ts
// 脚本路由
{
    path: '/scripts',
    name: 'scripts',
    component: () => import('@/views/ScriptList.vue'),
    meta: { requiresAuth: true, requiresAdmin: true }  // 新增 requiresAdmin: true
},
{
    path: '/scripts/new',
    name: 'script-new',
    component: () => import('@/views/ScriptForm.vue'),
    meta: { requiresAuth: true, requiresAdmin: true }  // 新增
},
{
    path: '/scripts/:id/edit',
    name: 'script-edit',
    component: () => import('@/views/ScriptForm.vue'),
    meta: { requiresAuth: true, requiresAdmin: true }  // 新增
},

// 文件路由
{
    path: '/files',
    name: 'files',
    component: () => import('@/views/FileList.vue'),
    meta: { requiresAuth: true, requiresAdmin: true }  // 新增
},
{
    path: '/files/new',
    name: 'file-new',
    component: () => import('@/views/FileForm.vue'),
    meta: { requiresAuth: true, requiresAdmin: true }  // 新增
},
{
    path: '/files/:id/edit',
    name: 'file-edit',
    component: () => import('@/views/FileForm.vue'),
    meta: { requiresAuth: true, requiresAdmin: true }  // 新增
},
```

#### 2.3 服务器详情页隐藏"执行脚本"按钮

修改服务器详情页（如 `ServerList.vue` 中的跳转逻辑），确保非 admin 用户无法从服务器列表进入脚本执行流程。

如果服务器详情页包含"执行脚本"或"上传文件"等入口，需要对这些入口添加 `v-if="authStore.isAdmin"` 条件。

### 3. 后端权限控制（可选加固）

当前后端 Handler 已通过中间件进行了 `requiresAdmin` 检查（见 `backend/internal/middleware/auth.go`），但需确认脚本和文件相关的 Handler 方法已正确应用管理员权限中间件：

```go
// backend/internal/handler/handler.go
// 确认脚本相关路由已应用 AdminRequired 中间件
scripts := protected.Group("")
scripts.Use(middleware.AdminRequired())
{
    scripts.GET("", h.ListScripts)
    scripts.POST("", h.CreateScript)
    // ...
}

// 确认文件相关路由已应用 AdminRequired 中间件
files := protected.Group("")
files.Use(middleware.AdminRequired())
{
    files.GET("", h.ListFiles)
    files.POST("", h.UploadFile)
    // ...
}
```

> **注意**：如果当前中间件尚未对脚本/文件路由单独分组，而是统一使用 `protected.Use(middleware.AuthRequired())`，则需要新增 `AdminRequired` 中间件并对脚本/文件路由进行分组保护。

### 4. Auth Store 增强（可选）

当前 `frontend/src/stores/auth.ts` 中已有 `isAdmin` 计算属性，可直接使用，无需修改。

如需新增 `canManageScripts` / `canManageFiles` 等语义化权限属性，可扩展：

```typescript
// stores/auth.ts
const canManageScripts = computed(() => user.value?.role === 'ROLE_ADMIN');
const canManageFiles = computed(() => user.value?.role === 'ROLE_ADMIN');
```

## 文件变更清单

| 文件 | 变更类型 | 说明 |
|------|---------|------|
| `frontend/src/components/NavBar.vue` | 修改 | 脚本和文件菜单添加 `v-if="authStore.isAdmin"` |
| `frontend/src/router/index.ts` | 修改 | 脚本/文件相关路由 `meta` 添加 `requiresAdmin: true` |
| `frontend/src/router/auth-guard.ts` | 修改 | 增强守卫逻辑，admin 验证已支持（非必要可跳过） |
| `frontend/src/views/ServerList.vue` | 修改 | 确认无"执行脚本"等危险入口暴露给普通用户 |
| `backend/internal/handler/handler.go` | 确认 | 脚本/文件 Handler 已应用管理员中间件 |
| `backend/internal/middleware/auth.go` | 确认 | 确认 `AdminRequired` 中间件存在且正确 |

## 验收标准

### 功能验收

- [ ] `ROLE_ADMIN` 用户登录后，导航栏显示：服务器、文件、脚本、任务、用户、个人
- [ ] `ROLE_USER` 用户登录后，导航栏显示：服务器、任务、个人（不显示文件、脚本、用户）
- [ ] `ROLE_USER` 用户直接访问 `/scripts`、`/scripts/new`、`/files`、`/files/new` 等路径时，被重定向到首页（`/`）
- [ ] `ROLE_ADMIN` 用户对所有页面有完全访问权限
- [ ] 服务器列表页无"执行脚本"等敏感入口暴露给普通用户

### 边界情况

- [ ] 直接 URL 访问脚本/文件页面（非导航栏入口）被正确拦截
- [ ] 切换用户角色后刷新页面，菜单状态正确更新
- [ ] 脚本/文件 API 接口（后端）在非 admin 用户调用时返回 403

### 兼容性

- [ ] 不影响现有 `ROLE_ADMIN` 用户的所有功能
- [ ] 不影响未登录用户的跳转逻辑

## 相关代码位置

- 前端导航栏：`frontend/src/components/NavBar.vue`
- 前端路由：`frontend/src/router/index.ts`
- 前端路由守卫：`frontend/src/router/auth-guard.ts`
- 前端权限状态：`frontend/src/stores/auth.ts`
- 后端认证中间件：`backend/internal/middleware/auth.go`
- 后端 Handler 路由注册：`backend/internal/handler/handler.go`
- 服务器列表页：`frontend/src/views/ServerList.vue`

## 注意事项

1. **双重防护**：前端路由守卫 + 后端中间件双重校验，前端防误触，后端防绕过
2. **用户体验**：非 admin 用户看不到菜单，但若通过 URL 直接访问，应被重定向而非显示 404
3. **服务器列表页**：如果服务器列表页或详情页有"执行脚本"按钮，需同步添加权限判断
4. **任务关联**：任务管理页面本身对所有用户可见，但如果任务关联了脚本，需确认普通用户无法创建/编辑涉及脚本的任务（建议任务创建时禁止选择脚本，仅 admin 可选）
