<template>
  <n-card class="tag-cloud">
    <template #header>
      <span class="card-header">标签云</span>
    </template>
    
    <div class="tags-container">
      <n-tag
        v-for="tag in tags"
        :key="tag.id"
        :bordered="false"
        :style="{ fontSize: getTagSize(tag.article_count) + 'px', cursor: 'pointer' }"
        @click="goToTag(tag)"
        class="tag-item"
      >
        {{ tag.name }}
        <template v-if="tag.article_count" #avatar>
          <span class="tag-count">{{ tag.article_count }}</span>
        </template>
      </n-tag>
    </div>
  </n-card>
</template>

<script setup>
import { useRouter } from 'vue-router'

const props = defineProps({
  tags: {
    type: Array,
    default: () => []
  }
})

const router = useRouter()

const getTagSize = (count) => {
  if (!count) return 12
  if (count > 20) return 16
  if (count > 10) return 14
  return 12
}

const goToTag = (tag) => {
  router.push({ name: 'Tag', params: { slug: tag.slug } })
}
</script>

<style scoped>
.tag-cloud {
  margin-bottom: 24px;
}

.card-header {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.tags-container {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag-item {
  transition: all 0.2s ease;
}

.tag-item:hover {
  transform: scale(1.05);
  background-color: #00ADD8 !important;
  color: #fff !important;
}

.tag-count {
  font-size: 10px;
  color: #999;
}
</style>
