# 脚本和任务列表分页

## 任务描述

为 ScriptList.vue 和 TaskList.vue 添加分页功能。

## 详细说明

### 1. ScriptList.vue 分页改造

参考 ServerList.vue 的实现：
- 添加 Pagination 组件
- 使用 useRouter/useRoute 管理 URL 参数
- 调用 store.fetchScripts(page, pageSize)
- 处理分页变化更新 URL

### 2. TaskList.vue 分页改造

参考 ServerList.vue 的实现：
- 添加 Pagination 组件
- 使用 useRouter/useRoute 管理 URL 参数
- 调用 store.fetchTasks(page, pageSize)
- 处理分页变化更新 URL

## 输入

- `frontend/src/views/ScriptList.vue`
- `frontend/src/views/TaskList.vue`
- `frontend/src/components/Pagination.vue`

## 输出

- ScriptList.vue 和 TaskList.vue 支持分页

## 依赖

- 05

## 验收标准

- [ ] ScriptList.vue 分页功能正常
- [ ] TaskList.vue 分页功能正常
- [ ] 三个列表页分页交互体验一致
