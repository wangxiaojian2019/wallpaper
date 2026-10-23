import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import request from '@/utils/request'
import { startTracking } from '@/utils/tracker'

export const useAppStore = defineStore('app', () => {
  const initialTheme = uni.getStorageSync('app_theme') || 'dark'
  const theme = ref(initialTheme) // 'dark' | 'light'
  
  // Check if user manually saved language previously
  const isLangManual = !!uni.getStorageSync('app_lang_manual')
  const savedLang = uni.getStorageSync('app_lang')
  const savedCurrency = uni.getStorageSync('app_currency')

  const lang = ref(savedLang || 'zh') // 'zh' | 'en'
  const currency = ref(savedCurrency || (lang.value === 'zh' ? 'CNY' : 'USD')) // 'CNY' | 'USD'
  const exchangeRate = ref(7.25)
  const userCountry = ref('')
  const userIP = ref('')
  const regionScope = ref('china') // 'china' | 'overseas'

  // If user hasn't explicitly set language preference, detect based on IP
  async function detectLocale() {
    try {
      const res = await request.get('/config/locale')
      if (res) {
        if (res.exchange_rate) exchangeRate.value = Number(res.exchange_rate)
        if (res.country) userCountry.value = res.country
        if (res.ip) userIP.value = res.ip
        if (res.region_scope) regionScope.value = res.region_scope
        
        // Only auto-override if user never manually switched language
        if (!isLangManual) {
          if (res.default_lang) {
            lang.value = res.default_lang
            uni.setStorageSync('app_lang', res.default_lang)
          }
          if (res.currency) {
            currency.value = res.currency
            uni.setStorageSync('app_currency', res.currency)
          } else {
            currency.value = lang.value === 'zh' ? 'CNY' : 'USD'
            uni.setStorageSync('app_currency', currency.value)
          }
        }
      }
    } catch (e) {
      // Fallback: browser language check if backend fails
      if (!isLangManual && typeof navigator !== 'undefined') {
        const navLang = (navigator.language || navigator.userLanguage || '').toLowerCase()
        if (!navLang.startsWith('zh')) {
          lang.value = 'en'
          currency.value = 'USD'
        }
      }
    }
  }

  // Dual-currency price calculations
  function getProductPriceNumber(item) {
    if (!item) return 0
    if (currency.value === 'USD') {
      if (item.price_usd && Number(item.price_usd) > 0) {
        return Number(item.price_usd)
      }
      if (item.price_cny && Number(item.price_cny) > 0) {
        return Number((Number(item.price_cny) / exchangeRate.value).toFixed(2))
      }
      if (item.price && Number(item.price) > 0) {
        return Number((Number(item.price) / exchangeRate.value).toFixed(2))
      }
      return 0
    } else {
      // CNY
      if (item.price_cny && Number(item.price_cny) > 0) {
        return Number(item.price_cny)
      }
      if (item.price_usd && Number(item.price_usd) > 0) {
        return Number((Number(item.price_usd) * exchangeRate.value).toFixed(2))
      }
      if (item.price && Number(item.price) > 0) {
        return Number(item.price)
      }
      return 0
    }
  }

  function formatPrice(item) {
    const num = getProductPriceNumber(item)
    if (currency.value === 'USD') {
      return `$${num.toFixed(2)}`
    }
    return `¥${num.toFixed(2)}`
  }

  // Tier requirement label - calculates lowest allowed VIP tier
  function getFreeTierLabel(item) {
    if (!item) return ''
    const p = item
    const isZero = Number(p.price || 0) === 0 && Number(p.price_cny || 0) === 0 && Number(p.price_usd || 0) === 0
    if (isZero) {
      return lang.value === 'en' ? '🎁 Free for All' : '🎁 全员免费'
    }

    let list = ['monthly', 'quarterly', 'yearly', 'lifetime']
    const raw = item.allowed_vip_levels
    if (raw && raw !== 'null' && raw !== '') {
      try {
        const parsed = typeof raw === 'string' ? JSON.parse(raw) : raw
        if (Array.isArray(parsed)) {
          list = parsed.map(x => String(x).toLowerCase().trim())
        }
      } catch (e) {}
    }

    if (Array.isArray(list) && list.length === 0) {
      return lang.value === 'en' ? 'Purchase Only' : '限单独购买'
    }

    const isEn = lang.value === 'en'
    if (list.includes('monthly')) {
      return isEn ? '👑 Monthly+ Free' : '👑 月度会员可免'
    }
    if (list.includes('quarterly')) {
      return isEn ? '👑 Quarterly+ Free' : '👑 季度会员可免'
    }
    if (list.includes('yearly')) {
      return isEn ? '👑 Yearly+ Free' : '👑 年度会员可免'
    }
    if (list.includes('lifetime')) {
      return isEn ? '👑 Lifetime Free' : '👑 终身会员可免'
    }

    return isEn ? '👑 VIP Free' : '👑 VIP会员可免'
  }

  // Trigger locale detection and visit tracking
  detectLocale()
  if (typeof window !== 'undefined') {
    startTracking(window.location.hash || window.location.pathname)
  }

  function applyTheme(val) {
    if (typeof document !== 'undefined') {
      document.documentElement.classList.remove('dark', 'light')
      document.documentElement.classList.add(val)
      document.body.classList.remove('dark', 'light')
      document.body.classList.add(val)
    }

    try {
      if (val === 'dark') {
        uni.setTabBarStyle({
          backgroundColor: '#121420',
          color: '#94a3b8',
          selectedColor: '#38bdf8',
          borderStyle: 'black'
        })
      } else {
        uni.setTabBarStyle({
          backgroundColor: '#ffffff',
          color: '#64748b',
          selectedColor: '#0f172a',
          borderStyle: 'white'
        })
      }
    } catch (e) {
      // ignore
    }
  }

  // Apply immediately upon initialization
  applyTheme(theme.value)

  let lastThemeToggle = 0
  function toggleTheme() {
    const now = Date.now()
    if (now - lastThemeToggle < 300) return
    lastThemeToggle = now
    theme.value = theme.value === 'dark' ? 'light' : 'dark'
    uni.setStorageSync('app_theme', theme.value)
    applyTheme(theme.value)
  }

  let lastLangToggle = 0
  function toggleLang() {
    const now = Date.now()
    if (now - lastLangToggle < 300) return
    lastLangToggle = now
    lang.value = lang.value === 'zh' ? 'en' : 'zh'
    currency.value = lang.value === 'zh' ? 'CNY' : 'USD'
    uni.setStorageSync('app_lang_manual', '1')
    uni.setStorageSync('app_lang', lang.value)
    uni.setStorageSync('app_currency', currency.value)
  }

  return {
    theme,
    lang,
    currency,
    exchangeRate,
    userCountry,
    userIP,
    regionScope,
    toggleTheme,
    toggleLang,
    detectLocale,
    getProductPriceNumber,
    formatPrice,
    getFreeTierLabel
  }
})
