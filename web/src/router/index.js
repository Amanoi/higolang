import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

import FrontLayout from '../views/layout/FrontLayout.vue'
import Home from '../views/Home.vue'
import ArticleDetail from '../views/ArticleDetail.vue'
import Category from '../views/Category.vue'
import Tag from '../views/Tag.vue'
import Search from '../views/Search.vue'
import About from '../views/About.vue'

import Login from '../views/admin/Login.vue'
import AdminLayout from '../views/layout/AdminLayout.vue'
import Dashboard from '../views/admin/Dashboard.vue'
import ArticleList from '../views/admin/ArticleList.vue'
import ArticleEdit from '../views/admin/ArticleEdit.vue'
import FeedSources from '../views/admin/FeedSources.vue'
import FetchLogs from '../views/admin/FetchLogs.vue'
import Settings from '../views/admin/Settings.vue'

const routes = [
  {
    path: '/',
    component: FrontLayout,
    children: [
      { path: '', name: 'Home', component: Home },
      { path: 'article/:slug', name: 'Article', component: ArticleDetail },
      { path: 'category/:slug', name: 'Category', component: Category },
      { path: 'tag/:slug', name: 'Tag', component: Tag },
      { path: 'search', name: 'Search', component: Search },
      { path: 'about', name: 'About', component: About },
    ]
  },
  {
    path: '/admin/login',
    name: 'AdminLogin',
    component: Login
  },
  {
    path: '/admin',
    component: AdminLayout,
    meta: { requiresAuth: true },
    children: [
      { path: '', name: 'Dashboard', component: Dashboard },
      { path: 'articles', name: 'ArticleManage', component: ArticleList },
      { path: 'articles/create', name: 'ArticleCreate', component: ArticleEdit },
      { path: 'articles/:id/edit', name: 'ArticleUpdate', component: ArticleEdit },
      { path: 'feed-sources', name: 'FeedSources', component: FeedSources },
      { path: 'fetch-logs', name: 'FetchLogs', component: FetchLogs },
      { path: 'settings', name: 'Settings', component: Settings },
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  }
})

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    next({ name: 'AdminLogin', query: { redirect: to.fullPath } })
  } else if (to.name === 'AdminLogin' && authStore.isLoggedIn) {
    next({ name: 'Dashboard' })
  } else {
    next()
  }
})

export default router
