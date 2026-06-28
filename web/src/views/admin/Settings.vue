<template>
  <div class="settings-page">
    <n-card title="系统设置">
      <n-spin :show="loading">
        <n-form
          ref="formRef"
          :model="form"
          label-placement="left"
          label-width="100"
          style="max-width: 600px;"
        >
          <n-divider title-placement="left">基本信息</n-divider>

          <n-form-item label="站点名称">
            <n-input v-model:value="form.site_name" placeholder="HiGolang" />
          </n-form-item>

          <n-form-item label="站点描述">
            <n-input
              v-model:value="form.description"
              type="textarea"
              placeholder="站点描述"
              :rows="2"
            />
          </n-form-item>

          <n-form-item label="Logo URL">
            <n-input v-model:value="form.logo_url" placeholder="站点 Logo 图片地址" />
          </n-form-item>

          <n-divider title-placement="left">SEO 设置</n-divider>

          <n-form-item label="Meta 描述">
            <n-input
              v-model:value="form.meta_description"
              type="textarea"
              placeholder="搜索引擎描述"
              :rows="2"
            />
          </n-form-item>

          <n-form-item label="关键词">
            <n-input v-model:value="form.meta_keywords" placeholder="关键词，用逗号分隔" />
          </n-form-item>

          <div class="form-actions">
            <n-button type="primary" :loading="saving" @click="handleSave">
              保存设置
            </n-button>
          </div>
        </n-form>
      </n-spin>
    </n-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useMessage } from 'naive-ui'
import { getAdminSettings, updateAdminSettings } from '../../api'

const message = useMessage()
const loading = ref(false)
const saving = ref(false)

const form = reactive({
  site_name: '',
  description: '',
  logo_url: '',
  meta_description: '',
  meta_keywords: ''
})

const loadSettings = async () => {
  loading.value = true
  try {
    const res = await getAdminSettings()
    if (res.data.code === 0) {
      const data = res.data.data
      Object.keys(form).forEach(key => {
        form[key] = data[key] || ''
      })
    }
  } catch (err) {
    console.error('Load settings failed:', err)
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  saving.value = true
  try {
    await updateAdminSettings({ settings: { ...form } })
    message.success('设置已保存')
  } catch (err) {
    message.error('保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(loadSettings)
</script>

<style scoped>
.form-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid #e8e8e8;
}
</style>
