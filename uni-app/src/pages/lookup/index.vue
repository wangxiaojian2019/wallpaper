<template>
  <view class="lookup-container" :class="[appStore.theme]">
    <AmbientCanvas />

    <view class="lookup-content">
      <!-- Back button -->
      <view class="top-nav">
        <view class="back-pill" @tap="goBack">
          <text class="back-arrow">←</text>
          <text class="back-text">返回</text>
        </view>
      </view>

      <view class="lookup-header">
        <view class="header-badge">
          <text class="badge-icon">🔍</text>
          <text class="badge-text">订单凭据提取通道</text>
        </view>
        <text class="lookup-title">一键找回已购壁纸与网盘提取码</text>
        <text class="lookup-sub">输入您购买时预留的邮箱地址，即刻同步获取全部 4K 原图下载口令</text>
      </view>

      <!-- Query Form -->
      <view class="form-box">
        <input
          v-model="email"
          placeholder="请输入下单时预留的邮箱地址"
          class="email-input"
          @confirm="handleQuery"
        />
        <button class="btn-query" :loading="loading" @tap="handleQuery">
          立即查询
        </button>
      </view>

      <!-- Query Results -->
      <view v-if="orders.length > 0" class="results-list">
        <text class="results-title">✦ 查询到 {{ orders.length }} 笔订单记录：</text>
        <view v-for="order in orders" :key="order.id" class="order-card">
          <view class="order-head">
            <text class="order-no">订单号：{{ order.id }}</text>
            <view class="status-pill" :class="{ 'status-pill-review': order.status === 'pending_review' }">
              <text class="status-dot">●</text>
              <text class="status-text">{{ order.status === 'pending_review' ? '⏳ 待审核' : '已支付' }}</text>
            </view>
          </view>
          <text class="order-name">{{ order.product_title || (order.order_type === 'membership' ? '👑 VIP 会员开通' : '精选壁纸 4K 原图') }}</text>
          
          <!-- Direct download links display if paid and available -->
          <view v-if="order.status === 'paid' && order.download_links && order.download_links.length > 0" class="direct-drives-list">
            <view v-for="(dl, dlIdx) in order.download_links" :key="dlIdx" class="direct-drive-item">
              <view class="direct-drive-info">
                <text class="direct-drive-plat">{{ dl.platform === 'quark' ? '夸克网盘' : (dl.platform === 'baidu' ? '百度网盘' : (dl.platform === 'google' ? 'Google Drive' : (dl.platform || '网盘'))) }}</text>
                <text v-if="dl.extract_code" class="direct-drive-code">提取码: {{ dl.extract_code }}</text>
              </view>
              <a :href="dl.url" target="_blank" rel="noopener noreferrer" class="direct-drive-link-str" @tap.stop>
                <text class="link-str-val">{{ dl.url }}</text>
                <text class="link-str-icon">↗</text>
              </a>
              <view class="direct-drive-btns">
                <button class="btn-open-sm" @tap.stop="openExternalUrl(dl.url)">打开 ↗</button>
                <button class="btn-copy-sm" @tap.stop="copyDriveLink(dl.url, dl.extract_code)">复制链接{{ dl.extract_code ? '与码' : '' }}</button>
              </view>
            </view>
          </view>

          <view class="order-footer">
            <text class="order-date">{{ order.created_at ? order.created_at.slice(0, 10) : '今日' }}</text>
            <button v-if="order.status === 'paid'" class="btn-get" @tap="goToDownload(order)">
              <text class="btn-get-icon">⚡</text>
              <text class="btn-get-text">进入网盘下载</text>
            </button>
            <view v-else class="review-tag-box">
              <text class="review-tag-text">⏳ 管理员核验到账后自动解锁</text>
            </view>
          </view>
        </view>
      </view>

      <view v-else-if="queried" class="empty-box">
        <text class="empty-icon">📂</text>
        <text class="empty-text">未查到该邮箱的购买记录，请确认邮箱是否正确</text>
        <text class="empty-sub">若付款后未收到，请联系 VIP 专属客服协助找回</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { orderAPI } from '@/api'
import { useAppStore } from '@/store/app'
import AmbientCanvas from '@/components/AmbientCanvas.vue'

const appStore = useAppStore()
const email = ref('')
const loading = ref(false)
const queried = ref(false)
const orders = ref([])

onLoad((options) => {
  if (options && options.email) {
    email.value = decodeURIComponent(options.email)
    handleQuery()
  }
})

function goBack() {
  const pages = typeof getCurrentPages === 'function' ? getCurrentPages() : []
  if (pages && pages.length > 1) {
    uni.navigateBack({
      delta: 1,
      fail: () => {
        uni.switchTab({ url: '/pages/index/index' })
      }
    })
    return
  }
  uni.switchTab({
    url: '/pages/index/index',
    fail: () => {
      uni.reLaunch({ url: '/pages/index/index' })
    }
  })
}

async function handleQuery() {
  if (!email.value || !email.value.includes('@')) {
    uni.showToast({ title: '请输入有效的邮箱', icon: 'none' })
    return
  }

  loading.value = true
  queried.value = true
  try {
    const res = await orderAPI.lookup(email.value.trim())
    orders.value = res.orders || res.data || []
    if (orders.value.length === 0) {
      uni.showToast({ title: '未查询到该邮箱的购买记录', icon: 'none' })
    }
  } catch (err) {
    orders.value = []
    uni.showToast({ title: err?.message || '未查询到订单记录', icon: 'none' })
  } finally {
    loading.value = false
  }
}

function goToDownload(order) {
  const pId = typeof order === 'object' ? order.product_id : order
  const oId = typeof order === 'object' ? (order.id || order.order_id) : ''
  const oEmail = (typeof order === 'object' ? order.user_email : '') || email.value

  if (oId) uni.setStorageSync('last_paid_order_id', oId)
  if (oEmail) uni.setStorageSync('last_buyer_email', oEmail)

  uni.navigateTo({
    url: `/pages/download/index?id=${pId}&order_id=${encodeURIComponent(oId || '')}&email=${encodeURIComponent(oEmail || '')}`
  })
}

function copyDriveLink(url, code) {
  if (!url) return
  let text = url
  if (code) {
    text = `链接: ${url} 提取码: ${code}`
  }
  uni.setClipboardData({
    data: text,
    success: () => {
      uni.showToast({ title: '网盘口令已复制！', icon: 'success' })
    }
  })
}

function openExternalUrl(url) {
  if (!url) return
  // #ifdef H5
  if (typeof window !== 'undefined') {
    window.open(url, '_blank')
    return
  }
  // #endif
  copyDriveLink(url)
}
</script>

<style scoped>
.lookup-container {
  min-height: 100vh;
  position: relative;
  box-sizing: border-box;
  transition: all 0.25s ease;
}

.lookup-content {
  position: relative;
  z-index: 2;
  max-width: 960rpx;
  margin: 0 auto;
  padding: calc(36rpx + env(safe-area-inset-top)) 36rpx 60rpx;
  box-sizing: border-box;
}

.top-nav {
  margin-bottom: 24rpx;
}

.back-pill {
  display: inline-flex;
  align-items: center;
  gap: 8rpx;
  padding: 10rpx 24rpx;
  border-radius: 999rpx;
  cursor: pointer;
}

.dark .back-pill {
  background: rgba(255, 255, 255, 0.08);
  border: 1rpx solid rgba(255, 255, 255, 0.12);
  color: #ffffff;
}

.light .back-pill {
  background: #ffffff;
  border: 1rpx solid rgba(0, 0, 0, 0.08);
  color: #0f172a;
}

.back-arrow {
  font-size: 24rpx;
}

.back-text {
  font-size: 22rpx;
  font-weight: 600;
}

.lookup-header {
  margin-bottom: 36rpx;
}

.header-badge {
  display: inline-flex;
  align-items: center;
  gap: 8rpx;
  padding: 6rpx 20rpx;
  border-radius: 999rpx;
  background: rgba(59, 130, 246, 0.14);
  border: 1rpx solid rgba(59, 130, 246, 0.3);
  margin-bottom: 16rpx;
}

.badge-icon {
  font-size: 22rpx;
}

.badge-text {
  font-size: 20rpx;
  font-weight: 700;
  color: #3b82f6;
}

.lookup-title {
  font-size: 36rpx;
  font-weight: 800;
  display: block;
  line-height: 1.35;
}

.dark .lookup-title {
  color: #ffffff;
}

.light .lookup-title {
  color: #0f172a;
}

.lookup-sub {
  font-size: 24rpx;
  line-height: 1.5;
  margin-top: 10rpx;
  display: block;
}

.dark .lookup-sub {
  color: rgba(255, 255, 255, 0.55);
}

.light .lookup-sub {
  color: #64748b;
}

.form-box {
  display: flex;
  gap: 16rpx;
  margin-bottom: 40rpx;
}

.email-input {
  flex: 1;
  height: 88rpx;
  border-radius: 20rpx;
  padding: 0 28rpx;
  font-size: 26rpx;
  box-sizing: border-box;
}

.dark .email-input {
  background: rgba(255, 255, 255, 0.06);
  border: 1rpx solid rgba(255, 255, 255, 0.15);
  color: #ffffff;
}

.light .email-input {
  background: #ffffff;
  border: 1rpx solid rgba(0, 0, 0, 0.12);
  color: #0f172a;
}

.btn-query {
  height: 88rpx;
  padding: 0 40rpx;
  background: linear-gradient(135deg, #2563eb, #1d4ed8);
  color: #ffffff;
  font-size: 26rpx;
  font-weight: 700;
  border-radius: 20rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 6rpx 20rpx rgba(37, 99, 235, 0.3);
  cursor: pointer;
}

.results-title {
  font-size: 26rpx;
  font-weight: 700;
  margin-bottom: 24rpx;
  display: block;
}

.dark .results-title {
  color: rgba(255, 255, 255, 0.75);
}

.light .results-title {
  color: #334155;
}

.order-card {
  border-radius: 24rpx;
  padding: 28rpx 32rpx;
  margin-bottom: 24rpx;
  box-sizing: border-box;
}

.dark .order-card {
  background: rgba(255, 255, 255, 0.05);
  border: 1rpx solid rgba(255, 255, 255, 0.1);
}

.light .order-card {
  background: #ffffff;
  border: 1rpx solid rgba(0, 0, 0, 0.06);
  box-shadow: 0 6rpx 24rpx rgba(0, 0, 0, 0.04);
}

.order-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 14rpx;
}

.order-no {
  font-size: 20rpx;
}

.dark .order-no {
  color: rgba(255, 255, 255, 0.45);
}

.light .order-no {
  color: #94a3b8;
}

.status-pill {
  display: inline-flex;
  align-items: center;
  gap: 6rpx;
  background: rgba(16, 185, 129, 0.15);
  padding: 4rpx 14rpx;
  border-radius: 999rpx;
}

.status-dot {
  font-size: 16rpx;
  color: #10b981;
}

.status-text {
  font-size: 20rpx;
  font-weight: 700;
  color: #10b981;
}

.order-name {
  font-size: 30rpx;
  font-weight: 800;
  display: block;
  margin-bottom: 20rpx;
}

.dark .order-name {
  color: #ffffff;
}

.light .order-name {
  color: #0f172a;
}

.order-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.order-date {
  font-size: 22rpx;
}

.dark .order-date {
  color: rgba(255, 255, 255, 0.4);
}

.light .order-date {
  color: #94a3b8;
}

.btn-get {
  height: 64rpx;
  padding: 0 28rpx;
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.25), rgba(37, 99, 235, 0.15));
  border: 1rpx solid rgba(59, 130, 246, 0.45);
  border-radius: 999rpx;
  display: flex;
  align-items: center;
  gap: 8rpx;
  cursor: pointer;
}

.btn-get-icon {
  font-size: 20rpx;
  color: #3b82f6;
}

.btn-get-text {
  font-size: 22rpx;
  font-weight: 700;
  color: #3b82f6;
}

.empty-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 80rpx 0;
  text-align: center;
}

.empty-icon {
  font-size: 72rpx;
  margin-bottom: 16rpx;
}

.empty-text {
  font-size: 26rpx;
  font-weight: 600;
  margin-bottom: 8rpx;
}

.dark .empty-text {
  color: rgba(255, 255, 255, 0.7);
}

.light .empty-text {
  color: #334155;
}

.empty-sub {
  font-size: 22rpx;
}

.dark .empty-sub {
  color: rgba(255, 255, 255, 0.4);
}

.light .empty-sub {
  color: #94a3b8;
}

.status-pill-review {
  background: rgba(245, 158, 11, 0.15) !important;
  color: #d97706 !important;
  border: 1px solid rgba(245, 158, 11, 0.3) !important;
}

.review-tag-box {
  padding: 6px 12px;
  background: rgba(245, 158, 11, 0.1);
  border-radius: 8px;
  border: 1px dashed rgba(245, 158, 11, 0.35);
}

.review-tag-text {
  font-size: 22rpx;
  font-weight: 700;
  color: #d97706;
}

/* Direct drive items in order card */
.direct-drives-list {
  display: flex;
  flex-direction: column;
  gap: 12rpx;
  margin-bottom: 20rpx;
  padding: 16rpx 20rpx;
  border-radius: 16rpx;
}
.dark .direct-drives-list {
  background: rgba(255, 255, 255, 0.04);
}
.light .direct-drives-list {
  background: #f8fafc;
}
.direct-drive-item {
  display: flex;
  flex-direction: column;
  gap: 8rpx;
  padding: 10rpx 0;
  border-bottom: 1rpx dashed rgba(0, 0, 0, 0.08);
}
.dark .direct-drive-item {
  border-bottom-color: rgba(255, 255, 255, 0.08);
}
.direct-drive-item:last-child {
  border-bottom: none;
}
.direct-drive-info {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.direct-drive-plat {
  font-size: 22rpx;
  font-weight: 700;
}
.dark .direct-drive-plat { color: #38bdf8; }
.light .direct-drive-plat { color: #0284c7; }
.direct-drive-code {
  font-size: 20rpx;
  font-weight: 700;
  color: #d97706;
  background: rgba(245, 158, 11, 0.15);
  padding: 2rpx 10rpx;
  border-radius: 6rpx;
}
.direct-drive-link-str {
  font-size: 22rpx;
  line-height: 1.4;
  word-break: break-all;
  display: flex;
  align-items: center;
  gap: 6rpx;
  text-decoration: underline;
  cursor: pointer;
  font-family: monospace;
}
.dark .direct-drive-link-str { color: #38bdf8; }
.light .direct-drive-link-str { color: #0284c7; }
.link-str-val { flex: 1; word-break: break-all; }
.link-str-icon { font-size: 22rpx; font-weight: 900; }
.direct-drive-btns {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12rpx;
  margin-top: 4rpx;
}
.btn-open-sm {
  font-size: 20rpx;
  font-weight: 700;
  padding: 6rpx 18rpx;
  border-radius: 999rpx;
  background: rgba(16, 185, 129, 0.15);
  color: #059669;
  border: none;
  cursor: pointer;
}
.btn-copy-sm {
  font-size: 20rpx;
  font-weight: 700;
  padding: 6rpx 18rpx;
  border-radius: 999rpx;
  background: rgba(59, 130, 246, 0.15);
  color: #2563eb;
  border: none;
  cursor: pointer;
}
</style>
