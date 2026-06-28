<template>
  <div class="article-edit-page">
    <n-card :title="isEdit ? '编辑文章' : '创建文章'">
      <n-spin :show="pageLoading">
        <n-form ref="formRef" :model="form" :rules="rules" label-placement="left" label-width="80">
          <n-grid :cols="2" :x-gap="16">
            <n-gi>
              <n-form-item label="标题" path="title">
                <n-input v-model:value="form.title" placeholder="文章标题" />
              </n-form-item>
            </n-gi>
            <n-gi>
              <n-form-item label="分类" path="category_id">
                <n-select
                  v-model:value="form.category_id"
                  :options="categoryOptions"
                  placeholder="选择分类"
                  clearable
                />
              </n-form-item>
            </n-gi>
          </n-grid>

          <n-form-item label="摘要">
            <n-input
              v-model:value="form.summary"
              type="textarea"
              placeholder="文章摘要（可选，留空自动截取正文前 200 字）"
              :rows="2"
            />
          </n-form-item>

          <n-grid :cols="2" :x-gap="16">
            <n-gi>
              <n-form-item label="标签">
                <n-select
                  v-model:value="form.tag_ids"
                  :options="tagOptions"
                  placeholder="选择标签"
                  multiple
                  clearable
                />
              </n-form-item>
            </n-gi>
            <n-gi>
              <n-form-item label="封面图">
                <n-input v-model:value="form.cover_url" placeholder="封面图 URL（可选）" />
              </n-form-item>
            </n-gi>
          </n-grid>

          <n-form-item label="正文" path="content">
            <div style="width: 100%;">
              <MdEditor
                v-model="form.content"
                :style="{ height: '400px' }"
                theme="light"
              />
            </div>
          </n-form-item>

          <n-grid :cols="2" :x-gap="16">
            <n-gi>
              <n-form-item label="状态">
                <n-radio-group v-model:value="form.status">
                  <n-radio-button value="published">发布</n-radio-button>
                  <n-radio-button value="draft">草稿</n-radio-button>
                </n-radio-group>
              </n-form-item>
            </n-gi>
            <n-gi>
              <n-form-item label="置顶">
                <n-switch v-model:value="form.is_pinned" />
              </n-form-item>
            </n-gi>
          </n-grid>

          <div class="form-actions">
            <n-button @click="$router.back()">取消</n-button>
            <n-button type="primary" :loading="saving" @click="handleSave">
              {{ isEdit ? '保存修改' : '创建文章' }}
            </n-button>
          </div>
        </n-form>
      </n-spin>
    </n-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import { getAdminArticle, createArticle, updateArticle, getAdminCategories, getAdminTags } from '../../api'

const route = useRoute()
const router = useRouter()
const message = useMessage()

const isEdit = computed(() => !!route.params.id)
const formRef = ref(null)
const pageLoading = ref(false)
const saving = ref(false)

const form = reactive({
  title: '',
  summary: '',
  content: '',
  cover_url: '',
  category_id: null,
  tag_ids: [],
  is_pinned: false,
  status: 'published'
})

const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入正文', trigger: 'change' }]
}

const categoryOptions = ref([])
const tagOptions = ref([])

const loadOptions = async () => {
  try {
    const [catRes, tagRes] = await Promise.all([getAdminCategories(), getAdminTags()])
    if (catRes.data.code === 0) {
      categoryOptions.value = catRes.data.data.map(c => ({ label: c.name, value: c.id }))
    }
    if (tagRes.data.code === 0) {
      tagOptions.value = tagRes.data.data.map(t => ({ label: t.name, value: t.id }))
    }
  } catch (err) {
    console.error('Load options failed:', err)
  }
}

const loadArticle = async () => {
  if (!isEdit.value) return
  pageLoading.value = true
  try {
    const res = await getAdminArticle(route.params.id)
    if (res.data.code === 0) {
      const article = res.data.data
      form.title = article.title
      form.summary = article.summary || ''
      form.content = article.content || ''
      form.cover_url = article.cover_url || ''
      form.category_id = article.category_id || null
      form.tag_ids = article.tags ? article.tags.map(t => t.id) : []
      form.is_pinned = article.is_pinned
      form.status = article.status
    }
  } catch (err) {
    message.error('加载文章失败')
    router.back()
  } finally {
    pageLoading.value = false
  }
}

const handleSave = async () => {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  saving.value = true
  try {
    const data = {
      title: form.title,
      summary: form.summary,
      content: form.content,
      cover_url: form.cover_url,
      category_id: form.category_id,
      tag_ids: form.tag_ids,
      is_pinned: form.is_pinned,
      status: form.status
    }

    if (isEdit.value) {
      await updateArticle(route.params.id, data)
      message.success('文章已更新')
    } else {
      await createArticle(data)
      message.success('文章已创建')
    }
    router.push({ name: 'ArticleManage' })
  } catch (err) {
    message.error(err.response?.data?.message || '保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadOptions()
  loadArticle()
})
</script>

<style scoped>
.article-edit-page {
  max-width: 1000px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid #e8e8e8;
}
</style>
