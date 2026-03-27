import axios, { type AxiosInstance, type AxiosResponse } from 'axios'
import type { Plugin, PluginListResponse, ExecuteResponse } from '@/types/plugin'

// API 基础配置
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'

// 创建 axios 实例
const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 响应拦截器 - 统一处理响应格式
apiClient.interceptors.response.use(
  (response: AxiosResponse) => {
    // 后端返回格式: { body: { ... } }
    if (response.data && response.data.body !== undefined) {
      return response.data.body
    }
    return response.data
  },
  (error) => {
    console.error('API Error:', error)
    return Promise.reject(error)
  }
)

// 插件 API
export const pluginApi = {
  // 获取插件列表
  getPluginList(keyword?: string, status?: string, page: number = 1, limit: number = 100): Promise<PluginListResponse> {
    const params = new URLSearchParams()
    if (keyword) params.append('keyword', keyword)
    if (status) params.append('status', status)
    params.append('page', page.toString())
    params.append('limit', limit.toString())

    return apiClient.get(`/plugins/list?${params.toString()}`)
  },

  // 获取单个插件
  getPlugin(id: string): Promise<Plugin> {
    return apiClient.get(`/plugins/${id}`)
  },

  // 执行插件方法
  executeMethod(pluginId: string, methodName: string, params?: Record<string, any>): Promise<ExecuteResponse> {
    return apiClient.post(`/plugins/${pluginId}/execute/${methodName}`, params || {})
  },

  // 检查插件健康
  checkHealth(pluginId: string): Promise<Plugin['health']> {
    return apiClient.get(`/plugins/${pluginId}/health`)
  },

  // 注销插件
  unregisterPlugin(pluginId: string): Promise<void> {
    return apiClient.delete(`/plugins/${pluginId}`)
  },
}

export default apiClient
