<template>
  <div class="plugin-detail">
    <div class="detail-header">
      <div class="header-info">
        <h2>{{ plugin.name }}</h2>
        <el-tag type="info">{{ plugin.version }}</el-tag>
      </div>
      <div class="header-actions">
        <el-tag :type="statusType">{{ statusText }}</el-tag>
      </div>
    </div>

    <div class="detail-section">
      <h4>基本信息</h4>
      <p class="description">{{ plugin.description || '暂无描述' }}</p>
      <div class="info-grid">
        <div class="info-item">
          <span class="label">端点</span>
          <span class="value">{{ plugin.endpoint }}</span>
        </div>
        <div class="info-item">
          <span class="label">状态</span>
          <span class="value">
            <HealthIndicator :status="plugin.health?.status" />
            {{ plugin.health?.status === 'healthy' ? '健康' : '异常' }}
            ({{ plugin.health?.latency }}ms)
          </span>
        </div>
        <div class="info-item">
          <span class="label">更新时间</span>
          <span class="value">{{ formatDate(plugin.updated_at) }}</span>
        </div>
      </div>
    </div>

    <div class="detail-section">
      <h4>方法列表 ({{ plugin.methods?.length || 0 }})</h4>
      <el-collapse v-model="activeMethods">
        <el-collapse-item
          v-for="method in plugin.methods"
          :key="method.name"
          :name="method.name"
          :title="method.name"
        >
          <MethodPanel :plugin-id="plugin.id" :method="method" />
        </el-collapse-item>
      </el-collapse>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { Plugin } from '@/types/plugin'
import HealthIndicator from './HealthIndicator.vue'
import MethodPanel from './MethodPanel.vue'

const props = defineProps<{
  plugin: Plugin
}>()

const activeMethods = ref<string[]>([])

const statusType = computed(() => {
  switch (props.plugin.status) {
    case 'active': return 'success'
    case 'inactive': return 'info'
    case 'error': return 'danger'
    default: return 'info'
  }
})

const statusText = computed(() => {
  switch (props.plugin.status) {
    case 'active': return '运行中'
    case 'inactive': return '已停止'
    case 'error': return '异常'
    default: return '未知'
  }
})

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleString('zh-CN')
}
</script>

<style scoped>
.plugin-detail {
  padding: 20px;
  color: #fff;
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding-bottom: 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.header-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-info h2 {
  margin: 0;
  font-size: 24px;
}

.detail-section {
  margin-bottom: 24px;
}

.detail-section h4 {
  color: #fff;
  margin-bottom: 16px;
  font-size: 16px;
}

.description {
  color: #aaa;
  margin-bottom: 16px;
  line-height: 1.6;
}

.info-grid {
  display: grid;
  gap: 12px;
}

.info-item {
  display: flex;
  gap: 12px;
}

.info-item .label {
  color: #888;
  min-width: 80px;
}

.info-item .value {
  color: #fff;
  display: flex;
  align-items: center;
  gap: 8px;
}

:deep(.el-collapse) {
  border-color: rgba(255, 255, 255, 0.1);
}

:deep(.el-collapse-item__header) {
  background: transparent;
  color: #fff;
  border-color: rgba(255, 255, 255, 0.1);
}

:deep(.el-collapse-item__wrap) {
  background: transparent;
  border-color: rgba(255, 255, 255, 0.1);
}

:deep(.el-collapse-item__content) {
  color: #fff;
  padding-bottom: 20px;
}
</style>
