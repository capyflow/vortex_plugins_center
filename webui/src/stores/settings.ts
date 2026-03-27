import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useSettingsStore = defineStore('settings', () => {
  // State
  const apiBaseUrl = ref(import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1')
  const autoRefresh = ref(true)
  const refreshInterval = ref(30000) // 30秒
  const darkMode = ref(true)

  // Actions
  function setApiBaseUrl(url: string) {
    apiBaseUrl.value = url
  }

  function setAutoRefresh(enabled: boolean) {
    autoRefresh.value = enabled
  }

  function setRefreshInterval(interval: number) {
    refreshInterval.value = interval
  }

  function toggleDarkMode() {
    darkMode.value = !darkMode.value
    // 应用主题
    if (darkMode.value) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }

  function initTheme() {
    if (darkMode.value) {
      document.documentElement.classList.add('dark')
    }
  }

  return {
    apiBaseUrl,
    autoRefresh,
    refreshInterval,
    darkMode,
    setApiBaseUrl,
    setAutoRefresh,
    setRefreshInterval,
    toggleDarkMode,
    initTheme,
  }
})
