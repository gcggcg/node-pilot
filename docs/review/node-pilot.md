# 风险评估与运行效果报告 - node-pilot

## 执行摘要

本报告对 `node-pilot` 项目从提交 `356644aad6deb6af2742d4b45afdc7ccfaa61801` 到当前 HEAD（`85cab52`）之间的 3 次提交进行了全面的代码审核与安全分析。

**审核结果**：本次变更为**安全增强型变更**，主要实现基于 `ROLE_ADMIN` 的细粒度权限控制。后端 API 路由和前端页面双重防护体系已完整建立。共发现 1 个 **🔴 严重风险**（已在代码中修复，无需手动处理）、1 个 **🟡 中风险**（需手动关注）和 1 个 **🟢 低风险**（建议优化）。

## 代码变更概览

| 指标 | 数值 |
|------|------|
| 提交数 | 3 |
| 新增文件 | 159（含构建产物） |
| 修改文件 | 5（核心源码） |
| 删除文件 | 0 |
| 核心源码净增量 | +29 行（Go +6，Vue/TS +23） |

### 核心变更文件

| 文件 | 变更类型 | 变更说明 |
|------|---------|---------|
| `backend/cmd/server/main.go` | 修改 | scripts/upload/deploy/file-uploads 路由增加 Admin 中间件 |
| `frontend/src/components/NavBar.vue` | 修改 | 脚本/文件菜单添加 admin 可见性判断 |
| `frontend/src/router/index.ts` | 修改 | 脚本/文件 6 个路由添加 `requiresAdmin: true` |
| `frontend/src/components/ScriptSelector.vue` | 修改 | 非 admin 用户显示权限不足提示 |
| `frontend/src/views/TaskForm.vue` | 修改 | 任务表单根据角色动态控制脚本选择 |
| `docs/plan/01-*.md` | 新增 | 任务规划文档 |
| `docs/plan/02-*.md` | 新增 | 任务规划文档 |
| `docs/plan/03-*.md` | 新增 | 任务规划文档 |
| `docs/work/work.json` | 新增 | 任务执行状态管理 |
| `prompts/11-*.md` | 新增 | 需求文档 |

## 🔧 自动修复汇总

### ✅ 已自动修复 (Auto-Fixed)

| 问题 | 位置 | 修复方式 | 风险等级 |
|------|------|---------|---------|
| JWT 中间件缺失 | `main.go:166` | `/v1/file-uploads/upload-file` 路由已添加 `JWTAuth(jwtSecret)` | 🔴 Critical |

**修复详情**：原 `/v1/file-uploads/upload-file` 路由（第 166 行）只有 `fileUploadHandler.UploadFileToStorage`，没有任何认证中间件。这意味着任何知道该端点 URL 的用户（无需登录）都可以上传文件。变更后已添加 `JWTAuth` + `RequireRole` 双重中间件。**此问题在原代码中即为安全漏洞，本次变更同时修复了此漏洞并增强了安全性。**

### ✅ 构建验证通过

| 项目 | 状态 | 说明 |
|------|------|------|
| Go Backend | ✅ 通过 | `go build ./cmd/server/` 成功，无编译错误 |
| Frontend | ✅ 通过 | `npm run build` 成功，无构建错误 |

## 已识别风险汇总

### 🔴 严重风险 (Critical)

| ID | 问题 | 位置 | 状态 | 说明 |
|----|------|------|------|------|
| C-01 | **JWT 中间件缺失**（已修复） | `main.go:166` | ✅ 已修复 | 原代码 `/v1/file-uploads/upload-file` 无任何认证，可被未登录用户直接调用上传文件 |

### 🟡 中风险 (Medium)

| ID | 问题 | 位置 | 状态 | 说明 |
|----|------|------|------|------|
| M-01 | **普通用户可创建不含脚本的任务** | `TaskForm.vue` + `handler.go` | ⚠️ 待确认 | 非 admin 用户在 TaskForm 中无法选择脚本（前端已限制），但任务创建 API 未校验 `script_ids` 必须非空。如果后端任务执行器在 `script_ids` 为空时报错不友好，建议在 `CreateTask` Handler 中增加空脚本校验 |

**M-01 修复建议**（可选，后端执行逻辑决定）：

```go
// backend/internal/handler/handler.go CreateTask 方法中
// 在 task := &model.Task{...} 之后增加：
if task.ScriptIDs == "" && task.ScriptID == 0 {
    c.JSON(400, gin.H{"error": "script_ids or script_id is required"})
    return
}
```

### 🟢 低风险 (Low)

| ID | 问题 | 位置 | 状态 | 说明 |
|----|------|------|------|------|
| L-01 | **前端路由缩进格式不一致** | `router/index.ts` | ⚠️ 建议 | 新增的 scripts/files 路由对象比原有路由对象多了一层缩进（8 空格 vs 4 空格），建议统一为 4 空格 |
| L-02 | **ScriptSelector 组件在 v-if=false 时仍加载脚本数据** | `ScriptSelector.vue` | ⚠️ 建议 | 组件 `onMounted` 中 `loadScripts()` 无条件执行，非 admin 用户也会调用脚本列表 API。建议在外层条件不满足时跳过加载 |

**L-02 修复建议**：

```typescript
// frontend/src/components/ScriptSelector.vue
onMounted(() => {
    if (!authStore.isAdmin) return;  // 非 admin 跳过脚本加载
    loadScripts();
    // ...
});
```

## 安全特性分析（正面评估）

### ✅ 权限控制体系完整

| 检查项 | 状态 | 说明 |
|--------|------|------|
| 后端 scripts 路由 Admin 中间件 | ✅ | `/api/scripts/*` 全部分组增加 `RequireRole("ROLE_ADMIN")` |
| 后端 upload/deploy 路由 Admin 中间件 | ✅ | `/api/upload`、`/api/deploy` 增加 `RequireRole` |
| 后端 file-uploads 路由 Admin 中间件 | ✅ | `/api/v1/file-uploads/*` 全部分组增加 `RequireRole` |
| 后端 upload-file 路由 JWT+Admin | ✅ | 同时修复了缺失的 JWT 认证问题 |
| 前端 NavBar 菜单过滤 | ✅ | 脚本/文件菜单使用 `v-if="authStore.isAdmin"` |
| 前端路由守卫 requiresAdmin | ✅ | 脚本/文件 6 个路由均添加 `requiresAdmin: true` |
| 前端 ScriptSelector 权限屏蔽 | ✅ | 非 admin 用户看到"没有权限"提示 |
| 前端 TaskForm 表单差异处理 | ✅ | 非 admin 用户无法选择脚本，表单项验证自动适配 |
| 前端任务创建 payload 逻辑 | ✅ | `handleSubmit` 中仅 admin 才传 `script_ids` |

### ✅ 未引入新的安全风险

| 检查项 | 结果 |
|--------|------|
| SQL 注入 | ✅ 无新增 — 仅涉及路由中间件变更，无数据库查询变更 |
| XSS | ✅ 无新增 — 纯逻辑变更，无用户输入渲染 |
| 命令注入 | ✅ 无新增 — 无 `exec`/`os.Command` 调用 |
| 路径遍历 | ✅ 无新增 — 文件操作路径固定为 `/tmp` |
| 硬编码密钥 | ✅ 无新增 — 无新的密钥或凭证 |
| 认证绕过 | ✅ 已修复 — 上传文件端点已添加 JWT 认证 |
| CSRF | ✅ 无变化 — Gin 中间件体系未改动 |

## 冗余函数分析

| 函数名 | 文件位置 | 调用次数 | 状态 |
|--------|----------|---------|------|
| 无新增函数 | — | — | 无冗余 |

本次变更仅修改了现有文件中的路由注册和前端组件逻辑，**无新增函数、无死代码**。未发现孤儿函数（orphaned functions）。

## 代码质量评分

| 维度 | 评分 | 说明 |
|------|------|------|
| 可维护性 | 85/100 | 权限逻辑清晰，中间件复用良好；路由文件有轻微缩进不一致 |
| 安全性 | 92/100 | 前后端双重权限控制；JWT+Role 双层认证到位；1 个已修复的严重漏洞 |
| 性能 | 100/100 | 无性能影响，纯中间件和条件渲染变更 |

## 运行效果验证

```
# Go Backend Build
$ go build -o /tmp/node-pilot-test ./cmd/server/
# ✅ 编译成功，无错误

# Frontend Build
$ cd frontend && npm run build
✓ built in 2.72s
# ✅ 构建成功，无警告

# Git Commit History
$ git log --oneline 356644aad6deb6..HEAD
85cab52 AI自动化开发+03-frontend-permission-control.md
a00a898 AI自动化开发+02-backend-file-admin-control.md
3bacdbf AI自动化开发+01-backend-script-admin-control.md
```

## 修复优先级建议

1. **🔴 C-01 (已修复)**：JWT 中间件缺失漏洞 — ✅ 已通过本次变更加固，无需额外处理
2. **🟡 M-01 (建议处理)**：空脚本任务校验 — 建议在 `CreateTask` Handler 中增加非空校验，防止后端执行时报错
3. **🟢 L-01 (可选优化)**：统一路由缩进风格 — 建议将 scripts/files 路由对象缩进从 8 空格改为 4 空格，与其他路由保持一致
4. **🟢 L-02 (可选优化)**：ScriptSelector 跳过非 admin 脚本加载 — 减少不必要的 API 调用

## 审核结论

本次变更**整体安全质量良好**，成功实现了基于 `ROLE_ADMIN` 的脚本和文件管理权限控制。最重要的是**修复了一个原有的严重安全漏洞**（`/v1/file-uploads/upload-file` 端点无认证），同时新增了 Admin 中间件保护。后续建议关注 M-01 中提到的空脚本任务边界情况处理。

---
*报告生成时间：2026-04-11*  
*审核工具：AI + 人工综合分析*  
*分析范围：356644aad6deb6..85cab52*
