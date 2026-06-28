<template>
  <div class="fetch-logs-page">
    <n-card>
      <div class="toolbar">
        <h3 style="margin: 0;">抓取日志</h3>
        <n-select
          v-model:value="sourceFilter"
          :options="sourceOptions"
          placeholder="按数据源筛选"
          clearable
          style="width: 240px;"
          @update:value="handleFilterChange"
        />
      </div>

      <n-data-table
        :columns="columns"
        :data="logs"
        :loading="loading"
        :bordered="false"
        style="margin-top: 16px;"
      />

      <div class="pagination-wrapper" v-if="total > pageSize">
        <n-pagination
          v-model:page="currentPage"
          :page-count="Math.ceil(total / pageSize)"
          @update:page="handlePageChange"
        />
      </div>
    </n-card>
  </div>
</template>

<script setup>
import { ref, h, onMounted } from 'vue'
import { NTag } from 'naive-ui'
import { getFetchLogs, getFeedSources } from '../../api'

const logs = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const loading = ref(false)
const sourceFilter = ref(null)
const sourceOptions = ref([])

const formatDate = (ts) => {
  if (!ts) return '-'
  return new Date(ts).toLocaleString('zh-CN')
}

const columns = [
  { title: '状态', key: 'status', width: 80, render(row) {
    return row.status === 'success'
      ? h(NTag, { type: 'success', size: 'small' }, { default: () => '成功' })
      : h(NTag, { type: 'error', size: 'small' }, { default: () => '失败' })
  }},
  { title: '抓取数量', key: 'fetched_count', width: 100 },
  { title: '消息', key: 'message', ellipsis: { tooltip: true }, render(row) { return row.message || '-' }},
  { title: '耗时', key: 'duration_ms', width: 100, render(row) { return `${row.duration_ms} ms` }},
  { title: '时间', key: 'created_at', width: 180, render(row) { return formatDate(row.created_at) }}
]

const fetchLogs = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    if (sourceFilter.value) params.source_id = sourceFilter.value

    const res = await getFetchLogs(params)
    if (res.data.code === 0) {
      logs.value = res.data.data.list
      total.value = res.data.data.total
    }
  } catch (err) {
    console.error('Fetch logs failed:', err)
  } finally {
    loading.value = false
  }
}

const loadSourceOptions = async () => {
  try {
    const res = await getFeedSources()
    if (res.data.code === 0) {
      sourceOptions.value = res.data.data.map(s => ({ label: s.name, value: s.id }))
    }
  } catch (err) {
    console.error('Load sources failed:', err)
  }
}

const handleFilterChange = () => {
  currentPage.value = 1
  fetchLogs()
}

const handlePageChange = (page) => {
  currentPage.value = page
  fetchLogs()
}

onMounted(() => {
  loadSourceOptions()
  fetchLogs()
})
</script>

<style scoped>
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.pagination-wrapper {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
