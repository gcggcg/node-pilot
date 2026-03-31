<template>
    <div id="app">
        <NavBar />
        <main class="main-content">
            <router-view />
        </main>
    </div>
</template>

<script setup lang="ts">
import NavBar from '@/components/NavBar.vue';
</script>

<style>
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    background: #f5f6f8;
    color: #333;
}

#app {
    min-height: 100vh;
}

.main-content {
    padding-top: 60px;
    min-height: calc(100vh - 60px);
}
</style>
