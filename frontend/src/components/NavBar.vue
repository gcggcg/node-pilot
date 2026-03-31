<template>
    <div class="navbar">
        <div class="logo">NodePilot</div>
        <nav class="nav-links">
            <router-link to="/servers" class="nav-link">服务器</router-link>
            <router-link to="/scripts" class="nav-link">脚本</router-link>
            <router-link to="/tasks" class="nav-link">任务</router-link>
        </nav>
    </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
.navbar {
    height: 60px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    display: flex;
    align-items: center;
    padding: 0 24px;
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    z-index: 100;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.logo {
    font-size: 1.5rem;
    font-weight: bold;
    color: white;
    margin-right: 48px;
}

.nav-links {
    display: flex;
    gap: 24px;
}

.nav-link {
    color: rgba(255, 255, 255, 0.85);
    text-decoration: none;
    font-weight: 500;
    padding: 8px 16px;
    border-radius: 6px;
    transition: all 0.2s;
}

.nav-link:hover {
    color: white;
    background: rgba(255, 255, 255, 0.1);
}

.nav-link.router-link-active {
    color: white;
    background: rgba(255, 255, 255, 0.2);
}
</style>
