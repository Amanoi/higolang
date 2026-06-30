import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const FrontLayout = () => import('../views/layout/FrontLayout.vue')
const Home = () => import('../views/Home.vue')
const ArticleDetail = () => import('../views/ArticleDetail.vue')
const Category = () => import('../views/Category.vue')
const Tag = () => import('../views/Tag.vue')
const Search = () => import('../views/Search.vue')
const About = () => import('../views/About.vue')

const Login = () => import('../views/admin/Login.vue')
const AdminLayout = () => import('../views/layout/AdminLayout.vue')
const Dashboard = () => import('../views/admin/Dashboard.vue')
const ArticleList = () => import('../views/admin/ArticleList.vue')
const ArticleEdit = () => import('../views/admin/ArticleEdit.vue')
const FeedSources = () => import('../views/admin/FeedSources.vue')
const FetchLogs = () => import('../views/admin/FetchLogs.vue')
const Settings = () => import('../views/admin/Settings.vue')

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
