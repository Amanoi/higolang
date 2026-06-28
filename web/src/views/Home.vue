<template>
  <div class="home-page">
    <!-- Hero Section -->
    <div class="hero">
      <h1 class="hero-title">HiGolang</h1>
      <p class="hero-subtitle">Go 语言最新动态</p>
    </div>

    <div class="content-layout">
      <!-- Main Content -->
      <div class="main-column">
        <!-- Pinned Articles -->
        <div v-if="pinnedArticles.length > 0" class="pinned-section">
          <h2 class="section-title">置顶文章</h2>
          <ArticleCard
            v-for="article in pinnedArticles"
            :key="article.id"
            :article="article"
          />
        </div>

        <!-- Article List -->
        <div class="articles-section">
          <h2 class="section-title">最新文章</h2>
          <n-spin :show="loading">
            <div v-if="articles.length > 0" class="articles-grid">
              <ArticleCard
                v-for="article in articles"
                :key="article.id"
                :article="article"
              />
            </div>
            <n-empty v-else description="暂无文章" />
          </n-spin>
        </div>

        <!-- Pagination -->
        <div class="pagination-wrapper" v-if="total > pageSize">
          <n-pagination
            v-model:page="currentPage"
            :page-count="Math.ceil(total / pageSize)"
            @update:page="handlePageChange"
          />
        </div>
      </div>

      <!-- Sidebar -->
      <div class="sidebar-column">
        <GoVersion />
        <TagCloud :tags="tags" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import ArticleCard from '../components/ArticleCard.vue'
import TagCloud from '../components/TagCloud.vue'
import GoVersion from '../components/GoVersion.vue'
import { getArticles, getTags } from '../api'

const articles = ref([])
const pinnedArticles = ref([])
const tags = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const loading = ref(false)

const fetchArticles = async () => {
  loading.value = true
  try {
    const response = await getArticles({
      page: currentPage.value,
      page_size: pageSize.value
    })
    if (response.data.code === 0) {
      const allArticles = response.data.data.list
      pinnedArticles.value = allArticles.filter(a => a.is_pinned)
      articles.value = allArticles.filter(a => !a.is_pinned)
      total.value = response.data.data.total
    }
  } catch (error) {
    console.error('Failed to fetch articles:', error)
  } finally {
    loading.value = false
  }
}

const fetchTags = async () => {
  try {
    const response = await getTags()
    if (response.data.code === 0) {
      tags.value = response.data.data
    }
  } catch (error) {
    console.error('Failed to fetch tags:', error)
  }
}

const handlePageChange = (page) => {
  currentPage.value = page
  fetchArticles()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => {
  fetchArticles()
  fetchTags()
})
</script>

<style scoped>
.home-page {
  padding: 24px 0;
}

.hero {
  text-align: center;
  padding: 48px 24px;
  background: linear-gradient(135deg, #00ADD8 0%, #00C5F5 100%);
  border-radius: 12px;
  margin-bottom: 32px;
  color: #fff;
}

.hero-title {
  font-size: 48px;
  font-weight: bold;
  margin-bottom: 12px;
  letter-spacing: -1px;
}

.hero-subtitle {
  font-size: 20px;
  opacity: 0.9;
}

.content-layout {
  display: grid;
  grid-template-columns: 1fr 300px;
  gap: 32px;
}

.main-column {
  min-width: 0;
}

.sidebar-column {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.section-title {
  font-size: 20px;
  font-weight: 600;
  color: #333;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 2px solid #00ADD8;
}

.pinned-section {
  margin-bottom: 32px;
}

.articles-grid {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.pagination-wrapper {
  margin-top: 32px;
  display: flex;
  justify-content: center;
}

@media (max-width: 768px) {
  .content-layout {
    grid-template-columns: 1fr;
  }

  .hero-title {
    font-size: 36px;
  }

  .hero-subtitle {
    font-size: 16px;
  }
}
</style>
