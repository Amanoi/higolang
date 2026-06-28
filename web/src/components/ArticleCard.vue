<template>
  <n-card class="article-card" hoverable @click="goToArticle">
    <div class="card-content">
      <div class="card-header">
        <h3 class="card-title">{{ article.title }}</h3>
        <n-tag v-if="article.is_pinned" type="warning" size="small" class="pin-tag">置顶</n-tag>
      </div>
      
      <p class="card-summary">{{ article.summary }}</p>
      
      <div class="card-meta">
        <div class="meta-left">
          <n-tag v-if="article.category" :bordered="false" size="small" type="info">
            {{ article.category.name }}
          </n-tag>
          
          <div v-if="article.tags && article.tags.length > 0" class="tags">
            <n-tag
              v-for="tag in article.tags.slice(0, 3)"
              :key="tag.id"
              :bordered="false"
              size="small"
            >
              {{ tag.name }}
            </n-tag>
          </div>
        </div>
        
        <div class="meta-right">
          <span class="date">{{ formatDate(article.published_at) }}</span>
          <span v-if="article.source_name" class="source">· {{ article.source_name }}</span>
        </div>
      </div>
    </div>
  </n-card>
</template>

<script setup>
import { useRouter } from 'vue-router'

const props = defineProps({
  article: {
    type: Object,
    required: true
  }
})

const router = useRouter()

const goToArticle = () => {
  router.push({ name: 'Article', params: { slug: props.article.slug } })
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })
}
</script>

<style scoped>
.article-card {
  cursor: pointer;
  transition: all 0.3s ease;
  margin-bottom: 16px;
}

.article-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 16px rgba(0, 173, 216, 0.15);
}

.card-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.card-header {
  display: flex;
  align-items: flex-start;
  gap: 8px;
}

.card-title {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin: 0;
  flex: 1;
  line-height: 1.4;
}

.pin-tag {
  flex-shrink: 0;
}

.card-summary {
  font-size: 14px;
  color: #666;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  margin: 0;
}

.card-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
}

.meta-left {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.tags {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.meta-right {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #999;
}

.date, .source {
  white-space: nowrap;
}

@media (max-width: 768px) {
  .card-meta {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
