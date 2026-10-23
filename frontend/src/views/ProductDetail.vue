<template>
  <div class="detail-page">
    <div v-if="loading" class="detail-loading">
      <span class="spinner" style="width:40px;height:40px;border-width:3px"></span>
    </div>

    <template v-else-if="product">
      <!-- Back to Home -->
      <div class="container detail-nav">
        <button class="btn btn-ghost btn-sm" @click="$router.push('/')">← {{ t('common.back') }}</button>
      </div>

      <div class="container detail-grid">
        <!-- Left: Image -->
        <!-- Left: Image & Embedded Video Player -->
        <div class="detail-image-wrap">
          <!-- Media Mode Tabs Switcher if TikTok is available -->
          <div v-if="tiktokInfo" class="media-tabs-bar">
            <button
              class="media-tab-btn"
              :class="{ active: activeMediaView === 'cover' }"
              @click="activeMediaView = 'cover'"
            >
              🖼️ {{ locale === 'zh' ? '静态壁纸封面' : 'Cover Image' }}
            </button>
            <button
              class="media-tab-btn live-tab-btn"
              :class="{ active: activeMediaView === 'tiktok' }"
              @click="activeMediaView = 'tiktok'"
            >
              <span class="pulse-live-dot"></span>
              🎬 {{ locale === 'zh' ? 'TikTok 实机动效 (内嵌播放)' : 'TikTok Live Video' }}
            </button>
          </div>

          <!-- View 1: Static Image Cover -->
          <div v-show="!tiktokInfo || activeMediaView === 'cover'" class="detail-image">
            <img v-if="product.cover_image" :src="product.cover_image" :alt="product.title_zh" />
            <div v-else class="image-placeholder" :style="{ background: product._color || 'linear-gradient(135deg,#667eea,#764ba2)' }">
              <span>🖼️</span>
            </div>
            <!-- Quick overlay button to switch to TikTok live play -->
            <button
              v-if="tiktokInfo"
              class="overlay-play-btn"
              @click="activeMediaView = 'tiktok'"
              :title="locale === 'zh' ? '点击切换为 TikTok 实机内嵌动效播放' : 'Play TikTok live demo'"
            >
              <span class="play-icon">▶</span>
              <span>{{ locale === 'zh' ? '点击内嵌播放动态实机效果' : 'Play Live Motion Demo' }}</span>
            </button>
          </div>

          <!-- View 2: Embedded TikTok Player (iPhone Smartphone Frame) -->
          <div v-if="tiktokInfo" v-show="activeMediaView === 'tiktok'" class="tiktok-player-card">
            <div class="tiktok-phone-bezel">
              <div class="phone-dynamic-island"></div>
              <iframe
                :src="`https://www.tiktok.com/player/v1/${tiktokInfo.id}?music_info=1&description=1`"
                class="tiktok-embed-frame"
                allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
                allowfullscreen
                loading="lazy"
              ></iframe>
            </div>
            <div class="tiktok-control-strip">
              <div class="control-strip-info">
                <span class="tiktok-badge">🎵 TikTok</span>
                <span class="control-strip-title">{{ locale === 'zh' ? '实机动态锁屏演示' : 'Live Lock Screen Motion' }}</span>
              </div>
              <div class="control-strip-actions">
                <button class="btn btn-ghost btn-xs" @click="activeMediaView = 'cover'">
                  🖼️ {{ locale === 'zh' ? '返回封面' : 'Cover' }}
                </button>
                <a :href="tiktokInfo.url" target="_blank" rel="noopener noreferrer" class="btn btn-outline btn-xs">
                  ↗ {{ locale === 'zh' ? '原平台打开' : 'TikTok App' }}
                </a>
              </div>
            </div>
          </div>

          <!-- Other Video Preview Platforms (TikTok, 抖音, 小红书, B站, etc.) -->
          <div v-if="videoLinks.length > 0" class="video-preview-card">
            <div class="video-preview-header">
              <span class="video-icon-anim">🎬</span>
              <span class="video-title-text">{{ locale === 'zh' ? '动态效果演示链接' : 'Live Wallpaper Previews' }}</span>
            </div>
            <div class="video-platforms-list">
              <template v-for="(v, idx) in videoLinks" :key="idx">
                <!-- If this is TikTok: offer an In-Page Play button + External Link -->
                <div v-if="v.platform === 'tiktok' || (v.url && v.url.includes('tiktok.com'))" class="video-platform-dual">
                  <button
                    class="video-platform-link platform-tiktok inpage-play-btn"
                    @click="activeMediaView = 'tiktok'"
                  >
                    <div class="v-platform-info">
                      <span class="v-platform-icon">🎬</span>
                      <span class="v-platform-name">{{ locale === 'zh' ? 'TikTok 内嵌直接播放' : 'Play TikTok Embedded' }}</span>
                    </div>
                    <span class="v-jump-badge live-badge">▶ {{ locale === 'zh' ? '立即播放' : 'Play' }}</span>
                  </button>
                  <a
                    :href="v.url"
                    target="_blank"
                    rel="noopener noreferrer"
                    class="video-platform-link-mini"
                    title="在 TikTok 网页/App 打开"
                  >
                    ↗
                  </a>
                </div>

                <!-- Other platforms (Douyin, Xiaohongshu, Bilibili, etc.) -->
                <a
                  v-else
                  :href="v.url"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="video-platform-link"
                  :class="`platform-${v.platform || 'direct'}`"
                >
                  <div class="v-platform-info">
                    <span class="v-platform-icon">{{ getPlatformMeta(v.platform).icon }}</span>
                    <span class="v-platform-name">{{ v.title || getPlatformMeta(v.platform).name }}</span>
                  </div>
                  <span class="v-jump-badge">{{ locale === 'zh' ? '前往查看 ↗' : 'Watch ↗' }}</span>
                </a>
              </template>
            </div>
          </div>
        </div>

        <!-- Right: Info -->
        <div class="detail-info">
          <!-- Category & Tags -->
          <div class="detail-meta">
            <span v-if="product.category" class="badge badge-neutral">
              {{ getCategoryLabel(product.category) }}
            </span>
            <span v-for="tag in parseTags(product.tags)" :key="tag" class="tag" style="font-size:0.72rem; padding:2px 10px;">
              {{ tag }}
            </span>
          </div>

          <h1 class="detail-title">{{ locale === 'zh' ? product.title_zh : product.title_en }}</h1>

          <!-- Stats -->
          <div class="detail-stats">
            <span>👁 {{ product.views || 0 }} {{ t('product.views') }}</span>
            <span>🛒 {{ product.sales || 0 }} {{ t('product.sales') }}</span>
          </div>

          <!-- Price block -->
          <div class="price-block">
            <template v-if="auth.isMember">
              <div class="member-price">
                <span class="badge badge-gold">👑 {{ locale === 'zh' ? '会员免费' : 'Member Free' }}</span>
              </div>
            </template>
            <template v-else>
              <div class="single-price">
                <span class="price-label">{{ t('product.single_price') }}</span>
                <div class="price-display">
                  <span class="price-currency">$</span>
                  <span class="price-amount">{{ product.price }}</span>
                  <span class="price-unit">USD</span>
                </div>
              </div>
            </template>
          </div>

          <!-- Action buttons -->
          <div class="detail-actions">
            <!-- Case 1: Free Wallpaper (price === 0) -->
            <template v-if="Number(product.price) === 0">
              <!-- Subcase 1A: Free & Unlocked -->
              <div v-if="isPurchased" class="unlocked-section free-unlocked">
                <div class="unlocked-header free-header">
                  <span>🎁</span>
                  <span>{{ locale === 'zh' ? '免费专享壁纸，无需支付，直接选择网盘下载：' : 'Free Wallpaper! Select Cloud Drive Below:' }}</span>
                </div>
                <DownloadLinks :product-id="product.id" />
              </div>

              <!-- Subcase 1B: Free requires login -->
              <div v-else-if="requireLoginForFree" class="unlocked-section free-prompt-box">
                <div class="free-prompt-inner">
                  <div class="free-prompt-icon">🔒</div>
                  <div class="free-prompt-body">
                    <h4>{{ locale === 'zh' ? '免费专享壁纸 · 需登录后下载' : 'Free Wallpaper - Login Required' }}</h4>
                    <p>{{ locale === 'zh' ? '本款壁纸为注册会员免费专享福利，登录已有账号或极速注册后即可免费获取百度/夸克/谷歌网盘下载地址与提取码！' : 'Please log in or register to get cloud drive download links for free.' }}</p>
                  </div>
                  <button class="btn btn-primary btn-lg" @click="openBuyWithMode('member')">
                    🔑 {{ locale === 'zh' ? '立即登录 / 极速注册免费下载' : 'Login / Register to Download' }}
                  </button>
                </div>
              </div>

              <!-- Subcase 1C: Daily limit reached -->
              <div v-else-if="dailyLimitReached" class="unlocked-section free-prompt-box limit-box">
                <div class="free-prompt-inner">
                  <div class="free-prompt-icon">⏳</div>
                  <div class="free-prompt-body">
                    <h4>{{ locale === 'zh' ? '今日免费下载次数已达上限' : 'Daily Free Download Limit Reached' }}</h4>
                    <p>{{ locale === 'zh' ? `普通注册用户每日限免下载 ${dailyFreeLimit} 次，请明日再来；或开通 VIP 会员立即畅享全场无限次极速免单下载！` : 'Daily limit reached. Upgrade to VIP for unlimited downloads.' }}</p>
                  </div>
                  <router-link to="/membership" class="btn btn-gold btn-lg">
                    👑 {{ locale === 'zh' ? '升级 VIP 尊享全场无限免单 →' : 'Upgrade VIP Unlimited →' }}
                  </router-link>
                </div>
              </div>

              <!-- Subcase 1D: Fallback direct downloads -->
              <div v-else class="unlocked-section free-unlocked">
                <div class="unlocked-header free-header">
                  <span>🎁</span>
                  <span>{{ locale === 'zh' ? '免费专享壁纸，无需支付，直接选择网盘下载：' : 'Free Wallpaper! Select Cloud Drive Below:' }}</span>
                </div>
                <DownloadLinks :product-id="product.id" />
              </div>
            </template>

            <!-- Case 2: VIP Member (free access) -->
            <template v-else-if="auth.isMember">
              <div class="unlocked-section vip-unlocked">
                <div class="unlocked-header vip-header">
                  <span>👑</span>
                  <span>{{ locale === 'zh' ? '尊享 VIP 会员特权：全场壁纸免单已生效，直接下载：' : 'VIP Member Benefit: Free unlock for you!' }}</span>
                </div>
                <DownloadLinks :product-id="product.id" />
                <div class="reopen-test-bar">
                  <span class="reopen-hint">{{ locale === 'zh' ? '💡 想要测试其他支付通道？' : 'Want to test checkout?' }}</span>
                  <button class="btn btn-ghost btn-xs" @click="openBuyWithMode(guestBuyEnabled ? 'guest' : 'member')">
                    🛒 {{ locale === 'zh' ? '打开收银台弹窗' : 'Open Cashier' }}
                  </button>
                </div>
              </div>
            </template>

            <!-- Case 3: Paid & Already Purchased -->
            <template v-else-if="isPurchased">
              <div class="unlocked-section paid-unlocked">
                <div class="unlocked-header paid-header">
                  <span>🎉</span>
                  <span>{{ locale === 'zh' ? '已成功解锁本款壁纸！请选择下方网盘下载：' : 'Unlocked! Select Cloud Drive Below:' }}</span>
                </div>
                <DownloadLinks :product-id="product.id" />
                <div class="reopen-test-bar">
                  <span class="reopen-hint">{{ locale === 'zh' ? '💡 想要体验其他购买方式？' : 'Want to test other checkout flow?' }}</span>
                  <button class="btn btn-outline btn-xs" @click="openBuyWithMode(guestBuyEnabled ? 'guest' : 'member')">
                    🛒 {{ locale === 'zh' ? '再次打开购买/支付收银台' : 'Re-open Buy Cashier' }}
                  </button>
                </div>
              </div>
            </template>

            <!-- Case 4: Paid & Not Yet Purchased: Dual Channels (临时购买 vs 会员登录购买) -->
            <template v-else>
              <div class="purchase-channels-grid" :class="{ 'single-channel': !guestBuyEnabled }">
                <!-- Channel 1: 临时购买 (免登录游客快捷下单) -->
                <div v-if="guestBuyEnabled" class="channel-card guest-channel-card">
                  <div class="channel-card-badge guest-badge">
                    ⚡ {{ locale === 'zh' ? '无需注册・极速通道' : 'Instant Guest' }}
                  </div>
                  <div class="channel-card-body">
                    <div class="channel-card-title-row">
                      <h4>{{ locale === 'zh' ? '临时购买 (免登录)' : 'Guest Checkout' }}</h4>
                      <span class="channel-card-price">${{ product.price }}</span>
                    </div>
                    <p class="channel-card-desc">
                      {{ locale === 'zh' ? '游客直接付款，无需注册账号，1秒快速解锁百度/夸克/谷歌网盘下载' : 'No account needed. Instant unlock with PayPal or test payment.' }}
                    </p>
                    <button class="btn btn-primary btn-block btn-lg channel-action-btn" @click="openBuyWithMode('guest')">
                      ⚡ {{ locale === 'zh' ? `立即临时购买 ($${product.price})` : `Guest Buy ($${product.price})` }}
                    </button>
                  </div>
                </div>

                <!-- Channel 2: 会员登陆购买 -->
                <div class="channel-card member-channel-card" :class="{ 'full-width': !guestBuyEnabled }">
                  <div class="channel-card-badge member-badge">
                    👑 {{ locale === 'zh' ? '专属特权・订单保存' : 'Member / VIP' }}
                  </div>
                  <div class="channel-card-body">
                    <div class="channel-card-title-row">
                      <h4>{{ locale === 'zh' ? '会员登录购买' : 'Member Checkout' }}</h4>
                      <span v-if="auth.isLoggedIn" class="channel-user-pill">{{ auth.user?.email }}</span>
                    </div>

                    <!-- Not Logged In -->
                    <template v-if="!auth.isLoggedIn">
                      <p class="channel-card-desc">
                        {{ locale === 'zh' ? '登录已有账号购买永久保存订单；或直接开通 VIP 会员全场所有壁纸免费下！' : 'Login to save order in account, or subscribe VIP for unlimited downloads.' }}
                      </p>
                      <div class="channel-dual-btns">
                        <button class="btn btn-dark btn-block channel-action-btn" @click="openBuyWithMode('member')">
                          🔑 {{ locale === 'zh' ? '会员登录购买' : 'Login & Buy' }}
                        </button>
                        <router-link to="/membership" class="btn btn-gold btn-block channel-action-btn">
                          👑 {{ locale === 'zh' ? '开通 VIP 全场免单' : 'VIP Unlimited ($9.99)' }}
                        </router-link>
                      </div>
                    </template>

                    <!-- Logged in as regular user -->
                    <template v-else>
                      <p class="channel-card-desc">
                        {{ locale === 'zh' ? '以当前账号购买此壁纸并写入个人订单中心，跨设备随时可反复下载' : 'Purchase with account to save order permanently.' }}
                      </p>
                      <div class="channel-dual-btns">
                        <button class="btn btn-dark btn-block channel-action-btn" @click="openBuyWithMode('member')">
                          💳 {{ locale === 'zh' ? `账号购买 ($${product.price})` : `Buy on Account` }}
                        </button>
                        <router-link to="/membership" class="btn btn-gold btn-block channel-action-btn">
                          👑 {{ locale === 'zh' ? '升级 VIP 全场免单' : 'Upgrade VIP' }}
                        </router-link>
                      </div>
                    </template>
                  </div>
                </div>
              </div>

              <!-- Extra redeem entry -->
              <div class="redeem-entry-bar">
                <button class="btn-redeem-link" @click="openBuyWithMode('code')">
                  🎫 {{ locale === 'zh' ? '使用提取码 / 卡密直接兑换解锁 →' : 'Redeem with gift code or voucher →' }}
                </button>
              </div>
            </template>
          </div>

          <!-- Description -->
          <div class="detail-desc">
            <h3>{{ locale === 'zh' ? '商品描述' : 'Description' }}</h3>
            <div class="desc-content" v-html="product.description || (locale === 'zh' ? '暂无描述' : 'No description')"></div>
          </div>

          <!-- Universal Setup Tutorial Template Section -->
          <div
            v-if="tutorialConfig && tutorialConfig.enabled !== false && tutorialConfig.items && tutorialConfig.items.length > 0"
            class="universal-tutorial-box card"
          >
            <div class="tut-header-row">
              <div class="tut-header-left">
                <span class="tut-header-icon">📖</span>
                <h3 class="tut-header-title">
                  {{ locale === 'zh' ? (tutorialConfig.title_zh || 'iOS 动态壁纸实机交互设置教程') : (tutorialConfig.title_en || 'iOS Live Wallpaper Setup Guide') }}
                </h3>
              </div>
              <span class="tut-header-badge">{{ locale === 'zh' ? '通用设置教程' : 'Setup Guide' }}</span>
            </div>

            <div v-if="tutorialNotice" class="tut-notice-bar">
              <span>{{ tutorialNotice }}</span>
            </div>

            <!-- Version Switcher Pills -->
            <div class="tut-versions-row">
              <button
                v-for="(item, idx) in tutorialConfig.items"
                :key="item.id || idx"
                class="tut-version-pill"
                :class="{
                  active: activeTutorialIdx === idx,
                  [item.badge_color || 'blue']: true
                }"
                @click="activeTutorialIdx = idx"
              >
                <span>{{ item.version_range }}</span>
                <span v-if="item.badge_zh" class="pill-sub-tag">
                  {{ locale === 'zh' ? item.badge_zh : (item.badge_en || item.badge_zh) }}
                </span>
              </button>
            </div>

            <!-- Active Version Content Card -->
            <div v-if="activeTutorialItem" class="tut-active-card">
              <div class="tut-active-head">
                <h4 class="tut-subtitle">
                  {{ locale === 'zh' ? activeTutorialItem.title_zh : (activeTutorialItem.title_en || activeTutorialItem.title_zh) }}
                </h4>
                <a
                  v-if="activeTutorialItem.video_url"
                  :href="activeTutorialItem.video_url"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="btn btn-primary btn-sm tut-video-btn"
                >
                  ▶ {{ locale === 'zh' ? '观看视频教程' : 'Watch Video' }} ↗
                </a>
              </div>

              <!-- Steps List -->
              <div v-if="activeTutorialItem.steps && activeTutorialItem.steps.length > 0" class="tut-steps-list">
                <div
                  v-for="(step, sIdx) in activeTutorialItem.steps"
                  :key="sIdx"
                  class="tut-step-row"
                >
                  <span class="tut-step-num">{{ sIdx + 1 }}</span>
                  <span class="tut-step-text">{{ step }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Related -->
      <section v-if="related.length" class="related-section">
        <div class="container">
          <h2>{{ t('product.related') }}</h2>
          <div class="related-grid">
            <WallpaperCard v-for="p in related" :key="p.id" :product="p" @buy="openBuy" />
          </div>
        </div>
      </section>
    </template>

    <div v-else class="not-found">
      <span>🔍</span>
      <p>{{ locale === 'zh' ? '壁纸不存在' : 'Wallpaper not found' }}</p>
      <router-link to="/" class="btn btn-primary">{{ t('common.back') }}</router-link>
    </div>

    <!-- Mobile Sticky Bottom Action Bar -->
    <div v-if="product" class="mobile-sticky-bar show-mobile-bar">
      <div class="sticky-bar-content">
        <!-- Left: Price / Status -->
        <div class="sticky-bar-info">
          <template v-if="auth.isMember">
            <span class="sticky-status-badge badge-gold">👑 {{ locale === 'zh' ? 'VIP免单' : 'VIP Free' }}</span>
          </template>
          <template v-else-if="Number(product.price) === 0">
            <span class="sticky-status-badge badge-success">🎁 {{ locale === 'zh' ? '免费壁纸' : 'Free' }}</span>
          </template>
          <template v-else-if="isPurchased">
            <span class="sticky-status-badge badge-success">🎉 {{ locale === 'zh' ? '已解锁' : 'Unlocked' }}</span>
          </template>
          <template v-else>
            <div class="sticky-price-row">
              <span class="sticky-price-currency">$</span>
              <span class="sticky-price-num">{{ product.price }}</span>
              <span class="sticky-price-unit">USD</span>
            </div>
          </template>
        </div>

        <!-- Right: Action Buttons -->
        <div class="sticky-bar-actions">
          <template v-if="auth.isMember || isPurchased">
            <button class="btn btn-primary sticky-btn" @click="scrollToDownloads">
              📥 {{ locale === 'zh' ? '查看网盘下载' : 'Downloads' }}
            </button>
          </template>
          <template v-else-if="Number(product.price) === 0">
            <button v-if="requireLoginForFree" class="btn btn-primary sticky-btn" @click="openBuyWithMode('member')">
              🔑 {{ locale === 'zh' ? '登录免费下载' : 'Login' }}
            </button>
            <button v-else-if="dailyLimitReached" class="btn btn-gold sticky-btn" @click="router.push('/membership')">
              👑 {{ locale === 'zh' ? '升级VIP' : 'VIP' }}
            </button>
            <button v-else class="btn btn-primary sticky-btn" @click="scrollToDownloads">
              📥 {{ locale === 'zh' ? '免费下载' : 'Download' }}
            </button>
          </template>
          <template v-else>
            <button v-if="guestBuyEnabled" class="btn btn-primary sticky-btn" @click="openBuyWithMode('guest')">
              ⚡ {{ locale === 'zh' ? '临时购买' : 'Guest' }}
            </button>
            <button class="btn btn-gold sticky-btn" @click="openBuyWithMode('member')">
              👑 {{ locale === 'zh' ? '会员购买' : 'VIP' }}
            </button>
          </template>
        </div>
      </div>
    </div>

    <!-- Buy modal -->
    <BuyModal
      v-if="showBuyModal"
      :product="product"
      :initial-mode="buyModalMode"
      @close="showBuyModal = false"
      @success="handleSuccess"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/store'
import { productAPI, orderAPI, configAPI, tutorialAPI } from '@/api'
import { useToast } from '@/composables/useToast'
import WallpaperCard from '@/components/WallpaperCard.vue'
import BuyModal from '@/components/BuyModal.vue'
import DownloadLinks from '@/components/DownloadLinks.vue'

const route = useRoute()
const router = useRouter()
const { t, locale } = useI18n()
const auth = useAuthStore()
const { success } = useToast()

const tutorialConfig = ref(null)
const activeTutorialIdx = ref(0)
const activeTutorialItem = computed(() => {
  if (!tutorialConfig.value?.items || tutorialConfig.value.items.length === 0) return null
  return tutorialConfig.value.items[activeTutorialIdx.value] || tutorialConfig.value.items[0]
})
const tutorialNotice = computed(() => {
  if (!tutorialConfig.value) return ''
  return locale.value === 'zh'
    ? (tutorialConfig.value.notice_zh || '')
    : (tutorialConfig.value.notice_en || tutorialConfig.value.notice_zh || '')
})

async function loadTutorialConfig() {
  try {
    const res = await tutorialAPI.getPublic()
    if (res && res.data) {
      tutorialConfig.value = res.data
    } else if (res) {
      tutorialConfig.value = res
    }
  } catch (e) {
    console.warn('Failed to load tutorial template', e)
  }
}

const product = ref(null)
const related = ref([])
const loading = ref(true)
const showBuyModal = ref(false)
const buyModalMode = ref('guest')
const isPurchased = ref(false)
const activeMediaView = ref('cover') // 'cover' | 'tiktok'
const requireLoginForFree = ref(false)
const dailyLimitReached = ref(false)
const dailyFreeLimit = ref(0)
const guestBuyEnabled = ref(true)

const platformMetaMap = {
  douyin: { name: '抖音 演示视频', icon: '🎶' },
  xiaohongshu: { name: '小红书 效果笔记', icon: '📕' },
  tiktok: { name: 'TikTok Preview', icon: '🎵' },
  bilibili: { name: '哔哩哔哩 高清演示', icon: '📺' },
  youtube: { name: 'YouTube 4K Preview', icon: '▶️' },
  direct: { name: '实机动效演示', icon: '🎬' }
}

function getPlatformMeta(platform) {
  return platformMetaMap[platform] || platformMetaMap.direct
}

const videoLinks = computed(() => {
  if (!product.value) return []
  const list = []
  if (product.value.preview_videos) {
    try {
      const parsed = typeof product.value.preview_videos === 'string'
        ? JSON.parse(product.value.preview_videos)
        : product.value.preview_videos
      if (Array.isArray(parsed) && parsed.length > 0) {
        return parsed.filter(v => v.url)
      }
    } catch { /* ignore */ }
  }
  if (product.value.preview_video_url) {
    let platform = product.value.preview_platform || 'douyin'
    const lower = product.value.preview_video_url.toLowerCase()
    if (lower.includes('tiktok.com')) platform = 'tiktok'
    else if (lower.includes('douyin.com') || lower.includes('v.douyin.com')) platform = 'douyin'
    else if (lower.includes('xiaohongshu.com') || lower.includes('xhslink.com')) platform = 'xiaohongshu'
    else if (lower.includes('bilibili.com') || lower.includes('b23.tv')) platform = 'bilibili'
    else if (lower.includes('youtube.com') || lower.includes('youtu.be')) platform = 'youtube'

    list.push({
      platform,
      url: product.value.preview_video_url,
      title: ''
    })
  }
  return list
})

function extractTikTokVideoId(urlOrHtml) {
  if (!urlOrHtml) return null
  const match = urlOrHtml.match(/data-video-id=["']?(\d+)["']?/i) ||
                urlOrHtml.match(/\/video\/(\d+)/i) || 
                urlOrHtml.match(/\/v\/(\d+)/i) || 
                urlOrHtml.match(/\/embed\/(?:v\d+\/)?(\d+)/i) ||
                urlOrHtml.match(/\/player\/v\d+\/(\d+)/i)
  return match ? match[1] : null
}

const tiktokInfo = computed(() => {
  if (!product.value) return null

  // 1. Check videoLinks
  for (const v of videoLinks.value) {
    if (v.platform === 'tiktok' || (v.url && v.url.toLowerCase().includes('tiktok.com'))) {
      const id = extractTikTokVideoId(v.url)
      if (id) {
        return { id, url: v.url, title: v.title || 'TikTok 演示视频' }
      }
    }
  }

  // 2. Check preview_video_url
  if (product.value.preview_video_url && product.value.preview_video_url.toLowerCase().includes('tiktok.com')) {
    const id = extractTikTokVideoId(product.value.preview_video_url)
    if (id) {
      return { id, url: product.value.preview_video_url, title: 'TikTok 演示视频' }
    }
  }

  // 3. Check description
  if (product.value.description && product.value.description.includes('tiktok')) {
    const id = extractTikTokVideoId(product.value.description)
    if (id) {
      const citeMatch = product.value.description.match(/cite=["']([^"']+)["']/i)
      const url = citeMatch ? citeMatch[1] : `https://www.tiktok.com/@user/video/${id}`
      return { id, url, title: 'TikTok 演示视频' }
    }
  }

  return null
})

function loadTikTokScript() {
  if (typeof window === 'undefined') return
  const existing = document.getElementById('tiktok-embed-script')
  if (existing) existing.remove()
  const script = document.createElement('script')
  script.id = 'tiktok-embed-script'
  script.src = 'https://www.tiktok.com/embed.js'
  script.async = true
  document.body.appendChild(script)
}

const categoryMap = {
  nature: { zh: '自然', en: 'Nature' },
  space: { zh: '宇宙', en: 'Space' },
  city: { zh: '城市', en: 'City' },
  abstract: { zh: '抽象', en: 'Abstract' },
  anime: { zh: '动漫', en: 'Anime' },
  minimal: { zh: '极简', en: 'Minimal' }
}

function getCategoryLabel(cat) {
  const item = categoryMap[cat]
  return item ? (locale.value === 'zh' ? item.zh : item.en) : cat
}

function parseTags(tags) {
  if (!tags) return []
  if (Array.isArray(tags)) return tags
  try { return JSON.parse(tags) } catch { return [] }
}

async function loadProduct() {
  loading.value = true
  try {
    try {
      const cfg = await configAPI.getPublic()
      if (cfg?.permissions) {
        guestBuyEnabled.value = cfg.permissions.guest_buy_enabled !== false
      }
    } catch {}

    const res = await productAPI.get(route.params.id)
    product.value = res.product || res
    nextTick(() => {
      loadTikTokScript()
    })

    requireLoginForFree.value = !!res.require_login_for_free
    dailyLimitReached.value = !!res.daily_limit_reached
    dailyFreeLimit.value = Number(res.daily_free_limit) || 0

    if (res.is_unlocked) {
      isPurchased.value = true
    } else if (Number(product.value?.price || 0) <= 0) {
      isPurchased.value = false
    } else if (auth.isLoggedIn) {
      const orders = await orderAPI.list()
      const list = orders.orders || orders || []
      isPurchased.value = list.some(o => o.product_id === product.value.id && o.status === 'paid')
    }
    // Load related
    const relRes = await productAPI.list({ category: product.value.category, page_size: 6, limit: 6 })
    const relItems = relRes.items || relRes.products || (Array.isArray(relRes) ? relRes : [])
    related.value = relItems.filter(p => p.id !== product.value.id).slice(0, 4)
  } catch {
    // Mock data fallback
    product.value = {
      id: route.params.id,
      title_zh: '精美壁纸示例',
      title_en: 'Premium Wallpaper Example',
      description: '<p>这是一张精美的交互壁纸，支持高清下载。</p>',
      price: '9.99',
      category: 'nature',
      tags: ['4K', '自然', '高清'],
      cover_image: null,
      preview_video_url: null,
      views: 1234,
      sales: 56,
      _color: 'linear-gradient(135deg,#43e97b,#38f9d7)'
    }
  } finally {
    loading.value = false
  }
}

function openBuyWithMode(mode = 'guest') {
  buyModalMode.value = mode
  showBuyModal.value = true
}

function scrollToDownloads() {
  const el = document.querySelector('.unlocked-section') || document.querySelector('.detail-actions')
  if (el) {
    el.scrollIntoView({ behavior: 'smooth', block: 'center' })
  }
}

function openBuy(payload) {
  if (payload?.product) {
    product.value = payload.product
    buyModalMode.value = payload.mode || 'guest'
  } else {
    product.value = payload
    buyModalMode.value = 'guest'
  }
  showBuyModal.value = true
}

function handleSuccess() {
  showBuyModal.value = false
  isPurchased.value = true
  success(t('payment.success'))
}

onMounted(() => {
  loadProduct()
  loadTutorialConfig()
})
</script>

<style scoped>
.detail-page { min-height: 100vh; padding-bottom: 80px; }

.detail-loading {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 60vh;
}

.detail-nav {
  padding-top: 32px;
  padding-bottom: 16px;
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr 1.2fr;
  gap: 60px;
  align-items: start;
  padding-top: 20px;
}

@media (max-width: 900px) {
  .detail-grid { grid-template-columns: 1fr; gap: 32px; }
  .detail-image-wrap { position: static; }
}

/* Image (iPhone 9:19.5 Mockup Frame) */
.detail-image-wrap {
  position: sticky;
  top: calc(var(--nav-h) + 20px);
  display: flex;
  flex-direction: column;
  gap: 16px;
  max-width: 380px;
  margin: 0 auto;
  width: 100%;
}

/* Media mode tabs */
.media-tabs-bar {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 6px;
  background: var(--c-bg, #f5f5f7);
  padding: 4px;
  border-radius: 14px;
  border: 1px solid var(--c-border, #e2e8f0);
}

[data-theme="dark"] .media-tabs-bar {
  background: rgba(255, 255, 255, 0.06);
  border-color: rgba(255, 255, 255, 0.1);
}

.media-tab-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 8px 10px;
  border-radius: 10px;
  border: none;
  background: transparent;
  font-size: 0.78rem;
  font-weight: 700;
  color: #666;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  white-space: nowrap;
}

[data-theme="dark"] .media-tab-btn {
  color: #aaa;
}

.media-tab-btn.active {
  background: #ffffff;
  color: #111;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

[data-theme="dark"] .media-tab-btn.active {
  background: #242428;
  color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.pulse-live-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #ff2d55;
  box-shadow: 0 0 0 0 rgba(255, 45, 85, 0.7);
  animation: pulseDot 1.6s infinite;
}

@keyframes pulseDot {
  0% { box-shadow: 0 0 0 0 rgba(255, 45, 85, 0.7); }
  70% { box-shadow: 0 0 0 6px rgba(255, 45, 85, 0); }
  100% { box-shadow: 0 0 0 0 rgba(255, 45, 85, 0); }
}

.detail-image {
  border-radius: 28px;
  overflow: hidden;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.2), 0 0 0 1px rgba(255, 255, 255, 0.12);
  aspect-ratio: 9 / 19.5;
  width: 100%;
  background: #111;
  position: relative;
}

.detail-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.image-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 4rem;
}

/* Overlay play button on cover */
.overlay-play-btn {
  position: absolute;
  bottom: 16px;
  left: 16px;
  right: 16px;
  padding: 10px 14px;
  border-radius: 14px;
  background: rgba(0, 0, 0, 0.78);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.25);
  color: #fff;
  font-size: 0.82rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  cursor: pointer;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.35);
  transition: all 0.2s ease;
  z-index: 5;
}

.overlay-play-btn:hover {
  background: rgba(0, 0, 0, 0.9);
  transform: translateY(-2px);
  box-shadow: 0 12px 28px rgba(0, 0, 0, 0.45);
}

.play-icon {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: #ff2d55;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 0.65rem;
  padding-left: 2px;
}

/* TikTok Embedded Player Card (iPhone Bezel Design) */
.tiktok-player-card {
  display: flex;
  flex-direction: column;
  background: #000;
  border-radius: 32px;
  overflow: hidden;
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.35);
  border: 5px solid #222;
}

.tiktok-phone-bezel {
  position: relative;
  width: 100%;
  aspect-ratio: 9 / 16;
  max-height: 560px;
  background: #000;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.phone-dynamic-island {
  position: absolute;
  top: 10px;
  left: 50%;
  transform: translateX(-50%);
  width: 76px;
  height: 18px;
  background: #111;
  border-radius: 999px;
  z-index: 10;
  pointer-events: none;
  box-shadow: 0 2px 6px rgba(0,0,0,0.6);
}

.tiktok-embed-frame {
  width: 100%;
  height: 100%;
  border: none;
  display: block;
}

.tiktok-control-strip {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 14px;
  background: #141416;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  color: #fff;
}

.control-strip-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.tiktok-badge {
  font-size: 0.72rem;
  font-weight: 800;
  background: #000;
  border: 1px solid rgba(255, 255, 255, 0.2);
  padding: 2px 8px;
  border-radius: 999px;
  color: #fff;
}

.control-strip-title {
  font-size: 0.76rem;
  color: #aaa;
}

.control-strip-actions {
  display: flex;
  gap: 6px;
}

.video-preview-card {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  padding: 14px 16px;
  box-shadow: var(--shadow-sm);
}

.video-preview-header {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.84rem;
  font-weight: 700;
  color: var(--text-secondary);
  margin-bottom: 10px;
}

.video-platforms-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.video-platform-dual {
  display: flex;
  gap: 6px;
  align-items: center;
}

.inpage-play-btn {
  flex: 1;
  border: 1px solid rgba(255, 45, 85, 0.4) !important;
  background: linear-gradient(135deg, #111, #222) !important;
  cursor: pointer;
}

.inpage-play-btn:hover {
  border-color: #ff2d55 !important;
  box-shadow: 0 4px 14px rgba(255, 45, 85, 0.25) !important;
}

.live-badge {
  background: #ff2d55 !important;
  color: #fff !important;
  font-weight: 700;
}

.video-platform-link-mini {
  width: 38px;
  height: 38px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-md, 12px);
  background: #111;
  border: 1px solid rgba(255, 255, 255, 0.15);
  color: #fff;
  text-decoration: none;
  font-weight: 700;
  font-size: 0.9rem;
  transition: all 0.2s;
  flex-shrink: 0;
}

.video-platform-link-mini:hover {
  background: #ff2d55;
  border-color: #ff2d55;
}

.video-platform-link {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 14px;
  border-radius: var(--radius-md);
  text-decoration: none;
  font-size: 0.86rem;
  font-weight: 600;
  transition: all 0.2s ease;
  color: white;
}

.video-platform-link.platform-douyin {
  background: linear-gradient(135deg, #161823, #222);
}
.video-platform-link.platform-xiaohongshu {
  background: linear-gradient(135deg, #ff2442, #ff4761);
}
.video-platform-link.platform-tiktok {
  background: linear-gradient(135deg, #000000, #111111);
  border: 1px solid rgba(37, 244, 238, 0.4);
}
.video-platform-link.platform-bilibili {
  background: linear-gradient(135deg, #00aeec, #008ac5);
}
.video-platform-link.platform-youtube {
  background: linear-gradient(135deg, #cc0000, #ff0000);
}
.video-platform-link.platform-direct {
  background: var(--text-primary);
  color: var(--bg-primary);
}

.video-platform-link:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 14px rgba(0,0,0,0.18);
  opacity: 0.95;
}

.v-platform-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.v-platform-icon {
  font-size: 1.15rem;
}

.v-jump-badge {
  font-size: 0.75rem;
  opacity: 0.9;
  background: rgba(255,255,255,0.22);
  padding: 2px 8px;
  border-radius: 999px;
}

/* Info */
.detail-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 16px;
}

.detail-title {
  font-size: clamp(1.6rem, 3vw, 2.2rem);
  margin-bottom: 12px;
  line-height: 1.25;
}

.detail-stats {
  display: flex;
  gap: 16px;
  font-size: 0.82rem;
  color: var(--c-text-tertiary);
  margin-bottom: 24px;
}

/* Price */
.price-block {
  padding: 20px;
  background: var(--c-bg);
  border-radius: var(--r-lg);
  margin-bottom: 20px;
}

.member-price { display: flex; align-items: center; gap: 10px; }
.price-label {
  font-size: 0.78rem;
  color: var(--c-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.06em;
  font-weight: 600;
  margin-bottom: 6px;
}
.price-display { display: flex; align-items: baseline; gap: 4px; }
.price-currency { font-size: 1.2rem; font-weight: 700; color: var(--c-text-secondary); }
.price-amount { font-size: 2.4rem; font-weight: 900; letter-spacing: -0.04em; }
.price-unit { font-size: 0.85rem; color: var(--c-text-secondary); }

/* Actions */
.detail-actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 32px;
}

/* Dual Channel Purchase Grid */
.purchase-channels-grid {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.channel-card {
  position: relative;
  background: var(--c-surface, #fff);
  border: 1.5px solid var(--c-border, #e2e8f0);
  border-radius: 18px;
  padding: 18px 20px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.04);
  transition: all 0.25s ease;
}

[data-theme="dark"] .channel-card {
  background: rgba(255, 255, 255, 0.04);
  border-color: rgba(255, 255, 255, 0.1);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
}

.guest-channel-card {
  border-color: rgba(0, 122, 255, 0.4);
  background: linear-gradient(180deg, rgba(0, 122, 255, 0.03) 0%, transparent 100%);
}

.member-channel-card {
  border-color: rgba(245, 158, 11, 0.4);
  background: linear-gradient(180deg, rgba(245, 158, 11, 0.03) 0%, transparent 100%);
}

.channel-card-badge {
  position: absolute;
  top: -10px;
  left: 18px;
  padding: 2px 10px;
  border-radius: 999px;
  font-size: 0.7rem;
  font-weight: 800;
  letter-spacing: 0.2px;
}

.guest-badge {
  background: #007aff;
  color: #fff;
  box-shadow: 0 2px 8px rgba(0, 122, 255, 0.35);
}

.member-badge {
  background: linear-gradient(135deg, #f59e0b, #d97706);
  color: #fff;
  box-shadow: 0 2px 8px rgba(245, 158, 11, 0.35);
}

.channel-card-body {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding-top: 4px;
}

.channel-card-title-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.channel-card-title-row h4 {
  font-size: 1.05rem;
  font-weight: 800;
  margin: 0;
  color: var(--c-text, #111);
}

.channel-card-price {
  font-size: 1.25rem;
  font-weight: 900;
  color: #007aff;
}

.channel-user-pill {
  font-size: 0.72rem;
  font-weight: 600;
  color: #d97706;
  background: rgba(245, 158, 11, 0.15);
  padding: 2px 8px;
  border-radius: 999px;
}

.channel-card-desc {
  font-size: 0.82rem;
  color: #666;
  margin: 0 0 6px;
  line-height: 1.45;
}

.channel-action-btn {
  font-weight: 800;
  border-radius: 12px;
}

.channel-dual-btns {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

@media (max-width: 540px) {
  .channel-dual-btns {
    grid-template-columns: 1fr;
  }
}

.btn-gold {
  background: linear-gradient(135deg, #f59e0b, #d97706);
  color: #fff;
  font-weight: 800;
  box-shadow: 0 4px 14px rgba(245, 158, 11, 0.35);
  border: none;
  text-decoration: none;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-gold:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 18px rgba(245, 158, 11, 0.45);
}

.redeem-entry-bar {
  text-align: center;
  margin-top: 6px;
}

.btn-redeem-link {
  background: none;
  border: none;
  font-size: 0.82rem;
  font-weight: 700;
  color: #666;
  cursor: pointer;
  padding: 6px 12px;
  border-radius: 8px;
  transition: all 0.15s;
}

.btn-redeem-link:hover {
  color: #007aff;
  background: rgba(0, 122, 255, 0.08);
}

/* Unlocked Sections */
.unlocked-section {
  border-radius: var(--r-lg, 16px);
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.free-unlocked {
  background: rgba(16, 185, 129, 0.08);
  border: 1px solid rgba(16, 185, 129, 0.25);
}

.free-header {
  color: #10b981;
}

.vip-unlocked {
  background: rgba(245, 158, 11, 0.08);
  border: 1px solid rgba(245, 158, 11, 0.25);
}

.vip-header {
  color: #d97706;
}

.paid-unlocked {
  background: rgba(16, 185, 129, 0.08);
  border: 1px solid rgba(16, 185, 129, 0.25);
}

.paid-header {
  color: #059669;
}

.unlocked-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 700;
  font-size: 0.92rem;
  line-height: 1.4;
}

.reopen-test-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 6px;
  padding-top: 10px;
  border-top: 1px dashed rgba(0, 0, 0, 0.1);
  font-size: 0.78rem;
  color: #777;
}

[data-theme="dark"] .reopen-test-bar {
  border-top-color: rgba(255, 255, 255, 0.1);
}

/* Description */
.detail-desc h3 {
  font-size: 1rem;
  margin-bottom: 12px;
  padding-bottom: 10px;
  border-bottom: 1px solid var(--c-border-light);
}

.desc-content {
  font-size: 0.92rem;
  color: var(--c-text-secondary);
  line-height: 1.8;
}
.desc-content :deep(p) { margin-bottom: 10px; }
.desc-content :deep(ul) { padding-left: 20px; }
.desc-content :deep(li) { margin-bottom: 6px; }
.desc-content :deep(strong) { color: var(--c-text-primary); }
.desc-content :deep(h1), .desc-content :deep(h2), .desc-content :deep(h3) {
  color: var(--c-text-primary);
  margin: 16px 0 8px;
}

/* Related */
.related-section {
  padding: 60px 0;
  background: var(--c-bg);
  margin-top: 60px;
}

.related-section h2 { margin-bottom: 28px; }

.related-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(215px, 1fr));
  gap: 20px;
}

@media (max-width: 768px) {
  .related-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
  }
}

/* Not found */
.not-found {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  min-height: 60vh;
  font-size: 3rem;
}
.not-found p { font-size: 1rem; color: var(--c-text-secondary); }

/* Free prompt box & single channel */
.free-prompt-box {
  background: linear-gradient(135deg, rgba(0, 113, 227, 0.06), rgba(16, 185, 129, 0.06));
  border: 2px dashed rgba(0, 113, 227, 0.35);
  border-radius: 22px;
  padding: 24px;
}
.free-prompt-box.limit-box {
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.08), rgba(239, 68, 68, 0.06));
  border-color: rgba(245, 158, 11, 0.4);
}
.free-prompt-inner {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 14px;
}
.free-prompt-icon {
  font-size: 2.4rem;
  line-height: 1;
}
.free-prompt-body h4 {
  font-size: 1.15rem;
  font-weight: 800;
  margin: 0 0 6px;
  color: var(--c-text);
}
.free-prompt-body p {
  font-size: 0.88rem;
  color: var(--c-text-secondary);
  line-height: 1.5;
  margin: 0;
  max-width: 520px;
}
.single-channel {
  grid-template-columns: 1fr !important;
}
.full-width {
  width: 100% !important;
}

/* Mobile Sticky Action Bar & Mobile Overrides */
.mobile-sticky-bar {
  display: none;
}

@media (max-width: 768px) {
  .detail-page {
    padding-bottom: calc(96px + env(safe-area-inset-bottom, 0px)) !important;
  }

  .detail-grid {
    grid-template-columns: 1fr !important;
    gap: 20px !important;
    padding-top: 10px !important;
  }

  .detail-image-wrap {
    max-width: min(310px, 84vw) !important;
    position: static !important;
    margin: 0 auto 20px !important;
  }

  .tiktok-player-card {
    max-width: min(310px, 84vw) !important;
    margin: 0 auto !important;
  }

  .detail-title {
    font-size: 1.45rem !important;
    line-height: 1.25 !important;
  }

  .purchase-channels-grid {
    grid-template-columns: 1fr !important;
    gap: 12px !important;
  }

  .channel-card-body {
    padding: 16px !important;
  }

  .channel-dual-btns {
    flex-direction: column !important;
  }

  .mobile-sticky-bar {
    display: block !important;
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    z-index: 150;
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(20px);
    -webkit-backdrop-filter: blur(20px);
    border-top: 1px solid rgba(0, 0, 0, 0.08);
    box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.1);
    padding: 10px 16px max(12px, env(safe-area-inset-bottom));
    animation: slideUp 0.3s ease;
  }

  [data-theme="dark"] .mobile-sticky-bar {
    background: rgba(24, 24, 28, 0.96);
    border-top-color: rgba(255, 255, 255, 0.1);
    box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.4);
  }

  .sticky-bar-content {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    max-width: 500px;
    margin: 0 auto;
  }

  .sticky-price-row {
    display: flex;
    align-items: baseline;
    gap: 2px;
  }

  .sticky-price-currency {
    font-size: 1rem;
    font-weight: 700;
    color: var(--c-accent);
  }

  .sticky-price-num {
    font-size: 1.45rem;
    font-weight: 900;
    color: var(--c-text);
    letter-spacing: -0.5px;
  }

  .sticky-price-unit {
    font-size: 0.72rem;
    font-weight: 600;
    color: var(--c-text-tertiary);
    margin-left: 2px;
  }

  .sticky-status-badge {
    display: inline-flex;
    align-items: center;
    font-size: 0.82rem;
    font-weight: 800;
    padding: 4px 10px;
    border-radius: 999px;
  }

  .sticky-bar-actions {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .sticky-btn {
    padding: 10px 14px;
    font-size: 0.82rem;
    font-weight: 800;
    border-radius: 999px;
    white-space: nowrap;
    min-height: 40px;
  }
}

/* ========================================================
   Universal Setup Tutorial Box
   ======================================================== */
.universal-tutorial-box {
  padding: 24px;
  margin-top: 24px;
  border-radius: 20px;
  background: var(--c-bg-card, rgba(255, 255, 255, 0.03));
  border: 1px solid var(--c-border, rgba(255, 255, 255, 0.08));
}

.tut-header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
}

.tut-header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.tut-header-icon {
  font-size: 1.35rem;
}

.tut-header-title {
  font-size: 1.15rem;
  font-weight: 800;
  color: var(--c-text);
  margin: 0;
}

.tut-header-badge {
  font-size: 0.75rem;
  font-weight: 700;
  padding: 3px 10px;
  border-radius: 999px;
  background: rgba(14, 165, 233, 0.15);
  color: #38bdf8;
  border: 1px solid rgba(14, 165, 233, 0.3);
}

.tut-notice-bar {
  padding: 10px 14px;
  margin-bottom: 16px;
  border-radius: 12px;
  background: rgba(245, 158, 11, 0.1);
  border: 1px solid rgba(245, 158, 11, 0.25);
  color: #fbbf24;
  font-size: 0.85rem;
  line-height: 1.5;
}

.tut-versions-row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 16px;
}

.tut-version-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  border-radius: 12px;
  border: 1px solid var(--c-border, rgba(255, 255, 255, 0.1));
  background: rgba(255, 255, 255, 0.04);
  color: var(--c-text-secondary);
  font-size: 0.88rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s ease;
}

.tut-version-pill:hover {
  background: rgba(255, 255, 255, 0.08);
}

.tut-version-pill.active {
  background: rgba(14, 165, 233, 0.15);
  border-color: #38bdf8;
  color: #38bdf8;
  box-shadow: 0 4px 14px rgba(14, 165, 233, 0.2);
}

.tut-version-pill.active.purple {
  background: rgba(168, 85, 247, 0.15);
  border-color: #c084fc;
  color: #c084fc;
  box-shadow: 0 4px 14px rgba(168, 85, 247, 0.2);
}

.tut-version-pill.active.emerald {
  background: rgba(16, 185, 129, 0.15);
  border-color: #34d399;
  color: #34d399;
  box-shadow: 0 4px 14px rgba(16, 185, 129, 0.2);
}

.tut-version-pill.active.amber {
  background: rgba(245, 158, 11, 0.15);
  border-color: #fbbf24;
  color: #fbbf24;
  box-shadow: 0 4px 14px rgba(245, 158, 11, 0.2);
}

.pill-sub-tag {
  font-size: 0.7rem;
  padding: 1px 6px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.15);
  font-weight: 800;
}

.tut-active-card {
  padding: 18px;
  border-radius: 14px;
  background: rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.05);
}

.tut-active-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
  flex-wrap: wrap;
}

.tut-subtitle {
  font-size: 1rem;
  font-weight: 700;
  color: var(--c-text);
  margin: 0;
  flex: 1;
}

.tut-video-btn {
  white-space: nowrap;
  font-weight: 800;
  padding: 6px 14px;
  border-radius: 999px;
  font-size: 0.82rem;
  text-decoration: none;
}

.tut-steps-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.tut-step-row {
  display: flex;
  align-items: flex-start;
  gap: 10px;
}

.tut-step-num {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: rgba(14, 165, 233, 0.15);
  border: 1px solid rgba(14, 165, 233, 0.3);
  color: #38bdf8;
  font-size: 0.75rem;
  font-weight: 800;
  flex-shrink: 0;
  margin-top: 1px;
}

.tut-step-text {
  font-size: 0.88rem;
  line-height: 1.5;
  color: var(--c-text-secondary);
}
</style>
