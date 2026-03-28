<template>
  <div class="plugin-card" :class="{ 'is-active': plugin.status === 'active', 'is-error': plugin.status === 'error' }">
    <div class="plugin-header">
      <div class="plugin-icon">
        <span class="icon">🔌</span>
      </div>
      <div class="plugin-info">
        <h3 class="plugin-name">{{ plugin.name }}</h3>
        <span class="plugin-version">v{{ plugin.version }}</span>
      </div>
      <div class="plugin-status" :class="plugin.status">
        <span class="status-dot"></span>
        {{ statusText }}
      </div>
    </div>

    <p v-if="plugin.summary" class="plugin-summary">{{ plugin.summary }}</p>

    <div class="plugin-meta">
      <div class="meta-item">
        <span class="meta-label">输出类型:</span>
        <span class="meta-value">{{ plugin.outputs || 'unknown' }}</span>
      </div>
      <div v-if="plugin.health" class="meta-item">
        <span class="meta-label">延迟:</span>
        <span class="meta-value" :class="getLatencyClass(plugin.health.latency)">
          {{ plugin.health.latency }}ms
        </span>
      </div>
      <div class="meta-item">
        <span class="meta-label">参数:</span>
        <span class="meta-value">{{ paramCount }} 个</span>
      </div>
    </div>

    <div class="plugin-params-preview" v-if="paramCount > 0">
      <div class="params-tags">
        <span
          v-for="(def, name) in plugin.params"
          :key="name"
          class="param-tag"
          :class="`type-${def.type}`"
        >
          {{ name }}: {{ def.type }}
          <span v-if="def.required" class="required-mark">*</span>
        </span>
      </div>
    </div>

    <div class="plugin-actions">
      <button class="btn btn-primary" @click="$emit('execute', plugin)">
        <span class="btn-icon">▶</span>
        执行
      </button>
      <button class="btn btn-secondary" @click="$emit('detail', plugin)">
        详情
      </button>
      <button
        v-if="plugin.status === 'active'"
        class="btn btn-danger"
        @click="$emit('unregister', plugin)"
      >
        注销
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Plugin } from '@/types/plugin'

interface Props {
  plugin: Plugin
}

const props = defineProps<Props>()
defineEmits<{
  execute: [plugin: Plugin]
  detail: [plugin: Plugin]
  unregister: [plugin: Plugin]
}>()

const statusText = computed(() => {
  const map: Record<string, string> = {
    active: '运行中',
    inactive: '已停止',
    error: '异常',
  }
  return map[props.plugin.status] || props.plugin.status
})

const paramCount = computed(() => {
  return Object.keys(props.plugin.params || {}).length
})

function getLatencyClass(latency: number): string {
  if (latency < 50) return 'latency-good'
  if (latency < 200) return 'latency-normal'
  return 'latency-slow'
}
</script>

<style scoped>
.plugin-card {
  background: var(--bg-secondary);
  border-radius: 12px;
  padding: 20px;
  border: 2px solid transparent;
  transition: all 0.3s ease;
}

.plugin-card:hover {
  border-color: var(--primary-color);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.plugin-card.is-error {
  border-color: #ef4444;
}

.plugin-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.plugin-icon {
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.plugin-info {
  flex: 1;
}

.plugin-name {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
  color: var(--text-primary);
}

.plugin-version {
  font-size: 12px;
  color: var(--text-secondary);
  background: var(--bg-tertiary);
  padding: 2px 8px;
  border-radius: 4px;
}

.plugin-status {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  padding: 4px 10px;
  border-radius: 20px;
  background: var(--bg-tertiary);
}

.plugin-status.active {
  background: #dcfce7;
  color: #166534;
}

.plugin-status.inactive {
  background: #f3f4f6;
  color: #6b7280;
}

.plugin-status.error {
  background: #fee2e2;
  color: #991b1b;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: currentColor;
}

.plugin-summary {
  color: var(--text-secondary);
  font-size: 14px;
  margin: 0 0 12px 0;
  line-height: 1.5;
}

.plugin-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  margin-bottom: 12px;
  padding: 12px;
  background: var(--bg-tertiary);
  border-radius: 8px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
}

.meta-label {
  color: var(--text-secondary);
}

.meta-value {
  color: var(--text-primary);
  font-weight: 500;
}

.latency-good {
  color: #16a34a;
}

.latency-normal {
  color: #ca8a04;
}

.latency-slow {
  color: #dc2626;
}

.plugin-params-preview {
  margin-bottom: 16px;
}

.params-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.param-tag {
  font-size: 12px;
  padding: 4px 10px;
  border-radius: 6px;
  background: var(--bg-tertiary);
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  gap: 4px;
}

.param-tag.type-string {
  background: #dbeafe;
  color: #1e40af;
}

.param-tag.type-number {
  background: #dcfce7;
  color: #166534;
}

.param-tag.type-file {
  background: #fce7f3;
  color: #9d174d;
}

.param-tag.type-boolean {
  background: #f3e8ff;
  color: #6b21a8;
}

.param-tag.type-select {
  background: #fef3c7;
  color: #92400e;
}

.required-mark {
  color: #ef4444;
  font-weight: bold;
}

.plugin-actions {
  display: flex;
  gap: 8px;
}

.btn {
  flex: 1;
  padding: 10px 16px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}

.btn-primary {
  background: var(--primary-color);
  color: white;
}

.btn-primary:hover {
  background: var(--primary-hover);
}

.btn-secondary {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.btn-secondary:hover {
  background: var(--border-color);
}

.btn-danger {
  background: #fee2e2;
  color: #dc2626;
}

.btn-danger:hover {
  background: #fecaca;
}

.btn-icon {
  font-size: 12px;
}
</style>
