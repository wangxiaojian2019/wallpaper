<template>
  <header class="navbar" :class="{ scrolled: isScrolled }">
    <div class="container navbar-inner">
      <!-- Logo -->
      <router-link to="/" class="navbar-logo">
        <span class="logo-icon">✦</span>
        <span class="logo-text">WallArt</span>
      </router-link>

      <!-- Desktop Nav -->
      <nav class="navbar-links hide-mobile">
        <router-link to="/" class="nav-link" :class="{ active: $route.name === 'Home' }">{{ t('nav.home') }}</router-link>
        <router-link to="/membership" class="nav-link" :class="{ active: $route.name === 'Membership' }">
          <span class="nav-link-crown">👑</span> {{ t('nav.membership') }}
        </router-link>
        <button class="nav-link nav-btn-link" @click="showLookup = true">
          <span>🔍</span> {{ locale === 'zh' ? '订单找回' : 'Order Lookup' }}
        </button>
      </nav>

      <!-- Right actions -->
      <div class="navbar-actions">
        <!-- Language Switch -->
        <button class="lang-btn" @click="toggleLang" :title="locale === 'zh' ? 'Switch to English' : '切换中文'">
          <span class="lang-icon">🌐</span>
          <span>{{ locale === 'zh' ? 'EN' : '中文' }}</span>
        </button>

        <!-- Logged in -->
        <template v-if="auth.isLoggedIn">
          <!-- Admin direct entrance button -->
          <router-link v-if="auth.isAdmin" to="/admin" class="btn btn-outline btn-sm hide-mobile" style="border-color: #38bdf8; color: #38bdf8; font-weight: 600;">
            ⚙️ {{ locale === 'zh' ? '管理后台' : 'Admin' }}
          </router-link>

          <div class="avatar-menu" ref="menuRef">
            <button class="avatar-btn" @click="showMenu = !showMenu">
              <span class="avatar-text">{{ avatarLetter }}</span>
            </button>
            <transition name="slide">
              <div v-if="showMenu" class="dropdown-menu">
                <div class="dropdown-header">
                  <span class="dropdown-email">{{ auth.user?.email }}</span>
                  <span v-if="auth.isMember" class="badge badge-gold dropdown-badge">
                    👑 {{ t(`membership.${auth.membershipType}`) }}
                  </span>
                </div>
                <div class="dropdown-divider"></div>
                <router-link to="/dashboard" class="dropdown-item" @click="showMenu = false">
                  <span>🏠</span> {{ t('nav.dashboard') }}
                </router-link>
                <router-link v-if="auth.isAdmin" to="/admin" class="dropdown-item" @click="showMenu = false">
                  <span>⚙️</span> {{ t('nav.admin') }}
                </router-link>
                <div class="dropdown-divider"></div>
                <button class="dropdown-item danger" @click="handleLogout">
                  <span>🚪</span> {{ t('nav.logout') }}
                </button>
              </div>
            </transition>
          </div>
        </template>

        <!-- Guest -->
        <template v-else>
          <router-link to="/membership" class="btn btn-outline btn-xs show-mobile mobile-vip-pill">
            👑 VIP
          </router-link>
          <router-link to="/login" class="btn btn-ghost btn-sm hide-mobile">{{ t('nav.login') }}</router-link>
          <router-link to="/register" class="btn btn-primary btn-sm hide-mobile">{{ t('nav.register') }}</router-link>
        </template>

        <!-- Mobile menu btn -->
        <button class="mobile-menu-btn show-mobile" @click="showMobile = !showMobile" :title="locale === 'zh' ? '展开导航' : 'Menu'">
          <span class="hamburger" :class="{ open: showMobile }"></span>
        </button>
      </div>
    </div>

    <!-- Mobile menu -->
    <transition name="slide">
      <div v-if="showMobile" class="mobile-menu">
        <router-link to="/" class="mobile-nav-link" @click="showMobile = false">{{ t('nav.home') }}</router-link>
        <router-link to="/membership" class="mobile-nav-link" @click="showMobile = false">👑 {{ t('nav.membership') }}</router-link>
        <template v-if="auth.isLoggedIn">
          <router-link to="/dashboard" class="mobile-nav-link" @click="showMobile = false">{{ t('nav.dashboard') }}</router-link>
          <router-link v-if="auth.isAdmin" to="/admin" class="mobile-nav-link" @click="showMobile = false">{{ t('nav.admin') }}</router-link>
          <button class="mobile-nav-link danger" @click="handleLogout">{{ t('nav.logout') }}</button>
        </template>
        <template v-else>
          <router-link to="/login" class="mobile-nav-link" @click="showMobile = false">{{ t('nav.login') }}</router-link>
          <router-link to="/register" class="mobile-nav-link accent" @click="showMobile = false">{{ t('nav.register') }}</router-link>
        </template>
        <button class="mobile-nav-link" @click="showLookup = true; showMobile = false">
          🔍 {{ locale === 'zh' ? '订单与下载找回' : 'Order Lookup' }}
        </button>
      </div>
    </transition>
    <OrderLookupModal v-model="showLookup" />
  </header>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store'
import { useToast } from '@/composables/useToast'
import OrderLookupModal from '@/components/OrderLookupModal.vue'

const { t, locale } = useI18n()
const router = useRouter()
const auth = useAuthStore()
const { success } = useToast()

const isScrolled = ref(false)
const showMenu = ref(false)
const showMobile = ref(false)
const showLookup = ref(false)
const menuRef = ref(null)

const avatarLetter = computed(() => {
  const email = auth.user?.email || ''
  return email.charAt(0).toUpperCase()
})

function toggleLang() {
  locale.value = locale.value === 'zh' ? 'en' : 'zh'
  localStorage.setItem('lang', locale.value)
}

function handleLogout() {
  auth.logout()
  showMenu.value = false
  showMobile.value = false
  success(t('auth.logout_success'))
  router.push('/')
}

function handleScroll() {
  isScrolled.value = window.scrollY > 20
}

function handleClickOutside(e) {
  if (menuRef.value && !menuRef.value.contains(e.target)) {
    showMenu.value = false
  }
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll)
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 200;
  height: var(--nav-h);
  background: rgba(255,255,255,0.85);
  backdrop-filter: saturate(180%) blur(20px);
  -webkit-backdrop-filter: saturate(180%) blur(20px);
  border-bottom: 1px solid transparent;
  transition: all var(--t-normal);
}

.navbar.scrolled {
  background: rgba(255,255,255,0.92);
  border-bottom-color: var(--c-border-light);
  box-shadow: 0 1px 20px rgba(0,0,0,0.06);
}

.navbar-inner {
  display: flex;
  align-items: center;
  height: 100%;
  gap: 24px;
}

/* Logo */
.navbar-logo {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 800;
  font-size: 1.25rem;
  color: var(--c-text-primary);
  letter-spacing: -0.03em;
  flex-shrink: 0;
}
.navbar-logo:hover { color: var(--c-text-primary); }
.logo-icon {
  color: var(--c-accent);
  font-size: 1.1rem;
}
.logo-text { color: var(--c-text-primary); }

/* Nav links */
.navbar-links {
  display: flex;
  align-items: center;
  gap: 4px;
  flex: 1;
}

.nav-link {
  padding: 6px 14px;
  border-radius: var(--r-full);
  font-size: 0.9rem;
  font-weight: 500;
  color: var(--c-text-secondary);
  transition: all var(--t-fast);
  display: flex;
  align-items: center;
  gap: 4px;
}
.nav-link:hover, .nav-link.active {
  color: var(--c-text-primary);
  background: var(--c-surface-hover);
}
.nav-link-crown { font-size: 0.85rem; }

.nav-btn-link {
  background: transparent;
  border: none;
  cursor: pointer;
  font-family: inherit;
}

/* Actions */
.navbar-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-left: auto;
}

/* Language btn */
.lang-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  border-radius: var(--r-full);
  background: transparent;
  border: 1.5px solid var(--c-border);
  font-size: 0.82rem;
  font-weight: 600;
  color: var(--c-text-secondary);
  cursor: pointer;
  transition: all var(--t-fast);
}
.lang-btn:hover {
  border-color: var(--c-text-primary);
  color: var(--c-text-primary);
}
.lang-icon { font-size: 0.9rem; }

/* Avatar menu */
.avatar-menu { position: relative; }

.avatar-btn {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: var(--g-accent);
  color: white;
  font-weight: 700;
  font-size: 0.9rem;
  cursor: pointer;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all var(--t-fast);
  box-shadow: var(--shadow-sm);
}
.avatar-btn:hover { transform: scale(1.05); }

.dropdown-menu {
  position: absolute;
  top: calc(100% + 10px);
  right: 0;
  width: 240px;
  background: var(--c-surface);
  border: 1px solid var(--c-border-light);
  border-radius: var(--r-xl);
  box-shadow: var(--shadow-lg);
  overflow: hidden;
  z-index: 300;
}

.dropdown-header {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.dropdown-email {
  font-size: 0.82rem;
  color: var(--c-text-secondary);
  font-weight: 500;
  word-break: break-all;
}

.dropdown-badge { align-self: flex-start; }

.dropdown-divider {
  height: 1px;
  background: var(--c-border-light);
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 11px 16px;
  font-size: 0.88rem;
  color: var(--c-text-primary);
  cursor: pointer;
  transition: background var(--t-fast);
  width: 100%;
  border: none;
  background: none;
  text-align: left;
  font-family: var(--font-sans);
}
.dropdown-item:hover { background: var(--c-surface-hover); }
.dropdown-item.danger { color: var(--c-error); }

/* Mobile menu button */
.mobile-menu-btn {
  background: none;
  border: none;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.hamburger,
.hamburger::before,
.hamburger::after {
  display: block;
  width: 20px;
  height: 2px;
  background: var(--c-text-primary);
  border-radius: 2px;
  transition: all 0.3s ease;
}
.hamburger { position: relative; }
.hamburger::before,
.hamburger::after {
  content: '';
  position: absolute;
  left: 0;
}
.hamburger::before { top: -6px; }
.hamburger::after { top: 6px; }

.hamburger.open { background: transparent; }
.hamburger.open::before { transform: rotate(45deg); top: 0; }
.hamburger.open::after { transform: rotate(-45deg); top: 0; }

/* Mobile menu */
.mobile-menu {
  background: rgba(255, 255, 255, 0.96);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-top: 1px solid var(--c-border-light);
  border-bottom: 1px solid var(--c-border-light);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.12);
  padding: 14px 16px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.mobile-nav-link {
  display: block;
  padding: 12px 16px;
  border-radius: var(--r-md);
  font-size: 0.95rem;
  font-weight: 500;
  color: var(--c-text-primary);
  transition: background var(--t-fast);
  border: none;
  background: none;
  text-align: left;
  cursor: pointer;
  font-family: var(--font-sans);
  width: 100%;
}
.mobile-nav-link:hover { background: var(--c-surface-hover); }
.mobile-nav-link.danger { color: var(--c-error); }
.mobile-nav-link.accent { color: var(--c-accent); font-weight: 600; }

@media (max-width: 768px) {
  .navbar-inner {
    gap: 8px !important;
  }

  .navbar-logo {
    font-size: 1.1rem;
    gap: 6px;
  }

  .logo-icon {
    font-size: 1rem;
  }

  .navbar-actions {
    gap: 6px;
  }

  .lang-btn {
    padding: 4px 8px;
    font-size: 0.74rem;
  }

  .mobile-vip-pill {
    padding: 4px 10px !important;
    font-size: 0.74rem !important;
    font-weight: 700 !important;
    border-radius: var(--r-full) !important;
    border: 1px solid var(--c-gold) !important;
    color: var(--c-gold) !important;
    background: var(--c-gold-light) !important;
    line-height: 1.2;
  }

  .avatar-btn {
    width: 32px;
    height: 32px;
    font-size: 0.8rem;
  }
}
</style>
