<template>
  <view class="vip-container" :class="[appStore.theme]">
    <!-- Dynamic Ambient Background -->
    <AmbientCanvas />

    <!-- Top Sticky Frosted Glass Header -->
    <AppHeader />

    <!-- Main Content Wrapped with proper Z-index -->
    <view class="vip-content">
      <!-- VIP Luxury Hero Header -->
      <view class="vip-hero">
        <view class="vip-crown-badge">
          <text class="crown-icon">👑</text>
          <text class="crown-text">VIP 尊享特权会员</text>
        </view>
        <text class="vip-hero-title">全场 4K/8K 臻选壁纸无限畅享</text>
        <text class="vip-hero-sub">一次开通 · 无需单张重复付费 · 独家原画每日首发更新</text>
      </view>

      <!-- Plan Cards Selection -->
      <view class="plans-list">
        <view
          v-for="plan in planList"
          :key="plan.type"
          class="plan-tier-card"
          :class="{ featured: plan.featured, active: selectedPlan === plan.type }"
          @click="selectedPlan = plan.type"
          @tap="selectedPlan = plan.type"
        >
          <!-- Top Right Ribbon -->
          <view v-if="plan.badge" class="plan-ribbon">
            <text class="ribbon-text">{{ plan.badge }}</text>
          </view>

          <!-- Checkmark for selected plan -->
          <view v-if="selectedPlan === plan.type" class="plan-check-badge">
            <text class="check-icon">✓</text>
          </view>

          <view class="plan-header-row">
            <view class="tier-title-wrap">
              <text class="tier-name">{{ plan.name }}</text>
              <text class="tier-badge-pill" v-if="plan.subTag">{{ plan.subTag }}</text>
            </view>
            <view class="tier-price-row">
              <text class="tier-cur">$</text>
              <text class="tier-num">{{ plan.price }}</text>
              <text class="tier-origin" v-if="plan.originalPrice">${{ plan.originalPrice }}</text>
            </view>
          </view>

          <text class="tier-desc">{{ plan.desc }}</text>

          <view class="plan-footer-row">
            <view class="roi-tag">
              <text class="roi-icon">⚡</text>
              <text class="roi-text">{{ plan.roi }}</text>
            </view>
            <text v-if="plan.originalPrice" class="save-tag">立省 ${{ (Number(plan.originalPrice) - Number(plan.price)).toFixed(1) }}</text>
          </view>
        </view>
      </view>

      <!-- Primary Action CTA Card (With Email & Payment Channel Selector) -->
      <view class="vip-action-card">
        <view class="action-summary-row">
          <view class="summary-left">
            <text class="sum-label">当前所选：</text>
            <text class="sum-name">{{ currentPlan?.name }}</text>
          </view>
          <view class="summary-price">
            <text class="sum-cur">$</text>
            <text class="sum-num">{{ currentPlan?.price }}</text>
            <text class="sum-unit">USD</text>
          </view>
        </view>

        <!-- Account / Receipt Email Input -->
        <view class="vip-input-group">
          <view class="vip-input-label-row">
            <text class="vip-input-label">接收账号 / 凭证邮箱</text>
            <text class="quick-email-link" @click="generateGuestEmail" @tap="generateGuestEmail">🎲 匿名凭据</text>
          </view>
          <input
            v-model="email"
            type="text"
            placeholder="请输入用于绑定 VIP 特权的邮箱"
            class="vip-email-input"
          />
        </view>

        <!-- Payment Method Options -->
        <view class="payment-method-group">
          <text class="pay-method-label">支付通道选择：</text>
          <view class="pay-methods-grid">
            <!-- Channel 1: PayPal Official -->
            <view
              class="pay-method-card"
              :class="{ selected: paymentMethod === 'paypal' }"
              @click="paymentMethod = 'paypal'"
              @tap="paymentMethod = 'paypal'"
            >
              <view class="pay-radio">
                <view v-if="paymentMethod === 'paypal'" class="pay-radio-inner"></view>
              </view>
              <view class="pay-method-info">
                <view class="pay-title-row">
                  <text class="pay-icon-brand">💳</text>
                  <text class="pay-title">国际 PayPal 官方支付</text>
                  <text class="pay-tag-recommend">推荐</text>
                </view>
                <text class="pay-desc">支持 Visa / MasterCard / 借记卡及 PayPal 账户</text>
              </view>
            </view>

            <!-- Channel 2: Sandbox Fast Test (if enabled) -->
            <view
              v-if="allowTestPay"
              class="pay-method-card test-card"
              :class="{ selected: paymentMethod === 'test' }"
              @click="paymentMethod = 'test'"
              @tap="paymentMethod = 'test'"
            >
              <view class="pay-radio">
                <view v-if="paymentMethod === 'test'" class="pay-radio-inner test-radio"></view>
              </view>
              <view class="pay-method-info">
                <view class="pay-title-row">
                  <text class="pay-icon-brand">⚡</text>
                  <text class="pay-title">沙盒免密秒级测试通道</text>
                  <text class="pay-tag-sandbox">DEV</text>
                </view>
                <text class="pay-desc">无需真实扣款，一键模拟已支付并激活 VIP 权益</text>
              </view>
            </view>
          </view>
        </view>

        <!-- Big Primary CTA Button -->
        <button class="btn-subscribe" :loading="submitting" @click="handleSubscribe" @tap="handleSubscribe">
          <text class="btn-sub-icon">{{ paymentMethod === 'test' ? '⚡' : '👑' }}</text>
          <text class="btn-sub-text">
            {{ paymentMethod === 'test' ? `沙盒一键模拟开通 ($${currentPlan?.price} USD)` : `通过 PayPal 立即开通 ($${currentPlan?.price} USD)` }}
          </text>
        </button>

        <!-- Trust Badges -->
        <view class="trust-row">
          <view class="trust-item">
            <text class="trust-icon">🔒</text>
            <text class="trust-text">PayPal 官方安全加密</text>
          </view>
          <view class="trust-dot">·</view>
          <view class="trust-item">
            <text class="trust-icon">⚡</text>
            <text class="trust-text">开通秒级自动生效</text>
          </view>
          <view class="trust-dot">·</view>
          <view class="trust-item">
            <text class="trust-icon">🛡️</text>
            <text class="trust-text">4K/8K 超清全免单</text>
          </view>
        </view>
      </view>

      <!-- Privileges Checklist Card -->
      <view class="privileges-card">
        <view class="privilege-header">
          <text class="privilege-star">✦</text>
          <text class="privilege-title">会员尊享 6 大核心特权</text>
        </view>
        <view class="privilege-grid">
          <view v-for="item in benefits" :key="item.title" class="benefit-item">
            <view class="benefit-icon-box">
              <text class="benefit-icon">{{ item.icon }}</text>
            </view>
            <view class="benefit-texts">
              <text class="benefit-name">{{ item.title }}</text>
              <text class="benefit-sub">{{ item.desc }}</text>
            </view>
          </view>
        </view>
      </view>

      <!-- FAQ Section -->
      <view class="faq-card">
        <text class="faq-main-title">💬 常见问题与解答</text>
        <view class="faq-item">
          <text class="faq-q">Q: 开通后如何获取原画壁纸？</text>
          <text class="faq-a">A: 开通后进入任意壁纸详情页，点击「VIP 免单下载」即可直接获取高速夸克/百度网盘原画直连与无损文件提取码。</text>
        </view>
        <view class="faq-item">
          <text class="faq-q">Q: 支持哪些设备使用？</text>
          <text class="faq-a">A: 手机（iPhone/Android）、平板（iPad）、电脑（Mac/PC）全平台通用，同一个账号支持跨端共享特权。</text>
        </view>
        <view class="faq-item">
          <text class="faq-q">Q: 更换设备后如何找回会员？</text>
          <text class="faq-a">A: 在「我的」页面使用购买时填写的邮箱登录即可自动恢复 VIP 身份与所有历史下载凭据。</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { useUserStore } from '@/store/user'
import { useAppStore } from '@/store/app'
import { paypalAPI, configAPI } from '@/api'
import AmbientCanvas from '@/components/AmbientCanvas.vue'
import AppHeader from '@/components/AppHeader.vue'

const userStore = useUserStore()
const appStore = useAppStore()
const selectedPlan = ref('quarterly')
const paymentMethod = ref('paypal') // 'paypal' or 'test'
const email = ref(userStore.userInfo?.email || '')
const submitting = ref(false)
const allowTestPay = ref(true)

const planList = ref([
  {
    type: 'quarterly',
    name: '季度 VIP 会员',
    subTag: '热门推荐',
    price: '24.99',
    originalPrice: '49.99',
    unit: '季',
    desc: '90天全场数千款 4K/8K 壁纸无限畅下 · 每日独家首发',
    roi: '折合仅 $0.27/天',
    featured: true,
    badge: '🔥 85% 会员首选'
  },
  {
    type: 'monthly',
    name: '月度 VIP 会员',
    subTag: '尝鲜首选',
    price: '9.99',
    originalPrice: '19.99',
    unit: '月',
    desc: '30天全场无限免单 · 支持超清原画与动态视频',
    roi: '下载 3 款即回本'
  },
  {
    type: 'yearly',
    name: '年度至尊 VIP',
    subTag: '尊享年卡',
    price: '59.99',
    originalPrice: '129.99',
    unit: '年',
    desc: '365天全站特权 · 独家优先定制与 VIP 专属求图通道',
    roi: '低至 $0.16/天 · 限时特惠',
    badge: '⚡ 限时4.6折'
  }
])

const benefits = [
  { icon: '♾️', title: '全场免单无限下载', desc: '免除单款付费，数千张 4K/8K 超清原图任下' },
  { icon: '🎬', title: 'Live 动态壁纸特权', desc: 'iOS 锁屏动态、视频桌面全部解锁' },
  { icon: '⚡', title: '高速网盘直连通道', desc: '独享夸克/百度/谷歌高速直链提取码' },
  { icon: '📱', title: '多设备全通用', desc: '手机、电脑、平板账号互通无障碍' },
  { icon: '🎨', title: '每日独家首发更新', desc: '每日精选新专栏第一时间极速推送到手' },
  { icon: '💎', title: 'VIP 专属客服求图', desc: '壁纸找寻与特殊尺寸定制一对一支持' }
]

const currentPlan = computed(() => {
  return planList.value.find(x => x.type === selectedPlan.value) || planList.value[0]
})

function generateGuestEmail() {
  const rand = Math.floor(100000 + Math.random() * 900000)
  email.value = `vip_${rand}@wallart.vip`
  uni.showToast({ title: '已生成凭证邮箱', icon: 'none' })
}

function getPlanTitle(p) {
  const mapZh = {
    monthly: '月度 VIP 会员',
    quarterly: '季度 VIP 会员',
    yearly: '年度至尊 VIP',
    lifetime: '终生至尊 VIP'
  }
  const mapEn = {
    monthly: 'Monthly VIP',
    quarterly: 'Quarterly VIP',
    yearly: 'Yearly VIP',
    lifetime: 'Lifetime VIP'
  }
  if (!p.name_zh || p.name_zh.includes('?')) {
    return appStore.lang === 'zh' ? (mapZh[p.type] || 'VIP 会员') : (mapEn[p.type] || 'VIP Member')
  }
  return appStore.lang === 'zh' ? p.name_zh : (p.name_en || p.name_zh)
}

function getPlanBadge(p) {
  if (p.badge && !p.badge.includes('?')) return p.badge
  if (p.featured) return '🔥 热门推荐'
  return ''
}

async function loadPublicConfig() {
  try {
    const res = await configAPI.getPublic()
    if (res?.membership_plans && res.membership_plans.length > 0) {
      planList.value = res.membership_plans.map(p => ({
        type: p.type,
        name: getPlanTitle(p),
        subTag: getPlanBadge(p),
        price: String(p.price),
        originalPrice: p.save_percent ? String((p.price * 100 / (100 - p.save_percent)).toFixed(2)) : '',
        unit: p.period === 'year' ? '年' : (p.period === 'quarter' ? '季' : '月'),
        desc: (!p.desc_zh || p.desc_zh.includes('?')) ? `${p.duration_days || 30}天全场数千款 4K/8K 壁纸与动态视频无限免单` : (appStore.lang === 'zh' ? p.desc_zh : (p.desc_en || p.desc_zh)),
        roi: `折合 $${(p.price / (p.duration_days || 30)).toFixed(2)}/天`,
        featured: p.featured,
        badge: getPlanBadge(p)
      }))
    }
    if (res?.permissions?.allow_test_pay !== undefined) {
      allowTestPay.value = Boolean(res.permissions.allow_test_pay)
    }
  } catch (e) {
    // Keep local default plans
  }
}

async function handleSubscribe() {
  if (!email.value || !email.value.includes('@')) {
    uni.showToast({ title: '请输入有效的凭证邮箱', icon: 'none' })
    return
  }

  submitting.value = true
  const isTest = paymentMethod.value === 'test'

  try {
    uni.showToast({ title: isTest ? '正在模拟支付...' : '正在连接 PayPal...', icon: 'loading', duration: 15000 })
    
    const res = await paypalAPI.createOrder({
      order_type: 'membership',
      membership_type: selectedPlan.value,
      user_email: email.value.trim(),
      is_test: isTest
    })

    uni.hideToast()

    if (isTest && (res.status === 'paid' || res.is_test)) {
      if (res.token) {
        uni.setStorageSync('token', res.token)
      }
      if (res.user) {
        userStore.userInfo = {
          ...(userStore.userInfo || {}),
          ...res.user,
          is_member: true
        }
        uni.setStorageSync('user_info', JSON.stringify(userStore.userInfo))
      }
      uni.showModal({
        title: '🎉 VIP 特权开通成功！',
        content: `恭喜您已开通【${currentPlan.value?.name}】，全场 4K/8K 壁纸与动态视频已全部免单解锁！`,
        showCancel: false,
        confirmText: '立即畅享',
        success: () => {
          uni.switchTab({
            url: '/pages/index/index',
            fail: () => {
              // #ifdef H5
              if (typeof window !== 'undefined' && window.location) {
                window.location.hash = '#/pages/index/index'
              }
              // #endif
            }
          })
        }
      })
      return
    }

    if (res.approval_url) {
      // Store pending order ID
      uni.setStorageSync('pending_paypal_order_id', res.order_id)
      
      // #ifdef H5
      if (typeof window !== 'undefined') {
        window.location.href = res.approval_url
        return
      }
      // #endif

      // Non-H5 fallback
      uni.setClipboardData({
        data: res.approval_url,
        success: () => {
          uni.showModal({
            title: 'PayPal 付款链接已生成',
            content: '已将 PayPal 支付安全链接复制到剪贴板，请在浏览器中打开完成付款。',
            showCancel: false
          })
        }
      })
    } else {
      uni.showToast({ title: res.message || '发起支付异常', icon: 'none' })
    }
  } catch (err) {
    uni.hideToast()
    uni.showModal({
      title: '支付发起提示',
      content: err.message || 'PayPal 网关连接超时，请确认 PayPal Client ID 与 Secret 配置无误，或选择沙盒免密秒级测试。',
      showCancel: false
    })
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadPublicConfig()
  if (!email.value && userStore.userInfo?.email) {
    email.value = userStore.userInfo.email
  }
})

onLoad((options) => {
  // If returning from PayPal directly to membership page
  const orderId = options?.order_id || options?.token
  if (orderId) {
    uni.navigateTo({
      url: `/pages/payment/success?order_id=${orderId}`
    })
  }
})
</script>

<style scoped>
.vip-container {
  min-height: 100vh;
  position: relative;
  box-sizing: border-box;
  transition: all 0.25s ease;
  padding-bottom: 40rpx;
}

.vip-content {
  position: relative;
  z-index: 2;
  max-width: 960rpx;
  margin: 0 auto;
  padding: calc(44rpx + env(safe-area-inset-top)) 32rpx calc(40rpx + env(safe-area-inset-bottom));
  box-sizing: border-box;
}

/* ========================================================
   Hero Header
   ======================================================== */
.vip-hero {
  text-align: center;
  margin-bottom: 44rpx;
}

.vip-crown-badge {
  display: inline-flex;
  align-items: center;
  gap: 8rpx;
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.2), rgba(212, 175, 55, 0.1));
  border: 1rpx solid rgba(245, 158, 11, 0.4);
  padding: 8rpx 24rpx;
  border-radius: 999rpx;
  margin-bottom: 20rpx;
  box-shadow: 0 4rpx 16rpx rgba(245, 158, 11, 0.15);
}

.crown-icon {
  font-size: 24rpx;
}

.crown-text {
  font-size: 22rpx;
  font-weight: 800;
  color: #f59e0b;
  letter-spacing: 0.5rpx;
}

.vip-hero-title {
  font-size: 40rpx;
  font-weight: 800;
  line-height: 1.35;
  display: block;
  letter-spacing: -0.5rpx;
}

.dark .vip-hero-title {
  color: #ffffff;
}

.light .vip-hero-title {
  color: #0f172a;
}

.vip-hero-sub {
  font-size: 24rpx;
  line-height: 1.5;
  margin-top: 12rpx;
  display: block;
}

.dark .vip-hero-sub {
  color: rgba(255, 255, 255, 0.6);
}

.light .vip-hero-sub {
  color: #64748b;
}

/* ========================================================
   Plan Cards
   ======================================================== */
.plans-list {
  display: flex;
  flex-direction: column;
  gap: 24rpx;
  margin-bottom: 36rpx;
}

.plan-tier-card {
  position: relative;
  border-radius: 28rpx;
  padding: 32rpx 32rpx 28rpx;
  cursor: pointer;
  transition: all 0.25s cubic-bezier(0.16, 1, 0.3, 1);
  box-sizing: border-box;
}

/* Dark mode card */
.dark .plan-tier-card {
  background: rgba(255, 255, 255, 0.05);
  border: 2rpx solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.2);
}

.dark .plan-tier-card.active {
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.16), rgba(212, 175, 55, 0.06));
  border-color: #f59e0b;
  box-shadow: 0 8rpx 32rpx rgba(245, 158, 11, 0.25);
}

/* Light mode card */
.light .plan-tier-card {
  background: rgba(255, 255, 255, 0.9);
  border: 2rpx solid rgba(0, 0, 0, 0.08);
  box-shadow: 0 6rpx 24rpx rgba(0, 0, 0, 0.05);
}

.light .plan-tier-card.active {
  background: linear-gradient(135deg, rgba(254, 243, 199, 0.8), rgba(255, 255, 255, 0.95));
  border-color: #d97706;
  box-shadow: 0 8rpx 30rpx rgba(217, 119, 6, 0.18);
}

.plan-ribbon {
  position: absolute;
  top: -14rpx;
  right: 28rpx;
  background: linear-gradient(135deg, #f59e0b, #ef4444);
  padding: 4rpx 18rpx;
  border-radius: 999rpx;
  box-shadow: 0 4rpx 14rpx rgba(239, 68, 68, 0.35);
  z-index: 3;
}

.ribbon-text {
  font-size: 20rpx;
  font-weight: 800;
  color: #ffffff;
}

.plan-check-badge {
  position: absolute;
  bottom: -10rpx;
  right: -10rpx;
  width: 44rpx;
  height: 44rpx;
  border-radius: 50%;
  background: #f59e0b;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4rpx 12rpx rgba(245, 158, 11, 0.4);
}

.check-icon {
  font-size: 26rpx;
  font-weight: 800;
  color: #000000;
}

.plan-header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12rpx;
}

.tier-title-wrap {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.tier-name {
  font-size: 32rpx;
  font-weight: 800;
}

.dark .tier-name {
  color: #ffffff;
}

.light .tier-name {
  color: #0f172a;
}

.tier-badge-pill {
  font-size: 20rpx;
  font-weight: 700;
  padding: 2rpx 12rpx;
  border-radius: 999rpx;
  background: rgba(245, 158, 11, 0.15);
  color: #f59e0b;
}

.tier-price-row {
  display: flex;
  align-items: baseline;
  gap: 4rpx;
}

.tier-cur {
  font-size: 26rpx;
  font-weight: 800;
  color: #f59e0b;
}

.tier-num {
  font-size: 46rpx;
  font-weight: 900;
  color: #f59e0b;
  letter-spacing: -1rpx;
}

.tier-origin {
  font-size: 22rpx;
  color: #94a3b8;
  text-decoration: line-through;
  margin-left: 6rpx;
}

.tier-desc {
  font-size: 24rpx;
  line-height: 1.5;
  display: block;
  margin-bottom: 18rpx;
}

.dark .tier-desc {
  color: rgba(255, 255, 255, 0.6);
}

.light .tier-desc {
  color: #475569;
}

.plan-footer-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.roi-tag {
  display: inline-flex;
  align-items: center;
  gap: 6rpx;
  background: rgba(245, 158, 11, 0.14);
  padding: 6rpx 16rpx;
  border-radius: 12rpx;
}

.roi-icon {
  font-size: 20rpx;
  color: #f59e0b;
}

.roi-text {
  font-size: 22rpx;
  color: #f59e0b;
  font-weight: 700;
}

.save-tag {
  font-size: 22rpx;
  font-weight: 700;
  color: #ef4444;
}

/* ========================================================
   Action CTA Card
   ======================================================== */
.vip-action-card {
  border-radius: 28rpx;
  padding: 32rpx;
  margin-bottom: 40rpx;
  box-sizing: border-box;
}

.dark .vip-action-card {
  background: rgba(255, 255, 255, 0.05);
  border: 1rpx solid rgba(255, 255, 255, 0.12);
  box-shadow: 0 8rpx 30rpx rgba(0, 0, 0, 0.35);
}

.light .vip-action-card {
  background: #ffffff;
  border: 1rpx solid rgba(0, 0, 0, 0.08);
  box-shadow: 0 8rpx 30rpx rgba(0, 0, 0, 0.06);
}

.action-summary-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24rpx;
}

.summary-left {
  display: flex;
  align-items: center;
  gap: 8rpx;
}

.sum-label {
  font-size: 24rpx;
}

.dark .sum-label {
  color: rgba(255, 255, 255, 0.5);
}

.light .sum-label {
  color: #64748b;
}

.sum-name {
  font-size: 28rpx;
  font-weight: 800;
}

.dark .sum-name {
  color: #ffffff;
}

.light .sum-name {
  color: #0f172a;
}

.summary-price {
  display: flex;
  align-items: baseline;
}

.sum-cur {
  font-size: 26rpx;
  font-weight: 800;
  color: #f59e0b;
}

.sum-num {
  font-size: 44rpx;
  font-weight: 900;
  color: #f59e0b;
}

.sum-unit {
  font-size: 20rpx;
  color: #f59e0b;
  margin-left: 6rpx;
  font-weight: 700;
}

.vip-input-group {
  margin-bottom: 24rpx;
}

.vip-input-label-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12rpx;
}

.vip-input-label {
  font-size: 22rpx;
  color: #94a3b8;
}

.quick-email-link {
  font-size: 22rpx;
  color: #3b82f6;
  cursor: pointer;
}

.vip-email-input {
  width: 100%;
  height: 80rpx;
  border-radius: 18rpx;
  padding: 0 24rpx;
  font-size: 26rpx;
  box-sizing: border-box;
}

.dark .vip-email-input {
  background: rgba(255, 255, 255, 0.08);
  border: 1rpx solid rgba(255, 255, 255, 0.16);
  color: #ffffff;
}

.light .vip-email-input {
  background: #f8fafc;
  border: 1rpx solid rgba(0, 0, 0, 0.1);
  color: #0f172a;
}

/* Payment Methods */
.payment-method-group {
  margin-bottom: 28rpx;
}

.pay-method-label {
  font-size: 22rpx;
  color: #94a3b8;
  display: block;
  margin-bottom: 14rpx;
}

.pay-methods-grid {
  display: flex;
  flex-direction: column;
  gap: 16rpx;
}

.pay-method-card {
  display: flex;
  align-items: center;
  border-radius: 20rpx;
  padding: 20rpx 22rpx;
  cursor: pointer;
  transition: all 0.2s ease;
  box-sizing: border-box;
}

.dark .pay-method-card {
  background: rgba(255, 255, 255, 0.04);
  border: 1rpx solid rgba(255, 255, 255, 0.1);
}

.light .pay-method-card {
  background: #f8fafc;
  border: 1rpx solid rgba(0, 0, 0, 0.08);
}

.pay-method-card.selected {
  border-color: #3b82f6 !important;
  background: rgba(59, 130, 246, 0.12) !important;
}

.pay-method-card.test-card.selected {
  border-color: #f59e0b !important;
  background: rgba(245, 158, 11, 0.12) !important;
}

.pay-radio {
  width: 34rpx;
  height: 34rpx;
  border-radius: 50%;
  border: 2rpx solid rgba(255, 255, 255, 0.3);
  margin-right: 18rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.light .pay-radio {
  border-color: rgba(0, 0, 0, 0.25);
}

.pay-method-card.selected .pay-radio {
  border-color: #3b82f6;
}

.pay-method-card.test-card.selected .pay-radio {
  border-color: #f59e0b;
}

.pay-radio-inner {
  width: 18rpx;
  height: 18rpx;
  border-radius: 50%;
  background: #3b82f6;
}

.test-radio {
  background: #f59e0b !important;
}

.pay-method-info {
  flex: 1;
}

.pay-title-row {
  display: flex;
  align-items: center;
  gap: 10rpx;
}

.pay-icon-brand {
  font-size: 26rpx;
}

.pay-title {
  font-size: 26rpx;
  font-weight: 700;
}

.dark .pay-title {
  color: #ffffff;
}

.light .pay-title {
  color: #0f172a;
}

.pay-tag-recommend {
  font-size: 18rpx;
  font-weight: 700;
  padding: 2rpx 10rpx;
  border-radius: 999rpx;
  background: rgba(59, 130, 246, 0.2);
  color: #3b82f6;
}

.pay-tag-sandbox {
  font-size: 18rpx;
  font-weight: 700;
  padding: 2rpx 10rpx;
  border-radius: 999rpx;
  background: rgba(245, 158, 11, 0.2);
  color: #f59e0b;
}

.pay-desc {
  font-size: 20rpx;
  color: #94a3b8;
  display: block;
  margin-top: 4rpx;
}

.btn-subscribe {
  width: 100%;
  height: 92rpx;
  background: linear-gradient(135deg, #f59e0b, #d97706);
  border-radius: 999rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12rpx;
  box-shadow: 0 10rpx 32rpx rgba(245, 158, 11, 0.4);
  cursor: pointer;
  transition: all 0.2s ease;
  margin-bottom: 24rpx;
}

.btn-subscribe:active {
  transform: scale(0.98);
}

.btn-sub-icon {
  font-size: 32rpx;
}

.btn-sub-text {
  font-size: 30rpx;
  font-weight: 800;
  color: #000000;
  letter-spacing: 0.5rpx;
}

.trust-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12rpx;
}

.trust-item {
  display: flex;
  align-items: center;
  gap: 6rpx;
}

.trust-icon {
  font-size: 20rpx;
}

.trust-text {
  font-size: 20rpx;
}

.dark .trust-text {
  color: rgba(255, 255, 255, 0.45);
}

.light .trust-text {
  color: #64748b;
}

.trust-dot {
  font-size: 20rpx;
  color: #94a3b8;
}

/* ========================================================
   Privileges Grid
   ======================================================== */
.privileges-card {
  border-radius: 28rpx;
  padding: 36rpx 32rpx;
  margin-bottom: 36rpx;
}

.dark .privileges-card {
  background: rgba(255, 255, 255, 0.04);
  border: 1rpx solid rgba(255, 255, 255, 0.08);
}

.light .privileges-card {
  background: #ffffff;
  border: 1rpx solid rgba(0, 0, 0, 0.06);
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.04);
}

.privilege-header {
  display: flex;
  align-items: center;
  gap: 10rpx;
  margin-bottom: 28rpx;
}

.privilege-star {
  font-size: 26rpx;
  color: #f59e0b;
}

.privilege-title {
  font-size: 30rpx;
  font-weight: 800;
}

.dark .privilege-title {
  color: #ffffff;
}

.light .privilege-title {
  color: #0f172a;
}

.privilege-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 28rpx 20rpx;
}

@media (max-width: 480px) {
  .privilege-grid {
    grid-template-columns: 1fr;
    gap: 20rpx;
  }
}

.benefit-item {
  display: flex;
  align-items: flex-start;
  gap: 16rpx;
}

.benefit-icon-box {
  width: 60rpx;
  height: 60rpx;
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.dark .benefit-icon-box {
  background: rgba(255, 255, 255, 0.06);
}

.light .benefit-icon-box {
  background: rgba(0, 0, 0, 0.04);
}

.benefit-icon {
  font-size: 30rpx;
}

.benefit-texts {
  display: flex;
  flex-direction: column;
}

.benefit-name {
  font-size: 26rpx;
  font-weight: 700;
}

.dark .benefit-name {
  color: #ffffff;
}

.light .benefit-name {
  color: #0f172a;
}

.benefit-sub {
  font-size: 20rpx;
  line-height: 1.45;
  margin-top: 6rpx;
}

.dark .benefit-sub {
  color: rgba(255, 255, 255, 0.45);
}

.light .benefit-sub {
  color: #64748b;
}

/* ========================================================
   FAQ Section
   ======================================================== */
.faq-card {
  border-radius: 28rpx;
  padding: 32rpx;
  margin-bottom: 24rpx;
}

.dark .faq-card {
  background: rgba(255, 255, 255, 0.03);
  border: 1rpx solid rgba(255, 255, 255, 0.06);
}

.light .faq-card {
  background: #ffffff;
  border: 1rpx solid rgba(0, 0, 0, 0.06);
}

.faq-main-title {
  font-size: 28rpx;
  font-weight: 800;
  display: block;
  margin-bottom: 24rpx;
}

.dark .faq-main-title {
  color: #ffffff;
}

.light .faq-main-title {
  color: #0f172a;
}

.faq-item {
  margin-bottom: 20rpx;
}

.faq-item:last-child {
  margin-bottom: 0;
}

.faq-q {
  font-size: 24rpx;
  font-weight: 700;
  display: block;
  margin-bottom: 6rpx;
}

.dark .faq-q {
  color: #fde68a;
}

.light .faq-q {
  color: #b45309;
}

.faq-a {
  font-size: 22rpx;
  line-height: 1.5;
  display: block;
}

.dark .faq-a {
  color: rgba(255, 255, 255, 0.5);
}

.light .faq-a {
  color: #64748b;
}
</style>
