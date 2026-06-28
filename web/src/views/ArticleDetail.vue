<template>
  <div class="article-detail">
    <n-spin :show="loading">
      <div v-if="article" class="article-container">
        <!-- Article Header -->
        <div class="article-header">
          <h1 class="article-title">{{ article.title }}</h1>
          
          <div class="article-meta">
            <div class="meta-item">
              <n-icon><calendar-icon /></n-icon>
              <span>{{ formatDate(article.published_at) }}</span>
            </div>
            <div class="meta-item" v-if="article.category">
              <n-icon><folder-icon /></n-icon>
              <router-link :to="{ name: 'Category', params: { slug: article.category.slug } }">
                {{ article.category.name }}
              </router-link>
            </div>
            <div class="meta-item" v-if="article.view_count !== undefined">
              <n-icon><eye-icon /></n-icon>
              <span>{{ article.view_count }} 次阅读</span>
            </div>
            <div class="meta-item" v-if="article.source_url">
              <n-icon><link-icon /></n-icon>
              <a :href="article.source_url" target="_blank" rel="noopener">原文链接</a>
            </div>
          </div>

          <div class="article-tags" v-if="article.tags && article.tags.length > 0">
            <n-tag
              v-for="tag in article.tags"
              :key="tag.id"
              :bordered="false"
              size="small"
              @click="$router.push({ name: 'Tag', params: { slug: tag.slug } })"
              style="cursor: pointer;"
            >
              {{ tag.name }}
            </n-tag>
          </div>
        </div>

        <!-- Article Content -->
        <div class="article-content markdown-body" v-html="renderedContent"></div>

        <!-- Article Actions -->
        <div class="article-actions">
          <n-button @click="copyLink" type="primary" ghost>
            <template #icon>
              <n-icon><copy-icon /></n-icon>
            </template>
            复制链接
          </n-button>
        </div>

        <!-- Navigation -->
        <div class="article-nav" v-if="prevArticle || nextArticle">
          <div class="nav-item prev" v-if="prevArticle">
            <span class="nav-label">上一篇</span>
            <router-link :to="{ name: 'Article', params: { slug: prevArticle.slug } }">
              {{ prevArticle.title }}
            </router-link>
          </div>
          <div class="nav-item next" v-if="nextArticle">
            <span class="nav-label">下一篇</span>
            <router-link :to="{ name: 'Article', params: { slug: nextArticle.slug } }">
              {{ nextArticle.title }}
            </router-link>
          </div>
        </div>
      </div>

      <n-empty v-else-if="!loading" description="文章不存在" />
    </n-spin>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useMessage } from 'naive-ui'
import {
  CalendarOutline as CalendarIcon,
  FolderOutline as FolderIcon,
  EyeOutline as EyeIcon,
  LinkOutline as LinkIcon,
  CopyOutline as CopyIcon
} from '@vicons/ionicons5'
import { marked } from 'marked'
import hljs from 'highlight.js'
import 'highlight.js/styles/atom-one-dark.css'
import { getArticleBySlug } from '../api'

const route = useRoute()
const message = useMessage()

const article = ref(null)
const loading = ref(false)
const prevArticle = ref(null)
const nextArticle = ref(null)

const renderedContent = computed(() => {
  if (!article.value?.content) return ''
  
  marked.setOptions({
    highlight: (code, lang) => {
      if (lang && hljs.getLanguage(lang)) {
        return hljs.highlight(code, { language: lang }).value
      }
      return hljs.highlightAuto(code).value
    }
  })
  
  return marked(article.value.content)
})

const fetchArticle = async () => {
  loading.value = true
  try {
    const response = await getArticleBySlug(route.params.slug)
    if (response.data.code === 0) {
      article.value = response.data.data
      prevArticle.value = response.data.data.prev || null
      nextArticle.value = response.data.data.next || null
    }
  } catch (error) {
    console.error('Failed to fetch article:', error)
    message.error('获取文章失败')
  } finally {
    loading.value = false
  }
}

const copyLink = async () => {
  try {
    await navigator.clipboard.writeText(window.location.href)
    message.success('链接已复制')
  } catch (error) {
    message.error('复制失败')
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

watch(() => route.params.slug, () => {
  fetchArticle()
  window.scrollTo({ top: 0, behavior: 'smooth' })
})

onMounted(() => {
  fetchArticle()
})
</script>

<style scoped>
.article-detail {
  padding: 24px 0;
}

.article-container {
  max-width: 800px;
  margin: 0 auto;
  background: #fff;
  border-radius: 12px;
  padding: 40px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
}

.article-header {
  margin-bottom: 32px;
  padding-bottom: 24px;
  border-bottom: 1px solid #eee;
}

.article-title {
  font-size: 32px;
  font-weight: bold;
  color: #1a1a1a;
  line-height: 1.3;
  margin-bottom: 16px;
}

.article-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  margin-bottom: 16px;
  color: #666;
  font-size: 14px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.meta-item a {
  color: #00ADD8;
}

.meta-item a:hover {
  color: #00C5F5;
}

.article-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.article-content {
  margin-bottom: 32px;
}

.article-actions {
  display: flex;
  justify-content: center;
  margin-bottom: 32px;
}

.article-nav {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  padding-top: 24px;
  border-top: 1px solid #eee;
}

.nav-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.nav-item.next {
  text-align: right;
}

.nav-label {
  font-size: 12px;
  color: #999;
  text-transform: uppercase;
}

.nav-item a {
  font-size: 14px;
  color: #00ADD8;
  line-height: 1.4;
}

.nav-item a:hover {
  color: #00C5F5;
}

@media (max-width: 768px) {
  .article-container {
    padding: 24px;
  }

  .article-title {
    font-size: 24px;
  }

  .article-nav {
    grid-template-columns: 1fr;
  }

  .nav-item.next {
    text-align: left;
  }
}
</style>
