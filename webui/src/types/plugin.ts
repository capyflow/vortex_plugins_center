/**
 * 插件平台类型定义
 */

// 参数定义
export interface ParamDef {
  type: string
  required?: boolean
  default?: any
  description?: string
  // 类型特定属性
  maxLength?: number
  min?: number
  max?: number
  step?: number
  accept?: string
  multiple?: boolean
  maxSize?: number
  options?: string[]
  rows?: number
  placeholder?: string
}

// 插件信息
export interface Plugin {
  id: string
  name: string
  version: string
  endpoint: string
  summary?: string
  params?: Record<string, ParamDef>
  outputs?: string
  metadata?: Record<string, string>
  status: 'active' | 'inactive' | 'error'
  health?: {
    status: string
    latency: number
    checked_at: string
  }
  created_at: string
  updated_at: string
}

// 注册请求
export interface RegisterPluginRequest {
  name: string
  version: string
  endpoint: string
  summary?: string
  params?: Record<string, ParamDef>
  outputs?: string
  metadata?: Record<string, string>
}

// 执行请求
export interface ExecutePluginRequest {
  plugin_id: string
  params?: Record<string, any>
}

// 执行响应
export interface ExecutePluginResponse {
  success: boolean
  result?: any
  error?: string
  latency: number
}

// 支持的参数类型
export const PARAM_TYPES = [
  { value: 'string', label: '文本', icon: 'text' },
  { value: 'number', label: '数字', icon: 'hash' },
  { value: 'file', label: '文件', icon: 'upload' },
  { value: 'boolean', label: '开关', icon: 'toggle' },
  { value: 'select', label: '选择', icon: 'list' },
  { value: 'textarea', label: '多行文本', icon: 'align-left' },
  { value: 'password', label: '密码', icon: 'lock' },
  { value: 'date', label: '日期', icon: 'calendar' },
  { value: 'object', label: '对象', icon: 'code' },
  { value: 'array', label: '数组', icon: 'list-ordered' },
] as const

// 参数类型对应的默认配置
export const PARAM_TYPE_DEFAULTS: Record<string, Partial<ParamDef>> = {
  string: { required: true },
  number: { required: true, min: 0 },
  file: { required: true, accept: '*' },
  boolean: { default: false },
  select: { options: [] },
  textarea: { required: true, rows: 3 },
  password: { required: true },
  date: { required: true },
  object: { required: true },
  array: { required: true },
}
