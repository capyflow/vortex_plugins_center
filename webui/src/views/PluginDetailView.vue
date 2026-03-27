<template>
  <div class="plugin-detail-page">
    <!-- 返回按钮 -->
    <div class="page-header">
      <el-button @click="goBack" class="back-btn">
        <el-icon><ArrowLeft /></el-icon>
        返回列表
      </el-button>
    </div>

    <PluginDetail v-if="plugin" :plugin="plugin" />
    <el-empty v-else description="插件不存在或已删除" />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'
import { usePluginStore } from '@/stores/plugin'
import PluginDetail from '@/components/PluginDetail.vue'

const route = useRoute()
const router = useRouter()
const pluginStore = usePluginStore()

const pluginId = computed(() => route.params.id as string)
const plugin = computed(() => pluginStore.plugins.find(p => p.id === pluginId.value))

onMounted(() => {
  if (!plugin.value) {
    pluginStore.fetchPlugins()
  }
})

function goBack() {
  router.push('/')
}
</script>

<style scoped>
.plugin-detail-page {
  padding: 24px;
}

.page-header {
  margin-bottom: 24px;
}

.back-btn {
  background: rgba(255, 255, 255, 0.1) !important;
  border-color: rgba(255, 255, 255, 0.2) !important;
  color: #fff !important;
}

.back-btn:hover {
  background: rgba(255, 255, 255, 0.2) !important;
}
</style>
