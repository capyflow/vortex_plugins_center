<template>
  <div class="plugin-list-page">
    <!-- 头部统计和搜索 -->
    <div class="page-header">
      <div class="header-left">
        <h1>插件列表</h1>
        <div class="stats">
          <el-statistic title="总插件" :value="pluginStore.plugins.length" />
          <el-statistic title="运行中" :value="pluginStore.activePlugins.length" class="stat-active" />
          <el-statistic title="异常" :value="pluginStore.errorPlugins.length" class="stat-error" />
        </div>
      </div>
      <div class="header-right">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索插件..."
          clearable
          class="search-input"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-button type="primary" @click="refresh">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <!-- 插件卡片网格 -->
    <div class="plugins-grid">
      <PluginCard
        v-for="plugin in filteredPlugins"
        :key="plugin.id"
        :plugin="plugin"
        @click="goToDetail(plugin)"
      />
    </div>

    <!-- 空状态 -->
    <el-empty v-if="filteredPlugins.length === 0" description="暂无插件" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { Search, Refresh } from '@element-plus/icons-vue'
import { usePluginStore } from '@/stores/plugin'
import { useSettingsStore } from '@/stores/settings'
import type { Plugin } from '@/types/plugin'
import PluginCard from '@/components/PluginCard.vue'

const pluginStore = usePluginStore()
const settingsStore = useSettingsStore()
const router = useRouter()

const searchKeyword = ref('')

// 过滤插件
const filteredPlugins = computed(() => {
  if (!searchKeyword.value) return pluginStore.plugins
  const keyword = searchKeyword.value.toLowerCase()
  return pluginStore.plugins.filter(p =>
    p.name.toLowerCase().includes(keyword) ||
    p.description?.toLowerCase().includes(keyword)
  )
})

// 跳转到插件详情页
function goToDetail(plugin: Plugin) {
  router.push(`/plugin/${plugin.id}`)
}

// 刷新
function refresh() {
  pluginStore.fetchPlugins()
}

// 自动刷新
let refreshTimer: ReturnType<typeof setInterval> | null = null

onMounted(() => {
  pluginStore.fetchPlugins()
  if (settingsStore.autoRefresh) {
    refreshTimer = setInterval(() => {
      pluginStore.fetchPlugins()
    }, settingsStore.refreshInterval)
  }
})

onUnmounted(() => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
  }
})
</script>

<style scoped>
.plugin-list-page {
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.header-left h1 {
  margin: 0 0 16px 0;
  font-size: 28px;
}

.stats {
  display: flex;
  gap: 32px;
}

:deep(.stat-active .el-statistic__content) {
  color: #2ed573;
}

:deep(.stat-error .el-statistic__content) {
  color: #ff4757;
}

.header-right {
  display: flex;
  gap: 12px;
}

.search-input {
  width: 280px;
}

.plugins-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}
</style>
