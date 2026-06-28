<template>
  <div class="dashboard-page">
    <n-spin :show="loading">
      <!-- Stats Cards -->
      <n-grid :cols="4" :x-gap="16" :y-gap="16" class="stats-grid">
        <n-gi>
          <n-card>
            <n-statistic label="文章总数" :value="stats.total_articles">
              <template #prefix>
                <n-icon :size="20" color="#00ADD8"><document-icon /></n-icon>
              </template>
            </n-statistic>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card>
            <n-statistic label="本周新增" :value="stats.week_count">
              <template #prefix>
                <n-icon :size="20" color="#18a058"><add-icon /></n-icon>
              </template>
            </n-statistic>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card>
            <n-statistic label="启用数据源" :value="stats.enabled_sources">
              <template #prefix>
                <n-icon :size="20" color="#f0a020"><rss-icon /></n-icon>
              </template>
            </n-statistic>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card>
            <n-statistic label="禁用数据源" :value="stats.disabled_sources">
              <template #prefix>
                <n-icon :size="20" color="#999"><pause-icon /></n-icon>
              </template>
            </n-statistic>
          </n-card>
        </n-gi>
      </n-grid>

      <n-grid :cols="2" :x-gap="16" :y-gap="16" style="margin-top: 16px;">
        <!-- Category Distribution -->
        <n-gi>
          <n-card title="分类统计">
            <div v-if="stats.category_stats" class="cat-stats">
              <div v-for="cat in stats.category_stats" :key="cat.name" class="cat-row">
                <span class="cat-name">{{ cat.name }}</span>
                <n-progress
                  type="line"
                  :percentage="getCatPercentage(cat.count)"
                  :show-indicator="false"
                  color="#00ADD8"
                  style="flex: 1; margin: 0 12px;"
                />
                <span class="cat-count">{{ cat.count }}</span>
              </div>
              <n-empty v-if="!stats.category_stats?.length" description="暂无数据" />
            </div>
          </n-card>
        </n-gi>

        <!-- Recent Fetch Logs -->
        <n-gi>
          <n-card title="最近抓取记录">
            <n-data-table
              :columns="logColumns"
              :data="stats.recent_logs || []"
              :bordered="false"
              :single-line="false"
              size="small"
              :max-height="300"
            />
            <n-empty v-if="!stats.recent_logs?.length" description="暂无抓取记录" />
          </n-card>
        </n-gi>
      </n-grid>
    </n-spin>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { NTag } from 'naive-ui'
import { DocumentTextOutline as DocumentIcon, AddCircleOutline as AddIcon, RssOutline as RssIcon, PauseCircleOutline as PauseIcon } from '@vicons/ionicons5'
import { getDashboard } from '../../api'

const loading = ref(false)
const stats = reactive({
  total_articles: 0,
  week_count: 0,
  enabled_sources: 0,
  disabled_sources: 0,
  category_stats: [],
  recent_logs: []
})

const maxCatCount = computed(() => {
  if (!stats.category_stats?.length) return 1
  return Math.max(...stats.category_stats.map(c => c.count), 1)
})

const getCatPercentage = (count) => {
  return Math.round((count / maxCatCount.value) * 100)
}

const formatDate = (ts) => {
  if (!ts) return '-'
  return new Date(ts).toLocaleString('zh-CN', { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })
}

const logColumns = [
  { title: '状态', key: 'status', width: 80, render(row) {
    return row.status === 'success'
      ? h(NTag, { type: 'success', size: 'small' }, { default: () => '成功' })
      : h(NTag, { type: 'error', size: 'small' }, { default: () => '失败' })
  }},
  { title: '数量', key: 'fetched_count', width: 60 },
  { title: '耗时', key: 'duration_ms', width: 80, render(row) { return `${row.duration_ms}ms` } },
  { title: '时间', key: 'created_at', width: 100, render(row) { return formatDate(row.created_at) } },
]

const fetchDashboard = async () => {
  loading.value = true
  try {
    const res = await getDashboard()
    if (res.data.code === 0) {
      Object.assign(stats, res.data.data)
    }
  } catch (err) {
    console.error('Dashboard fetch failed:', err)
  } finally {
    loading.value = false
  }
}

import { h } from 'vue'

onMounted(fetchDashboard)
</script>

<style scoped>
.dashboard-page {
  padding: 0;
}

.stats-grid {
  margin-bottom: 0;
}

.cat-stats {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.cat-row {
  display: flex;
  align-items: center;
}

.cat-name {
  width: 80px;
  font-size: 14px;
  color: #333;
  flex-shrink: 0;
}

.cat-count {
  width: 40px;
  text-align: right;
  font-size: 14px;
  font-weight: 600;
  color: #00ADD8;
  flex-shrink: 0;
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr) !important;
  }
}
</style>
