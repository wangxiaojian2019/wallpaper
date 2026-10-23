<template>
  <view class="download-container" :class="[appStore.theme]">
    <AmbientCanvas />

    <view class="download-content">
      <!-- Forbidden Notice Card (VIP Level Mismatch / Purchase Required) -->
      <view v-if="isForbidden" class="forbidden-card">
        <view class="forbidden-icon">👑</view>
        <text class="forbidden-title">需要升级会员等级</text>
        <text class="forbidden-msg">{{ forbiddenMessage }}</text>
        <button class="btn-upgrade-now" @tap="goToVip">
          立即升级更高等级 VIP
        </button>
        <button class="btn-back-detail" @tap="goBack">
          返回上一页
        </button>
      </view>

      <!-- Success Content (Only when authorized) -->
      <template v-else>
        <!-- Success Banner -->
        <view class="success-banner">
          <view class="success-icon-circle">
            <text class="success-check">✓</text>
          </view>
          <text class="banner-title">恭喜，已成功解锁 4K 原画！</text>
          <text class="banner-sub">请选择下方官方网盘通道获取无损原图与动态视频</text>
        </view>

        <!-- Cloud Drive Cards (Dynamically rendered based on admin configuration) -->
        <view v-if="hasAnyLink" class="drives-list">
          <!-- Quark Drive -->
          <view v-if="downloadLinks.quark_url" class="drive-card">
            <view class="drive-header">
              <view class="drive-brand">
                <view class="drive-icon-wrap quark-icon">⚡</view>
                <view class="drive-text-col">
                  <text class="drive-name">夸克网盘 (极速推荐 · 不限速)</text>
                  <text class="drive-desc">手机可直接无损保存原画至系统相册</text>
                </view>
              </view>
              <text class="drive-speed-badge">推荐直连</text>
            </view>

            <!-- URL Hyperlink Display Box -->
            <view class="drive-url-box">
              <view class="url-header-row">
                <text class="url-label">🔗 夸克网盘下载地址：</text>
                <view v-if="downloadLinks.quark_code" class="code-pill">
                  <text class="code-pill-label">提取码：</text>
                  <text class="code-pill-val">{{ downloadLinks.quark_code }}</text>
                </view>
              </view>
              <a :href="downloadLinks.quark_url" target="_blank" rel="noopener noreferrer" class="drive-url-link" @tap.stop>
                <text class="url-text-inner">{{ downloadLinks.quark_url }}</text>
                <text class="url-open-icon">↗</text>
              </a>
            </view>

            <view class="drive-action-row">
              <button class="btn-open-drive" @tap="openDriveUrl(downloadLinks.quark_url)">
                <text class="btn-icon">🌐</text>
                <text class="btn-txt">直接打开网盘 ↗</text>
              </button>
              <button class="btn-copy-drive" @tap="copyLink(downloadLinks.quark_url, downloadLinks.quark_code)">
                <text class="btn-icon">📋</text>
                <text class="btn-txt">复制链接{{ downloadLinks.quark_code ? '与提取码' : '' }}</text>
              </button>
            </view>
          </view>

          <!-- Baidu Drive -->
          <view v-if="downloadLinks.baidu_url" class="drive-card">
            <view class="drive-header">
              <view class="drive-brand">
                <view class="drive-icon-wrap baidu-icon">📦</view>
                <view class="drive-text-col">
                  <text class="drive-name">百度网盘</text>
                  <text class="drive-desc">长期有效 · 支持多端电脑手机同步转存</text>
                </view>
              </view>
              <text class="drive-speed-badge">全平台</text>
            </view>

            <!-- URL Hyperlink Display Box -->
            <view class="drive-url-box">
              <view class="url-header-row">
                <text class="url-label">🔗 百度网盘下载地址：</text>
                <view v-if="downloadLinks.baidu_code" class="code-pill">
                  <text class="code-pill-label">提取码：</text>
                  <text class="code-pill-val">{{ downloadLinks.baidu_code }}</text>
                </view>
              </view>
              <a :href="downloadLinks.baidu_url" target="_blank" rel="noopener noreferrer" class="drive-url-link" @tap.stop>
                <text class="url-text-inner">{{ downloadLinks.baidu_url }}</text>
                <text class="url-open-icon">↗</text>
              </a>
            </view>

            <view class="drive-action-row">
              <button class="btn-open-drive" @tap="openDriveUrl(downloadLinks.baidu_url)">
                <text class="btn-icon">🌐</text>
                <text class="btn-txt">直接打开网盘 ↗</text>
              </button>
              <button class="btn-copy-drive" @tap="copyLink(downloadLinks.baidu_url, downloadLinks.baidu_code)">
                <text class="btn-icon">📋</text>
                <text class="btn-txt">复制链接{{ downloadLinks.baidu_code ? '与提取码' : '' }}</text>
              </button>
            </view>
          </view>

          <!-- Google Drive (Overseas) -->
          <view v-if="downloadLinks.google_url" class="drive-card">
            <view class="drive-header">
              <view class="drive-brand">
                <view class="drive-icon-wrap google-icon">🌐</view>
                <view class="drive-text-col">
                  <text class="drive-name">Google Drive</text>
                  <text class="drive-desc">Direct high-speed worldwide download</text>
                </view>
              </view>
              <text class="drive-speed-badge">全球直连</text>
            </view>

            <!-- URL Hyperlink Display Box -->
            <view class="drive-url-box">
              <view class="url-header-row">
                <text class="url-label">🔗 Google Drive Download Link：</text>
              </view>
              <a :href="downloadLinks.google_url" target="_blank" rel="noopener noreferrer" class="drive-url-link" @tap.stop>
                <text class="url-text-inner">{{ downloadLinks.google_url }}</text>
                <text class="url-open-icon">↗</text>
              </a>
            </view>

            <view class="drive-action-row">
              <button class="btn-open-drive" @tap="openDriveUrl(downloadLinks.google_url)">
                <text class="btn-icon">🌐</text>
                <text class="btn-txt">Open Google Drive ↗</text>
              </button>
              <button class="btn-copy-drive" @tap="copyLink(downloadLinks.google_url, '')">
                <text class="btn-icon">📋</text>
                <text class="btn-txt">Copy Link</text>
              </button>
            </view>
          </view>

          <!-- Other Custom Platform Drives -->
          <view v-for="(item, idx) in customLinks" :key="idx" class="drive-card">
            <view class="drive-header">
              <view class="drive-brand">
                <view class="drive-icon-wrap custom-icon">💾</view>
                <view class="drive-text-col">
                  <text class="drive-name">{{ item.platform || '网盘高速下载' }}</text>
                  <text class="drive-desc">官方原画资源通道</text>
                </view>
              </view>
            </view>

            <!-- URL Hyperlink Display Box -->
            <view class="drive-url-box">
              <view class="url-header-row">
                <text class="url-label">🔗 网盘下载地址：</text>
                <view v-if="item.extract_code" class="code-pill">
                  <text class="code-pill-label">提取码：</text>
                  <text class="code-pill-val">{{ item.extract_code }}</text>
                </view>
              </view>
              <a :href="item.url" target="_blank" rel="noopener noreferrer" class="drive-url-link" @tap.stop>
                <text class="url-text-inner">{{ item.url }}</text>
                <text class="url-open-icon">↗</text>
              </a>
            </view>

            <view class="drive-action-row">
              <button class="btn-open-drive" @tap="openDriveUrl(item.url)">
                <text class="btn-icon">🌐</text>
                <text class="btn-txt">直接打开网盘 ↗</text>
              </button>
              <button class="btn-copy-drive" @tap="copyLink(item.url, item.extract_code)">
                <text class="btn-icon">📋</text>
                <text class="btn-txt">复制链接{{ item.extract_code ? '与提取码' : '' }}</text>
              </button>
            </view>
          </view>
        </view>

        <!-- Empty Links Notice (When authorized but links are still being uploaded) -->
        <view v-else-if="!loading" class="empty-links-card">
          <text class="empty-icon">📂</text>
          <text class="empty-title">网盘资源链接同步中</text>
          <text class="empty-desc">管理员正在为您上传与同步高速网盘链接，请稍后刷新页面或凭订单号联系客服直接获取原画文件。</text>
        </view>
      </template>

      <!-- Post-purchase Retention & Upsell: VIP Difference Upgrade -->
      <view class="vip-upgrade-card">
        <view class="vip-up-content">
          <view class="vip-up-tag">
            <text class="vip-tag-text">👑 尊享限时升舱</text>
          </view>
          <text class="vip-up-title">觉得好看？全场数千款壁纸今天全部免单</text>
          <text class="vip-up-sub">本单已自动抵扣，升级月度/季度全场无限免单 VIP 即享 8K 原画</text>
        </view>
        <button class="vip-up-btn" @click="goToVip" @tap="goToVip">
          立即升级 VIP
        </button>
      </view>

      <!-- Save to Photos Guide Tip -->
      <view class="guide-tip-box">
        <text class="tip-title">📱 手机如何设为锁屏壁纸？</text>
        <text class="tip-desc">1. 点击上方「复制链接」后在浏览器打开，转存或下载原画 PNG 保存至手机相册；</text>
        <text class="tip-desc">2. 打开手机「设置」>「墙纸」>「添加新墙纸」，选取该照片并设置景深效果。</text>
      </view>

      <!-- Back to Gallery -->
      <view class="bottom-action">
        <button class="btn-home" @click="goHome" @tap="goHome">
          ← 返回画廊挑选更多壁纸
        </button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { productAPI } from '@/api'
import { useAppStore } from '@/store/app'
import AmbientCanvas from '@/components/AmbientCanvas.vue'

const appStore = useAppStore()
const productId = ref('')
const wallpaperTitle = ref('')
const orderId = ref('')
const userEmail = ref('')
const unlockCode = ref('')
const loading = ref(true)
const isForbidden = ref(false)
const forbiddenMessage = ref('')

const downloadLinks = ref({
  quark_url: '',
  quark_code: '',
  baidu_url: '',
  baidu_code: '',
  google_url: ''
})
const customLinks = ref([])

const hasAnyLink = computed(() => {
  return Boolean(
    downloadLinks.value.quark_url ||
    downloadLinks.value.baidu_url ||
    downloadLinks.value.google_url ||
    customLinks.value.length > 0
  )
})

async function fetchLinks(id) {
  loading.value = true
  isForbidden.value = false
  forbiddenMessage.value = ''

  // Reset links
  downloadLinks.value = {
    quark_url: '',
    quark_code: '',
    baidu_url: '',
    baidu_code: '',
    google_url: ''
  }
  customLinks.value = []

  try {
    const params = {}
    if (orderId.value) params.order_id = orderId.value
    if (userEmail.value) params.email = userEmail.value
    if (unlockCode.value) params.code = unlockCode.value

    const res = await productAPI.getDownloads(id, params)
    const list = res.links || res.download_links || res.downloads || res.data || []
    if (Array.isArray(list) && list.length > 0) {
      list.forEach(item => {
        if (!item || !item.url) return
        const plat = (item.platform || '').toLowerCase().trim()
        if (plat === 'quark' || plat === '夸克' || plat.includes('quark')) {
          downloadLinks.value.quark_url = item.url
          downloadLinks.value.quark_code = item.extract_code || ''
        } else if (plat === 'baidu' || plat === '百度' || plat.includes('baidu')) {
          downloadLinks.value.baidu_url = item.url
          downloadLinks.value.baidu_code = item.extract_code || ''
        } else if (plat === 'google' || plat === '谷歌' || plat.includes('google')) {
          downloadLinks.value.google_url = item.url
        } else {
          customLinks.value.push(item)
        }
      })
    }
  } catch (err) {
    if (err?.statusCode === 403 || err?.status === 403 || err?.data?.vip_tier_mismatch) {
      isForbidden.value = true
      forbiddenMessage.value = err?.data?.message || err?.message || '您当前会员等级暂未包含此壁纸的下载权益'
    } else {
      console.warn('Failed to load download links:', err)
      uni.showToast({
        title: err?.message || '获取下载链接失败，请核验订单',
        icon: 'none'
      })
    }
  } finally {
    loading.value = false
  }
}

function goBack() {
  uni.navigateBack({
    fail: () => {
      uni.reLaunch({ url: '/pages/index/index' })
    }
  })
}

function copyLink(url, code) {
  if (!url) {
    uni.showToast({ title: '暂未配置该网盘链接', icon: 'none' })
    return
  }
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

function openDriveUrl(url) {
  if (!url) {
    uni.showToast({ title: '暂未配置该网盘链接', icon: 'none' })
    return
  }
  // #ifdef H5
  if (typeof window !== 'undefined') {
    window.open(url, '_blank')
    return
  }
  // #endif

  uni.setClipboardData({
    data: url,
    success: () => {
      uni.showToast({ title: '网盘地址已复制，请在浏览器中打开', icon: 'none', duration: 2500 })
    }
  })
}

function goToVip() {
  // #ifdef H5
  if (typeof window !== 'undefined' && window.location) {
    window.location.hash = '#/pages/membership/index'
    return
  }
  // #endif
  uni.switchTab({
    url: '/pages/membership/index',
    fail: () => {
      uni.reLaunch({ url: '/pages/membership/index' })
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

onLoad((options) => {
  productId.value = options?.id || ''
  wallpaperTitle.value = decodeURIComponent(options?.title || '壁纸')
  orderId.value = options?.order_id || uni.getStorageSync('last_paid_order_id') || ''
  userEmail.value = decodeURIComponent(options?.email || '') || uni.getStorageSync('last_buyer_email') || ''
  unlockCode.value = decodeURIComponent(options?.code || '')

  // #ifdef H5
  if (typeof window !== 'undefined' && window.location) {
    const searchStr = window.location.search || (window.location.hash.includes('?') ? window.location.hash.split('?')[1] : '')
    if (searchStr) {
      const urlParams = new URLSearchParams(searchStr)
      if (!orderId.value) orderId.value = urlParams.get('order_id') || ''
      if (!userEmail.value) userEmail.value = urlParams.get('email') || ''
      if (!unlockCode.value) unlockCode.value = urlParams.get('code') || ''
      if (!productId.value) productId.value = urlParams.get('id') || ''
    }
  }
  // #endif

  if (productId.value) {
    fetchLinks(productId.value)
  }
})
</script>

<style scoped>
.download-container {
  min-height: 100vh;
  position: relative;
  box-sizing: border-box;
  transition: all 0.25s ease;
}

.download-content {
  position: relative;
  z-index: 2;
  max-width: 960rpx;
  margin: 0 auto;
  padding: calc(40rpx + env(safe-area-inset-top)) 36rpx 60rpx;
  box-sizing: border-box;
}

/* Banner */
.success-banner {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  margin-bottom: 40rpx;
}

.success-icon-circle {
  width: 108rpx;
  height: 108rpx;
  border-radius: 50%;
  background: linear-gradient(135deg, #10b981, #059669);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 24rpx;
  box-shadow: 0 10rpx 32rpx rgba(16, 185, 129, 0.4);
}

.success-check {
  font-size: 56rpx;
  color: #ffffff;
  font-weight: 900;
}

.banner-title {
  font-size: 38rpx;
  font-weight: 800;
  display: block;
  line-height: 1.35;
}

.dark .banner-title {
  color: #ffffff;
}

.light .banner-title {
  color: #0f172a;
}

.banner-sub {
  font-size: 24rpx;
  margin-top: 10rpx;
  line-height: 1.5;
  display: block;
}

.dark .banner-sub {
  color: rgba(255, 255, 255, 0.55);
}

.light .banner-sub {
  color: #64748b;
}

/* Drive cards */
.drives-list {
  display: flex;
  flex-direction: column;
  gap: 24rpx;
  margin-bottom: 40rpx;
}

.drive-card {
  border-radius: 28rpx;
  padding: 28rpx 32rpx;
  box-sizing: border-box;
}

.dark .drive-card {
  background: rgba(255, 255, 255, 0.05);
  border: 1rpx solid rgba(255, 255, 255, 0.1);
}

.light .drive-card {
  background: #ffffff;
  border: 1rpx solid rgba(0, 0, 0, 0.06);
  box-shadow: 0 6rpx 24rpx rgba(0, 0, 0, 0.04);
}

.drive-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 22rpx;
}

.drive-brand {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.drive-icon-wrap {
  width: 64rpx;
  height: 64rpx;
  border-radius: 18rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32rpx;
}

.quark-icon { background: rgba(245, 158, 11, 0.16); color: #f59e0b; }
.baidu-icon { background: rgba(59, 130, 246, 0.16); color: #3b82f6; }
.google-icon { background: rgba(16, 185, 129, 0.16); color: #10b981; }

.drive-text-col {
  display: flex;
  flex-direction: column;
}

.drive-name {
  font-size: 28rpx;
  font-weight: 800;
}

.dark .drive-name {
  color: #ffffff;
}

.light .drive-name {
  color: #0f172a;
}

.drive-desc {
  font-size: 22rpx;
  margin-top: 4rpx;
}

.dark .drive-desc {
  color: rgba(255, 255, 255, 0.45);
}

.light .drive-desc {
  color: #64748b;
}

.drive-speed-badge {
  font-size: 20rpx;
  font-weight: 700;
  color: #10b981;
  background: rgba(16, 185, 129, 0.15);
  padding: 4rpx 14rpx;
  border-radius: 999rpx;
}

/* Direct URL Hyperlink Box */
.drive-url-box {
  margin: 16rpx 0 20rpx;
  padding: 16rpx 20rpx;
  border-radius: 18rpx;
  display: flex;
  flex-direction: column;
  gap: 10rpx;
  box-sizing: border-box;
}

.dark .drive-url-box {
  background: rgba(255, 255, 255, 0.04);
  border: 1rpx solid rgba(255, 255, 255, 0.08);
}

.light .drive-url-box {
  background: #f8fafc;
  border: 1rpx solid rgba(0, 0, 0, 0.06);
}

.url-header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12rpx;
}

.url-label {
  font-size: 22rpx;
  font-weight: 700;
}

.dark .url-label {
  color: rgba(255, 255, 255, 0.65);
}

.light .url-label {
  color: #475569;
}

.code-pill {
  display: inline-flex;
  align-items: center;
  padding: 4rpx 14rpx;
  border-radius: 8rpx;
  background: rgba(245, 158, 11, 0.15);
}

.code-pill-label {
  font-size: 20rpx;
  color: #d97706;
}

.code-pill-val {
  font-size: 22rpx;
  font-weight: 800;
  color: #d97706;
  letter-spacing: 1rpx;
}

.drive-url-link {
  font-size: 24rpx;
  line-height: 1.5;
  word-break: break-all;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12rpx;
  text-decoration: underline;
  text-underline-offset: 4rpx;
  cursor: pointer;
  font-family: -apple-system, BlinkMacSystemFont, 'SF Mono', 'Segoe UI Mono', monospace;
  padding: 8rpx 12rpx;
  border-radius: 10rpx;
  transition: all 0.2s ease;
}

.dark .drive-url-link {
  color: #38bdf8;
  background: rgba(56, 189, 248, 0.06);
}

.light .drive-url-link {
  color: #0284c7;
  background: rgba(2, 132, 199, 0.06);
}

.drive-url-link:hover {
  opacity: 0.85;
}

.url-text-inner {
  flex: 1;
  word-break: break-all;
}

.url-open-icon {
  font-size: 24rpx;
  font-weight: 900;
  flex-shrink: 0;
}

/* Action button row */
.drive-action-row {
  display: flex;
  align-items: center;
  gap: 16rpx;
  margin-top: 8rpx;
}

.btn-open-drive {
  flex: 1;
  height: 72rpx;
  border-radius: 999rpx;
  background: linear-gradient(135deg, #10b981, #059669);
  color: #ffffff;
  font-size: 24rpx;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8rpx;
  border: none;
  cursor: pointer;
  box-shadow: 0 6rpx 18rpx rgba(16, 185, 129, 0.25);
  transition: all 0.2s ease;
}

.btn-open-drive:active {
  transform: scale(0.98);
}

.btn-copy-drive {
  flex: 1;
  height: 72rpx;
  border-radius: 999rpx;
  background: linear-gradient(135deg, #2563eb, #1d4ed8);
  color: #ffffff;
  font-size: 24rpx;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8rpx;
  border: none;
  cursor: pointer;
  box-shadow: 0 6rpx 18rpx rgba(37, 99, 235, 0.25);
  transition: all 0.2s ease;
}

.btn-copy-drive:active {
  transform: scale(0.98);
}

.btn-icon {
  font-size: 22rpx;
}

.btn-txt {
  font-size: 24rpx;
}

.full-btn {
  width: 100%;
}

/* VIP Upgrade Card */
.vip-upgrade-card {
  border-radius: 28rpx;
  padding: 32rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 40rpx;
  gap: 16rpx;
}

.dark .vip-upgrade-card {
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.16), rgba(212, 175, 55, 0.06));
  border: 1rpx solid rgba(245, 158, 11, 0.35);
}

.light .vip-upgrade-card {
  background: linear-gradient(135deg, rgba(254, 243, 199, 0.8), rgba(255, 255, 255, 0.95));
  border: 1rpx solid rgba(217, 119, 6, 0.3);
  box-shadow: 0 6rpx 24rpx rgba(217, 119, 6, 0.1);
}

.vip-up-content {
  flex: 1;
}

.vip-up-tag {
  display: inline-block;
  background: #f59e0b;
  padding: 2rpx 14rpx;
  border-radius: 8rpx;
  margin-bottom: 10rpx;
}

.vip-tag-text {
  font-size: 18rpx;
  font-weight: 800;
  color: #000;
}

.vip-up-title {
  font-size: 28rpx;
  font-weight: 800;
  display: block;
}

.dark .vip-up-title {
  color: #fde68a;
}

.light .vip-up-title {
  color: #92400e;
}

.vip-up-sub {
  font-size: 22rpx;
  line-height: 1.45;
  margin-top: 8rpx;
  display: block;
}

.dark .vip-up-sub {
  color: rgba(255, 255, 255, 0.6);
}

.light .vip-up-sub {
  color: #78350f;
}

.vip-up-btn {
  height: 72rpx;
  padding: 0 32rpx;
  background: linear-gradient(135deg, #f59e0b, #d97706);
  color: #000;
  font-size: 24rpx;
  font-weight: 800;
  border-radius: 999rpx;
  white-space: nowrap;
  box-shadow: 0 6rpx 20rpx rgba(245, 158, 11, 0.35);
}

/* Guide tip */
.guide-tip-box {
  border-radius: 24rpx;
  padding: 28rpx 32rpx;
  margin-bottom: 40rpx;
}

.dark .guide-tip-box {
  background: rgba(255, 255, 255, 0.04);
  border: 1rpx solid rgba(255, 255, 255, 0.06);
}

.light .guide-tip-box {
  background: #ffffff;
  border: 1rpx solid rgba(0, 0, 0, 0.06);
}

.tip-title {
  font-size: 26rpx;
  font-weight: 700;
  margin-bottom: 12rpx;
  display: block;
}

.dark .tip-title {
  color: rgba(255, 255, 255, 0.85);
}

.light .tip-title {
  color: #0f172a;
}

.tip-desc {
  font-size: 22rpx;
  line-height: 1.6;
  display: block;
}

.dark .tip-desc {
  color: rgba(255, 255, 255, 0.5);
}

.light .tip-desc {
  color: #64748b;
}

/* Bottom Action */
.bottom-action {
  display: flex;
  justify-content: center;
}

.btn-home {
  width: 100%;
  height: 88rpx;
  border-radius: 999rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 26rpx;
  font-weight: 700;
  cursor: pointer;
}

.dark .btn-home {
  background: rgba(255, 255, 255, 0.06);
  border: 1rpx solid rgba(255, 255, 255, 0.14);
  color: #ffffff;
}

.light .btn-home {
  background: #ffffff;
  border: 1rpx solid rgba(0, 0, 0, 0.1);
  color: #0f172a;
}

/* Forbidden Card */
.forbidden-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding: 60rpx 40rpx;
  border-radius: 32rpx;
  margin-bottom: 40rpx;
}
.dark .forbidden-card {
  background: rgba(30, 20, 50, 0.7);
  border: 1rpx solid rgba(168, 85, 247, 0.3);
}
.light .forbidden-card {
  background: #ffffff;
  border: 1rpx solid rgba(168, 85, 247, 0.2);
  box-shadow: 0 10rpx 30rpx rgba(168, 85, 247, 0.08);
}
.forbidden-icon {
  font-size: 72rpx;
  margin-bottom: 20rpx;
}
.forbidden-title {
  font-size: 34rpx;
  font-weight: 800;
  margin-bottom: 12rpx;
}
.dark .forbidden-title { color: #ffffff; }
.light .forbidden-title { color: #0f172a; }
.forbidden-msg {
  font-size: 26rpx;
  line-height: 1.6;
  margin-bottom: 36rpx;
  max-width: 520rpx;
}
.dark .forbidden-msg { color: #c084fc; }
.light .forbidden-msg { color: #7c3aed; }
.btn-upgrade-now {
  width: 100%;
  max-width: 440rpx;
  height: 84rpx;
  border-radius: 999rpx;
  background: linear-gradient(135deg, #d4af37, #f59e0b);
  color: #ffffff;
  font-size: 28rpx;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 20rpx;
  border: none;
}
.btn-back-detail {
  width: 100%;
  max-width: 440rpx;
  height: 80rpx;
  border-radius: 999rpx;
  background: rgba(255, 255, 255, 0.08);
  border: 1rpx solid rgba(255, 255, 255, 0.15);
  font-size: 26rpx;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
}
.dark .btn-back-detail { color: #ffffff; }
.light .btn-back-detail { color: #0f172a; }

/* Empty Links Card */
.empty-links-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding: 60rpx 40rpx;
  border-radius: 28rpx;
  margin-bottom: 40rpx;
}

.dark .empty-links-card {
  background: rgba(255, 255, 255, 0.04);
  border: 1rpx dashed rgba(255, 255, 255, 0.15);
}

.light .empty-links-card {
  background: #ffffff;
  border: 1rpx dashed rgba(0, 0, 0, 0.15);
}

.empty-icon {
  font-size: 64rpx;
  margin-bottom: 16rpx;
}

.empty-title {
  font-size: 30rpx;
  font-weight: 700;
  margin-bottom: 10rpx;
}

.dark .empty-title { color: #ffffff; }
.light .empty-title { color: #0f172a; }

.empty-desc {
  font-size: 24rpx;
  line-height: 1.6;
  max-width: 520rpx;
}

.dark .empty-desc { color: rgba(255, 255, 255, 0.55); }
.light .empty-desc { color: #64748b; }

.custom-icon {
  background: rgba(147, 51, 234, 0.16);
  color: #a855f7;
}
</style>
