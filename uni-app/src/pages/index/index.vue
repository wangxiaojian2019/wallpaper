<template>
  <view class="home-container" :class="[appStore.theme]">
    <!-- Dynamic Ambient Background (Starry Sky & Meteors in Dark, Glass Raindrops in Light) -->
    <AmbientCanvas />

    <!-- Top Sticky Frosted Glass Header with Theme & Lang Switches -->
    <AppHeader />

    <!-- Main Container with Moderate PC Margins -->
    <view class="content-container">
      <!-- Spacious Ambient Hero -->
      <view class="hero-spotlight">
        <view class="spotlight-badge">
          <text class="badge-star">✦</text>
          <text class="badge-title">{{ appStore.lang === 'zh' ? '4K 极致视觉画廊' : 'Curated 4K Gallery' }}</text>
        </view>
        <text class="hero-headline">
          {{ appStore.lang === 'zh' ? '让每一块屏幕，都有惊艳时刻' : 'Transform Every Display with Pure Art' }}
        </text>

        <!-- Search bar -->
        <view class="search-box">
          <text class="search-icon">🔍</text>
          <input
            v-model="searchKeyword"
            :placeholder="appStore.lang === 'zh' ? '搜索壁纸、风格、色系...' : 'Search wallpapers, tags, styles...'"
            class="search-input"
            @confirm="onSearch"
          />
          <text v-if="searchKeyword" class="clear-icon" @click="clearSearch" @tap="clearSearch">✕</text>
        </view>
      </view>

      <!-- Device Filter Tabs -->
      <view class="filter-wrapper">
        <scroll-view scroll-x class="device-scroll" :show-scrollbar="false">
          <view class="device-row">
            <view
              v-for="d in deviceOptions"
              :key="d.value"
              class="device-chip"
              :class="{ active: currentDevice === d.value }"
              @click="selectDevice(d.value)"
              @tap="selectDevice(d.value)"
            >
              <text class="chip-icon">{{ d.icon }}</text>
              <text class="chip-text">{{ appStore.lang === 'zh' ? d.zh : d.en }}</text>
            </view>
          </view>
        </scroll-view>

        <!-- Category Filter Pills -->
        <scroll-view scroll-x class="category-scroll" :show-scrollbar="false">
          <view class="category-row">
            <view
              v-for="cat in categoryList"
              :key="cat.value"
              class="cat-chip"
              :class="{ active: currentCategory === cat.value }"
              @click="selectCategory(cat.value)"
              @tap="selectCategory(cat.value)"
            >
              <text class="cat-chip-text">{{ appStore.lang === 'zh' ? cat.zh : cat.en }}</text>
            </view>
          </view>
        </scroll-view>
      </view>

      <!-- Wallpaper Gallery Grid (Auto-fill Responsive: 4-5 cols on PC, 2 on mobile) -->
      <view class="gallery-wrapper">
        <!-- Loading skeleton -->
        <view v-if="loading && products.length === 0" class="wallpaper-grid">
          <view v-for="n in 8" :key="n" class="skeleton-card"></view>
        </view>

        <!-- Error state with retry -->
        <view v-else-if="loadError && products.length === 0" class="empty-state">
          <text class="empty-icon">⚠️</text>
          <text class="empty-text">
            {{ errorMessage || (appStore.lang === 'zh' ? '网络连接异常，未能加载壁纸' : 'Failed to connect. Please retry.') }}
          </text>
          <button class="btn-retry" @tap="loadData(true)">
            {{ appStore.lang === 'zh' ? '点击重新加载' : 'Retry' }}
          </button>
        </view>

        <!-- Empty state -->
        <view v-else-if="products.length === 0" class="empty-state">
          <text class="empty-icon">🖼️</text>
          <text class="empty-text">
            {{ appStore.lang === 'zh' ? '未找到相关壁纸，换个风格或关键词探索' : 'No wallpapers found. Try another query.' }}
          </text>
        </view>

        <!-- Products Grid -->
        <view v-else class="wallpaper-grid">
          <view
            v-for="item in products"
            :key="item.id"
            class="grid-item"
          >
            <WallpaperCard :product="item" />
          </view>
        </view>

        <!-- Load More / End state -->
        <view class="load-more-box">
          <text v-if="loading" class="load-text">
            {{ appStore.lang === 'zh' ? '正在探索更多画作...' : 'Loading more artworks...' }}
          </text>
          <text v-else-if="!hasMore && products.length > 0" class="load-text">
            {{ appStore.lang === 'zh' ? '✦ 已呈现全部精选壁纸 ✦' : '✦ End of Curated Collection ✦' }}
          </text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { onPullDownRefresh, onReachBottom } from '@dcloudio/uni-app'
import { productAPI, categoryAPI } from '@/api'
import { useAppStore } from '@/store/app'
import AppHeader from '@/components/AppHeader.vue'
import WallpaperCard from '@/components/WallpaperCard.vue'
import AmbientCanvas from '@/components/AmbientCanvas.vue'

const appStore = useAppStore()
const products = ref([])
const loading = ref(false)
const hasMore = ref(true)
const page = ref(1)
const loadError = ref(false)
const errorMessage = ref('')
const searchKeyword = ref('')
const currentDevice = ref('all')
const currentCategory = ref('all')

const deviceOptions = [
  { value: 'all', zh: '全部设备', en: 'All Devices', icon: '🌟' },
  { value: 'phone', zh: '手机锁屏', en: 'iPhone / Phone', icon: '📱' },
  { value: 'desktop', zh: '电脑桌面', en: 'Mac / Desktop', icon: '💻' },
  { value: 'live', zh: '动态 4K', en: 'Live 4K Motion', icon: '🎬' }
]

const categoryList = ref([
  { value: 'all', zh: '🔥 全部', en: '🔥 All' },
  { value: 'ios', zh: '📱 iOS 风格', en: '📱 iOS Theme' },
  { value: 'anime', zh: '✨ 动漫二次元', en: '✨ Anime' },
  { value: 'cyberpunk', zh: '🏙️ 赛博科技', en: '🏙️ Cyberpunk' },
  { value: 'nature', zh: '🌿 治愈风景', en: '🌿 Nature' },
  { value: 'minimal', zh: '🌸 极简冷淡', en: '🌸 Minimal' },
  { value: 'free', zh: '🎁 限时免费', en: '🎁 Free Access' }
])

async function loadCategories() {
  try {
    const res = await categoryAPI.list()
    const items = res.items || res.categories || res || []
    if (items.length > 0) {
      const dynamicCats = items.map(c => ({
        value: c.code,
        zh: `${c.icon || '📁'} ${c.name_zh}`.trim(),
        en: `${c.icon || '📁'} ${c.name_en || c.name_zh}`.trim()
      }))
      categoryList.value = [
        { value: 'all', zh: '🔥 全部', en: '🔥 All' },
        ...dynamicCats,
        { value: 'free', zh: '🎁 限时免费', en: '🎁 Free Access' }
      ]
    }
  } catch (e) {
    console.warn('Failed to fetch categories:', e)
  }
}

async function loadData(reset = false) {
  if (loading.value) return
  if (reset) {
    page.value = 1
    hasMore.value = true
    loadError.value = false
    errorMessage.value = ''
  }

  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: 20,
      search: searchKeyword.value ? searchKeyword.value.trim() : undefined,
      category: currentCategory.value === 'all' ? undefined : currentCategory.value
    }
    const res = await productAPI.list(params)
    const items = res.items || res.products || res.data?.products || []

    if (reset) {
      products.value = items
    } else {
      products.value.push(...items)
    }

    const total = res.total || 0
    hasMore.value = products.value.length < total
  } catch (err) {
    console.error('Failed to load products:', err)
    if (products.value.length === 0) {
      loadError.value = true
      errorMessage.value = err?.message || ''
    }
  } finally {
    loading.value = false
    uni.stopPullDownRefresh()
  }
}

function onSearch() {
  loadData(true)
}

function clearSearch() {
  searchKeyword.value = ''
  loadData(true)
}

function selectDevice(val) {
  currentDevice.value = val
  loadData(true)
}

function selectCategory(val) {
  currentCategory.value = val
  loadData(true)
}

onMounted(() => {
  loadCategories()
  loadData(true)
})

onPullDownRefresh(() => {
  loadData(true)
})

onReachBottom(() => {
  if (hasMore.value && !loading.value) {
    page.value++
    loadData(false)
  }
})
</script>

<style scoped>
.home-container {
  min-height: 100vh;
  position: relative;
  transition: background-color 0.25s ease, color 0.25s ease;
  padding-bottom: 60px;
}

.home-container.dark {
  background: transparent;
  color: #f1f2f6;
}

.home-container.light {
  background: transparent;
  color: #1a1a1e;
}

/* Page Container (Moderate Balanced PC Margins) */
.content-container {
  position: relative;
  z-index: 10;
  width: 100%;
  max-width: 1240px;
  margin: 0 auto;
  padding: 0 24px;
  box-sizing: border-box;
}

@media (max-width: 768px) {
  .content-container {
    padding: 0 12px;
  }
}

/* Hero Spotlight */
.hero-spotlight {
  padding: 48px 0 32px;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.spotlight-badge {
  display: inline-flex;
  align-items: center;
  padding: 6px 16px;
  border-radius: 999px;
  margin-bottom: 16px;
  backdrop-filter: blur(12px);
}

.dark .spotlight-badge {
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.12);
}

.light .spotlight-badge {
  background: rgba(0, 0, 0, 0.05);
  border: 1px solid rgba(0, 0, 0, 0.08);
}

.badge-star {
  font-size: 14px;
  color: #fbbf24;
  margin-right: 6px;
}

.badge-title {
  font-size: 13px;
  font-weight: 700;
  letter-spacing: 0.5px;
}

.dark .badge-title {
  color: rgba(255, 255, 255, 0.85);
}

.light .badge-title {
  color: rgba(0, 0, 0, 0.75);
}

.hero-headline {
  font-size: 32px;
  font-weight: 800;
  letter-spacing: -0.5px;
  line-height: 1.35;
  margin-bottom: 24px;
}

@media (max-width: 768px) {
  .hero-spotlight {
    padding: 20px 0 16px;
  }
  .spotlight-badge {
    padding: 4px 12px;
    margin-bottom: 10px;
  }
  .badge-star {
    font-size: 12px;
  }
  .badge-title {
    font-size: 11px;
  }
  .hero-headline {
    font-size: 18px !important;
    line-height: 1.35 !important;
    margin-bottom: 14px !important;
    letter-spacing: -0.2px !important;
  }
  .search-box {
    height: 40px !important;
    padding: 0 14px !important;
    margin-bottom: 8px !important;
  }
  .search-icon {
    font-size: 14px !important;
    margin-right: 8px !important;
  }
  .search-input {
    font-size: 13px !important;
  }
}

.dark .hero-headline {
  color: #ffffff;
}

.light .hero-headline {
  color: #111827;
}

/* Search bar */
.search-box {
  width: 100%;
  max-width: 520px;
  height: 48px;
  border-radius: 999px;
  display: flex;
  align-items: center;
  padding: 0 20px;
  box-sizing: border-box;
  transition: all 0.2s ease;
}

.dark .search-box {
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.14);
}

.light .search-box {
  background: #ffffff;
  border: 1px solid rgba(0, 0, 0, 0.12);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.05);
}

.search-icon {
  font-size: 16px;
  margin-right: 12px;
  opacity: 0.55;
}

.search-input {
  flex: 1;
  font-size: 15px;
}

.dark .search-input {
  color: #ffffff;
}

.light .search-input {
  color: #111827;
}

.clear-icon {
  font-size: 14px;
  opacity: 0.45;
  padding: 6px;
  cursor: pointer;
}

/* Filters */
.filter-wrapper {
  margin-bottom: 8px;
}

.device-scroll, .category-scroll {
  white-space: nowrap;
  width: 100%;
}

.device-row {
  display: inline-flex;
  gap: 12px;
  padding: 8px 0 14px;
}

.device-chip {
  display: inline-flex;
  align-items: center;
  padding: 8px 18px;
  border-radius: 999px;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.dark .device-chip {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.light .device-chip {
  background: #ffffff;
  border: 1px solid rgba(0, 0, 0, 0.08);
}

.dark .device-chip.active {
  background: #ffffff;
  border-color: #ffffff;
  box-shadow: 0 4px 14px rgba(255, 255, 255, 0.2);
}

.light .device-chip.active {
  background: #111827;
  border-color: #111827;
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.15);
}

.chip-icon {
  font-size: 14px;
  margin-right: 6px;
}

.chip-text {
  font-size: 13px;
  font-weight: 600;
}

.dark .chip-text {
  color: rgba(255, 255, 255, 0.7);
}

.light .chip-text {
  color: rgba(0, 0, 0, 0.65);
}

.dark .device-chip.active .chip-text {
  color: #111218;
}

.light .device-chip.active .chip-text {
  color: #ffffff;
}

/* Category row */
.category-row {
  display: inline-flex;
  gap: 10px;
  padding: 4px 0 16px;
}

.cat-chip {
  padding: 6px 14px;
  border-radius: 999px;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.dark .cat-chip {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid transparent;
}

.light .cat-chip {
  background: rgba(0, 0, 0, 0.04);
  border: 1px solid transparent;
}

.cat-chip.active {
  background: rgba(59, 130, 246, 0.15);
  border-color: rgba(59, 130, 246, 0.35);
}

.cat-chip-text {
  font-size: 13px;
  font-weight: 500;
}

.dark .cat-chip-text {
  color: rgba(255, 255, 255, 0.6);
}

.light .cat-chip-text {
  color: rgba(0, 0, 0, 0.6);
}

.cat-chip.active .cat-chip-text {
  color: #3b82f6;
  font-weight: 700;
}

/* ========================================================
   Wallpaper Gallery Grid (Legacy Gold Ratio: minmax(215px, 1fr))
   ======================================================== */
.gallery-wrapper {
  width: 100%;
}

.wallpaper-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(215px, 1fr));
  gap: 22px;
  width: 100%;
  padding: 16px 0 48px;
  box-sizing: border-box;
}

@media (max-width: 768px) {
  .device-row {
    gap: 8px !important;
    padding: 2px 0 8px !important;
  }
  .device-chip {
    padding: 6px 12px !important;
  }
  .chip-icon {
    font-size: 12px !important;
    margin-right: 4px !important;
  }
  .chip-text {
    font-size: 12px !important;
  }
  .category-row {
    gap: 6px !important;
    padding: 2px 0 10px !important;
  }
  .cat-chip {
    padding: 4px 10px !important;
  }
  .cat-chip-text {
    font-size: 12px !important;
  }
  .wallpaper-grid {
    grid-template-columns: repeat(2, 1fr) !important;
    gap: 10px !important;
    padding: 8px 0 28px !important;
  }
  .skeleton-card {
    border-radius: 16px !important;
  }
}

.grid-item {
  width: 100%;
}

.skeleton-card {
  width: 100%;
  aspect-ratio: 9 / 16;
  border-radius: 20px;
  animation: pulse 1.6s ease-in-out infinite;
}

.dark .skeleton-card {
  background: #191a24;
}

.light .skeleton-card {
  background: #e2e8f0;
}

@keyframes pulse {
  0%, 100% { opacity: 0.6; }
  50% { opacity: 0.25; }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 80px 0;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
  opacity: 0.4;
}

.empty-text {
  font-size: 15px;
  opacity: 0.5;
}

.btn-retry {
  margin-top: 18px;
  padding: 8px 24px;
  font-size: 13px;
  font-weight: 700;
  border-radius: 999px;
  background: linear-gradient(135deg, #3b82f6, #2563eb);
  color: #ffffff;
  border: none;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.3);
}

.load-more-box {
  display: flex;
  justify-content: center;
  padding: 40px 0 20px;
}

.load-text {
  font-size: 13px;
  opacity: 0.4;
  letter-spacing: 1px;
}
</style>
