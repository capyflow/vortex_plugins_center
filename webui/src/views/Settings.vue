<template>
  <div class="settings-page">
    <h1>系统设置</h1>
    <el-form label-width="120px" class="settings-form">
      <el-form-item label="API 地址">
        <el-input v-model="settings.apiBaseUrl" placeholder="http://localhost:8080/api/v1" />
      </el-form-item>
      <el-form-item label="自动刷新">
        <el-switch v-model="settings.autoRefresh" />
      </el-form-item>
      <el-form-item label="刷新间隔">
        <el-slider v-model="settings.refreshInterval" :min="5000" :max="60000" :step="1000" show-stops />
        <span>{{ settings.refreshInterval / 1000 }}秒</span>
      </el-form-item>
      <el-form-item label="暗黑模式">
        <el-switch v-model="settings.darkMode" @change="settings.toggleDarkMode()" />
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { useSettingsStore } from '@/stores/settings'

const settings = useSettingsStore()
</script>

<style scoped>
.settings-page {
  padding: 24px;
}

.settings-page h1 {
  margin-bottom: 24px;
}

.settings-form {
  max-width: 600px;
}

:deep(.el-form-item__label) {
  color: #fff;
}

:deep(.el-slider__runway) {
  background: rgba(255, 255, 255, 0.1);
}

:deep(.el-slider__button) {
  background: #667eea;
}
</style>
