<template>
  <view class="detail-container" :class="[appStore.theme]">
    <!-- Dynamic Ambient Background (Starry Sky & Meteors in Dark, Glass Raindrops in Light) -->
    <AmbientCanvas />

    <view class="detail-content-wrap">
      <!-- Navbar / Top Header Bar -->
      <view class="nav-bar">
        <view class="back-circle" @click="goBack" @tap="goBack">
          <text class="back-icon">←</text>
          <text class="back-label">{{ appStore.lang === 'zh' ? '返回画廊' : 'Back' }}</text>
        </view>
        <text class="nav-title">{{ appStore.lang === 'zh' ? (product?.title_zh || '壁纸预览') : (product?.title_en || product?.title_zh || 'Wallpaper Preview') }}</text>
        
        <!-- Controls in nav -->
        <view class="nav-right-controls">
          <view class="vip-capsule-btn" @click="goToVip" @tap="goToVip">
            <text class="vip-capsule-icon">👑</text>
            <text class="vip-capsule-text">VIP特权</text>
          </view>
          <view class="icon-circle my-nav-circle" @click="goToMy" @tap="goToMy" title="个人中心">
            <text class="nav-btn-icon">👤</text>
          </view>
          <view class="icon-circle" @tap="appStore.toggleLang">
            <text class="nav-btn-text">{{ appStore.lang === 'zh' ? 'EN' : '中文' }}</text>
          </view>
          <view class="icon-circle" @tap="appStore.toggleTheme">
            <text class="nav-btn-icon">{{ appStore.theme === 'dark' ? '☀️' : '🌙' }}</text>
          </view>
        </view>
      </view>

      <!-- Main Responsive Layout: 2-Cols on PC, Stacked on Mobile -->
      <view class="detail-main-layout">
        <!-- Left Column: Interactive Mockup Switcher & Live Previews -->
        <view class="mockup-column">
          <view class="mockup-wrap">
            <MockupPreview :image-url="product?.cover_image" :product="product || {}" />
          </view>

          <!-- Live Video Motion Preview Strip (TikTok, Douyin, XHS, etc.) -->
          <view v-if="videoLinks.length > 0" class="video-preview-strip">
            <view class="v-strip-header">
              <text class="v-strip-icon">🎬</text>
              <text class="v-strip-title">{{ appStore.lang === 'zh' ? '实机动态效果演示' : 'Live Motion Previews' }}</text>
            </view>
            <view class="v-strip-tags">
              <view
                v-for="(v, idx) in videoLinks"
                :key="idx"
                class="v-link-chip"
                @tap="openVideoLink(v.url)"
              >
                <text class="v-chip-icon">▶</text>
                <text class="v-chip-name">{{ v.title || getPlatformName(v.platform) }}</text>
                <text class="v-chip-jump">↗</text>
              </view>
            </view>
          </view>
        </view>

        <!-- Right Column: Info, Specs, Rich-Text Description, Action CTA -->
        <view class="info-column">
          <view class="info-block">
            <!-- Badges row -->
            <view class="tag-row">
              <text v-if="isInteractiveProduct" class="chip chip-interactive">✨ iOS 原生 .tendies 3D交互与景深</text>
              <text class="chip chip-cat">{{ product?.category || (appStore.lang === 'zh' ? '4K 精选' : 'Curated 4K') }}</text>
              <text class="chip chip-res">4K UHD 超清原图</text>
              <text class="chip chip-cloud">夸克/百度/谷歌网盘</text>
            </view>

            <!-- Product Title -->
            <text class="product-title">{{ appStore.lang === 'zh' ? (product?.title_zh || '4K 精选壁纸') : (product?.title_en || product?.title_zh || 'Curated 4K Wallpaper') }}</text>

            <!-- Main Action Card (Always Visible on All Devices: Mobile & Desktop) -->
            <view class="main-cta-box">
              <view class="cta-price-info">
                <template v-if="userStore.isMember">
                  <template v-if="isVIPTierAllowed">
                    <text class="cta-vip-tag">👑 {{ appStore.lang === 'zh' ? 'VIP 会员特权' : 'VIP Privilege' }}</text>
                    <text class="cta-price-free">{{ appStore.lang === 'zh' ? '全场免单畅享' : 'Free Download' }}</text>
                  </template>
                  <template v-else>
                    <text class="cta-vip-tag vip-tag-restricted">👑 {{ appStore.lang === 'zh' ? '专享会员壁纸' : 'VIP Exclusive' }}</text>
                    <text class="cta-price-sub tier-alert">{{ appStore.lang === 'zh' ? `限【${allowedVIPNames}】免单` : 'Tier Upgrade Required' }}</text>
                  </template>
                </template>
                <template v-else-if="Number(product?.price) === 0">
                  <text class="cta-free-tag">🎁 {{ appStore.lang === 'zh' ? '限时福利' : 'Free Gift' }}</text>
                  <text class="cta-price-free">{{ appStore.lang === 'zh' ? '免费下载' : 'Free' }}</text>
                </template>
                <template v-else>
                  <view class="cta-price-header">
                    <text class="cta-price-sub">{{ appStore.lang === 'zh' ? '单款永久解锁' : 'Single Unlock' }}</text>
                    <text class="cta-free-badge">{{ freeTierLabel }}</text>
                  </view>
                  <view class="cta-price-main">
                    <text class="cta-cur">{{ appStore.currency === 'USD' ? '$' : '¥' }}</text>
                    <text class="cta-num">{{ displayUnitPrice }}</text>
                  </view>
                </template>
              </view>

              <button
                class="cta-action-btn"
                :class="{
                  'btn-vip-gold': (userStore.isMember && isVIPTierAllowed) || Number(product?.price) === 0,
                  'btn-vip-upgrade': userStore.isMember && !isVIPTierAllowed
                }"
                @click="handleMainCTA"
                @tap="handleMainCTA"
              >
                <text class="cta-btn-icon">{{ ctaIcon }}</text>
                <text class="cta-btn-text">{{ ctaLabel }}</text>
              </button>
            </view>

            <!-- Specifications Grid -->
            <view class="specs-grid">
              <view class="spec-cell">
                <text class="spec-label">{{ appStore.lang === 'zh' ? '画质规格' : 'Resolution' }}</text>
                <text class="spec-val">3840 × 2160 (4K)</text>
              </view>
              <view class="spec-cell">
                <text class="spec-label">{{ appStore.lang === 'zh' ? '文件格式' : 'Format' }}</text>
                <text class="spec-val">无损原画 PNG / MP4</text>
              </view>
              <view class="spec-cell">
                <text class="spec-label">{{ appStore.lang === 'zh' ? '交付通道' : 'Delivery' }}</text>
                <text class="spec-val">高速网盘直连</text>
              </view>
              <view class="spec-cell" v-if="isInteractiveProduct">
                <text class="spec-label">{{ appStore.lang === 'zh' ? '实机交互' : 'Live Interaction' }}</text>
                <text class="spec-val highlight-purple">3D景深时钟穿透 / .tendies</text>
              </view>
            </view>

            <!-- Rich Text Description Section (Formatted HTML, NOT raw code) -->
            <view class="rich-desc-container">
              <view class="desc-heading-row">
                <text class="desc-heading-icon">📝</text>
                <text class="desc-heading-text">{{ appStore.lang === 'zh' ? '作品详述' : 'Art Details & Description' }}</text>
              </view>
              
              <!-- #ifdef H5 -->
              <view class="rich-html-wrap" v-html="formattedDescription"></view>
              <!-- #endif -->
              <!-- #ifndef H5 -->
              <rich-text class="rich-html-wrap" :nodes="formattedDescription"></rich-text>
              <!-- #endif -->

              <!-- Bottom In-Page Purchase Action -->
              <view class="bottom-inpage-cta">
                <button
                  class="cta-action-btn btn-bottom-large"
                  :class="{
                    'btn-vip-gold': (userStore.isMember && isVIPTierAllowed) || Number(product?.price) === 0,
                    'btn-vip-upgrade': userStore.isMember && !isVIPTierAllowed
                  }"
                  @click="handleMainCTA"
                  @tap="handleMainCTA"
                >
                  <text class="cta-btn-icon">{{ ctaIcon }}</text>
                  <text class="cta-btn-text">{{ ctaLabel }}</text>
                </button>
              </view>
            </view>

            <!-- Universal Setup Tutorial Template Section -->
            <view
              v-if="tutorialConfig && tutorialConfig.enabled !== false && tutorialConfig.items && tutorialConfig.items.length > 0"
              class="universal-tutorial-section"
            >
              <view class="tut-header-row">
                <view class="tut-header-left">
                  <text class="tut-header-icon">📖</text>
                  <text class="tut-header-title">
                    {{ appStore.lang === 'zh' ? (tutorialConfig.title_zh || 'iOS 动态壁纸实机交互设置教程') : (tutorialConfig.title_en || 'iOS Live Wallpaper Setup Guide') }}
                  </text>
                </view>
                <text class="tut-header-badge">{{ appStore.lang === 'zh' ? '通用设置指引' : 'Universal Guide' }}</text>
              </view>

              <!-- Notice bar if provided -->
              <view v-if="tutorialNotice" class="tut-notice-bar">
                <text class="tut-notice-text">{{ tutorialNotice }}</text>
              </view>

              <!-- Version Selector Tabs -->
              <scroll-view scroll-x class="tut-versions-scroll" :show-scrollbar="false">
                <view class="tut-versions-row">
                  <view
                    v-for="(item, idx) in tutorialConfig.items"
                    :key="item.id || idx"
                    class="tut-version-tab"
                    :class="{
                      active: activeTutorialIdx === idx,
                      [item.badge_color || 'blue']: true
                    }"
                    @click="activeTutorialIdx = idx"
                    @tap="activeTutorialIdx = idx"
                  >
                    <text class="tut-tab-ver">{{ item.version_range }}</text>
                    <text v-if="item.badge_zh" class="tut-tab-tag">
                      {{ appStore.lang === 'zh' ? item.badge_zh : (item.badge_en || item.badge_zh) }}
                    </text>
                  </view>
                </view>
              </scroll-view>

              <!-- Active Tutorial Content Card -->
              <view v-if="activeTutorialItem" class="tut-card-body">
                <view class="tut-card-title-row">
                  <text class="tut-card-subtitle">
                    {{ appStore.lang === 'zh' ? activeTutorialItem.title_zh : (activeTutorialItem.title_en || activeTutorialItem.title_zh) }}
                  </text>

                  <!-- Direct Video Link Jump Button -->
                  <view
                    v-if="activeTutorialItem.video_url"
                    class="tut-video-btn"
                    @click="openTutorialVideo(activeTutorialItem.video_url)"
                    @tap="openTutorialVideo(activeTutorialItem.video_url)"
                  >
                    <text class="tut-video-icon">▶</text>
                    <text class="tut-video-text">{{ appStore.lang === 'zh' ? '观看视频教程' : 'Watch Video' }}</text>
                    <text class="tut-video-arrow">↗</text>
                  </view>
                </view>

                <!-- Steps List -->
                <view v-if="activeTutorialItem.steps && activeTutorialItem.steps.length > 0" class="tut-steps-container">
                  <view
                    v-for="(step, sIdx) in activeTutorialItem.steps"
                    :key="sIdx"
                    class="tut-step-item"
                  >
                    <view class="tut-step-badge">
                      <text class="tut-step-num">{{ sIdx + 1 }}</text>
                    </view>
                    <text class="tut-step-content">{{ step }}</text>
                  </view>
                </view>
              </view>
            </view>
          </view>
        </view>
      </view>
    </view>

    <!-- Bottom Sticky Single-CTA Bar (Mobile Only: Hidden on Desktop) -->
    <view class="sticky-cta-bar" :class="{ 'cta-bar-hidden': showCashier }">
      <view class="cta-inner-wrap">
        <view class="cta-price-info">
          <template v-if="userStore.isMember">
            <template v-if="isVIPTierAllowed">
              <text class="cta-vip-tag">👑 {{ appStore.lang === 'zh' ? 'VIP 会员特权' : 'VIP Privilege' }}</text>
              <text class="cta-price-free">{{ appStore.lang === 'zh' ? '全场免单' : 'Free Download' }}</text>
            </template>
            <template v-else>
              <text class="cta-vip-tag vip-tag-restricted">👑 {{ appStore.lang === 'zh' ? '专享会员壁纸' : 'VIP Exclusive' }}</text>
              <text class="cta-price-sub tier-alert">{{ appStore.lang === 'zh' ? `限【${allowedVIPNames}】` : 'Upgrade Required' }}</text>
            </template>
          </template>
          <template v-else-if="Number(product?.price) === 0">
            <text class="cta-free-tag">🎁 {{ appStore.lang === 'zh' ? '限时福利' : 'Free Gift' }}</text>
            <text class="cta-price-free">{{ appStore.lang === 'zh' ? '免费下载' : 'Free' }}</text>
          </template>
          <template v-else>
            <view class="cta-price-header">
              <text class="cta-price-sub">{{ appStore.lang === 'zh' ? '单款永久解锁' : 'Single Unlock' }}</text>
              <text class="cta-free-badge">{{ freeTierLabel }}</text>
            </view>
            <view class="cta-price-main">
              <text class="cta-cur">{{ appStore.currency === 'USD' ? '$' : '¥' }}</text>
              <text class="cta-num">{{ displayUnitPrice }}</text>
            </view>
          </template>
        </view>

        <!-- Main Action Button -->
        <button
          class="cta-action-btn"
          :class="{
            'btn-vip-gold': (userStore.isMember && isVIPTierAllowed) || Number(product?.price) === 0,
            'btn-vip-upgrade': userStore.isMember && !isVIPTierAllowed
          }"
          @click="handleMainCTA"
          @tap="handleMainCTA"
        >
          <text class="cta-btn-icon">{{ ctaIcon }}</text>
          <text class="cta-btn-text">{{ ctaLabel }}</text>
        </button>
      </view>
    </view>

    <!-- Cashier Bottom Sheet -->
    <CashierModal
      :visible="showCashier"
      :product="product || {}"
      @close="showCashier = false"
      @success="handlePurchaseSuccess"
    />
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { productAPI, tutorialAPI } from '@/api'
import { useUserStore } from '@/store/user'
import { useAppStore } from '@/store/app'
import MockupPreview from '@/components/MockupPreview.vue'
import CashierModal from '@/components/CashierModal.vue'
import AmbientCanvas from '@/components/AmbientCanvas.vue'

const userStore = useUserStore()
const appStore = useAppStore()
const productId = ref('')
const product = ref(null)
const detailResponse = ref(null)
const showCashier = ref(false)

// Universal Tutorial Template State
const tutorialConfig = ref(null)
const activeTutorialIdx = ref(0)

const activeTutorialItem = computed(() => {
  if (!tutorialConfig.value || !tutorialConfig.value.items || tutorialConfig.value.items.length === 0) return null
  return tutorialConfig.value.items[activeTutorialIdx.value] || tutorialConfig.value.items[0]
})

const tutorialNotice = computed(() => {
  if (!tutorialConfig.value) return ''
  return appStore.lang === 'zh'
    ? (tutorialConfig.value.notice_zh || '')
    : (tutorialConfig.value.notice_en || tutorialConfig.value.notice_zh || '')
})

const displayUnitPrice = computed(() => {
  return appStore.getProductPriceNumber(product.value).toFixed(2)
})

const freeTierLabel = computed(() => {
  return appStore.getFreeTierLabel(product.value)
})

const isVIPTierAllowed = computed(() => {
  if (Number(product.value?.price) === 0) return true
  if (userStore.userInfo?.role === 'admin') return true
  if (!userStore.isMember) return false
  if (detailResponse.value?.vip_allowed !== undefined) {
    return Boolean(detailResponse.value.vip_allowed)
  }
  const raw = product.value?.allowed_vip_levels
  if (!raw || raw === 'null' || raw === '') return true
  try {
    const list = typeof raw === 'string' ? JSON.parse(raw) : raw
    if (Array.isArray(list)) {
      if (list.length === 0) return false
      const userType = userStore.membershipType || userStore.userInfo?.membership_type
      return list.includes(userType)
    }
  } catch (e) {}
  return true
})

const allowedVIPNames = computed(() => {
  if (detailResponse.value?.allowed_vip_names) {
    return detailResponse.value.allowed_vip_names
  }
  const raw = product.value?.allowed_vip_levels
  if (!raw || raw === 'null' || raw === '') return '全等级 VIP'
  try {
    const list = typeof raw === 'string' ? JSON.parse(raw) : raw
    if (Array.isArray(list)) {
      const map = {
        monthly: '月度会员',
        quarterly: '季度会员',
        yearly: '年度会员',
        lifetime: '终生尊享'
      }
      return list.map(t => map[t] || t).join('、')
    }
  } catch (e) {}
  return '全等级 VIP'
})

const userPlanName = computed(() => {
  const type = userStore.membershipType || userStore.userInfo?.membership_type
  const map = {
    monthly: '月度会员',
    quarterly: '季度会员',
    yearly: '年度会员',
    lifetime: '终生尊享'
  }
  return map[type] || '普通会员'
})

const ctaIcon = computed(() => {
  if (userStore.isMember) {
    return isVIPTierAllowed.value ? '👑' : '⭐'
  }
  if (Number(product.value?.price) === 0) return '📥'
  return '⚡'
})

const ctaLabel = computed(() => {
  if (detailResponse.value?.is_unlocked) {
    return appStore.lang === 'zh' ? '已解锁 · 前往下载' : 'Unlocked · Download'
  }
  if (userStore.isMember) {
    if (isVIPTierAllowed.value) {
      return appStore.lang === 'zh' ? 'VIP 免单下载' : 'VIP Free Download'
    } else {
      return appStore.lang === 'zh' ? '升级特权免费享' : 'Upgrade VIP Plan'
    }
  }
  if (Number(product.value?.price) === 0) return appStore.lang === 'zh' ? '免费直接获取' : 'Get Free Wallpaper'
  return appStore.lang === 'zh' ? '立即获取此壁纸' : 'Unlock Wallpaper'
})

// Whether interactive 3D simulation is active strictly depends on
// backend configuration and whether relevant files were uploaded
const isInteractiveProduct = computed(() => {
  const p = product.value
  if (!p || !p.is_interactive || !p.interactive_config) return false
  try {
    const cfg = typeof p.interactive_config === 'string' ? JSON.parse(p.interactive_config) : p.interactive_config
    return Boolean(
      (cfg.foreground_url && cfg.foreground_url.trim()) ||
      (cfg.background_url && cfg.background_url.trim()) ||
      (cfg.motion_url && cfg.motion_url.trim()) ||
      (cfg.tendies_file && cfg.tendies_file.trim())
    )
  } catch (e) {
    return false
  }
})

// Video previews list
const videoLinks = computed(() => {
  const p = product.value
  if (!p) return []
  const list = []
  if (p.preview_video_url) {
    list.push({
      platform: p.preview_platform || 'direct',
      url: p.preview_video_url,
      title: appStore.lang === 'zh' ? '实机动态效果演示' : 'Live Motion Demo'
    })
  }
  if (p.preview_videos) {
    try {
      const parsed = typeof p.preview_videos === 'string' ? JSON.parse(p.preview_videos) : p.preview_videos
      if (Array.isArray(parsed)) {
        list.push(...parsed)
      }
    } catch (e) {}
  }
  return list
})

function getPlatformName(platform) {
  const map = {
    douyin: '抖音',
    tiktok: 'TikTok',
    xiaohongshu: '小红书',
    bilibili: '哔哩哔哩',
    direct: '在线视频'
  }
  return map[platform] || '在线视频'
}

function openVideoLink(url) {
  if (!url) return
  // #ifdef H5
  window.open(url, '_blank')
  // #endif
  // #ifndef H5
  uni.setClipboardData({
    data: url,
    success: () => {
      uni.showToast({ title: '演示链接已复制', icon: 'none' })
    }
  })
  // #endif
}

// Formatted Rich Text for Product Description (HTML Rendering)
const formattedDescription = computed(() => {
  const p = product.value
  let desc = appStore.lang === 'zh' ? (p?.description || '') : (p?.description_en || p?.description || '')
  if (!desc) {
    return appStore.lang === 'zh'
      ? '<p>本款壁纸经由专业色彩校准与超分质感渲染，完美适配各大品牌手机与高分视网膜屏幕。</p>'
      : '<p>Engineered with pristine color grading and ultra-resolution detail, meticulously calibrated for all modern Retina displays.</p>'
  }

  // If plain text (no HTML tags detected), split by lines and wrap into <p>
  if (!/<[a-z][\s\S]*>/i.test(desc)) {
    return desc
      .split('\n')
      .map(line => line.trim())
      .filter(Boolean)
      .map(line => `<p>${line}</p>`)
      .join('')
  }

  // Ensure embedded images in rich text have max-width & aesthetic styling
  desc = desc.replace(/<img\s+/gi, '<img style="max-width:100%;height:auto;border-radius:12px;margin:10px 0;display:block;" ')
  return desc
})

function goBack() {
  // Always navigate directly to home gallery tab instead of navigating back in browser history
  uni.switchTab({
    url: '/pages/index/index',
    fail: () => {
      uni.reLaunch({
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

function handleMainCTA() {
  if (detailResponse.value?.is_unlocked || Number(product.value?.price) === 0) {
    const lastOrderId = uni.getStorageSync('last_paid_order_id') || ''
    const lastEmail = uni.getStorageSync('last_buyer_email') || ''
    uni.navigateTo({
      url: `/pages/download/index?id=${product.value?.id}&title=${encodeURIComponent(product.value?.title_zh || '')}&order_id=${encodeURIComponent(lastOrderId)}&email=${encodeURIComponent(lastEmail)}`
    })
    return
  }

  if (userStore.isMember) {
    if (isVIPTierAllowed.value) {
      uni.navigateTo({
        url: `/pages/download/index?id=${product.value?.id}&title=${encodeURIComponent(product.value?.title_zh || '')}`
      })
    } else {
      uni.showModal({
        title: '👑 会员等级特权提示',
        content: `本款壁纸仅对【${allowedVIPNames.value}】免费开放下载，您当前为【${userPlanName.value}】。点击确定前往升级会员，立即解锁全站高阶壁纸！`,
        confirmText: '升级会员',
        cancelText: '单款购买',
        success: (res) => {
          if (res.confirm) {
            goToMembership()
          } else {
            showCashier.value = true
          }
        }
      })
    }
  } else {
    showCashier.value = true
  }
}

function handlePurchaseSuccess() {
  showCashier.value = false
}

async function fetchDetail(id) {
  try {
    const params = {}
    const lastOrderId = uni.getStorageSync('last_paid_order_id')
    const lastEmail = uni.getStorageSync('last_buyer_email')
    if (lastOrderId) params.order_id = lastOrderId
    if (lastEmail) params.email = lastEmail

    const res = await productAPI.get(id, params)
    detailResponse.value = res
    product.value = res.product || res.data || res
  } catch (err) {
    console.error('Failed to load product detail:', err)
    uni.showToast({
      title: err?.message || (appStore.lang === 'zh' ? '壁纸详情加载失败，请重试' : 'Failed to load details'),
      icon: 'none',
      duration: 3000
    })
  }
}

async function loadTutorialConfig() {
  try {
    const res = await tutorialAPI.get()
    if (res && res.enabled !== false) {
      tutorialConfig.value = res
    }
  } catch (err) {
    console.warn('Failed to load universal tutorial config:', err)
  }
}

function openTutorialVideo(url) {
  if (!url) return
  // #ifdef H5
  if (typeof window !== 'undefined') {
    window.open(url, '_blank')
    return
  }
  // #endif
  uni.setClipboardData({
    data: url,
    success: () => {
      uni.showToast({
        title: appStore.lang === 'zh' ? '已复制视频链接，请到浏览器打开' : 'Link copied to clipboard',
        icon: 'none',
        duration: 2500
      })
    }
  })
}

onLoad(async (options) => {
  productId.value = options.id || ''
  await Promise.all([
    fetchDetail(productId.value),
    loadTutorialConfig()
  ])
  if (options.action === 'buy') {
    handleMainCTA()
  }
})
</script>

<style scoped>
.detail-container {
  min-height: 100vh;
  position: relative;
  transition: background-color 0.25s ease, color 0.25s ease;
  padding-bottom: calc(140px + env(safe-area-inset-bottom));
}

@media (min-width: 900px) {
  .detail-container {
    padding-bottom: 60px;
  }
}

.detail-container.dark {
  background: transparent;
  color: #f1f2f6;
}

.detail-container.light {
  background: transparent;
  color: #1a1a1e;
}

.detail-content-wrap {
  position: relative;
  z-index: 2;
  width: 100%;
  max-width: 1240px;
  margin: 0 auto;
  padding: 0 24px;
  box-sizing: border-box;
}

@media (max-width: 768px) {
  .detail-content-wrap {
    padding: 0 12px;
  }
  .nav-bar {
    padding: 12px 0 10px !important;
  }
  .back-circle {
    height: 34px !important;
    padding: 0 12px !important;
  }
  .icon-circle {
    height: 32px !important;
    padding: 0 10px !important;
  }
}

/* Nav Bar */
.nav-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24px 0 16px;
}

.back-circle {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  height: 40px;
  padding: 0 16px;
  border-radius: 999px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.dark .back-circle {
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.light .back-circle {
  background: rgba(0, 0, 0, 0.05);
  border: 1px solid rgba(0, 0, 0, 0.08);
}

.back-circle:hover {
  transform: translateX(-2px);
}

.back-icon {
  font-size: 16px;
}

.back-label {
  font-size: 14px;
  font-weight: 600;
}

.nav-title {
  font-size: 16px;
  font-weight: 700;
  max-width: 420px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

@media (max-width: 768px) {
  .nav-title {
    display: none;
  }
}

.nav-right-controls {
  display: flex;
  align-items: center;
  gap: 10px;
}

.vip-capsule-btn {
  height: 36px;
  padding: 0 14px;
  border-radius: 999px;
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.22), rgba(212, 175, 55, 0.12));
  border: 1px solid rgba(245, 158, 11, 0.45);
  transition: all 0.2s ease;
  box-shadow: 0 2px 10px rgba(245, 158, 11, 0.15);
}

.vip-capsule-btn:active {
  transform: scale(0.96);
}

.vip-capsule-icon {
  font-size: 14px;
}

.vip-capsule-text {
  font-size: 13px;
  font-weight: 800;
  color: #f59e0b;
}

.my-nav-circle {
  padding: 0 10px !important;
}

.icon-circle {
  height: 36px;
  padding: 0 14px;
  border-radius: 999px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
}

.dark .icon-circle {
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.12);
  color: #ffffff;
}

.light .icon-circle {
  background: #ffffff;
  border: 1px solid rgba(0, 0, 0, 0.08);
  color: #111827;
}

.nav-btn-text {
  font-size: 13px;
  font-weight: 700;
}

.nav-btn-icon {
  font-size: 16px;
}

/* ========================================================
   Main Layout: 2-Columns on Desktop, Stacked on Mobile
   ======================================================== */
.detail-main-layout {
  display: flex;
  flex-direction: column;
}

@media (min-width: 900px) {
  .detail-main-layout {
    display: grid;
    grid-template-columns: 480px 1fr;
    gap: 48px;
    align-items: start;
    padding-top: 16px;
  }

  .mockup-column {
    position: sticky;
    top: 24px;
  }
}

/* Mockup wrap */
.mockup-wrap {
  padding: 10px 0 16px;
}

/* Live Video Preview Strip */
.video-preview-strip {
  margin-top: 12px;
  padding: 14px 18px;
  border-radius: 16px;
  backdrop-filter: blur(12px);
}

.dark .video-preview-strip {
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.light .video-preview-strip {
  background: rgba(0, 0, 0, 0.03);
  border: 1px solid rgba(0, 0, 0, 0.06);
}

.v-strip-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
}

.v-strip-icon {
  font-size: 16px;
}

.v-strip-title {
  font-size: 13px;
  font-weight: 700;
  letter-spacing: 0.5px;
}

.dark .v-strip-title {
  color: #e2e8f0;
}

.light .v-strip-title {
  color: #334155;
}

.v-strip-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.v-link-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  border-radius: 999px;
  background: rgba(59, 130, 246, 0.14);
  border: 1px solid rgba(59, 130, 246, 0.28);
  cursor: pointer;
  transition: all 0.2s ease;
}

.v-link-chip:hover {
  background: rgba(59, 130, 246, 0.25);
  transform: translateY(-1px);
}

.v-chip-icon {
  font-size: 11px;
  color: #3b82f6;
}

.v-chip-name {
  font-size: 12px;
  font-weight: 600;
  color: #3b82f6;
}

.v-chip-jump {
  font-size: 11px;
  color: #3b82f6;
  opacity: 0.7;
}

/* Info block */
.info-block {
  padding: 10px 0 40px;
}

.tag-row {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-bottom: 16px;
}

.chip {
  padding: 4px 12px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 600;
}

.chip-cat {
  background: rgba(59, 130, 246, 0.15);
  color: #60a5fa;
}

.chip-res {
  background: rgba(16, 185, 129, 0.15);
  color: #34d399;
}

.chip-cloud {
  background: rgba(255, 255, 255, 0.06);
  opacity: 0.8;
}

.chip-interactive {
  background: linear-gradient(135deg, rgba(139, 92, 246, 0.25), rgba(99, 102, 241, 0.25));
  border: 1px solid rgba(139, 92, 246, 0.45);
  color: #a78bfa;
  font-weight: 700;
  box-shadow: 0 0 10px rgba(139, 92, 246, 0.25);
}

.highlight-purple {
  color: #a78bfa !important;
  font-weight: 700 !important;
}

.product-title {
  font-size: 28px;
  font-weight: 800;
  line-height: 1.35;
  display: block;
  margin-bottom: 24px;
}

@media (max-width: 768px) {
  .product-title {
    font-size: 22px;
  }
}

/* Main Inline Action Card (Visible on All Devices) */
.main-cta-box {
  display: flex !important;
  align-items: center;
  justify-content: space-between;
  padding: 18px 22px;
  border-radius: 20px;
  margin-bottom: 22px;
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  box-sizing: border-box;
}

.dark .main-cta-box {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.12);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.25);
}

.light .main-cta-box {
  background: #ffffff;
  border: 1px solid rgba(0, 0, 0, 0.08);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.08);
}

@media (max-width: 768px) {
  .main-cta-box {
    flex-direction: column;
    align-items: stretch;
    gap: 14px;
    padding: 16px;
    border-radius: 16px;
  }
  .main-cta-box .cta-action-btn {
    width: 100%;
    justify-content: center;
    height: 48px;
    font-size: 15px;
  }
}

.bottom-inpage-cta {
  margin-top: 24px;
  display: flex;
  justify-content: center;
}

.btn-bottom-large {
  width: 100%;
  max-width: 480px;
  height: 52px;
  justify-content: center;
  font-size: 16px;
}

/* Specs Grid */
.specs-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  padding: 20px;
  border-radius: 20px;
  margin-bottom: 28px;
}

.dark .specs-grid {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
}

.light .specs-grid {
  background: #ffffff;
  border: 1px solid rgba(0, 0, 0, 0.06);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.04);
}

.spec-cell {
  display: flex;
  flex-direction: column;
}

.spec-label {
  font-size: 11px;
  opacity: 0.55;
  margin-bottom: 6px;
}

.spec-val {
  font-size: 13px;
  font-weight: 700;
}

/* ========================================================
   Rich Text Description Section (Formatted HTML)
   ======================================================== */
.rich-desc-container {
  padding-top: 24px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
}

.light .rich-desc-container {
  border-top-color: rgba(0, 0, 0, 0.08);
}

.desc-heading-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 14px;
}

.desc-heading-icon {
  font-size: 16px;
}

.desc-heading-text {
  font-size: 15px;
  font-weight: 700;
  letter-spacing: 0.5px;
}

.dark .desc-heading-text {
  color: #f1f5f9;
}

.light .desc-heading-text {
  color: #1e293b;
}

.rich-html-wrap {
  font-size: 14px;
  line-height: 1.85;
  letter-spacing: 0.2px;
  word-break: break-word;
}

.dark .rich-html-wrap {
  color: #94a3b8;
}

.light .rich-html-wrap {
  color: #475569;
}

.rich-html-wrap :deep(p) {
  margin: 0 0 12px 0;
}

.rich-html-wrap :deep(strong),
.rich-html-wrap :deep(b) {
  font-weight: 700;
}

.dark .rich-html-wrap :deep(strong),
.dark .rich-html-wrap :deep(b) {
  color: #f8fafc;
}

.light .rich-html-wrap :deep(strong),
.light .rich-html-wrap :deep(b) {
  color: #0f172a;
}

.rich-html-wrap :deep(ul) {
  padding-left: 22px;
  margin: 8px 0 14px 0;
}

.rich-html-wrap :deep(ol) {
  padding-left: 22px;
  margin: 8px 0 14px 0;
}

.rich-html-wrap :deep(li) {
  margin-bottom: 6px;
}

.rich-html-wrap :deep(h1),
.rich-html-wrap :deep(h2),
.rich-html-wrap :deep(h3),
.rich-html-wrap :deep(h4) {
  margin: 20px 0 10px 0;
  font-weight: 700;
  line-height: 1.4;
}

.dark .rich-html-wrap :deep(h1),
.dark .rich-html-wrap :deep(h2),
.dark .rich-html-wrap :deep(h3),
.dark .rich-html-wrap :deep(h4) {
  color: #ffffff;
}

.light .rich-html-wrap :deep(h1),
.light .rich-html-wrap :deep(h2),
.light .rich-html-wrap :deep(h3),
.light .rich-html-wrap :deep(h4) {
  color: #0f172a;
}

.rich-html-wrap :deep(blockquote) {
  border-left: 3px solid #3b82f6;
  padding: 8px 16px;
  margin: 14px 0;
  border-radius: 0 8px 8px 0;
}

.dark .rich-html-wrap :deep(blockquote) {
  background: rgba(59, 130, 246, 0.08);
  color: #cbd5e1;
}

.light .rich-html-wrap :deep(blockquote) {
  background: rgba(59, 130, 246, 0.06);
  color: #334155;
}

.rich-html-wrap :deep(img) {
  max-width: 100% !important;
  height: auto !important;
  border-radius: 12px;
  margin: 12px 0;
  display: block;
}

/* ========================================================
   Bottom Sticky CTA Bar (Mobile Only: Hidden on Desktop)
   ======================================================== */
.sticky-cta-bar {
  position: fixed !important;
  bottom: 0 !important;
  left: 0 !important;
  right: 0 !important;
  width: 100% !important;
  height: calc(72px + env(safe-area-inset-bottom)) !important;
  backdrop-filter: blur(28px) !important;
  -webkit-backdrop-filter: blur(28px) !important;
  padding: 0 16px env(safe-area-inset-bottom) !important;
  box-sizing: border-box !important;
  z-index: 1000 !important;
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  box-shadow: 0 -4px 24px rgba(0, 0, 0, 0.15) !important;
  transition: transform 0.28s cubic-bezier(0.16, 1, 0.3, 1), opacity 0.22s ease !important;
}

.sticky-cta-bar.cta-bar-hidden {
  transform: translateY(115%) !important;
  opacity: 0 !important;
  pointer-events: none !important;
  visibility: hidden !important;
}

@media (min-width: 900px) {
  .sticky-cta-bar {
    display: none !important;
  }
}

.dark .sticky-cta-bar {
  background: rgba(18, 19, 25, 0.96) !important;
  border-top: 1px solid rgba(255, 255, 255, 0.12) !important;
}

.light .sticky-cta-bar {
  background: rgba(255, 255, 255, 0.96) !important;
  border-top: 1px solid rgba(0, 0, 0, 0.08) !important;
}

.cta-inner-wrap {
  width: 100%;
  max-width: 600px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.cta-price-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.cta-price-header {
  display: flex;
  align-items: center;
  gap: 6px;
}

.cta-free-badge {
  font-size: 10px;
  font-weight: 700;
  color: #f59e0b;
  background: rgba(245, 158, 11, 0.12);
  border: 1px solid rgba(245, 158, 11, 0.28);
  padding: 1px 6px;
  border-radius: 999px;
  white-space: nowrap;
}

.cta-price-sub {
  font-size: 11px;
  opacity: 0.55;
}

.cta-price-main {
  display: flex;
  align-items: baseline;
}

.cta-cur {
  font-size: 14px;
  color: #3b82f6;
  font-weight: 700;
}

.cta-num {
  font-size: 24px;
  font-weight: 800;
}

.cta-vip-tag {
  font-size: 11px;
  color: #fbbf24;
  font-weight: 700;
}

.cta-free-tag {
  font-size: 11px;
  color: #10b981;
  font-weight: 700;
}

.cta-price-free {
  font-size: 20px;
  font-weight: 800;
}

.cta-action-btn {
  height: 48px;
  padding: 0 28px;
  background: #2563eb;
  color: #ffffff;
  border-radius: 999px;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  box-shadow: 0 8px 24px rgba(37, 99, 235, 0.35);
  cursor: pointer;
  transition: all 0.2s ease;
  border: none;
}

.cta-action-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 28px rgba(37, 99, 235, 0.45);
}

.btn-vip-gold {
  background: linear-gradient(135deg, #d4af37, #f59e0b);
  box-shadow: 0 8px 24px rgba(212, 175, 55, 0.35);
}

.btn-vip-gold:hover {
  box-shadow: 0 12px 28px rgba(212, 175, 55, 0.45);
}

.btn-vip-upgrade {
  background: linear-gradient(135deg, #8b5cf6, #ec4899);
  box-shadow: 0 8px 24px rgba(139, 92, 246, 0.35);
}

.btn-vip-upgrade:hover {
  box-shadow: 0 12px 28px rgba(139, 92, 246, 0.45);
}

.vip-tag-restricted {
  color: #c084fc;
}

.tier-alert {
  font-size: 13px;
  font-weight: 700;
  color: #fbbf24;
}

.cta-btn-icon {
  font-size: 16px;
}

.cta-btn-text {
  font-size: 15px;
  font-weight: 700;
  color: #ffffff;
}

/* ========================================================
   Universal Setup Tutorial Template Section
   ======================================================== */
.universal-tutorial-section {
  margin-top: 24px;
  padding: 22px;
  border-radius: 20px;
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  transition: all 0.25s ease;
}

.dark .universal-tutorial-section {
  background: rgba(30, 41, 59, 0.45);
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.25);
}

.light .universal-tutorial-section {
  background: rgba(255, 255, 255, 0.85);
  border: 1px solid rgba(0, 0, 0, 0.08);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.06);
}

.tut-header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.tut-header-left {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.tut-header-icon {
  font-size: 18px;
}

.tut-header-title {
  font-size: 15px;
  font-weight: 800;
  letter-spacing: -0.2px;
}

.dark .tut-header-title {
  color: #f8fafc;
}

.light .tut-header-title {
  color: #0f172a;
}

.tut-header-badge {
  font-size: 11px;
  font-weight: 700;
  padding: 3px 10px;
  border-radius: 999px;
  background: rgba(56, 189, 248, 0.15);
  color: #0ea5e9;
  letter-spacing: 0.3px;
}

.tut-notice-bar {
  padding: 9px 14px;
  border-radius: 10px;
  background: rgba(245, 158, 11, 0.12);
  border: 1px solid rgba(245, 158, 11, 0.25);
  margin-bottom: 14px;
}

.tut-notice-text {
  font-size: 12px;
  color: #d97706;
  line-height: 1.4;
}

.dark .tut-notice-text {
  color: #fbbf24;
}

/* Version Tabs Scroll */
.tut-versions-scroll {
  width: 100%;
  white-space: nowrap;
  margin-bottom: 16px;
}

.tut-versions-row {
  display: inline-flex;
  gap: 8px;
  padding: 2px 0;
}

.tut-version-tab {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 7px 14px;
  border-radius: 999px;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  white-space: nowrap;
}

.dark .tut-version-tab {
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.08);
  color: #94a3b8;
}

.light .tut-version-tab {
  background: rgba(0, 0, 0, 0.04);
  border: 1px solid rgba(0, 0, 0, 0.06);
  color: #64748b;
}

.tut-version-tab.active {
  background: linear-gradient(135deg, rgba(14, 165, 233, 0.2), rgba(2, 132, 199, 0.3));
  border-color: rgba(14, 165, 233, 0.6);
  color: #0284c7;
  transform: translateY(-1px);
}

.dark .tut-version-tab.active {
  color: #38bdf8;
  box-shadow: 0 4px 14px rgba(14, 165, 233, 0.25);
}

.tut-version-tab.active.purple {
  background: linear-gradient(135deg, rgba(168, 85, 247, 0.2), rgba(147, 51, 234, 0.3));
  border-color: rgba(168, 85, 247, 0.6);
  color: #a855f7;
}

.dark .tut-version-tab.active.purple {
  color: #c084fc;
  box-shadow: 0 4px 14px rgba(168, 85, 247, 0.25);
}

.tut-version-tab.active.emerald {
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.2), rgba(5, 150, 105, 0.3));
  border-color: rgba(16, 185, 129, 0.6);
  color: #059669;
}

.dark .tut-version-tab.active.emerald {
  color: #34d399;
  box-shadow: 0 4px 14px rgba(16, 185, 129, 0.25);
}

.tut-version-tab.active.amber {
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.2), rgba(217, 119, 6, 0.3));
  border-color: rgba(245, 158, 11, 0.6);
  color: #d97706;
}

.dark .tut-version-tab.active.amber {
  color: #fbbf24;
  box-shadow: 0 4px 14px rgba(245, 158, 11, 0.25);
}

.tut-tab-ver {
  font-size: 13px;
  font-weight: 700;
}

.tut-tab-tag {
  font-size: 10px;
  font-weight: 800;
  padding: 1px 6px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.2);
}

/* Card Body */
.tut-card-body {
  border-radius: 14px;
  padding: 16px;
  transition: all 0.2s ease;
}

.dark .tut-card-body {
  background: rgba(15, 23, 42, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.06);
}

.light .tut-card-body {
  background: rgba(241, 245, 249, 0.7);
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.tut-card-title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
  flex-wrap: wrap;
}

.tut-card-subtitle {
  font-size: 14px;
  font-weight: 700;
  flex: 1;
}

.dark .tut-card-subtitle {
  color: #e2e8f0;
}

.light .tut-card-subtitle {
  color: #1e293b;
}

.tut-video-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  border-radius: 999px;
  cursor: pointer;
  background: linear-gradient(135deg, #0284c7, #2563eb);
  box-shadow: 0 3px 12px rgba(37, 99, 235, 0.35);
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.tut-video-btn:active {
  transform: scale(0.97);
}

.tut-video-icon {
  font-size: 11px;
  color: #ffffff;
}

.tut-video-text {
  font-size: 12px;
  font-weight: 800;
  color: #ffffff;
}

.tut-video-arrow {
  font-size: 11px;
  color: #ffffff;
}

/* Steps List */
.tut-steps-container {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.tut-step-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
}

.tut-step-badge {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: rgba(14, 165, 233, 0.15);
  border: 1px solid rgba(14, 165, 233, 0.3);
  flex-shrink: 0;
  margin-top: 1px;
}

.tut-step-num {
  font-size: 11px;
  font-weight: 800;
  color: #0ea5e9;
  line-height: 1;
}

.tut-step-content {
  font-size: 13px;
  line-height: 1.5;
}

.dark .tut-step-content {
  color: #cbd5e1;
}

.light .tut-step-content {
  color: #334155;
}
</style>
