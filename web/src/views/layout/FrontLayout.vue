<template>
  <div class="front-layout">
    <n-layout position="absolute">
      <n-layout-header bordered class="header">
        <div class="header-content">
          <div class="logo" @click="$router.push({ name: 'Home' })">
            <span class="logo-text">HiGolang</span>
          </div>
          
          <n-menu
            mode="horizontal"
            :value="activeMenu"
            :options="menuOptions"
            @update:value="handleMenuClick"
            class="nav-menu"
          />
          
          <div class="search-box">
            <n-input
              v-model:value="searchKeyword"
              placeholder="搜索文章..."
              @keyup.enter="handleSearch"
              clearable
            >
              <template #prefix>
                <n-icon><search-icon /></n-icon>
              </template>
            </n-input>
          </div>
        </div>
      </n-layout-header>
      
      <n-layout-content content-style="padding: 24px;" class="main-content">
        <div class="content-wrapper">
          <router-view />
        </div>
      </n-layout-content>
      
      <n-layout-footer bordered class="footer">
        <div class="footer-content">
          <p>Powered by <strong>HiGolang</strong> — Go 语言新闻聚合博客</p>
          <p class="footer-meta">
            <go-version />
          </p>
        </div>
      </n-layout-footer>
    </n-layout>
  </div>
</template>

<script setup>
import { ref, computed, h, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { NIcon, useMessage } from 'naive-ui'
import { Search as SearchIcon } from '@vicons/ionicons5'
import GoVersion from '../../components/GoVersion.vue'
import { getCategories } from '../../api'

const router = useRouter()
const route = useRoute()
const message = useMessage()

const searchKeyword = ref('')
const categories = ref([])

const activeMenu = computed(() => {
  if (route.name === 'Home') return 'home'
  if (route.name === 'Category') return `category-${route.params.slug}`
  if (route.name === 'Tag') return 'tag'
  if (route.name === 'About') return 'about'
  return null
})

const menuOptions = computed(() => {
  const options = [
    { label: '首页', key: 'home' }
  ]
  
  if (categories.value.length > 0) {
    options.push({
      label: '分类',
      key: 'categories',
      children: categories.value.map(cat => ({
        label: cat.name,
        key: `category-${cat.slug}`
      }))
    })
  }
  
  options.push(
    { label: '标签', key: 'tag' },
    { label: '关于', key: 'about' }
  )
  
  return options
})

const handleMenuClick = (key) => {
  if (key === 'home') {
    router.push({ name: 'Home' })
  } else if (key.startsWith('category-')) {
    const slug = key.replace('category-', '')
    router.push({ name: 'Category', params: { slug } })
  } else if (key === 'tag') {
    router.push({ name: 'Tag', params: { slug: 'all' } })
  } else if (key === 'about') {
    router.push({ name: 'About' })
  }
}

const handleSearch = () => {
  if (searchKeyword.value.trim()) {
    router.push({ name: 'Search', query: { q: searchKeyword.value.trim() } })
  }
}

onMounted(async () => {
  try {
    const response = await getCategories()
    if (response.data.code === 0) {
      categories.value = response.data.data
    }
  } catch (error) {
    console.error('Failed to load categories:', error)
  }
})
</script>

<style scoped>
.front-layout {
  min-height: 100vh;
  background-color: #f5f7fa;
}

.header {
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
  height: 64px;
  display: flex;
  align-items: center;
  gap: 32px;
}

.logo {
  cursor: pointer;
  flex-shrink: 0;
}

.logo-text {
  font-size: 24px;
  font-weight: bold;
  color: #00ADD8;
  letter-spacing: -0.5px;
}

.nav-menu {
  flex: 1;
}

.search-box {
  width: 240px;
  flex-shrink: 0;
}

.main-content {
  min-height: calc(100vh - 64px - 80px);
}

.content-wrapper {
  max-width: 1200px;
  margin: 0 auto;
}

.footer {
  background-color: #fff;
  padding: 24px;
  text-align: center;
}

.footer-content {
  max-width: 1200px;
  margin: 0 auto;
  color: #666;
  font-size: 14px;
}

.footer-content p {
  margin: 4px 0;
}

.footer-meta {
  font-size: 12px;
  color: #999;
}
</style>
