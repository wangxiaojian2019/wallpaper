<template>
  <nav v-if="showTabBar" class="mobile-tab-bar">
    <router-link
      to="/"
      class="tab-item"
      :class="{ active: $route.path === '/' }"
    >
      <span class="tab-icon">🌟</span>
      <span class="tab-label">{{ locale === 'zh' ? '发现' : 'Explore' }}</span>
    </router-link>

    <router-link
      to="/membership"
      class="tab-item"
      :class="{ active: $route.path === '/membership' }"
    >
      <span class="tab-icon vip-icon">👑</span>
      <span class="tab-label">{{ locale === 'zh' ? 'VIP会员' : 'VIP' }}</span>
      <span v-if="auth.isMember" class="tab-dot"></span>
    </router-link>

    <router-link
      to="/dashboard"
      class="tab-item"
      :class="{ active: $route.path === '/dashboard' }"
    >
      <span class="tab-icon">📦</span>
      <span class="tab-label">{{ locale === 'zh' ? '订单/下载' : 'Orders' }}</span>
    </router-link>

    <router-link
      :to="auth.isLoggedIn ? '/dashboard' : '/login'"
      class="tab-item"
      :class="{ active: $route.path === '/login' || $route.path === '/register' }"
    >
      <span class="tab-icon">
        <template v-if="auth.isLoggedIn">👤</template>
        <template v-else>🔑</template>
      </span>
      <span class="tab-label">
        {{ auth.isLoggedIn ? (locale === 'zh' ? '我的' : 'Me') : (locale === 'zh' ? '登录' : 'Login') }}
      </span>
    </router-link>
  </nav>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/store'

const route = useRoute()
const { locale } = useI18n()
const auth = useAuthStore()

// Hide tab bar on admin views and detail pages (detail page has its own sticky action bar)
const showTabBar = computed(() => {
  const path = route.path || ''
  if (path.startsWith('/admin')) return false
  if (path.startsWith('/wallpaper/')) return false
  return true
})
</script>

<style scoped>
.mobile-tab-bar {
  display: none;
}

@media (max-width: 768px) {
  .mobile-tab-bar {
    display: flex;
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    z-index: 180;
    height: calc(56px + env(safe-area-inset-bottom, 0px));
    padding-bottom: env(safe-area-inset-bottom, 0px);
    background: rgba(255, 255, 255, 0.92);
    backdrop-filter: saturate(180%) blur(20px);
    -webkit-backdrop-filter: saturate(180%) blur(20px);
    border-top: 0.5px solid rgba(0, 0, 0, 0.12);
    box-shadow: 0 -2px 16px rgba(0, 0, 0, 0.05);
    align-items: center;
    justify-content: space-around;
  }

  [data-theme="dark"] .mobile-tab-bar {
    background: rgba(24, 24, 28, 0.94);
    border-top-color: rgba(255, 255, 255, 0.12);
    box-shadow: 0 -2px 16px rgba(0, 0, 0, 0.35);
  }

  .tab-item {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    color: #8e8e93;
    text-decoration: none;
    position: relative;
    transition: transform 0.15s ease, color 0.15s ease;
    -webkit-tap-highlight-color: transparent;
  }

  .tab-item:active {
    transform: scale(0.92);
  }

  .tab-icon {
    font-size: 1.25rem;
    line-height: 1;
    margin-bottom: 2px;
    transition: transform 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  }

  .tab-label {
    font-size: 0.68rem;
    font-weight: 600;
    letter-spacing: -0.2px;
  }

  .tab-item.active {
    color: var(--c-accent, #0071e3);
  }

  .tab-item.active .tab-icon {
    transform: translateY(-2px) scale(1.1);
  }

  .vip-icon {
    filter: drop-shadow(0 1px 2px rgba(245, 158, 11, 0.3));
  }

  .tab-item.active .vip-icon {
    filter: drop-shadow(0 2px 6px rgba(245, 158, 11, 0.5));
  }

  .tab-dot {
    position: absolute;
    top: 6px;
    right: calc(50% - 14px);
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: #ff9500;
    box-shadow: 0 0 4px #ff9500;
  }
}
</style>
