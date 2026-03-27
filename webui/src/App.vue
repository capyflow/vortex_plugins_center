<template>
  <router-view />
</template>

<script setup lang="ts">
</script>

<style>
/* 全局样式覆盖 */
html, body {
  margin: 0;
  padding: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

#app {
  min-height: 100vh;
}

/* Element Plus 暗黑主题覆盖 */
.dark {
  --el-bg-color: transparent;
  --el-bg-color-page: transparent;
  --el-text-color-primary: #fff;
  --el-text-color-regular: #fff;
  --el-text-color-secondary: #aaa;
  --el-border-color: rgba(255, 255, 255, 0.1);
  --el-fill-color-blank: rgba(255, 255, 255, 0.05);
  --el-fill-color-light: rgba(255, 255, 255, 0.1);
}

/* 滚动条样式 */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.1);
}

::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.3);
}
</style>
