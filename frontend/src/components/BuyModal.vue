<template>
  <transition name="fade">
    <div class="modal-overlay" @click.self="$emit('close')">
      <div class="modal-box buy-modal">
        <!-- Mobile Bottom Sheet Handle -->
        <div class="sheet-drag-handle show-mobile-handle"></div>
        <!-- Header -->
        <div class="modal-header">
          <div class="modal-title-wrap">
            <span class="modal-badge-icon">{{ isUnlocked ? '🎉' : (purchaseChannel === 'guest' ? '⚡' : '👑') }}</span>
            <h3>{{ modalTitle }}</h3>
          </div>
          <button class="close-btn" @click="$emit('close')">✕</button>
        </div>

        <div class="modal-body">
          <!-- Product Summary Card -->
          <div class="product-summary">
            <div class="product-thumb">
              <img v-if="product?.cover_image" :src="product.cover_image" :alt="product.title_zh" />
              <div v-else class="thumb-placeholder" :style="{ background: product?._color || '#667eea' }">🖼️</div>
            </div>
            <div class="product-info">
              <p class="product-name">{{ locale === 'zh' ? product?.title_zh : (product?.title_en || product?.title_zh) }}</p>
              <div class="product-price-row">
                <template v-if="auth.isMember">
                  <span class="badge badge-gold">👑 {{ locale === 'zh' ? 'VIP会员免单' : 'Member Free' }}</span>
                </template>
                <template v-else-if="isFree">
                  <span class="badge badge-success">🎁 {{ locale === 'zh' ? '免费专享壁纸' : 'Free Wallpaper' }}</span>
                </template>
                <template v-else>
                  <span class="price-currency">$</span>
                  <span class="price-amount">{{ product?.price }}</span>
                  <span class="price-unit"> USD</span>
                </template>
              </div>
            </div>
          </div>

          <!-- Unlocked / Free View: Show Download Links Directly -->
          <template v-if="isUnlocked">
            <div class="unlocked-view">
              <div class="unlocked-banner">
                <span class="banner-icon">🎉</span>
                <div>
                  <h4 class="banner-title">{{ locale === 'zh' ? '已成功解锁！请选择网盘下载' : 'Unlocked! Choose Cloud Drive' }}</h4>
                  <p class="banner-sub">{{ locale === 'zh' ? '包含百度网盘、夸克网盘和Google Drive，点击可一键复制提取码' : 'Download links ready below' }}</p>
                </div>
              </div>
              <DownloadLinks :product-id="product?.id" />
            </div>
          </template>

          <!-- Paid & Not Unlocked: Purchase Channel Selector -->
          <template v-else>
            <!-- Mode Tabs: 临时购买 vs 会员登录购买 -->
            <div v-if="guestBuyEnabled" class="channel-mode-switch">
              <button
                class="channel-tab-btn"
                :class="{ active: purchaseChannel === 'guest' }"
                @click="switchChannel('guest')"
              >
                <span class="tab-icon">⚡</span>
                <div class="tab-text">
                  <span class="tab-main">{{ locale === 'zh' ? '临时购买 (免登录)' : 'Guest Checkout' }}</span>
                  <span class="tab-sub">{{ locale === 'zh' ? '无需注册・极速解锁' : 'Instant & No Login' }}</span>
                </div>
              </button>

              <button
                class="channel-tab-btn"
                :class="{ active: purchaseChannel === 'member' }"
                @click="switchChannel('member')"
              >
                <span class="tab-icon">👑</span>
                <div class="tab-text">
                  <span class="tab-main">{{ locale === 'zh' ? '会员登录购买' : 'Member Checkout' }}</span>
                  <span class="tab-sub">{{ locale === 'zh' ? '永久保存订单 / VIP免单' : 'Saved in Account / VIP' }}</span>
                </div>
              </button>
            </div>

            <!-- CHANNEL 1: 临时购买 (免登录游客快捷支付) -->
            <div v-if="purchaseChannel === 'guest'" class="channel-content guest-content">
              <div class="channel-tip-box">
                <span class="tip-badge">⚡ {{ locale === 'zh' ? '临时游客快捷通道' : 'Guest Quick Pay' }}</span>
                <p>{{ locale === 'zh' ? '无需注册或登录账号，付款后当前窗口立即展示网盘提取码及下载地址！' : 'No login required. Complete payment to reveal cloud drive links immediately.' }}</p>
              </div>

              <!-- Guest Email Input -->
              <div class="guest-email-box">
                <div class="email-label-row">
                  <label class="input-label">
                    <span>📧 {{ locale === 'zh' ? '接收邮箱 (用于保存订单凭据)' : 'Your Email (To receive order)' }}</span>
                  </label>
                  <button type="button" class="btn-text-link" @click="generateRandomEmail">
                    🎲 {{ locale === 'zh' ? '快捷生成匿名邮箱' : 'Generate Guest Email' }}
                  </button>
                </div>
                <input
                  v-model="guestEmail"
                  type="email"
                  class="form-input"
                  :placeholder="locale === 'zh' ? '例如: yourname@gmail.com' : 'e.g. yourname@gmail.com'"
                />
              </div>

              <!-- Payment Method Sub-tabs -->
              <div class="payment-tabs">
                <button
                  class="pay-tab-btn"
                  :class="{ active: guestPayTab === 'test' }"
                  @click="guestPayTab = 'test'"
                >
                  ⚡ {{ locale === 'zh' ? '快速测试支付' : 'Fast Test Pay' }}
                </button>
                <button
                  class="pay-tab-btn"
                  :class="{ active: guestPayTab === 'paypal' }"
                  @click="guestPayTab = 'paypal'; initPayPal()"
                >
                  💳 PayPal {{ locale === 'zh' ? '官方结账' : 'Checkout' }}
                </button>
                <button
                  class="pay-tab-btn"
                  :class="{ active: guestPayTab === 'code' }"
                  @click="guestPayTab = 'code'"
                >
                  🎫 {{ locale === 'zh' ? '卡密兑换' : 'Redeem Code' }}
                </button>
              </div>

              <!-- Sub-tab 1: Fast Sandbox / Test Pay -->
              <div v-if="guestPayTab === 'test'" class="tab-pane test-pane">
                <div class="test-notice">
                  <span class="notice-badge">⚡ {{ locale === 'zh' ? '本地免扣款体验' : 'Demo Mode' }}</span>
                  <p>{{ locale === 'zh' ? '一键秒级模拟支付成功，真实体验从付款到下载解锁的全闭环！' : '1-click simulated payment for testing. Instant download unlock without charge.' }}</p>
                </div>
                <button class="btn btn-primary btn-block btn-lg test-pay-btn" :disabled="processing" @click="handleGuestTestPay">
                  <span v-if="processing" class="spinner"></span>
                  <span v-else>⚡ {{ locale === 'zh' ? `模拟支付 $${product?.price} 立即解锁` : `Simulate $${product?.price} Pay & Unlock` }}</span>
                </button>
              </div>

              <!-- Sub-tab 2: PayPal Official Payment -->
              <div v-if="guestPayTab === 'paypal'" class="tab-pane paypal-pane">
                <div id="paypal-button-container" ref="paypalContainer"></div>
                <div v-if="paypalLoading" class="paypal-loading-ph">
                  <span class="spinner"></span>
                  <p>{{ locale === 'zh' ? '正在连接 PayPal 安全支付网关...' : 'Connecting to PayPal secure gateway...' }}</p>
                </div>
                <div v-if="paypalLoadFailed" class="paypal-fallback">
                  <p>⚠️ {{ locale === 'zh' ? 'PayPal 网关连接超时（可能是当前网络环境未开放），推荐使用快速测试支付：' : 'PayPal timeout. Try fast test checkout:' }}</p>
                  <button class="btn btn-outline btn-block btn-sm" @click="guestPayTab = 'test'">
                    ⚡ {{ locale === 'zh' ? '切换为快速模拟测试支付' : 'Switch to Fast Test Pay' }}
                  </button>
                </div>
              </div>

              <!-- Sub-tab 3: Verify Code / Redeem Token -->
              <div v-if="guestPayTab === 'code'" class="tab-pane code-pane">
                <div class="code-box">
                  <p class="code-desc">
                    {{ locale === 'zh' ? '如果您有作者赠送的提取码、活动口令或卡密（内置可用口令: loveyou, 666, 1194, VIP888）：' : 'Enter author gift code or redeem token:' }}
                  </p>
                  <div class="code-input-row">
                    <input
                      v-model="unlockCode"
                      type="text"
                      class="form-input code-input"
                      :placeholder="locale === 'zh' ? '输入口令 (如 666 或 loveyou)' : 'Enter code (e.g. 666)'"
                      @keyup.enter="handleCodeUnlock"
                    />
                    <button class="btn btn-dark" :disabled="!unlockCode.trim() || processing" @click="handleCodeUnlock">
                      <span v-if="processing" class="spinner"></span>
                      <span v-else>{{ locale === 'zh' ? '兑换解锁' : 'Redeem' }}</span>
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <!-- CHANNEL 2: 会员登录购买 -->
            <div v-else class="channel-content member-content">
              <!-- Case A: NOT Logged In -->
              <div v-if="!auth.isLoggedIn" class="member-auth-box">
                <div class="member-intro-banner">
                  <span class="intro-icon">👑</span>
                  <div>
                    <h4 class="intro-title">{{ locale === 'zh' ? '会员专享特权通道' : 'Member Exclusive Channel' }}</h4>
                    <p class="intro-sub">{{ locale === 'zh' ? '登录账号后购买，壁纸将永久保存在【个人中心 - 我的订单】中；或开通 VIP 畅享全场免单！' : 'Login to permanently save orders in your account, or subscribe VIP for unlimited downloads.' }}</p>
                  </div>
                </div>

                <!-- Mini Login / Register Tabs -->
                <div class="auth-sub-tabs">
                  <button
                    class="auth-sub-btn"
                    :class="{ active: authMode === 'login' }"
                    @click="authMode = 'login'"
                  >
                    🔑 {{ locale === 'zh' ? '账号登录' : 'Login' }}
                  </button>
                  <button
                    class="auth-sub-btn"
                    :class="{ active: authMode === 'register' }"
                    @click="authMode = 'register'"
                  >
                    ✨ {{ locale === 'zh' ? '极速注册' : 'Register' }}
                  </button>
                </div>

                <!-- Login Form -->
                <form v-if="authMode === 'login'" @submit.prevent="handleQuickLogin" class="auth-mini-form">
                  <div class="form-group">
                    <input
                      v-model="loginEmail"
                      type="email"
                      class="form-input"
                      :placeholder="locale === 'zh' ? '登录邮箱' : 'Email'"
                      required
                    />
                  </div>
                  <div class="form-group">
                    <input
                      v-model="loginPassword"
                      type="password"
                      class="form-input"
                      :placeholder="locale === 'zh' ? '密码' : 'Password'"
                      required
                    />
                  </div>
                  <div class="form-actions-col">
                    <button type="submit" class="btn btn-dark btn-block btn-lg" :disabled="authLoading">
                      <span v-if="authLoading" class="spinner"></span>
                      <span v-else>🔑 {{ locale === 'zh' ? '登录并继续购买' : 'Login & Continue' }}</span>
                    </button>
                    <button type="button" class="btn btn-ghost btn-xs btn-block fill-demo-btn" @click="fillAdminAccount">
                      ⚡ {{ locale === 'zh' ? '一键填入管理员测试账号 (admin@wallpaper.com)' : 'Fill Admin Demo Account' }}
                    </button>
                  </div>
                </form>

                <!-- Register Form -->
                <form v-else @submit.prevent="handleQuickRegister" class="auth-mini-form">
                  <div class="form-group">
                    <input
                      v-model="registerEmail"
                      type="email"
                      class="form-input"
                      :placeholder="locale === 'zh' ? '注册邮箱' : 'Email'"
                      required
                    />
                  </div>
                  <div class="form-group">
                    <input
                      v-model="registerPassword"
                      type="password"
                      class="form-input"
                      :placeholder="locale === 'zh' ? '设置密码 (至少6位)' : 'Password (min 6 chars)'"
                      required
                      minlength="6"
                    />
                  </div>
                  <button type="submit" class="btn btn-primary btn-block btn-lg" :disabled="authLoading">
                    <span v-if="authLoading" class="spinner"></span>
                    <span v-else>✨ {{ locale === 'zh' ? '注册并登录' : 'Register & Login' }}</span>
                  </button>
                </form>

                <!-- VIP Upsell link -->
                <div class="direct-vip-bar">
                  <span>{{ locale === 'zh' ? '不想单张购买？' : 'Want all wallpapers?' }}</span>
                  <router-link to="/membership" class="direct-vip-link" @click="$emit('close')">
                    👑 {{ locale === 'zh' ? '开通 VIP 会员全场壁纸免费下 ($9.99/月起) →' : 'Get VIP Membership →' }}
                  </router-link>
                </div>
              </div>

              <!-- Case B: ALREADY Logged In -->
              <div v-else class="member-logged-box">
                <!-- User Account Info Card -->
                <div class="user-card-info">
                  <div class="user-avatar-badge">{{ auth.isMember ? '👑' : '👤' }}</div>
                  <div class="user-details">
                    <span class="user-email-text">{{ auth.user?.email }}</span>
                    <span class="user-role-tag" :class="{ 'tag-vip': auth.isMember }">
                      {{ auth.isMember ? (locale === 'zh' ? 'VIP尊贵会员 (全场免单)' : 'VIP Member') : (locale === 'zh' ? '普通注册会员' : 'Regular Member') }}
                    </span>
                  </div>
                </div>

                <!-- If already VIP: celebrate & direct unlock -->
                <div v-if="auth.isMember" class="vip-benefit-pane">
                  <p class="vip-notice">
                    👑 {{ locale === 'zh' ? '您的 VIP 会员特权已生效！本款壁纸已为您全场免单，无需再支付任何费用。' : 'Your VIP membership is active! This wallpaper is 100% free for you.' }}
                  </p>
                  <button class="btn btn-gold btn-block btn-lg" @click="handleVipDirectUnlock">
                    📥 {{ locale === 'zh' ? '立即查看并下载网盘资源' : 'Access Downloads Now' }}
                  </button>
                </div>

                <!-- If logged in as regular user -->
                <div v-else class="regular-member-pane">
                  <p class="member-buy-desc">
                    {{ locale === 'zh' ? `以当前账号下单购买，壁纸将永久存入您的个人中心订单中。` : 'Purchasing with your account permanently saves this item to your orders.' }}
                  </p>

                  <div class="member-actions-stack">
                    <button class="btn btn-primary btn-block btn-lg" :disabled="processing" @click="handleMemberTestPay">
                      <span v-if="processing" class="spinner"></span>
                      <span v-else>⚡ {{ locale === 'zh' ? `模拟账号支付 $${product?.price} (写入我的订单)` : `Test Pay $${product?.price} to Account` }}</span>
                    </button>

                    <div class="or-divider">
                      <span>{{ locale === 'zh' ? '或者' : 'OR' }}</span>
                    </div>

                    <router-link to="/membership" class="btn btn-gold btn-block btn-lg" @click="$emit('close')">
                      👑 {{ locale === 'zh' ? '升级 VIP 全场壁纸免费下载 ($9.99/月起)' : 'Upgrade to VIP Unlimited' }}
                    </router-link>
                  </div>
                </div>
              </div>
            </div>
          </template>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup>
import { ref, computed, watch, onMounted, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/store'
import { paypalAPI, downloadAPI, configAPI, authAPI } from '@/api'
import { useToast } from '@/composables/useToast'
import DownloadLinks from './DownloadLinks.vue'
import { loadScript } from '@paypal/paypal-js'

const props = defineProps({
  product: { type: Object, default: null },
  initialMode: { type: String, default: 'guest' } // 'guest' | 'member' | 'code'
})
const emit = defineEmits(['close', 'success'])

const { t, locale } = useI18n()
const auth = useAuthStore()
const { error, success } = useToast()

const paypalContainer = ref(null)
const processing = ref(false)
const paypalLoading = ref(false)
const paypalLoadFailed = ref(false)
const unlockedLocally = ref(false)
const guestBuyEnabled = ref(true)

// Channels & Tabs
const purchaseChannel = ref('guest') // 'guest' | 'member'
const guestPayTab = ref('test') // 'test' | 'paypal' | 'code'
const guestEmail = ref('')
const unlockCode = ref('')

async function loadPublicSettings() {
  try {
    const cfg = await configAPI.getPublic()
    if (cfg?.permissions) {
      guestBuyEnabled.value = cfg.permissions.guest_buy_enabled !== false
      if (!guestBuyEnabled.value) {
        purchaseChannel.value = 'member'
      }
    }
  } catch {}
}

onMounted(() => {
  loadPublicSettings()
})

// Member inline auth
const authMode = ref('login') // 'login' | 'register'
const loginEmail = ref('')
const loginPassword = ref('')
const registerEmail = ref('')
const registerPassword = ref('')
const authLoading = ref(false)

watch(
  () => [props.initialMode, guestBuyEnabled.value],
  ([mode, guestAllowed]) => {
    if (!guestAllowed || mode === 'member') {
      purchaseChannel.value = 'member'
    } else if (mode === 'code') {
      purchaseChannel.value = 'guest'
      guestPayTab.value = 'code'
    } else {
      purchaseChannel.value = 'guest'
    }
  },
  { immediate: true }
)

const isFree = computed(() => Number(props.product?.price || 0) <= 0)
const isUnlocked = computed(() => auth.isMember || isFree.value || unlockedLocally.value)

const modalTitle = computed(() => {
  if (isUnlocked.value) return locale.value === 'zh' ? '壁纸下载资源' : 'Wallpaper Downloads'
  if (purchaseChannel.value === 'member') return locale.value === 'zh' ? '会员登录购买' : 'Member Purchase'
  return locale.value === 'zh' ? '临时购买 (免登录)' : 'Guest Purchase'
})

function switchChannel(channel) {
  purchaseChannel.value = channel
  if (channel === 'guest' && guestPayTab.value === 'paypal') {
    nextTick(() => initPayPal())
  }
}

function generateRandomEmail() {
  const rand = Math.random().toString(36).substring(2, 8)
  guestEmail.value = `buyer_${rand}@wallpaper.com`
  success(locale.value === 'zh' ? '已自动生成临时匿名邮箱！' : 'Anonymous guest email generated!')
}

function fillAdminAccount() {
  loginEmail.value = 'admin@wallpaper.com'
  loginPassword.value = 'admin123456'
  success(locale.value === 'zh' ? '已自动填入管理员测试账号' : 'Admin credentials filled')
}

// 1. Guest Test Pay
async function handleGuestTestPay() {
  if (processing.value || !props.product) return
  processing.value = true
  try {
    const res = await paypalAPI.createOrder({
      product_id: props.product.id,
      amount: props.product.price,
      order_type: 'single',
      is_test: true,
      user_email: guestEmail.value
    })
    if (res.token) {
      if (res.user) {
        auth.setAuth(res.user, res.token)
      } else {
        localStorage.setItem('token', res.token)
        auth.token = res.token
      }
    }
    unlockedLocally.value = true
    success(locale.value === 'zh' ? '🎉 临时支付模拟成功，已为您解锁下载资源！' : 'Payment successful! Downloads unlocked.')
    emit('success')
  } catch (err) {
    error(err.message || (locale.value === 'zh' ? '支付测试失败' : 'Test payment failed'))
  } finally {
    processing.value = false
  }
}

// 2. Member Test Pay
async function handleMemberTestPay() {
  if (processing.value || !props.product) return
  processing.value = true
  try {
    const res = await paypalAPI.createOrder({
      product_id: props.product.id,
      amount: props.product.price,
      order_type: 'single',
      is_test: true,
      user_email: auth.user?.email || guestEmail.value
    })
    if (res.token) {
      if (res.user) {
        auth.setAuth(res.user, res.token)
      } else {
        localStorage.setItem('token', res.token)
        auth.token = res.token
      }
    }
    unlockedLocally.value = true
    success(locale.value === 'zh' ? '🎉 订单已存入您的会员账号，网盘资源已成功解锁！' : 'Order saved to your account! Downloads unlocked.')
    emit('success')
  } catch (err) {
    error(err.message || (locale.value === 'zh' ? '会员支付失败' : 'Payment failed'))
  } finally {
    processing.value = false
  }
}

// 3. Unlock with promo/verify code
async function handleCodeUnlock() {
  if (!unlockCode.value.trim() || processing.value || !props.product) return
  processing.value = true
  try {
    const res = await downloadAPI.unlockWithCode(props.product.id, {
      code: unlockCode.value.trim(),
      user_email: guestEmail.value || auth.user?.email
    })
    if (res.token) {
      if (res.user) {
        auth.setAuth(res.user, res.token)
      } else {
        localStorage.setItem('token', res.token)
        auth.token = res.token
      }
    }
    unlockedLocally.value = true
    success(res.message || (locale.value === 'zh' ? '🎉 验证码正确，壁纸已解锁！' : 'Unlocked successfully!'))
    emit('success')
  } catch (err) {
    error(err.message || (locale.value === 'zh' ? '验证码无效或错误' : 'Invalid code'))
  } finally {
    processing.value = false
  }
}

// 4. Quick Member Login
async function handleQuickLogin() {
  if (!loginEmail.value || !loginPassword.value) return
  authLoading.value = true
  try {
    const res = await authAPI.login({
      email: loginEmail.value.trim(),
      password: loginPassword.value
    })
    auth.setAuth(res.user, res.token)
    success(locale.value === 'zh' ? `欢迎回来，${res.user.email}！` : 'Login successful!')
    if (auth.isMember) {
      unlockedLocally.value = true
      emit('success')
    }
  } catch (err) {
    error(err.message || (locale.value === 'zh' ? '登录失败，请检查账号密码' : 'Login failed'))
  } finally {
    authLoading.value = false
  }
}

// 5. Quick Member Register
async function handleQuickRegister() {
  if (!registerEmail.value || !registerPassword.value) return
  authLoading.value = true
  try {
    const res = await authAPI.register({
      email: registerEmail.value.trim(),
      password: registerPassword.value
    })
    auth.setAuth(res.user, res.token)
    success(locale.value === 'zh' ? '注册成功，已自动登录！' : 'Registered successfully!')
  } catch (err) {
    error(err.message || (locale.value === 'zh' ? '注册失败，该邮箱可能已被占用' : 'Registration failed'))
  } finally {
    authLoading.value = false
  }
}

function handleVipDirectUnlock() {
  unlockedLocally.value = true
  emit('success')
}

// 6. PayPal Setup
async function initPayPal() {
  if (isUnlocked.value || !props.product) return
  paypalLoading.value = true
  paypalLoadFailed.value = false

  try {
    let clientId = 'sb'
    try {
      const cfg = await configAPI.getPublic()
      if (cfg?.paypal_client_id) clientId = cfg.paypal_client_id
    } catch {}

    const paypal = await loadScript({
      clientId: clientId,
      currency: 'USD'
    })
    if (!paypal || !paypalContainer.value) return

    await nextTick()
    paypalContainer.value.innerHTML = ''

    await paypal.Buttons({
      style: {
        layout: 'vertical',
        color: 'gold',
        shape: 'pill',
        label: 'pay',
        height: 44
      },
      createOrder: async () => {
        if (purchaseChannel.value === 'guest' && !auth.isLoggedIn) {
          if (!guestEmail.value || !guestEmail.value.includes('@') || !guestEmail.value.includes('.')) {
            error(locale.value === 'zh' ? '请先填写有效的联系邮箱以接收订单与网盘提取码' : 'Please enter a valid email to receive order and downloads')
            throw new Error('Valid email required')
          }
        }
        processing.value = true
        try {
          const res = await paypalAPI.createOrder({
            product_id: props.product.id,
            amount: props.product.price,
            order_type: 'single',
            user_email: guestEmail.value || auth.user?.email
          })
          if (res.token) localStorage.setItem('token', res.token)
          return res.paypal_order_id || res.id
        } catch (e) {
          error(e.message || t('payment.failed'))
          throw e
        } finally {
          processing.value = false
        }
      },
      onApprove: async (data) => {
        processing.value = true
        try {
          const res = await paypalAPI.captureOrder(data.orderID)
          if (res.token) {
            if (res.user) {
              auth.setAuth(res.user, res.token)
            } else {
              localStorage.setItem('token', res.token)
              auth.token = res.token
            }
          }
          unlockedLocally.value = true
          success(t('payment.success'))
          emit('success')
        } catch (e) {
          error(t('payment.failed'))
        } finally {
          processing.value = false
        }
      },
      onCancel: () => {
        error(t('payment.cancel'))
      },
      onError: (err) => {
        console.error('PayPal button error:', err)
        paypalLoadFailed.value = true
      }
    }).render(paypalContainer.value)
  } catch (e) {
    console.warn('PayPal script load warning:', e)
    paypalLoadFailed.value = true
  } finally {
    paypalLoading.value = false
  }
}

onMounted(() => {
  if (auth.user?.email) {
    guestEmail.value = auth.user.email
  }
})
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.68);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
}

.buy-modal {
  width: 100%;
  max-width: 520px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  background: var(--c-surface, #ffffff);
  border-radius: 24px;
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.25);
  border: 1px solid rgba(255, 255, 255, 0.2);
  overflow: hidden;
  animation: modalPop 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

[data-theme="dark"] .buy-modal {
  background: #18181c;
  border-color: rgba(255, 255, 255, 0.1);
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.5);
}

@keyframes modalPop {
  from { opacity: 0; transform: scale(0.95) translateY(10px); }
  to { opacity: 1; transform: scale(1) translateY(0); }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--c-border, #eee);
}

[data-theme="dark"] .modal-header {
  border-color: rgba(255, 255, 255, 0.08);
}

.modal-title-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
}

.modal-badge-icon {
  font-size: 1.25rem;
}

.modal-header h3 {
  font-size: 1.12rem;
  font-weight: 700;
  margin: 0;
  color: var(--c-text, #111);
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.15rem;
  color: #888;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 8px;
  transition: all 0.15s;
}

.close-btn:hover {
  background: rgba(0, 0, 0, 0.05);
  color: #111;
}

[data-theme="dark"] .close-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
}

.modal-body {
  padding: 18px 20px;
  overflow-y: auto;
}

/* Product summary pill */
.product-summary {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 14px;
  background: var(--c-bg, #f5f5f7);
  border-radius: 16px;
  margin-bottom: 16px;
}

[data-theme="dark"] .product-summary {
  background: rgba(255, 255, 255, 0.05);
}

.product-thumb {
  width: 52px;
  height: 64px;
  border-radius: 10px;
  overflow: hidden;
  flex-shrink: 0;
  background: #111;
}

.product-thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.thumb-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.4rem;
}

.product-info {
  flex: 1;
  min-width: 0;
}

.product-name {
  font-size: 0.95rem;
  font-weight: 700;
  margin: 0 0 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: var(--c-text, #111);
}

.product-price-row {
  display: flex;
  align-items: baseline;
  gap: 2px;
}

.price-currency {
  font-size: 0.85rem;
  font-weight: 700;
  color: #007aff;
}

.price-amount {
  font-size: 1.3rem;
  font-weight: 800;
  color: #007aff;
}

.price-unit {
  font-size: 0.72rem;
  font-weight: 600;
  color: #888;
  margin-left: 2px;
}

/* Channel Mode Switch */
.channel-mode-switch {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
  padding: 4px;
  background: var(--c-bg, #f0f0f4);
  border-radius: 16px;
  margin-bottom: 16px;
}

[data-theme="dark"] .channel-mode-switch {
  background: rgba(255, 255, 255, 0.06);
}

.channel-tab-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  border-radius: 12px;
  border: none;
  background: transparent;
  cursor: pointer;
  text-align: left;
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

.channel-tab-btn.active {
  background: #ffffff;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

[data-theme="dark"] .channel-tab-btn.active {
  background: #28282e;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.4);
}

.tab-icon {
  font-size: 1.25rem;
  flex-shrink: 0;
}

.tab-text {
  display: flex;
  flex-direction: column;
}

.tab-main {
  font-size: 0.88rem;
  font-weight: 700;
  color: var(--c-text, #111);
  line-height: 1.2;
}

.tab-sub {
  font-size: 0.68rem;
  color: #777;
  margin-top: 2px;
}

/* Tip Box */
.channel-tip-box {
  background: rgba(0, 122, 255, 0.08);
  border: 1px solid rgba(0, 122, 255, 0.2);
  border-radius: 12px;
  padding: 10px 12px;
  margin-bottom: 14px;
  font-size: 0.78rem;
  color: var(--c-text, #333);
  line-height: 1.45;
}

.tip-badge {
  display: inline-block;
  font-weight: 700;
  color: #007aff;
  margin-bottom: 4px;
}

.channel-tip-box p {
  margin: 0;
}

/* Guest Email Box */
.guest-email-box {
  margin-bottom: 14px;
}

.email-label-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.input-label {
  font-size: 0.76rem;
  font-weight: 700;
  color: var(--c-text, #444);
}

.btn-text-link {
  background: none;
  border: none;
  font-size: 0.72rem;
  font-weight: 600;
  color: #007aff;
  cursor: pointer;
  padding: 0;
}

.btn-text-link:hover {
  text-decoration: underline;
}

/* Payment Sub-tabs */
.payment-tabs {
  display: flex;
  gap: 6px;
  margin-bottom: 14px;
}

.pay-tab-btn {
  flex: 1;
  padding: 8px 10px;
  border-radius: 10px;
  font-size: 0.78rem;
  font-weight: 700;
  border: 1px solid var(--c-border, #e2e8f0);
  background: var(--c-surface, #fff);
  color: var(--c-text, #4b5563);
  cursor: pointer;
  transition: all 0.15s;
  text-align: center;
}

[data-theme="dark"] .pay-tab-btn {
  border-color: rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.04);
  color: #9ca3af;
}

.pay-tab-btn.active {
  background: #111;
  color: #fff;
  border-color: #111;
}

[data-theme="dark"] .pay-tab-btn.active {
  background: #007aff;
  border-color: #007aff;
  color: #fff;
}

/* Panes */
.tab-pane {
  padding: 4px 0;
}

.test-notice {
  background: rgba(16, 185, 129, 0.1);
  border: 1px solid rgba(16, 185, 129, 0.25);
  border-radius: 12px;
  padding: 10px 12px;
  margin-bottom: 14px;
  font-size: 0.78rem;
  color: var(--c-text, #333);
  line-height: 1.4;
}

.notice-badge {
  display: inline-block;
  font-weight: 700;
  color: #10b981;
  margin-bottom: 3px;
}

.test-notice p {
  margin: 0;
}

.test-pay-btn {
  background: linear-gradient(135deg, #007aff, #00c6ff);
  border: none;
  font-weight: 800;
  letter-spacing: 0.2px;
  box-shadow: 0 4px 14px rgba(0, 122, 255, 0.35);
}

.test-pay-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 18px rgba(0, 122, 255, 0.45);
}

.paypal-loading-ph {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 24px 0;
  color: #666;
  font-size: 0.84rem;
  gap: 8px;
}

.paypal-fallback {
  background: rgba(245, 158, 11, 0.1);
  border: 1px solid rgba(245, 158, 11, 0.3);
  border-radius: 12px;
  padding: 12px;
  font-size: 0.8rem;
  color: var(--c-text, #333);
  margin-top: 10px;
}

.code-box {
  background: var(--c-bg, #f8f9fa);
  padding: 14px;
  border-radius: 14px;
}

[data-theme="dark"] .code-box {
  background: rgba(255, 255, 255, 0.04);
}

.code-desc {
  font-size: 0.78rem;
  color: #666;
  margin: 0 0 10px;
  line-height: 1.45;
}

.code-input-row {
  display: flex;
  gap: 8px;
}

.code-input {
  font-size: 0.88rem;
  font-weight: 700;
  letter-spacing: 1px;
}

/* Member Channel Elements */
.member-intro-banner {
  display: flex;
  gap: 10px;
  padding: 12px 14px;
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.12), rgba(217, 119, 6, 0.06));
  border: 1px solid rgba(245, 158, 11, 0.25);
  border-radius: 14px;
  margin-bottom: 14px;
}

.intro-icon {
  font-size: 1.5rem;
  flex-shrink: 0;
}

.intro-title {
  font-size: 0.88rem;
  font-weight: 700;
  margin: 0 0 3px;
  color: #d97706;
}

.intro-sub {
  font-size: 0.74rem;
  color: #666;
  margin: 0;
  line-height: 1.4;
}

.auth-sub-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
}

.auth-sub-btn {
  flex: 1;
  padding: 8px 12px;
  border-radius: 10px;
  font-size: 0.8rem;
  font-weight: 700;
  border: 1px solid var(--c-border, #eee);
  background: transparent;
  cursor: pointer;
  color: #666;
  transition: all 0.15s;
}

.auth-sub-btn.active {
  background: #111;
  color: #fff;
  border-color: #111;
}

[data-theme="dark"] .auth-sub-btn.active {
  background: #fff;
  color: #111;
  border-color: #fff;
}

.auth-mini-form {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.form-actions-col {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-top: 4px;
}

.fill-demo-btn {
  color: #007aff;
  font-size: 0.72rem;
  padding: 4px;
}

.direct-vip-bar {
  margin-top: 14px;
  text-align: center;
  font-size: 0.76rem;
  color: #777;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.direct-vip-link {
  font-weight: 700;
  color: #d97706;
  text-decoration: none;
}

.direct-vip-link:hover {
  text-decoration: underline;
}

/* User Card Info (Logged in) */
.user-card-info {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 14px;
  background: var(--c-bg, #f5f5f7);
  border-radius: 14px;
  margin-bottom: 14px;
}

[data-theme="dark"] .user-card-info {
  background: rgba(255, 255, 255, 0.05);
}

.user-avatar-badge {
  font-size: 1.6rem;
}

.user-details {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.user-email-text {
  font-size: 0.88rem;
  font-weight: 700;
  color: var(--c-text, #111);
}

.user-role-tag {
  font-size: 0.72rem;
  font-weight: 600;
  color: #666;
}

.user-role-tag.tag-vip {
  color: #d97706;
  font-weight: 700;
}

.vip-notice {
  font-size: 0.84rem;
  color: #d97706;
  line-height: 1.45;
  background: rgba(245, 158, 11, 0.1);
  padding: 12px;
  border-radius: 12px;
  margin: 0 0 14px;
}

.member-buy-desc {
  font-size: 0.82rem;
  color: #666;
  margin: 0 0 14px;
  line-height: 1.4;
}

.member-actions-stack {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.or-divider {
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  margin: 2px 0;
}

.or-divider::before {
  content: '';
  position: absolute;
  left: 0;
  right: 0;
  height: 1px;
  background: var(--c-border, #eee);
}

[data-theme="dark"] .or-divider::before {
  background: rgba(255, 255, 255, 0.08);
}

.or-divider span {
  position: relative;
  background: var(--c-surface, #fff);
  padding: 0 10px;
  font-size: 0.72rem;
  color: #999;
}

[data-theme="dark"] .or-divider span {
  background: #18181c;
}

.btn-gold {
  background: linear-gradient(135deg, #f59e0b, #d97706);
  color: #fff;
  font-weight: 800;
  box-shadow: 0 4px 14px rgba(245, 158, 11, 0.35);
  border: none;
}

.btn-gold:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 18px rgba(245, 158, 11, 0.45);
}

/* Unlocked View */
.unlocked-banner {
  display: flex;
  gap: 12px;
  padding: 12px 14px;
  background: rgba(16, 185, 129, 0.1);
  border: 1px solid rgba(16, 185, 129, 0.3);
  border-radius: 14px;
  margin-bottom: 16px;
}

.banner-icon {
  font-size: 1.6rem;
  flex-shrink: 0;
}

.banner-title {
  font-size: 0.95rem;
  font-weight: 700;
  color: #059669;
  margin: 0 0 3px;
}

.banner-sub {
  font-size: 0.75rem;
  color: #666;
  margin: 0;
}

/* Mobile Responsive Bottom Sheet */
.show-mobile-handle {
  display: none;
}

@media (max-width: 640px) {
  .modal-overlay {
    align-items: flex-end !important;
    padding: 0 !important;
  }

  .buy-modal {
    max-width: 100% !important;
    border-radius: 26px 26px 0 0 !important;
    max-height: 88vh !important;
    animation: sheetSlideUp 0.32s cubic-bezier(0.16, 1, 0.3, 1) !important;
  }

  @keyframes sheetSlideUp {
    from { transform: translateY(100%); }
    to { transform: translateY(0); }
  }

  .show-mobile-handle {
    display: block;
    width: 44px;
    height: 5px;
    background: #d1d5db;
    border-radius: 999px;
    margin: 10px auto 0;
  }

  [data-theme="dark"] .show-mobile-handle {
    background: #4b5563;
  }

  .modal-header {
    padding: 12px 18px !important;
  }

  .modal-body {
    padding: 16px 18px max(24px, env(safe-area-inset-bottom)) !important;
  }

  .channel-mode-switch {
    gap: 8px !important;
  }

  .channel-tab-btn {
    padding: 10px 8px !important;
    min-height: 52px !important;
  }

  .tab-icon {
    font-size: 1.3rem !important;
  }

  .tab-main {
    font-size: 0.85rem !important;
  }

  .tab-sub {
    font-size: 0.7rem !important;
  }

  .payment-tabs {
    gap: 6px !important;
  }

  .pay-tab-btn {
    padding: 10px 6px !important;
    font-size: 0.8rem !important;
    min-height: 42px !important;
  }

  .btn-lg {
    padding: 14px 20px !important;
    font-size: 0.98rem !important;
  }
}
</style>
