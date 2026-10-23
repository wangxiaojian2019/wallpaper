<template>
  <view v-if="visible" class="cashier-overlay" @tap.self="close">
    <view class="cashier-sheet" :class="[appStore.theme]">
      <!-- Drag Handle -->
      <view class="drag-bar"></view>

      <!-- Header -->
      <view class="sheet-header">
        <view class="title-row">
          <text class="sheet-title">{{ isSubmittedReview ? (appStore.lang === 'zh' ? '⏳ 支付凭据已提交审核' : '⏳ Payment Under Review') : (showSecondaryConfirm ? (appStore.lang === 'zh' ? '📱 官方收款码直接付款' : '📱 Direct Payment') : (appStore.lang === 'zh' ? '获取高清壁纸原图' : 'Unlock Original Wallpaper')) }}</text>
          <text class="sheet-sub">{{ isSubmittedReview ? (appStore.lang === 'zh' ? '管理员核对账单后将为您审核下发，可凭邮箱查询进度' : 'Admin is verifying payment. Links will be unlocked once approved.') : (showSecondaryConfirm ? (appStore.lang === 'zh' ? '请使用手机扫描下方收款码付款，完成后点击按钮提交审核' : 'Scan QR code and tap button below to submit') : (appStore.lang === 'zh' ? '解锁百度网盘 / 夸克网盘 / 谷歌网盘原画高速直链' : 'Direct high-speed Cloud Drive & original 4K/8K downloads')) }}</text>
        </view>
        <view class="close-btn" @tap="close">
          <text class="close-icon">✕</text>
        </view>
      </view>

      <!-- Submitted Review Status View (待审核状态卡片) -->
      <view v-if="isSubmittedReview" class="submitted-review-view">
        <view class="review-status-card">
          <view class="review-icon-wrap">
            <text class="review-icon">⏳</text>
          </view>
          <text class="review-title">{{ appStore.lang === 'zh' ? '付款凭据已提交，等待人工审核' : 'Payment Submitted for Review' }}</text>
          <text class="review-desc">
            {{ appStore.lang === 'zh' ? '微信 / 支付宝收款码需人工对账。管理员核验转账记录与预留邮箱后，将立即为您审核通过并解锁！' : 'Manual verification is required for direct QR payments. Downloads will be unlocked once approved.' }}
          </text>

          <!-- Order Credential Box -->
          <view class="review-credential-box">
            <view class="review-row">
              <text class="r-label">{{ appStore.lang === 'zh' ? '订单编号：' : 'Order ID: ' }}</text>
              <view class="r-val-wrap" @tap="copyToClipboard(pendingOrder?.order_id, appStore.lang === 'zh' ? '订单号已复制' : 'Order ID copied')">
                <text class="r-val font-mono font-bold">{{ pendingOrder?.order_id }}</text>
                <text class="r-copy">📋 复制</text>
              </view>
            </view>
            <view class="review-row">
              <text class="r-label">{{ appStore.lang === 'zh' ? '购买项目：' : 'Product: ' }}</text>
              <text class="r-val font-bold">{{ selectedPlan === 'single' ? cleanProductTitle : currentPlanName }}</text>
            </view>
            <view class="review-row">
              <text class="r-label">{{ appStore.lang === 'zh' ? '支付金额：' : 'Amount: ' }}</text>
              <text class="r-val font-bold text-amber">¥{{ currencySymbol === '¥' ? currentPlanPrice : exchangeRatePrice }}</text>
            </view>
            <view class="review-row">
              <text class="r-label">{{ appStore.lang === 'zh' ? '交付邮箱：' : 'Email: ' }}</text>
              <text class="r-val font-bold">{{ email }}</text>
            </view>
            <view class="review-row">
              <text class="r-label">{{ appStore.lang === 'zh' ? '审核时效：' : 'ETA: ' }}</text>
              <text class="r-val text-emerald font-bold">{{ appStore.lang === 'zh' ? '通常 1 ~ 5 分钟内完成' : 'Usually within 1-5 mins' }}</text>
            </view>
          </view>

          <!-- Action Buttons -->
          <view class="review-actions-col">
            <button class="btn-lookup-progress" @tap="goToLookup">
              <text class="btn-icon">🔍</text>
              <text class="btn-text">{{ appStore.lang === 'zh' ? '去订单中心查询审核进度与提取码' : 'Check Order & Download Progress' }}</text>
            </button>
            <button class="btn-copy-kefu" @tap="copyToClipboard(pendingOrder?.order_id, appStore.lang === 'zh' ? '订单号已复制，可发给客服' : 'Order ID copied')">
              <text class="btn-icon">📋</text>
              <text class="btn-text">{{ appStore.lang === 'zh' ? '复制订单号（可发给客服加快审核）' : 'Copy Order ID for Support' }}</text>
            </button>
            <button class="btn-close-review" @tap="close">
              <text class="btn-text">{{ appStore.lang === 'zh' ? '完成并关闭' : 'Done & Close' }}</text>
            </button>
          </view>
        </view>
      </view>

      <!-- Primary Plan & Payment Method Selection View -->
      <view v-else-if="!showSecondaryConfirm" class="cashier-primary-view">
        <!-- Member Registration Incentive Banner (Only shown if NOT logged in) -->
        <view v-if="!userStore.isLoggedIn" class="guest-register-incentive">
          <view class="incentive-left">
            <text class="incentive-crown">👑</text>
            <view class="incentive-text-wrap">
              <text class="incentive-title">{{ appStore.lang === 'zh' ? '强烈建议注册/登录账号购买' : 'Recommended: Register an Account' }}</text>
              <text class="incentive-sub">{{ appStore.lang === 'zh' ? '已购壁纸将永久保存在【我的已购】，换设备或清理缓存绝不丢失；游客易丢失下发链接。' : 'Purchases saved permanently to cloud. Guest orders risk losing access upon cache clear.' }}</text>
            </view>
          </view>
          <button class="btn-incentive-reg" @tap="showAuthModal = true">
            <text class="btn-reg-text">{{ appStore.lang === 'zh' ? '免费注册' : 'Register' }}</text>
          </button>
        </view>

        <!-- Wallpaper Mini Info -->
        <view class="product-strip">
          <image
            v-if="product.cover_image"
            :src="product.cover_image"
            mode="aspectFill"
            class="strip-thumb"
          />
          <view v-else class="strip-thumb strip-placeholder"></view>
          <view class="strip-info">
            <text class="strip-title">{{ cleanProductTitle }}</text>
            <text class="strip-tag">{{ appStore.lang === 'zh' ? '4K/8K 超高清 · 永久有效 · 极速网盘直下' : '4K/8K UHD · Lifetime Validity · Fast Cloud Drive' }}</text>
          </view>
        </view>

        <!-- Plan Selector (All 5 Options: Single, Monthly, Quarterly, Yearly, Lifetime) -->
        <view class="plan-options">
          <!-- Plan 1: Single Buy -->
          <view
            class="plan-card plan-single-card"
            :class="{ selected: selectedPlan === 'single' }"
            @click="selectedPlan = 'single'"
          >
            <view class="plan-radio">
              <view v-if="selectedPlan === 'single'" class="radio-inner"></view>
            </view>
            <view class="plan-body">
              <view class="plan-title-line">
                <text class="plan-name">{{ appStore.lang === 'zh' ? '单张永久解锁' : 'Single Wallpaper Unlock' }}</text>
                <text class="plan-badge-mini">{{ appStore.lang === 'zh' ? '当前单款' : 'Current' }}</text>
              </view>
              <text class="plan-desc">{{ appStore.lang === 'zh' ? '仅本款壁纸 4K/8K 超清原图与动态源文件永久下载' : 'Lifetime download access for this artwork' }}</text>
            </view>
            <view class="plan-price-wrap">
              <text class="plan-symbol">{{ currencySymbol }}</text>
              <text class="plan-price">{{ singlePrice }}</text>
            </view>
          </view>

          <!-- VIP Section Header -->
          <view v-if="availableVipTiers.length > 0" class="vip-section-divider">
            <view class="divider-line"></view>
            <view class="divider-label">
              <text class="divider-crown">👑</text>
              <text class="divider-text">{{ appStore.lang === 'zh' ? '开通 VIP 会员 · 本壁纸免费享' : 'VIP Membership · Free Access' }}</text>
            </view>
            <view class="divider-line"></view>
          </view>

          <!-- VIP Dynamic Grid: Only displays tiers allowed for this wallpaper -->
          <view v-if="availableVipTiers.length > 0" class="vip-grid" :class="{ 'grid-three': availableVipTiers.length === 3, 'grid-single': availableVipTiers.length === 1 }">
            <view
              v-for="tier in availableVipTiers"
              :key="tier.type"
              class="vip-tier-card"
              :class="{ selected: selectedPlan === tier.type, featured: tier.isFeatured }"
              @click="selectedPlan = tier.type"
            >
              <view v-if="tier.isFeatured" class="featured-ribbon">
                <text class="ribbon-text">{{ tier.ribbonText }}</text>
              </view>
              <view class="vip-tier-header">
                <text class="vip-tier-name">{{ tier.name }}</text>
                <text class="tier-pill-small" :class="tier.pillClass">{{ tier.days }}</text>
              </view>
              <view class="vip-tier-price-row">
                <text class="vip-cur" :class="{ 'gold-cur': tier.isFeatured }">{{ currencySymbol }}</text>
                <text class="vip-price-num" :class="{ 'gold-price': tier.isFeatured }">{{ tier.price }}</text>
              </view>
              <text class="vip-tier-hint" :class="{ 'gold-hint': tier.isFeatured }">{{ tier.hint }}</text>
            </view>
          </view>

          <!-- Notice when no VIP tier is allowed (exclusive individual buy) -->
          <view v-else class="vip-exclusive-box">
            <text class="vip-exclusive-icon">💎</text>
            <text class="vip-exclusive-text">
              {{ appStore.lang === 'zh' ? '本款壁纸已设为专属单张解锁，不参与 VIP 免单' : 'This artwork is single unlock only and does not participate in VIP free downloads.' }}
            </text>
          </view>
        </view>

        <!-- Delivery / Account Target -->
        <view class="delivery-box">
          <view class="delivery-label-row">
            <text class="delivery-label">
              {{ userStore.isLoggedIn ? (appStore.lang === 'zh' ? '绑定交付账号' : 'Account Email') : (appStore.lang === 'zh' ? '接收邮箱 (订单凭证与发货链接)' : 'Delivery Email') }}
            </text>
            <text v-if="!userStore.isLoggedIn" class="quick-email-btn" @click="generateGuestEmail">
              🎲 {{ appStore.lang === 'zh' ? '匿名生成凭据' : 'Random Guest' }}
            </text>
          </view>
          <input
            v-model="email"
            type="text"
            :disabled="userStore.isLoggedIn"
            :placeholder="appStore.lang === 'zh' ? '请输入用于接收和查询订单的邮箱' : 'Enter email for order confirmation'"
            class="delivery-input"
          />
        </view>

        <!-- Redeem Code Entrance -->
        <view class="redeem-row" @click="showRedeemInput = !showRedeemInput">
          <text class="redeem-toggle">🎫 {{ appStore.lang === 'zh' ? '有兑换口令或卡密？点击在此输入' : 'Have a redeem code? Tap here' }}</text>
        </view>
        <view v-if="showRedeemInput" class="redeem-box">
          <input
            v-model="redeemCode"
            :placeholder="appStore.lang === 'zh' ? '请输入兑换口令 / 卡密' : 'Enter redeem code'"
            class="delivery-input"
          />
          <button class="redeem-btn" :disabled="!redeemCode" @click="handleRedeem">
            {{ appStore.lang === 'zh' ? '立即兑换' : 'Redeem' }}
          </button>
        </view>

        <!-- Dynamic Payment Channels Selector (Mobile Friendly) -->
        <view class="cashier-pay-methods">
          <view class="channel-header-row">
            <text class="c-section-label">{{ appStore.lang === 'zh' ? '支付方式选择：' : 'Select Payment Method:' }}</text>
            <text class="geo-region-tag">
              {{ appStore.regionScope === 'china' ? '🇨🇳 中国通道' : '🌐 国际通道' }}
            </text>
          </view>
          <view class="c-methods-grid">
            <view
              v-for="ch in activeChannels"
              :key="ch.code"
              class="c-method-chip"
              :class="{ active: selectedChannel === ch.code, 'test-chip': ch.code === 'test' }"
              @click="selectedChannel = ch.code"
            >
              <text class="c-method-icon">{{ ch.icon || getChannelDefaultIcon(ch.code) }}</text>
              <view class="c-method-text-wrap">
                <text class="c-method-title">{{ appStore.lang === 'zh' ? ch.name : (ch.name_en || ch.name) }}</text>
                <text class="c-method-sub">{{ getChannelSubtitle(ch) }}</text>
              </view>
              <view class="c-method-radio">
                <view v-if="selectedChannel === ch.code" class="c-radio-dot"></view>
              </view>
            </view>
          </view>
        </view>

        <!-- Channel Direct Pay Hint Card (If WeChat/Alipay or QR channel is selected) -->
        <view v-if="selectedChannel !== 'test' && selectedChannel !== 'paypal'" class="direct-pay-hint-card">
          <view class="hint-card-header">
            <view class="hint-badge-left">
              <text class="hint-badge-icon">{{ currentChannelInfo?.icon || '📱' }}</text>
              <text class="hint-badge-title">{{ currentChannelInfo?.name || (appStore.lang === 'zh' ? '直接付款通道' : 'Direct Payment') }}</text>
            </view>
            <text class="hint-badge-tag">{{ currentQrCodeUrl ? (appStore.lang === 'zh' ? '📷 收款码已就绪' : 'QR Ready') : (appStore.lang === 'zh' ? '⚡ 极速对账' : 'Instant') }}</text>
          </view>
          
          <view class="hint-card-content">
            <image
              v-if="currentQrCodeUrl"
              :src="currentQrCodeUrl"
              mode="aspectFit"
              class="hint-qr-thumb"
              @tap="handlePreviewQR(currentQrCodeUrl)"
            />
            <view class="hint-desc-wrap">
              <text class="hint-desc-main">{{ appStore.lang === 'zh' ? '支持微信 / 支付宝扫码直接付款' : 'Direct scan payment supported' }}</text>
              <text class="hint-desc-sub">{{ currentChannelInfo?.account_info || (appStore.lang === 'zh' ? '点击下方【进入直接付款】将展示专属收款码，付款后一键极速自动解锁！' : 'Tap button below to view QR code and unlock.') }}</text>
            </view>
          </view>
        </view>

        <!-- Checkout Action Button -->
        <view class="sheet-action">
          <button
            class="btn-pay"
            :class="{
              'btn-pay-vip': selectedPlan !== 'single',
              'btn-pay-direct': selectedChannel === 'wechat' || selectedChannel === 'alipay' || selectedChannel === 'custom'
            }"
            :loading="submitting"
            @click="handlePay"
            @tap="handlePay"
          >
            <text class="btn-pay-icon">{{ selectedChannel === 'test' ? '⚡' : (selectedChannel === 'wechat' || selectedChannel === 'alipay' ? '📱' : (selectedPlan !== 'single' ? '👑' : '💳')) }}</text>
            <text class="btn-pay-text">{{ payButtonText }}</text>
          </button>

          <!-- Quick Secondary Confirmation Entrance for Mobile Visitors (Only for WeChat, Alipay, Custom QR) -->
          <view v-if="selectedChannel === 'wechat' || selectedChannel === 'alipay' || selectedChannel === 'custom'" class="manual-confirm-entrance" @click="handlePay" @tap="handlePay">
            <text class="manual-confirm-text">
              {{ appStore.lang === 'zh' ? '📱 点击进入直接付款界面（扫码与核验解锁） →' : 'Tap to enter Direct Payment & Unlock Screen →' }}
            </text>
          </view>

          <text class="security-tip">
            🔒 {{ appStore.lang === 'zh' ? '付款成功后立即显示百度/夸克/谷歌网盘链接与提取码' : 'Instant Cloud Drive download link provided upon payment completion' }}
          </text>
        </view>
      </view>

      <!-- Direct Payment View (官方收款码直接付款与核验界面) -->
      <view v-else class="secondary-confirm-view">
        <view class="confirm-card">
          <!-- Top Header: Channel Badge + Order ID with Copy -->
          <view class="confirm-header">
            <view class="confirm-badge">
              <text class="confirm-badge-icon">{{ currentChannelInfo?.icon || '💳' }}</text>
              <text class="confirm-badge-text">{{ currentChannelInfo?.name || pendingOrder?.channel_name || (appStore.lang === 'zh' ? '官方收款码直接付款' : 'Direct QR Payment') }}</text>
            </view>
            <view class="order-no-wrap" @tap="copyToClipboard(pendingOrder?.order_id, appStore.lang === 'zh' ? '订单号已复制' : 'Order ID copied')">
              <text class="order-no-text">{{ appStore.lang === 'zh' ? '订单: ' : 'Order: ' }}{{ pendingOrder?.order_id?.slice(0, 14) }}...</text>
              <text class="copy-tag">📋 {{ appStore.lang === 'zh' ? '复制' : 'Copy' }}</text>
            </view>
          </view>

          <!-- HERO: Direct Payment QR Code Card -->
          <view class="direct-pay-hero-box">
            <!-- Amount Display Tag -->
            <view class="hero-amount-tag">
              <text class="hero-amount-label">{{ appStore.lang === 'zh' ? '需付金额' : 'Amount to Pay' }}</text>
              <view class="hero-amount-val">
                <text class="hero-cur">¥</text>
                <text class="hero-num">{{ currencySymbol === '¥' ? currentPlanPrice : exchangeRatePrice }}</text>
                <text v-if="currencySymbol === '$'" class="hero-sub">({{ currencySymbol }}{{ currentPlanPrice }})</text>
              </view>
            </view>

            <!-- QR Code Image or Fallback -->
            <view v-if="currentQrCodeUrl" class="qr-container-card">
              <view class="qr-img-frame" @tap="handlePreviewQR(currentQrCodeUrl)">
                <image
                  :src="currentQrCodeUrl"
                  mode="aspectFit"
                  class="direct-qr-image"
                />
              </view>
              <view class="qr-scan-tips" @tap="handlePreviewQR(currentQrCodeUrl)">
                <text class="qr-scan-icon">🔍</text>
                <text class="qr-scan-txt">{{ appStore.lang === 'zh' ? '点击可放大收款码 · 手机端支持长按或截屏识别' : 'Tap to enlarge · Long press or screenshot to scan' }}</text>
              </view>
            </view>

            <view v-else class="qr-fallback-card">
              <text class="fallback-icon">{{ currentChannelInfo?.icon || '📱' }}</text>
              <text class="fallback-title">{{ currentChannelInfo?.name || '官方收款通道' }}</text>
              <text class="fallback-desc">{{ pendingOrder?.account_info || currentChannelInfo?.account_info || (appStore.lang === 'zh' ? '请使用对应 App 向官方账户转账完成付款' : 'Please transfer to official account') }}</text>
            </view>

            <!-- Transfer Memo Reminder -->
            <view class="transfer-memo-card">
              <view class="memo-title-line">
                <text class="memo-icon">💡</text>
                <text class="memo-title">{{ appStore.lang === 'zh' ? '付款转账备注提示 (极速自动对账)' : 'Payment Note (Fast Match)' }}</text>
              </view>
              <view class="memo-content-line">
                <text class="memo-desc">{{ appStore.lang === 'zh' ? '转账备注栏请填写：' : 'Note email in transfer: ' }}</text>
                <text class="memo-email" @tap="copyToClipboard(email, appStore.lang === 'zh' ? '邮箱已复制' : 'Email copied')">{{ email }}</text>
                <text class="memo-copy-btn" @tap="copyToClipboard(email, appStore.lang === 'zh' ? '邮箱已复制' : 'Email copied')">📋 {{ appStore.lang === 'zh' ? '复制' : 'Copy' }}</text>
              </view>
            </view>
          </view>

          <!-- Order Summary Details -->
          <view class="confirm-summary">
            <view class="summary-line">
              <text class="s-label">{{ appStore.lang === 'zh' ? '购买项目' : 'Product' }}</text>
              <text class="s-val font-bold">{{ selectedPlan === 'single' ? cleanProductTitle : currentPlanName }}</text>
            </view>
            <view class="summary-line">
              <text class="s-label">{{ appStore.lang === 'zh' ? '支付方式' : 'Payment Method' }}</text>
              <text class="s-val channel-val">{{ currentChannelInfo?.name || pendingOrder?.channel_name }}</text>
            </view>
            <view class="summary-line">
              <text class="s-label">{{ appStore.lang === 'zh' ? '交付邮箱' : 'Delivery Email' }}</text>
              <text class="s-val text-email">{{ email }}</text>
            </view>
          </view>

          <!-- Action Buttons -->
          <view class="confirm-actions">
            <button
              class="btn-confirm-pay"
              :class="{ 'btn-confirm-pay-review': selectedChannel !== 'test' }"
              :loading="confirmingPayment"
              :disabled="confirmingPayment"
              @tap="handleConfirmPaid"
            >
              <text class="btn-confirm-icon">{{ selectedChannel === 'test' ? '⚡' : '📤' }}</text>
              <text class="btn-confirm-text">
                {{ confirmingPayment ? (appStore.lang === 'zh' ? '正在提交凭据...' : 'Submitting Proof...') : (selectedChannel === 'test' ? (appStore.lang === 'zh' ? '⚡ 模拟测试支付（秒级免密开通）' : 'Simulate Sandbox Instant Pay') : (appStore.lang === 'zh' ? '我已扫码付款，提交核验 / 审核' : 'I Have Paid, Submit for Review')) }}
              </text>
            </button>
            <button class="btn-cancel-confirm" :disabled="confirmingPayment" @tap="showSecondaryConfirm = false">
              <text class="btn-cancel-text">{{ appStore.lang === 'zh' ? '← 返回更换支付方式 / 方案' : '← Back to Payment Methods' }}</text>
            </button>
          </view>
        </view>
      </view>
    </view>

    <!-- Embedded Auth Modal for In-Cashier Registration -->
    <AuthModal
      :visible="showAuthModal"
      initial-mode="register"
      :initial-email="email"
      @close="showAuthModal = false"
      @success="handleAuthSuccess"
    />

    <!-- Dedicated Full-Screen QR Code Magnifier Lightbox Modal (解决二维码放大在弹窗后面的问题) -->
    <view v-if="showQRLightbox" class="qr-lightbox-overlay" @tap="showQRLightbox = false">
      <view class="qr-lightbox-dialog" @tap.stop>
        <view class="lightbox-close-btn" @tap="showQRLightbox = false">
          <text class="lightbox-close-icon">✕</text>
        </view>
        <view class="lightbox-channel-header">
          <text class="lightbox-channel-icon">{{ currentChannelInfo?.icon || '📱' }}</text>
          <text class="lightbox-channel-name">{{ currentChannelInfo?.name || '官方收款码' }}</text>
          <text class="lightbox-badge-direct">{{ appStore.lang === 'zh' ? '扫码直接付款' : 'Scan to Pay' }}</text>
        </view>
        <view class="lightbox-img-wrap">
          <image
            v-if="currentQrCodeUrl"
            :src="currentQrCodeUrl"
            mode="aspectFit"
            class="lightbox-img"
          />
          <view v-else class="lightbox-img-ph">
            <text style="font-size: 32px;">📱</text>
            <text style="font-size: 12px; opacity: 0.6; margin-top: 6px;">{{ appStore.lang === 'zh' ? '暂未配置收款码' : 'No QR Code' }}</text>
          </view>
        </view>
        <view class="lightbox-amount-pill">
          <text class="lightbox-amount-label">{{ appStore.lang === 'zh' ? '请扫码支付金额：' : 'Amount: ' }}</text>
          <text class="lightbox-amount-val">¥{{ currencySymbol === '¥' ? currentPlanPrice : exchangeRatePrice }}</text>
          <text v-if="currencySymbol === '$'" class="lightbox-amount-sub">({{ currencySymbol }}{{ currentPlanPrice }})</text>
        </view>
        <text class="lightbox-hint">{{ appStore.lang === 'zh' ? '💡 支持微信/支付宝扫一扫 · 手机端长按可保存至相册' : 'Scan with WeChat or Alipay. Long press to save image.' }}</text>
        <button class="lightbox-btn-done" @tap="showQRLightbox = false">
          {{ appStore.lang === 'zh' ? '✓ 我已扫码，返回核验' : 'I Have Scanned, Return' }}
        </button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useUserStore } from '@/store/user'
import { useAppStore } from '@/store/app'
import { productAPI, paypalAPI } from '@/api'
import request from '@/utils/request'
import AuthModal from '@/components/AuthModal.vue'

const props = defineProps({
  visible: Boolean,
  product: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['close', 'success'])

const userStore = useUserStore()
const appStore = useAppStore()

const selectedPlan = ref('quarterly')
const selectedChannel = ref('test')
const email = ref(userStore.userInfo?.email || '')
const showRedeemInput = ref(false)
const redeemCode = ref('')
const submitting = ref(false)
const showAuthModal = ref(false)
const fetchedChannels = ref([])
const showQRLightbox = ref(false)

// Secondary Payment Confirmation View States
const showSecondaryConfirm = ref(false)
const isSubmittedReview = ref(false)
const pendingOrder = ref(null)
const confirmingPayment = ref(false)

// Pricing Computations
const currencySymbol = computed(() => (appStore.currency === 'USD' ? '$' : '¥'))

const singlePrice = computed(() => {
  return appStore.getProductPriceNumber(props.product).toFixed(2)
})

const monthlyPrice = computed(() => (appStore.currency === 'USD' ? '9.99' : '19.9'))
const quarterlyPrice = computed(() => (appStore.currency === 'USD' ? '24.99' : '49.9'))
const yearlyPrice = computed(() => (appStore.currency === 'USD' ? '59.99' : '129.0'))
const lifetimePrice = computed(() => (appStore.currency === 'USD' ? '99.99' : '199.0'))

// VIP Membership Plans dynamically filtered by props.product.allowed_vip_levels
const availableVipTiers = computed(() => {
  const allTiers = [
    {
      type: 'monthly',
      name: appStore.lang === 'zh' ? '月度 VIP' : 'Monthly VIP',
      days: appStore.lang === 'zh' ? '30天' : '30d',
      price: monthlyPrice.value,
      hint: appStore.lang === 'zh' ? '全场畅下 · 尝鲜体验' : 'All-access pass',
      pillClass: '',
      isFeatured: false,
      ribbonText: ''
    },
    {
      type: 'quarterly',
      name: appStore.lang === 'zh' ? '季度 VIP' : 'Quarterly VIP',
      days: appStore.lang === 'zh' ? '90天' : '90d',
      price: quarterlyPrice.value,
      hint: appStore.lang === 'zh' ? '85% 会员首选 · 省17%' : 'Best choice · Save 17%',
      pillClass: 'gold',
      isFeatured: true,
      ribbonText: appStore.lang === 'zh' ? '🔥 推荐' : '🔥 Popular'
    },
    {
      type: 'yearly',
      name: appStore.lang === 'zh' ? '年度 VIP' : 'Yearly VIP',
      days: appStore.lang === 'zh' ? '365天' : '365d',
      price: yearlyPrice.value,
      hint: appStore.lang === 'zh' ? '立省 33% · 全年无忧' : 'Save 33% all year',
      pillClass: '',
      isFeatured: false,
      ribbonText: ''
    },
    {
      type: 'lifetime',
      name: appStore.lang === 'zh' ? '终生尊享' : 'Lifetime VIP',
      days: appStore.lang === 'zh' ? '永久' : 'Life',
      price: lifetimePrice.value,
      hint: appStore.lang === 'zh' ? '买断特权 · 终身免单' : 'Permanent access',
      pillClass: 'diamond',
      isFeatured: false,
      ribbonText: ''
    }
  ]

  const raw = props.product?.allowed_vip_levels
  if (raw === undefined || raw === null || raw === '' || raw === 'null') {
    return allTiers
  }

  let list = []
  try {
    list = typeof raw === 'string' ? JSON.parse(raw) : raw
  } catch (e) {
    list = ['monthly', 'quarterly', 'yearly', 'lifetime']
  }

  if (!Array.isArray(list)) {
    return allTiers
  }

  return allTiers.filter(t => list.includes(t.type))
})

const currentPlanPrice = computed(() => {
  if (selectedPlan.value === 'single') return singlePrice.value
  const found = availableVipTiers.value.find(t => t.type === selectedPlan.value)
  if (found) return found.price
  switch (selectedPlan.value) {
    case 'monthly': return monthlyPrice.value
    case 'quarterly': return quarterlyPrice.value
    case 'yearly': return yearlyPrice.value
    case 'lifetime': return lifetimePrice.value
    default: return singlePrice.value
  }
})

const exchangeRatePrice = computed(() => {
  const rate = Number(appStore.exchangeRate || 7.25)
  const price = Number(currentPlanPrice.value) || 0
  return (price * rate).toFixed(2)
})

const currentPlanName = computed(() => {
  if (selectedPlan.value === 'single') return appStore.lang === 'zh' ? '单张解锁' : 'Single Unlock'
  const found = availableVipTiers.value.find(t => t.type === selectedPlan.value)
  if (found) return found.name
  switch (selectedPlan.value) {
    case 'monthly': return appStore.lang === 'zh' ? '月度 VIP' : 'Monthly VIP'
    case 'quarterly': return appStore.lang === 'zh' ? '季度 VIP' : 'Quarterly VIP'
    case 'yearly': return appStore.lang === 'zh' ? '年度 VIP' : 'Yearly VIP'
    case 'lifetime': return appStore.lang === 'zh' ? '终生尊享' : 'Lifetime VIP'
    default: return ''
  }
})

const cleanProductTitle = computed(() => {
  const p = props.product
  return appStore.lang === 'zh' ? (p.title_zh || p.title_en) : (p.title_en || p.title_zh)
})

// Dynamic Payment Channels
const activeChannels = computed(() => {
  if (fetchedChannels.value && fetchedChannels.value.length > 0) {
    const currentScope = appStore.regionScope || 'china'
    const filtered = fetchedChannels.value.filter(c => {
      const isEnabled = c.enabled !== false && c.is_active !== false
      if (!isEnabled) return false
      if (c.scope === 'global' || !c.scope) return true
      if (currentScope === 'china' && c.scope === 'china_only') return true
      if (currentScope === 'overseas' && c.scope === 'overseas_only') return true
      return false
    })
    if (filtered.length > 0) return filtered
  }

  // Graceful Fallback for China & Global visitors
  if (appStore.regionScope === 'china' || appStore.lang === 'zh') {
    return [
      { code: 'alipay', name: '支付宝支付', icon: '🔵', scope: 'china_only', subtitle: '推荐使用 · 自动下发' },
      { code: 'wechat', name: '微信支付', icon: '🟢', scope: 'china_only', subtitle: '微信扫码 · 官方安全' },
      { code: 'paypal', name: 'PayPal 国际', icon: '💳', scope: 'global', subtitle: '信用卡 / Debit Card' },
      { code: 'test', name: '沙盒免密秒测', icon: '⚡', scope: 'global', subtitle: '无需扣款 · 即刻体验' }
    ]
  } else {
    return [
      { code: 'paypal', name: 'PayPal Official', icon: '💳', scope: 'overseas_only', subtitle: 'Credit / Debit Card' },
      { code: 'alipay', name: 'Alipay Global', icon: '🔵', scope: 'global', subtitle: 'Scan QR / Web' },
      { code: 'test', name: 'Sandbox Instant Test', icon: '⚡', scope: 'global', subtitle: 'No real charge · Demo' }
    ]
  }
})

const currentChannelInfo = computed(() => {
  return activeChannels.value.find(c => c.code === selectedChannel.value)
})

const payButtonText = computed(() => {
  const price = currentPlanPrice.value
  const sym = currencySymbol.value
  const planName = currentPlanName.value

  if (selectedChannel.value === 'test') {
    return appStore.lang === 'zh' ? `沙盒免密秒级开通 · ${planName} (${sym}${price})` : `Sandbox Instant Unlock · ${planName} (${sym}${price})`
  }
  if (selectedChannel.value === 'alipay' || selectedChannel.value === 'wechat' || selectedChannel.value === 'custom') {
    const ch = currentChannelInfo.value
    const chName = ch ? ch.name : (appStore.lang === 'zh' ? '直接付款' : 'Direct Pay')
    return appStore.lang === 'zh' ? `进入直接付款 · ${chName} (${sym}${price})` : `Direct Pay · ${chName} (${sym}${price})`
  }
  if (selectedChannel.value === 'paypal') {
    return appStore.lang === 'zh' ? `唤醒 PayPal 官方支付 · ${planName} ($${price})` : `Wake up PayPal Checkout · ${planName} ($${price})`
  }
  const ch = currentChannelInfo.value
  const chName = ch ? ch.name : (appStore.lang === 'zh' ? '立即支付' : 'Pay')
  return `${chName} · ${planName} (${sym}${price})`
})

const currentQrCodeUrl = computed(() => {
  return pendingOrder.value?.qr_code_url || currentChannelInfo.value?.qr_code_url || ''
})

function handlePreviewQR(url) {
  if (url) {
    showQRLightbox.value = true
  }
}

function copyToClipboard(text, msg) {
  if (!text) return
  uni.setClipboardData({
    data: text,
    success: () => {
      uni.showToast({
        title: msg || (appStore.lang === 'zh' ? '已复制到剪贴板' : 'Copied to clipboard'),
        icon: 'none'
      })
    }
  })
}

function getChannelDefaultIcon(code) {
  switch (code) {
    case 'alipay': return '🔵'
    case 'wechat': return '🟢'
    case 'paypal': return '💳'
    case 'test': return '⚡'
    default: return '💰'
  }
}

function getChannelSubtitle(ch) {
  if (ch.subtitle) return ch.subtitle
  switch (ch.code) {
    case 'alipay': return appStore.lang === 'zh' ? '支付宝直接支付 / 扫码' : 'Alipay direct checkout'
    case 'wechat': return appStore.lang === 'zh' ? '微信零钱 / 借记卡' : 'WeChat Pay'
    case 'paypal': return appStore.lang === 'zh' ? '信用卡 / 借记卡 / 账户' : 'Credit / Debit Card'
    case 'test': return appStore.lang === 'zh' ? '无需真实扣款 · 秒级下发' : 'No real charge · Instant test'
    default: return ch.scope === 'china_only' ? '国内通道' : '全球通道'
  }
}

async function loadPaymentChannels() {
  try {
    const res = await request.get('/payment/channels')
    if (res && res.channels && Array.isArray(res.channels)) {
      fetchedChannels.value = res.channels
    } else if (Array.isArray(res)) {
      fetchedChannels.value = res
    }
  } catch (err) {
    console.warn('Failed to load payment channels:', err)
  } finally {
    // Keep selection valid
    if (activeChannels.value.length > 0 && !activeChannels.value.some(c => c.code === selectedChannel.value)) {
      selectedChannel.value = activeChannels.value[0].code
    }
  }
}

watch(() => [props.visible, props.product], ([val, prod]) => {
  if (val) {
    showSecondaryConfirm.value = false
    isSubmittedReview.value = false
    pendingOrder.value = null
    if (userStore.userInfo?.email) {
      email.value = userStore.userInfo.email
    }
    loadPaymentChannels()

    // Ensure selectedPlan aligns with the tiers allowed for this wallpaper
    const allowed = availableVipTiers.value.map(t => t.type)
    if (allowed.length === 0) {
      selectedPlan.value = 'single'
    } else if (selectedPlan.value !== 'single' && !allowed.includes(selectedPlan.value)) {
      selectedPlan.value = allowed.includes('quarterly') ? 'quarterly' : allowed[0]
    }
  }
}, { immediate: true })

watch(activeChannels, (channels) => {
  if (channels.length > 0 && !channels.some(c => c.code === selectedChannel.value)) {
    selectedChannel.value = channels[0].code
  }
}, { immediate: true })

onMounted(() => {
  loadPaymentChannels()
})

function generateGuestEmail() {
  const rand = Math.floor(100000 + Math.random() * 900000)
  email.value = `guest_${rand}@wallart.vip`
  uni.showToast({ title: appStore.lang === 'zh' ? '已生成匿名凭证' : 'Guest credential generated', icon: 'none' })
}

function handleAuthSuccess(user) {
  if (user?.email) {
    email.value = user.email
  }
}

function close() {
  emit('close')
}

async function handleRedeem() {
  const code = redeemCode.value.trim()
  if (!code) {
    uni.showToast({ title: appStore.lang === 'zh' ? '请输入兑换口令' : 'Please enter code', icon: 'none' })
    return
  }
  submitting.value = true
  try {
    const res = await request.post('/redeem/apply', {
      code: code,
      product_id: props.product.id,
      user_email: email.value.trim()
    })
    uni.showToast({ title: res.message || (appStore.lang === 'zh' ? '兑换成功！' : 'Redeemed successfully!'), icon: 'success' })
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
    emit('success', res)
    close()
    uni.navigateTo({
      url: `/pages/download/index?id=${props.product.id}&title=${encodeURIComponent(props.product.title_zh || '')}&code=${encodeURIComponent(unlockCode.value || '')}`
    })
  } catch (err) {
    uni.showToast({ title: err.message || (appStore.lang === 'zh' ? '口令无效或已过期' : 'Invalid code'), icon: 'none' })
  } finally {
    submitting.value = false
  }
}

async function handlePay() {
  // If email is empty or invalid, automatically generate a guest credential so the user is NEVER blocked!
  if (!email.value || !email.value.includes('@')) {
    generateGuestEmail()
  }

  const isDirectPay = selectedChannel.value === 'wechat' || selectedChannel.value === 'alipay' || selectedChannel.value === 'custom'

  // For WeChat / Alipay direct payment channels, IMMEDIATELY switch to direct pay view!
  // This guarantees 100% instant UI response!
  if (isDirectPay) {
    pendingOrder.value = {
      order_id: pendingOrder.value?.order_id || ('ORD-' + Date.now().toString().slice(-8)),
      paypal_order_id: pendingOrder.value?.paypal_order_id || ('ORD-' + Date.now().toString().slice(-8)),
      channel_name: currentChannelInfo.value?.name || selectedChannel.value,
      account_info: currentChannelInfo.value?.account_info || '',
      qr_code_url: currentChannelInfo.value?.qr_code_url || '',
      amount: currentPlanPrice.value
    }
    showSecondaryConfirm.value = true
  }

  submitting.value = true
  const isTest = selectedChannel.value === 'test'
  const isVip = selectedPlan.value !== 'single'

  try {
    if (!isDirectPay) {
      uni.showToast({
        title: isTest ? (appStore.lang === 'zh' ? '正在沙盒免密模拟...' : 'Simulating in sandbox...') : (appStore.lang === 'zh' ? '正在连接网关...' : 'Connecting gateway...'),
        icon: 'loading',
        duration: 15000
      })
    }

    const curCurrency = appStore.currency || (selectedChannel.value === 'alipay' || selectedChannel.value === 'wechat' ? 'CNY' : 'USD')
    const curAmount = Number(currentPlanPrice.value) || 0

    const res = await paypalAPI.createOrder({
      product_id: props.product?.id || '',
      order_type: isVip ? 'membership' : 'single',
      membership_type: isVip ? selectedPlan.value : '',
      user_email: email.value.trim(),
      payment_method: selectedChannel.value,
      is_test: isTest,
      amount: curAmount,
      currency: curCurrency
    })

    if (!isDirectPay) {
      uni.hideToast()
    }

    // 1. Sandbox instant unlock flow
    if (isTest && (res.status === 'paid' || res.is_test)) {
      if (res.token) {
        uni.setStorageSync('token', res.token)
      }
      if (res.user) {
        userStore.userInfo = {
          ...(userStore.userInfo || {}),
          ...res.user,
          is_member: Boolean(isVip || res.user.membership_type !== 'none')
        }
        uni.setStorageSync('user_info', JSON.stringify(userStore.userInfo))
      }

      uni.showToast({ title: appStore.lang === 'zh' ? '解锁成功！' : 'Unlocked successfully!', icon: 'success' })
      emit('success', { order_type: selectedPlan.value })
      close()

      if (res.order_id) uni.setStorageSync('last_paid_order_id', res.order_id)
      if (email.value) uni.setStorageSync('last_buyer_email', email.value)

      uni.navigateTo({
        url: `/pages/download/index?id=${props.product.id}&title=${encodeURIComponent(props.product.title_zh || '')}&order_id=${encodeURIComponent(res.order_id || '')}&email=${encodeURIComponent(email.value || '')}`
      })
      return
    }

    // 2. PayPal Gateway with approval link
    if (selectedChannel.value === 'paypal') {
      if (res.approval_url) {
        uni.setStorageSync('pending_paypal_order_id', res.order_id)
        
        // #ifdef H5
        if (typeof window !== 'undefined') {
          window.location.href = res.approval_url
          return
        }
        // #endif

        uni.setClipboardData({
          data: res.approval_url,
          success: () => {
            uni.showModal({
              title: appStore.lang === 'zh' ? 'PayPal 支付链接已生成' : 'PayPal Link Ready',
              content: appStore.lang === 'zh' ? '已将 PayPal 支付收银链接复制到剪贴板，请在浏览器中打开完成付款。' : 'PayPal link copied to clipboard.',
              showCancel: false
            })
          }
        })
        return
      } else {
        throw new Error(res.message || (appStore.lang === 'zh' ? '未获取到 PayPal 支付链接，请检查系统后台 PayPal 参数配置' : 'No PayPal approval URL returned'))
      }
    }

    // 3. Alipay, WeChat, or custom channels:
    // Update pendingOrder with the persisted backend order ID, channel info, and QR code URL
    pendingOrder.value = {
      order_id: res.order_id || pendingOrder.value?.order_id,
      paypal_order_id: res.paypal_order_id || res.order_id || pendingOrder.value?.order_id,
      channel_name: currentChannelInfo.value?.name || res.channel_name || selectedChannel.value,
      account_info: res.account_info || currentChannelInfo.value?.account_info || '',
      qr_code_url: res.qr_code_url || currentChannelInfo.value?.qr_code_url || '',
      amount: res.amount || currentPlanPrice.value
    }
    showSecondaryConfirm.value = true

  } catch (err) {
    if (!isDirectPay) {
      uni.hideToast()
      uni.showModal({
        title: appStore.lang === 'zh' ? '支付发起提示' : 'Payment Alert',
        content: err.message || (appStore.lang === 'zh' ? '连接支付通道超时，请选择沙盒免密秒级测试快速体验。' : 'Payment timed out. Please try Sandbox test mode.'),
        showCancel: false
      })
    } else {
      console.warn('Direct pay background order creation notice:', err)
    }
  } finally {
    submitting.value = false
  }
}

async function handleConfirmPaid() {
  const orderId = pendingOrder.value?.order_id || pendingOrder.value?.paypal_order_id
  if (!orderId) {
    uni.showToast({ title: appStore.lang === 'zh' ? '订单号异常，请重新发起' : 'Order ID error', icon: 'none' })
    return
  }

  confirmingPayment.value = true
  try {
    uni.showToast({
      title: appStore.lang === 'zh' ? '正在核验支付...' : 'Verifying payment...',
      icon: 'loading',
      duration: 10000
    })

    const res = await paypalAPI.captureOrder(orderId)
    uni.hideToast()

    if (res.status === 'paid' || res.is_test || res.order?.status === 'paid') {
      if (res.token) {
        uni.setStorageSync('token', res.token)
      }
      if (res.user) {
        userStore.userInfo = {
          ...(userStore.userInfo || {}),
          ...res.user,
          is_member: Boolean(selectedPlan.value !== 'single' || res.user.membership_type !== 'none')
        }
        uni.setStorageSync('user_info', JSON.stringify(userStore.userInfo))
      }

      uni.showToast({
        title: appStore.lang === 'zh' ? '🎉 支付核验成功，壁纸已解锁！' : 'Payment confirmed & unlocked!',
        icon: 'success'
      })
      emit('success', { order_type: selectedPlan.value, order_id: orderId })
      showSecondaryConfirm.value = false
      close()

      if (orderId) uni.setStorageSync('last_paid_order_id', orderId)
      if (email.value) uni.setStorageSync('last_buyer_email', email.value)

      uni.navigateTo({
        url: `/pages/download/index?id=${props.product.id}&title=${encodeURIComponent(props.product.title_zh || '')}&order_id=${encodeURIComponent(orderId || '')}&email=${encodeURIComponent(email.value || '')}`
      })
    } else if (res.status === 'pending_review' || res.pending_review || res.order?.status === 'pending_review') {
      isSubmittedReview.value = true
      uni.showToast({
        title: appStore.lang === 'zh' ? '凭据已提交审核！' : 'Submitted for review',
        icon: 'success'
      })
    } else {
      uni.showModal({
        title: appStore.lang === 'zh' ? '核验提示' : 'Verification Alert',
        content: res.message || (appStore.lang === 'zh' ? '尚未查询到到账凭据，请确认已完成转账并在备注中填写邮箱。' : 'Payment not detected yet.'),
        showCancel: false
      })
    }
  } catch (err) {
    uni.hideToast()
    uni.showModal({
      title: appStore.lang === 'zh' ? '核验提示' : 'Verification Alert',
      content: err.message || (appStore.lang === 'zh' ? '核验超时或订单未支付，请稍后再试或联系客服。' : 'Verification failed'),
      showCancel: false
    })
  } finally {
    confirmingPayment.value = false
  }
}

function goToLookup() {
  close()
  uni.navigateTo({
    url: `/pages/lookup/index?email=${encodeURIComponent(email.value || '')}`
  })
}
</script>

<style scoped>
.cashier-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.76);
  backdrop-filter: blur(14px);
  -webkit-backdrop-filter: blur(14px);
  display: flex;
  align-items: flex-end;
  justify-content: center;
  z-index: 100000;
}

.cashier-sheet {
  width: 100%;
  max-width: 520px;
  max-height: 92vh;
  overflow-y: auto;
  border-radius: 28px 28px 0 0;
  padding: 16px 20px calc(24px + env(safe-area-inset-bottom));
  box-sizing: border-box;
  animation: sheetSlideUp 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes sheetSlideUp {
  from { transform: translateY(100%); }
  to { transform: translateY(0); }
}

.cashier-sheet.dark {
  background: #141622;
  border-top: 1px solid rgba(255, 255, 255, 0.12);
  color: #f1f5f9;
  box-shadow: 0 -12px 48px rgba(0, 0, 0, 0.7);
}

.cashier-sheet.light {
  background: #ffffff;
  border-top: 1px solid rgba(0, 0, 0, 0.08);
  color: #0f172a;
  box-shadow: 0 -8px 36px rgba(0, 0, 0, 0.15);
}

.drag-bar {
  width: 44px;
  height: 4px;
  background: rgba(148, 163, 184, 0.35);
  border-radius: 999px;
  margin: 0 auto 12px;
}

.sheet-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}

.title-row {
  display: flex;
  flex-direction: column;
}

.sheet-title {
  font-size: 18px;
  font-weight: 800;
  letter-spacing: 0.3px;
}

.sheet-sub {
  font-size: 11px;
  opacity: 0.6;
  margin-top: 2px;
}

.close-btn {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: rgba(148, 163, 184, 0.14);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.close-icon {
  font-size: 13px;
  color: #94a3b8;
}

/* Guest Registration Incentive Banner */
.guest-register-incentive {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.14), rgba(217, 119, 6, 0.08));
  border: 1px solid rgba(245, 158, 11, 0.35);
  border-radius: 14px;
  padding: 10px 12px;
  margin-bottom: 12px;
  gap: 10px;
}

.incentive-left {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  flex: 1;
}

.incentive-crown {
  font-size: 16px;
  flex-shrink: 0;
  margin-top: 2px;
}

.incentive-text-wrap {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.incentive-title {
  font-size: 12px;
  font-weight: 700;
  color: #fbbf24;
}

.incentive-sub {
  font-size: 10px;
  line-height: 1.35;
  opacity: 0.8;
}

.btn-incentive-reg {
  height: 28px;
  padding: 0 10px;
  background: #f59e0b;
  color: #111827;
  border-radius: 999px;
  border: none;
  font-size: 11px;
  font-weight: 800;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(245, 158, 11, 0.35);
}

.btn-reg-text {
  font-size: 11px;
  font-weight: 800;
  color: #111827;
  white-space: nowrap;
}

/* Product Mini Strip */
.product-strip {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  background: rgba(148, 163, 184, 0.08);
  border-radius: 12px;
  margin-bottom: 14px;
}

.strip-thumb {
  width: 44px;
  height: 44px;
  border-radius: 8px;
  flex-shrink: 0;
}

.strip-placeholder {
  background: linear-gradient(135deg, #1e293b, #334155);
}

.strip-info {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.strip-title {
  font-size: 13px;
  font-weight: 700;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.strip-tag {
  font-size: 10px;
  opacity: 0.55;
  margin-top: 2px;
}

/* Plan Cards */
.plan-options {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 12px;
}

.plan-card {
  display: flex;
  align-items: center;
  padding: 10px 12px;
  border-radius: 14px;
  border: 1.5px solid rgba(148, 163, 184, 0.2);
  background: rgba(148, 163, 184, 0.06);
  cursor: pointer;
  position: relative;
  transition: all 0.2s ease;
}

.plan-card.selected {
  border-color: #3b82f6;
  background: rgba(59, 130, 246, 0.1);
}

.plan-title-line {
  display: flex;
  align-items: center;
  gap: 6px;
}

.plan-badge-mini {
  font-size: 9px;
  padding: 1px 5px;
  border-radius: 999px;
  background: rgba(59, 130, 246, 0.15);
  color: #38bdf8;
  font-weight: 700;
}

.plan-radio {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  border: 2px solid rgba(148, 163, 184, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 10px;
  flex-shrink: 0;
}

.plan-card.selected .plan-radio {
  border-color: #3b82f6;
}

.radio-inner {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #3b82f6;
}

.plan-body {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.plan-name {
  font-size: 13px;
  font-weight: 700;
}

.plan-desc {
  font-size: 10px;
  opacity: 0.55;
  margin-top: 2px;
}

.plan-price-wrap {
  display: flex;
  align-items: baseline;
  gap: 2px;
}

.plan-symbol {
  font-size: 12px;
  font-weight: 700;
  color: #3b82f6;
}

.plan-price {
  font-size: 17px;
  font-weight: 800;
}

/* VIP Section Divider */
.vip-section-divider {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 4px 0 8px;
}

.divider-line {
  flex: 1;
  height: 1px;
  background: rgba(148, 163, 184, 0.18);
}

.divider-label {
  display: flex;
  align-items: center;
  gap: 5px;
}

.divider-crown {
  font-size: 13px;
}

.divider-text {
  font-size: 11px;
  font-weight: 700;
  color: #f59e0b;
}

/* VIP 4-Tier Grid */
.vip-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
  margin-bottom: 12px;
}

.vip-tier-card {
  position: relative;
  padding: 10px 11px;
  border-radius: 14px;
  border: 1.5px solid rgba(148, 163, 184, 0.2);
  background: rgba(148, 163, 184, 0.05);
  cursor: pointer;
  display: flex;
  flex-direction: column;
  transition: all 0.2s ease;
  overflow: hidden;
}

.vip-tier-card.selected {
  border-color: #f59e0b;
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.14), rgba(217, 119, 6, 0.06));
  box-shadow: 0 4px 16px rgba(245, 158, 11, 0.16);
}

.vip-tier-card.featured {
  border-color: rgba(245, 158, 11, 0.35);
}

.featured-ribbon {
  position: absolute;
  top: 0;
  right: 0;
  background: linear-gradient(135deg, #f59e0b, #d97706);
  padding: 1px 6px;
  border-bottom-left-radius: 8px;
}

.featured-ribbon .ribbon-text {
  font-size: 9px;
  font-weight: 800;
  color: #111827;
}

.vip-tier-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 4px;
}

.vip-tier-name {
  font-size: 12px;
  font-weight: 700;
}

.vip-tier-card.selected .vip-tier-name {
  color: #fbbf24;
}

.tier-pill-small {
  font-size: 9px;
  font-weight: 700;
  padding: 1px 5px;
  border-radius: 999px;
  background: rgba(148, 163, 184, 0.15);
  color: #94a3b8;
}

.tier-pill-small.gold {
  background: rgba(245, 158, 11, 0.2);
  color: #f59e0b;
}

.tier-pill-small.diamond {
  background: rgba(168, 85, 247, 0.2);
  color: #c084fc;
}

.vip-tier-price-row {
  display: flex;
  align-items: baseline;
  gap: 2px;
  margin-bottom: 2px;
}

.vip-cur {
  font-size: 11px;
  font-weight: 700;
  color: #f59e0b;
}

.vip-price-num {
  font-size: 17px;
  font-weight: 800;
  color: #f59e0b;
}

.vip-tier-hint {
  font-size: 9px;
  opacity: 0.6;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.gold-hint {
  color: #fbbf24;
  opacity: 0.9;
}

.vip-grid.grid-three {
  grid-template-columns: repeat(3, 1fr);
  gap: 6px;
}

.vip-grid.grid-single {
  grid-template-columns: 1fr;
}

.vip-grid.grid-three .vip-tier-card {
  padding: 8px 6px;
}

.vip-grid.grid-three .vip-tier-name {
  font-size: 11px;
}

.vip-grid.grid-three .tier-pill-small {
  font-size: 8px;
  padding: 1px 3px;
}

.vip-grid.grid-three .vip-price-num {
  font-size: 15px;
}

.vip-grid.grid-three .vip-tier-hint {
  font-size: 8.5px;
}

.vip-exclusive-box {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  margin-bottom: 12px;
  border-radius: 14px;
  background: rgba(148, 163, 184, 0.08);
  border: 1px dashed rgba(148, 163, 184, 0.25);
}

.vip-exclusive-icon {
  font-size: 18px;
  flex-shrink: 0;
}

.vip-exclusive-text {
  font-size: 12px;
  color: #94a3b8;
  line-height: 1.5;
}

/* Delivery Box */
.delivery-box {
  margin-bottom: 10px;
}

.delivery-label-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.delivery-label {
  font-size: 11px;
  font-weight: 600;
  opacity: 0.7;
}

.quick-email-btn {
  font-size: 11px;
  color: #38bdf8;
  cursor: pointer;
}

.delivery-input {
  width: 100%;
  height: 40px;
  box-sizing: border-box;
  border-radius: 12px;
  padding: 0 12px;
  font-size: 13px;
  border: 1px solid rgba(148, 163, 184, 0.25);
  background: rgba(148, 163, 184, 0.08);
}

/* Redeem */
.redeem-row {
  margin-bottom: 8px;
}

.redeem-toggle {
  font-size: 11px;
  color: #38bdf8;
  cursor: pointer;
}

.redeem-box {
  display: flex;
  gap: 8px;
  margin-bottom: 10px;
}

.redeem-btn {
  height: 40px;
  padding: 0 16px;
  border-radius: 12px;
  background: #3b82f6;
  color: #ffffff;
  border: none;
  font-size: 12px;
  font-weight: 700;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

/* Payment Methods */
.cashier-pay-methods {
  margin-bottom: 14px;
}

.channel-header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.c-section-label {
  font-size: 12px;
  font-weight: 700;
  opacity: 0.8;
}

.geo-region-tag {
  font-size: 10px;
  padding: 1px 8px;
  border-radius: 999px;
  background: rgba(59, 130, 246, 0.14);
  color: #38bdf8;
  font-weight: 600;
}

.c-methods-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
}

@media (max-width: 400px) {
  .c-methods-grid {
    grid-template-columns: 1fr;
  }
}

.c-method-chip {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 10px;
  border-radius: 12px;
  border: 1.5px solid rgba(148, 163, 184, 0.2);
  background: rgba(148, 163, 184, 0.06);
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
}

.c-method-chip.active {
  border-color: #3b82f6;
  background: rgba(59, 130, 246, 0.12);
}

.c-method-chip.test-chip {
  border-style: dashed;
}

.c-method-icon {
  font-size: 16px;
  flex-shrink: 0;
}

.c-method-text-wrap {
  display: flex;
  flex-direction: column;
  flex: 1;
  overflow: hidden;
}

.c-method-title {
  font-size: 11px;
  font-weight: 700;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.c-method-sub {
  font-size: 9px;
  opacity: 0.55;
  margin-top: 1px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.c-method-radio {
  width: 14px;
  height: 14px;
  border-radius: 50%;
  border: 1.5px solid rgba(148, 163, 184, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.c-method-chip.active .c-method-radio {
  border-color: #3b82f6;
}

.c-radio-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #3b82f6;
}

/* Channel info card */
.channel-info-card {
  padding: 8px 10px;
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.25);
  border-radius: 12px;
  margin-bottom: 12px;
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.channel-info-title {
  font-size: 11px;
  font-weight: 700;
  color: #38bdf8;
}

.channel-info-desc {
  font-size: 10px;
  line-height: 1.4;
  opacity: 0.85;
  white-space: pre-wrap;
}

/* Action button */
.sheet-action {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.btn-pay {
  height: 46px;
  border-radius: 14px;
  background: #2563eb;
  color: #ffffff;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  box-shadow: 0 6px 20px rgba(37, 99, 235, 0.35);
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-pay.btn-pay-vip {
  background: linear-gradient(135deg, #d4af37, #f59e0b);
  box-shadow: 0 6px 20px rgba(245, 158, 11, 0.35);
}

.btn-pay-icon {
  font-size: 16px;
}

.btn-pay-text {
  font-size: 14px;
  font-weight: 800;
  color: #ffffff;
}

.manual-confirm-entrance {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 6px 0;
  cursor: pointer;
}

.manual-confirm-text {
  font-size: 11px;
  color: #38bdf8;
  font-weight: 600;
  text-decoration: underline;
  opacity: 0.9;
}

/* Direct Pay Hint Card in Primary View */
.direct-pay-hint-card {
  background: rgba(59, 130, 246, 0.08);
  border: 1px solid rgba(59, 130, 246, 0.25);
  border-radius: 14px;
  padding: 10px 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.hint-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.hint-badge-left {
  display: flex;
  align-items: center;
  gap: 6px;
}

.hint-badge-icon {
  font-size: 14px;
}

.hint-badge-title {
  font-size: 12px;
  font-weight: 700;
  color: #38bdf8;
}

.hint-badge-tag {
  font-size: 10px;
  padding: 2px 8px;
  border-radius: 999px;
  background: rgba(16, 185, 129, 0.15);
  border: 1px solid rgba(16, 185, 129, 0.3);
  color: #10b981;
  font-weight: 600;
}

.hint-card-content {
  display: flex;
  align-items: center;
  gap: 10px;
}

.hint-qr-thumb {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  background: #ffffff;
  padding: 2px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  cursor: pointer;
  flex-shrink: 0;
}

.hint-desc-wrap {
  display: flex;
  flex-direction: column;
  gap: 2px;
  flex: 1;
}

.hint-desc-main {
  font-size: 12px;
  font-weight: 600;
}

.hint-desc-sub {
  font-size: 10px;
  opacity: 0.75;
  line-height: 1.4;
}

.btn-pay.btn-pay-direct {
  background: linear-gradient(135deg, #0284c7, #2563eb);
  box-shadow: 0 6px 20px rgba(2, 132, 199, 0.35);
}

/* Secondary Payment Confirmation & Direct Pay View */
.secondary-confirm-view {
  display: flex;
  flex-direction: column;
  gap: 12px;
  animation: fadeIn 0.25s ease;
  padding-bottom: 8px;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(6px); }
  to { opacity: 1; transform: translateY(0); }
}

.confirm-card {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.confirm-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 8px;
  border-bottom: 1px dashed rgba(148, 163, 184, 0.25);
}

.confirm-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: rgba(59, 130, 246, 0.15);
  border: 1px solid rgba(59, 130, 246, 0.35);
  padding: 4px 12px;
  border-radius: 999px;
}

.confirm-badge-icon {
  font-size: 14px;
}

.confirm-badge-text {
  font-size: 12px;
  font-weight: 700;
  color: #38bdf8;
}

.order-no-wrap {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  background: rgba(148, 163, 184, 0.1);
  padding: 3px 8px;
  border-radius: 6px;
}

.order-no-text {
  font-size: 11px;
  font-family: monospace;
  opacity: 0.75;
}

.copy-tag {
  font-size: 10px;
  color: #38bdf8;
  font-weight: 600;
}

/* Hero Box for Direct Payment */
.direct-pay-hero-box {
  display: flex;
  flex-direction: column;
  gap: 10px;
  background: rgba(15, 23, 42, 0.4);
  border: 1px solid rgba(59, 130, 246, 0.25);
  border-radius: 18px;
  padding: 14px;
}

.cashier-sheet.light .direct-pay-hero-box {
  background: rgba(241, 245, 249, 0.8);
  border-color: rgba(59, 130, 246, 0.2);
}

.hero-amount-tag {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  border-radius: 12px;
  background: rgba(245, 158, 11, 0.1);
  border: 1px solid rgba(245, 158, 11, 0.25);
}

.hero-amount-label {
  font-size: 12px;
  font-weight: 600;
  color: #fbbf24;
}

.hero-amount-val {
  display: flex;
  align-items: baseline;
  gap: 3px;
}

.hero-cur {
  font-size: 14px;
  font-weight: 700;
  color: #f59e0b;
}

.hero-num {
  font-size: 22px;
  font-weight: 900;
  color: #f59e0b;
}

.hero-sub {
  font-size: 11px;
  opacity: 0.65;
  margin-left: 4px;
}

/* QR Code Card */
.qr-container-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 10px 0;
}

.qr-img-frame {
  width: 190px;
  height: 190px;
  background: #ffffff;
  border-radius: 16px;
  padding: 10px;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.25), 0 0 0 1px rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: transform 0.2s ease;
}

.qr-img-frame:active {
  transform: scale(0.98);
}

.direct-qr-image {
  width: 100%;
  height: 100%;
}

.qr-scan-tips {
  display: flex;
  align-items: center;
  gap: 6px;
}

.qr-scan-icon {
  font-size: 12px;
}

.qr-scan-txt {
  font-size: 11px;
  color: #38bdf8;
  font-weight: 600;
  text-align: center;
}

/* Fallback Card */
.qr-fallback-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 18px 12px;
  text-align: center;
  border-radius: 12px;
  background: rgba(148, 163, 184, 0.08);
  border: 1px dashed rgba(148, 163, 184, 0.3);
}

.fallback-icon {
  font-size: 28px;
}

.fallback-title {
  font-size: 14px;
  font-weight: 700;
}

.fallback-desc {
  font-size: 12px;
  opacity: 0.8;
  line-height: 1.5;
  word-break: break-all;
}

/* Transfer Memo Card */
.transfer-memo-card {
  display: flex;
  flex-direction: column;
  gap: 6px;
  background: rgba(16, 185, 129, 0.08);
  border: 1px solid rgba(16, 185, 129, 0.25);
  border-radius: 12px;
  padding: 8px 12px;
}

.memo-title-line {
  display: flex;
  align-items: center;
  gap: 6px;
}

.memo-icon {
  font-size: 13px;
}

.memo-title {
  font-size: 11px;
  font-weight: 700;
  color: #10b981;
}

.memo-content-line {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.memo-desc {
  font-size: 11px;
  opacity: 0.75;
}

.memo-email {
  font-size: 12px;
  font-family: monospace;
  font-weight: 700;
  color: #10b981;
  background: rgba(16, 185, 129, 0.15);
  padding: 2px 6px;
  border-radius: 4px;
  cursor: pointer;
}

.memo-copy-btn {
  font-size: 10px;
  color: #38bdf8;
  font-weight: 600;
  cursor: pointer;
  background: rgba(56, 189, 248, 0.1);
  padding: 2px 6px;
  border-radius: 4px;
}

.confirm-summary {
  display: flex;
  flex-direction: column;
  gap: 8px;
  background: rgba(148, 163, 184, 0.08);
  border: 1px solid rgba(148, 163, 184, 0.15);
  border-radius: 14px;
  padding: 10px 14px;
}

.summary-line {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
}

.s-label {
  opacity: 0.65;
}

.s-val {
  font-weight: 600;
  max-width: 220px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.channel-val {
  color: #3b82f6;
  font-weight: 700;
}

.text-email {
  font-family: monospace;
  color: #10b981;
}

.confirm-actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-top: 6px;
}

.btn-confirm-pay {
  height: 48px;
  border-radius: 14px;
  background: linear-gradient(135deg, #10b981, #059669);
  color: #ffffff;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  box-shadow: 0 6px 20px rgba(16, 185, 129, 0.35);
  cursor: pointer;
  font-size: 15px;
  font-weight: 800;
  transition: all 0.2s ease;
}

.btn-confirm-icon {
  font-size: 18px;
  font-weight: 800;
}

.btn-confirm-text {
  font-size: 15px;
  font-weight: 800;
  color: #ffffff;
}

.btn-cancel-confirm {
  height: 40px;
  border-radius: 12px;
  background: transparent;
  color: inherit;
  border: 1px solid rgba(148, 163, 184, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.btn-cancel-text {
  font-size: 12px;
  opacity: 0.7;
}

.security-tip {
  font-size: 10px;
  text-align: center;
  opacity: 0.5;
}

/* Dedicated QR Lightbox Magnifier Dialog (High Z-Index above Cashier Modal) */
.qr-lightbox-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.88);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100001;
  padding: 20px;
  box-sizing: border-box;
  animation: fadeIn 0.2s ease;
}

.qr-lightbox-dialog {
  width: 100%;
  max-width: 360px;
  background: #ffffff;
  color: #0f172a;
  border-radius: 24px;
  padding: 24px 20px 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  position: relative;
  box-shadow: 0 24px 64px rgba(0, 0, 0, 0.6), 0 0 0 1px rgba(255, 255, 255, 0.15);
  animation: zoomIn 0.25s cubic-bezier(0.16, 1, 0.3, 1);
  box-sizing: border-box;
}

@keyframes zoomIn {
  from { opacity: 0; transform: scale(0.92); }
  to { opacity: 1; transform: scale(1); }
}

.lightbox-close-btn {
  position: absolute;
  top: 14px;
  right: 14px;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: rgba(15, 23, 42, 0.08);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 10;
}

.lightbox-close-icon {
  font-size: 15px;
  color: #64748b;
  font-weight: 700;
}

.lightbox-channel-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.lightbox-channel-icon {
  font-size: 20px;
}

.lightbox-channel-name {
  font-size: 16px;
  font-weight: 800;
  color: #0f172a;
}

.lightbox-badge-direct {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 999px;
  background: rgba(16, 185, 129, 0.15);
  color: #059669;
  font-weight: 700;
}

.lightbox-img-wrap {
  width: 240px;
  height: 240px;
  background: #ffffff;
  border-radius: 16px;
  padding: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(0, 0, 0, 0.08);
  box-sizing: border-box;
}

.lightbox-img {
  width: 100%;
  height: 100%;
}

.lightbox-img-ph {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.lightbox-amount-pill {
  display: flex;
  align-items: baseline;
  gap: 4px;
  background: rgba(245, 158, 11, 0.12);
  border: 1px solid rgba(245, 158, 11, 0.3);
  padding: 6px 14px;
  border-radius: 999px;
}

.lightbox-amount-label {
  font-size: 12px;
  color: #d97706;
  font-weight: 600;
}

.lightbox-amount-val {
  font-size: 20px;
  font-weight: 900;
  color: #d97706;
}

.lightbox-amount-sub {
  font-size: 11px;
  color: #b45309;
}

.lightbox-hint {
  font-size: 11px;
  color: #64748b;
  text-align: center;
  line-height: 1.4;
}

.lightbox-btn-done {
  width: 100%;
  height: 44px;
  border-radius: 12px;
  background: linear-gradient(135deg, #10b981, #059669);
  color: #ffffff;
  font-size: 14px;
  font-weight: 800;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 4px;
  box-shadow: 0 4px 14px rgba(16, 185, 129, 0.35);
}

/* Force any uni-preview-image to stay on top */
:deep(uni-preview-image),
uni-preview-image {
  z-index: 999999 !important;
}
/* Submitted Review View Styling */
.submitted-review-view {
  padding: 10px 0 20px;
  animation: fadeIn 0.3s ease;
}

.review-status-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 12px;
}

.review-icon-wrap {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: rgba(245, 158, 11, 0.12);
  border: 2px solid rgba(245, 158, 11, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
}

.review-icon {
  font-size: 32px;
}

.review-title {
  font-size: 18px;
  font-weight: 800;
  color: #d97706;
}

.review-desc {
  font-size: 13px;
  color: #64748b;
  line-height: 1.5;
  padding: 0 10px;
}

.review-credential-box {
  width: 100%;
  background: rgba(0, 0, 0, 0.03);
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: 16px;
  padding: 14px 16px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  text-align: left;
  box-sizing: border-box;
  margin-top: 6px;
}

.review-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 13px;
}

.r-label {
  color: #64748b;
}

.r-val-wrap {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
}

.r-copy {
  font-size: 11px;
  color: #0071e3;
  background: rgba(0, 113, 227, 0.08);
  padding: 2px 6px;
  border-radius: 4px;
}

.text-amber {
  color: #d97706;
}

.text-emerald {
  color: #059669;
}

.review-actions-col {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-top: 10px;
}

.btn-lookup-progress {
  width: 100%;
  height: 46px;
  border-radius: 12px;
  background: linear-gradient(135deg, #0071e3, #005bb5);
  color: #ffffff;
  font-weight: 700;
  font-size: 14px;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  box-shadow: 0 4px 14px rgba(0, 113, 227, 0.3);
}

.btn-copy-kefu {
  width: 100%;
  height: 42px;
  border-radius: 12px;
  background: rgba(0, 113, 227, 0.08);
  color: #0071e3;
  font-weight: 700;
  font-size: 13px;
  border: 1px solid rgba(0, 113, 227, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}

.btn-close-review {
  width: 100%;
  height: 40px;
  border-radius: 12px;
  background: transparent;
  color: #64748b;
  font-weight: 600;
  font-size: 13px;
  border: 1px solid rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
