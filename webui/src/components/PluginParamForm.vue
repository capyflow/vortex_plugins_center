<template>
  <div class="plugin-param-form">
    <div
      v-for="(def, name) in params"
      :key="name"
      class="form-item"
    >
      <label class="form-label">
        {{ name }}
        <span v-if="def.required" class="required">*</span>
        <span v-if="def.description" class="description">({{ def.description }})</span>
      </label>

      <!-- 文本输入 -->
      <input
        v-if="def.type === 'string'"
        v-model="formData[name]"
        type="text"
        :placeholder="def.placeholder || `请输入 ${name}`"
        :maxlength="def.maxLength"
        :required="def.required"
        class="form-input"
      />

      <!-- 数字输入 -->
      <input
        v-else-if="def.type === 'number'"
        v-model.number="formData[name]"
        type="number"
        :min="def.min"
        :max="def.max"
        :step="def.step || 1"
        :placeholder="def.placeholder"
        :required="def.required"
        class="form-input"
      />

      <!-- 文件上传 -->
      <div v-else-if="def.type === 'file'" class="file-upload">
        <input
          type="file"
          :accept="def.accept"
          :multiple="def.multiple"
          @change="handleFileChange($event, name)"
          :required="def.required"
          class="file-input"
        />
        <div v-if="def.accept && def.accept !== '*'" class="file-hint">
          支持格式: {{ def.accept }}
        </div>
        <div v-if="def.maxSize" class="file-hint">
          最大大小: {{ formatFileSize(def.maxSize) }}
        </div>
      </div>

      <!-- 开关 -->
      <label v-else-if="def.type === 'boolean'" class="switch-label">
        <input
          v-model="formData[name]"
          type="checkbox"
          class="switch-input"
        />
        <span class="switch-slider"></span>
      </label>

      <!-- 下拉选择 -->
      <select
        v-else-if="def.type === 'select'"
        v-model="formData[name]"
        :multiple="def.multiple"
        :required="def.required"
        class="form-select"
      >
        <option v-for="opt in def.options" :key="opt" :value="opt">
          {{ opt }}
        </option>
      </select>

      <!-- 多行文本 -->
      <textarea
        v-else-if="def.type === 'textarea'"
        v-model="formData[name]"
        :rows="def.rows || 3"
        :maxlength="def.maxLength"
        :placeholder="def.placeholder || `请输入 ${name}`"
        :required="def.required"
        class="form-textarea"
      />

      <!-- 密码 -->
      <input
        v-else-if="def.type === 'password'"
        v-model="formData[name]"
        type="password"
        :placeholder="def.placeholder || `请输入 ${name}`"
        :required="def.required"
        class="form-input"
      />

      <!-- 日期 -->
      <input
        v-else-if="def.type === 'date'"
        v-model="formData[name]"
        type="date"
        :min="def.min"
        :max="def.max"
        :required="def.required"
        class="form-input"
      />

      <!-- 对象/JSON -->
      <div v-else-if="def.type === 'object' || def.type === 'array'" class="json-editor">
        <textarea
          v-model="jsonInputs[name]"
          @blur="validateJson(name)"
          rows="5"
          :placeholder="`请输入 ${def.type === 'object' ? 'JSON 对象' : 'JSON 数组'}`"
          :required="def.required"
          class="form-textarea json-textarea"
          :class="{ 'json-error': jsonErrors[name] }"
        />
        <div v-if="jsonErrors[name]" class="json-error-msg">
          {{ jsonErrors[name] }}
        </div>
      </div>

      <!-- 未知类型 -->
      <div v-else class="unsupported-type">
        不支持的参数类型: {{ def.type }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, reactive } from 'vue'
import type { ParamDef } from '@/types/plugin'

interface Props {
  params: Record<string, ParamDef>
  modelValue: Record<string, any>
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:modelValue': [value: Record<string, any>]
}>()

const formData = reactive<Record<string, any>>({})
const jsonInputs = reactive<Record<string, string>>({})
const jsonErrors = reactive<Record<string, string>>({})
const fileData = reactive<Record<string, File | File[]>>({})

// 初始化表单数据
watch(() => props.params, (params) => {
  Object.entries(params).forEach(([name, def]) => {
    if (!(name in formData)) {
      // 设置默认值
      if (def.default !== undefined) {
        formData[name] = def.default
      } else {
        switch (def.type) {
          case 'boolean':
            formData[name] = false
            break
          case 'number':
            formData[name] = def.min || 0
            break
          case 'select':
            formData[name] = def.options?.[0] || ''
            break
          case 'object':
            formData[name] = {}
            jsonInputs[name] = '{}'
            break
          case 'array':
            formData[name] = []
            jsonInputs[name] = '[]'
            break
          default:
            formData[name] = ''
        }
      }
    }
  })
}, { immediate: true })

// 监听表单变化，更新父组件
watch(formData, (val) => {
  // 合并文件数据
  const result = { ...val, ...fileData }
  emit('update:modelValue', result)
}, { deep: true })

// 处理文件选择
function handleFileChange(event: Event, name: string) {
  const input = event.target as HTMLInputElement
  const def = props.params[name]

  if (def.multiple) {
    fileData[name] = Array.from(input.files || [])
  } else {
    fileData[name] = input.files?.[0]
  }

  // 触发更新
  emit('update:modelValue', { ...formData, ...fileData })
}

// 验证 JSON
function validateJson(name: string) {
  const input = jsonInputs[name]
  try {
    const parsed = JSON.parse(input)
    formData[name] = parsed
    jsonErrors[name] = ''
  } catch (e: any) {
    jsonErrors[name] = e.message
  }
}

// 格式化文件大小
function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}
</script>

<style scoped>
.plugin-param-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-label {
  font-weight: 500;
  color: var(--text-primary);
}

.form-label .required {
  color: #ef4444;
  margin-left: 4px;
}

.form-label .description {
  color: var(--text-secondary);
  font-size: 12px;
  font-weight: normal;
  margin-left: 8px;
}

.form-input,
.form-select,
.form-textarea {
  padding: 8px 12px;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  font-size: 14px;
  background: var(--bg-primary);
  color: var(--text-primary);
  transition: border-color 0.2s;
}

.form-input:focus,
.form-select:focus,
.form-textarea:focus {
  outline: none;
  border-color: var(--primary-color);
}

.form-textarea {
  resize: vertical;
  font-family: monospace;
}

.file-upload {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.file-input {
  padding: 8px 0;
}

.file-hint {
  font-size: 12px;
  color: var(--text-secondary);
}

/* 开关样式 */
.switch-label {
  position: relative;
  display: inline-block;
  width: 44px;
  height: 24px;
  cursor: pointer;
}

.switch-input {
  opacity: 0;
  width: 0;
  height: 0;
}

.switch-slider {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #ccc;
  border-radius: 24px;
  transition: 0.3s;
}

.switch-slider:before {
  content: '';
  position: absolute;
  height: 18px;
  width: 18px;
  left: 3px;
  bottom: 3px;
  background-color: white;
  border-radius: 50%;
  transition: 0.3s;
}

.switch-input:checked + .switch-slider {
  background-color: var(--primary-color);
}

.switch-input:checked + .switch-slider:before {
  transform: translateX(20px);
}

/* JSON 编辑器 */
.json-editor {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.json-textarea {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
}

.json-textarea.json-error {
  border-color: #ef4444;
}

.json-error-msg {
  color: #ef4444;
  font-size: 12px;
}

/* 不支持类型 */
.unsupported-type {
  padding: 12px;
  background: #fef3c7;
  border-radius: 6px;
  color: #92400e;
  font-size: 14px;
}
</style>