# 创建登录页面

## 任务描述

创建具有深色主题、科幻控制台风格的登录页面，采用 Glassmorphism 和霓虹灯效果。

## 详细说明

### 1. 创建 Login.vue

创建 `frontend/src/views/Login.vue`：

```vue
<template>
    <div class="login-container">
        <!-- 背景动画 -->
        <div class="bg-animation">
            <div class="particle" v-for="i in 50" :key="i" :style="particleStyle(i)"></div>
        </div>
        
        <div class="login-panel">
            <div class="logo-section">
                <h1 class="logo-text">NodePilot</h1>
                <p class="logo-subtitle">批量服务器管理平台</p>
            </div>

            <form @submit.prevent="handleLogin" class="login-form">
                <div class="input-group">
                    <div class="input-wrapper">
                        <span class="input-icon">❯</span>
                        <input 
                            v-model="username" 
                            type="text" 
                            placeholder="用户名"
                            autocomplete="username"
                            required
                        />
                        <div class="input-glow"></div>
                    </div>
                </div>

                <div class="input-group">
                    <div class="input-wrapper">
                        <span class="input-icon">***</span>
                        <input 
                            v-model="password" 
                            :type="showPassword ? 'text' : 'password'" 
                            placeholder="密码"
                            autocomplete="current-password"
                            required
                        />
                        <button type="button" class="toggle-password" @click="showPassword = !showPassword">
                            {{ showPassword ? '👁' : '👁‍🗨' }}
                        </button>
                        <div class="input-glow"></div>
                    </div>
                </div>

                <div v-if="error" class="error-message">
                    {{ error }}
                </div>

                <button type="submit" class="login-button" :disabled="loading">
                    <span v-if="!loading">接入系统</span>
                    <span v-else class="loading-dots">连接中<span>.</span><span>.</span><span>.</span></span>
                </button>
            </form>

            <div class="system-info">
                <span class="version">v1.0.0</span>
                <span class="status-indicator"></span>
                <span class="status-text">系统就绪</span>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

const router = useRouter();
const authStore = useAuthStore();

const username = ref('');
const password = ref('');
const showPassword = ref(false);
const loading = ref(false);
const error = ref('');

function particleStyle(i: number) {
    return {
        left: `${Math.random() * 100}%`,
        animationDelay: `${Math.random() * 10}s`,
        animationDuration: `${10 + Math.random() * 10}s`,
    };
}

async function handleLogin() {
    error.value = '';
    loading.value = true;
    
    try {
        const success = await authStore.login({
            username: username.value,
            password: password.value,
        });
        
        if (success) {
            router.push('/');
        } else {
            error.value = '用户名或密码错误';
        }
    } catch (e) {
        error.value = '登录失败，请检查网络连接';
    } finally {
        loading.value = false;
    }
}
</script>

<style scoped>
.login-container {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #0a0a0a 0%, #1a1a2e 50%, #0a0a0a 100%);
    position: relative;
    overflow: hidden;
}

/* 粒子动画背景 */
.bg-animation {
    position: absolute;
    inset: 0;
    overflow: hidden;
}

.particle {
    position: absolute;
    width: 2px;
    height: 2px;
    background: rgba(102, 126, 234, 0.6);
    border-radius: 50%;
    animation: float linear infinite;
    box-shadow: 0 0 10px rgba(102, 126, 234, 0.8);
}

@keyframes float {
    0% { transform: translateY(100vh) scale(0); opacity: 0; }
    10% { opacity: 1; }
    90% { opacity: 1; }
    100% { transform: translateY(-100vh) scale(1); opacity: 0; }
}

/* 登录面板 - 玻璃拟态效果 */
.login-panel {
    background: rgba(20, 20, 40, 0.6);
    backdrop-filter: blur(20px);
    border: 1px solid rgba(102, 126, 234, 0.3);
    border-radius: 20px;
    padding: 48px;
    width: 100%;
    max-width: 420px;
    position: relative;
    box-shadow: 
        0 0 40px rgba(102, 126, 234, 0.2),
        inset 0 0 60px rgba(102, 126, 234, 0.05);
}

.login-panel::before {
    content: '';
    position: absolute;
    inset: -1px;
    border-radius: 20px;
    padding: 1px;
    background: linear-gradient(135deg, rgba(102, 126, 234, 0.5), transparent, rgba(118, 75, 162, 0.5));
    -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
    mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
    -webkit-mask-composite: xor;
    mask-composite: exclude;
}

/* Logo 区域 */
.logo-section {
    text-align: center;
    margin-bottom: 40px;
}

.logo-text {
    font-size: 2.5rem;
    font-weight: bold;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    text-shadow: 0 0 30px rgba(102, 126, 234, 0.5);
    margin: 0;
}

.logo-subtitle {
    color: rgba(255, 255, 255, 0.6);
    font-size: 0.9rem;
    margin-top: 8px;
}

/* 输入框组 */
.input-group {
    margin-bottom: 24px;
}

.input-wrapper {
    position: relative;
    display: flex;
    align-items: center;
}

.input-icon {
    position: absolute;
    left: 16px;
    color: rgba(102, 126, 234, 0.8);
    font-family: monospace;
    font-size: 0.9rem;
}

.input-wrapper input {
    width: 100%;
    padding: 16px 16px 16px 48px;
    background: rgba(10, 10, 30, 0.8);
    border: 2px solid rgba(102, 126, 234, 0.3);
    border-radius: 12px;
    color: #fff;
    font-size: 1rem;
    transition: all 0.3s ease;
}

.input-wrapper input:focus {
    outline: none;
    border-color: rgba(102, 126, 234, 0.8);
    box-shadow: 0 0 20px rgba(102, 126, 234, 0.3);
}

.input-wrapper input::placeholder {
    color: rgba(255, 255, 255, 0.3);
}

/* 霓虹灯光晕动画 */
.input-glow {
    position: absolute;
    inset: -2px;
    border-radius: 14px;
    background: transparent;
    opacity: 0;
    transition: opacity 0.3s ease;
    pointer-events: none;
}

.input-wrapper input:focus + .input-glow {
    opacity: 1;
    animation: glow-pulse 2s ease-in-out infinite;
}

@keyframes glow-pulse {
    0%, 100% { box-shadow: 0 0 20px rgba(102, 126, 234, 0.3); }
    50% { box-shadow: 0 0 40px rgba(102, 126, 234, 0.6); }
}

.toggle-password {
    position: absolute;
    right: 16px;
    background: none;
    border: none;
    color: rgba(255, 255, 255, 0.5);
    cursor: pointer;
    font-size: 1.1rem;
}

/* 错误消息 */
.error-message {
    background: rgba(255, 82, 82, 0.1);
    border: 1px solid rgba(255, 82, 82, 0.3);
    border-radius: 8px;
    padding: 12px;
    color: #ff5252;
    text-align: center;
    margin-bottom: 20px;
    font-size: 0.9rem;
}

/* 登录按钮 */
.login-button {
    width: 100%;
    padding: 16px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border: none;
    border-radius: 12px;
    color: #fff;
    font-size: 1.1rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s ease;
    position: relative;
    overflow: hidden;
}

.login-button:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 10px 30px rgba(102, 126, 234, 0.4);
}

.login-button:disabled {
    opacity: 0.7;
    cursor: not-allowed;
}

/* 流光动效 */
.login-button::before {
    content: '';
    position: absolute;
    top: 0;
    left: -100%;
    width: 100%;
    height: 100%;
    background: linear-gradient(90deg, transparent, rgba(255,255,255,0.2), transparent);
    transition: left 0.5s ease;
}

.login-button:hover::before {
    left: 100%;
}

.loading-dots span {
    animation: blink 1.4s infinite both;
}
.loading-dots span:nth-child(2) { animation-delay: 0.2s; }
.loading-dots span:nth-child(3) { animation-delay: 0.4s; }

@keyframes blink {
    0%, 80%, 100% { opacity: 0; }
    40% { opacity: 1; }
}

/* 系统信息 */
.system-info {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 8px;
    margin-top: 32px;
    font-size: 0.8rem;
    color: rgba(255, 255, 255, 0.4);
}

.status-indicator {
    width: 8px;
    height: 8px;
    background: #4ade80;
    border-radius: 50%;
    animation: pulse 2s infinite;
}

@keyframes pulse {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.5; }
}

/* 响应式 */
@media (max-width: 480px) {
    .login-panel {
        margin: 20px;
        padding: 32px 24px;
    }
    
    .logo-text {
        font-size: 2rem;
    }
}
</style>
```

## 输入

- 需求文档 `05-添加授权登录.md`
- 现有的 `frontend/src/views/*.vue` 页面组件

## 输出

- 新建 `frontend/src/views/Login.vue` - 登录页面

## 依赖

- 05-frontend-auth-api-store.md (Auth Store 实现)

## 验收标准

- [ ] 深色主题背景，带动态粒子效果
- [ ] 玻璃拟态登录面板，边框发光效果
- [ ] 霓虹灯风格输入框，聚焦时光晕扩散
- [ ] 登录按钮渐变色，悬停流光动效
- [ ] 响应式设计适配移动端
- [ ] 显示系统状态指示器
- [ ] 登录失败显示错误消息
- [ ] 登录成功跳转到首页
