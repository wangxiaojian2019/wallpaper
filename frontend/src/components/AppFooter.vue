<template>
  <footer class="footer">
    <div class="container">
      <div class="footer-grid">
        <!-- Brand -->
        <div class="footer-brand">
          <div class="brand-logo">
            <span class="logo-icon">✦</span>
            <span>WallArt</span>
          </div>
          <p class="brand-desc">{{ t('home.hero_subtitle') }}</p>
          <div class="social-links">
            <a href="#" class="social-btn" title="Twitter/X">𝕏</a>
            <a href="#" class="social-btn" title="Instagram">📷</a>
          </div>
        </div>

        <!-- Links -->
        <div class="footer-col">
          <h4>{{ locale === 'zh' ? '导航' : 'Navigation' }}</h4>
          <ul>
            <li><router-link to="/">{{ t('nav.home') }}</router-link></li>
            <li><router-link to="/membership">{{ t('nav.membership') }}</router-link></li>
            <li><a href="#" @click.prevent="showLookup = true">{{ locale === 'zh' ? '🔍 订单与下载找回' : '🔍 Order Lookup' }}</a></li>
          </ul>
        </div>

        <div class="footer-col">
          <h4>{{ locale === 'zh' ? '账户' : 'Account' }}</h4>
          <ul>
            <li><router-link to="/login">{{ t('nav.login') }}</router-link></li>
            <li><router-link to="/register">{{ t('nav.register') }}</router-link></li>
            <li v-if="auth.isLoggedIn"><router-link to="/dashboard">{{ t('nav.dashboard') }}</router-link></li>
          </ul>
        </div>

        <div class="footer-col">
          <h4>{{ locale === 'zh' ? '支付方式' : 'Payment' }}</h4>
          <div class="payment-badges">
            <span class="payment-badge">PayPal</span>
          </div>
          <h4 style="margin-top:16px">{{ locale === 'zh' ? '网盘支持' : 'Cloud Drives' }}</h4>
          <div class="drive-badges">
            <span class="drive-badge">百度网盘</span>
            <span class="drive-badge">夸克</span>
            <span class="drive-badge">Google Drive</span>
          </div>
        </div>
      </div>

      <div class="footer-bottom">
        <p>© {{ year }} WallArt. {{ locale === 'zh' ? '保留所有权利' : 'All rights reserved' }}.</p>
        <div class="footer-bottom-links">
          <a href="#">{{ locale === 'zh' ? '隐私政策' : 'Privacy Policy' }}</a>
          <a href="#">{{ locale === 'zh' ? '使用条款' : 'Terms of Service' }}</a>
        </div>
      </div>
    </div>
    <OrderLookupModal v-model="showLookup" />
  </footer>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/store'
import OrderLookupModal from '@/components/OrderLookupModal.vue'

const { t, locale } = useI18n()
const auth = useAuthStore()
const year = computed(() => new Date().getFullYear())
const showLookup = ref(false)
</script>

<style scoped>
.footer {
  background: var(--c-text-primary);
  color: white;
  padding: 64px 0 0;
  margin-top: auto;
}

.footer-grid {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1.5fr;
  gap: 40px;
  margin-bottom: 48px;
}

@media (max-width: 900px) {
  .footer-grid {
    grid-template-columns: 1fr 1fr;
  }
}

@media (max-width: 560px) {
  .footer-grid {
    grid-template-columns: 1fr;
  }
}

.brand-logo {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 1.3rem;
  font-weight: 800;
  letter-spacing: -0.03em;
  margin-bottom: 12px;
}
.logo-icon { color: #5eb5ff; }

.brand-desc {
  color: rgba(255,255,255,0.5);
  font-size: 0.88rem;
  line-height: 1.6;
  margin-bottom: 20px;
}

.social-links { display: flex; gap: 8px; }
.social-btn {
  width: 36px;
  height: 36px;
  background: rgba(255,255,255,0.1);
  border-radius: var(--r-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.9rem;
  transition: background var(--t-fast);
}
.social-btn:hover { background: rgba(255,255,255,0.2); }

.footer-col h4 {
  font-size: 0.82rem;
  font-weight: 700;
  color: rgba(255,255,255,0.5);
  text-transform: uppercase;
  letter-spacing: 0.08em;
  margin-bottom: 14px;
  color: white;
}

.footer-col ul { display: flex; flex-direction: column; gap: 8px; }
.footer-col a {
  color: rgba(255,255,255,0.65);
  font-size: 0.9rem;
  transition: color var(--t-fast);
}
.footer-col a:hover { color: white; }

.payment-badges, .drive-badges {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.payment-badge, .drive-badge {
  padding: 4px 10px;
  background: rgba(255,255,255,0.1);
  border: 1px solid rgba(255,255,255,0.15);
  border-radius: var(--r-full);
  font-size: 0.75rem;
  font-weight: 600;
  color: rgba(255,255,255,0.8);
}

.footer-bottom {
  border-top: 1px solid rgba(255,255,255,0.1);
  padding: 24px 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.footer-bottom p {
  font-size: 0.82rem;
  color: rgba(255,255,255,0.4);
}

.footer-bottom-links {
  display: flex;
  gap: 20px;
}
.footer-bottom-links a {
  font-size: 0.82rem;
  color: rgba(255,255,255,0.4);
  transition: color var(--t-fast);
}
.footer-bottom-links a:hover { color: rgba(255,255,255,0.8); }

@media (max-width: 560px) {
  .footer-bottom {
    flex-direction: column;
    text-align: center;
  }
}
</style>
