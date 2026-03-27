<template>
  <div class="method-panel">
    <div class="panel-header">
      <div>
        <h3 class="method-name">{{ method.name }}</h3>
        <p class="method-desc">{{ method.description }}</p>
      </div>
      <el-tag type="info">{{ method.method }} {{ method.path }}</el-tag>
    </div>

    <div v-if="hasParams" class="params-section">
      <h4>参数</h4>
      <el-form label-position="top">
        <el-row :gutter="20">
          <el-col v-for="param in method.parameters" :key="param.name" :span="12">
            <el-form-item>
              <template #label>
                <span>{{ param.name }}</span>
                <el-tag v-if="param.required" type="danger" size="small" class="required-tag">必填</el-tag>
                <span class="param-type">{{ param.type }}</span>
              </template>
              <el-input
                v-model="formData[param.name]"
                :placeholder="param.description || `${param.type}类型`"
                clearable
              />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </div>

    <div class="action-section">
      <el-button type="primary" @click="execute" :loading="executing">
        <el-icon><VideoPlay /></el-icon>
        执行
      </el-button>
    </div>

    <div v-if="result" class="result-section" :class="{ success: result.success, error: !result.success }">
      <div class="result-header">
        <span>{{ result.success ? '执行成功' : '执行失败' }}</span>
        <span v-if="result.latency" class="latency">{{ result.latency }}ms</span>
      </div>
      <pre class="result-body">{{ JSON.stringify(result.success ? result.result : result.error, null, 2) }}</pre>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { VideoPlay } from '@element-plus/icons-vue'
import type { Method, ExecuteResponse } from '@/types/plugin'
import { usePluginStore } from '@/stores/plugin'

const props = defineProps<{
  pluginId: string
  method: Method
}>()

const pluginStore = usePluginStore()

const formData = ref<Record<string, any>>({})
const executing = ref(false)
const result = ref<ExecuteResponse | null>(null)

const hasParams = computed(() => props.method.parameters && props.method.parameters.length > 0)

async function execute() {
  executing.value = true
  result.value = null
  
  try {
    const params: Record<string, any> = {}
    props.method.parameters?.forEach(p => {
      const value = formData.value[p.name]
      if (value !== undefined && value !== '') {
        params[p.name] = p.type === 'number' ? parseFloat(value) : value
      }
    })
    
    result.value = await pluginStore.executeMethod(props.pluginId, props.method.name, params)
  } catch (err) {
    result.value = err as ExecuteResponse
  } finally {
    executing.value = false
  }
}
</script>

<style scoped>
.method-panel {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 12px;
  padding: 24px;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: start;
  margin-bottom: 20px;
}

.method-name {
  font-size: 18px;
  font-weight: 600;
  color: #667eea;
  margin: 0 0 8px 0;
}

.method-desc {
  color: #aaa;
  font-size: 14px;
  margin: 0;
}

.params-section {
  margin-bottom: 20px;
}

.params-section h4 {
  color: #fff;
  margin-bottom: 16px;
}

.required-tag {
  margin-left: 8px;
}

.param-type {
  color: #888;
  font-size: 12px;
  margin-left: 8px;
}

.action-section {
  margin-bottom: 20px;
}

.result-section {
  border-radius: 8px;
  overflow: hidden;
}

.result-section.success {
  border: 1px solid rgba(46, 213, 115, 0.3);
}

.result-section.error {
  border: 1px solid rgba(255, 71, 87, 0.3);
}

.result-header {
  padding: 12px 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.result-section.success .result-header {
  background: rgba(46, 213, 115, 0.1);
  color: #2ed573;
}

.result-section.error .result-header {
  background: rgba(255, 71, 87, 0.1);
  color: #ff4757;
}

.latency {
  font-size: 12px;
  opacity: 0.8;
}

.result-body {
  padding: 16px;
  background: rgba(0, 0, 0, 0.3);
  color: #fff;
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 13px;
  line-height: 1.6;
  overflow-x: auto;
  margin: 0;
}
</style>
