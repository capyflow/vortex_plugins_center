<template>
  <div class="app-layout" :class="{ dark: settings.darkMode }">
    <el-container class="layout-container">
      <el-aside width="220px" class="sidebar">
        <div class="logo">
          <el-icon :size="28" class="logo-icon"><Connection /></el-icon>
          <span class="logo-text">Plugin Center</span>
        </div>
        <el-menu
          :default-active="$route.path"
          router
          class="sidebar-menu"
          :collapse="false"
          background-color="transparent"
          text-color="#fff"
          active-text-color="#667eea"
        >
          <el-menu-item index="/">
            <el-icon><Grid /></el-icon>
            <span>插件列表</span>
          </el-menu-item>
          <el-menu-item index="/history">
            <el-icon><Clock /></el-icon>
            <span>调用历史</span>
          </el-menu-item>
          <el-menu-item index="/settings">
            <el-icon><Setting /></el-icon>
            <span>系统设置</span>
          </el-menu-item>
        </el-menu>
        <div class="sidebar-footer">
          <el-button circle @click="settings.toggleDarkMode()" class="theme-btn">
            <el-icon><component :is="settings.darkMode ? Sunny : Moon" /></el-icon>
          </el-button>
        </div>
      </el-aside>
      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Grid, Clock, Setting, Connection, Moon, Sunny } from '@element-plus/icons-vue'
import { useSettingsStore } from '@/stores/settings'

const route = useRoute()
const settings = useSettingsStore()

onMounted(() => {
  settings.initTheme()
})
</script>

<style scoped>
.app-layout {
  min-height: 100vh;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
  color: #fff;
}

.layout-container {
  min-height: 100vh;
}

.sidebar {
  background: rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(10px);
  border-right: 1px solid rgba(255, 255, 255, 0.1);
}

.logo {
  padding: 24px 20px;
  display: flex;
  align-items: center;
  gap: 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.logo-icon {
  color: #667eea;
}

.logo-text {
  font-size: 20px;
  font-weight: 600;
  background: linear-gradient(135deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.sidebar-menu {
  border-right: none !important;
}

.sidebar-footer {
  padding: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  justify-content: center;
}

.theme-btn {
  background: rgba(255, 255, 255, 0.1) !important;
  border-color: rgba(255, 255, 255, 0.2) !important;
  color: #fff !important;
}

.theme-btn:hover {
  background: rgba(255, 255, 255, 0.2) !important;
}

.main-content {
  padding: 0;
  background: transparent;
}

:deep(.el-menu-item) {
  margin: 4px 12px;
  border-radius: 8px;
  transition: all 0.3s;
}

:deep(.el-menu-item:hover) {
  background: rgba(102, 126, 234, 0.1) !important;
}

:deep(.el-menu-item.is-active) {
  background: rgba(102, 126, 234, 0.2) !important;
}
</style>
