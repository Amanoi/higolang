<template>
  <div class="tag-page">
    <div class="page-header">
      <h1 class="page-title">标签：{{ tagName }}</h1>
    </div>

    <n-spin :show="loading">
      <div v-if="articles.length > 0" class="articles-list">
        <ArticleCard
          v-for="article in articles"
          :key="article.id"
          :article="article"
        />
      </div>
      <n-empty v-else description="该标签下暂无文章" />
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
import { getArticles, getTags } from '../api'

const route = useRoute()

const articles = ref([])
const tagName = ref('')
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
      tag: route.params.slug
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

const fetchTagName = async () => {
  try {
    const response = await getTags()
    if (response.data.code === 0) {
      const tag = response.data.data.find(t => t.slug === route.params.slug)
      if (tag) {
        tagName.value = tag.name
      }
    }
  } catch (error) {
    console.error('Failed to fetch tag:', error)
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
  fetchTagName()
})

onMounted(() => {
  fetchArticles()
  fetchTagName()
})
</script>

<style scoped>
.tag-page {
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
