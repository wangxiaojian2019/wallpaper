<template>
  <div class="home-page">
    <!-- Hero Section -->
    <section class="hero">
      <div class="container">
        <div class="hero-content">
          <div class="hero-badge">
            <span>✦</span> {{ locale === 'zh' ? '精选壁纸商城' : 'Premium Wallpaper Store' }}
          </div>
          <h1 class="hero-title">{{ t('home.hero_title') }}</h1>
          <p class="hero-subtitle">{{ t('home.hero_subtitle') }}</p>
          <div class="hero-cta">
            <a href="#gallery" class="btn btn-dark btn-lg">
              {{ t('home.explore') }} →
            </a>
            <router-link to="/membership" class="btn btn-outline btn-lg">
              👑 {{ t('nav.membership') }}
            </router-link>
          </div>

          <!-- Stats -->
          <div class="hero-stats">
            <div class="stat">
              <span class="stat-num">{{ stats.total }}+</span>
              <span class="stat-label">{{ locale === 'zh' ? '壁纸' : 'Wallpapers' }}</span>
            </div>
            <div class="stat-divider"></div>
            <div class="stat">
              <span class="stat-num">{{ stats.users }}+</span>
              <span class="stat-label">{{ locale === 'zh' ? '用户' : 'Users' }}</span>
            </div>
            <div class="stat-divider"></div>
            <div class="stat">
              <span class="stat-num">4K</span>
              <span class="stat-label">{{ locale === 'zh' ? '高清' : 'Resolution' }}</span>
            </div>
          </div>
        </div>

        <!-- Hero Preview Grid -->
        <div class="hero-preview">
          <div v-for="(item, i) in previewItems" :key="i" class="preview-card" :class="`preview-card-${i}`">
            <img v-if="item.cover" :src="item.cover" :alt="item.title_zh" loading="lazy" />
            <div v-else class="preview-placeholder" :style="{ background: item.color }">
              <span>{{ item.emoji }}</span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Filter Bar -->
    <section id="gallery" class="gallery-section">
      <div class="container">
        <div class="filter-bar">
          <div class="filter-tags">
            <button
              v-for="cat in categories"
              :key="cat.value"
              class="tag"
              :class="{ active: activeCategory === cat.value }"
              @click="setCategory(cat.value)"
            >
              {{ cat.icon }} {{ locale === 'zh' ? cat.zh : cat.en }}
            </button>
          </div>
          <div class="filter-right">
            <div class="search-box">
              <span class="search-icon">🔍</span>
              <input
                v-model="searchQuery"
                type="text"
                class="search-input"
                :placeholder="t('nav.search')"
                @input="handleSearch"
              />
            </div>
            <select v-model="sortBy" class="form-input form-select sort-select" @change="loadProducts">
              <option value="newest">{{ t('home.newest') }}</option>
              <option value="popular">{{ t('home.popular') }}</option>
              <option value="price_asc">{{ locale === 'zh' ? '价格↑' : 'Price↑' }}</option>
              <option value="price_desc">{{ locale === 'zh' ? '价格↓' : 'Price↓' }}</option>
            </select>
          </div>
        </div>

        <!-- Grid -->
        <div v-if="loading && products.length === 0" class="wallpaper-grid">
          <div v-for="n in 8" :key="n" class="wallpaper-grid-item">
            <div class="skeleton" :style="{ height: `${200 + Math.random() * 150}px`, borderRadius: '16px' }"></div>
          </div>
        </div>

        <div v-else-if="products.length === 0" class="empty-state">
          <span class="empty-icon">🖼️</span>
          <p>{{ t('common.no_data') }}</p>
        </div>

        <div v-else class="wallpaper-grid">
          <div v-for="product in products" :key="product.id" class="wallpaper-grid-item">
            <WallpaperCard :product="product" @buy="openBuyModal" />
          </div>
        </div>

        <!-- Load More -->
        <div v-if="hasMore" class="load-more">
          <button class="btn btn-outline" :disabled="loading" @click="loadMore">
            <span v-if="loading" class="spinner"></span>
            {{ loading ? t('home.loading') : (locale === 'zh' ? '加载更多' : 'Load More') }}
          </button>
        </div>

        <div v-if="!hasMore && products.length > 0" class="no-more">
          {{ t('home.no_more') }} ✦
        </div>
      </div>
    </section>

    <!-- Membership Banner -->
    <section class="member-banner">
      <div class="container">
        <div class="member-banner-inner">
          <div class="member-banner-content">
            <span class="member-banner-icon">👑</span>
            <div>
              <h3>{{ locale === 'zh' ? '成为会员，解锁全部壁纸' : 'Become a Member, Unlock Everything' }}</h3>
              <p>{{ locale === 'zh' ? '月度/季度/年度会员，无限下载所有壁纸' : 'Monthly/Quarterly/Yearly plans with unlimited downloads' }}</p>
            </div>
          </div>
          <router-link to="/membership" class="btn btn-white btn-lg">
            {{ locale === 'zh' ? '查看方案' : 'View Plans' }} →
          </router-link>
        </div>
      </div>
    </section>

    <!-- Buy Modal -->
    <BuyModal
      v-if="showBuyModal"
      :product="selectedProduct"
      :initial-mode="buyModalMode"
      @close="showBuyModal = false"
      @success="handlePurchaseSuccess"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useDebounceFn } from '@vueuse/core'
import { productAPI } from '@/api'
import { useAuthStore } from '@/store'
import { useToast } from '@/composables/useToast'
import WallpaperCard from '@/components/WallpaperCard.vue'
import BuyModal from '@/components/BuyModal.vue'

const { t, locale } = useI18n()
const auth = useAuthStore()
const { success, error } = useToast()

const products = ref([])
const loading = ref(false)
const hasMore = ref(true)
const page = ref(1)
const searchQuery = ref('')
const activeCategory = ref('all')
const sortBy = ref('newest')
const showBuyModal = ref(false)
const selectedProduct = ref(null)

const stats = ref({ total: 200, users: 5000 })

const previewItems = ref([
  { color: 'linear-gradient(135deg,#667eea,#764ba2)', emoji: '🌌', cover: null },
  { color: 'linear-gradient(135deg,#f093fb,#f5576c)', emoji: '🌸', cover: null },
  { color: 'linear-gradient(135deg,#4facfe,#00f2fe)', emoji: '🌊', cover: null },
  { color: 'linear-gradient(135deg,#43e97b,#38f9d7)', emoji: '🌿', cover: null },
  { color: 'linear-gradient(135deg,#fa709a,#fee140)', emoji: '🌅', cover: null },
  { color: 'linear-gradient(135deg,#a18cd1,#fbc2eb)', emoji: '✨', cover: null }
])

const categories = [
  { value: 'all', zh: '全部壁纸', en: 'All', icon: '🌟' },
  { value: '苹果壁纸', zh: '苹果壁纸', en: 'iOS Themes', icon: '📱' },
  { value: '动态壁纸', zh: '动态壁纸', en: 'Live 4K', icon: '🎬' },
  { value: '电脑壁纸', zh: '电脑壁纸', en: 'Desktop', icon: '💻' },
  { value: '动漫', zh: '动漫卡通', en: 'Anime', icon: '✨' },
  { value: '赛博朋克', zh: '赛博朋克', en: 'Cyberpunk', icon: '🏙️' },
  { value: '教程分享', zh: '教程工具', en: 'Tutorials', icon: '🛠️' }
]

async function loadProducts(reset = false) {
  if (loading.value) return
  if (reset) {
    products.value = []
    page.value = 1
    hasMore.value = true
  }
  loading.value = true
  try {
    const res = await productAPI.list({
      page: page.value,
      page_size: 24,
      category: activeCategory.value === 'all' ? '' : activeCategory.value,
      search: searchQuery.value,
      sort: sortBy.value
    })
    const newItems = res.items || res.products || res.data?.products || []
    if (reset) {
      products.value = newItems
    } else {
      products.value.push(...newItems)
    }
    const total = res.total || 0
    hasMore.value = products.value.length < total
    // Also update preview with real images
    if (products.value.length > 0) {
      products.value.slice(0, 6).forEach((p, i) => {
        if (p.cover_image && previewItems.value[i]) {
          previewItems.value[i].cover = p.cover_image
        }
      })
      stats.value.total = total
    }
  } catch (err) {
    console.error('Error loading products:', err)
  } finally {
    loading.value = false
  }
}

function loadMockData(reset) {
  const mockCategories = ['nature', 'space', 'city', 'abstract', 'anime', 'minimal']
  const mockData = Array.from({ length: 12 }, (_, i) => ({
    id: `mock-${page.value}-${i}`,
    title_zh: `精美壁纸 ${page.value * 12 + i + 1}`,
    title_en: `Wallpaper ${page.value * 12 + i + 1}`,
    price: (4.99 + i * 0.5).toFixed(2),
    cover_image: null,
    category: mockCategories[i % mockCategories.length],
    is_published: true,
    _mock: true,
    _color: [
      'linear-gradient(135deg,#667eea,#764ba2)',
      'linear-gradient(135deg,#f093fb,#f5576c)',
      'linear-gradient(135deg,#4facfe,#00f2fe)',
      'linear-gradient(135deg,#43e97b,#38f9d7)',
      'linear-gradient(135deg,#fa709a,#fee140)',
      'linear-gradient(135deg,#a18cd1,#fbc2eb)',
      'linear-gradient(135deg,#ffecd2,#fcb69f)',
      'linear-gradient(135deg,#2af598,#009efd)',
      'linear-gradient(135deg,#feada6,#f5efef)',
      'linear-gradient(135deg,#c471f5,#fa71cd)',
      'linear-gradient(135deg,#30cfd0,#330867)',
      'linear-gradient(135deg,#89f7fe,#66a6ff)'
    ][i % 12]
  }))
  if (reset) {
    products.value = mockData
  } else {
    products.value.push(...mockData)
  }
  hasMore.value = page.value < 3
}

function setCategory(cat) {
  activeCategory.value = cat
  loadProducts(true)
}

const handleSearch = useDebounceFn(() => {
  loadProducts(true)
}, 400)

function loadMore() {
  page.value++
  loadProducts()
}

const buyModalMode = ref('guest')

function openBuyModal(payload) {
  if (payload?.product) {
    selectedProduct.value = payload.product
    buyModalMode.value = payload.mode || 'guest'
  } else {
    selectedProduct.value = payload
    buyModalMode.value = 'guest'
  }
  showBuyModal.value = true
}

function handlePurchaseSuccess() {
  showBuyModal.value = false
  success(t('payment.success'))
}

onMounted(() => {
  loadProducts(true)
})
</script>

<style scoped>
/* Hero */
.hero {
  background: var(--g-hero);
  padding: 80px 0 60px;
  overflow: hidden;
}

.hero .container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 60px;
  align-items: center;
}

.hero-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: var(--c-accent-light);
  color: var(--c-accent);
  padding: 6px 16px;
  border-radius: var(--r-full);
  font-size: 0.82rem;
  font-weight: 600;
  margin-bottom: 20px;
}

.hero-title {
  margin-bottom: 16px;
  background: linear-gradient(135deg, #1d1d1f 0%, #3a3a3c 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.hero-subtitle {
  font-size: 1.1rem;
  color: var(--c-text-secondary);
  margin-bottom: 32px;
  max-width: 440px;
}

.hero-cta {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  margin-bottom: 40px;
}

.hero-stats {
  display: flex;
  align-items: center;
  gap: 24px;
}

.stat {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.stat-num {
  font-size: 1.5rem;
  font-weight: 800;
  letter-spacing: -0.03em;
  color: var(--c-text-primary);
}

.stat-label {
  font-size: 0.78rem;
  color: var(--c-text-tertiary);
  font-weight: 500;
}

.stat-divider {
  width: 1px;
  height: 32px;
  background: var(--c-border);
}

/* Hero Preview */
.hero-preview {
  position: relative;
  height: 400px;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  grid-template-rows: repeat(2, 1fr);
  gap: 12px;
  border-radius: var(--r-2xl);
  overflow: hidden;
}

.preview-card {
  border-radius: var(--r-lg);
  overflow: hidden;
  transition: transform var(--t-normal);
}
.preview-card:hover { transform: scale(1.02); z-index: 2; }
.preview-card-0 { grid-column: span 2; }
.preview-card-3 { grid-column: span 2; }

.preview-card img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.preview-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2rem;
}

/* Gallery */
.gallery-section {
  padding: 60px 0 80px;
}

.filter-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 32px;
  flex-wrap: wrap;
}

.filter-tags {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.filter-right {
  display: flex;
  gap: 10px;
  align-items: center;
}

.search-box {
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 12px;
  font-size: 0.85rem;
  pointer-events: none;
  z-index: 1;
}

.search-input {
  padding: 10px 16px 10px 36px;
  border: 1.5px solid var(--c-border);
  border-radius: var(--r-full);
  font-size: 0.88rem;
  background: var(--c-white);
  color: var(--c-text-primary);
  width: 200px;
  transition: all var(--t-fast);
  font-family: var(--font-sans);
}
.search-input:focus {
  border-color: var(--c-accent);
  box-shadow: 0 0 0 3px rgba(0,113,227,0.1);
  width: 240px;
}

.sort-select {
  padding: 10px 36px 10px 14px;
  border-radius: var(--r-full);
  font-size: 0.88rem;
  width: auto;
  cursor: pointer;
}

.empty-state {
  text-align: center;
  padding: 80px 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}
.empty-icon { font-size: 3rem; }
.empty-state p { font-size: 1rem; }

.load-more, .no-more {
  text-align: center;
  margin-top: 48px;
}
.no-more {
  font-size: 0.88rem;
  color: var(--c-text-tertiary);
}

/* Member Banner */
.member-banner {
  background: var(--c-text-primary);
  padding: 48px 0;
}

.member-banner-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
}

.member-banner-content {
  display: flex;
  align-items: center;
  gap: 20px;
}

.member-banner-icon { font-size: 2.5rem; }

.member-banner-content h3 {
  color: white;
  font-size: 1.3rem;
  margin-bottom: 4px;
}

.member-banner-content p {
  color: rgba(255,255,255,0.6);
  font-size: 0.9rem;
}

.btn-white {
  background: white;
  color: var(--c-text-primary);
  font-weight: 700;
}
.btn-white:hover {
  background: var(--c-bg);
  transform: translateY(-1px);
}

@media (max-width: 900px) {
  .hero .container {
    grid-template-columns: 1fr;
  }
  .hero-preview {
    display: none;
  }
  .member-banner-inner {
    flex-direction: column;
    text-align: center;
  }
}

@media (max-width: 640px) {
  .hero {
    padding: 36px 0 24px;
  }
  .hero-title {
    font-size: 1.85rem !important;
    line-height: 1.2 !important;
    margin-bottom: 12px;
  }
  .hero-subtitle {
    font-size: 0.92rem;
    margin-bottom: 20px;
  }
  .hero-cta {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 8px;
    width: 100%;
    margin-bottom: 24px;
  }
  .hero-cta .btn {
    padding: 10px 14px;
    font-size: 0.88rem;
    width: 100%;
  }
  .hero-stats {
    width: 100%;
    justify-content: space-around;
    gap: 8px;
  }
  .stat-num {
    font-size: 1.25rem;
  }
  .stat-label {
    font-size: 0.72rem;
  }

  .gallery-section {
    padding: 24px 0 50px;
  }

  .filter-bar {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
    margin-bottom: 16px;
  }

  .filter-tags {
    display: flex;
    overflow-x: auto;
    white-space: nowrap;
    flex-wrap: nowrap;
    -webkit-overflow-scrolling: touch;
    scrollbar-width: none;
    padding-bottom: 4px;
    margin: 0 -12px;
    padding-left: 12px;
    padding-right: 12px;
  }

  .filter-tags::-webkit-scrollbar {
    display: none;
  }

  .filter-tags .tag {
    flex-shrink: 0;
    padding: 7px 15px !important;
    font-size: 0.82rem !important;
    border-radius: 999px !important;
  }

  .filter-right {
    display: flex;
    width: 100%;
    gap: 8px;
  }

  .search-box {
    flex: 1;
    min-width: 0;
  }

  .search-input {
    width: 100% !important;
    font-size: 16px !important;
  }

  .sort-select {
    font-size: 16px !important;
    flex-shrink: 0;
  }

  .member-banner {
    padding: 28px 0 !important;
    margin: 10px 0 30px !important;
    border-radius: 20px !important;
  }

  .member-banner-content {
    flex-direction: column !important;
    text-align: center !important;
    gap: 8px !important;
  }

  .member-banner-icon {
    font-size: 2rem !important;
  }

  .member-banner-content h3 {
    font-size: 1.15rem !important;
  }

  .member-banner-content p {
    font-size: 0.82rem !important;
  }

  .member-banner-inner .btn {
    width: 100% !important;
    padding: 12px 20px !important;
    font-size: 0.9rem !important;
  }
}
</style>
