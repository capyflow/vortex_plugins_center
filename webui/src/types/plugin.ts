// 参数定义
export interface Parameter {
  name: string
  type: 'string' | 'number' | 'boolean' | 'object' | 'array'
  required: boolean
  default?: string
  description?: string
}

// 返回类型
export interface ReturnType {
  type: string
  description: string
}

// 插件方法
export interface Method {
  name: string
  description: string
  path: string
  method: string
  parameters?: Parameter[]
  returns?: ReturnType
}

// 插件健康状态
export interface PluginHealth {
  status: 'healthy' | 'unhealthy'
  latency: number
  checked_at: string
}

// 插件状态
export type PluginStatus = 'active' | 'inactive' | 'error'

// 插件定义
export interface Plugin {
  id: string
  name: string
  version: string
  description: string
  endpoint: string
  methods: Method[]
  metadata: Record<string, string>
  status: PluginStatus
  health: PluginHealth
  created_at: string
  updated_at: string
  last_seen: string
}

// 插件列表响应
export interface PluginListResponse {
  total: number
  plugins: Plugin[]
}

// 执行响应
export interface ExecuteResponse {
  success: boolean
  result: any
  error?: string
  latency: number
}

// 调用历史记录
export interface CallHistory {
  id: string
  pluginId: string
  pluginName: string
  methodName: string
  params: Record<string, any>
  response: ExecuteResponse
  timestamp: number
}
