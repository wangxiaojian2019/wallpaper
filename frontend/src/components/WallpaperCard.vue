<template>
  <div class="wallpaper-card" @click="goToDetail">
    <!-- Card Inner (9:19.5 Smartphone Screen Ratio) -->
    <div class="card-cover">
      <!-- Image / Placeholder -->
      <img
        v-if="product.cover_image && !imgError"
        :src="product.cover_image"
        :alt="locale === 'zh' ? product.title_zh : product.title_en"
        loading="lazy"
        class="card-img"
        @error="imgError = true"
      />
      <div
        v-else
        class="cover-placeholder"
        :style="{ background: product._color || 'linear-gradient(135deg,#007AFF,#00C6FF)' }"
      >
        <span class="placeholder-icon">📱</span>
      </div>

      <!-- Top Badges -->
      <div class="top-badges">
        <!-- Pinned Badge -->
        <div v-if="product.is_pinned" class="pinned-badge">
          📌 {{ locale === 'zh' ? '置顶' : 'TOP' }}
        </div>

        <!-- Video Platform Badge (TikTok, 抖音, 小红书, etc.) -->
        <div v-if="videoBadge" class="video-badge" :title="videoBadge.text">
          <span class="v-icon">{{ videoBadge.icon }}</span>
          <span class="v-text">{{ videoBadge.text }}</span>
        </div>

        <!-- Price Tag (Gold glass pill) -->
        <div class="card-price" :class="{ 'is-free': Number(product.price) === 0, 'is-member': auth.isMember }">
          <template v-if="auth.isMember">
            <span class="price-member-icon">👑</span> {{ locale === 'zh' ? '免费' : 'Free' }}
          </template>
          <template v-else-if="Number(product.price) === 0">
            {{ locale === 'zh' ? '免费' : 'Free' }}
          </template>
          <template v-else>
            $ {{ Number(product.price).toFixed(2) }}
          </template>
        </div>
      </div>

      <!-- Hover Overlay Actions -->
      <div class="card-overlay">
        <div class="overlay-actions">
          <button class="overlay-btn view-btn" @click.stop="goToDetail">
            {{ locale === 'zh' ? '预览详情' : 'Preview' }}
          </button>
          <button class="overlay-btn buy-btn" @click.stop="handleBuy">
            <span v-if="auth.isMember">👑 {{ locale === 'zh' ? '解锁' : 'Unlock' }}</span>
            <span v-else-if="Number(product.price) === 0">📥 {{ locale === 'zh' ? '下载' : 'Download' }}</span>
            <span v-else>⚡ {{ t('product.buy') }}</span>
          </button>
        </div>
      </div>

      <!-- Signature iOS Floating Frosted Glass Deck at bottom -->
      <div class="glass-deck">
        <div class="deck-title" :title="locale === 'zh' ? product.title_zh : product.title_en">
          {{ locale === 'zh' ? product.title_zh : (product.title_en || product.title_zh) }}
        </div>
        <div class="deck-meta">
          <span class="meta-cat">
            <span class="cat-dot"></span>
            {{ getCategoryLabel(product.category) }}
          </span>
          <span class="meta-views" v-if="product.views_count">
            👁 {{ product.views_count }}
          </span>
        </div>

        <!-- Visible Action Button on the card -->
        <div class="deck-actions" @click.stop>
          <button
            v-if="auth.isMember"
            class="deck-btn deck-btn-member"
            @click="goToDetail"
            title="VIP会员免单"
          >
            👑 {{ locale === 'zh' ? 'VIP免单下载' : 'VIP Download' }}
          </button>
          <button
            v-else-if="Number(product.price) === 0"
            class="deck-btn deck-btn-free"
            @click="goToDetail"
            title="免费壁纸，点击直接下载"
          >
            📥 {{ locale === 'zh' ? '免费下载' : 'Free Download' }}
          </button>
          <div v-else class="deck-btn-dual">
            <button
              class="deck-btn deck-btn-buy"
              @click="handleQuickBuy('guest')"
              :title="locale === 'zh' ? '临时购买 (免登录游客快捷支付)' : 'Guest Buy'"
            >
              ⚡ ${{ Number(product.price).toFixed(2) }} {{ locale === 'zh' ? '购买' : 'Buy' }}
            </button>
            <button
              class="deck-btn deck-btn-vip"
              @click="handleQuickBuy('member')"
              :title="locale === 'zh' ? '会员登录购买' : 'Member Buy'"
            >
              👑 {{ locale === 'zh' ? '会员' : 'VIP' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/store'

const props = defineProps({
  product: { type: Object, required: true }
})
const emit = defineEmits(['buy'])

const { t, locale } = useI18n()
const router = useRouter()
const auth = useAuthStore()
const imgError = ref(false)

function goToDetail() {
  router.push(`/wallpaper/${props.product.id}`)
}

function handleBuy() {
  if (Number(props.product.price) === 0 || auth.isMember) {
    goToDetail()
    return
  }
  emit('buy', { product: props.product, mode: 'guest' })
}

function handleQuickBuy(mode = 'guest') {
  if (Number(props.product.price) === 0 || auth.isMember) {
    goToDetail()
    return
  }
  emit('buy', { product: props.product, mode })
}

function getCategoryLabel(cat) {
  if (!cat) return locale.value === 'zh' ? '精选壁纸' : 'Wallpaper'
  const map = {
    '苹果壁纸': { zh: '苹果壁纸', en: 'iOS Theme' },
    '动态壁纸': { zh: '动态壁纸', en: 'Live 4K' },
    '电脑壁纸': { zh: '电脑壁纸', en: 'Desktop 5K' },
    '教程分享': { zh: '教程工具', en: 'Tutorial' },
    '动漫': { zh: '动漫', en: 'Anime' },
    'anime': { zh: '动漫', en: 'Anime' },
    'nature': { zh: '自然', en: 'Nature' },
    'space': { zh: '宇宙', en: 'Space' },
    'city': { zh: '城市', en: 'City' },
    'cyberpunk': { zh: '赛博朋克', en: 'Cyberpunk' },
    'minimal': { zh: '极简', en: 'Minimal' }
  }
  const item = map[cat]
  if (item) return locale.value === 'zh' ? item.zh : item.en
  return cat
}

const videoBadge = computed(() => {
  const prod = props.product
  if (!prod.preview_video_url && (!prod.preview_videos || prod.preview_videos === '[]')) return null
  const p = (prod.preview_platform || '').toLowerCase()
  const url = (prod.preview_video_url || '').toLowerCase()
  if (p === 'tiktok' || url.includes('tiktok.com')) return { icon: '🎵', text: 'TikTok' }
  if (p === 'douyin' || url.includes('douyin.com')) return { icon: '🎶', text: '抖音' }
  if (p === 'xiaohongshu' || url.includes('xiaohongshu.com') || url.includes('xhslink.com')) return { icon: '📕', text: '小红书' }
  if (p === 'bilibili' || url.includes('bilibili.com')) return { icon: '📺', text: 'B站' }
  if (p === 'youtube' || url.includes('youtube.com')) return { icon: '▶️', text: 'YouTube' }
  return { icon: '▶', text: locale.value === 'zh' ? '实机演示' : 'Video' }
})
</script>

<style scoped>
.wallpaper-card {
  position: relative;
  width: 100%;
  border-radius: 22px;
  overflow: hidden;
  background: var(--card-bg, #fff);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.06);
  transition: transform 0.3s cubic-bezier(0.16, 1, 0.3, 1), box-shadow 0.3s cubic-bezier(0.16, 1, 0.3, 1);
  cursor: pointer;
  user-select: none;
}

.wallpaper-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 16px 36px rgba(0, 0, 0, 0.15);
}

/* Smartphone screen 9:19.5 aspect ratio */
.card-cover {
  position: relative;
  width: 100%;
  aspect-ratio: 9 / 19.5;
  overflow: hidden;
  background: #111;
}

.card-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s ease;
  will-change: transform;
  display: block;
}

.wallpaper-card:hover .card-img {
  transform: scale(1.06);
}

.cover-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 3rem;
  color: white;
}

/* Top badges bar */
.top-badges {
  position: absolute;
  top: 10px;
  left: 10px;
  right: 10px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  z-index: 5;
  pointer-events: none;
}

/* Pinned badge */
.pinned-badge {
  display: inline-flex;
  align-items: center;
  height: 24px;
  padding: 0 9px;
  background: linear-gradient(135deg, #ef4444, #dc2626);
  border: 1px solid rgba(255, 255, 255, 0.4);
  border-radius: 999px;
  color: #fff;
  font-size: 0.7rem;
  font-weight: 800;
  box-shadow: 0 4px 10px rgba(239, 68, 68, 0.35);
  letter-spacing: 0.3px;
}

/* Video badge */
.video-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  height: 24px;
  padding: 0 9px;
  background: rgba(0, 0, 0, 0.72);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 999px;
  color: #fff;
  font-size: 0.7rem;
  font-weight: 700;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.2);
}

.v-icon { font-size: 0.78rem; }
.v-text { letter-spacing: 0.2px; }

/* Price badge (signature gold frosted glass) */
.card-price {
  margin-left: auto;
  background: rgba(0, 0, 0, 0.68);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 215, 0, 0.35);
  color: #ffd700;
  font-size: 0.75rem;
  font-weight: 800;
  padding: 3px 10px;
  border-radius: 999px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.15);
  letter-spacing: 0.3px;
}

.card-price.is-free {
  border-color: rgba(16, 185, 129, 0.4);
  color: #10b981;
}

.card-price.is-member {
  border-color: rgba(255, 215, 0, 0.6);
  background: rgba(255, 215, 0, 0.2);
  color: #fff;
}

/* Floating iOS Glass Deck at the bottom */
.glass-deck {
  position: absolute;
  bottom: 8px;
  left: 8px;
  right: 8px;
  background: rgba(255, 255, 255, 0.78);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border: 1px solid rgba(255, 255, 255, 0.65);
  border-radius: 16px;
  padding: 10px 12px;
  color: #1a1a1a;
  display: flex;
  flex-direction: column;
  gap: 5px;
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.12);
  z-index: 4;
  pointer-events: auto;
  transition: transform 0.25s ease, background 0.25s ease;
}

[data-theme="dark"] .glass-deck {
  background: rgba(24, 24, 28, 0.82);
  border-color: rgba(255, 255, 255, 0.14);
  color: #f5f5f7;
}

.deck-title {
  font-size: 0.84rem;
  font-weight: 700;
  line-height: 1.35;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.6);
}

[data-theme="dark"] .deck-title {
  text-shadow: 0 1px 3px rgba(0, 0, 0, 0.8);
}

.deck-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.7rem;
  font-weight: 600;
  color: #4b5563;
  margin-top: 1px;
}

[data-theme="dark"] .deck-meta {
  color: #9ca3af;
}

.meta-cat {
  display: flex;
  align-items: center;
  gap: 4px;
}

.cat-dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: #007aff;
}

.meta-views {
  opacity: 0.85;
}

/* Action button bar directly on card deck */
.deck-actions {
  margin-top: 4px;
  pointer-events: auto;
}

.deck-btn-dual {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 6px;
}

.deck-btn {
  width: 100%;
  padding: 6px 10px;
  border-radius: 10px;
  font-size: 0.75rem;
  font-weight: 700;
  border: none;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  transition: all 0.2s ease;
  line-height: 1.2;
}

.deck-btn-buy {
  background: linear-gradient(135deg, #007aff, #00c6ff);
  color: #fff;
  box-shadow: 0 2px 8px rgba(0, 122, 255, 0.35);
}

.deck-btn-buy:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 122, 255, 0.45);
}

.deck-btn-vip {
  background: rgba(245, 158, 11, 0.15);
  color: #d97706;
  border: 1px solid rgba(245, 158, 11, 0.4);
  padding: 6px 8px;
}

[data-theme="dark"] .deck-btn-vip {
  background: rgba(245, 158, 11, 0.25);
  color: #fbbf24;
  border-color: rgba(245, 158, 11, 0.5);
}

.deck-btn-vip:hover {
  background: rgba(245, 158, 11, 0.3);
  transform: translateY(-1px);
}

.deck-btn-free {
  background: linear-gradient(135deg, #10b981, #059669);
  color: #fff;
  box-shadow: 0 2px 8px rgba(16, 185, 129, 0.35);
}

.deck-btn-free:hover {
  transform: translateY(-1px);
}

.deck-btn-member {
  background: linear-gradient(135deg, #f59e0b, #d97706);
  color: #fff;
  box-shadow: 0 2px 8px rgba(245, 158, 11, 0.35);
}

/* Hover overlay with quick action buttons */
.card-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(180deg, rgba(0,0,0,0.1) 0%, rgba(0,0,0,0.7) 100%);
  opacity: 0;
  transition: opacity 0.25s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  z-index: 6;
}

.wallpaper-card:hover .card-overlay {
  opacity: 1;
}

.overlay-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
  width: 85%;
  transform: translateY(10px);
  transition: transform 0.25s ease;
}

.wallpaper-card:hover .overlay-actions {
  transform: translateY(0);
}

.overlay-btn {
  padding: 10px 14px;
  border-radius: 12px;
  font-size: 0.82rem;
  font-weight: 700;
  cursor: pointer;
  border: none;
  transition: all 0.2s;
  text-align: center;
}

.view-btn {
  background: rgba(255, 255, 255, 0.9);
  color: #111;
  backdrop-filter: blur(8px);
}

.view-btn:hover {
  background: #fff;
  transform: scale(1.02);
}

.buy-btn {
  background: linear-gradient(135deg, #007aff, #00c6ff);
  color: white;
  box-shadow: 0 4px 14px rgba(0, 122, 255, 0.4);
}

.buy-btn:hover {
  opacity: 0.92;
  transform: scale(1.02);
}

/* Mobile Responsive Optimization (2-column phone screens) */
@media (max-width: 640px) {
  .wallpaper-card {
    border-radius: 16px;
  }

  .card-overlay {
    display: none !important; /* Disable hover overlay on mobile touch screens to prevent double-tap issues */
  }

  .top-badges {
    top: 6px;
    left: 6px;
    right: 6px;
  }

  .video-badge {
    height: 22px;
    padding: 0 6px;
    font-size: 0.65rem;
  }

  .card-price {
    font-size: 0.7rem;
    padding: 2px 7px;
  }

  .glass-deck {
    bottom: 6px;
    left: 6px;
    right: 6px;
    padding: 8px 8px;
    border-radius: 12px;
    gap: 4px;
  }

  .deck-title {
    font-size: 0.78rem;
    line-height: 1.25;
    -webkit-line-clamp: 1;
  }

  .deck-meta {
    font-size: 0.65rem;
  }

  .deck-btn-dual {
    gap: 4px;
    grid-template-columns: 1fr auto;
  }

  .deck-btn {
    padding: 6px 6px !important;
    font-size: 0.68rem !important;
    border-radius: 8px !important;
    white-space: nowrap !important;
  }

  .deck-btn-vip {
    padding: 6px 6px !important;
    font-size: 0.68rem !important;
  }
}
</style>
