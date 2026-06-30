<template>
  <div class="admin-layout">
    <n-layout has-sider position="absolute">
      <n-layout-sider
        bordered
        collapse-mode="width"
        :collapsed-width="64"
        :width="240"
        show-trigger
        :native-scrollbar="false"
        class="sidebar"
      >
        <div class="sidebar-logo">
          <span class="logo-text">HiGolang</span>
          <span class="logo-sub">管理后台</span>
        </div>
        
        <n-menu
          :value="activeMenu"
          :options="menuOptions"
          @update:value="handleMenuClick"
          class="sidebar-menu"
        />
      </n-layout-sider>
      
      <n-layout>
        <n-layout-header bordered class="admin-header">
          <div class="header-content">
            <n-breadcrumb>
              <n-breadcrumb-item>
                <span @click="$router.push({ name: 'Dashboard' })" style="cursor: pointer;">管理后台</span>
              </n-breadcrumb-item>
              <n-breadcrumb-item v-if="currentRouteName">
                {{ currentRouteName }}
              </n-breadcrumb-item>
            </n-breadcrumb>
            
            <div class="header-actions">
              <n-dropdown :options="userMenuOptions" @select="handleUserMenu">
                <n-button quaternary>
                  <template #icon>
                    <n-icon><person-icon /></n-icon>
                  </template>
                  管理员
                </n-button>
              </n-dropdown>
            </div>
          </div>
        </n-layout-header>
        
        <n-layout-content content-style="padding: 24px;" class="admin-content">
          <div class="content-wrapper">
            <router-view />
          </div>
        </n-layout-content>
      </n-layout>
    </n-layout>
  </div>
</template>

<script setup>
import { computed, h } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { NIcon } from 'naive-ui'
import { 
  GridOutline, 
  DocumentTextOutline, 
  FolderOutline, 
  PricetagsOutline,
  NewspaperOutline,
  ListOutline,
  SettingsOutline,
  PersonOutline as PersonIcon,
  LogOutOutline
} from '@vicons/ionicons5'
import { useAuthStore } from '../../stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const activeMenu = computed(() => {
  const name = route.name
  if (name === 'Dashboard') return 'dashboard'
  if (name?.includes('Article')) return 'articles'
  if (name === 'FeedSources') return 'feed-sources'
  if (name === 'FetchLogs') return 'fetch-logs'
  if (name === 'Settings') return 'settings'
  return null
})

const currentRouteName = computed(() => {
  const name = route.name
  const map = {
    'Dashboard': '仪表盘',
    'ArticleManage': '文章管理',
    'ArticleCreate': '创建文章',
    'ArticleUpdate': '编辑文章',
    'FeedSources': '资讯源管理',
    'FetchLogs': '抓取日志',
    'Settings': '系统设置'
  }
  return map[name] || ''
})

const renderIcon = (icon) => {
  return () => h(NIcon, null, { default: () => h(icon) })
}

const menuOptions = [
  { label: '仪表盘', key: 'dashboard', icon: renderIcon(GridOutline) },
  { label: '文章管理', key: 'articles', icon: renderIcon(DocumentTextOutline) },
  { label: '分类管理', key: 'categories', icon: renderIcon(FolderOutline) },
  { label: '标签管理', key: 'tags', icon: renderIcon(PricetagsOutline) },
  { label: '资讯源管理', key: 'feed-sources', icon: renderIcon(NewspaperOutline) },
  { label: '抓取日志', key: 'fetch-logs', icon: renderIcon(ListOutline) },
  { label: '系统设置', key: 'settings', icon: renderIcon(SettingsOutline) }
]

const userMenuOptions = [
  { label: '退出登录', key: 'logout', icon: renderIcon(LogOutOutline) }
]

const handleMenuClick = (key) => {
  const routes = {
    'dashboard': { name: 'Dashboard' },
    'articles': { name: 'ArticleManage' },
    'categories': { name: 'ArticleManage', query: { tab: 'categories' } },
    'tags': { name: 'ArticleManage', query: { tab: 'tags' } },
    'feed-sources': { name: 'FeedSources' },
    'fetch-logs': { name: 'FetchLogs' },
    'settings': { name: 'Settings' }
  }
  
  if (routes[key]) {
    router.push(routes[key])
  }
}

const handleUserMenu = (key) => {
  if (key === 'logout') {
    authStore.logout()
  }
}
</script>

<style scoped>
.admin-layout {
  min-height: 100vh;
  background-color: #f5f7fa;
}

.sidebar {
  background-color: #fff;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.06);
}

.sidebar-logo {
  padding: 20px 16px;
  border-bottom: 1px solid #e8e8e8;
  text-align: center;
}

.logo-text {
  display: block;
  font-size: 20px;
  font-weight: bold;
  color: #00ADD8;
  letter-spacing: -0.5px;
}

.logo-sub {
  display: block;
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.sidebar-menu {
  padding: 16px 0;
}

.admin-header {
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.header-content {
  height: 64px;
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.admin-content {
  min-height: calc(100vh - 64px);
}

.content-wrapper {
  max-width: 1400px;
  margin: 0 auto;
}
</style>
