<template>
  <div class="download-links">
    <div v-if="loading" class="dl-loading">
      <span class="spinner"></span>
    </div>
    <div v-else-if="links.length === 0" class="dl-empty">
      {{ locale === 'zh' ? '暂无下载链接' : 'No download links available' }}
    </div>
    <div v-else class="dl-list">
      <a
        v-for="link in links"
        :key="link.id"
        :href="link.url"
        target="_blank"
        rel="noopener noreferrer"
        class="dl-item"
      >
        <div class="dl-icon" :class="`dl-icon-${link.platform}`">
          {{ platformIcons[link.platform] || '📁' }}
        </div>
        <div class="dl-info">
          <span class="dl-name">{{ platformNames[link.platform] }}</span>
          <span v-if="link.extract_code" class="dl-code">
            {{ locale === 'zh' ? '提取码' : 'Code' }}: <strong>{{ link.extract_code }}</strong>
            <button class="copy-btn" @click.prevent="copyCode(link.extract_code)">
              {{ copied === link.extract_code ? '✓' : '复制' }}
            </button>
          </span>
        </div>
        <span class="dl-arrow">→</span>
      </a>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { downloadAPI } from '@/api'

const props = defineProps({
  productId: { type: String, required: true }
})

const { locale } = useI18n()
const loading = ref(true)
const links = ref([])
const copied = ref('')

const platformIcons = {
  baidu: '🔵',
  quark: '⚡',
  google: '🔴'
}

const platformNames = {
  baidu: '百度网盘',
  quark: '夸克网盘',
  google: 'Google Drive'
}

async function loadLinks() {
  loading.value = true
  try {
    const res = await downloadAPI.getLinks(props.productId)
    links.value = res.links || res || []
  } catch {
    links.value = []
  } finally {
    loading.value = false
  }
}

async function copyCode(code) {
  await navigator.clipboard.writeText(code)
  copied.value = code
  setTimeout(() => { copied.value = '' }, 2000)
}

onMounted(loadLinks)
</script>

<style scoped>
.dl-loading {
  display: flex;
  justify-content: center;
  padding: 20px;
}

.dl-empty {
  text-align: center;
  padding: 16px;
  color: var(--c-text-tertiary);
  font-size: 0.88rem;
}

.dl-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.dl-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 14px;
  background: var(--c-surface);
  border: 1.5px solid var(--c-border-light);
  border-radius: var(--r-md);
  transition: all var(--t-fast);
  color: var(--c-text-primary);
}
.dl-item:hover {
  border-color: var(--c-accent);
  background: var(--c-accent-light);
  transform: translateX(3px);
}

.dl-icon {
  width: 36px;
  height: 36px;
  border-radius: var(--r-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.2rem;
  flex-shrink: 0;
}

.dl-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.dl-name {
  font-size: 0.88rem;
  font-weight: 600;
  color: var(--c-text-primary);
}

.dl-code {
  font-size: 0.78rem;
  color: var(--c-text-secondary);
  display: flex;
  align-items: center;
  gap: 6px;
}

.copy-btn {
  background: var(--c-accent-light);
  color: var(--c-accent);
  border: none;
  padding: 2px 8px;
  border-radius: var(--r-full);
  font-size: 0.72rem;
  font-weight: 600;
  cursor: pointer;
  transition: all var(--t-fast);
  font-family: var(--font-sans);
}
.copy-btn:hover { background: var(--c-accent); color: white; }

.dl-arrow {
  font-size: 0.9rem;
  color: var(--c-text-tertiary);
  transition: transform var(--t-fast);
}
.dl-item:hover .dl-arrow { transform: translateX(2px); }
</style>
