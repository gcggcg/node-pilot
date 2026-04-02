# AGENTS.md - NodePilot Development Guide

> This file provides context for AI coding agents operating in the NodePilot repository.

## Project Overview

**NodePilot** is a batch server management platform (批量服务器管理平台) built with:
- **Backend**: Go + Gin + SQLite + WebSocket
- **Frontend**: Vue 3 + TypeScript + Vite + Pinia

## Directory Structure

```
node-pilot/
├── backend/
│   ├── cmd/server/main.go        # Entry point
│   ├── internal/
│   │   ├── config/              # Configuration
│   │   ├── handler/              # HTTP handlers (Gin)
│   │   ├── logger/              # Logging (DEBUG/INFO/WARN/ERROR)
│   │   ├── model/               # Data models
│   │   ├── repository/           # Database access (SQLite)
│   │   ├── service/              # Business logic (SSH, TaskExecutor)
│   │   └── websocket/            # WebSocket hub
│   └── web/                     # Embedded frontend assets
├── frontend/
│   ├── src/
│   │   ├── api/                 # Axios API wrappers
│   │   ├── components/          # Vue components (Pagination.vue)
│   │   ├── router/              # Vue Router (7 routes)
│   │   ├── stores/              # Pinia stores (server, script, task)
│   │   ├── types/               # TypeScript interfaces
│   │   ├── views/               # Page components
│   │   ├── App.vue
│   │   └── main.ts
│   ├── dist/                    # Built frontend
│   └── node_modules/
├── docs/
│   ├── plan/                   # Task plan documents
│   └── review/                  # Code review reports
└── data/                        # SQLite database
```

---

## Build, Lint, and Test Commands

### Backend (Go)

```bash
# Navigate to backend
cd backend

# Build binary
go build -o node-pilot ./cmd/server

# Run (from backend directory)
./node-pilot --db ../data/servers.db --listen :8080 --debug

# Run tests
go test ./...

# Run single test
go test -v -run TestFunctionName ./internal/repository/

# Run tests with coverage
go test -cover ./...

# Format code
go fmt ./...

# Vet code
go vet ./...
```

### Frontend (Node.js)

```bash
# Navigate to frontend
cd frontend

# Install dependencies
npm install

# Development server (proxies /api and /ws to localhost:8080)
npm run dev

# Type check
npm run build  # runs vue-tsc first

# Build for production
npm run build

# Preview production build
npm run preview

# Single test (if tests exist)
npm run test -- --run
```

### Full Stack Development

```bash
# Build frontend and copy to backend web directory
cd frontend && npm run build && cp -r dist/* ../backend/web/

# Run backend with embedded frontend
cd backend && ./node-pilot --db ../data/servers.db --listen :8080
```

---

## Code Style Guidelines

### Go (Backend)

**Formatting & Style**
- Use `gofmt` for formatting (standard Go style)
- 4-space indentation, tabs not allowed
- Maximum line length: 120 characters (soft limit)
- Blank lines between functions and type definitions

**Imports**
- Group imports: stdlib → third-party → internal
- Use `goimports` or `gofmt` to manage import order
- Alias packages when names conflict (e.g., `ws "node-pilot/internal/websocket"`)

```go
import (
    "encoding/json"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"

    "node-pilot/internal/model"
    "node-pilot/internal/repository"
)
```

**Naming Conventions**
- PascalCase for exported names (types, functions, constants)
- camelCase for unexported names (variables, local functions)
- Acronyms: `ID`, `URL`, `JSON`, `API` (not `Id`, `Url`)
- Package names: lowercase, no underscores
- Interface names: often `-er` suffix (e.g., `Handler`, `Repository`)

**Error Handling**
- Always handle errors explicitly; never ignore with `_`
- Return early on error in handlers
- Use `logger.Error()` for logged errors
- Return appropriate HTTP status codes: 400 (bad request), 404 (not found), 500 (server error)

```go
// GOOD
if err != nil {
    c.JSON(404, gin.H{"error": "server not found"})
    return
}

// BAD - ignoring error
id, _ := strconv.ParseInt(idStr, 10, 64)
```

**Struct Tags**
- JSON field tags: `json:"field_name"`
- Use `json:"-"` for sensitive fields that should not be serialized

```go
type Server struct {
    PasswordEncrypted string `json:"-"`  // Never expose to client
    ConnectionStatus  string `json:"connection_status"`
}
```

**Context Usage**
- Pass `*gin.Context` as first parameter in handlers
- Use `c.Request.Context()` for cancellation/timeout handling

### TypeScript & Vue (Frontend)

**Formatting**
- Use default TypeScript/Vue conventions
- 2-space indentation
- Single quotes for strings
- Trailing commas in multiline

**Imports**
- Order: Vue core → Vue ecosystem → third-party → local/relative
- Use absolute paths with `@/` alias (configured in vite.config.ts)

```typescript
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { useServerStore } from '@/stores/server'
```

**TypeScript**
- Enable strict mode (tsconfig.json: `"strict": true`)
- Define interfaces for all API responses and forms
- Use `ref<T>()` and `computed<T>()` with explicit types
- Avoid `any` type; use `unknown` when type is truly unknown

```typescript
// GOOD
interface Server {
    id: number
    name: string
    host: string
}

const servers = ref<Server[]>([])

// BAD
const servers = ref<any[]>([])
```

**Vue Component Structure**
```vue
<script setup lang="ts">
// Imports
// Props/Emits definitions
// Composables/stores
// Component logic
</script>

<template>
  <!-- HTML structure -->
</template>

<style scoped>
/* Component-scoped styles */
</style>
```

**Pinia Stores**
- Use composition API style (`defineStore` with `() => { ... }`)
- Return reactive state and methods
- **Important**: When using pagination, ensure `loadData()` uses `async/await` and sync local pagination state after mutations

```typescript
export const useServerStore = defineStore('server', () => {
    const servers = ref<Server[]>([])
    const loading = ref(false)
    const pagination = ref({ page: 1, pageSize: 10, total: 0 })

    async function fetchServers(page = 1, pageSize = 10) {
        loading.value = true
        try {
            const res = await serverApi.list({ page, pageSize })
            servers.value = res.data
            pagination.value = { page: res.page, pageSize: res.pageSize, total: res.total }
        } finally {
            loading.value = false
        }
    }

    return { servers, loading, pagination, fetchServers }
})
```

**Naming Conventions**
- Components: PascalCase (e.g., `ServerList.vue`, `OutputPanel.vue`, `Pagination.vue`)
- Composables/stores: camelCase with use prefix (e.g., `useServerStore`)
- Types/Interfaces: PascalCase (e.g., `Server`, `TaskForm`, `PaginatedResponse`)
- CSS classes: kebab-case with BEM-style for components

---

## TypeScript Interfaces

```typescript
interface Server {
    id: number
    name: string
    host: string
    port: number
    username: string
    password_encrypted?: string
    connection_status: 'online' | 'offline' | 'unknown'
    created_at: string
    updated_at: string
}

interface Script {
    id: number
    name: string
    description: string
    content: string
    target_path: string
    created_at: string
    updated_at: string
}

interface Task {
    id: number
    script_id: number
    name: string
    status: 'pending' | 'running' | 'completed' | 'cancelled' | 'failed'
    created_at: string
    started_at?: string
    finished_at?: string
}

interface WSMessage {
    type: 'output' | 'server_start' | 'server_done' | 'task_start' | 'task_done'
    task_id: number
    server_id?: number
    server_name?: string
    content?: string
    status?: string
    exit_code?: number
    timestamp: string
}

// Pagination interfaces
interface PaginatedResponse<T> {
    data: T[]
    total: number
    page: number
    pageSize: number
}

interface PaginationParams {
    page?: number
    pageSize?: number
}
```

---

## Common Patterns

### Backend Handler Pattern
```go
func (h *Handler) HandlerName(c *gin.Context) {
    // 1. Parse input
    // 2. Call repository/service
    // 3. Return response or error
    // 4. Always return on error (no fallthrough)
}
```

### Backend Pagination Handler Pattern
```go
func (h *Handler) ListWithPagination(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

    if page < 1 {
        page = 1
    }
    if pageSize < 1 || pageSize > 100 {
        pageSize = 10
    }

    items, total, err := h.repo.ListWithPagination(page, pageSize)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{
        "data":     items,
        "total":    total,
        "page":     page,
        "pageSize": pageSize,
    })
}
```

### Frontend API Call Pattern
```typescript
async function fetchData() {
    try {
        loading.value = true
        const response = await api.get('/endpoint')
        data.value = response.data
    } catch (error) {
        error.value = 'Failed to fetch data'
        console.error(error)
    } finally {
        loading.value = false
    }
}
```

### Frontend Pagination Pattern
```typescript
// In view component
const pagination = ref({
    page: Number(route.query.page) || 1,
    pageSize: Number(route.query.pageSize) || 10,
    total: 0
})

async function loadData() {
    await store.fetchTasks(pagination.value.page, pagination.value.pageSize)
    pagination.value.total = store.pagination.total  // Sync after async!
}

function handlePageChange(payload: { page: number; pageSize: number }) {
    pagination.value.page = payload.page
    pagination.value.pageSize = payload.pageSize
    router.replace({ query: { page: payload.page, pageSize: payload.pageSize } })
    store.fetchTasks(payload.page, payload.pageSize)
}

// After mutations (create/delete), always sync:
pagination.value.total = store.pagination.total
```

### WebSocket Usage
```typescript
function connectWebSocket(taskId: number) {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = window.location.host
    ws.value = new WebSocket(`${protocol}//${host}/ws?task_id=${taskId}`)

    ws.value.onmessage = (event) => {
        const data = JSON.parse(event.data) as WSMessage
        handleMessage(data)
    }
}
```

---

## Security Notes

- **Sensitive Fields**: Use `json:"-"` tag to exclude from JSON serialization
- **Passwords**: Stored encrypted with AES-256-GCM; decrypted only when needed for SSH
- **CORS**: Currently allows all origins (`CheckOrigin: func() bool { return true }`) for development
- **Input Validation**: Always validate and sanitize user input before database operations

---

## Database

- SQLite database at `data/servers.db`
- Tables: `servers`, `scripts`, `tasks`, `task_servers`
- Automatic migration on startup
- Foreign keys cascade for task_servers deletion

---

## API Endpoints

| Method | Endpoint | Description | Pagination |
|--------|----------|-------------|------------|
| GET | `/api/servers` | List servers | ✅ page, pageSize |
| POST | `/api/servers` | Create server | |
| PUT | `/api/servers/:id` | Update server | |
| DELETE | `/api/servers/:id` | Delete server | |
| POST | `/api/servers/:id/test` | Test SSH connection | |
| GET | `/api/scripts` | List scripts | ✅ page, pageSize |
| POST | `/api/scripts` | Create script | |
| PUT | `/api/scripts/:id` | Update script | |
| DELETE | `/api/scripts/:id` | Delete script | |
| GET | `/api/tasks` | List tasks | ✅ page, pageSize |
| POST | `/api/tasks` | Create and execute task | |
| DELETE | `/api/tasks/:id` | Cancel task | |
| GET | `/ws?task_id=N` | WebSocket for real-time output | |

### Pagination Response Format

All list endpoints support pagination and return:
```json
{
  "data": [...],
  "total": 100,
  "page": 1,
  "pageSize": 10
}
```

<!-- gitnexus:start -->
# GitNexus — Code Intelligence

This project is indexed by GitNexus as **node-pilot** (5017 symbols, 14834 relationships, 300 execution flows). Use the GitNexus MCP tools to understand code, assess impact, and navigate safely.

> If any GitNexus tool warns the index is stale, run `npx gitnexus analyze` in terminal first.

## Always Do

- **MUST run impact analysis before editing any symbol.** Before modifying a function, class, or method, run `gitnexus_impact({target: "symbolName", direction: "upstream"})` and report the blast radius (direct callers, affected processes, risk level) to the user.
- **MUST run `gitnexus_detect_changes()` before committing** to verify your changes only affect expected symbols and execution flows.
- **MUST warn the user** if impact analysis returns HIGH or CRITICAL risk before proceeding with edits.
- When exploring unfamiliar code, use `gitnexus_query({query: "concept"})` to find execution flows instead of grepping. It returns process-grouped results ranked by relevance.
- When you need full context on a specific symbol — callers, callees, which execution flows it participates in — use `gitnexus_context({name: "symbolName"})`.

## When Debugging

1. `gitnexus_query({query: "<error or symptom>"})` — find execution flows related to the issue
2. `gitnexus_context({name: "<suspect function>"})` — see all callers, callees, and process participation
3. `READ gitnexus://repo/node-pilot/process/{processName}` — trace the full execution flow step by step
4. For regressions: `gitnexus_detect_changes({scope: "compare", base_ref: "main"})` — see what your branch changed

## When Refactoring

- **Renaming**: MUST use `gitnexus_rename({symbol_name: "old", new_name: "new", dry_run: true})` first. Review the preview — graph edits are safe, text_search edits need manual review. Then run with `dry_run: false`.
- **Extracting/Splitting**: MUST run `gitnexus_context({name: "target"})` to see all incoming/outgoing refs, then `gitnexus_impact({target: "target", direction: "upstream"})` to find all external callers before moving code.
- After any refactor: run `gitnexus_detect_changes({scope: "all"})` to verify only expected files changed.

## Never Do

- NEVER edit a function, class, or method without first running `gitnexus_impact` on it.
- NEVER ignore HIGH or CRITICAL risk warnings from impact analysis.
- NEVER rename symbols with find-and-replace — use `gitnexus_rename` which understands the call graph.
- NEVER commit changes without running `gitnexus_detect_changes()` to check affected scope.

## Tools Quick Reference

| Tool | When to use | Command |
|------|-------------|---------|
| `query` | Find code by concept | `gitnexus_query({query: "auth validation"})` |
| `context` | 360-degree view of one symbol | `gitnexus_context({name: "validateUser"})` |
| `impact` | Blast radius before editing | `gitnexus_impact({target: "X", direction: "upstream"})` |
| `detect_changes` | Pre-commit scope check | `gitnexus_detect_changes({scope: "staged"})` |
| `rename` | Safe multi-file rename | `gitnexus_rename({symbol_name: "old", new_name: "new", dry_run: true})` |
| `cypher` | Custom graph queries | `gitnexus_cypher({query: "MATCH ..."})` |

## Impact Risk Levels

| Depth | Meaning | Action |
|-------|---------|--------|
| d=1 | WILL BREAK — direct callers/importers | MUST update these |
| d=2 | LIKELY AFFECTED — indirect deps | Should test |
| d=3 | MAY NEED TESTING — transitive | Test if critical path |

## Resources

| Resource | Use for |
|----------|---------|
| `gitnexus://repo/node-pilot/context` | Codebase overview, check index freshness |
| `gitnexus://repo/node-pilot/clusters` | All functional areas |
| `gitnexus://repo/node-pilot/processes` | All execution flows |
| `gitnexus://repo/node-pilot/process/{name}` | Step-by-step execution trace |

## Self-Check Before Finishing

Before completing any code modification task, verify:
1. `gitnexus_impact` was run for all modified symbols
2. No HIGH/CRITICAL risk warnings were ignored
3. `gitnexus_detect_changes()` confirms changes match expected scope
4. All d=1 (WILL BREAK) dependents were updated

## Keeping the Index Fresh

After committing code changes, the GitNexus index becomes stale. Re-run analyze to update it:

```bash
npx gitnexus analyze
```

If the index previously included embeddings, preserve them by adding `--embeddings`:

```bash
npx gitnexus analyze --embeddings
```

To check whether embeddings exist, inspect `.gitnexus/meta.json` — the `stats.embeddings` field shows the count (0 means no embeddings). **Running analyze without `--embeddings` will delete any previously generated embeddings.**

> Claude Code users: A PostToolUse hook handles this automatically after `git commit` and `git merge`.

## CLI

| Task | Read this skill file |
|------|---------------------|
| Understand architecture / "How does X work?" | `.claude/skills/gitnexus/gitnexus-exploring/SKILL.md` |
| Blast radius / "What breaks if I change X?" | `.claude/skills/gitnexus/gitnexus-impact-analysis/SKILL.md` |
| Trace bugs / "Why is X failing?" | `.claude/skills/gitnexus/gitnexus-debugging/SKILL.md` |
| Rename / extract / split / refactor | `.claude/skills/gitnexus/gitnexus-refactoring/SKILL.md` |
| Tools, resources, schema reference | `.claude/skills/gitnexus/gitnexus-guide/SKILL.md` |
| Index, status, clean, wiki CLI commands | `.claude/skills/gitnexus/gitnexus-cli/SKILL.md` |

<!-- gitnexus:end -->
