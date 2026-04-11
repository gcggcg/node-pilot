# 任务一：后端脚本管理 Admin 权限控制

## 任务描述

在后端对脚本管理相关的所有 API 路由添加 `RequireRole("ROLE_ADMIN")` 中间件，确保非 admin 用户通过 API 直接调用脚本相关接口时返回 403 Forbidden。

## 详细说明

### 1. 修改路由注册文件

文件：`backend/cmd/server/main.go`

找到第 127-136 行的 scripts 路由组，将：
```go
scripts := api.Group("/scripts")
scripts.Use(middleware.JWTAuth(jwtSecret))
{
    scripts.GET("", h.ListScripts)
    scripts.GET("/:id", h.GetScript)
    scripts.POST("", h.CreateScript)
    scripts.PUT("/:id", h.UpdateScript)
    scripts.DELETE("/:id", h.DeleteScript)
    scripts.POST("/batch-delete", h.DeleteScripts)
}
```

改为：
```go
scripts := api.Group("/scripts")
scripts.Use(middleware.JWTAuth(jwtSecret))
scripts.Use(middleware.RequireRole("ROLE_ADMIN"))
{
    scripts.GET("", h.ListScripts)
    scripts.GET("/:id", h.GetScript)
    scripts.POST("", h.CreateScript)
    scripts.PUT("/:id", h.UpdateScript)
    scripts.DELETE("/:id", h.DeleteScript)
    scripts.POST("/batch-delete", h.DeleteScripts)
}
```

### 2. 验证中间件存在

确认 `backend/internal/middleware/auth.go` 中已有 `RequireRole` 函数（已存在，无需修改）。

## 输入

- `backend/cmd/server/main.go` 第 127-136 行
- `backend/internal/middleware/auth.go`

## 输出

- 修改后的 `backend/cmd/server/main.go`，scripts 路由组增加了 `.Use(middleware.RequireRole("ROLE_ADMIN"))`

## 依赖

- 无（第一个任务，无依赖）

## 验收标准

- [ ] scripts 路由组中增加了 `.Use(middleware.RequireRole("ROLE_ADMIN"))`
- [ ] `middleware/auth.go` 中的 `RequireRole` 函数未被修改（确认存在）
- [ ] Go 编译通过（`cd backend && go build ./...`）
