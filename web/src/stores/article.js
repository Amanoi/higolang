import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../api'

export const useArticleStore = defineStore('article', () => {
  const articles = ref([])
  const total = ref(0)
  const currentPage = ref(1)
  const pageSize = ref(10)
  const loading = ref(false)

  const fetchArticles = async (params = {}) => {
    loading.value = true
    try {
      const response = await api.get('/api/v1/articles', {
        params: {
          page: currentPage.value,
          page_size: pageSize.value,
          ...params
        }
      })
      if (response.data.code === 0) {
        articles.value = response.data.data.list
        total.value = response.data.data.total
      }
    } finally {
      loading.value = false
    }
  }

  const fetchArticleBySlug = async (slug) => {
    const response = await api.get(`/api/v1/articles/${slug}`)
    if (response.data.code === 0) {
      return response.data.data
    }
    throw new Error(response.data.message || '获取文章失败')
  }

  return {
    articles,
    total,
    currentPage,
    pageSize,
    loading,
    fetchArticles,
    fetchArticleBySlug
  }
})
