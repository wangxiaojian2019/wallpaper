<template>
  <view class="wallpaper-card" :class="[appStore.theme]" @tap="handleTap">
    <view class="card-inner">
      <!-- Image cover -->
      <image
        v-if="product.cover_image"
        :src="product.cover_image"
        mode="aspectFill"
        class="card-img"
        lazy-load
      />
      <view
        v-else
        class="card-placeholder"
        :style="{ background: placeholderBg }"
      >
        <text class="placeholder-emoji">{{ placeholderEmoji }}</text>
      </view>

      <!-- Top Badges (Tech spec & Free tier eligibility) -->
      <view class="badge-row">
        <view class="badge-left-group">
          <!-- Pinned badge -->
          <view v-if="product.is_pinned" class="badge badge-pinned">
            <text class="badge-pinned-text">📌 {{ appStore.lang === 'zh' ? '置顶' : 'TOP' }}</text>
          </view>
          <!-- Tech badge -->
          <view class="badge badge-tech">
            <text class="badge-tech-text">{{ techBadge }}</text>
          </view>
        </view>

        <!-- Tier / Free Eligibility Badge (Shows which VIP level is free) -->
        <view
          class="badge badge-tier"
          :class="{
            'is-free': isZeroPrice,
            'is-vip': !isZeroPrice && !isPurchaseOnly,
            'is-purchase-only': isPurchaseOnly
          }"
        >
          <text class="badge-tier-text">{{ freeTierLabel }}</text>
        </view>
      </view>

      <!-- Signature iOS 26 Floating Frosted Glass Dock -->
      <view class="glass-dock" @tap.stop="handleTap">
        <view class="dock-title-wrap">
          <text class="dock-title">{{ cleanTitle }}</text>
        </view>

        <view class="dock-meta-row">
          <view class="dock-left-meta">
            <view class="dock-cat-pill">
              <view class="cat-dot"></view>
              <text class="cat-name">{{ categoryLabel }}</text>
            </view>

            <!-- Unit Price Badge -->
            <view class="dock-price-pill" :class="{ 'is-zero': isZeroPrice }">
              <text class="price-symbol">{{ appStore.currency === 'USD' ? '$' : '¥' }}</text>
              <text class="price-num">{{ displayUnitPrice }}</text>
            </view>
          </view>

          <!-- Quick Action / Buy Button -->
          <button
            class="dock-buy-btn"
            :class="{ 'is-vip': isUserAllowedFree, 'is-free': isZeroPrice }"
            @tap.stop="handleBuyTap"
          >
            <text class="buy-btn-icon">{{ buyBtnIcon }}</text>
            <text class="buy-btn-text">{{ buyBtnText }}</text>
          </button>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed } from 'vue'
import { useUserStore } from '@/store/user'
import { useAppStore } from '@/store/app'

const props = defineProps({
  product: {
    type: Object,
    required: true
  }
})

const userStore = useUserStore()
const appStore = useAppStore()

// Unit price calculation
const displayUnitPrice = computed(() => {
  return appStore.getProductPriceNumber(props.product).toFixed(2)
})

const isZeroPrice = computed(() => {
  const p = props.product
  const pCny = Number(p.price_cny || 0)
  const pUsd = Number(p.price_usd || 0)
  const pBase = Number(p.price || 0)
  return pCny === 0 && pUsd === 0 && pBase === 0
})

// Check if product is set to purchase-only (empty allowed_vip_levels)
const isPurchaseOnly = computed(() => {
  if (isZeroPrice.value) return false
  const raw = props.product?.allowed_vip_levels
  if (!raw || raw === 'null' || raw === '') return false
  try {
    const list = typeof raw === 'string' ? JSON.parse(raw) : raw
    if (Array.isArray(list) && list.length === 0) return true
  } catch (e) {}
  return false
})

// Free Tier Requirement (e.g. 月度会员可免 / 季度会员可免 / 年度会员可免 / 终身会员可免)
const freeTierLabel = computed(() => {
  if (isZeroPrice.value) {
    return appStore.lang === 'zh' ? '🎁 全员免费' : '🎁 Free for All'
  }
  return appStore.getFreeTierLabel(props.product)
})

// Check if current user qualifies for free download
const isUserAllowedFree = computed(() => {
  if (isZeroPrice.value) return true
  if (!userStore.isMember) return false
  if (userStore.userInfo?.role === 'admin') return true
  if (userStore.membershipType === 'lifetime') return true

  const raw = props.product?.allowed_vip_levels
  if (!raw || raw === 'null' || raw === '') return true
  try {
    const list = typeof raw === 'string' ? JSON.parse(raw) : raw
    if (Array.isArray(list)) {
      if (list.length === 0) return false
      return list.includes(userStore.membershipType)
    }
  } catch (e) {}
  return true
})

const cleanTitle = computed(() => {
  const p = props.product
  let title = appStore.lang === 'zh' ? (p.title_zh || p.title_en) : (p.title_en || p.title_zh)
  if (!title || title.startsWith('???') || title.includes('????')) {
    title = p.title_en || '4K 精选壁纸'
  }
  return title
})

const techBadge = computed(() => {
  const cat = props.product.category || ''
  if (cat.includes('动态') || props.product.preview_platform) return '🎬 Live'
  if (props.product.tags?.includes('8K')) return '8K'
  return '4K UHD'
})

const buyBtnIcon = computed(() => {
  if (userStore.isMember && isUserAllowedFree.value) return '👑'
  if (isZeroPrice.value) return '📥'
  return '⚡'
})

const buyBtnText = computed(() => {
  if (userStore.isMember && isUserAllowedFree.value) {
    return appStore.lang === 'zh' ? '免单' : 'Free'
  }
  if (isZeroPrice.value) {
    return appStore.lang === 'zh' ? '免费' : 'Free'
  }
  return appStore.lang === 'zh' ? '购买' : 'Buy'
})

const categoryLabel = computed(() => {
  const cat = props.product.category || ''
  if (!cat) return appStore.lang === 'zh' ? '精选' : 'Curated'
  const map = {
    '苹果壁纸': { zh: '苹果风', en: 'iOS' },
    '动态壁纸': { zh: '动态4K', en: 'Live' },
    '电脑壁纸': { zh: '桌面', en: 'Desktop' },
    '动漫': { zh: '动漫', en: 'Anime' },
    'anime': { zh: '动漫', en: 'Anime' },
    '赛博朋克': { zh: '赛博', en: 'Cyber' },
    'cyberpunk': { zh: '赛博', en: 'Cyber' },
    '自然': { zh: '自然', en: 'Nature' },
    'nature': { zh: '自然', en: 'Nature' },
    '极简': { zh: '极简', en: 'Minimal' },
    'minimal': { zh: '极简', en: 'Minimal' },
    'City': { zh: '城市', en: 'City' }
  }
  const item = map[cat]
  if (item) return appStore.lang === 'zh' ? item.zh : item.en
  return cat
})

const placeholderBg = computed(() => {
  const gradients = [
    'linear-gradient(135deg, #1e293b, #334155)',
    'linear-gradient(135deg, #18181b, #27272a)',
    'linear-gradient(135deg, #1e1b4b, #312e81)',
    'linear-gradient(135deg, #2e1065, #4c1d95)'
  ]
  const idx = Math.abs((props.product.id || '').toString().charCodeAt(0) || 0) % gradients.length
  return gradients[idx]
})

const placeholderEmoji = computed(() => {
  const emojis = ['🌌', '✨', '🌊', '🏔️', '🏙️', '🎨']
  const idx = Math.abs((props.product.id || '').toString().charCodeAt(0) || 0) % emojis.length
  return emojis[idx]
})

function handleTap() {
  uni.navigateTo({
    url: `/pages/detail/index?id=${props.product.id}`
  })
}

function handleBuyTap() {
  uni.navigateTo({
    url: `/pages/detail/index?id=${props.product.id}&action=buy`
  })
}
</script>

<style scoped>
.wallpaper-card {
  width: 100%;
  box-sizing: border-box;
  cursor: pointer;
}

.card-inner {
  position: relative;
  width: 100%;
  aspect-ratio: 9 / 16;
  border-radius: 20px;
  overflow: hidden;
  transition: transform 0.25s cubic-bezier(0.16, 1, 0.3, 1), box-shadow 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}

.dark .card-inner {
  background: #181920;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.35);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.light .card-inner {
  background: #ffffff;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.07);
  border: 1px solid rgba(0, 0, 0, 0.06);
}

/* #ifdef H5 */
.card-inner:hover {
  transform: translateY(-6px);
}
.dark .card-inner:hover {
  box-shadow: 0 16px 36px rgba(0, 0, 0, 0.5), 0 0 20px rgba(59, 130, 246, 0.2);
}
.light .card-inner:hover {
  box-shadow: 0 12px 28px rgba(0, 0, 0, 0.12);
}
/* #endif */

.card-img {
  width: 100%;
  height: 100%;
  display: block;
}

.card-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.placeholder-emoji {
  font-size: 48px;
  opacity: 0.7;
}

/* Top Glass Badges (Compact & Sleek) */
.badge-row {
  position: absolute;
  top: 10px;
  left: 10px;
  right: 10px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  z-index: 5;
  pointer-events: none;
}

.badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  white-space: nowrap;
  flex-shrink: 0;
  height: 24px;
  padding: 0 10px;
  border-radius: 999px;
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
}

.badge-left-group {
  display: flex;
  align-items: center;
  gap: 6px;
}

.badge-pinned {
  background: linear-gradient(135deg, #ef4444, #dc2626);
  border: 1px solid rgba(255, 255, 255, 0.4);
  box-shadow: 0 2px 8px rgba(239, 68, 68, 0.45);
}

.badge-pinned-text {
  font-size: 11px;
  font-weight: 800;
  color: #ffffff;
  letter-spacing: 0.3px;
  line-height: 1;
}

.badge-tech {
  background: rgba(0, 0, 0, 0.55);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.badge-tech-text {
  font-size: 11px;
  font-weight: 700;
  color: #ffffff;
  letter-spacing: 0.5px;
  line-height: 1;
}

.badge-tier {
  background: rgba(255, 255, 255, 0.94);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.16);
}

.badge-tier.is-free {
  background: linear-gradient(135deg, #10b981, #059669);
}

.badge-tier.is-free .badge-tier-text {
  color: #ffffff;
}

.badge-tier.is-vip {
  background: linear-gradient(135deg, #f59e0b, #d97706);
}

.badge-tier.is-vip .badge-tier-text {
  color: #ffffff;
}

.badge-tier.is-purchase-only {
  background: rgba(15, 23, 42, 0.75);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.badge-tier.is-purchase-only .badge-tier-text {
  color: #cbd5e1;
}

.badge-tier-text {
  font-size: 11px;
  font-weight: 800;
  color: #1a1a1a;
  line-height: 1;
  white-space: nowrap;
}

/* ========================================================
   iOS 26 Signature Floating Frosted Glass Dock
   ======================================================== */
.glass-dock {
  position: absolute;
  bottom: 10px;
  left: 10px;
  right: 10px;
  border-radius: 16px;
  padding: 10px 12px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  backdrop-filter: blur(28px) saturate(210%);
  -webkit-backdrop-filter: blur(28px) saturate(210%);
  z-index: 10;
  box-sizing: border-box;
  transition: all 0.25s ease;
}

.dark .glass-dock {
  background: rgba(14, 16, 24, 0.55);
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.45), inset 0 1px 1.5px rgba(255, 255, 255, 0.25);
}

.light .glass-dock {
  background: rgba(255, 255, 255, 0.65);
  border: 1px solid rgba(255, 255, 255, 0.7);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.1), inset 0 1px 2px rgba(255, 255, 255, 0.9);
}

.dock-title-wrap {
  width: 100%;
  overflow: hidden;
}

.dock-title {
  font-size: 13px;
  font-weight: 700;
  line-height: 1.35;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  display: block;
}

.dark .dock-title {
  color: #ffffff;
  text-shadow: 0 1px 4px rgba(0, 0, 0, 0.9);
}

.light .dock-title {
  color: #111827;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.8);
}

.dock-meta-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 6px;
}

.dock-left-meta {
  display: flex;
  align-items: center;
  gap: 6px;
  min-width: 0;
}

.dock-cat-pill {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  background: rgba(0, 122, 255, 0.16);
  padding: 2px 7px;
  border-radius: 999px;
  white-space: nowrap;
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
}

.cat-dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: #38bdf8;
  flex-shrink: 0;
}

.cat-name {
  font-size: 10px;
  font-weight: 700;
  color: #38bdf8;
  line-height: 1;
  white-space: nowrap;
}

/* Price Pill */
.dock-price-pill {
  display: inline-flex;
  align-items: baseline;
  background: rgba(16, 185, 129, 0.18);
  border: 1px solid rgba(16, 185, 129, 0.4);
  padding: 2px 7px;
  border-radius: 999px;
  font-weight: 800;
  color: #10b981;
  white-space: nowrap;
}

.dark .dock-price-pill {
  background: rgba(16, 185, 129, 0.22);
  border-color: rgba(52, 211, 153, 0.45);
  color: #34d399;
}

.price-symbol {
  font-size: 10px;
  margin-right: 1px;
}

.price-num {
  font-size: 12px;
  font-weight: 800;
  letter-spacing: -0.2px;
  line-height: 1;
}

/* Buy Action Button */
.dock-buy-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 3px;
  height: 24px;
  padding: 0 9px;
  border-radius: 999px;
  background: #2563eb;
  border: none;
  box-shadow: 0 2px 8px rgba(37, 99, 235, 0.4);
  cursor: pointer;
  transition: all 0.2s ease;
  flex-shrink: 0;
  margin: 0;
  line-height: 1;
}

.dock-buy-btn.is-vip {
  background: linear-gradient(135deg, #f59e0b, #d97706);
  box-shadow: 0 2px 8px rgba(245, 158, 11, 0.4);
}

.dock-buy-btn.is-free {
  background: linear-gradient(135deg, #10b981, #059669);
  box-shadow: 0 2px 8px rgba(16, 185, 129, 0.4);
}

.buy-btn-icon {
  font-size: 11px;
  color: #ffffff;
  line-height: 1;
}

.buy-btn-text {
  font-size: 11px;
  font-weight: 700;
  color: #ffffff;
  white-space: nowrap;
  line-height: 1;
}

/* ========================================================
   Mobile Optimization (<768px)
   ======================================================== */
@media (max-width: 768px) {
  .wallpaper-card {
    border-radius: 16px;
  }

  .card-inner {
    border-radius: 16px;
  }

  .badge-row {
    top: 6px;
    left: 6px;
    right: 6px;
  }

  .badge {
    height: 20px;
    padding: 0 6px;
  }

  .badge-tech-text,
  .badge-tier-text {
    font-size: 9px;
  }

  .glass-dock {
    bottom: 6px;
    left: 6px;
    right: 6px;
    padding: 6px 8px;
    border-radius: 12px;
    gap: 4px;
  }

  .dock-title {
    font-size: 11px;
    line-height: 1.25;
  }

  .dock-cat-pill {
    padding: 1px 4px;
    gap: 2px;
  }

  .cat-dot {
    width: 4px;
    height: 4px;
  }

  .cat-name {
    font-size: 9px;
  }

  .dock-price-pill {
    padding: 1px 5px;
  }

  .price-symbol {
    font-size: 9px;
  }

  .price-num {
    font-size: 11px;
  }

  .dock-buy-btn {
    height: 22px;
    padding: 0 7px;
    gap: 2px;
  }

  .buy-btn-icon {
    font-size: 9px;
  }

  .buy-btn-text {
    font-size: 10px;
  }
}
</style>
