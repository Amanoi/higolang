<template>
  <div class="feed-sources-page">
    <n-card>
      <div class="toolbar">
        <h3 style="margin: 0;">资讯源管理</h3>
        <n-button type="primary" @click="openModal(null)">添加数据源</n-button>
      </div>

      <n-data-table
        :columns="columns"
        :data="sources"
        :loading="loading"
        :bordered="false"
        style="margin-top: 16px;"
      />
    </n-card>

    <!-- Add/Edit Modal -->
    <n-modal v-model:show="showModal" preset="dialog" :title="editingSource ? '编辑数据源' : '添加数据源'" style="width: 500px;">
      <n-form ref="modalFormRef" :model="modalForm" :rules="modalRules" label-placement="left" label-width="80" style="margin-top: 16px;">
        <n-form-item label="名称" path="name">
          <n-input v-model:value="modalForm.name" placeholder="数据源名称" />
        </n-form-item>
        <n-form-item label="URL" path="url">
          <n-input v-model:value="modalForm.url" placeholder="RSS/Atom/JSON 地址" />
        </n-form-item>
        <n-form-item label="类型" path="feed_type">
          <n-select v-model:value="modalForm.feed_type" :options="typeOptions" placeholder="选择类型" />
        </n-form-item>
        <n-form-item label="分类">
          <n-select v-model:value="modalForm.category_id" :options="categoryOptions" placeholder="选择分类" clearable />
        </n-form-item>
        <n-form-item label="间隔(分)">
          <n-input-number v-model:value="modalForm.fetch_interval" :min="5" :max="1440" />
        </n-form-item>
        <n-form-item label="启用">
          <n-switch v-model:value="modalForm.is_enabled" />
        </n-form-item>
      </n-form>
      <template #action>
        <n-button @click="showModal = false">取消</n-button>
        <n-button type="primary" :loading="saving" @click="handleSave">保存</n-button>
      </template>
    </n-modal>

    <!-- Delete Confirmation -->
    <n-modal v-model:show="showDeleteModal" preset="dialog" title="确认删除" positive-text="删除" negative-text="取消" @positive-click="confirmDelete">
      <p>确定要删除数据源「{{ deleteTarget?.name }}」吗？</p>
    </n-modal>
  </div>
</template>

<script setup>
import { ref, h, onMounted } from 'vue'
import { NButton, NTag, NSpace, NSwitch } from 'naive-ui'
import { getFeedSources, createFeedSource, updateFeedSource, deleteFeedSource, fetchFeedSource, getAdminCategories } from '../../api'

const sources = ref([])
const loading = ref(false)
const saving = ref(false)
const showModal = ref(false)
const showDeleteModal = ref(false)
const editingSource = ref(null)
const deleteTarget = ref(null)
const modalFormRef = ref(null)
const categoryOptions = ref([])

const typeOptions = [
  { label: 'RSS', value: 'rss' },
  { label: 'Atom', value: 'atom' },
  { label: 'JSON', value: 'json' }
]

const modalForm = ref({
  name: '',
  url: '',
  feed_type: 'rss',
  category_id: null,
  fetch_interval: 30,
  is_enabled: true
})

const modalRules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  url: [{ required: true, message: '请输入 URL', trigger: 'blur' }],
  feed_type: [{ required: true, message: '请选择类型', trigger: 'change' }]
}

const formatDate = (ts) => {
  if (!ts) return '从未'
  return new Date(ts).toLocaleString('zh-CN')
}

const columns = [
  { title: '名称', key: 'name', width: 160 },
  { title: 'URL', key: 'url', ellipsis: { tooltip: true } },
  { title: '类型', key: 'feed_type', width: 70, render(row) {
    return h(NTag, { size: 'small' }, { default: () => row.feed_type.toUpperCase() })
  }},
  { title: '间隔', key: 'fetch_interval', width: 70, render(row) { return `${row.fetch_interval}分` }},
  { title: '上次抓取', key: 'last_fetched_at', width: 140, render(row) { return formatDate(row.last_fetched_at) }},
  { title: '状态', key: 'is_enabled', width: 70, render(row) {
    return row.is_enabled
      ? h(NTag, { size: 'small', type: 'success' }, { default: () => '启用' })
      : h(NTag, { size: 'small', type: 'default' }, { default: () => '禁用' })
  }},
  { title: '操作', key: 'actions', width: 200, render(row) {
    return h(NSpace, { size: 'small' }, {
      default: () => [
        h(NButton, { size: 'small', type: 'info', loading: row._fetching, onClick: () => handleFetch(row) },
          { default: () => '抓取' }),
        h(NButton, { size: 'small', onClick: () => openModal(row) },
          { default: () => '编辑' }),
        h(NButton, { size: 'small', type: 'error', onClick: () => handleDelete(row) },
          { default: () => '删除' })
      ]
    })
  }}
]

const fetchSources = async () => {
  loading.value = true
  try {
    const res = await getFeedSources()
    if (res.data.code === 0) {
      sources.value = res.data.data.map(s => ({ ...s, _fetching: false }))
    }
  } catch (err) {
    console.error('Fetch sources failed:', err)
  } finally {
    loading.value = false
  }
}

const loadCategories = async () => {
  try {
    const res = await getAdminCategories()
    if (res.data.code === 0) {
      categoryOptions.value = res.data.data.map(c => ({ label: c.name, value: c.id }))
    }
  } catch (err) {
    console.error('Load categories failed:', err)
  }
}

const openModal = (source) => {
  editingSource.value = source
  if (source) {
    modalForm.value = { ...source }
  } else {
    modalForm.value = { name: '', url: '', feed_type: 'rss', category_id: null, fetch_interval: 30, is_enabled: true }
  }
  showModal.value = true
}

const handleSave = async () => {
  try {
    await modalFormRef.value?.validate()
  } catch { return }

  saving.value = true
  try {
    if (editingSource.value) {
      await updateFeedSource(editingSource.value.id, modalForm.value)
    } else {
      await createFeedSource(modalForm.value)
    }
    showModal.value = false
    fetchSources()
  } catch (err) {
    console.error('Save failed:', err)
  } finally {
    saving.value = false
  }
}

const handleFetch = async (row) => {
  row._fetching = true
  try {
    const res = await fetchFeedSource(row.id)
    if (res.data.code === 0) {
      const data = res.data.data
      if (data.error) {
        alert(`抓取完成，新增 ${data.fetched} 篇（有错误: ${data.error}）`)
      } else {
        alert(`抓取完成，新增 ${data.fetched} 篇文章`)
      }
      fetchSources()
    }
  } catch (err) {
    console.error('Fetch failed:', err)
  } finally {
    row._fetching = false
  }
}

const handleDelete = (row) => {
  deleteTarget.value = row
  showDeleteModal.value = true
}

const confirmDelete = async () => {
  if (!deleteTarget.value) return
  try {
    await deleteFeedSource(deleteTarget.value.id)
    fetchSources()
  } catch (err) {
    console.error('Delete failed:', err)
  }
}

onMounted(() => {
  fetchSources()
  loadCategories()
})
</script>

<style scoped>
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
