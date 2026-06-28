<template>
  <div class="category-page">
    <div class="page-header">
      <h1 class="page-title">分类：{{ categoryName }}</h1>
    </div>

    <n-spin :show="loading">
      <div v-if="articles.length > 0" class="articles-list">
        <ArticleCard
          v-for="article in articles"
          :key="article.id"
          :article="article"
        />
      </div>
      <n-empty v-else description="该分类下暂无文章" />
    </n-spin>

    <div class="pagination-wrapper" v-if="total > pageSize">
      <n-pagination
        v-model:page="currentPage"
        :page-count="Math.ceil(total / pageSize)"
        @update:page="handlePageChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import ArticleCard from '../components/ArticleCard.vue'
import { getArticles, getCategories } from '../api'

const route = useRoute()

const articles = ref([])
const categoryName = ref('')
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const loading = ref(false)

const fetchArticles = async () => {
  loading.value = true
  try {
    const response = await getArticles({
      page: currentPage.value,
      page_size: pageSize.value,
      category: route.params.slug
    })
    if (response.data.code === 0) {
      articles.value = response.data.data.list
      total.value = response.data.data.total
    }
  } catch (error) {
    console.error('Failed to fetch articles:', error)
  } finally {
    loading.value = false
  }
}

const fetchCategoryName = async () => {
  try {
    const response = await getCategories()
    if (response.data.code === 0) {
      const category = response.data.data.find(c => c.slug === route.params.slug)
      if (category) {
        categoryName.value = category.name
      }
    }
  } catch (error) {
    console.error('Failed to fetch category:', error)
  }
}

const handlePageChange = (page) => {
  currentPage.value = page
  fetchArticles()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

watch(() => route.params.slug, () => {
  currentPage.value = 1
  fetchArticles()
  fetchCategoryName()
})

onMounted(() => {
  fetchArticles()
  fetchCategoryName()
})
</script>

<style scoped>
.category-page {
  padding: 24px 0;
}

.page-header {
  margin-bottom: 24px;
}

.page-title {
  font-size: 28px;
  font-weight: bold;
  color: #333;
}

.articles-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.pagination-wrapper {
  margin-top: 32px;
  display: flex;
  justify-content: center;
}
</style>
