<template>
  <div class="article-list-page">
    <!-- Toolbar -->
    <n-card>
      <div class="toolbar">
        <div class="toolbar-left">
          <n-input
            v-model:value="searchKeyword"
            placeholder="搜索文章标题..."
            clearable
            style="width: 240px;"
            @clear="handleSearch"
            @keyup.enter="handleSearch"
          />
          <n-select
            v-model:value="statusFilter"
            :options="statusOptions"
            style="width: 140px;"
            @update:value="handleSearch"
          />
          <n-button type="primary" @click="handleSearch">搜索</n-button>
        </div>
        <n-button type="primary" @click="$router.push({ name: 'ArticleCreate' })">
          新建文章
        </n-button>
      </div>
    </n-card>

    <!-- Article Table -->
    <n-card style="margin-top: 16px;">
      <n-data-table
        :columns="columns"
        :data="articles"
        :loading="loading"
        :bordered="false"
        :row-key="row => row.id"
      />

      <div class="pagination-wrapper" v-if="total > pageSize">
        <n-pagination
          v-model:page="currentPage"
          :page-count="Math.ceil(total / pageSize)"
          show-size-picker
          :page-sizes="[10, 20, 50]"
          @update:page="handlePageChange"
          @update:page-size="handlePageSizeChange"
        >
          <template #prefix="{ itemCount }">共 {{ itemCount }} 篇</template>
        </n-pagination>
      </div>
    </n-card>

    <!-- Delete Confirmation -->
    <n-modal v-model:show="showDeleteModal" preset="dialog" title="确认删除" positive-text="删除" negative-text="取消" @positive-click="confirmDelete">
      <p>确定要删除文章「{{ deleteTarget?.title }}」吗？此操作不可恢复。</p>
    </n-modal>
  </div>
</template>

<script setup>
import { ref, h, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { NButton, NTag, NSpace, NSwitch } from 'naive-ui'
import { getAdminArticles, deleteArticle } from '../../api'

const router = useRouter()

const articles = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const loading = ref(false)
const searchKeyword = ref('')
const statusFilter = ref(null)
const showDeleteModal = ref(false)
const deleteTarget = ref(null)

const statusOptions = [
  { label: '全部', value: null },
  { label: '已发布', value: 'published' },
  { label: '草稿', value: 'draft' }
]

const formatDate = (ts) => {
  if (!ts) return '-'
  return new Date(ts).toLocaleDateString('zh-CN')
}

const columns = [
  { title: '标题', key: 'title', ellipsis: { tooltip: true }, render(row) {
    return h(NButton, { text: true, type: 'primary', onClick: () => router.push({ name: 'ArticleUpdate', params: { id: row.id } }) },
      { default: () => row.title })
  }},
  { title: '分类', key: 'category', width: 100, render(row) {
    return row.category ? h(NTag, { size: 'small', type: 'info' }, { default: () => row.category.name }) : '-'
  }},
  { title: '状态', key: 'status', width: 80, render(row) {
    return row.status === 'published'
      ? h(NTag, { size: 'small', type: 'success' }, { default: () => '已发布' })
      : h(NTag, { size: 'small', type: 'warning' }, { default: () => '草稿' })
  }},
  { title: '置顶', key: 'is_pinned', width: 60, render(row) {
    return row.is_pinned ? h(NTag, { size: 'small', type: 'warning' }, { default: () => '是' }) : '-'
  }},
  { title: '浏览', key: 'view_count', width: 60 },
  { title: '发布时间', key: 'published_at', width: 120, render(row) { return formatDate(row.published_at) }},
  { title: '操作', key: 'actions', width: 140, render(row) {
    return h(NSpace, { size: 'small' }, {
      default: () => [
        h(NButton, { size: 'small', onClick: () => router.push({ name: 'ArticleUpdate', params: { id: row.id } }) },
          { default: () => '编辑' }),
        h(NButton, { size: 'small', type: 'error', onClick: () => handleDelete(row) },
          { default: () => '删除' })
      ]
    })
  }}
]

const fetchArticles = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    if (searchKeyword.value) params.search = searchKeyword.value
    if (statusFilter.value) params.status = statusFilter.value

    const res = await getAdminArticles(params)
    if (res.data.code === 0) {
      articles.value = res.data.data.list
      total.value = res.data.data.total
    }
  } catch (err) {
    console.error('Fetch articles failed:', err)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchArticles()
}

const handlePageChange = (page) => {
  currentPage.value = page
  fetchArticles()
}

const handlePageSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
  fetchArticles()
}

const handleDelete = (row) => {
  deleteTarget.value = row
  showDeleteModal.value = true
}

const confirmDelete = async () => {
  if (!deleteTarget.value) return
  try {
    await deleteArticle(deleteTarget.value.id)
    fetchArticles()
  } catch (err) {
    console.error('Delete failed:', err)
  }
}

onMounted(fetchArticles)
</script>

<style scoped>
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.pagination-wrapper {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
