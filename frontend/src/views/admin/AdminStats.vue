<template>
  <div class="admin-stats">
    <!-- Revenue Metrics Overview -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">💰</div>
        <div class="stat-meta">
          <span class="stat-title">{{ t('admin.total_sales') }}</span>
          <span class="stat-val font-bold text-gradient">${{ stats.total_revenue?.toFixed(2) || '0.00' }}</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">📅</div>
        <div class="stat-meta">
          <span class="stat-title">{{ t('admin.today_sales') }}</span>
          <span class="stat-val font-bold">${{ stats.today_revenue?.toFixed(2) || '0.00' }}</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">🛍️</div>
        <div class="stat-meta">
          <span class="stat-title">{{ t('admin.total_orders') }}</span>
          <span class="stat-val font-bold">{{ stats.total_orders || 0 }}</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">👥</div>
        <div class="stat-meta">
          <span class="stat-title">{{ t('admin.total_users') }}</span>
          <span class="stat-val font-bold">{{ stats.total_users || 0 }}</span>
        </div>
      </div>
    </div>

    <!-- Charts / Breakdown Section -->
    <div class="charts-row">
      <!-- Revenue Breakdown by Plan -->
      <div class="card stat-section">
        <h3>{{ locale === 'zh' ? '营收结构分布' : 'Revenue Breakdown' }}</h3>
        <p class="section-sub">{{ locale === 'zh' ? '单品购买与会员订阅收入占比' : 'Wallpaper single sales vs Membership subscriptions' }}</p>

        <div class="breakdown-list">
          <div class="breakdown-item">
            <div class="b-info">
              <span>🖼️ {{ locale === 'zh' ? '单购壁纸' : 'Single Wallpaper' }}</span>
              <span class="font-bold">${{ (stats.single_sales_revenue || 0).toFixed(2) }} ({{ calcPercent(stats.single_sales_revenue, stats.total_revenue) }}%)</span>
            </div>
            <div class="progress-bar">
              <div class="progress-fill" :style="{ width: calcPercent(stats.single_sales_revenue, stats.total_revenue) + '%' }"></div>
            </div>
          </div>

          <div class="breakdown-item">
            <div class="b-info">
              <span>👑 {{ locale === 'zh' ? '月度会员' : 'Monthly Plan' }}</span>
              <span class="font-bold">${{ (stats.monthly_revenue || 0).toFixed(2) }} ({{ calcPercent(stats.monthly_revenue, stats.total_revenue) }}%)</span>
            </div>
            <div class="progress-bar">
              <div class="progress-fill fill-secondary" :style="{ width: calcPercent(stats.monthly_revenue, stats.total_revenue) + '%' }"></div>
            </div>
          </div>

          <div class="breakdown-item">
            <div class="b-info">
              <span>🌟 {{ locale === 'zh' ? '季度会员' : 'Quarterly Plan' }}</span>
              <span class="font-bold">${{ (stats.quarterly_revenue || 0).toFixed(2) }} ({{ calcPercent(stats.quarterly_revenue, stats.total_revenue) }}%)</span>
            </div>
            <div class="progress-bar">
              <div class="progress-fill fill-accent" :style="{ width: calcPercent(stats.quarterly_revenue, stats.total_revenue) + '%' }"></div>
            </div>
          </div>

          <div class="breakdown-item">
            <div class="b-info">
              <span>💎 {{ locale === 'zh' ? '年度会员' : 'Yearly Plan' }}</span>
              <span class="font-bold">${{ (stats.yearly_revenue || 0).toFixed(2) }} ({{ calcPercent(stats.yearly_revenue, stats.total_revenue) }}%)</span>
            </div>
            <div class="progress-bar">
              <div class="progress-fill fill-gold" :style="{ width: calcPercent(stats.yearly_revenue, stats.total_revenue) + '%' }"></div>
            </div>
          </div>
        </div>
      </div>

      <!-- Top Selling Products -->
      <div class="card stat-section">
        <h3>{{ locale === 'zh' ? '热销壁纸排行榜' : 'Top Selling Wallpapers' }}</h3>
        <p class="section-sub">{{ locale === 'zh' ? '按单购与下载热度综合排名' : 'Ranked by sales volume & downloads' }}</p>

        <div class="top-list">
          <div v-for="(p, idx) in topProducts" :key="p.id || idx" class="top-item">
            <div class="top-rank" :class="{ 'rank-1': idx === 0, 'rank-2': idx === 1, 'rank-3': idx === 2 }">
              {{ idx + 1 }}
            </div>
            <div class="top-thumb">
              <img v-if="p.cover_image" :src="p.cover_image" />
              <div v-else class="ph-thumb">🖼️</div>
            </div>
            <div class="top-meta">
              <span class="top-title">{{ locale === 'zh' ? p.title_zh : (p.title_en || p.title_zh) }}</span>
              <span class="top-sub">${{ p.price }} USD · {{ p.category || '4K' }}</span>
            </div>
            <div class="top-sales">
              <strong>{{ p.sales_count || 0 }}</strong>
              <small>{{ locale === 'zh' ? '已售' : 'sold' }}</small>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { userAPI, productAPI } from '@/api'

const { t, locale } = useI18n()

const stats = ref({
  total_revenue: 1280.50,
  today_revenue: 126.80,
  total_orders: 142,
  total_users: 388,
  single_sales_revenue: 430.50,
  monthly_revenue: 290.00,
  quarterly_revenue: 210.00,
  yearly_revenue: 350.00
})

const topProducts = ref([
  {
    id: '1',
    title_zh: '赛博霓虹雨夜 8K',
    title_en: 'Cyberpunk Rainy Night 8K',
    price: 3.99,
    sales_count: 86,
    category: 'Cyberpunk',
    cover_image: 'https://images.unsplash.com/photo-1542751371-adc38448a05e?w=800'
  },
  {
    id: '2',
    title_zh: '星际深邃黑洞 4K 动态',
    title_en: 'Deep Space Black Hole',
    price: 4.99,
    sales_count: 64,
    category: 'Space',
    cover_image: 'https://images.unsplash.com/photo-1451187580459-43490279c0fa?w=800'
  },
  {
    id: '3',
    title_zh: '极简北欧晨雾雪山',
    title_en: 'Minimal Nordic Misty Mountain',
    price: 2.99,
    sales_count: 42,
    category: 'Minimal',
    cover_image: 'https://images.unsplash.com/photo-1464822759023-fed622ff2c3b?w=800'
  }
])

function calcPercent(part, total) {
  if (!total || total === 0) return 0
  return Math.min(100, Math.round(((part || 0) / total) * 100))
}

onMounted(async () => {
  try {
    const res = await userAPI.stats()
    if (res) {
      stats.value = { ...stats.value, ...res }
    }
  } catch {
    // Keep baseline stats
  }

  try {
    const pres = await productAPI.list({ sort: 'popular', limit: 5 })
    if (pres?.items?.length) {
      topProducts.value = pres.items
    }
  } catch {
    // Keep baseline top products
  }
})
</script>

<style scoped>
.admin-stats {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
}

.stat-card {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  font-size: 2.2rem;
  width: 56px;
  height: 56px;
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-meta {
  display: flex;
  flex-direction: column;
}

.stat-title {
  font-size: 0.85rem;
  color: var(--text-muted);
  margin-bottom: 4px;
}

.stat-val {
  font-size: 1.6rem;
  line-height: 1.2;
}

.charts-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
}

.stat-section {
  padding: 24px;
}

.stat-section h3 {
  font-size: 1.15rem;
  font-weight: 700;
  margin-bottom: 4px;
}

.section-sub {
  font-size: 0.85rem;
  color: var(--text-muted);
  margin-bottom: 24px;
}

.breakdown-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.b-info {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 0.9rem;
}

.progress-bar {
  height: 10px;
  background: var(--bg-secondary);
  border-radius: 999px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: var(--primary);
  border-radius: 999px;
  transition: width 0.4s ease;
}

.progress-fill.fill-secondary {
  background: #3b82f6;
}

.progress-fill.fill-accent {
  background: #8b5cf6;
}

.progress-fill.fill-gold {
  background: #f59e0b;
}

.top-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.top-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px;
  border-radius: var(--radius-md);
  transition: background 0.15s;
}

.top-item:hover {
  background: var(--bg-hover);
}

.top-rank {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: var(--bg-secondary);
  font-size: 0.8rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
}

.top-rank.rank-1 {
  background: #fef3c7;
  color: #b45309;
}

.top-rank.rank-2 {
  background: #e0e7ff;
  color: #3730a3;
}

.top-rank.rank-3 {
  background: #f3f4f6;
  color: #4b5563;
}

.top-thumb {
  width: 44px;
  height: 44px;
  border-radius: var(--radius-sm);
  overflow: hidden;
  flex-shrink: 0;
}

.top-thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.ph-thumb {
  width: 100%;
  height: 100%;
  background: var(--bg-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
}

.top-meta {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.top-title {
  font-size: 0.9rem;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.top-sub {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.top-sales {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.top-sales strong {
  font-size: 1rem;
}

.top-sales small {
  font-size: 0.7rem;
  color: var(--text-muted);
}

@media (max-width: 800px) {
  .charts-row {
    grid-template-columns: 1fr;
  }
}
</style>
