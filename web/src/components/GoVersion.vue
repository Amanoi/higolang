<template>
  <n-tag v-if="version" :bordered="false" type="success" size="small" class="go-version">
    <template #icon>
      <n-icon><code-icon /></n-icon>
    </template>
    Go {{ version }}
  </n-tag>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { CodeOutline as CodeIcon } from '@vicons/ionicons5'
import { getGoVersion } from '../api'

const version = ref('')

onMounted(async () => {
  try {
    const response = await getGoVersion()
    if (response.data.code === 0) {
      version.value = response.data.data.version
    }
  } catch (error) {
    console.error('Failed to load Go version:', error)
  }
})
</script>

<style scoped>
.go-version {
  font-weight: 500;
}
</style>
