<template>
  <div class="search-page">
    <div class="search-header">
      <n-input
        v-model:value="keyword"
        placeholder="搜索文章..."
        size="large"
        clearable
        @keyup.enter="handleSearch"
        @clear="handleSearch"
      >
        <template #prefix>
          <n-icon><search-icon /></n-icon>
        </template>
        <template #suffix>
          <n-button text @click="handleSearch">搜索</n-button>
        </template>
      </n-input>
    </div>

    <div class="search-results">
      <n-spin :show="loading">
        <template v-if="searched">
          <p class="result-count" v-if="total > 0">
            找到 <strong>{{ total }}</strong> 篇与「<strong>{{ displayKeyword }}</strong>」相关的文章
          </p>
          <div v-if="articles.length > 0" class="articles-list">
            <ArticleCard
              v-for="article in articles"
              :key="article.id"
              :article="article"
            />
          </div>
          <n-empty v-else description="没有找到相关文章" />
        </template>
        <template v-else>
          <n-empty description="输入关键词开始搜索" />
        </template>
      </n-spin>
    </div>

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
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Search as SearchIcon } from '@vicons/ionicons5'
import ArticleCard from '../components/ArticleCard.vue'
import { getArticles } from '../api'

const route = useRoute()
const router = useRouter()

const keyword = ref('')
const displayKeyword = ref('')
const articles = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const loading = ref(false)
const searched = ref(false)

const fetchArticles = async () => {
  if (!keyword.value.trim()) return
  loading.value = true
  searched.value = true
  try {
    const response = await getArticles({
      page: currentPage.value,
      page_size: pageSize.value,
      search: keyword.value.trim()
    })
    if (response.data.code === 0) {
      articles.value = response.data.data.list
      total.value = response.data.data.total
      displayKeyword.value = keyword.value.trim()
    }
  } catch (error) {
    console.error('Search failed:', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  router.replace({ name: 'Search', query: { q: keyword.value.trim() } })
  fetchArticles()
}

const handlePageChange = (page) => {
  currentPage.value = page
  fetchArticles()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => {
  if (route.query.q) {
    keyword.value = route.query.q
    fetchArticles()
  }
})

watch(() => route.query.q, (val) => {
  if (val) {
    keyword.value = val
    fetchArticles()
  }
})
</script>

<style scoped>
.search-page {
  padding: 24px 0;
  max-width: 800px;
  margin: 0 auto;
}

.search-header {
  margin-bottom: 32px;
}

.result-count {
  font-size: 14px;
  color: #666;
  margin-bottom: 20px;
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
