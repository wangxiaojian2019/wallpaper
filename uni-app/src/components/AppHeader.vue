<template>
  <view class="app-header" :class="[appStore.theme]">
    <view class="header-inner">
      <!-- Brand Logo -->
      <view class="brand-logo" @click="goHome">
        <text class="logo-sparkle">✦</text>
        <text class="logo-text">WallArt</text>
        <text class="logo-sub">4K</text>
      </view>

      <!-- Desktop Nav Links (PC only) -->
      <view class="desktop-nav-links">
        <view class="nav-link-item" :class="{ active: activeNav === 'home' }" @click="goHome">
          <text class="nav-link-text">{{ appStore.lang === 'zh' ? '画廊' : 'Gallery' }}</text>
        </view>
        <view class="nav-link-item" :class="{ active: activeNav === 'vip' }" @click="goToVip">
          <text class="nav-link-text">{{ appStore.lang === 'zh' ? 'VIP 特权' : 'VIP Access' }}</text>
        </view>
        <view class="nav-link-item" :class="{ active: activeNav === 'my' }" @click="goToMy">
          <text class="nav-link-text">{{ appStore.lang === 'zh' ? '个人中心' : 'Account' }}</text>
        </view>
      </view>

      <!-- Right Controls: Theme + Lang + Lookup -->
      <view class="header-controls">
        <!-- Mobile Quick VIP Button -->
        <view class="control-btn mobile-vip-btn" @click.stop="goToVip">
          <text class="control-icon">👑</text>
          <text class="control-text vip-gold-text">VIP</text>
        </view>

        <!-- Mobile Quick My Button -->
        <view class="control-btn mobile-my-btn" @click.stop="goToMy" title="个人中心">
          <text class="control-icon">👤</text>
        </view>

        <!-- Language Switcher -->
        <view class="control-btn lang-btn" @click.stop="appStore.toggleLang">
          <text class="control-text">{{ appStore.lang === 'zh' ? 'EN' : '中文' }}</text>
        </view>

        <!-- Theme Switcher Capsule (Clearly Labeled) -->
        <view class="control-btn theme-toggle-capsule" @click.stop="appStore.toggleTheme" :class="[appStore.theme]" :title="appStore.theme === 'dark' ? '点击切换为浅色雨露模式' : '点击切换为深色星空模式'">
          <text class="control-icon">{{ appStore.theme === 'dark' ? '🌙' : '☀️' }}</text>
          <text class="control-text theme-label-text">
            {{ appStore.theme === 'dark' ? (appStore.lang === 'zh' ? '深色星空' : 'Dark Space') : (appStore.lang === 'zh' ? '浅色雨露' : 'Light Rain') }}
          </text>
        </view>

        <!-- Order Lookup -->
        <view class="control-btn lookup-btn" @click.stop="goToLookup">
          <text class="control-icon">🔍</text>
          <text class="control-text">{{ appStore.lang === 'zh' ? '找回订单' : 'Lookup' }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed } from 'vue'
import { useAppStore } from '@/store/app'

const appStore = useAppStore()

const activeNav = computed(() => {
  // #ifdef H5
  if (typeof window !== 'undefined' && window.location && window.location.hash) {
    const hash = window.location.hash
    if (hash.includes('membership')) return 'vip'
    if (hash.includes('my')) return 'my'
    return 'home'
  }
  // #endif
  try {
    const pages = getCurrentPages()
    if (pages && pages.length > 0) {
      const route = pages[pages.length - 1]?.route || ''
      if (route.includes('membership')) return 'vip'
      if (route.includes('my')) return 'my'
    }
  } catch (e) {}
  return 'home'
})

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

function goToLookup() {
  uni.navigateTo({ url: '/pages/lookup/index' })
}
</script>

<style scoped>
.app-header {
  position: sticky;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  backdrop-filter: blur(24px) saturate(180%);
  -webkit-backdrop-filter: blur(24px) saturate(180%);
  transition: background 0.25s ease, border-color 0.25s ease;
}

.app-header.dark {
  background: rgba(14, 16, 24, 0.62);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.app-header.light {
  background: rgba(255, 255, 255, 0.65);
  border-bottom: 1px solid rgba(0, 0, 0, 0.07);
}

.header-inner {
  max-width: 1240px;
  margin: 0 auto;
  height: calc(56px + env(safe-area-inset-top));
  padding: env(safe-area-inset-top) 24px 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-sizing: border-box;
}

/* Brand */
.brand-logo {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.logo-sparkle {
  font-size: 16px;
  color: #3b82f6;
}

.logo-text {
  font-size: 18px;
  font-weight: 800;
  letter-spacing: -0.5px;
}

.dark .logo-text {
  color: #ffffff;
}

.light .logo-text {
  color: #111827;
}

.logo-sub {
  font-size: 11px;
  font-weight: 700;
  padding: 1px 6px;
  border-radius: 4px;
  background: linear-gradient(135deg, #3b82f6, #06b6d4);
  color: #ffffff;
}

/* Desktop Nav Links */
.desktop-nav-links {
  display: none;
  align-items: center;
  gap: 24px;
}

@media (min-width: 768px) {
  .desktop-nav-links {
    display: flex;
  }
}

.nav-link-item {
  cursor: pointer;
  padding: 6px 12px;
  border-radius: 999px;
  transition: all 0.2s ease;
}

.nav-link-text {
  font-size: 14px;
  font-weight: 600;
}

.dark .nav-link-text {
  color: rgba(255, 255, 255, 0.7);
}

.light .nav-link-text {
  color: rgba(0, 0, 0, 0.65);
}

.nav-link-item.active .nav-link-text,
.nav-link-item:hover .nav-link-text {
  color: #3b82f6;
}

/* Controls */
.header-controls {
  display: flex;
  align-items: center;
  gap: 10px;
}

.control-btn {
  height: 34px;
  padding: 0 14px;
  border-radius: 999px;
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.mobile-vip-btn {
  display: flex;
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.2), rgba(212, 175, 55, 0.1)) !important;
  border: 1px solid rgba(245, 158, 11, 0.4) !important;
}

.mobile-my-btn {
  display: flex;
  padding: 0 10px !important;
}

.vip-gold-text {
  color: #f59e0b !important;
  font-weight: 800 !important;
}

@media (min-width: 768px) {
  .mobile-vip-btn,
  .mobile-my-btn {
    display: none !important;
  }
}

.dark .control-btn {
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #f3f4f6;
}

.light .control-btn {
  background: rgba(0, 0, 0, 0.04);
  border: 1px solid rgba(0, 0, 0, 0.08);
  color: #1f2937;
}

.dark .theme-toggle-capsule {
  background: rgba(30, 41, 59, 0.7);
  border: 1px solid rgba(56, 189, 248, 0.35);
  box-shadow: 0 2px 10px rgba(56, 189, 248, 0.15);
}

.light .theme-toggle-capsule {
  background: #ffffff;
  border: 1px solid rgba(15, 23, 42, 0.15);
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.08);
}

.control-btn:hover {
  transform: translateY(-1px);
}

.control-icon {
  font-size: 13px;
}

.control-text {
  font-size: 12px;
  font-weight: 600;
}

/* ========================================================
   Mobile Optimization (<768px & <360px)
   ======================================================== */
@media (max-width: 768px) {
  .app-header {
    position: sticky;
    top: calc(8px + env(safe-area-inset-top));
    left: 0;
    right: 0;
    z-index: 100;
    background: transparent !important;
    border-bottom: none !important;
    backdrop-filter: none !important;
    -webkit-backdrop-filter: none !important;
    padding: 0 12px;
    box-sizing: border-box;
  }

  .header-inner {
    height: 48px;
    padding: 0 14px;
    margin: 0 auto;
    border-radius: 999px; /* Dock-like curved capsule */
    backdrop-filter: blur(28px) saturate(210%);
    -webkit-backdrop-filter: blur(28px) saturate(210%);
    transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
    box-sizing: border-box;
  }

  .dark .header-inner {
    background: rgba(18, 20, 32, 0.85);
    border: 1px solid rgba(255, 255, 255, 0.18);
    box-shadow: 0 8px 30px rgba(0, 0, 0, 0.4), inset 0 1px 1px rgba(255, 255, 255, 0.15);
  }

  .light .header-inner {
    background: rgba(255, 255, 255, 0.88);
    border: 1px solid rgba(0, 0, 0, 0.1);
    box-shadow: 0 8px 26px rgba(0, 0, 0, 0.08), inset 0 1px 2px rgba(255, 255, 255, 0.9);
  }

  .brand-logo {
    gap: 5px;
  }

  .logo-sparkle {
    font-size: 14px;
  }

  .logo-text {
    font-size: 16px;
  }

  .logo-sub {
    font-size: 9px;
    padding: 1px 4px;
  }

  .header-controls {
    gap: 6px;
  }

  .control-btn {
    height: 30px;
    padding: 0 8px;
    gap: 4px;
  }

  .control-btn.theme-toggle-capsule {
    padding: 0 8px;
  }

  .theme-label-text {
    display: none !important;
  }

  .lookup-btn .control-text {
    display: none !important;
  }

  .lookup-btn {
    padding: 0 8px;
  }

  .lang-btn {
    padding: 0 8px;
    font-size: 11px;
  }
}

@media (max-width: 360px) {
  .logo-sub {
    display: none !important;
  }
  .header-inner {
    padding: env(safe-area-inset-top) 8px 0;
  }
  .header-controls {
    gap: 4px;
  }
}
</style>
