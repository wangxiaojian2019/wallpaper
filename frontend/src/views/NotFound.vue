<template>
  <div class="not-found-page">
    <div class="not-found-content">
      <div class="error-code">404</div>
      <h2>{{ isZh ? '页面未找到' : 'Page Not Found' }}</h2>
      <p class="text-secondary">
        {{ isZh ? '抱歉，您访问的页面不存在或已被移除。' : 'Sorry, the page you are looking for does not exist or has been removed.' }}
      </p>
      <router-link to="/" class="btn btn-primary btn-lg">
        {{ isZh ? '返回首页' : 'Back to Home' }}
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { locale } = useI18n()
const isZh = computed(() => locale.value === 'zh')
</script>

<style scoped>
.not-found-page {
  min-height: 75vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  text-align: center;
}

.error-code {
  font-size: 7rem;
  font-weight: 900;
  line-height: 1;
  background: linear-gradient(135deg, var(--primary), #a855f7);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin-bottom: 16px;
  letter-spacing: -2px;
}

.not-found-content h2 {
  font-size: 2rem;
  font-weight: 700;
  margin-bottom: 12px;
}

.not-found-content p {
  font-size: 1.1rem;
  margin-bottom: 32px;
  max-width: 460px;
}
</style>
