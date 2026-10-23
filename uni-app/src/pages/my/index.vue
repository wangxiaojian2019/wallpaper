<template>
  <view class="my-container" :class="[appStore.theme]">
    <!-- Dynamic Ambient Background -->
    <AmbientCanvas />

    <!-- Top Sticky Frosted Glass Header -->
    <AppHeader />

    <view class="my-content">
      <!-- Luxury User Profile Header Card -->
      <view class="user-hero-card">
        <view class="user-main-row">
          <view class="avatar-box" :class="{ 'avatar-vip': userStore.isMember }">
            <text class="avatar-char">{{ userAvatarChar }}</text>
            <view v-if="userStore.isMember" class="avatar-crown-badge">👑</view>
          </view>

          <view class="user-info-col">
            <view class="name-row">
              <text class="user-title">{{ userDisplayName }}</text>
              <view v-if="userStore.isMember" class="vip-glow-pill">
                <text class="vip-pill-icon">👑</text>
                <text class="vip-pill-text">VIP 会员</text>
              </view>
              <view v-else class="free-pill">
                <text class="free-pill-text">普通用户</text>
              </view>
            </view>
            <text class="user-subtext">
              {{ userStore.isMember ? '已尊享全场壁纸免单与 8K 原图特权' : (userStore.isLoggedIn ? '已登录 · 享有每日精选与限免壁纸' : '未登录 · 支持输入邮箱直接提取已购订单') }}
            </text>
          </view>
        </view>

        <!-- VIP Promo Banner Inside Card (for non-members) -->
        <view v-if="!userStore.isMember" class="vip-banner-inside" @click="goToVip" @tap="goToVip">
          <view class="banner-left">
            <text class="banner-crown">👑</text>
            <view class="banner-texts">
              <text class="banner-t1">升级至尊 VIP 会员</text>
              <text class="banner-t2">全场数千款 4K/8K 臻选原图免单畅下</text>
            </view>
          </view>
          <view class="banner-btn">
            <text class="banner-btn-text">立即开通</text>
            <text class="banner-btn-arrow">›</text>
          </view>
        </view>
      </view>

      <!-- Quick 3-Metric Stats Bar -->
      <view class="stats-row">
        <view class="stat-card" @click="goToLookup" @tap="goToLookup">
          <text class="stat-num">{{ orderCount }}</text>
          <text class="stat-label">已购壁纸</text>
        </view>
        <view class="stat-card" @click="goToVip" @tap="goToVip">
          <text class="stat-num" :class="{ 'gold-text': userStore.isMember }">{{ userStore.isMember ? 'VIP' : '开通' }}</text>
          <text class="stat-label">会员特权</text>
        </view>
        <view class="stat-card">
          <text class="stat-num">4K/8K</text>
          <text class="stat-label">交付画质</text>
        </view>
      </view>

      <!-- Menu Group 1: 订单与下载资产 -->
      <view class="menu-group">
        <view class="group-title-row">
          <text class="group-title">订单与交付资产</text>
        </view>
        <view class="menu-box">
          <view class="menu-item" @click="goToLookup" @tap="goToLookup">
            <view class="menu-left">
              <view class="item-icon-wrap icon-blue">
                <text class="item-icon">🔍</text>
              </view>
              <view class="item-texts">
                <text class="item-title">游客订单与提取码找回</text>
                <text class="item-desc">输入购买时填写的邮箱，一键找回已购壁纸与提取码</text>
              </view>
            </view>
            <text class="menu-arrow">›</text>
          </view>

          <view class="menu-item" @click="goToVip" @tap="goToVip">
            <view class="menu-left">
              <view class="item-icon-wrap icon-gold">
                <text class="item-icon">👑</text>
              </view>
              <view class="item-texts">
                <text class="item-title">开通 / 续费全场 VIP</text>
                <text class="item-desc">免单任下数千款 4K/8K 超清原图及动态视频</text>
              </view>
            </view>
            <text class="menu-arrow">›</text>
          </view>

          <!-- Redeem Code Entry -->
          <view class="menu-item" @click="showRedeemModal = true">
            <view class="menu-left">
              <view class="item-icon-wrap icon-emerald">
                <text class="item-icon">🎫</text>
              </view>
              <view class="item-texts">
                <text class="item-title">口令兑换与会员卡密</text>
                <text class="item-desc">输入礼品卡密或兑换口令，即刻激活 VIP 或领壁纸</text>
              </view>
            </view>
            <text class="menu-arrow">›</text>
          </view>

          <view class="menu-item" @tap="showHelpModal">
            <view class="menu-left">
              <view class="item-icon-wrap icon-amber">
                <text class="item-icon">💡</text>
              </view>
              <view class="item-texts">
                <text class="item-title">锁屏 / 桌面壁纸设置指南</text>
                <text class="item-desc">iOS 景深锁屏与 Mac/Android 壁纸配置技巧</text>
              </view>
            </view>
            <text class="menu-arrow">›</text>
          </view>
        </view>
      </view>

      <!-- Menu Group 2: 系统偏好与个性化 -->
      <view class="menu-group">
        <view class="group-title-row">
          <text class="group-title">个性化与偏好</text>
        </view>
        <view class="menu-box">
          <!-- Theme toggle -->
          <view class="menu-item" @tap="appStore.toggleTheme">
            <view class="menu-left">
              <view class="item-icon-wrap icon-purple">
                <text class="item-icon">{{ appStore.theme === 'dark' ? '🌙' : '☀️' }}</text>
              </view>
              <view class="item-texts">
                <text class="item-title">界面视觉外观</text>
                <text class="item-desc">在「深色星空」与「浅色雨露」模式之间一键切换</text>
              </view>
            </view>
            <view class="status-badge">
              <text class="status-badge-text">{{ appStore.theme === 'dark' ? '深色星空' : '浅色雨露' }}</text>
              <text class="menu-arrow">›</text>
            </view>
          </view>

          <!-- Language toggle -->
          <view class="menu-item" @tap="appStore.toggleLang">
            <view class="menu-left">
              <view class="item-icon-wrap icon-cyan">
                <text class="item-icon">🌐</text>
              </view>
              <view class="item-texts">
                <text class="item-title">语言切换 (Language)</text>
                <text class="item-desc">切换中文与英文界面显示</text>
              </view>
            </view>
            <view class="status-badge">
              <text class="status-badge-text">{{ appStore.lang === 'zh' ? '简体中文' : 'English' }}</text>
              <text class="menu-arrow">›</text>
            </view>
          </view>

          <!-- Cache clear -->
          <view class="menu-item" @tap="handleClearCache">
            <view class="menu-left">
              <view class="item-icon-wrap icon-gray">
                <text class="item-icon">🧹</text>
              </view>
              <view class="item-texts">
                <text class="item-title">清理本地缓存</text>
                <text class="item-desc">释放临时预览图片与离线数据</text>
              </view>
            </view>
            <view class="status-badge">
              <text class="status-badge-text">{{ cacheSize }}</text>
              <text class="menu-arrow">›</text>
            </view>
          </view>
        </view>
      </view>

      <!-- Menu Group 3: 服务与反馈 -->
      <view class="menu-group">
        <view class="group-title-row">
          <text class="group-title">服务支持与条款</text>
        </view>
        <view class="menu-box">
          <view class="menu-item" @tap="showSupportModal">
            <view class="menu-left">
              <view class="item-icon-wrap icon-green">
                <text class="item-icon">💬</text>
              </view>
              <view class="item-texts">
                <text class="item-title">VIP 专属求图与客服反馈</text>
                <text class="item-desc">找不到心仪尺寸？联系客服为您专属寻图</text>
              </view>
            </view>
            <text class="menu-arrow">›</text>
          </view>

          <view class="menu-item" @tap="showTermsModal">
            <view class="menu-left">
              <view class="item-icon-wrap icon-slate">
                <text class="item-icon">📜</text>
              </view>
              <view class="item-texts">
                <text class="item-title">用户服务协议与隐私条款</text>
                <text class="item-desc">关于版权、交付合规与隐私保护说明</text>
              </view>
            </view>
            <text class="menu-arrow">›</text>
          </view>
        </view>
      </view>

      <!-- Action Button: Login / Logout -->
      <view class="account-action-wrap">
        <button v-if="userStore.isLoggedIn" class="btn-account-logout" @tap="handleLogout">
          退出当前登录
        </button>
        <button v-else class="btn-account-login" @tap="showLoginPrompt">
          <text class="login-sparkle">✦</text>
          <text class="login-btn-text">登录账号 / 绑定邮箱凭据</text>
        </button>
      </view>

      <!-- App Branding Footer -->
      <view class="app-branding-footer">
        <text class="brand-line-1">WallArt 4K · 让每一块屏幕都有惊艳时刻</text>
        <text class="brand-line-2">Version 2.4.0 · Curated Art & High Definition</text>
      </view>
    </view>

    <!-- Redeem Code Modal -->
    <view v-if="showRedeemModal" class="custom-modal-overlay" @click.self="showRedeemModal = false">
      <view class="redeem-modal-card">
        <view class="redeem-modal-header">
          <view class="r-title-wrap">
            <text class="r-title-icon">🎫</text>
            <text class="r-title-text">{{ appStore.lang === 'zh' ? '口令兑换 / 会员卡密' : 'Redeem Code' }}</text>
          </view>
          <view class="r-close-btn" @click="showRedeemModal = false">✕</view>
        </view>
        <text class="r-desc-text">
          {{ appStore.lang === 'zh' ? '输入兑换口令即可激活 VIP 会员特权、单款壁纸解锁或限额免单礼包。' : 'Redeem VIP pass, wallpaper unlock vouchers, or promo tokens.' }}
        </text>

        <view class="r-input-group">
          <text class="r-input-label">{{ appStore.lang === 'zh' ? '兑换口令 / 卡密' : 'Redeem Code' }}</text>
          <input
            v-model="redeemCodeInput"
            type="text"
            :placeholder="appStore.lang === 'zh' ? '请输入兑换口令（如 VIP888 等）' : 'Enter redeem code'"
            class="r-input"
          />
        </view>

        <view v-if="!userStore.isLoggedIn" class="r-input-group">
          <text class="r-input-label">{{ appStore.lang === 'zh' ? '绑定账户邮箱 (可选)' : 'Account Email (Optional)' }}</text>
          <input
            v-model="redeemEmailInput"
            type="text"
            :placeholder="appStore.lang === 'zh' ? '请输入用于保存权益的邮箱' : 'Enter email to link'"
            class="r-input"
          />
        </view>

        <view class="r-btn-row">
          <button class="r-btn-cancel" @click="showRedeemModal = false">
            {{ appStore.lang === 'zh' ? '取消' : 'Cancel' }}
          </button>
          <button class="r-btn-confirm" :loading="redeeming" @click="submitRedeemCode">
            {{ appStore.lang === 'zh' ? '立即兑换' : 'Confirm' }}
          </button>
        </view>
      </view>
    </view>

    <!-- Auth Modal for Login and Registration -->
    <AuthModal
      :visible="showAuth"
      @close="showAuth = false"
      @success="handleAuthDone"
    />
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useUserStore } from '@/store/user'
import { useAppStore } from '@/store/app'
import request from '@/utils/request'
import AmbientCanvas from '@/components/AmbientCanvas.vue'
import AppHeader from '@/components/AppHeader.vue'
import AuthModal from '@/components/AuthModal.vue'

const userStore = useUserStore()
const appStore = useAppStore()
const showAuth = ref(false)
const cacheSize = ref('16.4 MB')
const orderCount = ref('0')

// Redeem code state
const showRedeemModal = ref(false)
const redeemCodeInput = ref('')
const redeemEmailInput = ref('')
const redeeming = ref(false)

async function submitRedeemCode() {
  const code = redeemCodeInput.value.trim()
  if (!code) {
    uni.showToast({ title: appStore.lang === 'zh' ? '请输入兑换口令' : 'Please enter redeem code', icon: 'none' })
    return
  }

  redeeming.value = true
  try {
    const res = await request.post('/redeem/apply', {
      code: code,
      user_email: redeemEmailInput.value.trim() || userStore.userInfo?.email || ''
    })

    uni.showToast({
      title: res.message || (appStore.lang === 'zh' ? '兑换成功！' : 'Redeemed successfully!'),
      icon: 'success'
    })

    if (res.token) {
      uni.setStorageSync('token', res.token)
    }
    if (res.user) {
      userStore.userInfo = {
        ...(userStore.userInfo || {}),
        ...res.user
      }
      uni.setStorageSync('user_info', JSON.stringify(userStore.userInfo))
    }

    showRedeemModal.value = false
    redeemCodeInput.value = ''
    redeemEmailInput.value = ''

    if (userStore.fetchUserInfo) {
      userStore.fetchUserInfo()
    }
  } catch (err) {
    uni.showToast({
      title: err.message || (appStore.lang === 'zh' ? '兑换失败，口令无效或已过期' : 'Invalid or expired code'),
      icon: 'none'
    })
  } finally {
    redeeming.value = false
  }
}

const userDisplayName = computed(() => {
  if (userStore.isLoggedIn && userStore.userInfo?.email) {
    return userStore.userInfo.email
  }
  return '游客探索者'
})

const userAvatarChar = computed(() => {
  if (userStore.isLoggedIn && userStore.userInfo?.email) {
    return userStore.userInfo.email.slice(0, 1).toUpperCase()
  }
  return '👤'
})

function goToLookup() {
  uni.navigateTo({ url: '/pages/lookup/index' })
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

function showHelpModal() {
  uni.showModal({
    title: '壁纸使用指南',
    content: '1. 复制网盘链接与提取码下载原画 PNG\n2. iPhone 用户进入系统「设置」>「墙纸」选取照片并开启景深\n3. Mac 用户右键图片选择「设定为桌面图片」',
    showCancel: false
  })
}

function showSupportModal() {
  uni.showModal({
    title: 'VIP 专属客服与求图',
    content: '如遇到网盘链接失效、寻找特定风格壁纸或特殊屏幕比例定制，请联系客服邮箱：support@wallart.cc，工作时间内 1 小时内响应。',
    showCancel: false
  })
}

function showTermsModal() {
  uni.showModal({
    title: '服务协议与说明',
    content: '1. WallArt 画廊壁纸经由画师原创及正版开源授权精修\n2. 下载凭据永久有效，支持随时多端提取\n3. 尊重原创数字艺术版权，禁止商业转售。',
    showCancel: false
  })
}

function handleClearCache() {
  uni.showLoading({ title: '正在清理...' })
  setTimeout(() => {
    uni.hideLoading()
    cacheSize.value = '0 KB'
    uni.showToast({ title: '缓存已释放', icon: 'success' })
  }, 600)
}

function handleLogout() {
  userStore.logout()
  uni.showToast({ title: '已退出登录', icon: 'none' })
}

function showLoginPrompt() {
  showAuth.value = true
}

function handleAuthDone() {
  showAuth.value = false
  userStore.fetchProfile()
}
</script>

<style scoped>
.my-container {
  min-height: 100vh;
  position: relative;
  box-sizing: border-box;
  transition: all 0.25s ease;
  padding-bottom: 40rpx;
}

.my-content {
  position: relative;
  z-index: 2;
  max-width: 960rpx;
  margin: 0 auto;
  padding: calc(44rpx + env(safe-area-inset-top)) 32rpx calc(40rpx + env(safe-area-inset-bottom));
  box-sizing: border-box;
}

/* ========================================================
   Hero User Card
   ======================================================== */
.user-hero-card {
  border-radius: 32rpx;
  padding: 36rpx 32rpx 28rpx;
  margin-bottom: 28rpx;
  box-sizing: border-box;
  transition: all 0.25s ease;
}

.dark .user-hero-card {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.08), rgba(255, 255, 255, 0.03));
  border: 1rpx solid rgba(255, 255, 255, 0.12);
  box-shadow: 0 10rpx 36rpx rgba(0, 0, 0, 0.35);
}

.light .user-hero-card {
  background: #ffffff;
  border: 1rpx solid rgba(0, 0, 0, 0.08);
  box-shadow: 0 10rpx 32rpx rgba(0, 0, 0, 0.05);
}

.user-main-row {
  display: flex;
  align-items: center;
  gap: 24rpx;
}

.avatar-box {
  width: 108rpx;
  height: 108rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  flex-shrink: 0;
}

.dark .avatar-box {
  background: linear-gradient(135deg, #1e293b, #334155);
  border: 3rpx solid rgba(255, 255, 255, 0.2);
}

.light .avatar-box {
  background: linear-gradient(135deg, #e2e8f0, #cbd5e1);
  border: 3rpx solid #ffffff;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.08);
}

.avatar-vip {
  border-color: #f59e0b !important;
  box-shadow: 0 0 20rpx rgba(245, 158, 11, 0.35) !important;
}

.avatar-char {
  font-size: 48rpx;
  font-weight: 800;
}

.dark .avatar-char {
  color: #ffffff;
}

.light .avatar-char {
  color: #0f172a;
}

.avatar-crown-badge {
  position: absolute;
  bottom: -4rpx;
  right: -4rpx;
  width: 36rpx;
  height: 36rpx;
  border-radius: 50%;
  background: #f59e0b;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20rpx;
  border: 2rpx solid #000;
}

.user-info-col {
  flex: 1;
}

.name-row {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 12rpx;
}

.user-title {
  font-size: 34rpx;
  font-weight: 800;
  max-width: 380rpx;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.dark .user-title {
  color: #ffffff;
}

.light .user-title {
  color: #0f172a;
}

.vip-glow-pill {
  display: inline-flex;
  align-items: center;
  gap: 6rpx;
  padding: 4rpx 16rpx;
  border-radius: 999rpx;
  background: linear-gradient(135deg, #f59e0b, #d97706);
  box-shadow: 0 4rpx 12rpx rgba(245, 158, 11, 0.3);
}

.vip-pill-icon {
  font-size: 20rpx;
}

.vip-pill-text {
  font-size: 20rpx;
  font-weight: 800;
  color: #000000;
}

.free-pill {
  padding: 4rpx 14rpx;
  border-radius: 999rpx;
}

.dark .free-pill {
  background: rgba(255, 255, 255, 0.08);
}

.light .free-pill {
  background: rgba(0, 0, 0, 0.06);
}

.free-pill-text {
  font-size: 20rpx;
  font-weight: 600;
}

.dark .free-pill-text {
  color: rgba(255, 255, 255, 0.6);
}

.light .free-pill-text {
  color: #64748b;
}

.user-subtext {
  font-size: 22rpx;
  line-height: 1.45;
  margin-top: 10rpx;
  display: block;
}

.dark .user-subtext {
  color: rgba(255, 255, 255, 0.5);
}

.light .user-subtext {
  color: #64748b;
}

/* VIP Banner Inside Card */
.vip-banner-inside {
  margin-top: 24rpx;
  padding: 18rpx 24rpx;
  border-radius: 20rpx;
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.16), rgba(212, 175, 55, 0.08));
  border: 1rpx solid rgba(245, 158, 11, 0.35);
  display: flex;
  align-items: center;
  justify-content: space-between;
  cursor: pointer;
  transition: all 0.2s ease;
}

.vip-banner-inside:active {
  transform: scale(0.99);
}

.banner-left {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.banner-crown {
  font-size: 36rpx;
}

.banner-texts {
  display: flex;
  flex-direction: column;
}

.banner-t1 {
  font-size: 24rpx;
  font-weight: 800;
  color: #f59e0b;
}

.banner-t2 {
  font-size: 20rpx;
  margin-top: 2rpx;
}

.dark .banner-t2 {
  color: rgba(255, 255, 255, 0.6);
}

.light .banner-t2 {
  color: #475569;
}

.banner-btn {
  display: flex;
  align-items: center;
  gap: 4rpx;
  background: #f59e0b;
  padding: 8rpx 20rpx;
  border-radius: 999rpx;
}

.banner-btn-text {
  font-size: 22rpx;
  font-weight: 800;
  color: #000000;
}

.banner-btn-arrow {
  font-size: 24rpx;
  font-weight: 800;
  color: #000000;
}

/* ========================================================
   Quick 3-Metric Stats Row
   ======================================================== */
.stats-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16rpx;
  margin-bottom: 32rpx;
}

.stat-card {
  border-radius: 24rpx;
  padding: 24rpx 16rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s ease;
}

.dark .stat-card {
  background: rgba(255, 255, 255, 0.04);
  border: 1rpx solid rgba(255, 255, 255, 0.08);
}

.light .stat-card {
  background: #ffffff;
  border: 1rpx solid rgba(0, 0, 0, 0.06);
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.03);
}

.stat-num {
  font-size: 32rpx;
  font-weight: 800;
  margin-bottom: 6rpx;
}

.dark .stat-num {
  color: #ffffff;
}

.light .stat-num {
  color: #0f172a;
}

.gold-text {
  color: #f59e0b !important;
}

.stat-label {
  font-size: 20rpx;
}

.dark .stat-label {
  color: rgba(255, 255, 255, 0.5);
}

.light .stat-label {
  color: #64748b;
}

/* ========================================================
   Menu Groups
   ======================================================== */
.menu-group {
  margin-bottom: 28rpx;
}

.group-title-row {
  padding: 0 8rpx 14rpx;
}

.group-title {
  font-size: 22rpx;
  font-weight: 700;
  letter-spacing: 0.5rpx;
  text-transform: uppercase;
}

.dark .group-title {
  color: rgba(255, 255, 255, 0.45);
}

.light .group-title {
  color: #94a3b8;
}

.menu-box {
  border-radius: 28rpx;
  overflow: hidden;
}

.dark .menu-box {
  background: rgba(255, 255, 255, 0.04);
  border: 1rpx solid rgba(255, 255, 255, 0.08);
}

.light .menu-box {
  background: #ffffff;
  border: 1rpx solid rgba(0, 0, 0, 0.06);
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.03);
}

.menu-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 28rpx 24rpx;
  cursor: pointer;
  transition: background 0.15s ease;
}

.dark .menu-item {
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.05);
}

.light .menu-item {
  border-bottom: 1rpx solid rgba(0, 0, 0, 0.04);
}

.menu-item:last-child {
  border-bottom: none;
}

.dark .menu-item:active {
  background: rgba(255, 255, 255, 0.04);
}

.light .menu-item:active {
  background: #f8fafc;
}

.menu-left {
  display: flex;
  align-items: center;
  gap: 18rpx;
  flex: 1;
}

.item-icon-wrap {
  width: 64rpx;
  height: 64rpx;
  border-radius: 18rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.icon-blue { background: rgba(59, 130, 246, 0.14); }
.icon-gold { background: rgba(245, 158, 11, 0.16); }
.icon-amber { background: rgba(245, 158, 11, 0.12); }
.icon-purple { background: rgba(168, 85, 247, 0.14); }
.icon-cyan { background: rgba(6, 182, 212, 0.14); }
.icon-gray { background: rgba(148, 163, 184, 0.14); }
.icon-green { background: rgba(16, 185, 129, 0.14); }
.icon-slate { background: rgba(100, 116, 139, 0.14); }

.item-icon {
  font-size: 32rpx;
}

.item-texts {
  display: flex;
  flex-direction: column;
}

.item-title {
  font-size: 26rpx;
  font-weight: 700;
}

.dark .item-title {
  color: #ffffff;
}

.light .item-title {
  color: #0f172a;
}

.item-desc {
  font-size: 20rpx;
  line-height: 1.4;
  margin-top: 4rpx;
}

.dark .item-desc {
  color: rgba(255, 255, 255, 0.45);
}

.light .item-desc {
  color: #64748b;
}

.status-badge {
  display: flex;
  align-items: center;
  gap: 6rpx;
}

.status-badge-text {
  font-size: 22rpx;
  font-weight: 600;
}

.dark .status-badge-text {
  color: rgba(255, 255, 255, 0.45);
}

.light .status-badge-text {
  color: #64748b;
}

.menu-arrow {
  font-size: 30rpx;
  font-weight: 600;
}

.dark .menu-arrow {
  color: rgba(255, 255, 255, 0.25);
}

.light .menu-arrow {
  color: #cbd5e1;
}

/* ========================================================
   Account Actions
   ======================================================== */
.account-action-wrap {
  margin-top: 36rpx;
  margin-bottom: 40rpx;
}

.btn-account-login {
  width: 100%;
  height: 88rpx;
  border-radius: 999rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10rpx;
  cursor: pointer;
  transition: all 0.2s ease;
}

.dark .btn-account-login {
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.3), rgba(37, 99, 235, 0.15));
  border: 1rpx solid rgba(59, 130, 246, 0.5);
  box-shadow: 0 8rpx 24rpx rgba(59, 130, 246, 0.2);
}

.light .btn-account-login {
  background: linear-gradient(135deg, #2563eb, #1d4ed8);
  box-shadow: 0 8rpx 24rpx rgba(37, 99, 235, 0.25);
}

.btn-account-login:active {
  transform: scale(0.98);
}

.login-sparkle {
  font-size: 24rpx;
  color: #60a5fa;
}

.light .login-sparkle {
  color: #93c5fd;
}

.login-btn-text {
  font-size: 28rpx;
  font-weight: 800;
  color: #ffffff;
}

.btn-account-logout {
  width: 100%;
  height: 84rpx;
  border-radius: 999rpx;
  font-size: 26rpx;
  font-weight: 600;
  cursor: pointer;
}

.dark .btn-account-logout {
  background: rgba(255, 255, 255, 0.06);
  border: 1rpx solid rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.6);
}

.light .btn-account-logout {
  background: rgba(0, 0, 0, 0.04);
  border: 1rpx solid rgba(0, 0, 0, 0.08);
  color: #64748b;
}

/* ========================================================
   Branding Footer
   ======================================================== */
.app-branding-footer {
  text-align: center;
  padding: 16rpx 0 24rpx;
}

.brand-line-1 {
  font-size: 22rpx;
  font-weight: 600;
  display: block;
  margin-bottom: 6rpx;
}

.dark .brand-line-1 {
  color: rgba(255, 255, 255, 0.35);
}

.light .brand-line-1 {
  color: #94a3b8;
}

.brand-line-2 {
  font-size: 18rpx;
  display: block;
}

.dark .brand-line-2 {
  color: rgba(255, 255, 255, 0.2);
}

.light .brand-line-2 {
  color: #cbd5e1;
}

/* Emerald Icon */
.icon-emerald {
  background: rgba(16, 185, 129, 0.14);
  color: #10b981;
}

/* Custom Redeem Modal */
.custom-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 24rpx;
  box-sizing: border-box;
}

.redeem-modal-card {
  width: 100%;
  max-width: 600rpx;
  background: #181926;
  border: 1rpx solid rgba(255, 255, 255, 0.12);
  border-radius: 28rpx;
  padding: 32rpx;
  box-sizing: border-box;
  box-shadow: 0 20rpx 60rpx rgba(0, 0, 0, 0.6);
  color: #f1f5f9;
}

.light .redeem-modal-card {
  background: #ffffff;
  border-color: rgba(0, 0, 0, 0.08);
  box-shadow: 0 16rpx 48rpx rgba(0, 0, 0, 0.12);
  color: #0f172a;
}

.redeem-modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12rpx;
}

.r-title-wrap {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.r-title-icon {
  font-size: 32rpx;
}

.r-title-text {
  font-size: 30rpx;
  font-weight: 800;
}

.r-close-btn {
  font-size: 28rpx;
  color: #94a3b8;
  width: 44rpx;
  height: 44rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.r-desc-text {
  font-size: 22rpx;
  opacity: 0.65;
  line-height: 1.45;
  display: block;
  margin-bottom: 24rpx;
}

.r-input-group {
  display: flex;
  flex-direction: column;
  gap: 8rpx;
  margin-bottom: 18rpx;
}

.r-input-label {
  font-size: 22rpx;
  font-weight: 700;
  opacity: 0.8;
}

.r-input {
  width: 100%;
  height: 80rpx;
  box-sizing: border-box;
  border-radius: 16rpx;
  padding: 0 20rpx;
  font-size: 26rpx;
  border: 1rpx solid rgba(148, 163, 184, 0.25);
  background: rgba(148, 163, 184, 0.08);
}

.r-btn-row {
  display: flex;
  gap: 16rpx;
  margin-top: 28rpx;
}

.r-btn-cancel {
  flex: 1;
  height: 76rpx;
  border-radius: 999rpx;
  background: rgba(148, 163, 184, 0.12);
  border: none;
  font-size: 26rpx;
  font-weight: 700;
  color: #94a3b8;
  display: flex;
  align-items: center;
  justify-content: center;
}

.r-btn-confirm {
  flex: 1.5;
  height: 76rpx;
  border-radius: 999rpx;
  background: linear-gradient(135deg, #10b981, #059669);
  border: none;
  font-size: 26rpx;
  font-weight: 800;
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 6rpx 20rpx rgba(16, 185, 129, 0.35);
}
</style>
