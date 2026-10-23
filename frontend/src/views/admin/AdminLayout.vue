<template>
  <div class="admin-layout">
    <!-- Mobile Drawer Backdrop -->
    <div v-if="mobileDrawerOpen" class="sidebar-backdrop" @click="mobileDrawerOpen = false"></div>

    <!-- Sidebar -->
    <aside class="admin-sidebar" :class="{ collapsed: sidebarCollapsed, 'mobile-open': mobileDrawerOpen }">
      <div class="sidebar-header">
        <div class="sidebar-logo">
          <span class="logo-icon">✦</span>
          <span v-if="!sidebarCollapsed">WallArt Admin</span>
        </div>
        <button class="collapse-btn hide-mobile-admin" @click="sidebarCollapsed = !sidebarCollapsed">
          {{ sidebarCollapsed ? '→' : '←' }}
        </button>
        <button class="collapse-btn show-mobile-admin" @click="mobileDrawerOpen = false">
          ✕
        </button>
      </div>

      <nav class="sidebar-nav">
        <router-link v-for="item in navItems" :key="item.to"
          :to="item.to" class="sidebar-link"
          :class="{ active: $route.path.startsWith(item.to) && (item.to !== '/admin' || $route.path === '/admin') }"
          :title="sidebarCollapsed ? (locale === 'zh' ? item.zh : item.en) : ''"
          @click="mobileDrawerOpen = false"
        >
          <span class="nav-icon">{{ item.icon }}</span>
          <span v-if="!sidebarCollapsed" class="nav-label">{{ locale === 'zh' ? item.zh : item.en }}</span>
          <span v-if="item.to === '/admin/orders' && pendingReviewCount > 0" class="sidebar-pulse-badge" :title="`有 ${pendingReviewCount} 条待审核订单`">
            {{ pendingReviewCount }}
          </span>
        </router-link>
      </nav>

      <div class="sidebar-footer">
        <a :href="frontendUrl" target="_blank" class="sidebar-link" style="color: #38bdf8; font-weight: 600;" :title="sidebarCollapsed ? '打开 Uni-app 新版前台 (带星空/雨滴双模式)' : ''">
          <span class="nav-icon">✨</span>
          <span v-if="!sidebarCollapsed">Uni-app 前台 ↗</span>
        </a>
        <router-link to="/" class="sidebar-link" :title="sidebarCollapsed ? (locale === 'zh' ? '旧版前台' : 'Legacy Frontend') : ''" @click="mobileDrawerOpen = false">
          <span class="nav-icon">🏠</span>
          <span v-if="!sidebarCollapsed">{{ locale === 'zh' ? '旧版前台' : 'Frontend' }}</span>
        </router-link>
      </div>
    </aside>

    <!-- Main content -->
    <div class="admin-main">
      <!-- Top bar -->
      <header class="admin-topbar">
        <div class="topbar-left">
          <button class="mobile-toggle-btn show-mobile-admin" @click="mobileDrawerOpen = !mobileDrawerOpen">
            ☰
          </button>
          <h2 class="page-title">{{ currentPageTitle }}</h2>
        </div>
        <div class="topbar-right">
          <button class="lang-btn" @click="toggleLang">
            🌐 {{ locale === 'zh' ? 'EN' : '中文' }}
          </button>
          <div class="admin-user">
            <span class="admin-avatar">{{ avatarLetter }}</span>
            <span class="admin-email hide-mobile">{{ auth.user?.email }}</span>
          </div>
        </div>
      </header>

      <!-- Content -->
      <div class="admin-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </div>
    </div>

    <!-- New Order Floating Alert Notification -->
    <transition name="toast-slide">
      <div v-if="newOrderAlert" class="order-alert-toast" @click.stop>
        <div class="toast-header">
          <div class="toast-title-row">
            <span class="bell-icon-pulse">🔔</span>
            <span class="toast-headline">{{ locale === 'zh' ? '新订单通知！' : 'New Order Alert!' }}</span>
            <span class="toast-now-pill">{{ locale === 'zh' ? '刚刚' : 'Just now' }}</span>
          </div>
          <button class="toast-close-btn" @click="newOrderAlert = null">✕</button>
        </div>
        <div class="toast-content">
          <div class="toast-info-line">
            <span class="t-k">{{ locale === 'zh' ? '买家' : 'Buyer' }}：</span>
            <span class="t-v font-bold">{{ newOrderAlert.user_email || '游客' }}</span>
          </div>
          <div class="toast-info-line">
            <span class="t-k">{{ locale === 'zh' ? '内容' : 'Item' }}：</span>
            <span class="t-v text-truncate">{{ newOrderAlert.product_title || (newOrderAlert.membership_type ? 'VIP会员订阅' : '单张壁纸') }}</span>
          </div>
          <div class="toast-info-line price-highlight">
            <span class="t-k">{{ locale === 'zh' ? '金额' : 'Amount' }}：</span>
            <span class="t-price">
              {{ (newOrderAlert.currency || '').toUpperCase() === 'CNY' ? '¥' : '$' }}{{ Number(newOrderAlert.amount || 0).toFixed(2) }}
              <span class="t-cur-badge" :class="(newOrderAlert.currency || '').toUpperCase() === 'CNY' ? 't-cny' : 't-usd'">
                {{ newOrderAlert.currency || 'USD' }}
              </span>
            </span>
          </div>
          <div class="toast-info-line">
            <span class="t-k">{{ locale === 'zh' ? '状态' : 'Status' }}：</span>
            <span class="t-status-tag" :class="newOrderAlert.status">
              {{ newOrderAlert.status === 'pending_review' ? (locale === 'zh' ? '⏳ 待审核' : 'Pending Review') : (locale === 'zh' ? '✓ 已支付' : newOrderAlert.status) }}
            </span>
          </div>
        </div>
        <div class="toast-actions">
          <button class="btn btn-sm btn-primary flex-1" @click="viewOrderFromAlert(newOrderAlert)">
            {{ locale === 'zh' ? '🔍 立即查看' : 'View' }}
          </button>
          <button class="btn btn-sm btn-outline" @click="newOrderAlert = null">
            {{ locale === 'zh' ? '稍后处理' : 'Dismiss' }}
          </button>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/store'
import { orderAPI } from '@/api'

const { t, locale } = useI18n()
const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const sidebarCollapsed = ref(false)
const mobileDrawerOpen = ref(false)
const avatarLetter = computed(() => (auth.user?.email || 'A').charAt(0).toUpperCase())

const pendingReviewCount = ref(0)
const newOrderAlert = ref(null)
let latestOrderId = null
let isFirstCheck = true
let pollTimer = null

function playOrderChime() {
  try {
    const AudioContextClass = window.AudioContext || window.webkitAudioContext
    if (!AudioContextClass) return
    const ctx = new AudioContextClass()
    if (ctx.state === 'suspended') {
      ctx.resume().catch(() => {})
    }
    const now = ctx.currentTime
    const osc = ctx.createOscillator()
    const gain = ctx.createGain()
    osc.type = 'triangle'
    osc.frequency.setValueAtTime(659.25, now) // E5
    osc.frequency.exponentialRampToValueAtTime(987.77, now + 0.12) // B5
    gain.gain.setValueAtTime(0.2, now)
    gain.gain.exponentialRampToValueAtTime(0.001, now + 0.6)
    osc.connect(gain)
    gain.connect(ctx.destination)
    osc.start(now)
    osc.stop(now + 0.6)
  } catch (e) {}
}

async function checkNewOrders() {
  if (!auth.token) return
  try {
    const res = await orderAPI.adminList({ page: 1, page_size: 5 })
    const orders = res.orders || res.items || res || []
    if (res.pending_review_total !== undefined) {
      pendingReviewCount.value = Number(res.pending_review_total) || 0
    } else {
      pendingReviewCount.value = orders.filter(o => o.status === 'pending_review').length
    }

    if (orders.length > 0) {
      const newest = orders[0]
      if (isFirstCheck) {
        latestOrderId = newest.id
        isFirstCheck = false
        return
      }

      if (latestOrderId && newest.id !== latestOrderId) {
        latestOrderId = newest.id
        newOrderAlert.value = newest
        playOrderChime()
        // Auto-dismiss after 15 seconds if unhandled
        setTimeout(() => {
          if (newOrderAlert.value?.id === newest.id) {
            newOrderAlert.value = null
          }
        }, 15000)
      }
    }
  } catch (err) {
    // silent catch
  }
}

function viewOrderFromAlert(order) {
  newOrderAlert.value = null
  router.push('/admin/orders')
}

onMounted(() => {
  checkNewOrders()
  pollTimer = setInterval(checkNewOrders, 12000)
})

onBeforeUnmount(() => {
  if (pollTimer) clearInterval(pollTimer)
})

function toggleLang() {
  locale.value = locale.value === 'zh' ? 'en' : 'zh'
  localStorage.setItem('lang', locale.value)
}

const navItems = [
  { to: '/admin', icon: '📊', zh: '数据看板', en: 'Dashboard' },
  { to: '/admin/products', icon: '🖼️', zh: '商品管理', en: 'Products' },
  { to: '/admin/categories', icon: '📁', zh: '分类管理', en: 'Categories' },
  { to: '/admin/tutorial', icon: '📖', zh: '教程模板', en: 'Tutorial Guide' },
  { to: '/admin/orders', icon: '📦', zh: '订单管理', en: 'Orders' },
  { to: '/admin/users', icon: '👥', zh: '会员管理', en: 'Users' },
  { to: '/admin/visitors', icon: '👁️', zh: '访客管理', en: 'Visitors' },
  { to: '/admin/stats', icon: '📈', zh: '数据统计', en: 'Statistics' },
  { to: '/admin/settings', icon: '⚙️', zh: '支付与系统设置', en: 'Settings' },
  { to: '/admin/redeem-codes', icon: '🎫', zh: '口令兑换', en: 'Redeem Codes' }
]

const currentPageTitle = computed(() => {
  const item = navItems.find(n =>
    n.to === '/admin' ? route.path === '/admin' : route.path.startsWith(n.to)
  )
  return item ? (locale.value === 'zh' ? item.zh : item.en) : 'Admin'
})

const frontendUrl = computed(() => {
  if (typeof window === 'undefined') return 'https://wallpaper.icangoanyway.com/#/'
  const host = window.location.hostname
  const protocol = window.location.protocol
  if (host === 'localhost' || host === '127.0.0.1') {
    return 'http://localhost:5173/#/'
  }
  if (host.includes('admin.icangoanyway.com')) {
    return `${protocol}//wallpaper.icangoanyway.com/#/`
  }
  if (window.location.port === '5174') {
    return `${protocol}//${host}:5173/#/`
  }
  return `${protocol}//wallpaper.icangoanyway.com/#/`
})
</script>

<style scoped>
.admin-layout {
  display: flex;
  min-height: 100vh;
  background: var(--c-bg);
}

/* Sidebar */
.admin-sidebar {
  width: 240px;
  background: var(--c-text-primary);
  color: white;
  display: flex;
  flex-direction: column;
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  z-index: 100;
  transition: width var(--t-normal), transform var(--t-normal);
  overflow: hidden;
}

.admin-sidebar.collapsed { width: 64px; }

.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 16px;
  border-bottom: 1px solid rgba(255,255,255,0.1);
  height: var(--nav-h);
}

.sidebar-logo {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 800;
  font-size: 0.95rem;
  white-space: nowrap;
  overflow: hidden;
}
.logo-icon { color: #5eb5ff; font-size: 1rem; flex-shrink: 0; }

.collapse-btn {
  background: rgba(255,255,255,0.1);
  border: none;
  color: rgba(255,255,255,0.7);
  width: 28px;
  height: 28px;
  border-radius: var(--r-sm);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.8rem;
  flex-shrink: 0;
  transition: all var(--t-fast);
}
.collapse-btn:hover { background: rgba(255,255,255,0.2); color: white; }

.sidebar-nav {
  flex: 1;
  padding: 12px 8px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  overflow-y: auto;
}

.sidebar-link {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: var(--r-md);
  color: rgba(255,255,255,0.65);
  font-size: 0.88rem;
  font-weight: 500;
  transition: all var(--t-fast);
  white-space: nowrap;
  overflow: hidden;
  text-decoration: none;
}
.sidebar-link:hover { background: rgba(255,255,255,0.1); color: white; }
.sidebar-link.active { background: rgba(255,255,255,0.15); color: white; }

.nav-icon { font-size: 1rem; flex-shrink: 0; }

.sidebar-footer {
  padding: 8px;
  border-top: 1px solid rgba(255,255,255,0.1);
}

/* Main */
.admin-main {
  flex: 1;
  margin-left: 240px;
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  transition: margin-left var(--t-normal);
}

.admin-sidebar.collapsed ~ .admin-main { margin-left: 64px; }

.admin-topbar {
  height: var(--nav-h);
  background: var(--c-surface);
  border-bottom: 1px solid var(--c-border-light);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  position: sticky;
  top: 0;
  z-index: 50;
}

.topbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.page-title { font-size: 1.1rem; font-weight: 700; }

.topbar-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.lang-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  border-radius: var(--r-full);
  background: var(--c-bg);
  border: 1px solid var(--c-border);
  font-size: 0.8rem;
  font-weight: 600;
  color: var(--c-text-secondary);
  cursor: pointer;
  transition: all var(--t-fast);
}
.lang-btn:hover { border-color: var(--c-accent); color: var(--c-accent); }

.admin-user {
  display: flex;
  align-items: center;
  gap: 8px;
}

.admin-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--g-accent);
  color: white;
  font-weight: 700;
  font-size: 0.85rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.admin-email {
  font-size: 0.82rem;
  color: var(--c-text-secondary);
}

.admin-content {
  flex: 1;
  padding: 28px;
  overflow-y: auto;
}

.sidebar-backdrop {
  display: none;
}

.show-mobile-admin {
  display: none;
}

.hide-mobile-admin {
  display: flex;
}

@media (max-width: 768px) {
  .show-mobile-admin {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 34px;
    height: 34px;
    background: transparent;
    border: none;
    font-size: 1.25rem;
    color: var(--c-text-primary);
    cursor: pointer;
  }

  .hide-mobile-admin {
    display: none !important;
  }

  .sidebar-backdrop {
    display: block;
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.45);
    backdrop-filter: blur(4px);
    z-index: 199;
  }

  .admin-sidebar {
    width: 250px !important;
    transform: translateX(-100%);
    transition: transform 0.3s cubic-bezier(0.16, 1, 0.3, 1) !important;
    z-index: 200;
  }

  .admin-sidebar.mobile-open {
    transform: translateX(0) !important;
  }

  .admin-main {
    margin-left: 0 !important;
    width: 100% !important;
  }

  .admin-topbar {
    padding: 0 16px;
  }

  .admin-content {
    padding: 16px 12px !important;
  }
}

/* Sidebar Pulse Badge */
.sidebar-pulse-badge {
  margin-left: auto;
  background: #ef4444;
  color: white;
  font-size: 11px;
  font-weight: 800;
  padding: 1px 7px;
  border-radius: 999px;
  box-shadow: 0 0 10px rgba(239, 68, 68, 0.6);
  animation: badgePulse 2s infinite;
}

@keyframes badgePulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.18); }
}

/* Order Alert Toast */
.order-alert-toast {
  position: fixed;
  bottom: 24px;
  right: 24px;
  width: 340px;
  max-width: calc(100vw - 32px);
  background: #1e293b;
  color: #f8fafc;
  border: 1.5px solid rgba(245, 158, 11, 0.4);
  border-radius: 16px;
  padding: 16px;
  box-shadow: 0 12px 36px rgba(0, 0, 0, 0.5), 0 0 24px rgba(245, 158, 11, 0.25);
  z-index: 99999;
  backdrop-filter: blur(12px);
}

.toast-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  padding-bottom: 8px;
}

.toast-title-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.bell-icon-pulse {
  font-size: 18px;
  animation: bellRing 1.5s infinite;
}

@keyframes bellRing {
  0%, 100% { transform: rotate(0); }
  15% { transform: rotate(14deg); }
  30% { transform: rotate(-14deg); }
  45% { transform: rotate(8deg); }
  60% { transform: rotate(-4deg); }
  75% { transform: rotate(0); }
}

.toast-headline {
  font-size: 14px;
  font-weight: 800;
  color: #fbbf24;
}

.toast-now-pill {
  font-size: 10px;
  padding: 1px 6px;
  border-radius: 999px;
  background: rgba(245, 158, 11, 0.2);
  color: #f59e0b;
}

.toast-close-btn {
  background: none;
  border: none;
  color: #94a3b8;
  font-size: 16px;
  cursor: pointer;
  padding: 2px 6px;
  border-radius: 6px;
}

.toast-close-btn:hover {
  color: white;
  background: rgba(255, 255, 255, 0.1);
}

.toast-content {
  display: flex;
  flex-direction: column;
  gap: 6px;
  font-size: 12px;
  margin-bottom: 14px;
}

.toast-info-line {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.toast-info-line.price-highlight {
  padding: 5px 8px;
  background: rgba(245, 158, 11, 0.08);
  border-radius: 8px;
  margin: 4px 0;
}

.t-k {
  color: #94a3b8;
  flex-shrink: 0;
}

.t-v {
  color: #e2e8f0;
}

.t-price {
  font-size: 16px;
  font-weight: 800;
  color: #fbbf24;
  display: flex;
  align-items: center;
  gap: 4px;
}

.t-cur-badge {
  font-size: 10px;
  padding: 1px 4px;
  border-radius: 4px;
  font-weight: 700;
}

.t-cur-badge.t-cny {
  background: rgba(245, 158, 11, 0.2);
  color: #fbbf24;
}

.t-cur-badge.t-usd {
  background: rgba(16, 185, 129, 0.2);
  color: #34d399;
}

.t-status-tag {
  font-size: 11px;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 999px;
}

.t-status-tag.pending_review {
  background: rgba(245, 158, 11, 0.2);
  color: #fbbf24;
}

.t-status-tag.paid {
  background: rgba(16, 185, 129, 0.2);
  color: #34d399;
}

.toast-actions {
  display: flex;
  gap: 8px;
}

/* Toast Transitions */
.toast-slide-enter-active,
.toast-slide-leave-active {
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

.toast-slide-enter-from,
.toast-slide-leave-to {
  opacity: 0;
  transform: translateY(30px) scale(0.95);
}
</style>
