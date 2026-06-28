import axios from 'axios'
import { useAuthStore } from '../stores/auth'

const api = axios.create({
  baseURL: '',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor
api.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    if (authStore.token) {
      config.headers.Authorization = `Bearer ${authStore.token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      const authStore = useAuthStore()
      authStore.logout()
    }
    return Promise.reject(error)
  }
)

// Public API
export const getArticles = (params) => api.get('/api/v1/articles', { params })
export const getArticleBySlug = (slug) => api.get(`/api/v1/articles/${slug}`)
export const getCategories = () => api.get('/api/v1/categories')
export const getTags = () => api.get('/api/v1/tags')
export const getGoVersion = () => api.get('/api/v1/go-version')
export const getPublicSettings = () => api.get('/api/v1/settings/public')

// Admin API
export const adminLogin = (data) => api.post('/api/v1/admin/login', data)
export const getDashboard = () => api.get('/api/v1/admin/dashboard')
export const getAdminArticles = (params) => api.get('/api/v1/admin/articles', { params })
export const getAdminArticle = (id) => api.get(`/api/v1/admin/articles/${id}`)
export const createArticle = (data) => api.post('/api/v1/admin/articles', data)
export const updateArticle = (id, data) => api.put(`/api/v1/admin/articles/${id}`, data)
export const deleteArticle = (id) => api.delete(`/api/v1/admin/articles/${id}`)

export const getAdminCategories = () => api.get('/api/v1/admin/categories')
export const createCategory = (data) => api.post('/api/v1/admin/categories', data)
export const updateCategory = (id, data) => api.put(`/api/v1/admin/categories/${id}`, data)
export const deleteCategory = (id) => api.delete(`/api/v1/admin/categories/${id}`)

export const getAdminTags = () => api.get('/api/v1/admin/tags')
export const createTag = (data) => api.post('/api/v1/admin/tags', data)
export const updateTag = (id, data) => api.put(`/api/v1/admin/tags/${id}`, data)
export const deleteTag = (id) => api.delete(`/api/v1/admin/tags/${id}`)

export const getFeedSources = () => api.get('/api/v1/admin/feed-sources')
export const createFeedSource = (data) => api.post('/api/v1/admin/feed-sources', data)
export const updateFeedSource = (id, data) => api.put(`/api/v1/admin/feed-sources/${id}`, data)
export const deleteFeedSource = (id) => api.delete(`/api/v1/admin/feed-sources/${id}`)
export const fetchFeedSource = (id) => api.post(`/api/v1/admin/feed-sources/${id}/fetch`)

export const getFetchLogs = (params) => api.get('/api/v1/admin/fetch-logs', { params })

export const getAdminSettings = () => api.get('/api/v1/admin/settings')
export const updateAdminSettings = (data) => api.put('/api/v1/admin/settings', data)

export default api
