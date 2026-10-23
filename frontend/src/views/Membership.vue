<template>
  <div class="membership-page">
    <!-- Header -->
    <section class="membership-hero">
      <div class="container">
        <div class="hero-label">👑 {{ locale === 'zh' ? '会员方案' : 'Membership Plans' }}</div>
        <h1>{{ t('membership.title') }}</h1>
        <p>{{ t('membership.subtitle') }}</p>

        <!-- Current status -->
        <div v-if="auth.isMember" class="current-plan-alert">
          <span>✅</span>
          <span>
            {{ locale === 'zh' ? '您当前是' : 'Current plan: ' }}
            <strong>{{ t(`membership.${auth.membershipType}`) }}</strong>
            {{ locale === 'zh' ? '，到期时间：' : ', expires: ' }}
            <strong>{{ formatDate(auth.user?.membership_expire) }}</strong>
          </span>
        </div>
      </div>
    </section>

    <!-- Plans -->
    <section class="plans-section">
      <div class="container">
        <div class="plans-grid">
          <div
            v-for="plan in plans"
            :key="plan.type"
            class="plan-card"
            :class="{ featured: plan.featured, current: auth.membershipType === plan.type }"
          >
            <div v-if="plan.featured || plan.badge" class="plan-badge">{{ plan.badge || t('membership.recommended') }}</div>
            <div v-if="auth.membershipType === plan.type" class="plan-badge current-badge">
              {{ t('dashboard.active') }}
            </div>

            <div class="plan-header">
              <div class="plan-icon">{{ plan.icon || '👑' }}</div>
              <h3 class="plan-name">{{ locale === 'zh' ? (plan.name_zh || t(`membership.${plan.type}`)) : (plan.name_en || t(`membership.${plan.type}`)) }}</h3>
              <div class="plan-price">
                <span class="price-currency">$</span>
                <span class="price-num">{{ plan.price }}</span>
                <span class="price-period">{{ formatPeriod(plan) }}</span>
              </div>
              <div v-if="plan.save_percent || plan.savePercent" class="plan-save">
                {{ t('membership.save') }} {{ plan.save_percent || plan.savePercent }}%
              </div>
            </div>

            <div class="plan-features">
              <div v-for="feature in features" :key="feature" class="feature-item">
                <span class="feature-check">✓</span>
                <span>{{ t(`membership.features.${feature}`) }}</span>
              </div>
            </div>

            <button
              class="btn plan-btn"
              :class="plan.featured ? 'btn-primary' : 'btn-outline'"
              :disabled="auth.membershipType === plan.type"
              @click="handleSubscribe(plan)"
            >
              <span v-if="subscribing === plan.type" class="spinner"></span>
              <span v-else-if="auth.membershipType === plan.type">{{ t('membership.current_plan') }}</span>
              <span v-else>{{ t('membership.subscribe') }}</span>
            </button>
          </div>
        </div>

        <!-- PayPal note -->
        <div class="payment-info">
          <span>🔒</span>
          <span>
            {{ locale === 'zh'
              ? '通过 PayPal 安全订阅，随时可以取消'
              : 'Secure subscription via PayPal. Cancel anytime.' }}
          </span>
          <div class="paypal-badge">PayPal</div>
        </div>
      </div>
    </section>

    <!-- FAQ -->
    <section class="faq-section">
      <div class="container">
        <h2>{{ locale === 'zh' ? '常见问题' : 'FAQ' }}</h2>
        <div class="faq-list">
          <div v-for="(faq, i) in faqs" :key="i" class="faq-item">
            <button class="faq-q" @click="toggleFaq(i)">
              {{ locale === 'zh' ? faq.q_zh : faq.q_en }}
              <span class="faq-arrow" :class="{ open: openFaq === i }">›</span>
            </button>
            <transition name="slide">
              <div v-if="openFaq === i" class="faq-a">
                {{ locale === 'zh' ? faq.a_zh : faq.a_en }}
              </div>
            </transition>
          </div>
        </div>
      </div>
    </section>

    <!-- Subscribe Modal (PayPal & Sandbox Test) -->
    <div v-if="showPayPalModal" class="modal-overlay mobile-sheet-overlay" @click.self="showPayPalModal = false">
      <div class="modal-box sub-modal-sheet">
        <div class="sheet-drag-handle show-mobile-handle"></div>
        <div class="modal-header">
          <h3>{{ locale === 'zh' ? '开通 VIP 会员' : 'Subscribe' }} - {{ locale === 'zh' ? (selectedPlan?.name_zh || t(`membership.${selectedPlan?.type}`)) : (selectedPlan?.name_en || t(`membership.${selectedPlan?.type}`)) }}</h3>
          <button class="close-btn" @click="showPayPalModal = false">✕</button>
        </div>
        <div class="modal-body" style="padding: 24px;">
          <div class="subscribe-summary" style="display:flex;align-items:center;gap:14px;padding:14px 18px;background:var(--c-bg);border-radius:16px;margin-bottom:18px;">
            <span class="sub-icon" style="font-size:2rem;">{{ selectedPlan?.icon || '👑' }}</span>
            <div>
              <p class="sub-plan" style="font-weight:700;margin:0 0 4px;">{{ locale === 'zh' ? (selectedPlan?.name_zh || t(`membership.${selectedPlan?.type}`)) : (selectedPlan?.name_en || t(`membership.${selectedPlan?.type}`)) }}</p>
              <p class="sub-price" style="font-size:1.2rem;font-weight:800;color:var(--c-accent);margin:0;">$ {{ selectedPlan?.price }} {{ formatPeriod(selectedPlan) }}</p>
            </div>
          </div>

          <!-- Checkout Tabs -->
          <div style="display:flex;gap:6px;background:var(--c-bg);padding:4px;border-radius:12px;margin-bottom:18px;">
            <button
              class="tab-btn"
              :class="{ active: subTab === 'test' }"
              style="flex:1;padding:8px;font-size:0.82rem;font-weight:700;border:none;border-radius:9px;cursor:pointer;"
              @click="subTab = 'test'"
            >
              ⚡ {{ locale === 'zh' ? '快速模拟开通' : 'Fast Test Pay' }}
            </button>
            <button
              class="tab-btn"
              :class="{ active: subTab === 'paypal' }"
              style="flex:1;padding:8px;font-size:0.82rem;font-weight:700;border:none;border-radius:9px;cursor:pointer;"
              @click="subTab = 'paypal'; initPayPalOrder(selectedPlan)"
            >
              💳 PayPal {{ locale === 'zh' ? '官方结算' : 'Checkout' }}
            </button>
          </div>

          <!-- Test Mode -->
          <div v-if="subTab === 'test'">
            <div style="background:rgba(0,113,227,0.08);border:1px dashed rgba(0,113,227,0.3);padding:14px;border-radius:12px;margin-bottom:18px;font-size:0.84rem;color:#4b5563;line-height:1.4;">
              <strong style="color:var(--c-accent);display:block;margin-bottom:4px;">⚡ {{ locale === 'zh' ? '免扣款测试通道' : 'Instant Test Subscription' }}</strong>
              {{ locale === 'zh' ? '点击下方按钮将立即模拟支付成功，并为您授予该方案对应的全站壁纸无限下载权益！' : '1-click simulated payment for testing. Instantly grants VIP benefits.' }}
            </div>
            <button class="btn btn-primary btn-block btn-lg" :disabled="subscribing" @click="handleTestSubscribe">
              <span v-if="subscribing" class="spinner"></span>
              <span v-else>⚡ {{ locale === 'zh' ? `模拟支付 $${selectedPlan?.price} 并立即开通` : `Simulate $${selectedPlan?.price} & Activate` }}</span>
            </button>
          </div>

          <!-- PayPal Mode -->
          <div v-if="subTab === 'paypal'">
            <div id="paypal-subscription-container" ref="paypalSubContainer" style="min-height:80px;"></div>
            <p v-if="paypalLoading" style="text-align:center;font-size:0.84rem;color:#888;padding:16px 0;">
              <span class="spinner"></span> {{ locale === 'zh' ? '正在加载 PayPal 支付组件...' : 'Loading PayPal...' }}
            </p>
            <p class="payment-note" style="margin-top:12px;font-size:0.78rem;color:#888;text-align:center;">
              🔒 {{ locale === 'zh' ? 'PayPal 个人卖家安全结算，支付成功即刻生效' : 'PayPal secure payment. Access granted immediately.' }}
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, nextTick, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store'
import { paypalAPI, configAPI } from '@/api'
import { useToast } from '@/composables/useToast'
import { loadScript } from '@paypal/paypal-js'

const { t, locale } = useI18n()
const auth = useAuthStore()
const router = useRouter()
const { success, error } = useToast()

const subscribing = ref(false)
const showPayPalModal = ref(false)
const selectedPlan = ref(null)
const paypalSubContainer = ref(null)
const openFaq = ref(null)
const subTab = ref('test')
const paypalLoading = ref(false)

const plans = ref([
  { type: 'monthly', name_zh: '月度会员', name_en: 'Monthly VIP', price: 9.99, period: 'month', duration_days: 30, icon: '🌙', featured: false, badge: '入门尝鲜' },
  { type: 'quarterly', name_zh: '季度会员', name_en: 'Quarterly VIP', price: 24.99, period: 'quarter', duration_days: 90, icon: '⭐', featured: true, badge: '热销推荐 / 省17%' },
  { type: 'yearly', name_zh: '年度会员', name_en: 'Yearly VIP', price: 79.99, period: 'year', duration_days: 365, icon: '👑', featured: false, badge: '超值之选 / 省33%' },
  { type: 'lifetime', name_zh: '终生尊享', name_en: 'Lifetime VIP', price: 199.99, period: 'lifetime', duration_days: 36500, icon: '💎', featured: false, badge: '终身买断' }
])

function formatPeriod(plan) {
  if (!plan) return ''
  if (plan.period === 'lifetime' || (plan.duration_days && plan.duration_days >= 3650)) {
    return locale.value === 'zh' ? '/ 终身买断' : '/ Lifetime'
  }
  if (plan.period === 'month' || (plan.duration_days && plan.duration_days <= 31)) {
    return locale.value === 'zh' ? '/ 月' : '/ Month'
  }
  if (plan.period === 'quarter' || (plan.duration_days && plan.duration_days <= 92)) {
    return locale.value === 'zh' ? '/ 季' : '/ Quarter'
  }
  if (plan.period === 'year' || (plan.duration_days && plan.duration_days <= 366)) {
    return locale.value === 'zh' ? '/ 年' : '/ Year'
  }
  return `/${plan.duration_days} ${locale.value === 'zh' ? '天' : 'days'}`
}

async function loadPlans() {
  try {
    const cfg = await configAPI.getPublic()
    if (cfg?.membership_plans && Array.isArray(cfg.membership_plans) && cfg.membership_plans.length > 0) {
      plans.value = cfg.membership_plans.filter(p => p.enabled !== false)
    }
  } catch (err) {
    console.warn('Failed to load dynamic plans:', err)
  }
}

onMounted(() => {
  loadPlans()
})

const features = ['all_wallpapers', 'hd_download', 'early_access', 'multi_drive', 'priority_support']

const faqs = [
  {
    q_zh: '购买会员后如何下载壁纸？',
    a_zh: '成为会员后，您可以在任何壁纸详情页直接看到所有网盘下载链接，包括百度网盘、夸克网盘和Google Drive。',
    q_en: 'How do I download wallpapers after purchasing membership?',
    a_en: 'After becoming a member, you can directly see all cloud drive download links on any wallpaper detail page, including Baidu Drive, Quark Drive, and Google Drive.'
  },
  {
    q_zh: '可以随时取消订阅吗？',
    a_zh: '是的，您可以随时通过 PayPal 账户取消订阅，取消后您的会员权益将持续到当前计费周期结束。',
    q_en: 'Can I cancel my subscription at any time?',
    a_en: 'Yes, you can cancel your subscription anytime through your PayPal account. Your membership benefits will continue until the end of the current billing period.'
  },
  {
    q_zh: '支持哪些支付方式？',
    a_zh: '目前支持 PayPal 支付，支持信用卡、借记卡及 PayPal 账户余额。',
    q_en: 'What payment methods are supported?',
    a_en: 'We currently support PayPal payments, which accepts credit cards, debit cards, and PayPal account balance.'
  },
  {
    q_zh: '壁纸可以商用吗？',
    a_zh: '所有壁纸仅供个人使用，不允许用于商业用途、再销售或二次分发。',
    q_en: 'Can wallpapers be used commercially?',
    a_en: 'All wallpapers are for personal use only. Commercial use, resale, or redistribution is not permitted.'
  }
]

function formatDate(dateStr) {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString(locale.value === 'zh' ? 'zh-CN' : 'en-US')
}

function toggleFaq(i) {
  openFaq.value = openFaq.value === i ? null : i
}

async function handleSubscribe(plan) {
  selectedPlan.value = plan
  showPayPalModal.value = true
}

async function handleTestSubscribe() {
  if (subscribing.value || !selectedPlan.value) return
  subscribing.value = true
  try {
    const res = await paypalAPI.createOrder({
      order_type: 'membership',
      membership_type: selectedPlan.value.type,
      is_test: true,
      user_email: auth.user?.email || 'user@wallpaper.com'
    })
    if (res.token) {
      localStorage.setItem('token', res.token)
    }
    if (res.user) {
      auth.user = res.user
    } else {
      // Refresh current user
      await auth.fetchMe()
    }
    success(locale.value === 'zh' ? '👑 成功开通 VIP 会员！全场壁纸已免单！' : 'VIP Membership Activated!')
    showPayPalModal.value = false
    router.push('/')
  } catch (err) {
    error(err.message || (locale.value === 'zh' ? '开通失败' : 'Failed to activate'))
  } finally {
    subscribing.value = false
  }
}

async function initPayPalOrder(plan) {
  paypalLoading.value = true
  await nextTick()
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
    if (!paypal || !paypalSubContainer.value) return

    paypalSubContainer.value.innerHTML = ''
    await paypal.Buttons({
      style: {
        layout: 'vertical',
        color: 'gold',
        shape: 'pill',
        label: 'pay',
        height: 44
      },
      createOrder: async () => {
        subscribing.value = true
        try {
          const res = await paypalAPI.createOrder({
            order_type: 'membership',
            membership_type: plan.type,
            user_email: auth.user?.email
          })
          if (res.token) localStorage.setItem('token', res.token)
          return res.paypal_order_id || res.id
        } catch (e) {
          error(e.message || t('payment.failed'))
          throw e
        } finally {
          subscribing.value = false
        }
      },
      onApprove: async (data) => {
        try {
          const res = await paypalAPI.captureOrder(data.orderID)
          if (res.token) localStorage.setItem('token', res.token)
          if (res.user) auth.user = res.user
          else await auth.fetchMe()
          success(locale.value === 'zh' ? '👑 恭喜成为 VIP 会员！' : 'Welcome to VIP!')
          showPayPalModal.value = false
          router.push('/')
        } catch (e) {
          error(t('payment.failed'))
        }
      },
      onError: (err) => {
        console.error('PayPal error:', err)
        error(t('payment.failed'))
      }
    }).render(paypalSubContainer.value)
  } catch (e) {
    console.warn('PayPal load warning:', e)
  } finally {
    paypalLoading.value = false
  }
}
</script>

<style scoped>
.membership-page {
  min-height: 100vh;
}

/* Hero */
.membership-hero {
  padding: 80px 0 60px;
  text-align: center;
  background: var(--g-hero);
}

.hero-label {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: var(--c-gold-light);
  color: var(--c-gold);
  padding: 6px 16px;
  border-radius: var(--r-full);
  font-size: 0.82rem;
  font-weight: 600;
  margin-bottom: 20px;
}

.membership-hero h1 { margin-bottom: 12px; }
.membership-hero p { font-size: 1.1rem; color: var(--c-text-secondary); }

.current-plan-alert {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  margin-top: 20px;
  padding: 10px 20px;
  background: rgba(52,199,89,0.1);
  border: 1px solid rgba(52,199,89,0.3);
  border-radius: var(--r-full);
  font-size: 0.88rem;
  color: var(--c-success);
}

/* Plans */
.plans-section {
  padding: 60px 0 80px;
}

.plans-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 20px;
  margin-bottom: 40px;
}

@media (max-width: 768px) {
  .membership-hero {
    padding: 40px 0 24px !important;
  }
  .membership-hero h1 {
    font-size: 1.85rem !important;
  }
  .plans-section {
    padding: 24px 0 50px !important;
  }
  .plans-grid {
    grid-template-columns: 1fr !important;
    max-width: 100% !important;
    gap: 16px !important;
    margin: 0 auto 30px !important;
  }
  .plan-card {
    padding: 24px 20px !important;
    border-radius: 20px !important;
    gap: 18px !important;
  }
  .plan-card.featured {
    transform: none !important;
  }
  .plan-card.featured:hover {
    transform: none !important;
  }
}

.plan-card {
  background: var(--c-surface);
  border: 2px solid var(--c-border-light);
  border-radius: var(--r-2xl);
  padding: 32px 28px;
  position: relative;
  transition: all var(--t-normal);
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.plan-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-4px);
}

.plan-card.featured {
  border-color: var(--c-accent);
  box-shadow: 0 0 0 1px var(--c-accent), var(--shadow-md);
  transform: scale(1.03);
}

.plan-card.featured:hover {
  transform: scale(1.03) translateY(-4px);
}

.plan-card.current {
  border-color: var(--c-success);
}

.plan-badge {
  position: absolute;
  top: -14px;
  left: 50%;
  transform: translateX(-50%);
  background: var(--c-accent);
  color: white;
  padding: 4px 16px;
  border-radius: var(--r-full);
  font-size: 0.75rem;
  font-weight: 700;
  white-space: nowrap;
}

.current-badge { background: var(--c-success); }

.plan-icon { font-size: 2rem; margin-bottom: 8px; }

.plan-name {
  font-size: 1.2rem;
  font-weight: 700;
  margin-bottom: 12px;
}

.plan-price {
  display: flex;
  align-items: baseline;
  gap: 2px;
  margin-bottom: 4px;
}

.price-currency {
  font-size: 1.2rem;
  font-weight: 700;
  color: var(--c-text-secondary);
}

.price-num {
  font-size: 2.5rem;
  font-weight: 900;
  letter-spacing: -0.04em;
  color: var(--c-text-primary);
}

.price-period {
  font-size: 0.85rem;
  color: var(--c-text-secondary);
  margin-left: 4px;
}

.plan-save {
  display: inline-block;
  padding: 3px 10px;
  background: var(--c-gold-light);
  color: var(--c-gold);
  border-radius: var(--r-full);
  font-size: 0.75rem;
  font-weight: 700;
}

.plan-features {
  display: flex;
  flex-direction: column;
  gap: 10px;
  flex: 1;
}

.feature-item {
  display: flex;
  gap: 10px;
  align-items: flex-start;
  font-size: 0.88rem;
  color: var(--c-text-secondary);
}

.feature-check {
  color: var(--c-success);
  font-weight: 700;
  flex-shrink: 0;
}

.plan-btn {
  width: 100%;
  padding: 14px;
  font-size: 0.95rem;
}

/* Payment info */
.payment-info {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  font-size: 0.85rem;
  color: var(--c-text-tertiary);
}

.paypal-badge {
  background: #003087;
  color: #00aeef;
  padding: 3px 10px;
  border-radius: var(--r-sm);
  font-size: 0.78rem;
  font-weight: 800;
  letter-spacing: 0.02em;
}

/* FAQ */
.faq-section {
  padding: 60px 0 80px;
  background: var(--c-bg);
}

.faq-section h2 {
  text-align: center;
  margin-bottom: 40px;
}

.faq-list {
  max-width: 720px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.faq-item {
  background: var(--c-surface);
  border-radius: var(--r-lg);
  overflow: hidden;
}

.faq-q {
  width: 100%;
  padding: 18px 20px;
  background: none;
  border: none;
  text-align: left;
  font-size: 0.95rem;
  font-weight: 600;
  color: var(--c-text-primary);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  font-family: var(--font-sans);
  transition: background var(--t-fast);
}
.faq-q:hover { background: var(--c-bg); }

.faq-arrow {
  font-size: 1.3rem;
  color: var(--c-text-tertiary);
  transition: transform var(--t-fast);
  flex-shrink: 0;
}
.faq-arrow.open { transform: rotate(90deg); }

.faq-a {
  padding: 0 20px 18px;
  font-size: 0.9rem;
  color: var(--c-text-secondary);
  line-height: 1.7;
}

/* Subscribe modal */
.close-btn {
  width: 32px;
  height: 32px;
  border-radius: var(--r-sm);
  background: var(--c-bg);
  border: none;
  cursor: pointer;
  font-size: 0.9rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.subscribe-summary {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: var(--c-gold-light);
  border-radius: var(--r-md);
  margin-bottom: 20px;
}
.sub-icon { font-size: 2rem; }
.sub-plan { font-weight: 700; color: var(--c-text-primary); margin-bottom: 4px; }
.sub-price { font-size: 0.88rem; color: var(--c-text-secondary); }

.payment-note {
  margin-top: 16px;
  font-size: 0.78rem;
  color: var(--c-text-tertiary);
  text-align: center;
}

.sub-modal-sheet {
  max-width: 480px;
  border-radius: 24px;
  overflow: hidden;
}

.show-mobile-handle {
  display: none;
}

@media (max-width: 640px) {
  .modal-overlay.mobile-sheet-overlay {
    align-items: flex-end !important;
    padding: 0 !important;
  }

  .sub-modal-sheet {
    max-width: 100% !important;
    border-radius: 26px 26px 0 0 !important;
    max-height: 90vh !important;
    overflow-y: auto !important;
    animation: sheetSlideUp 0.32s cubic-bezier(0.16, 1, 0.3, 1) !important;
    margin: 0 !important;
    width: 100% !important;
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
    margin: 10px auto 4px;
  }

  .modal-header {
    padding: 12px 18px 8px !important;
  }

  .modal-body {
    padding: 14px 18px max(24px, env(safe-area-inset-bottom)) !important;
  }

  .tab-btn {
    min-height: 42px;
  }
}
</style>
