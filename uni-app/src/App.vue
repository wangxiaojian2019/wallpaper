<script>
import { useAppStore } from '@/store/app'
import { startTracking } from '@/utils/tracker'

export default {
  onLaunch: function () {
    console.log('WallArt App Launched')
    try {
      const appStore = useAppStore()
      if (typeof document !== 'undefined') {
        const theme = appStore.theme || 'dark'
        document.documentElement.classList.remove('dark', 'light')
        document.documentElement.classList.add(theme)
        document.body.classList.remove('dark', 'light')
        document.body.classList.add(theme)
      }
    } catch (e) {
      // ignore
    }
  },
  onShow: function () {
    console.log('WallArt App Shown')
    try {
      if (typeof window !== 'undefined') {
        startTracking(window.location.hash || window.location.pathname)
      }
    } catch (e) {}
  },
  onHide: function () {
    console.log('WallArt App Hidden')
  }
}
</script>

<style>
/* Global CSS variables & spacious relaxed theme */
page {
  background-color: transparent;
  color: #f1f2f6;
  font-family: -apple-system, BlinkMacSystemFont, "SF Pro Text", "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  box-sizing: border-box;
}

/* Reset button default ugly styles in uni-app/mini-program */
button {
  border: none;
  margin: 0;
  padding: 0;
  line-height: normal;
  background: transparent;
  outline: none;
}
button::after {
  border: none;
}

/* Smooth scrolling */
html {
  scroll-behavior: smooth;
}

/* Hide scrollbars for clean look */
::-webkit-scrollbar {
  display: none;
  width: 0 !important;
  height: 0 !important;
  -webkit-appearance: none;
  background: transparent;
}

/* Hide mobile tabbar on PC desktop view (768px+) */
/* #ifdef H5 */
@media (min-width: 768px) {
  uni-tabbar,
  .uni-tabbar,
  .uni-tabbar-bottom,
  .uni-app--showtabbar uni-page-wrapper::after {
    display: none !important;
  }
  body, page, uni-page-body {
    padding-bottom: 0 !important;
  }
}

/* Signature iOS Floating Glass Dock Bottom TabBar on Mobile (<768px) */
@media (max-width: 768px) {
  body, page, uni-page-body, .home-container, .vip-container, .my-container {
    padding-bottom: calc(96px + env(safe-area-inset-bottom)) !important;
  }

  uni-tabbar,
  uni-tabbar.uni-tabbar-bottom,
  .uni-tabbar-bottom {
    position: fixed !important;
    bottom: 0 !important;
    left: 0 !important;
    right: 0 !important;
    width: 100% !important;
    height: auto !important;
    display: flex !important;
    justify-content: center !important;
    pointer-events: none !important;
    z-index: 9999 !important;
    background: transparent !important;
    border: none !important;
  }

  uni-tabbar .uni-tabbar,
  .uni-tabbar-bottom .uni-tabbar {
    position: fixed !important;
    bottom: calc(14px + env(safe-area-inset-bottom)) !important;
    left: 50% !important;
    transform: translateX(-50%) !important;
    width: calc(100% - 36px) !important;
    max-width: 380px !important;
    height: 56px !important;
    margin: 0 !important;
    padding: 0 8px !important;
    border-radius: 999px !important;
    pointer-events: auto !important;
    backdrop-filter: blur(28px) saturate(210%) !important;
    -webkit-backdrop-filter: blur(28px) saturate(210%) !important;
    display: flex !important;
    align-items: center !important;
    justify-content: space-around !important;
    box-sizing: border-box !important;
    transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1) !important;
    overflow: hidden !important;
  }

  /* Hide default top 1px border */
  .uni-tabbar-border {
    display: none !important;
  }

  .uni-tabbar__item {
    flex: 1 !important;
    height: 100% !important;
    display: flex !important;
    flex-direction: column !important;
    align-items: center !important;
    justify-content: center !important;
    padding: 4px 0 !important;
    cursor: pointer !important;
  }

  .uni-tabbar__bd {
    display: flex !important;
    flex-direction: column !important;
    align-items: center !important;
    justify-content: center !important;
  }

  /* Modern glyph icons via CSS pseudo elements (.uni-tabbar-border is child 1) */
  .uni-tabbar__item:nth-child(2) .uni-tabbar__bd::before {
    content: '✦';
    font-size: 15px;
    line-height: 1;
    margin-bottom: 3px;
    transition: transform 0.2s ease;
  }

  .uni-tabbar__item:nth-child(3) .uni-tabbar__bd::before {
    content: '👑';
    font-size: 15px;
    line-height: 1;
    margin-bottom: 3px;
    transition: transform 0.2s ease;
  }

  .uni-tabbar__item:nth-child(4) .uni-tabbar__bd::before {
    content: '👤';
    font-size: 15px;
    line-height: 1;
    margin-bottom: 3px;
    transition: transform 0.2s ease;
  }

  .uni-tabbar__item:active .uni-tabbar__bd::before {
    transform: scale(1.18);
  }

  /* Dark mode bottom dock */
  body.dark uni-tabbar .uni-tabbar,
  html.dark uni-tabbar .uni-tabbar,
  body.dark .uni-tabbar-bottom .uni-tabbar,
  html.dark .uni-tabbar-bottom .uni-tabbar {
    background: rgba(18, 20, 32, 0.88) !important;
    border: 1px solid rgba(255, 255, 255, 0.16) !important;
    box-shadow: 0 12px 40px rgba(0, 0, 0, 0.55), inset 0 1px 1px rgba(255, 255, 255, 0.18) !important;
  }

  body.dark .uni-tabbar__label,
  html.dark .uni-tabbar__label,
  body.dark .uni-tabbar__bd::before,
  html.dark .uni-tabbar__bd::before {
    color: rgba(255, 255, 255, 0.55) !important;
    font-size: 11px !important;
    font-weight: 600 !important;
  }

  body.dark .uni-tabbar__item.active .uni-tabbar__label,
  body.dark .uni-tabbar__item[style*="rgb(255, 255, 255)"] .uni-tabbar__label,
  body.dark .uni-tabbar__label[style*="rgb(255, 255, 255)"] {
    color: #38bdf8 !important;
    font-weight: 700 !important;
  }

  body.dark .uni-tabbar__item.active .uni-tabbar__bd::before,
  body.dark .uni-tabbar__item[style*="rgb(255, 255, 255)"] .uni-tabbar__bd::before {
    color: #38bdf8 !important;
  }

  /* VIP gold highlight when VIP tab is active */
  body.dark .uni-tabbar__item:nth-child(3).active .uni-tabbar__label,
  body.dark .uni-tabbar__item:nth-child(3).active .uni-tabbar__bd::before {
    color: #f59e0b !important;
  }

  /* Light mode bottom dock */
  body.light uni-tabbar .uni-tabbar,
  html.light uni-tabbar .uni-tabbar,
  body.light .uni-tabbar-bottom .uni-tabbar,
  html.light .uni-tabbar-bottom .uni-tabbar {
    background: rgba(255, 255, 255, 0.92) !important;
    border: 1px solid rgba(0, 0, 0, 0.08) !important;
    box-shadow: 0 10px 32px rgba(0, 0, 0, 0.1), inset 0 1px 2px rgba(255, 255, 255, 0.95) !important;
  }

  body.light .uni-tabbar__label,
  html.light .uni-tabbar__label,
  body.light .uni-tabbar__bd::before,
  html.light .uni-tabbar__bd::before {
    color: #64748b !important;
    font-size: 11px !important;
    font-weight: 600 !important;
  }

  body.light .uni-tabbar__item.active .uni-tabbar__label,
  body.light .uni-tabbar__item[style*="rgb(255, 255, 255)"] .uni-tabbar__label,
  body.light .uni-tabbar__label[style*="rgb(255, 255, 255)"] {
    color: #0f172a !important;
    font-weight: 700 !important;
  }

  body.light .uni-tabbar__item.active .uni-tabbar__bd::before,
  body.light .uni-tabbar__item[style*="rgb(255, 255, 255)"] .uni-tabbar__bd::before {
    color: #0f172a !important;
  }

  body.light .uni-tabbar__item:nth-child(3).active .uni-tabbar__label,
  body.light .uni-tabbar__item:nth-child(3).active .uni-tabbar__bd::before {
    color: #d97706 !important;
  }
}
/* #endif */
</style>
