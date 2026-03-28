<template>
  <div v-if="visible" class="dialog-overlay" @click.self="close">
    <div class="dialog-container">
      <div class="dialog-header">
        <h3 class="dialog-title">
          <span class="plugin-icon">🔌</span>
          执行 {{ plugin?.name }}
        </h3>
        <button class="close-btn" @click="close">×</button>
      </div>

      <div class="dialog-body">
        <p v-if="plugin?.summary" class="plugin-desc">{{ plugin.summary }}</p>

        <!-- 参数表单 -->
        <div v-if="hasParams" class="form-section">
          <h4 class="section-title">输入参数</h4>
          <PluginParamForm
            :params="plugin?.params || {}"
            v-model="formData"
          />
        </div>

        <div v-else class="no-params">
          <span class="info-icon">ℹ</span>
          此插件无需输入参数
        </div>

        <!-- 执行结果 -->
        <div v-if="result" class="result-section" :class="{ 'is-error': result.error }">
          <h4 class="section-title">
            {{ result.error ? '执行失败' : '执行结果' }}
            <span class="latency-badge" v-if="result.latency">
              {{ result.latency }}ms
            </span>
          </h4>
          <div class="result-content">
            <pre v-if="!result.error">{{ JSON.stringify(result.result, null, 2) }}</pre>
            <div v-else class="error-message">{{ result.error }}</div>
          </div>
        </div>
      </div>

      <div class="dialog-footer">
        <button class="btn btn-secondary" @click="close">取消</button>
        <button
          class="btn btn-primary"
          :disabled="loading"
          @click="execute"
        >
          <span v-if="loading" class="loading-spinner"></span>
          <span v-else class="btn-icon">▶</span>
          {{ loading ? '执行中...' : '执行' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import type { Plugin, ExecutePluginResponse } from '@/types/plugin'
import { PluginService } from '@/services/pluginService'
import PluginParamForm from './PluginParamForm.vue'

interface Props {
  visible: boolean
  plugin: Plugin | null
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:visible': [value: boolean]
  executed: [result: ExecutePluginResponse]
}>()

const formData = ref<Record<string, any>>({})
const loading = ref(false)
const result = ref<ExecutePluginResponse | null>(null)

const hasParams = computed(() => {
  return Object.keys(props.plugin?.params || {}).length > 0
})

// 关闭时重置状态
watch(() => props.visible, (visible) => {
  if (!visible) {
    formData.value = {}
    result.value = null
    loading.value = false
  }
})

function close() {
  emit('update:visible', false)
}

async function execute() {
  if (!props.plugin) return

  loading.value = true
  result.value = null

  try {
    const response = await PluginService.executePlugin({
      plugin_id: props.plugin.id,
      params: formData.value,
    })
    result.value = response
    emit('executed', response)
  } catch (error: any) {
    result.value = {
      success: false,
      error: error.message || '执行失败',
      latency: 0,
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.dialog-container {
  background: var(--bg-primary);
  border-radius: 16px;
  width: 100%;
  max-width: 600px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid var(--border-color);
}

.dialog-title {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 10px;
}

.plugin-icon {
  font-size: 24px;
}

.close-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  font-size: 24px;
  color: var(--text-secondary);
  cursor: pointer;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.close-btn:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.dialog-body {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

.plugin-desc {
  color: var(--text-secondary);
  margin: 0 0 20px 0;
  padding: 12px 16px;
  background: var(--bg-secondary);
  border-radius: 8px;
  font-size: 14px;
}

.form-section {
  margin-bottom: 24px;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 16px 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.no-params {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 20px;
  background: var(--bg-secondary);
  border-radius: 8px;
  color: var(--text-secondary);
  font-size: 14px;
}

.info-icon {
  width: 20px;
  height: 20px;
  background: var(--primary-color);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
}

.result-section {
  margin-top: 24px;
  padding: 16px;
  background: #f0fdf4;
  border: 1px solid #86efac;
  border-radius: 8px;
}

.result-section.is-error {
  background: #fef2f2;
  border-color: #fca5a5;
}

.latency-badge {
  font-size: 12px;
  padding: 2px 8px;
  background: var(--bg-primary);
  border-radius: 4px;
  color: var(--text-secondary);
  font-weight: normal;
}

.result-content {
  margin-top: 12px;
}

.result-content pre {
  margin: 0;
  padding: 16px;
  background: var(--bg-primary);
  border-radius: 6px;
  font-size: 13px;
  overflow-x: auto;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

.error-message {
  padding: 12px 16px;
  background: #fee2e2;
  color: #991b1b;
  border-radius: 6px;
  font-size: 14px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid var(--border-color);
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 6px;
}

.btn-primary {
  background: var(--primary-color);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: var(--primary-hover);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-secondary {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.btn-secondary:hover {
  background: var(--border-color);
}

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid transparent;
  border-top-color: currentColor;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.btn-icon {
  font-size: 12px;
}
</style>
