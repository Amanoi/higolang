<template>
  <div class="about-page">
    <div class="about-hero">
      <h1>关于 HiGolang</h1>
      <p class="about-desc">
        HiGolang 是一个专注于 Go 语言最新动态的资讯聚合站点。我们自动收集来自 Go 官方博客、GitHub 发布、社区讨论等来源的最新内容，帮助 Go 开发者第一时间掌握语言生态变化。
      </p>
    </div>

    <n-card title="站点功能" class="section-card">
      <p>HiGolang 提供以下功能：</p>
      <ul class="feature-list">
        <li>自动聚合 Go 官方博客、GitHub Releases、Reddit 等来源的资讯</li>
        <li>按分类和标签浏览文章，快速定位感兴趣的内容</li>
        <li>全文搜索，精确查找所需信息</li>
        <li>管理员后台，支持手动发布原创文章</li>
        <li>可配置的资讯源管理，自定义抓取频率</li>
      </ul>
    </n-card>

    <n-card title="Go 语言相关链接" class="section-card">
      <n-list bordered>
        <n-list-item v-for="link in usefulLinks" :key="link.name">
          <template #prefix>
            <n-icon :size="20" color="#00ADD8"><link-icon /></n-icon>
          </template>
          <template #suffix>
            <n-button text tag="a" :href="link.url" target="_blank" type="primary">
              访问
            </n-button>
          </template>
          <n-thing :title="link.name" :description="link.desc" />
        </n-list-item>
      </n-list>
    </n-card>

    <n-card title="技术栈" class="section-card">
      <n-space>
        <n-tag type="info" size="large">Go 1.21</n-tag>
        <n-tag type="info" size="large">Gin</n-tag>
        <n-tag type="info" size="large">GORM</n-tag>
        <n-tag type="info" size="large">SQLite</n-tag>
        <n-tag type="success" size="large">Vue 3</n-tag>
        <n-tag type="success" size="large">Vite</n-tag>
        <n-tag type="success" size="large">Naive UI</n-tag>
        <n-tag type="warning" size="large">Docker</n-tag>
      </n-space>
    </n-card>
  </div>
</template>

<script setup>
import { Link as LinkIcon } from '@vicons/ionicons5'

const usefulLinks = [
  { name: 'Go 官方网站', url: 'https://go.dev', desc: 'Go 语言官网，文档、教程、博客' },
  { name: 'Go GitHub', url: 'https://github.com/golang/go', desc: 'Go 语言源代码仓库' },
  { name: 'Go 文档', url: 'https://go.dev/doc/', desc: '官方文档、Effective Go、代码规范' },
  { name: 'Go 博客', url: 'https://go.dev/blog/', desc: 'Go 团队的技术博客' },
  { name: 'Go 周报', url: 'https://golangweekly.com/', desc: '每周 Go 语言精选资讯' },
  { name: 'Reddit r/golang', url: 'https://www.reddit.com/r/golang/', desc: 'Go 语言 Reddit 社区' },
  { name: 'Go Playground', url: 'https://go.dev/play/', desc: '在线 Go 代码运行环境' },
  { name: 'pkg.go.dev', url: 'https://pkg.go.dev/', desc: 'Go 包搜索和文档平台' },
]
</script>

<style scoped>
.about-page {
  padding: 24px 0;
  max-width: 800px;
  margin: 0 auto;
}

.about-hero {
  text-align: center;
  padding: 48px 24px;
  background: linear-gradient(135deg, #00ADD8 0%, #00C5F5 100%);
  border-radius: 12px;
  margin-bottom: 32px;
  color: #fff;
}

.about-hero h1 {
  font-size: 36px;
  margin-bottom: 16px;
}

.about-desc {
  font-size: 16px;
  line-height: 1.8;
  opacity: 0.95;
  max-width: 600px;
  margin: 0 auto;
}

.section-card {
  margin-bottom: 24px;
}

.feature-list {
  padding-left: 20px;
  line-height: 2;
  color: #555;
}

@media (max-width: 768px) {
  .about-hero h1 {
    font-size: 28px;
  }
}
</style>
