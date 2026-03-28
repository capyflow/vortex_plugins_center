import type { Plugin, RegisterPluginRequest, ExecutePluginRequest, ExecutePluginResponse } from '@/types/plugin'

const API_BASE = '/api/v1'

export class PluginService {
  // 获取插件列表
  static async listPlugins(keyword?: string, status?: string): Promise<{ total: number; plugins: Plugin[] }> {
    const params = new URLSearchParams()
    if (keyword) params.append('keyword', keyword)
    if (status) params.append('status', status)

    const response = await fetch(`${API_BASE}/plugins/list?${params}`)
    const data = await response.json()
    return data.body
  }

  // 获取插件详情
  static async getPlugin(id: string): Promise<Plugin> {
    const response = await fetch(`${API_BASE}/plugins/${id}`)
    const data = await response.json()
    return data.body
  }

  // 注册插件
  static async registerPlugin(plugin: RegisterPluginRequest): Promise<{ id: string; status: string }> {
    const response = await fetch(`${API_BASE}/plugins/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(plugin),
    })
    const data = await response.json()
    return data.body
  }

  // 注销插件
  static async unregisterPlugin(id: string): Promise<void> {
    await fetch(`${API_BASE}/plugins/${id}`, { method: 'DELETE' })
  }

  // 执行插件
  static async executePlugin(request: ExecutePluginRequest): Promise<ExecutePluginResponse> {
    const response = await fetch(`${API_BASE}/plugins/${request.plugin_id}/execute`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(request.params || {}),
    })
    const data = await response.json()
    return data.body
  }

  // 健康检查
  static async healthCheck(id: string): Promise<{ status: string; latency: number }> {
    const response = await fetch(`${API_BASE}/plugins/${id}/health`)
    const data = await response.json()
    return data.body
  }
}
