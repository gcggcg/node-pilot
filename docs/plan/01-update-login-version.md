# 前端登录页版本号升级至 v1.1.0

## 任务描述

将前端登录页面 `Login.vue` 中的版本号从 `v1.0.0` 修改为 `v1.1.0`，并重新构建前端资源以更新嵌入式产物。

## 详细说明

1. **修改 Login.vue**：
   - 文件路径：`frontend/src/views/Login.vue`
   - 定位第 102 行（或搜索 `<span class="version">v1.0.0</span>`）
   - 将 `v1.0.0` 改为 `v1.1.0`

2. **重新构建前端资源**（可选但推荐）：
   - 进入 `frontend/` 目录执行 `npm run build`
   - 构建产物会更新至 `backend/web/assets/` 目录
   - 注意：`backend/web/assets/` 下的 `Login-*.js` 是构建产物，旧版本文件属正常现象，构建后会被新版本替换

3. **验证**：
   - 检查 `Login.vue` 第 102 行已改为 `v1.1.0`
   - 如已构建，确认新构建产物存在于 `backend/web/assets/`

## 输入

- `frontend/src/views/Login.vue`（当前版本号为 v1.0.0）

## 输出

- 修改后的 `frontend/src/views/Login.vue`（版本号为 v1.1.0）
- （可选）重新构建的 `backend/web/assets/` 下的前端产物

## 依赖

- 无前置任务依赖

## 验收标准

- [ ] `frontend/src/views/Login.vue` 第 102 行已更新为 `<span class="version">v1.1.0</span>`
- [ ] 文件中不再存在 `v1.0.0` 的版本号字样（版本号区域）
- [ ] （可选）前端构建成功，`backend/web/assets/` 中存在最新构建产物
