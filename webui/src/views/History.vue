<template>
  <div class="history-page">
    <h1>调用历史</h1>
    <el-table :data="pluginStore.callHistory" class="history-table">
      <el-table-column prop="timestamp" label="时间" width="180">
        <template #default="{ row }">
          {{ formatTime(row.timestamp) }}
        </template>
      </el-table-column>
      <el-table-column prop="pluginName" label="插件" width="150" />
      <el-table-column prop="methodName" label="方法" width="150" />
      <el-table-column prop="params" label="参数">
        <template #default="{ row }">
          <code>{{ JSON.stringify(row.params) }}</code>
        </template>
      </el-table-column>
      <el-table-column prop="response.success" label="结果" width="100">
        <template #default="{ row }">
          <el-tag :type="row.response.success ? 'success' : 'danger'">
            {{ row.response.success ? '成功' : '失败' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="response.latency" label="耗时" width="100">
        <template #default="{ row }">
          {{ row.response.latency }}ms
        </template>
      </el-table-column>
    </el-table>
    <el-empty v-if="pluginStore.callHistory.length === 0" description="暂无调用记录" />
  </div>
</template>

<script setup lang="ts">
import { usePluginStore } from '@/stores/plugin'

const pluginStore = usePluginStore()

function formatTime(timestamp: number) {
  return new Date(timestamp).toLocaleString('zh-CN')
}
</script>

<style scoped>
.history-page {
  padding: 24px;
}

.history-page h1 {
  margin-bottom: 20px;
}

.history-table {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 8px;
}

:deep(.el-table) {
  background: transparent;
}

:deep(.el-table__header-wrapper th) {
  background: rgba(0, 0, 0, 0.3) !important;
  color: #fff;
}

:deep(.el-table__row) {
  background: transparent !important;
  color: #fff;
}

:deep(.el-table__row:hover > td) {
  background: rgba(255, 255, 255, 0.05) !important;
}

code {
  background: rgba(0, 0, 0, 0.3);
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 12px;
}
</style>
