<template>
  <div class="plugins-view">
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">
          <span class="title-icon">🔌</span>
          插件中心
        </h1>
        <p class="page-desc">管理和执行插件，扩展平台功能</p>
      </div>
      <div class="header-actions">
        <div class="search-box">
          <input
            v-model="searchKeyword"
            type="text"
            placeholder="搜索插件..."
            class="search-input"
            @input="handleSearch"
          />
          <span class="search-icon">🔍</span>
        </div>
        <button class="btn btn-primary" @click="showRegisterDialog = true">
          <span class="btn-icon">+</span>
          注册插件
        </button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-value">{{ stats.total }}</div>
        <div class="stat-label">总插件数</div>
      </div>
      <div class="stat-card active">
        <div class="stat-value">{{ stats.active }}</div>
        <div class="stat-label">运行中</div>
      </div>
      <div class="stat-card error" v-if="stats.error > 0">
        <div class="stat-value">{{ stats.error }}</div>
        <div class="stat-label">异常</div>
      </div>
    </div>

    <!-- 插件列表 -->
    <div v-if="loading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>加载中...</p>
    </div>

    <div v-else-if="plugins.length === 0" class="empty-state">
      <div class="empty-icon">🔌</div>
      <h3>暂无插件</h3>
      <p>点击右上角注册插件按钮添加第一个插件</p>
    </div>

    <div v-else class="plugins-grid">
      <PluginCard
        v-for="plugin in filteredPlugins"
        :key="plugin.id"
        :plugin="plugin"
        @execute="handleExecute"
        @detail="handleDetail"
        @unregister="handleUnregister"
      />
    </div>

    <!-- 执行对话框 -->
    <PluginExecuteDialog
      v-model:visible="executeDialogVisible"
      :plugin="selectedPlugin"
      @executed="handleExecuted"
    />

    <!-- 注册对话框 -->
    <div v-if="showRegisterDialog" class="dialog-overlay" @click.self="showRegisterDialog = false">
      <div class="dialog-container register-dialog">
        <div class="dialog-header">
          <h3 class="dialog-title">注册新插件</h3>
          <button class="close-btn" @click="showRegisterDialog = false">×</button>
        </div>
        <div class="dialog-body">
          <div class="form-group">
            <label>插件名称 <span class="required">*</span></label>
            <input v-model="registerForm.name" type="text" placeholder="例如: calculator" />
          </div>
          <div class="form-group">
            <label>版本 <span class="required">*</span></label>
            <input v-model="registerForm.version" type="text" placeholder="例如: 1.0.0" />
          </div>
          <div class="form-group">
            <label>端点地址 <span class="required">*</span></label>
            <input v-model="registerForm.endpoint" type="text" placeholder="例如: http://localhost:8001" />
          </div>
          <div class="form-group">
            <label>描述</label>
            <input v-model="registerForm.summary" type="text" placeholder="插件功能简介" />
          </div>
          <div class="form-group">
            <label>输出类型</label>
            <select v-model="registerForm.outputs">
              <option value="">请选择</option>
              <option value="string">字符串</option>
              <option value="number">数字</option>
              <option value="boolean">布尔值</option>
              <option value="object">对象</option>
              <option value="array">数组</option>
            </select>
          </div>
        </div>
        <div class="dialog-footer">
          <button class="btn btn-secondary" @click="showRegisterDialog = false">取消</button>
          <button class="btn btn-primary" :disabled="registering" @click="handleRegister">
            {{ registering ? '注册中...' : '注册' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { Plugin } from '@/types/plugin'
import { PluginService } from '@/services/pluginService'
import PluginCard from '@/components/PluginCard.vue'
import PluginExecuteDialog from '@/components/PluginExecuteDialog.vue'

const plugins = ref<Plugin[]>([])
const loading = ref(false)
const searchKeyword = ref('')
const executeDialogVisible = ref(false)
const selectedPlugin = ref<Plugin | null>(null)
const showRegisterDialog = ref(false)
const registering = ref(false)

const registerForm = ref({
  name: '',
  version: '',
  endpoint: '',
  summary: '',
  outputs: '',
})

const stats = computed(() => {
  return {
    total: plugins.value.length,
    active: plugins.value.filter(p => p.status === 'active').length,
    error: plugins.value.filter(p => p.status === 'error').length,
  }
})

const filteredPlugins = computed(() => {
  if (!searchKeyword.value) return plugins.value
  const keyword = searchKeyword.value.toLowerCase()
  return plugins.value.filter(p =>
    p.name.toLowerCase().includes(keyword) ||
    p.summary?.toLowerCase().includes(keyword)
  )
})

onMounted(() => {
  loadPlugins()
})

async function loadPlugins() {
  loading.value = true
  try {
    const result = await PluginService.listPlugins()
    plugins.value = result.plugins
  } catch (error) {
    console.error('Failed to load plugins:', error)
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  // 防抖搜索可以在这里实现
}

function handleExecute(plugin: Plugin) {
  selectedPlugin.value = plugin
  executeDialogVisible.value = true
}

function handleDetail(plugin: Plugin) {
  // 可以跳转到详情页或打开详情对话框
  console.log('Plugin detail:', plugin)
}

async function handleUnregister(plugin: Plugin) {
  if (!confirm(`确定要注销插件 "${plugin.name}" 吗？`)) return

  try {
    await PluginService.unregisterPlugin(plugin.id)
    await loadPlugins()
  } catch (error) {
    alert('注销失败: ' + error)
  }
}

function handleExecuted(result: any) {
  console.log('Execution result:', result)
}

async function handleRegister() {
  if (!registerForm.value.name || !registerForm.value.version || !registerForm.value.endpoint) {
    alert('请填写必填项')
    return
  }

  registering.value = true
  try {
    await PluginService.registerPlugin(registerForm.value)
    showRegisterDialog.value = false
    await loadPlugins()
    // 重置表单
    registerForm.value = {
      name: '',
      version: '',
      endpoint: '',
      summary: '',
      outputs: '',
    }
  } catch (error: any) {
    alert('注册失败: ' + error.message)
  } finally {
    registering.value = false
  }
}
</script>

<style scoped>
.plugins-view {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
}

.header-content {
  flex: 1;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  margin: 0 0 8px 0;
  display: flex;
  align-items: center;
  gap: 12px;
}

.title-icon {
  font-size: 32px;
}

.page-desc {
  color: var(--text-secondary);
  margin: 0;
  font-size: 14px;
}

.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.search-box {
  position: relative;
}

.search-input {
  width: 280px;
  padding: 10px 16px 10px 40px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 14px;
  background: var(--bg-primary);
  color: var(--text-primary);
}

.search-input:focus {
  outline: none;
  border-color: var(--primary-color);
}

.search-icon {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-secondary);
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 6px;
}

.btn-primary {
  background: var(--primary-color);
  color: white;
}

.btn-primary:hover {
  background: var(--primary-hover);
}

.btn-secondary {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.btn-secondary:hover {
  background: var(--border-color);
}

.btn-icon {
  font-size: 16px;
}

/* 统计卡片 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  background: var(--bg-secondary);
  border-radius: 12px;
  padding: 20px;
  text-align: center;
}

.stat-card.active {
  background: linear-gradient(135deg, #dcfce7, #bbf7d0);
}

.stat-card.error {
  background: linear-gradient(135deg, #fee2e2, #fecaca);
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.stat-card.active .stat-value {
  color: #166534;
}

.stat-card.error .stat-value {
  color: #991b1b;
}

.stat-label {
  font-size: 14px;
  color: var(--text-secondary);
}

/* 加载和空状态 */
.loading-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  color: var(--text-secondary);
}

.loading-spinner {
  width: 48px;
  height: 48px;
  border: 4px solid var(--border-color);
  border-top-color: var(--primary-color);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-state h3 {
  margin: 0 0 8px 0;
  color: var(--text-primary);
}

.empty-state p {
  margin: 0;
}

/* 插件网格 */
.plugins-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(360px, 1fr));
  gap: 20px;
}

/* 对话框样式 */
.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.dialog-container {
  background: var(--bg-primary);
  border-radius: 16px;
  width: 100%;
  max-width: 500px;
  max-height: 90vh;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid var(--border-color);
}

.dialog-title {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
}

.close-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  font-size: 24px;
  color: var(--text-secondary);
  cursor: pointer;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.close-btn:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.dialog-body {
  padding: 24px;
  max-height: 60vh;
  overflow-y: auto;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 6px;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

.form-group .required {
  color: #ef4444;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 14px;
  background: var(--bg-primary);
  color: var(--text-primary);
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: var(--primary-color);
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid var(--border-color);
}
</style>