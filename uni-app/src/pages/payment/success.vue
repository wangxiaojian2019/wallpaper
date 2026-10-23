<template>
  <view class="payment-result-container" :class="[appStore.theme]">
    <!-- Ambient Canvas Background -->
    <AmbientCanvas />

    <!-- Top Sticky Frosted Glass Header -->
    <AppHeader />

    <view class="result-content">
      <!-- Loading State -->
      <view v-if="loading" class="result-card loading-card">
        <view class="spinner"></view>
        <text class="state-title">正在确认支付订单...</text>
        <text class="state-sub">正在与 PayPal 官方网关核验交易凭证与扣款状态，请稍候</text>
      </view>

      <!-- Success State -->
      <view v-else-if="success" class="result-card success-card">
        <view class="icon-circle icon-success">
          <text class="success-glyph">✓</text>
        </view>

        <text class="state-title">{{ isMembership ? '👑 VIP 特权开通成功！' : '🎉 恭喜！壁纸已解锁' }}</text>
        <text class="state-sub">
          {{ isMembership ? '感谢您的支持！您已获得全站数千款 4K/8K 壁纸与动态视频免单畅享特权。' : '订单已支付确认，高清无损原画网盘下载链接与提取码已就绪。' }}
        </text>

        <!-- Order Summary Pill -->
        <view class="order-summary-box">
          <view class="summary-line">
            <text class="s-label">交易单号</text>
            <text class="s-val font-mono">{{ orderInfo?.id || orderId }}</text>
          </view>
          <view class="summary-line">
            <text class="s-label">结算项目</text>
            <text class="s-val">{{ isMembership ? '全场 VIP 尊享特权会员' : (orderInfo?.product_title || '4K 精选超清壁纸') }}</text>
          </view>
          <view class="summary-line">
            <text class="s-label">支付金额</text>
            <text class="s-val s-gold font-bold">${{ Number(orderInfo?.amount || 0).toFixed(2) }} {{ orderInfo?.currency || 'USD' }}</text>
          </view>
          <view v-if="isMembership && userStore.userInfo?.membership_expire" class="summary-line">
            <text class="s-label">特权有效期至</text>
            <text class="s-val s-green">{{ formatDate(userStore.userInfo.membership_expire) }}</text>
          </view>
        </view>

        <!-- If Single Wallpaper with Download Links -->
        <view v-if="!isMembership && downloadLinks.length > 0" class="downloads-box">
          <text class="dl-box-title">📥 原画网盘高速直连</text>
          <view v-for="(dl, idx) in downloadLinks" :key="idx" class="dl-item-row">
            <view class="dl-item-left">
              <text class="dl-platform-badge">{{ getPlatformLabel(dl.platform) }}</text>
              <text class="dl-code" v-if="dl.extract_code">提取码: {{ dl.extract_code }}</text>
            </view>
            <button class="btn-copy-link" @tap="copyLink(dl.url, dl.extract_code)">
              复制口令直达
            </button>
          </view>
        </view>

        <!-- Action Buttons -->
        <view class="action-buttons-group">
          <button v-if="!isMembership && orderInfo?.product_id" class="btn-action btn-gold" @click="goToDownload" @tap="goToDownload">
            <text class="btn-icon">📥</text>
            <text class="btn-text">前往专属下载页面</text>
          </button>
          <button v-else class="btn-action btn-gold" @click="goHome" @tap="goHome">
            <text class="btn-icon">✦</text>
            <text class="btn-text">立即探索全场 4K 画廊</text>
          </button>

          <button class="btn-action btn-secondary" @click="goToMy" @tap="goToMy">
            <text class="btn-icon">👤</text>
            <text class="btn-text">前往个人中心</text>
          </button>
        </view>
      </view>

      <!-- Failure or Pending State -->
      <view v-else class="result-card error-card">
        <view class="icon-circle icon-error">
          <text class="error-glyph">✕</text>
        </view>

        <text class="state-title">订单支付尚未完成</text>
        <text class="state-sub">{{ errorMessage || 'PayPal 尚未返回扣款成功凭证，或您取消了付款' }}</text>

        <view class="action-buttons-group">
          <button class="btn-action btn-gold" @click="retryCapture" @tap="retryCapture">
            <text class="btn-icon">🔄</text>
            <text class="btn-text">重新核验支付状态</text>
          </button>
          <button class="btn-action btn-secondary" @click="goHome" @tap="goHome">
            <text class="btn-icon">🏠</text>
            <text class="btn-text">返回画廊首页</text>
          </button>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { paypalAPI, productAPI } from '@/api'
import { useUserStore } from '@/store/user'
import { useAppStore } from '@/store/app'
import AmbientCanvas from '@/components/AmbientCanvas.vue'
import AppHeader from '@/components/AppHeader.vue'

const userStore = useUserStore()
const appStore = useAppStore()

const loading = ref(true)
const success = ref(false)
const errorMessage = ref('')
const orderId = ref('')
const orderInfo = ref(null)
const downloadLinks = ref([])

const isMembership = computed(() => {
  return orderInfo.value?.order_type === 'membership' || userStore.isMember
})

function formatDate(isoStr) {
  if (!isoStr) return '永久有效'
  try {
    const d = new Date(isoStr)
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
  } catch {
    return isoStr
  }
}

function getPlatformLabel(p) {
  const map = {
    quark: '夸克网盘',
    baidu: '百度网盘',
    google: 'Google Drive'
  }
  return map[p] || '网盘直连'
}

function copyLink(url, code) {
  const fullText = code ? `${url} 提取码: ${code}` : url
  uni.setClipboardData({
    data: fullText,
    success: () => {
      uni.showToast({ title: '已复制链接与提取码', icon: 'success' })
    }
  })
}

function goHome() {
  // #ifdef H5
  if (typeof window !== 'undefined' && window.location) {
    window.location.hash = '#/pages/index/index'
    return
  }
  // #endif
  uni.switchTab({
    url: '/pages/index/index',
    fail: () => {
      uni.reLaunch({ url: '/pages/index/index' })
    }
  })
}

function goToMy() {
  // #ifdef H5
  if (typeof window !== 'undefined' && window.location) {
    window.location.hash = '#/pages/my/index'
    return
  }
  // #endif
  uni.switchTab({
    url: '/pages/my/index',
    fail: () => {
      uni.reLaunch({ url: '/pages/my/index' })
    }
  })
}

function goToDownload() {
  if (orderInfo.value?.product_id) {
    const oId = orderInfo.value.id || orderInfo.value.order_id || orderId.value || ''
    const oEmail = orderInfo.value.user_email || ''
    if (oId) uni.setStorageSync('last_paid_order_id', oId)
    if (oEmail) uni.setStorageSync('last_buyer_email', oEmail)
    uni.navigateTo({
      url: `/pages/download/index?id=${orderInfo.value.product_id}&title=${encodeURIComponent(orderInfo.value.product_title || '')}&order_id=${encodeURIComponent(oId)}&email=${encodeURIComponent(oEmail)}`
    })
  }
}

async function doCapture(idOrToken) {
  loading.value = true
  errorMessage.value = ''
  try {
    const res = await paypalAPI.captureOrder(idOrToken)
    if (res.status === 'paid' || res.status === 'COMPLETED' || res.order) {
      success.value = true
      orderInfo.value = res.order || res
      if (res.token) {
        uni.setStorageSync('token', res.token)
      }
      if (res.user) {
        userStore.userInfo = {
          ...(userStore.userInfo || {}),
          ...res.user,
          is_member: Boolean(res.user.membership_type && res.user.membership_type !== 'none')
        }
        uni.setStorageSync('user_info', JSON.stringify(userStore.userInfo))
      }
      if (orderInfo.value?.product_id) {
        try {
          const oId = orderInfo.value.id || orderInfo.value.order_id || orderId.value || ''
          const oEmail = orderInfo.value.user_email || ''
          const dlRes = await productAPI.getDownloads(orderInfo.value.product_id, { order_id: oId, email: oEmail })
          downloadLinks.value = dlRes.links || dlRes.download_links || dlRes.downloads || []
        } catch (e) {}
      }
    } else {
      errorMessage.value = res.message || 'PayPal 订单付款尚未确认'
    }
  } catch (err) {
    console.error('Capture error:', err)
    errorMessage.value = err.message || 'PayPal 订单验证异常'
  } finally {
    loading.value = false
  }
}

function retryCapture() {
  if (orderId.value) {
    doCapture(orderId.value)
  }
}

onLoad((options) => {
  // Extract order_id or token from options or URL search
  let targetId = options?.order_id || options?.token || ''
  // #ifdef H5
  if (!targetId && typeof window !== 'undefined' && window.location) {
    const params = new URLSearchParams(window.location.search || window.location.hash.split('?')[1] || '')
    targetId = params.get('order_id') || params.get('token') || ''
  }
  // #endif

  if (!targetId) {
    targetId = uni.getStorageSync('pending_paypal_order_id') || ''
  }

  orderId.value = targetId
  if (targetId) {
    doCapture(targetId)
  } else {
    loading.value = false
    errorMessage.value = '未找到订单号或交易凭据'
  }
})
</script>

<style scoped>
.payment-result-container {
  min-height: 100vh;
  position: relative;
  box-sizing: border-box;
  padding-bottom: 60rpx;
}

.result-content {
  position: relative;
  z-index: 2;
  max-width: 800rpx;
  margin: 40rpx auto 0;
  padding: 0 32rpx;
  box-sizing: border-box;
}

.result-card {
  border-radius: 36rpx;
  padding: 48rpx 36rpx;
  text-align: center;
  box-sizing: border-box;
  backdrop-filter: blur(28px) saturate(200%);
  -webkit-backdrop-filter: blur(28px) saturate(200%);
  animation: cardFadeIn 0.35s ease;
}

@keyframes cardFadeIn {
  from { opacity: 0; transform: translateY(16rpx); }
  to { opacity: 1; transform: translateY(0); }
}

.dark .result-card {
  background: rgba(18, 20, 32, 0.88);
  border: 1rpx solid rgba(255, 255, 255, 0.16);
  box-shadow: 0 16rpx 48rpx rgba(0, 0, 0, 0.55);
}

.light .result-card {
  background: rgba(255, 255, 255, 0.95);
  border: 1rpx solid rgba(0, 0, 0, 0.08);
  box-shadow: 0 16rpx 48rpx rgba(0, 0, 0, 0.08);
}

.spinner {
  width: 72rpx;
  height: 72rpx;
  border: 6rpx solid rgba(245, 158, 11, 0.2);
  border-top-color: #f59e0b;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin: 20rpx auto 30rpx;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.icon-circle {
  width: 120rpx;
  height: 120rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 28rpx;
}

.icon-success {
  background: rgba(16, 185, 129, 0.16);
  border: 2rpx solid rgba(16, 185, 129, 0.4);
}

.success-glyph {
  font-size: 56rpx;
  color: #10b981;
  font-weight: 800;
}

.icon-error {
  background: rgba(239, 68, 68, 0.16);
  border: 2rpx solid rgba(239, 68, 68, 0.4);
}

.error-glyph {
  font-size: 56rpx;
  color: #ef4444;
  font-weight: 800;
}

.state-title {
  font-size: 38rpx;
  font-weight: 900;
  display: block;
  letter-spacing: -0.5rpx;
  margin-bottom: 14rpx;
}

.dark .state-title {
  color: #ffffff;
}

.light .state-title {
  color: #0f172a;
}

.state-sub {
  font-size: 24rpx;
  line-height: 1.6;
  display: block;
  margin-bottom: 36rpx;
  padding: 0 20rpx;
}

.dark .state-sub {
  color: rgba(255, 255, 255, 0.65);
}

.light .state-sub {
  color: #64748b;
}

.order-summary-box {
  border-radius: 24rpx;
  padding: 24rpx 28rpx;
  margin-bottom: 36rpx;
  text-align: left;
}

.dark .order-summary-box {
  background: rgba(255, 255, 255, 0.04);
  border: 1rpx solid rgba(255, 255, 255, 0.08);
}

.light .order-summary-box {
  background: rgba(0, 0, 0, 0.03);
  border: 1rpx solid rgba(0, 0, 0, 0.06);
}

.summary-line {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10rpx 0;
  border-bottom: 1rpx dashed rgba(255, 255, 255, 0.08);
}

.light .summary-line {
  border-bottom-color: rgba(0, 0, 0, 0.06);
}

.summary-line:last-child {
  border-bottom: none;
}

.s-label {
  font-size: 22rpx;
  color: #94a3b8;
}

.s-val {
  font-size: 24rpx;
  font-weight: 600;
}

.dark .s-val {
  color: #ffffff;
}

.light .s-val {
  color: #0f172a;
}

.s-gold {
  color: #f59e0b !important;
}

.s-green {
  color: #10b981 !important;
}

.downloads-box {
  border-radius: 24rpx;
  padding: 24rpx 28rpx;
  margin-bottom: 36rpx;
  text-align: left;
}

.dark .downloads-box {
  background: rgba(59, 130, 246, 0.08);
  border: 1rpx solid rgba(59, 130, 246, 0.2);
}

.light .downloads-box {
  background: rgba(59, 130, 246, 0.06);
  border: 1rpx solid rgba(59, 130, 246, 0.15);
}

.dl-box-title {
  font-size: 24rpx;
  font-weight: 700;
  color: #3b82f6;
  display: block;
  margin-bottom: 16rpx;
}

.dl-item-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12rpx 0;
}

.dl-platform-badge {
  font-size: 22rpx;
  font-weight: 700;
  margin-right: 12rpx;
}

.dark .dl-platform-badge {
  color: #ffffff;
}

.light .dl-platform-badge {
  color: #0f172a;
}

.dl-code {
  font-size: 20rpx;
  color: #f59e0b;
}

.btn-copy-link {
  height: 48rpx;
  padding: 0 20rpx;
  border-radius: 999rpx;
  background: #3b82f6;
  color: #ffffff;
  font-size: 20rpx;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
}

.action-buttons-group {
  display: flex;
  flex-direction: column;
  gap: 18rpx;
}

.btn-action {
  width: 100%;
  height: 88rpx;
  border-radius: 999rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12rpx;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-action:active {
  transform: scale(0.98);
}

.btn-gold {
  background: linear-gradient(135deg, #f59e0b, #d97706);
  box-shadow: 0 8rpx 28rpx rgba(245, 158, 11, 0.35);
}

.btn-gold .btn-text {
  color: #000000;
  font-weight: 800;
  font-size: 28rpx;
}

.btn-secondary {
  border: 1rpx solid rgba(255, 255, 255, 0.16);
  background: rgba(255, 255, 255, 0.06);
}

.dark .btn-secondary .btn-text {
  color: #ffffff;
}

.light .btn-secondary {
  border-color: rgba(0, 0, 0, 0.1);
  background: rgba(0, 0, 0, 0.04);
}

.light .btn-secondary .btn-text {
  color: #0f172a;
}

.btn-icon {
  font-size: 28rpx;
}
</style>
