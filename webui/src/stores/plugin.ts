import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Plugin, ExecuteResponse, CallHistory } from '@/types/plugin'
import { pluginApi } from '@/api/plugin'

export const usePluginStore = defineStore('plugin', () => {
  // State
  const plugins = ref<Plugin[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const selectedPlugin = ref<Plugin | null>(null)
  const callHistory = ref<CallHistory[]>([])

  // Getters
  const activePlugins = computed(() =>
    plugins.value.filter(p => p.status === 'active')
  )

  const inactivePlugins = computed(() =>
    plugins.value.filter(p => p.status === 'inactive')
  )

  const errorPlugins = computed(() =>
    plugins.value.filter(p => p.status === 'error')
  )

  // Actions
  // 加载插件列表
  async function fetchPlugins(keyword?: string, status?: string) {
    loading.value = true
    error.value = null
    try {
      const response = await pluginApi.getPluginList(keyword, status)
      plugins.value = response.plugins || []
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch plugins'
      console.error('Failed to fetch plugins:', err)
    } finally {
      loading.value = false
    }
  }

  // 选择插件
  function selectPlugin(plugin: Plugin | null) {
    selectedPlugin.value = plugin
  }

  // 执行方法
  async function executeMethod(
    pluginId: string,
    methodName: string,
    params: Record<string, any>
  ): Promise<ExecuteResponse> {
    try {
      const response = await pluginApi.executeMethod(pluginId, methodName, params)

      // 记录调用历史
      const plugin = plugins.value.find(p => p.id === pluginId)
      if (plugin) {
        const history: CallHistory = {
          id: `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
          pluginId,
          pluginName: plugin.name,
          methodName,
          params,
          response,
          timestamp: Date.now(),
        }
        callHistory.value.unshift(history)
        // 只保留最近 50 条历史
        if (callHistory.value.length > 50) {
          callHistory.value = callHistory.value.slice(0, 50)
        }
      }

      return response
    } catch (err) {
      const errorResponse: ExecuteResponse = {
        success: false,
        result: null,
        error: err instanceof Error ? err.message : 'Execution failed',
        latency: 0,
      }
      throw errorResponse
    }
  }

  // 检查健康状态
  async function checkPluginHealth(pluginId: string) {
    try {
      const health = await pluginApi.checkHealth(pluginId)
      const plugin = plugins.value.find(p => p.id === pluginId)
      if (plugin) {
        plugin.health = health
      }
      return health
    } catch (err) {
      console.error('Health check failed:', err)
      return null
    }
  }

  // 刷新单个插件
  async function refreshPlugin(pluginId: string) {
    try {
      const plugin = await pluginApi.getPlugin(pluginId)
      const index = plugins.value.findIndex(p => p.id === pluginId)
      if (index !== -1) {
        plugins.value[index] = plugin
      }
      return plugin
    } catch (err) {
      console.error('Failed to refresh plugin:', err)
      return null
    }
  }

  // 清除历史记录
  function clearHistory() {
    callHistory.value = []
  }

  return {
    plugins,
    loading,
    error,
    selectedPlugin,
    callHistory,
    activePlugins,
    inactivePlugins,
    errorPlugins,
    fetchPlugins,
    selectPlugin,
    executeMethod,
    checkPluginHealth,
    refreshPlugin,
    clearHistory,
  }
})
