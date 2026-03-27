import { createRouter, createWebHistory } from 'vue-router'
import AppLayout from '@/components/AppLayout.vue'
import PluginList from '@/views/PluginList.vue'

const routes = [
  {
    path: '/',
    component: AppLayout,
    children: [
      {
        path: '',
        name: 'PluginList',
        component: PluginList
      },
      {
        path: '/history',
        name: 'History',
        component: () => import('@/views/History.vue')
      },
      {
        path: '/plugin/:id',
        name: 'PluginDetail',
        component: () => import('@/views/PluginDetailView.vue')
      },
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
