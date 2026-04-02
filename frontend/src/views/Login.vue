<template>
    <div class="login-container">
        <div class="bg-animation">
            <!-- Animated grid background -->
            <div class="grid-bg"></div>
            <!-- Glowing orbs -->
            <div class="orb orb-1"></div>
            <div class="orb orb-2"></div>
            <div class="orb orb-3"></div>
            <!-- Circuit lines -->
            <svg class="circuit-svg" viewBox="0 0 100 100" preserveAspectRatio="none">
                <defs>
                    <linearGradient id="circuitGrad" x1="0%" y1="0%" x2="100%" y2="100%">
                        <stop offset="0%" style="stop-color:#00d4ff;stop-opacity:0.6" />
                        <stop offset="100%" style="stop-color:#667eea;stop-opacity:0.2" />
                    </linearGradient>
                </defs>
                <g class="circuit-lines">
                    <line x1="0" y1="20" x2="40" y2="20" />
                    <line x1="40" y1="20" x2="40" y2="50" />
                    <line x1="40" y1="50" x2="70" y2="50" />
                    <line x1="70" y1="50" x2="70" y2="80" />
                    <line x1="70" y1="80" x2="100" y2="80" />
                    <circle cx="40" cy="20" r="2" fill="#00d4ff" />
                    <circle cx="40" cy="50" r="2" fill="#00d4ff" />
                    <circle cx="70" cy="50" r="2" fill="#667eea" />
                    <circle cx="70" cy="80" r="2" fill="#667eea" />
                </g>
                <g class="circuit-lines circuit-lines-2">
                    <line x1="100" y1="10" x2="60" y2="10" />
                    <line x1="60" y1="10" x2="60" y2="40" />
                    <line x1="60" y1="40" x2="30" y2="40" />
                    <line x1="30" y1="40" x2="30" y2="70" />
                    <circle cx="60" cy="10" r="2" fill="#764ba2" />
                    <circle cx="60" cy="40" r="2" fill="#764ba2" />
                    <circle cx="30" cy="40" r="2" fill="#00d4ff" />
                </g>
                <g class="circuit-lines circuit-lines-3">
                    <line x1="0" y1="60" x2="25" y2="60" />
                    <line x1="25" y1="60" x2="25" y2="85" />
                    <line x1="25" y1="85" x2="55" y2="85" />
                    <line x1="55" y1="85" x2="55" y2="95" />
                    <circle cx="25" cy="60" r="2" fill="#00d4ff" />
                    <circle cx="55" cy="85" r="2" fill="#667eea" />
                </g>
            </svg>
            <!-- Data particles -->
            <div class="particle particle-cyan" v-for="i in 15" :key="'c'+i" :style="particleStyle(i, 'cyan')"></div>
            <div class="particle particle-purple" v-for="i in 15" :key="'p'+i" :style="particleStyle(i, 'purple')"></div>
            <div class="particle particle-blue" v-for="i in 10" :key="'b'+i" :style="particleStyle(i, 'blue')"></div>
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

function particleStyle(_i: number, color: string) {
    const colors: Record<string, string> = {
        cyan: '#00d4ff',
        purple: '#667eea',
        blue: '#764ba2'
    };
    return {
        left: `${Math.random() * 100}%`,
        top: `${Math.random() * 100}%`,
        background: colors[color] || colors.cyan,
        boxShadow: `0 0 10px ${colors[color]}, 0 0 20px ${colors[color]}`,
        animationDelay: `${Math.random() * 8}s`,
        animationDuration: `${8 + Math.random() * 12}s`,
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
    background: linear-gradient(135deg, #0d1b2a 0%, #1b263b 40%, #0d1b2a 70%, #1a1a2e 100%);
    position: relative;
    overflow: hidden;
}

.bg-animation {
    position: absolute;
    inset: 0;
    overflow: hidden;
}

/* Animated grid background */
.grid-bg {
    position: absolute;
    inset: 0;
    background-image: 
        linear-gradient(rgba(0, 212, 255, 0.03) 1px, transparent 1px),
        linear-gradient(90deg, rgba(0, 212, 255, 0.03) 1px, transparent 1px);
    background-size: 50px 50px;
    animation: gridMove 20s linear infinite;
}

@keyframes gridMove {
    0% { transform: translate(0, 0); }
    100% { transform: translate(50px, 50px); }
}

/* Glowing orbs */
.orb {
    position: absolute;
    border-radius: 50%;
    filter: blur(80px);
    opacity: 0.4;
    animation: orbFloat 8s ease-in-out infinite;
}

.orb-1 {
    width: 400px;
    height: 400px;
    background: radial-gradient(circle, rgba(0, 212, 255, 0.3) 0%, transparent 70%);
    top: -100px;
    right: -100px;
    animation-delay: 0s;
}

.orb-2 {
    width: 350px;
    height: 350px;
    background: radial-gradient(circle, rgba(102, 126, 234, 0.3) 0%, transparent 70%);
    bottom: -80px;
    left: -80px;
    animation-delay: -3s;
}

.orb-3 {
    width: 300px;
    height: 300px;
    background: radial-gradient(circle, rgba(118, 75, 162, 0.25) 0%, transparent 70%);
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    animation-delay: -5s;
}

@keyframes orbFloat {
    0%, 100% { transform: translate(0, 0) scale(1); }
    33% { transform: translate(30px, -30px) scale(1.1); }
    66% { transform: translate(-20px, 20px) scale(0.95); }
}

.orb-3 {
    animation-name: orbFloatCenter;
}

@keyframes orbFloatCenter {
    0%, 100% { transform: translate(-50%, -50%) scale(1); }
    50% { transform: translate(-50%, -50%) scale(1.2); }
}

/* Circuit SVG */
.circuit-svg {
    position: absolute;
    width: 100%;
    height: 100%;
    opacity: 0.6;
}

.circuit-lines {
    stroke: url(#circuitGrad);
    stroke-width: 0.3;
    fill: none;
    stroke-dasharray: 200;
    stroke-dashoffset: 200;
    animation: circuitDraw 8s ease-in-out infinite;
}

.circuit-lines-2 {
    animation-delay: -2s;
    stroke-dashoffset: 200;
}

.circuit-lines-3 {
    animation-delay: -4s;
    stroke-dashoffset: 200;
}

@keyframes circuitDraw {
    0% { stroke-dashoffset: 200; opacity: 0; }
    10% { opacity: 1; }
    90% { opacity: 1; }
    100% { stroke-dashoffset: 0; opacity: 0; }
}

/* Data particles */
.particle {
    position: absolute;
    width: 3px;
    height: 3px;
    border-radius: 50%;
    animation: dataStream linear infinite;
}

.particle-cyan {
    box-shadow: 0 0 6px #00d4ff, 0 0 12px #00d4ff;
}

.particle-purple {
    box-shadow: 0 0 6px #667eea, 0 0 12px #667eea;
}

.particle-blue {
    box-shadow: 0 0 6px #764ba2, 0 0 12px #764ba2;
}

@keyframes dataStream {
    0% { 
        transform: translateY(100vh) scale(0);
        opacity: 0;
    }
    10% { opacity: 1; transform: scale(1); }
    90% { opacity: 0.8; }
    100% { 
        transform: translateY(-20vh) scale(0.5);
        opacity: 0;
    }
}

.login-panel {
    background: rgba(13, 27, 42, 0.75);
    backdrop-filter: blur(20px);
    border: 1px solid rgba(0, 212, 255, 0.2);
    border-radius: 20px;
    padding: 48px;
    width: 100%;
    max-width: 420px;
    position: relative;
    box-shadow: 
        0 0 60px rgba(0, 212, 255, 0.15),
        0 0 100px rgba(102, 126, 234, 0.1),
        inset 0 0 60px rgba(0, 212, 255, 0.03);
}

.login-panel::before {
    content: '';
    position: absolute;
    inset: -1px;
    border-radius: 20px;
    padding: 1px;
    background: linear-gradient(135deg, rgba(0, 212, 255, 0.5), rgba(102, 126, 234, 0.3), rgba(118, 75, 162, 0.5));
    -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
    mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
    -webkit-mask-composite: xor;
    mask-composite: exclude;
}

.logo-section {
    text-align: center;
    margin-bottom: 40px;
}

.logo-text {
    font-size: 2.5rem;
    font-weight: bold;
    background: linear-gradient(135deg, #00d4ff 0%, #667eea 50%, #764ba2 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    text-shadow: 0 0 40px rgba(0, 212, 255, 0.5);
    margin: 0;
}

.logo-subtitle {
    color: rgba(255, 255, 255, 0.6);
    font-size: 0.9rem;
    margin-top: 8px;
}

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
    color: rgba(0, 212, 255, 0.8);
    font-family: monospace;
    font-size: 0.9rem;
}

.input-wrapper input {
    width: 100%;
    padding: 16px 16px 16px 48px;
    background: rgba(13, 27, 42, 0.9);
    border: 2px solid rgba(0, 212, 255, 0.2);
    border-radius: 12px;
    color: #fff;
    font-size: 1rem;
    transition: all 0.3s ease;
}

.input-wrapper input:focus {
    outline: none;
    border-color: rgba(0, 212, 255, 0.6);
    box-shadow: 0 0 25px rgba(0, 212, 255, 0.25), inset 0 0 15px rgba(0, 212, 255, 0.05);
}

.input-wrapper input::placeholder {
    color: rgba(255, 255, 255, 0.3);
}

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
    0%, 100% { box-shadow: 0 0 20px rgba(0, 212, 255, 0.3); }
    50% { box-shadow: 0 0 40px rgba(0, 212, 255, 0.6); }
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

.login-button {
    width: 100%;
    padding: 16px;
    background: linear-gradient(135deg, #00d4ff 0%, #667eea 50%, #764ba2 100%);
    border: none;
    border-radius: 12px;
    color: #fff;
    font-size: 1.1rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s ease;
    position: relative;
    overflow: hidden;
    box-shadow: 0 4px 20px rgba(0, 212, 255, 0.3);
}

.login-button:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 8px 35px rgba(0, 212, 255, 0.4), 0 0 60px rgba(102, 126, 234, 0.2);
}

.login-button:disabled {
    opacity: 0.7;
    cursor: not-allowed;
}

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

.system-info {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 8px;
    margin-top: 32px;
    font-size: 0.8rem;
    color: rgba(0, 212, 255, 0.5);
}

.status-indicator {
    width: 8px;
    height: 8px;
    background: #00d4ff;
    border-radius: 50%;
    box-shadow: 0 0 10px #00d4ff;
    animation: pulse 2s infinite;
}

@keyframes pulse {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.5; }
}

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
