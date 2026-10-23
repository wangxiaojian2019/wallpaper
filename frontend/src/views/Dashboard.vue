<template>
  <div class="dashboard-page">
    <div class="container">
      <div class="dashboard-header">
        <div class="user-avatar">{{ avatarLetter }}</div>
        <div>
          <h1>{{ t('dashboard.title') }}</h1>
          <p class="user-email">{{ auth.user?.email }}</p>
        </div>
      </div>

      <!-- Member status card -->
      <div class="status-cards">
        <div class="status-card" :class="auth.isMember ? 'status-gold' : 'status-neutral'">
          <div class="status-icon">{{ auth.isMember ? '👑' : '🔓' }}</div>
          <div>
            <p class="status-label">{{ t('dashboard.membership') }}</p>
            <p class="status-value">
              <template v-if="auth.isMember">
                <span class="badge badge-gold">{{ t(`membership.${auth.membershipType}`) }}</span>
                <span class="expire-date">
                  {{ t('membership.expires') }}: {{ formatDate(auth.user?.membership_expire) }}
                </span>
              </template>
              <template v-else>
                <span class="no-member">{{ t('dashboard.no_member') }}</span>
                <router-link to="/membership" class="btn btn-gold btn-sm" style="margin-top:8px">
                  👑 {{ locale === 'zh' ? '立即升级' : 'Upgrade Now' }}
                </router-link>
              </template>
            </p>
          </div>
        </div>

        <div class="status-card">
          <div class="status-icon">🖼️</div>
          <div>
            <p class="status-label">{{ t('dashboard.purchases') }}</p>
            <p class="status-value">
              <strong>{{ purchases.length }}</strong>
              {{ locale === 'zh' ? '件壁纸' : 'wallpapers' }}
            </p>
          </div>
        </div>
      </div>

      <!-- Purchased wallpapers -->
      <div class="section-block">
        <h2>{{ t('dashboard.purchases') }}</h2>

        <div v-if="loadingPurchases" class="loading-state">
          <span class="spinner"></span>
        </div>

        <div v-else-if="purchases.length === 0" class="empty-purchases">
          <span>🛍️</span>
          <p>{{ t('dashboard.no_purchases') }}</p>
          <router-link to="/" class="btn btn-primary btn-sm">
            {{ locale === 'zh' ? '去逛逛 →' : 'Browse Wallpapers →' }}
          </router-link>
        </div>

        <div v-else class="purchase-grid">
          <div v-for="order in purchases" :key="order.id" class="purchase-item">
            <div class="purchase-thumb">
              <img v-if="order.product?.cover_image" :src="order.product.cover_image" :alt="order.product?.title_zh" />
              <div v-else class="thumb-placeholder" :style="{ background: order.product?._color || '#667eea' }">🖼️</div>
            </div>
            <div class="purchase-info">
              <p class="purchase-name">
                {{ locale === 'zh' ? order.product?.title_zh : order.product?.title_en }}
              </p>
              <p class="purchase-date">{{ formatDate(order.created_at) }}</p>
              <p class="purchase-amount">$ {{ order.amount }}</p>
            </div>
            <router-link :to="`/wallpaper/${order.product_id}`" class="btn btn-ghost btn-sm">
              {{ locale === 'zh' ? '查看' : 'View' }}
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/store'
import { orderAPI } from '@/api'

const { t, locale } = useI18n()
const auth = useAuthStore()

const purchases = ref([])
const loadingPurchases = ref(true)

const avatarLetter = computed(() => (auth.user?.email || '?').charAt(0).toUpperCase())

function formatDate(dateStr) {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString(locale.value === 'zh' ? 'zh-CN' : 'en-US')
}

async function loadPurchases() {
  loadingPurchases.value = true
  try {
    const res = await orderAPI.list()
    purchases.value = (res.orders || res || []).filter(o => o.status === 'paid' && o.order_type === 'single')
  } catch { purchases.value = [] }
  finally { loadingPurchases.value = false }
}

onMounted(loadPurchases)
</script>

<style scoped>
.dashboard-page { padding: 40px 0 80px; min-height: 100vh; }

.dashboard-header {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 40px;
  padding-bottom: 32px;
  border-bottom: 1px solid var(--c-border-light);
}

.user-avatar {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  background: var(--g-accent);
  color: white;
  font-size: 1.8rem;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.dashboard-header h1 { font-size: 1.8rem; margin-bottom: 4px; }
.user-email { color: var(--c-text-secondary); font-size: 0.9rem; }

.status-cards {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  margin-bottom: 48px;
}

@media (max-width: 600px) { .status-cards { grid-template-columns: 1fr; } }

.status-card {
  background: var(--c-surface);
  border: 1.5px solid var(--c-border-light);
  border-radius: var(--r-xl);
  padding: 24px;
  display: flex;
  gap: 16px;
  align-items: flex-start;
}

.status-card.status-gold {
  background: var(--c-gold-light);
  border-color: rgba(191,138,48,0.2);
}

.status-icon { font-size: 2rem; flex-shrink: 0; }
.status-label { font-size: 0.78rem; color: var(--c-text-tertiary); text-transform: uppercase; letter-spacing: 0.06em; font-weight: 600; margin-bottom: 6px; }
.status-value { display: flex; flex-direction: column; gap: 4px; }
.status-value strong { font-size: 1.5rem; font-weight: 800; }
.expire-date { font-size: 0.78rem; color: var(--c-text-secondary); margin-top: 4px; }
.no-member { font-size: 0.88rem; color: var(--c-text-secondary); }

.section-block h2 { margin-bottom: 20px; font-size: 1.3rem; }

.loading-state { display: flex; justify-content: center; padding: 40px; }

.empty-purchases {
  text-align: center;
  padding: 60px 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  font-size: 2.5rem;
}
.empty-purchases p { font-size: 0.95rem; color: var(--c-text-secondary); }

.purchase-grid { display: flex; flex-direction: column; gap: 10px; }

.purchase-item {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px;
  background: var(--c-surface);
  border: 1px solid var(--c-border-light);
  border-radius: var(--r-lg);
  transition: all var(--t-fast);
}
.purchase-item:hover { box-shadow: var(--shadow-sm); }

.purchase-thumb {
  width: 48px;
  height: 64px;
  border-radius: var(--r-sm);
  overflow: hidden;
  flex-shrink: 0;
}
.purchase-thumb img { width: 100%; height: 100%; object-fit: cover; }
.thumb-placeholder {
  width: 100%; height: 100%;
  display: flex; align-items: center; justify-content: center;
  font-size: 1.2rem; border-radius: var(--r-sm);
}

.purchase-info { flex: 1; }
.purchase-name { font-weight: 600; font-size: 0.9rem; color: var(--c-text-primary); margin-bottom: 4px; }
.purchase-date { font-size: 0.75rem; color: var(--c-text-tertiary); margin-bottom: 4px; }
.purchase-amount { font-size: 0.82rem; font-weight: 600; color: var(--c-accent); }
</style>
