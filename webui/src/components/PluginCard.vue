<template>
  <div
    class="plugin-card"
    :class="{ inactive: plugin.status === 'inactive' }"
    @click="emit('click')"
  >
    <div class="card-header">
      <div class="plugin-info">
        <h3 class="plugin-name">{{ plugin.name }}</h3>
        <el-tag size="small" type="info" class="version-tag">{{ plugin.version }}</el-tag>
      </div>
      <HealthIndicator :status="plugin.health?.status" />
    </div>
    
    <p class="plugin-desc">{{ plugin.description || '暂无描述' }}</p>
    
    <div class="card-footer">
      <div class="meta-info">
        <el-tag size="small" :type="statusType">{{ statusText }}</el-tag>
        <span class="method-count">{{ plugin.methods?.length || 0 }} 个方法</span>
      </div>
      <div v-if="plugin.health?.latency" class="latency">
        <el-icon><Timer /></el-icon>
        <span>{{ plugin.health.latency }}ms</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Timer } from '@element-plus/icons-vue'
import type { Plugin } from '@/types/plugin'
import HealthIndicator from './HealthIndicator.vue'

const props = defineProps<{
  plugin: Plugin
}>()

const emit = defineEmits<{
  click: []
}>()

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
</script>

<style scoped>
.plugin-card {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.3s;
}

.plugin-card:hover {
  background: rgba(255, 255, 255, 0.1);
  transform: translateY(-4px);
  box-shadow: 0 10px 30px rgba(102, 126, 234, 0.3);
}

.plugin-card.inactive {
  opacity: 0.6;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: start;
  margin-bottom: 12px;
}

.plugin-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.plugin-name {
  font-size: 18px;
  font-weight: 600;
  color: #fff;
  margin: 0;
}

.version-tag {
  background: rgba(102, 126, 234, 0.2) !important;
  color: #667eea !important;
  border: none !important;
}

.plugin-desc {
  font-size: 14px;
  color: #aaa;
  margin: 0 0 16px 0;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.meta-info {
  display: flex;
  gap: 10px;
  align-items: center;
}

.method-count {
  font-size: 12px;
  color: #888;
}

.latency {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #2ed573;
}
</style>
