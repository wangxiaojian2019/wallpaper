<template>
  <div class="admin-dashboard">
    <!-- Stats Cards -->
    <div class="stats-grid">
      <div v-for="stat in stats" :key="stat.key" class="stat-card">
        <div class="stat-icon">{{ stat.icon }}</div>
        <div class="stat-body">
          <p class="stat-label">{{ locale === 'zh' ? stat.zh : stat.en }}</p>
          <p class="stat-value">{{ stat.prefix }}{{ stat.value }}{{ stat.suffix }}</p>
          <p class="stat-change" :class="stat.change > 0 ? 'positive' : 'negative'">
            {{ stat.change > 0 ? '+' : '' }}{{ stat.change }}% {{ locale === 'zh' ? '较昨日' : 'vs yesterday' }}
          </p>
        </div>
      </div>
    </div>

    <!-- Recent orders + Top products -->
    <div class="dashboard-grid">
      <div class="dash-card">
        <div class="dash-card-header">
          <h3>{{ locale === 'zh' ? '最近订单' : 'Recent Orders' }}</h3>
          <router-link to="/admin/orders" class="view-all">{{ locale === 'zh' ? '查看全部' : 'View all' }} →</router-link>
        </div>
        <div class="orders-table">
          <div v-for="order in recentOrders" :key="order.id" class="order-row">
            <div class="order-info">
              <span class="order-id">#{{ order.id?.slice(-6) || '000001' }}</span>
              <span class="order-user">{{ order.user_email || 'user@example.com' }}</span>
            </div>
            <span class="badge" :class="getStatusClass(order.status)">
              {{ getStatusLabel(order.status) }}
            </span>
            <span class="order-amount">$ {{ order.amount }}</span>
          </div>
          <div v-if="recentOrders.length === 0" class="empty-row">
            {{ t('common.no_data') }}
          </div>
        </div>
      </div>

      <div class="dash-card">
        <div class="dash-card-header">
          <h3>{{ locale === 'zh' ? '热门壁纸' : 'Top Wallpapers' }}</h3>
          <router-link to="/admin/products" class="view-all">{{ locale === 'zh' ? '查看全部' : 'View all' }} →</router-link>
        </div>
        <div class="top-products">
          <div v-for="(product, i) in topProducts" :key="product.id" class="top-product-row">
            <span class="rank">{{ i + 1 }}</span>
            <div class="product-thumb-sm">
              <img v-if="product.cover_image" :src="product.cover_image" />
              <div v-else class="thumb-ph" :style="{ background: product._color || '#667eea' }"></div>
            </div>
            <div class="product-name-sm">
              <p>{{ locale === 'zh' ? product.title_zh : product.title_en }}</p>
              <span>$ {{ product.price }}</span>
            </div>
            <span class="sales-count">{{ product.sales || 0 }} {{ locale === 'zh' ? '销量' : 'sales' }}</span>
          </div>
          <div v-if="topProducts.length === 0" class="empty-row">{{ t('common.no_data') }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { userAPI, orderAPI, productAPI } from '@/api'
import { useToast } from '@/composables/useToast'

const { t, locale } = useI18n()
const { error } = useToast()

const recentOrders = ref([])
const topProducts = ref([])

const stats = ref([
  { key: 'sales', icon: '💰', zh: '总销售额', en: 'Total Sales', prefix: '$', value: '0', suffix: '', change: 0 },
  { key: 'orders', icon: '📦', zh: '总订单', en: 'Total Orders', prefix: '', value: '0', suffix: '', change: 0 },
  { key: 'users', icon: '👥', zh: '总用户', en: 'Total Users', prefix: '', value: '0', suffix: '', change: 0 },
  { key: 'today', icon: '📈', zh: '今日销售', en: "Today's Sales", prefix: '$', value: '0', suffix: '', change: 0 }
])

function getStatusClass(status) {
  const map = { paid: 'badge-success', pending: 'badge-warning', cancelled: 'badge-error' }
  return map[status] || 'badge-neutral'
}

function getStatusLabel(status) {
  if (locale.value === 'zh') {
    const map = { paid: '已支付', pending: '待支付', cancelled: '已取消' }
    return map[status] || status
  }
  return status
}

async function loadData() {
  try {
    const res = await userAPI.stats()
    const s = res.stats || res
    stats.value[0].value = s.total_sales?.toFixed(2) || '0'
    stats.value[1].value = s.total_orders || '0'
    stats.value[2].value = s.total_users || '0'
    stats.value[3].value = s.today_sales?.toFixed(2) || '0'
    stats.value[0].change = s.sales_change || 0
    stats.value[1].change = s.orders_change || 0
    stats.value[2].change = s.users_change || 0
    stats.value[3].change = s.today_change || 0
  } catch {
    // Mock
    stats.value[0].value = '1,280.50'; stats.value[0].change = 12
    stats.value[1].value = '128'; stats.value[1].change = 8
    stats.value[2].value = '346'; stats.value[2].change = 5
    stats.value[3].value = '89.99'; stats.value[3].change = 15
  }

  try {
    const res = await orderAPI.adminList({ limit: 5 })
    recentOrders.value = res.orders || []
  } catch {
    recentOrders.value = [
      { id: 'abc123', user_email: 'alice@mail.com', status: 'paid', amount: '9.99' },
      { id: 'def456', user_email: 'bob@mail.com', status: 'pending', amount: '24.99' },
      { id: 'ghi789', user_email: 'carol@mail.com', status: 'paid', amount: '9.99' }
    ]
  }

  try {
    const res = await productAPI.list({ sort: 'popular', page_size: 5, limit: 5 })
    topProducts.value = res.items || res.products || []
  } catch {
    topProducts.value = [
      { id: '1', title_zh: '星空壁纸', title_en: 'Starry Night', price: '9.99', sales: 42, _color: 'linear-gradient(135deg,#667eea,#764ba2)' },
      { id: '2', title_zh: '森林绿意', title_en: 'Forest Green', price: '9.99', sales: 36, _color: 'linear-gradient(135deg,#43e97b,#38f9d7)' },
      { id: '3', title_zh: '城市夜景', title_en: 'City Night', price: '12.99', sales: 28, _color: 'linear-gradient(135deg,#f093fb,#f5576c)' }
    ]
  }
}

onMounted(loadData)
</script>

<style scoped>
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

@media (max-width: 900px) { .stats-grid { grid-template-columns: repeat(2, 1fr); } }
@media (max-width: 480px) { .stats-grid { grid-template-columns: 1fr; } }

.stat-card {
  background: var(--c-surface);
  border-radius: var(--r-xl);
  padding: 20px;
  display: flex;
  gap: 14px;
  align-items: flex-start;
  box-shadow: var(--shadow-xs);
}
.stat-icon { font-size: 1.8rem; flex-shrink: 0; }
.stat-label { font-size: 0.78rem; color: var(--c-text-tertiary); font-weight: 600; text-transform: uppercase; letter-spacing: 0.04em; margin-bottom: 4px; }
.stat-value { font-size: 1.6rem; font-weight: 800; letter-spacing: -0.03em; color: var(--c-text-primary); }
.stat-change { font-size: 0.75rem; font-weight: 600; margin-top: 4px; }
.stat-change.positive { color: var(--c-success); }
.stat-change.negative { color: var(--c-error); }

.dashboard-grid {
  display: grid;
  grid-template-columns: 1.2fr 1fr;
  gap: 16px;
}
@media (max-width: 900px) { .dashboard-grid { grid-template-columns: 1fr; } }

.dash-card {
  background: var(--c-surface);
  border-radius: var(--r-xl);
  padding: 20px;
  box-shadow: var(--shadow-xs);
}

.dash-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}
.dash-card-header h3 { font-size: 1rem; }
.view-all { font-size: 0.8rem; color: var(--c-accent); font-weight: 600; }

.order-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 0;
  border-bottom: 1px solid var(--c-border-light);
}
.order-row:last-child { border-bottom: none; }
.order-info { flex: 1; display: flex; flex-direction: column; gap: 2px; }
.order-id { font-size: 0.8rem; font-weight: 700; font-family: monospace; }
.order-user { font-size: 0.75rem; color: var(--c-text-tertiary); }
.order-amount { font-size: 0.88rem; font-weight: 700; white-space: nowrap; }
.empty-row { text-align: center; padding: 20px; color: var(--c-text-tertiary); font-size: 0.88rem; }

.top-product-row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 0;
  border-bottom: 1px solid var(--c-border-light);
}
.top-product-row:last-child { border-bottom: none; }
.rank { font-size: 0.78rem; font-weight: 800; color: var(--c-text-tertiary); width: 16px; flex-shrink: 0; }
.product-thumb-sm { width: 32px; height: 40px; border-radius: var(--r-xs); overflow: hidden; flex-shrink: 0; }
.product-thumb-sm img { width: 100%; height: 100%; object-fit: cover; }
.thumb-ph { width: 100%; height: 100%; }
.product-name-sm { flex: 1; }
.product-name-sm p { font-size: 0.82rem; font-weight: 600; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; max-width: 120px; }
.product-name-sm span { font-size: 0.72rem; color: var(--c-text-tertiary); }
.sales-count { font-size: 0.75rem; color: var(--c-text-secondary); white-space: nowrap; }
</style>
