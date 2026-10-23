<template>
  <view class="mockup-container" :class="[appStore.theme]">
    <!-- Mode Switcher Tabs -->
    <view class="mockup-tabs">
      <!-- ⚡ 3D Interactive Tab: Only rendered when hasInteractive is true! -->
      <view
        v-if="hasInteractive"
        class="tab-pill tab-pill-interactive"
        :class="{ active: currentMode === 'interactive' }"
        @tap="setMode('interactive')"
      >
        <text class="tab-icon">✨</text>
        <text class="tab-label">{{ appStore.lang === 'zh' ? '3D实机交互' : '3D Live' }}</text>
        <text class="live-pulse-dot"></text>
      </view>

      <view
        class="tab-pill"
        :class="{ active: currentMode === 'phone' }"
        @tap="setMode('phone')"
      >
        <text class="tab-icon">📱</text>
        <text class="tab-label">{{ appStore.lang === 'zh' ? '锁屏实机' : 'Phone' }}</text>
      </view>
      <view
        class="tab-pill"
        :class="{ active: currentMode === 'desktop' }"
        @tap="setMode('desktop')"
      >
        <text class="tab-icon">💻</text>
        <text class="tab-label">{{ appStore.lang === 'zh' ? '桌面效果' : 'Desktop' }}</text>
      </view>
      <view
        class="tab-pill"
        :class="{ active: currentMode === 'clean' }"
        @tap="setMode('clean')"
      >
        <text class="tab-icon">🖼️</text>
        <text class="tab-label">{{ appStore.lang === 'zh' ? '原图纯享' : 'Full HD' }}</text>
      </view>
    </view>

    <!-- Main Mockup Viewport -->
    <view class="viewport-wrapper" :class="[`mode-${currentMode}`]">
      <!-- MODE 0: 3D Interactive .tendies / Parallax Simulation Frame -->
      <view
        v-if="currentMode === 'interactive' && hasInteractive"
        class="interactive-stage"
        @mousemove="onMouseMove"
        @mouseleave="onMouseLeave"
      >
        <!-- Interactive Instruction Badge -->
        <view class="desktop-interact-tip">
          <view class="tip-badge">
            <text class="tip-dot"></text>
            <text class="tip-text" v-if="isSleeping">{{ appStore.lang === 'zh' ? '✦ 屏幕已熄屏待机：点击屏幕或底部按键即可点亮唤醒 ✦' : '✦ Display Sleeping: Tap screen or button to wake ✦' }}</text>
            <text class="tip-text" v-else-if="isUnlocked">{{ appStore.lang === 'zh' ? '✦ 已解锁进入 iOS 桌面：点击底部按键或横条可重新锁屏 ✦' : '✦ Unlocked iOS Screen: Tap Relock or bar to test ✦' }}</text>
            <text class="tip-text" v-else-if="isPlayingAction">{{ appStore.lang === 'zh' ? '✦ 正在连贯播放 3D 动作帧序列 (30 FPS) ✦' : '✦ Playing 3D Motion Frame Sequence (30 FPS) ✦' }}</text>
            <text class="tip-text" v-else-if="isDragging">{{ appStore.lang === 'zh' ? '✦ 实时交互中：上下拖动感受 3D 空间视差 ✦' : '✦ Live Interaction: Drag up/down to feel 3D parallax ✦' }}</text>
            <text class="tip-text" v-else-if="isHolding">{{ appStore.lang === 'zh' ? '✦ 正在触发 Live 动效 / 触控交互 ✦' : '✦ Live Motion / Touch Active ✦' }}</text>
            <text class="tip-text" v-else>{{ appStore.lang === 'zh' ? '✦ 实机交互演示：按住上滑实时解锁 ｜ 移动感受 3D 视差 ✦' : '✦ Live Demo: Swipe up to unlock | Move to feel 3D parallax ✦' }}</text>
          </view>
        </view>

        <!-- Interactive 3D iPhone Body -->
        <view
          class="interactive-phone-frame"
          :class="{
            'aod-active': isAod,
            'is-pressing': isHolding,
            'is-unlocked': unlockProgress > 0.5,
            'is-sleeping': isSleeping,
            'is-dragging': isDragging
          }"
          :style="[ { aspectRatio: frameAspectRatio }, { transform: phoneContainerTransform } ]"
          @touchstart.stop.prevent="handleTouchStart"
          @touchmove.stop.prevent="handleTouchMove"
          @touchend.stop="handleTouchEnd"
          @touchcancel.stop="handleTouchEnd"
          @mousedown="handleMouseDown"
          @mouseup="handleMouseUp"
        >
          <!-- Layer 1: Background Layer(s) (CAML multi-layer or sequential action frame) -->
          <template v-if="hasTendiesLayers && camlBgLayers.length">
            <view
              v-for="layer in camlBgLayers"
              :key="layer.id"
              class="caml-sprite-layer caml-bg-layer"
              :style="getLayerStyle(layer)"
            >
              <view class="caml-anim-pos" :style="getPosAnimStyle(layer)">
                <view class="caml-anim-rot" :style="getRotAnimStyle(layer)">
                  <image
                    :src="layer.resolvedUrl"
                    mode="scaleToFill"
                    class="caml-sprite-inner-img"
                    :style="getImgAnimStyle(layer)"
                  />
                </view>
              </view>
            </view>
          </template>
          <image
            v-else-if="activeFrameUrl"
            :src="activeFrameUrl"
            mode="aspectFill"
            class="parallax-layer bg-layer"
            :class="{ 'bg-unlocked-zoom': unlockProgress > 0.5 }"
            :style="{ transform: bgLayerTransform }"
          />
          <view v-else class="parallax-layer bg-layer placeholder-bg"></view>

          <!-- Layer 1.5: Optional Long-Press Motion Video Stream -->
          <video
            v-if="motionAssetUrl"
            :src="motionAssetUrl"
            class="parallax-layer motion-video-layer"
            :class="{ active: isHolding || isPlayingAction }"
            :autoplay="true"
            :loop="true"
            :muted="true"
            :controls="false"
            :show-center-play-btn="false"
            :enable-progress-gesture="false"
            object-fit="cover"
          ></video>

          <!-- Dynamic Island -->
          <view class="dynamic-island" :class="{ 'island-aod': isAod, 'island-hidden': isSleeping }">
            <view class="camera-lens"></view>
          </view>

          <!-- Top iOS Status Bar (Horizontally aligned with Dynamic Island on the same centerline) -->
          <view class="phone-status-bar" :class="{ 'status-hidden': isSleeping }">
            <text class="status-time" :style="{ opacity: unlockProgress > 0.15 ? 1 : 0 }">{{ currentTimeShort }}</text>
            <view class="status-right-icons">
              <view class="ios-signal-bars">
                <view class="signal-bar bar-1"></view>
                <view class="signal-bar bar-2"></view>
                <view class="signal-bar bar-3"></view>
                <view class="signal-bar bar-4"></view>
              </view>
              <text class="ios-5g-text">5G</text>
              <view class="ios-battery-group">
                <text class="ios-battery-pct">63%</text>
                <view class="ios-battery-pill">
                  <view class="ios-battery-level" style="width: 63%"></view>
                  <view class="ios-battery-terminal"></view>
                </view>
              </view>
            </view>
          </view>

          <!-- Layer 2.5: Unlocked 9 SD Gundams Cutout Layer with Flight Path Trajectories -->
          <view
            v-if="hasVArchFlight"
            class="varch-flight-container"
            :style="flightContainerStyle"
          >
            <!-- Left Wing (Zeta, Unicorn, Strike, Exia) flying from left -->
            <image
              :src="vArchLeftUrl"
              mode="aspectFit"
              class="varch-wing varch-wing-left"
              :style="leftWingStyle"
            />
            <!-- Center Mini RX-78-2 flying from center -->
            <image
              :src="vArchCenterUrl"
              mode="aspectFit"
              class="varch-wing varch-wing-center"
              :style="centerWingStyle"
            />
            <!-- Right Wing (Freedom, Barbatos, Nu, Wing) flying from right -->
            <image
              :src="vArchRightUrl"
              mode="aspectFit"
              class="varch-wing varch-wing-right"
              :style="rightWingStyle"
            />
          </view>
          <image
            v-else-if="unlockedFgUrl"
            :src="unlockedFgUrl"
            mode="aspectFit"
            class="parallax-layer unlocked-varch-layer"
            :class="{ 'varch-active': unlockProgress > 0.05 }"
            :style="singleVArchStyle"
          />

          <!-- ============================================== -->
          <!-- Layer 2: iOS 17 Depth Clock & Date (Between BG and Floating/FG layers) -->
          <!-- ============================================== -->
          <view
            class="phone-lock-header interactive-clock-layer"
            :class="{ 'clock-fade-out': unlockProgress >= 0.99 }"
            :style="clockHeaderStyle"
          >
            <view class="lock-date-row">
              <text class="lock-icon" @tap.stop="toggleSleep">{{ unlockProgress > 0.5 ? '🔓' : '🔒' }}</text>
              <text class="lock-date-text">{{ currentDateStr }}</text>
            </view>
            <text class="lock-time-text" :class="{ 'clock-aod': isAod }">{{ currentTimeStr }}</text>
          </view>

          <!-- Layer 3: Foreground Floating Layer(s) (CAML multi-layer, Clock Ring, or Cutout) -->
          <template v-if="hasTendiesLayers && camlFloatingLayers.length">
            <view
              v-for="layer in camlFloatingLayers"
              :key="layer.id"
              class="caml-sprite-layer caml-floating-layer"
              :style="getLayerStyle(layer)"
            >
              <view class="caml-anim-pos" :style="getPosAnimStyle(layer)">
                <view class="caml-anim-rot" :style="getRotAnimStyle(layer)">
                  <image
                    :src="layer.resolvedUrl"
                    mode="scaleToFill"
                    class="caml-sprite-inner-img"
                    :style="getImgAnimStyle(layer)"
                  />
                </view>
              </view>
            </view>
          </template>
          <view
            v-else-if="fgAssetUrl && fgType === 'clock_ring'"
            class="clock-ring-wrapper"
            :class="{ 'ring-hidden': unlockProgress >= 0.99 }"
            :style="clockRingWrapperStyle"
          >
            <image
              :src="fgAssetUrl"
              mode="aspectFit"
              class="clock-ring-image"
            />
          </view>
          <image
            v-else-if="fgAssetUrl"
            :src="fgAssetUrl"
            mode="aspectFill"
            class="parallax-layer fg-cutout-layer"
            :style="{ transform: fgLayerTransform }"
          />

          <!-- Lock Screen Bottom Shortcuts & Swipe-Up Unlock Trigger -->
          <view class="phone-lock-footer" :class="{ 'footer-hidden': isAod }" :style="lockFooterStyle">
            <view class="footer-btn">
              <text class="footer-icon">🔦</text>
            </view>
            <view class="swipe-unlock-trigger" @tap.stop="triggerUnlockAnimation">
              <text class="swipe-arrow">↑</text>
              <text class="swipe-text">{{ isDragging ? (appStore.lang === 'zh' ? '滑动中...' : 'Dragging...') : (appStore.lang === 'zh' ? '向上轻扫解锁' : 'Swipe up to unlock') }}</text>
              <view class="home-indicator"></view>
            </view>
            <view class="footer-btn">
              <text class="footer-icon">📷</text>
            </view>
          </view>

          <!-- ============================================== -->
          <!-- STATE B: Unlocked iOS 17 Home Screen (Springboard) -->
          <!-- ============================================== -->
          <view
            v-show="unlockProgress > 0.005 || isUnlocked"
            class="homescreen-ui-wrapper"
            :class="{ 'ui-show': unlockProgress > 0.01, 'mode-clean': !showAppGrid }"
            :style="homescreenUiStyle"
          >
            <!-- Springboard App Grid (optional toggle, default false so wallpaper is uncluttered) -->
            <view v-if="showAppGrid" class="springboard-grid">
              <!-- Row 1 -->
              <view class="app-icon-item">
                <view class="app-icon icon-calendar">
                  <text class="cal-day">周二</text>
                  <text class="cal-num">20</text>
                </view>
                <text class="app-name">日历</text>
              </view>
              <view class="app-icon-item">
                <view class="app-icon icon-photos">
                  <text class="emoji-icon">🌸</text>
                </view>
                <text class="app-name">照片</text>
              </view>
              <view class="app-icon-item">
                <view class="app-icon icon-camera">
                  <text class="emoji-icon">📷</text>
                </view>
                <text class="app-name">相机</text>
              </view>
              <view class="app-icon-item">
                <view class="app-icon icon-clock">
                  <text class="emoji-icon">⏰</text>
                </view>
                <text class="app-name">时钟</text>
              </view>

              <!-- Row 2 -->
              <view class="app-icon-item">
                <view class="app-icon icon-weather">
                  <text class="emoji-icon">🌤️</text>
                </view>
                <text class="app-name">天气</text>
              </view>
              <view class="app-icon-item">
                <view class="app-icon icon-maps">
                  <text class="emoji-icon">🗺️</text>
                </view>
                <text class="app-name">地图</text>
              </view>
              <view class="app-icon-item">
                <view class="app-icon icon-notes">
                  <text class="emoji-icon">📝</text>
                </view>
                <text class="app-name">备忘录</text>
              </view>
              <view class="app-icon-item">
                <view class="app-icon icon-settings">
                  <text class="emoji-icon">⚙️</text>
                </view>
                <text class="app-name">设置</text>
              </view>

              <!-- Row 3 -->
              <view class="app-icon-item">
                <view class="app-icon icon-appstore">
                  <text class="appstore-letter">A</text>
                </view>
                <text class="app-name">App Store</text>
              </view>
              <view class="app-icon-item">
                <view class="app-icon icon-wechat">
                  <text class="emoji-icon">💬</text>
                </view>
                <text class="app-name">微信</text>
              </view>
              <view class="app-icon-item">
                <view class="app-icon icon-music">
                  <text class="emoji-icon">🎵</text>
                </view>
                <text class="app-name">音乐</text>
              </view>
              <view class="app-icon-item">
                <view class="app-icon icon-safari">
                  <text class="emoji-icon">🧭</text>
                </view>
                <text class="app-name">Safari</text>
              </view>
            </view>
            <view v-else class="clean-springboard-spacer"></view>

            <!-- Bottom Section: Search Pill & Real iOS Dock (Matching User's Image 2!) -->
            <view class="springboard-bottom-group" :style="springboardBottomStyle">
              <view class="ios-search-pill">
                <text class="search-glass-icon">🔍</text>
                <text class="search-text">搜索</text>
              </view>

              <view class="iphone-dock-bar">
                <view class="dock-app icon-phone">
                  <text class="dock-app-icon">📞</text>
                  <view class="app-badge"><text class="badge-num">3</text></view>
                </view>
                <view class="dock-app icon-chrome">
                  <text class="dock-app-icon">🌐</text>
                </view>
                <view class="dock-app icon-wechat">
                  <text class="dock-app-icon">💬</text>
                  <view class="app-badge"><text class="badge-num">1</text></view>
                </view>
                <view class="dock-app icon-camera">
                  <text class="dock-app-icon">📷</text>
                </view>
              </view>

              <!-- Unlocked Home Bar Indicator -->
              <view class="unlocked-home-bar" @tap.stop="triggerLockAnimation">
                <view class="home-indicator"></view>
              </view>
            </view>
          </view>

          <!-- Layer 3.5: Apple Screen Glass Glare / Specular Reflection -->
          <view class="phone-glass-glare" :style="glareStyle"></view>

          <!-- Layer 3.8: OLED True Black Sleep Mode Overlay (Tap to wake!) -->
          <view
            class="screen-sleep-overlay"
            :class="{ active: isSleeping }"
            @tap.stop="wakeScreen"
          >
            <view class="sleep-tip-content">
              <text class="sleep-glow-dot"></text>
              <text class="sleep-text">{{ appStore.lang === 'zh' ? '轻触屏幕点亮唤醒' : 'Tap to wake' }}</text>
            </view>
          </view>

          <!-- Layer 4: Anti-theft Preview Security Watermark (Only visible when locked, hidden when unlocked) -->
          <view
            v-if="!isUnlocked && unlockProgress <= 0.05"
            class="interactive-security-watermark"
          >
            <text class="wm-text">{{ appStore.lang === 'zh' ? '✨ 3D实机模拟演示' : '✨ 3D Live Simulation' }}</text>
          </view>
        </view>

        <!-- Interactive Control Dock Toolbar Under iPhone -->
        <view class="interactive-dock-toolbar">
          <!-- Unlock / Relock Action Button -->
          <view
            class="dock-tool-btn"
            :class="{ active: isUnlocked }"
            @tap="toggleLockState"
            :title="isUnlocked ? (appStore.lang === 'zh' ? '点击锁定屏幕' : 'Tap to lock screen') : (appStore.lang === 'zh' ? '模拟 iPhone 向上轻扫解锁' : 'Simulate iPhone swipe to unlock')"
          >
            <text class="tool-icon">{{ isUnlocked ? '🔒' : '🔓' }}</text>
            <text class="tool-text">{{ isUnlocked ? (appStore.lang === 'zh' ? '重新锁屏' : 'Lock') : (appStore.lang === 'zh' ? '上滑解锁' : 'Unlock') }}</text>
          </view>

          <!-- Toggle Desktop Icon Mode (When Unlocked) -->
          <view
            v-if="isUnlocked"
            class="dock-tool-btn"
            :class="{ active: showAppGrid }"
            @tap="toggleAppGrid"
            :title="appStore.lang === 'zh' ? '切换纯净桌面/应用图标桌面' : 'Toggle clean desktop / app icons'"
          >
            <text class="tool-icon">📱</text>
            <text class="tool-text">{{ showAppGrid ? (appStore.lang === 'zh' ? '纯净桌面' : 'Clean') : (appStore.lang === 'zh' ? '显示图标' : 'Icons') }}</text>
          </view>

          <!-- Action Frames Play / Pause Button (if sequence exists) -->
          <view
            v-if="totalFrames > 1"
            class="dock-tool-btn btn-action-play"
            :class="{ active: isPlayingAction }"
            @tap="togglePlayAction"
            :title="appStore.lang === 'zh' ? '播放 3D 动作连贯帧动画' : 'Play 3D motion sequence'"
          >
            <text class="tool-icon">{{ isPlayingAction ? '⏸️' : '🎬' }}</text>
            <text class="tool-text">{{ isPlayingAction ? (appStore.lang === 'zh' ? '停止动作' : 'Pause') : (appStore.lang === 'zh' ? '连贯动作' : 'Motion') }}</text>
          </view>

          <!-- OLED Sleep / Wake Button -->
          <view
            class="dock-tool-btn"
            :class="{ active: isSleeping }"
            @tap="toggleSleep"
            :title="appStore.lang === 'zh' ? '模拟 OLED 熄屏待机，轻触亮屏' : 'Simulate OLED display sleep, tap to wake'"
          >
            <text class="tool-icon">{{ isSleeping ? '☀️' : '🌙' }}</text>
            <text class="tool-text">{{ isSleeping ? (appStore.lang === 'zh' ? '唤醒亮屏' : 'Wake') : (appStore.lang === 'zh' ? '熄屏待机' : 'Sleep') }}</text>
          </view>

          <!-- Reset Tilt Button -->
          <view class="dock-tool-btn" @tap="resetTilt" :title="appStore.lang === 'zh' ? '复位 3D 视角' : 'Reset 3D view'">
            <text class="tool-icon">🔄</text>
            <text class="tool-text">{{ appStore.lang === 'zh' ? '重置视角' : 'Reset' }}</text>
          </view>
        </view>
      </view>

      <!-- MODE 1: Standard iPhone Lock Screen Simulator -->
      <view v-else-if="currentMode === 'phone'" class="phone-frame">
        <!-- Wallpaper Background Image -->
        <image
          v-if="resolvedImageUrl"
          :src="resolvedImageUrl"
          mode="aspectFill"
          class="mockup-img"
        />
        <view v-else class="mockup-img placeholder-bg"></view>

        <!-- Dynamic Island -->
        <view class="dynamic-island">
          <view class="camera-lens"></view>
        </view>

        <!-- Status Bar on Phone Frame (Aligned with Dynamic Island) -->
        <view class="phone-status-bar">
          <text class="status-time" style="opacity: 0">{{ currentTimeShort }}</text>
          <view class="status-right-icons">
            <view class="ios-signal-bars">
              <view class="signal-bar bar-1"></view>
              <view class="signal-bar bar-2"></view>
              <view class="signal-bar bar-3"></view>
              <view class="signal-bar bar-4"></view>
            </view>
            <text class="ios-5g-text">5G</text>
            <view class="ios-battery-group">
              <text class="ios-battery-pct">63%</text>
              <view class="ios-battery-pill">
                <view class="ios-battery-level" style="width: 63%"></view>
                <view class="ios-battery-terminal"></view>
              </view>
            </view>
          </view>
        </view>

        <!-- Lock Screen Header (Date & Big Live Clock) -->
        <view class="phone-lock-header">
          <view class="lock-date-row">
            <text class="lock-icon">🔒</text>
            <text class="lock-date-text">{{ currentDateStr }}</text>
          </view>
          <text class="lock-time-text">{{ currentTimeStr }}</text>
        </view>

        <!-- Lock Screen Bottom Shortcuts -->
        <view class="phone-lock-footer">
          <view class="footer-btn">
            <text class="footer-icon">🔦</text>
          </view>
          <view class="home-indicator"></view>
          <view class="footer-btn">
            <text class="footer-icon">📷</text>
          </view>
        </view>
      </view>

      <!-- MODE 2: MacBook Desktop Simulator -->
      <view v-else-if="currentMode === 'desktop'" class="desktop-frame">
        <!-- Top macOS Menu Bar -->
        <view class="mac-menu-bar">
          <view class="mac-left">
            <text class="mac-apple-logo">💻</text>
            <text class="mac-menu-item">{{ appStore.lang === 'zh' ? '访达' : 'Finder' }}</text>
            <text class="mac-menu-item">{{ appStore.lang === 'zh' ? '文件' : 'File' }}</text>
            <text class="mac-menu-item">{{ appStore.lang === 'zh' ? '显示' : 'View' }}</text>
          </view>
          <view class="mac-right">
            <text class="mac-status-icon">🔋 100%</text>
            <text class="mac-status-time">{{ currentTimeShort }}</text>
          </view>
        </view>

        <!-- Wallpaper in desktop ratio -->
        <image
          v-if="resolvedImageUrl"
          :src="resolvedImageUrl"
          mode="aspectFill"
          class="mockup-img mac-wallpaper"
        />
        <view v-else class="mockup-img mac-wallpaper placeholder-bg"></view>

        <!-- Bottom macOS Frosted Glass Dock -->
        <view class="mac-dock">
          <view class="dock-item"><text class="dock-emoji"> Finder </text></view>
          <view class="dock-item"><text class="dock-emoji">🧭</text></view>
          <view class="dock-item"><text class="dock-emoji">💬</text></view>
          <view class="dock-item"><text class="dock-emoji">🎵</text></view>
          <view class="dock-item"><text class="dock-emoji">⚙️</text></view>
        </view>
      </view>

      <!-- MODE 3: Clean Wallpaper Mode -->
      <view v-else class="clean-frame">
        <image
          v-if="resolvedImageUrl"
          :src="resolvedImageUrl"
          mode="aspectFill"
          class="mockup-img clean-wallpaper"
        />
        <view v-else class="mockup-img clean-wallpaper placeholder-bg"></view>
        <view class="clean-hint">
          <text class="hint-text">✦ 4K 超高清画质全屏原画 ✦</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useAppStore } from '@/store/app'
import { resolveAndDecryptAssetUrl } from '@/utils/assetCrypto'

const appStore = useAppStore()

const props = defineProps({
  imageUrl: {
    type: String,
    default: ''
  },
  product: {
    type: Object,
    default: () => ({})
  }
})

// Resolve absolute URL for backend-uploaded assets
function resolveMediaUrl(path) {
  if (!path) return ''
  if (path.startsWith('http://') || path.startsWith('https://')) return path
  // #ifdef H5
  if (typeof window !== 'undefined' && window.location) {
    const base = import.meta.env?.VITE_API_URL
    const origin = base ? base.replace(/\/api\/?$/, '') : window.location.origin
    return `${origin}${path.startsWith('/') ? path : '/' + path}`
  }
  const base = import.meta.env?.VITE_API_URL || 'https://wallpaper.icangoanyway.com/api'
  const origin = base.replace(/\/api\/?$/, '')
  return `${origin}${path.startsWith('/') ? path : '/' + path}`
  // #endif
  // #ifndef H5
  return `https://wallpaper.icangoanyway.com${path.startsWith('/') ? path : '/' + path}`
  // #endif
}

const resolvedImageUrl = computed(() => resolveMediaUrl(props.imageUrl))

// Strict condition: Whether interactive 3D simulation is possible depends on
// whether the product has is_interactive enabled AND uploaded interactive files in backend!
const interactiveConfig = computed(() => {
  const p = props.product
  if (!p || !p.is_interactive) return null
  if (!p.interactive_config) return null
  try {
    const cfg = typeof p.interactive_config === 'string' ? JSON.parse(p.interactive_config) : p.interactive_config
    const hasFiles = Boolean(
      (cfg.foreground_url && cfg.foreground_url.trim()) ||
      (cfg.background_url && cfg.background_url.trim()) ||
      (cfg.motion_url && cfg.motion_url.trim()) ||
      (cfg.tendies_file && cfg.tendies_file.trim()) ||
      (cfg.frames_dir && cfg.frames_dir.trim()) ||
      (Array.isArray(cfg.bg_layers) && cfg.bg_layers.length > 0) ||
      (Array.isArray(cfg.floating_layers) && cfg.floating_layers.length > 0)
    )
    if (hasFiles) {
      return cfg
    }
  } catch (e) {
    console.warn('Failed to parse interactive_config', e)
  }
  return null
})

const hasInteractive = computed(() => !!interactiveConfig.value)

// Default to iPhone 14 Pro / 15 / 16 (393 x 852, 19.5:9) viewport
const canvasW = computed(() => Number(interactiveConfig.value?.canvas_width) || 393)
const canvasH = computed(() => Number(interactiveConfig.value?.canvas_height) || 852)
const frameAspectRatio = computed(() => `${canvasW.value} / ${canvasH.value}`)

const hasTendiesLayers = computed(() => {
  const cfg = interactiveConfig.value
  if (!cfg) return false
  const bgCount = Array.isArray(cfg.bg_layers) ? cfg.bg_layers.length : 0
  const flCount = Array.isArray(cfg.floating_layers) ? cfg.floating_layers.length : 0
  return (bgCount + flCount) > 0
})

function isLayerBackplate(layer) {
  if (!layer) return false
  const cW = canvasW.value || 393
  const cH = canvasH.value || 852
  const w = layer.width || 0
  const h = layer.height || 0
  const name = (layer.name || '').toLowerCase()
  const isSizeFull = (w >= cW * 0.85) && (h >= cH * 0.85)
  const isNameBg = name.includes('back') || name.includes('bg')
  return isSizeFull || (isNameBg && w >= cW * 0.5)
}

const decryptedLayerUrls = ref({})

// Prefetch & decrypt any encrypted .dat assets in memory
watch(() => interactiveConfig.value, (cfg) => {
  if (!cfg) return
  const allUrls = [
    ...(Array.isArray(cfg.bg_layers) ? cfg.bg_layers.map(l => l.url) : []),
    ...(Array.isArray(cfg.floating_layers) ? cfg.floating_layers.map(l => l.url) : []),
    cfg.background_url,
    cfg.foreground_url
  ].filter(Boolean)

  allUrls.forEach(url => {
    if (url.includes('.dat')) {
      const fullUrl = resolveMediaUrl(url)
      if (!decryptedLayerUrls.value[url]) {
        resolveAndDecryptAssetUrl(fullUrl).then(blobUrl => {
          decryptedLayerUrls.value[url] = blobUrl
        })
      }
    }
  })
}, { immediate: true })

const camlBgLayers = computed(() => {
  const cfg = interactiveConfig.value
  if (!cfg) return []
  const bg = Array.isArray(cfg.bg_layers) ? cfg.bg_layers : []
  const fl = Array.isArray(cfg.floating_layers) ? cfg.floating_layers : []
  // Route any backplates mistakenly stored in floating_layers to bg_layers
  const flBg = fl.filter(l => isLayerBackplate(l)).map(l => ({ ...l, category: 'background' }))
  return [...bg, ...flBg].map(l => ({
    ...l,
    resolvedUrl: decryptedLayerUrls.value[l.url] || resolveMediaUrl(l.url)
  }))
})

const camlFloatingLayers = computed(() => {
  const cfg = interactiveConfig.value
  if (!cfg || !Array.isArray(cfg.floating_layers)) return []
  // Only keep true floating cutouts, never full screen backgrounds
  return cfg.floating_layers.filter(l => !isLayerBackplate(l)).map(l => ({
    ...l,
    resolvedUrl: decryptedLayerUrls.value[l.url] || resolveMediaUrl(l.url)
  }))
})

function getLayerStyle(layer) {
  const cW = canvasW.value || 393
  const cH = canvasH.value || 852
  const p = unlockProgress.value
  const depthInt = interactiveConfig.value?.depth_intensity !== undefined && interactiveConfig.value?.depth_intensity !== null
    ? Number(interactiveConfig.value.depth_intensity)
    : 16

  // 1. Base locked state values
  const lockState = layer.states?.locked

  // Position animations in CAML specify absolute keyframes.
  // tendies_parser computed DeltaX = (rawVal - layer.pos_x) and DeltaY = (screenY - layer.pos_y).
  // Therefore, translate3d(dx, dy, 0) is relative to layer.pos_x / layer.pos_y.
  // For axes with animations, we use layer base pos to avoid double offset (which caused the ball overflow).
  // For other properties and non-animated layers, strictly respect lockState (e.g. wtksoi offscreen/hidden).
  const hasPosXAnim = Array.isArray(layer.animations) && layer.animations.some(a => a.key_path === 'position.x' || a.key_path === 'position')
  const hasPosYAnim = Array.isArray(layer.animations) && layer.animations.some(a => a.key_path === 'position.y' || a.key_path === 'position')
  const hasRotAnim = Array.isArray(layer.animations) && layer.animations.some(a => a.key_path === 'transform.rotation.z')

  let lockX = (lockState && lockState.pos_x !== undefined) ? lockState.pos_x : layer.pos_x
  let lockY = (lockState && lockState.pos_y !== undefined) ? lockState.pos_y : layer.pos_y
  let lockW = (lockState && lockState.width !== undefined && lockState.width > 0) ? lockState.width : layer.width
  let lockH = (lockState && lockState.height !== undefined && lockState.height > 0) ? lockState.height : layer.height
  let lockOp = (lockState && lockState.opacity !== undefined) ? lockState.opacity : (layer.opacity !== undefined ? layer.opacity : 1.0)
  let lockRot = (!hasRotAnim && lockState && lockState.rotation !== undefined) ? lockState.rotation : (layer.rotation || 0)
  let lockScale = (lockState && lockState.scale !== undefined && lockState.scale > 0) ? lockState.scale : (layer.scale || 1.0)

  // 2. Base unlock state values
  const unlockState = layer.states?.unlock
  const unlockX = (unlockState && unlockState.pos_x !== undefined) ? unlockState.pos_x : lockX
  const unlockY = (unlockState && unlockState.pos_y !== undefined) ? unlockState.pos_y : lockY
  const unlockW = (unlockState && unlockState.width !== undefined && unlockState.width > 0) ? unlockState.width : lockW
  const unlockH = (unlockState && unlockState.height !== undefined && unlockState.height > 0) ? unlockState.height : lockH
  const unlockOp = (unlockState && unlockState.opacity !== undefined) ? unlockState.opacity : lockOp
  const unlockRot = (unlockState && unlockState.rotation !== undefined) ? unlockState.rotation : lockRot
  const unlockScale = (unlockState && unlockState.scale !== undefined && unlockState.scale > 0) ? unlockState.scale : lockScale

  // 3. Interpolated state
  let curX = lockX + (unlockX - lockX) * p
  let curY = lockY + (unlockY - lockY) * p
  let curW = lockW + (unlockW - lockW) * p
  let curH = lockH + (unlockH - lockH) * p
  let curOp = lockOp + (unlockOp - lockOp) * p
  let curRot = lockRot + (unlockRot - lockRot) * p
  let curScale = lockScale + (unlockScale - lockScale) * p

  // 4. Sleep state override if sleeping
  if (isSleeping.value) {
    const sleepState = layer.states?.sleep
    if (sleepState) {
      if (sleepState.pos_x !== undefined) curX = sleepState.pos_x
      if (sleepState.pos_y !== undefined) curY = sleepState.pos_y
      if (sleepState.width !== undefined && sleepState.width > 0) curW = sleepState.width
      if (sleepState.height !== undefined && sleepState.height > 0) curH = sleepState.height
      if (sleepState.opacity !== undefined) curOp = sleepState.opacity
      if (sleepState.rotation !== undefined) curRot = sleepState.rotation
      if (sleepState.scale !== undefined && sleepState.scale > 0) curScale = sleepState.scale
    } else {
      curOp = curOp * 0.35
    }
  }

  // 5. Parallax offset
  let parallaxX = 0
  let parallaxY = 0
  if (isGyroActive.value && !isSleeping.value && depthInt !== 0) {
    if (layer.category === 'background') {
      const factor = (layer.order + 1) * 0.12
      parallaxX = -mousePos.value.x * depthInt * factor
      parallaxY = -mousePos.value.y * depthInt * factor
    } else {
      const factor = 0.5 + (layer.order + 1) * 0.18
      parallaxX = mousePos.value.x * depthInt * factor
      parallaxY = mousePos.value.y * depthInt * factor
    }
  }

  // 6. Check if full-canvas layer (avoids sub-pixel seam/displacement for full-frame layers)
  const isFullCanvas = (curW >= cW * 0.96) && (curH >= cH * 0.96) && Math.abs(curX - cW / 2) < 4 && Math.abs(curY - cH / 2) < 4 && Math.abs(curRot) < 0.001 && Math.abs(curScale - 1.0) < 0.001
  const isBackplate = isFullCanvas || isLayerBackplate(layer) || layer.category === 'background'

  // 7. Coordinates to percentages
  const leftPct = (curX / cW) * 100
  const topPct = (curY / cH) * 100
  const widthPct = (curW / cW) * 100
  const heightPct = (curH / cH) * 100

  // 8. Z-Index and 3D depth calculation:
  // Background: 10 + Math.min(layer.order || 0, 20) (strictly behind clock, Z: 5px)
  // Clock: 50 (strictly in front of background, Z: 30px)
  // Floating: 60 + Math.min(layer.order || 0, 15) (in front of clock for authentic iOS 17 depth effect, Z: 38px)
  // Foreground: 78 + Math.min(layer.order || 0, 15) (Z: 50px)
  let zIndex = 10 + Math.min(layer.order || 0, 20)
  let zTranslate = 5 + Math.min(layer.order || 0, 10)
  if (!isBackplate) {
    if (layer.category === 'floating') {
      zIndex = 60 + Math.min(layer.order || 0, 15)
      zTranslate = 38 + Math.min(layer.order || 0, 10)
    } else if (layer.category === 'foreground') {
      zIndex = 78 + Math.min(layer.order || 0, 15)
      zTranslate = 50 + Math.min(layer.order || 0, 10)
    }
  }

  const rotDeg = curRot * (180 / Math.PI)
  const rotStr = Math.abs(rotDeg) > 0.01 ? ` rotate(${rotDeg.toFixed(2)}deg)` : ''
  const scaleStr = Math.abs(curScale - 1.0) > 0.001 ? ` scale(${curScale.toFixed(3)})` : ''

  if (isFullCanvas) {
    return {
      position: 'absolute',
      left: '0',
      top: '0',
      width: '100%',
      height: '100%',
      transform: `translate3d(${parallaxX.toFixed(2)}px, ${parallaxY.toFixed(2)}px, ${zTranslate}px)`,
      transformOrigin: '50% 50%',
      opacity: Math.max(0, Math.min(1, curOp)),
      visibility: curOp < 0.005 ? 'hidden' : 'visible',
      zIndex: zIndex,
      pointerEvents: 'none',
      transition: isDragging.value ? 'none' : 'transform 0.08s ease-out, opacity 0.3s ease'
    }
  }

  return {
    position: 'absolute',
    left: `${leftPct}%`,
    top: `${topPct}%`,
    width: `${widthPct}%`,
    height: `${heightPct}%`,
    transform: `translate(-50%, -50%) translate3d(${parallaxX.toFixed(2)}px, ${parallaxY.toFixed(2)}px, ${zTranslate}px)${scaleStr}${rotStr}`,
    transformOrigin: '50% 50%',
    opacity: Math.max(0, Math.min(1, curOp)),
    visibility: curOp < 0.005 ? 'hidden' : 'visible',
    zIndex: zIndex,
    pointerEvents: 'none',
    transition: isDragging.value ? 'none' : 'transform 0.08s ease-out, opacity 0.3s ease'
  }
}

function getPosAnimStyle(layer) {
  if (!layer.animations || !layer.animations.length) return { width: '100%', height: '100%' }
  const anim = layer.animations.find(a => a.key_path === 'position.y' || a.key_path === 'position.x' || a.key_path === 'position')
  if (!anim) return { width: '100%', height: '100%' }
  const dir = anim.autoreverses ? 'alternate' : 'normal'
  const timing = anim.key_path === 'position' && !anim.autoreverses ? 'linear' : 'ease-in-out'
  return {
    animation: `ca-pos-${layer.id} ${anim.duration || 5}s ${timing} infinite ${dir}`,
    width: '100%',
    height: '100%'
  }
}

function getRotAnimStyle(layer) {
  if (!layer.animations || !layer.animations.length) return { width: '100%', height: '100%' }
  const anim = layer.animations.find(a => a.key_path === 'transform.rotation.z')
  if (!anim) return { width: '100%', height: '100%' }
  const dir = anim.autoreverses ? 'alternate' : 'normal'
  return {
    animation: `ca-rot-${layer.id} ${anim.duration || 5}s ease-in-out infinite ${dir}`,
    width: '100%',
    height: '100%'
  }
}

function getImgAnimStyle(layer) {
  if (!layer.animations || !layer.animations.length) return { width: '100%', height: '100%' }
  const anim = layer.animations.find(a => a.key_path === 'opacity')
  if (!anim) return { width: '100%', height: '100%' }
  const dir = anim.autoreverses ? 'alternate' : 'normal'
  return {
    animation: `ca-op-${layer.id} ${anim.duration || 5}s ease-in-out infinite ${dir}`,
    width: '100%',
    height: '100%'
  }
}

const generatedKeyframesCSS = computed(() => {
  let css = ''
  const allLayers = [...camlBgLayers.value, ...camlFloatingLayers.value]
  for (const layer of allLayers) {
    if (!layer.animations || !layer.animations.length) continue
    for (const anim of layer.animations) {
      if (anim.key_path === 'position.y' || anim.key_path === 'position.x' || anim.key_path === 'position') {
        const count = Math.max(anim.delta_x?.length || 0, anim.delta_y?.length || 0)
        if (count > 0) {
          let midX = 0
          if (anim.delta_x && anim.delta_x.length > 0) {
            const minX = Math.min(...anim.delta_x)
            const maxX = Math.max(...anim.delta_x)
            midX = (minX + maxX) / 2
          }
          let midY = 0
          if (anim.delta_y && anim.delta_y.length > 0) {
            const minY = Math.min(...anim.delta_y)
            const maxY = Math.max(...anim.delta_y)
            midY = (minY + maxY) / 2
          }
          css += `@keyframes ca-pos-${layer.id} {\n`
          for (let i = 0; i < count; i++) {
            const pct = count === 1 ? 100 : Math.round((i / (count - 1)) * 100)
            const rawDx = anim.delta_x && anim.delta_x[i] !== undefined ? anim.delta_x[i] : 0
            const rawDy = anim.delta_y && anim.delta_y[i] !== undefined ? anim.delta_y[i] : 0
            const dx = rawDx - midX
            const dy = rawDy - midY
            css += `  ${pct}% { transform: translate3d(${dx.toFixed(2)}px, ${dy.toFixed(2)}px, 0); }\n`
          }
          css += `}\n`
        }
      } else if (anim.key_path === 'transform.rotation.z') {
        const count = anim.rot_deg?.length || 0
        if (count > 0) {
          css += `@keyframes ca-rot-${layer.id} {\n`
          for (let i = 0; i < count; i++) {
            const pct = count === 1 ? 100 : Math.round((i / (count - 1)) * 100)
            const deg = anim.rot_deg[i] || 0
            css += `  ${pct}% { transform: rotate(${deg.toFixed(2)}deg); }\n`
          }
          css += `}\n`
        }
      } else if (anim.key_path === 'opacity') {
        const count = anim.opacity?.length || 0
        if (count > 0) {
          css += `@keyframes ca-op-${layer.id} {\n`
          for (let i = 0; i < count; i++) {
            const pct = count === 1 ? 100 : Math.round((i / (count - 1)) * 100)
            const op = anim.opacity[i] || 0
            css += `  ${pct}% { opacity: ${op.toFixed(3)}; }\n`
          }
          css += `}\n`
        }
      }
    }
  }
  return css
})

watch(generatedKeyframesCSS, (newCSS) => {
  // #ifdef H5
  if (typeof document !== 'undefined') {
    let styleEl = document.getElementById('caml-keyframes-style')
    if (!styleEl) {
      styleEl = document.createElement('style')
      styleEl.id = 'caml-keyframes-style'
      document.head.appendChild(styleEl)
    }
    styleEl.textContent = newCSS
  }
  // #endif
}, { immediate: true })

onUnmounted(() => {
  // #ifdef H5
  if (typeof document !== 'undefined') {
    const el = document.getElementById('caml-keyframes-style')
    if (el) el.remove()
  }
  // #endif
})

const bgAssetUrl = computed(() => {
  const cfg = interactiveConfig.value
  const url = cfg?.background_url || props.imageUrl
  if (!url) return ''
  return decryptedLayerUrls.value[url] || resolveMediaUrl(url)
})

const fgAssetUrl = computed(() => {
  const cfg = interactiveConfig.value
  if (!cfg?.foreground_url) return ''
  return decryptedLayerUrls.value[cfg.foreground_url] || resolveMediaUrl(cfg.foreground_url)
})

const motionAssetUrl = computed(() => {
  const cfg = interactiveConfig.value
  if (cfg?.motion_url) return resolveMediaUrl(cfg.motion_url)
  return ''
})

const unlockedFgUrl = computed(() => {
  const cfg = interactiveConfig.value
  if (cfg?.unlocked_fg_url) return resolveMediaUrl(cfg.unlocked_fg_url)
  return ''
})

const vArchLeftUrl = computed(() => resolveMediaUrl('/uploads/interactive/floating_assets/v_arch_left.png'))
const vArchRightUrl = computed(() => resolveMediaUrl('/uploads/interactive/floating_assets/v_arch_right.png'))
const vArchCenterUrl = computed(() => resolveMediaUrl('/uploads/interactive/floating_assets/v_arch_center.png'))
const hasVArchFlight = computed(() => {
  if (hasTendiesLayers.value) return false
  return Boolean(unlockedFgUrl.value || fgType.value === 'clock_ring')
})

const fgType = computed(() => {
  const cfg = interactiveConfig.value
  if (cfg?.fg_type) return cfg.fg_type
  if (fgAssetUrl.value && (
    fgAssetUrl.value.includes('xuanzhuan') ||
    fgAssetUrl.value.includes('ring') ||
    fgAssetUrl.value.includes('halo') ||
    fgAssetUrl.value.includes('clock')
  )) {
    return 'clock_ring'
  }
  return 'subject'
})

const showAppGrid = ref(false) // 默认纯净桌面，展现无遮挡原生壁纸与战队拱形动效

function toggleAppGrid() {
  showAppGrid.value = !showAppGrid.value
}

// Sequential Action Frames (e.g. Gundam 0-94 frames)
const totalFrames = computed(() => Number(interactiveConfig.value?.frame_count) || 0)
const framesDir = computed(() => interactiveConfig.value?.frames_dir ? resolveMediaUrl(interactiveConfig.value.frames_dir) : '')
const framePattern = computed(() => interactiveConfig.value?.frame_pattern || 'gundam_%d.jpg')
const currentFrame = ref(0)
const isPlayingAction = ref(false)

const activeFrameUrl = computed(() => {
  if (framesDir.value && totalFrames.value > 1) {
    const filename = framePattern.value.replace('%d', currentFrame.value)
    return `${framesDir.value}/${filename}`
  }
  return bgAssetUrl.value
})

const currentMode = ref('phone') // 'interactive' | 'phone' | 'desktop' | 'clean'
const currentTimeStr = ref('09:41')
const currentTimeShort = ref('9:41')
const currentDateStr = ref('10月20日 星期二')

const mousePos = ref({ x: 0, y: 0 })
const isHolding = ref(false)
const isAod = ref(false)
const isSleeping = ref(false)
const isUnlocked = ref(false)
const unlockProgress = ref(0) // 0 (locked) -> 1 (unlocked)

// Touch/drag gesture state
const dragStartY = ref(0)
const dragStartProgress = ref(0)
const isDragging = ref(false)

// Auto switch to interactive mode when interactive assets are present
watch(hasInteractive, (val) => {
  if (val) {
    currentMode.value = 'interactive'
  } else if (currentMode.value === 'interactive') {
    currentMode.value = 'phone'
  }
}, { immediate: true })

let timer = null
let actionTimer = null
let animFrameId = null

function updateClock() {
  const now = new Date()
  const hours = String(now.getHours()).padStart(2, '0')
  const minutes = String(now.getMinutes()).padStart(2, '0')
  currentTimeStr.value = `${hours}:${minutes}`
  currentTimeShort.value = `${now.getHours()}:${minutes}`

  const weekdaysZh = ['星期日', '星期一', '星期二', '星期三', '星期四', '星期五', '星期六']
  const weekdaysEn = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']
  const monthsEn = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']
  const month = now.getMonth() + 1
  const date = now.getDate()
  if (appStore.lang === 'zh') {
    currentDateStr.value = `${month}月${date}日 ${weekdaysZh[now.getDay()]}`
  } else {
    currentDateStr.value = `${weekdaysEn[now.getDay()]}, ${monthsEn[now.getMonth()]} ${date}`
  }
}

watch(() => appStore.lang, () => {
  updateClock()
})

function setMode(mode) {
  currentMode.value = mode
}

function toggleAOD() {
  isAod.value = !isAod.value
  if (isAod.value) {
    unlockProgress.value = 0
    isUnlocked.value = false
    isSleeping.value = false
  }
}

function animateProgressTo(targetValue, duration = 440, onDone = null) {
  if (animFrameId) {
    cancelAnimationFrame(animFrameId)
    animFrameId = null
  }
  const startVal = unlockProgress.value
  const diff = targetValue - startVal
  if (Math.abs(diff) < 0.002) {
    unlockProgress.value = targetValue
    isUnlocked.value = targetValue > 0.5
    if (onDone) onDone()
    return
  }
  const startTime = (typeof performance !== 'undefined' ? performance.now() : Date.now())
  function step(now) {
    const currentTime = now || (typeof performance !== 'undefined' ? performance.now() : Date.now())
    const elapsed = currentTime - startTime
    const t = Math.min(1, elapsed / duration)
    // Smooth cubic ease out
    const ease = 1 - Math.pow(1 - t, 3)
    unlockProgress.value = Math.max(0, Math.min(1, startVal + diff * ease))

    // Synchronize Gundam rifle action frame during animation
    if (!isPlayingAction.value && totalFrames.value > 1) {
      currentFrame.value = Math.min(totalFrames.value - 1, Math.round(unlockProgress.value * 48))
    }

    if (t < 1) {
      animFrameId = requestAnimationFrame(step)
    } else {
      unlockProgress.value = targetValue
      isUnlocked.value = targetValue > 0.5
      animFrameId = null
      if (onDone) onDone()
    }
  }
  animFrameId = requestAnimationFrame(step)
}

function triggerUnlockAnimation() {
  if (isSleeping.value) {
    isSleeping.value = false
    return
  }
  animateProgressTo(1, 460, () => {
    isUnlocked.value = true
  })
}

function triggerLockAnimation() {
  animateProgressTo(0, 420, () => {
    isUnlocked.value = false
  })
}

function toggleLockState() {
  if (isSleeping.value) {
    isSleeping.value = false
    return
  }
  if (isUnlocked.value || unlockProgress.value > 0.5) {
    triggerLockAnimation()
  } else {
    triggerUnlockAnimation()
  }
}

function unlockPhone() {
  triggerUnlockAnimation()
}

function lockPhone() {
  triggerLockAnimation()
}

function toggleSleep() {
  isSleeping.value = !isSleeping.value
  if (isSleeping.value) {
    unlockProgress.value = 0
    isUnlocked.value = false
  }
}

function wakeScreen() {
  isSleeping.value = false
}

function togglePlayAction() {
  isPlayingAction.value = !isPlayingAction.value
  if (isPlayingAction.value) {
    if (actionTimer) clearInterval(actionTimer)
    actionTimer = setInterval(() => {
      currentFrame.value = (currentFrame.value + 1) % Math.max(1, totalFrames.value)
    }, 1000 / 30)
  } else {
    if (actionTimer) clearInterval(actionTimer)
    actionTimer = null
  }
}

function resetTilt() {
  mousePos.value = { x: 0, y: 0 }
  currentFrame.value = 0
}

// Parallax calculations
const depthIntensity = computed(() => {
  if (interactiveConfig.value?.depth_intensity === undefined || interactiveConfig.value?.depth_intensity === null) {
    return 16
  }
  const int = Number(interactiveConfig.value.depth_intensity)
  return isNaN(int) ? 16 : int
})

const isGyroActive = computed(() => {
  return interactiveConfig.value?.gyro_enabled !== false
})

const bgLayerTransform = computed(() => {
  if (!isGyroActive.value || depthIntensity.value === 0) return 'none'
  const int = depthIntensity.value
  const mx = -mousePos.value.x * (int * 0.7)
  const my = -mousePos.value.y * (int * 0.7)
  return `scale(1.05) translate3d(${mx.toFixed(1)}px, ${my.toFixed(1)}px, 0)`
})

const fgLayerTransform = computed(() => {
  if (!isGyroActive.value || depthIntensity.value === 0) return 'none'
  const int = depthIntensity.value
  const mx = mousePos.value.x * (int * 1.5)
  const my = mousePos.value.y * (int * 1.5)
  return `scale(1.05) translate3d(${mx.toFixed(1)}px, ${my.toFixed(1)}px, 0)`
})

// V-Arch 9 SD Gundams Flight Trajectory Styles
const flightContainerStyle = computed(() => {
  const p = unlockProgress.value
  const int = isGyroActive.value ? depthIntensity.value : 0
  const mx = mousePos.value.x * (int * 0.5)
  const my = mousePos.value.y * (int * 0.5)
  return {
    opacity: p > 0.01 ? 1 : 0,
    pointerEvents: 'none',
    transform: `translate3d(${mx.toFixed(1)}px, ${my.toFixed(1)}px, 0)`,
    transition: isDragging.value ? 'none' : 'opacity 0.3s ease'
  }
})

// Left wing: Zeta, Unicorn, Strike, Exia fly in from left offscreen
const leftWingStyle = computed(() => {
  const p = unlockProgress.value
  const translateX = (1 - p) * -130 // in px
  const translateY = (1 - p) * 45   // in px
  const scale = 0.65 + p * 0.35     // 0.65 -> 1.0
  const opacity = Math.min(1, p * 1.6)
  return {
    transform: `translate3d(${translateX.toFixed(1)}px, ${translateY.toFixed(1)}px, 35px) scale(${scale.toFixed(3)})`,
    opacity: opacity.toFixed(3),
    transition: isDragging.value ? 'none' : 'transform 0.45s cubic-bezier(0.16, 1, 0.3, 1), opacity 0.35s ease'
  }
})

// Right wing: Freedom, Barbatos, Nu, Wing fly in from right offscreen
const rightWingStyle = computed(() => {
  const p = unlockProgress.value
  const translateX = (1 - p) * 130 // in px
  const translateY = (1 - p) * 45  // in px
  const scale = 0.65 + p * 0.35    // 0.65 -> 1.0
  const opacity = Math.min(1, p * 1.6)
  return {
    transform: `translate3d(${translateX.toFixed(1)}px, ${translateY.toFixed(1)}px, 35px) scale(${scale.toFixed(3)})`,
    opacity: opacity.toFixed(3),
    transition: isDragging.value ? 'none' : 'transform 0.45s cubic-bezier(0.16, 1, 0.3, 1), opacity 0.35s ease'
  }
})

// Center mini RX-78-2: swoops up and zooms from center
const centerWingStyle = computed(() => {
  const p = unlockProgress.value
  const translateY = (1 - p) * 60 // in px
  const scale = 0.35 + p * 0.65    // 0.35 -> 1.0
  const opacity = Math.min(1, p * 2.2)
  return {
    transform: `translate3d(0, ${translateY.toFixed(1)}px, 35px) scale(${scale.toFixed(3)})`,
    opacity: opacity.toFixed(3),
    transition: isDragging.value ? 'none' : 'transform 0.45s cubic-bezier(0.16, 1, 0.3, 1), opacity 0.35s ease'
  }
})

const singleVArchStyle = computed(() => {
  const p = unlockProgress.value
  const translateY = (1 - p) * 35 // in px
  const scale = 0.8 + p * 0.2
  return {
    transform: `translate3d(0, ${translateY.toFixed(1)}px, 35px) scale(${scale.toFixed(3)})`,
    opacity: p.toFixed(3),
    transition: isDragging.value ? 'none' : 'transform 0.45s cubic-bezier(0.16, 1, 0.3, 1), opacity 0.35s ease'
  }
})

// Rotating Clock Halo: shoots upwards and fades out on unlock
const clockRingWrapperStyle = computed(() => {
  const p = unlockProgress.value
  const int = isGyroActive.value ? depthIntensity.value : 0
  const mx = mousePos.value.x * (int * 0.4)
  const my = mousePos.value.y * (int * 0.4)
  const translateY = -p * 140
  const scale = 1 + p * 0.15
  const opacity = Math.max(0, 1 - p * 1.8)
  return {
    transform: `translate3d(${mx.toFixed(1)}px, ${(my + translateY).toFixed(1)}px, 40px) scale(${scale.toFixed(3)})`,
    opacity: opacity.toFixed(3),
    transition: isDragging.value ? 'none' : 'transform 0.45s cubic-bezier(0.16, 1, 0.3, 1), opacity 0.35s ease'
  }
})

// Depth Clock & Date: slides up and fades out (clean px coordinates and positive 30px Z depth)
const clockHeaderStyle = computed(() => {
  const p = unlockProgress.value
  const int = isGyroActive.value ? depthIntensity.value : 0
  const mx = -mousePos.value.x * (int * 0.25)
  const my = -mousePos.value.y * (int * 0.25)
  const translateY = -p * 110 // pixel translation
  const opacity = Math.max(0, 1 - p * 1.8)
  return {
    transform: `translate3d(${mx.toFixed(1)}px, ${(my + translateY).toFixed(1)}px, 30px)`,
    opacity: opacity.toFixed(3),
    transition: isDragging.value ? 'none' : 'transform 0.4s cubic-bezier(0.16, 1, 0.3, 1), opacity 0.3s ease',
    zIndex: 50
  }
})

// Lockscreen Footer (Flashlight, Camera, swipe bar)
const lockFooterStyle = computed(() => {
  const p = unlockProgress.value
  const translateY = p * 30
  const opacity = Math.max(0, 1 - p * 2.2)
  return {
    transform: `translate3d(0, ${translateY.toFixed(1)}px, 60px)`,
    opacity: opacity.toFixed(3),
    transition: isDragging.value ? 'none' : 'transform 0.35s ease, opacity 0.25s ease'
  }
})

// Lockscreen UI Wrapper
const lockscreenUiStyle = computed(() => {
  const p = unlockProgress.value
  return {
    pointerEvents: p >= 0.85 ? 'none' : 'auto'
  }
})

// Homescreen UI Wrapper
const homescreenUiStyle = computed(() => {
  const p = unlockProgress.value
  const isHidden = p <= 0.005
  const opacity = Math.max(0, (p - 0.12) / 0.88)
  const scale = 0.92 + p * 0.08
  return {
    opacity: isHidden ? '0' : opacity.toFixed(3),
    visibility: isHidden ? 'hidden' : 'visible',
    transform: `translate3d(0, 0, ${isHidden ? 0 : 100}px) scale(${scale.toFixed(3)})`,
    pointerEvents: p >= 0.8 ? 'auto' : 'none',
    zIndex: isHidden ? -1 : 200,
    transition: isDragging.value ? 'none' : 'opacity 0.4s ease, transform 0.4s cubic-bezier(0.16, 1, 0.3, 1)'
  }
})

// Homescreen Top Status Bar
const homeStatusStyle = computed(() => {
  const p = unlockProgress.value
  const isHidden = p <= 0.005
  const translateY = (1 - p) * -15
  const opacity = Math.max(0, (p - 0.25) / 0.75)
  return {
    transform: `translate3d(0, ${translateY.toFixed(1)}px, ${isHidden ? 0 : 150}px)`,
    opacity: isHidden ? '0' : opacity.toFixed(3),
    visibility: isHidden ? 'hidden' : 'visible',
    zIndex: isHidden ? -1 : 260,
    transition: isDragging.value ? 'none' : 'transform 0.4s ease, opacity 0.3s ease'
  }
})

// Springboard Bottom Group (Search Pill + Dock)
const springboardBottomStyle = computed(() => {
  const p = unlockProgress.value
  const isHidden = p <= 0.005
  // Translate 320px down so dock is completely below the bottom edge of the phone frame
  const translateY = (1 - p) * 320
  const opacity = Math.max(0, (p - 0.18) / 0.82)
  return {
    transform: `translate3d(0, ${translateY.toFixed(1)}px, ${isHidden ? 0 : 30}px)`,
    opacity: isHidden ? '0' : opacity.toFixed(3),
    visibility: isHidden ? 'hidden' : 'visible',
    zIndex: isHidden ? -1 : 230,
    transition: isDragging.value ? 'none' : 'transform 0.45s cubic-bezier(0.16, 1, 0.3, 1), opacity 0.35s ease'
  }
})

const phoneContainerTransform = computed(() => {
  if (!isGyroActive.value) {
    return isHolding.value ? 'scale(0.96)' : 'none'
  }
  const rotX = (-mousePos.value.y * 18).toFixed(1)
  const rotY = (mousePos.value.x * 18).toFixed(1)
  const scale = isHolding.value ? 0.96 : 1.0
  return `perspective(1000px) rotateX(${rotX}deg) rotateY(${rotY}deg) scale(${scale})`
})

const glareStyle = computed(() => {
  if (!isGyroActive.value) return { display: 'none' }
  const px = ((mousePos.value.x + 0.5) * 100).toFixed(0)
  const py = ((mousePos.value.y + 0.5) * 100).toFixed(0)
  return {
    background: `radial-gradient(circle at ${px}% ${py}%, rgba(255, 255, 255, 0.25) 0%, rgba(255, 255, 255, 0.05) 40%, transparent 70%)`
  }
})

// Update mouse position for gyro parallax
function updateMousePosition(dx, dy) {
  if (!isGyroActive.value) {
    mousePos.value = { x: 0, y: 0 }
    return
  }
  mousePos.value = {
    x: Math.max(-0.6, Math.min(0.6, dx)),
    y: Math.max(-0.6, Math.min(0.6, dy))
  }

  // Scrub action frames based on horizontal movement if not dragging vertically or playing
  if (!isDragging.value && !isPlayingAction.value && totalFrames.value > 1) {
    const norm = Math.max(0, Math.min(1, (mousePos.value.x + 0.6) / 1.2))
    currentFrame.value = Math.floor(norm * (totalFrames.value - 1))
  }
}

// Mouse tracking
function onMouseMove(e) {
  // #ifdef H5
  if (!isGyroActive.value) return
  const frame = document.querySelector('.interactive-phone-frame')
  if (frame) {
    const rect = frame.getBoundingClientRect()
    const centerX = rect.left + rect.width / 2
    const centerY = rect.top + rect.height / 2
    const dx = (e.clientX - centerX) / (rect.width * 1.2)
    const dy = (e.clientY - centerY) / (rect.height * 1.2)
    updateMousePosition(dx, dy)
  }
  // #endif
}

function handleGlobalMouseMove(e) {
  if (currentMode.value !== 'interactive') return

  if (isDragging.value) {
    const deltaY = dragStartY.value - e.clientY // drag upwards > 0
    if (dragStartProgress.value < 0.5) {
      // Swiping up to unlock
      const newP = Math.max(0, Math.min(1, deltaY / 170))
      unlockProgress.value = newP
    } else {
      // Already unlocked, swiping down to relock
      const newP = Math.max(0, Math.min(1, 1 - (-deltaY / 170)))
      unlockProgress.value = newP
    }
    if (!isPlayingAction.value && totalFrames.value > 1) {
      currentFrame.value = Math.min(totalFrames.value - 1, Math.round(unlockProgress.value * 48))
    }
  }

  if (!isGyroActive.value) return
  // #ifdef H5
  const frame = document.querySelector('.interactive-phone-frame')
  if (!frame) return
  const rect = frame.getBoundingClientRect()
  const centerX = rect.left + rect.width / 2
  const centerY = rect.top + rect.height / 2
  const dx = (e.clientX - centerX) / (rect.width * 1.5)
  const dy = (e.clientY - centerY) / (rect.height * 1.5)
  updateMousePosition(dx, dy)
  // #endif
}

function onMouseLeave() {
  mousePos.value = { x: 0, y: 0 }
  isHolding.value = false
}

function handleMouseDown(e) {
  if (isSleeping.value) {
    wakeScreen()
    return
  }
  isHolding.value = true
  dragStartY.value = e.clientY
  dragStartProgress.value = unlockProgress.value
  isDragging.value = true
}

function handleMouseUp() {
  isHolding.value = false
  if (isDragging.value) {
    isDragging.value = false
    if (dragStartProgress.value < 0.5) {
      if (unlockProgress.value >= 0.3) {
        triggerUnlockAnimation()
      } else {
        triggerLockAnimation()
      }
    } else {
      if (unlockProgress.value <= 0.7) {
        triggerLockAnimation()
      } else {
        triggerUnlockAnimation()
      }
    }
  }
}

function handleTouchStart(e) {
  if (e && e.cancelable && typeof e.preventDefault === 'function') {
    e.preventDefault()
  }
  if (isSleeping.value) {
    wakeScreen()
    return
  }
  isHolding.value = true
  const touch = (e.touches && e.touches.length > 0) ? e.touches[0] : e
  dragStartY.value = touch.clientY
  dragStartProgress.value = unlockProgress.value
  isDragging.value = true
}

function handleTouchMove(e) {
  if (e && e.cancelable && typeof e.preventDefault === 'function') {
    e.preventDefault()
  }
  const touch = (e.touches && e.touches.length > 0) ? e.touches[0] : e
  if (isDragging.value && touch) {
    const deltaY = dragStartY.value - touch.clientY
    if (dragStartProgress.value < 0.5) {
      unlockProgress.value = Math.max(0, Math.min(1, deltaY / 170))
    } else {
      unlockProgress.value = Math.max(0, Math.min(1, 1 - (-deltaY / 170)))
    }
    if (!isPlayingAction.value && totalFrames.value > 1) {
      currentFrame.value = Math.min(totalFrames.value - 1, Math.round(unlockProgress.value * 48))
    }
  }

  // #ifdef H5
  if (!isGyroActive.value) return
  const rect = e.currentTarget?.getBoundingClientRect?.()
  if (touch && rect) {
    const dx = (touch.clientX - (rect.left + rect.width / 2)) / rect.width
    const dy = (touch.clientY - (rect.top + rect.height / 2)) / rect.height
    updateMousePosition(dx, dy)
  }
  // #endif
}

function handleTouchEnd() {
  isHolding.value = false
  if (isDragging.value) {
    isDragging.value = false
    if (dragStartProgress.value < 0.5) {
      if (unlockProgress.value >= 0.3) {
        triggerUnlockAnimation()
      } else {
        triggerLockAnimation()
      }
    } else {
      if (unlockProgress.value <= 0.7) {
        triggerLockAnimation()
      } else {
        triggerUnlockAnimation()
      }
    }
  }
}

// Device orientation tilt listener for mobile browsers
function handleOrientation(e) {
  if (!isGyroActive.value) return
  const gamma = e.gamma || 0
  const beta = e.beta || 0
  const x = Math.max(-0.6, Math.min(0.6, gamma / 50))
  const y = Math.max(-0.6, Math.min(0.6, (beta - 45) / 50))
  updateMousePosition(x, y)
}

onMounted(() => {
  updateClock()
  timer = setInterval(updateClock, 1000)
  // #ifdef H5
  if (typeof window !== 'undefined' && window.addEventListener) {
    window.addEventListener('deviceorientation', handleOrientation, true)
    window.addEventListener('mousemove', handleGlobalMouseMove, { passive: true })
    window.addEventListener('mouseup', handleMouseUp)
  }
  // #endif
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
  if (actionTimer) clearInterval(actionTimer)
  // #ifdef H5
  if (typeof window !== 'undefined' && window.removeEventListener) {
    window.removeEventListener('deviceorientation', handleOrientation, true)
    window.removeEventListener('mousemove', handleGlobalMouseMove)
    window.removeEventListener('mouseup', handleMouseUp)
  }
  // #endif
})
</script>

<style scoped>
.mockup-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
  padding: 20rpx 0;
}

/* Tabs */
.mockup-tabs {
  display: flex;
  padding: 5px 6px;
  border-radius: 999px;
  margin-bottom: 20px;
  backdrop-filter: blur(16px);
  gap: 4px;
}

.mockup-container.dark .mockup-tabs {
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.mockup-container.light .mockup-tabs {
  background: rgba(0, 0, 0, 0.05);
  border: 1px solid rgba(0, 0, 0, 0.08);
}

.tab-pill {
  display: flex;
  align-items: center;
  padding: 6px 16px;
  border-radius: 999px;
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
}

.tab-icon {
  font-size: 14px;
  margin-right: 6px;
}

.tab-label {
  font-size: 13px;
  font-weight: 600;
  color: #4b5563;
  transition: all 0.2s ease;
}

.mockup-container.dark .tab-label {
  color: rgba(255, 255, 255, 0.75);
}

.mockup-container.light .tab-label {
  color: #1f2937;
}

.mockup-container.dark .tab-pill.active {
  background: #ffffff;
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.25);
}

.mockup-container.dark .tab-pill.active .tab-label {
  color: #111218 !important;
  font-weight: 700;
}

.mockup-container.light .tab-pill.active {
  background: #111827;
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.15);
}

.mockup-container.light .tab-pill.active .tab-label {
  color: #ffffff !important;
  font-weight: 700;
}

/* ⚡ Interactive Tab Special Glow Styling */
.tab-pill-interactive {
  background: linear-gradient(135deg, rgba(139, 92, 246, 0.15), rgba(99, 102, 241, 0.15));
  border: 1px solid rgba(139, 92, 246, 0.3);
}

.tab-pill-interactive.active {
  background: linear-gradient(135deg, #8b5cf6, #6366f1) !important;
  box-shadow: 0 4px 18px rgba(139, 92, 246, 0.45) !important;
}

.tab-pill-interactive.active .tab-label {
  color: #ffffff !important;
}

.live-pulse-dot {
  width: 6px;
  height: 6px;
  background: #10b981;
  border-radius: 50%;
  margin-left: 6px;
  box-shadow: 0 0 8px #10b981;
  animation: pulseDot 1.8s infinite;
}

@keyframes pulseDot {
  0% { transform: scale(0.9); opacity: 0.8; }
  50% { transform: scale(1.4); opacity: 1; }
  100% { transform: scale(0.9); opacity: 0.8; }
}

/* Viewport */
.viewport-wrapper {
  display: flex;
  justify-content: center;
  width: 100%;
}

/* ======================================================== */
/* ⚡ MODE 0: 3D Interactive .tendies Simulator */
/* ======================================================== */
.interactive-stage {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
  padding: 10rpx 0;
  touch-action: pan-y;
  -webkit-user-select: none;
  user-select: none;
}

.desktop-interact-tip {
  margin-bottom: 20rpx;
  display: flex;
  justify-content: center;
}

.tip-badge {
  display: inline-flex;
  align-items: center;
  gap: 8rpx;
  padding: 8rpx 22rpx;
  border-radius: 999rpx;
  background: rgba(139, 92, 246, 0.12);
  border: 1rpx solid rgba(139, 92, 246, 0.28);
}

.tip-dot {
  width: 8rpx;
  height: 8rpx;
  border-radius: 50%;
  background: #10b981;
  box-shadow: 0 0 8rpx #10b981;
  animation: pulseDot 1.8s infinite;
}

.tip-text {
  font-size: 20rpx;
  color: #a78bfa;
  font-weight: 600;
  letter-spacing: 0.5rpx;
}

.interactive-phone-frame {
  position: relative;
  width: 440rpx;
  height: auto;
  aspect-ratio: 393 / 852;
  border-radius: 64rpx;
  border: 8rpx solid #2a2c38;
  overflow: hidden;
  /* iOS WebKit clipping fix for border-radius + 3D transformed child elements */
  -webkit-mask-image: -webkit-radial-gradient(white, black);
  mask-image: radial-gradient(white, black);
  -webkit-backface-visibility: hidden;
  backface-visibility: hidden;
  box-shadow: 0 32rpx 100rpx rgba(0, 0, 0, 0.65), 0 0 36px rgba(139, 92, 246, 0.25);
  background: #000;
  cursor: grab;
  user-select: none;
  -webkit-user-select: none;
  touch-action: none;
  overscroll-behavior: contain;
  transition: transform 0.1s cubic-bezier(0.2, 0.8, 0.2, 1), filter 0.4s ease;
  transform-style: preserve-3d;
}

.interactive-phone-frame:active,
.interactive-phone-frame.is-dragging {
  cursor: grabbing;
}

.interactive-phone-frame.is-pressing {
  box-shadow: 0 16rpx 40rpx rgba(0, 0, 0, 0.7), 0 0 45px rgba(139, 92, 246, 0.4);
}

.interactive-phone-frame.aod-active {
  filter: brightness(0.38) contrast(1.15);
}

.interactive-phone-frame.is-sleeping {
  box-shadow: 0 20rpx 60rpx rgba(0, 0, 0, 0.9), 0 0 10px rgba(0, 0, 0, 0.5);
}

.parallax-layer {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  transition: transform 0.08s ease-out;
}

.bg-layer {
  z-index: 1;
  transition: transform 0.4s cubic-bezier(0.2, 0.8, 0.2, 1);
}

.bg-unlocked-zoom {
  transform: scale(1.04) !important;
}

.motion-video-layer {
  z-index: 2;
  opacity: 0;
  transition: opacity 0.3s ease;
  pointer-events: none;
}

.motion-video-layer.active {
  opacity: 1;
}

/* ======================================================== */
/* Lock Screen Layer (Header, Cutout, Shortcuts) */
/* ======================================================== */
.lockscreen-ui-wrapper {
  position: absolute;
  inset: 0;
  z-index: 5;
  transition: opacity 0.4s cubic-bezier(0.2, 0.8, 0.2, 1), transform 0.45s cubic-bezier(0.2, 0.8, 0.2, 1);
}

.lockscreen-ui-wrapper.ui-fade-out {
  opacity: 0;
  transform: translateY(-80rpx) scale(0.95);
  pointer-events: none;
}

.interactive-clock-layer {
  z-index: 50 !important;
  transform-style: preserve-3d;
}

.clock-fade-out {
  opacity: 0 !important;
  pointer-events: none !important;
}

.clock-aod {
  color: rgba(255, 255, 255, 0.7) !important;
  text-shadow: 0 0 10px rgba(255, 255, 255, 0.2);
}

.fg-cutout-layer {
  z-index: 80; /* In front of clock! Provides authentic iOS 17 depth occlusion */
  transform: translate3d(0, 0, 48px);
}

/* Clock Halo / Ring Layer (e.g. 9 SD Gundams circling around the clock) */
.clock-ring-wrapper {
  position: absolute;
  top: -205rpx;
  left: 50%;
  width: 380rpx;
  height: 380rpx;
  margin-left: -190rpx;
  z-index: 60; /* In front of clock */
  pointer-events: none;
  transition: transform 0.45s cubic-bezier(0.2, 0.8, 0.2, 1), opacity 0.4s ease;
  transform-origin: center center;
}

.clock-ring-wrapper.ring-hidden {
  opacity: 0;
  transform: translateY(-160rpx) scale(0.85);
}

.clock-ring-image {
  width: 100%;
  height: 100%;
  display: block;
  animation: spinClockHalo 36s linear infinite;
  transform-origin: 50% 50%;
}

@keyframes spinClockHalo {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

/* V-Arch Flight Trajectory Container (9 SD Gundams) */
.varch-flight-container {
  position: absolute;
  top: 70rpx;
  left: 0;
  width: 100%;
  height: calc(100% - 70rpx);
  z-index: 70;
  pointer-events: none;
}

.varch-wing {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  transform-origin: center center;
  will-change: transform, opacity;
}

.varch-wing-left {
  transform-origin: left bottom;
}

.varch-wing-right {
  transform-origin: right bottom;
}

.varch-wing-center {
  transform-origin: center center;
}

/* OLED True Black Sleep Mode Overlay */
.screen-sleep-overlay {
  position: absolute;
  inset: 0;
  background: #000000;
  z-index: 400;
  transform: translate3d(0, 0, 220px);
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.35s cubic-bezier(0.2, 0.8, 0.2, 1);
  display: flex;
  align-items: center;
  justify-content: center;
}

.screen-sleep-overlay.active {
  opacity: 1;
  pointer-events: auto;
  cursor: pointer;
}

.sleep-tip-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14rpx;
  animation: sleepPulse 2.6s ease-in-out infinite;
}

.sleep-glow-dot {
  width: 14rpx;
  height: 14rpx;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 50%;
  box-shadow: 0 0 18rpx rgba(255, 255, 255, 0.95);
}

.sleep-text {
  font-size: 22rpx;
  color: rgba(255, 255, 255, 0.6);
  font-weight: 500;
  letter-spacing: 2rpx;
}

@keyframes sleepPulse {
  0%, 100% {
    opacity: 0.35;
    transform: scale(0.96);
  }
  50% {
    opacity: 0.9;
    transform: scale(1.03);
  }
}

.dynamic-island.island-hidden {
  opacity: 0;
  transition: opacity 0.3s ease;
}

/* Fallback Unlocked 9 SD Gundams V-Arch Layer */
.unlocked-varch-layer {
  position: absolute;
  top: 70rpx;
  left: 0;
  width: 100%;
  height: calc(100% - 70rpx);
  z-index: 70;
  pointer-events: none;
  opacity: 0;
  transform: translateY(35rpx) scale(0.92);
  transition: opacity 0.5s cubic-bezier(0.2, 0.8, 0.2, 1), transform 0.5s cubic-bezier(0.2, 0.8, 0.2, 1);
}

.unlocked-varch-layer.varch-active {
  opacity: 1;
  transform: translateY(0) scale(1);
  animation: floatVArch 6s ease-in-out infinite alternate;
}

@keyframes floatVArch {
  0% {
    transform: translateY(0);
  }
  100% {
    transform: translateY(-8rpx);
  }
}

.phone-glass-glare {
  position: absolute;
  inset: 0;
  pointer-events: none;
  z-index: 300;
  transform: translate3d(0, 0, 200px);
  mix-blend-mode: overlay;
  transition: opacity 0.2s ease;
}

.dynamic-island.island-aod {
  opacity: 0.3;
}

.interactive-security-watermark {
  position: absolute;
  bottom: 140rpx;
  left: 50%;
  transform: translateX(-50%) translate3d(0, 0, 210px);
  background: rgba(0, 0, 0, 0.22);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  padding: 6rpx 22rpx;
  border-radius: 999rpx;
  border: 1rpx solid rgba(255, 255, 255, 0.16);
  z-index: 350;
  pointer-events: none;
  display: inline-flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  white-space: nowrap;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.18);
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.wm-text {
  font-size: 18rpx;
  color: rgba(255, 255, 255, 0.82);
  letter-spacing: 0.8rpx;
  white-space: nowrap;
  word-break: keep-all;
  display: inline-block;
  line-height: 1.2;
}

.footer-hidden {
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.3s ease;
}

.swipe-unlock-trigger {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6rpx;
  cursor: pointer;
  padding: 10rpx;
}

.swipe-arrow {
  font-size: 20rpx;
  color: rgba(255, 255, 255, 0.85);
  animation: bounceUp 1.8s infinite;
}

.swipe-text {
  font-size: 18rpx;
  color: rgba(255, 255, 255, 0.75);
  font-weight: 500;
  letter-spacing: 1rpx;
}

@keyframes bounceUp {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-8rpx); }
}

/* ======================================================== */
/* Unlocked iOS Home Screen Springboard */
/* ======================================================== */
.homescreen-ui-wrapper {
  position: absolute;
  inset: 0;
  z-index: 200 !important;
  opacity: 0;
  visibility: hidden;
  transform: translate3d(0, 0, 100px) scale(0.92);
  pointer-events: none;
  transition: opacity 0.4s cubic-bezier(0.2, 0.8, 0.2, 1), transform 0.4s cubic-bezier(0.2, 0.8, 0.2, 1);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 70rpx 28rpx 24rpx;
  box-sizing: border-box;
  transform-style: preserve-3d;
}

.homescreen-ui-wrapper.ui-show {
  opacity: 1;
  visibility: visible;
  transform: translate3d(0, 0, 100px) scale(1);
  pointer-events: auto;
}

.homescreen-ui-wrapper:not(.ui-show) {
  visibility: hidden !important;
  pointer-events: none !important;
  opacity: 0 !important;
}

.homescreen-ui-wrapper:not(.ui-show) .springboard-bottom-group,
.homescreen-ui-wrapper:not(.ui-show) .iphone-dock-bar,
.homescreen-ui-wrapper:not(.ui-show) .ios-search-pill {
  visibility: hidden !important;
  opacity: 0 !important;
}

.phone-status-bar {
  position: absolute;
  top: 14rpx;
  left: 0;
  right: 0;
  width: 100% !important;
  height: 38rpx;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 16rpx 0 28rpx;
  box-sizing: border-box;
  z-index: 260;
  transform: translate3d(0, 0, 150px);
  pointer-events: none;
  transition: opacity 0.3s ease;
}

.phone-status-bar.status-hidden {
  opacity: 0;
}

.home-status-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 10rpx;
}

.status-time {
  font-size: 20rpx;
  font-weight: 600;
  color: #ffffff;
  line-height: 1;
  font-family: -apple-system, BlinkMacSystemFont, "SF Pro Text", "Helvetica Neue", sans-serif;
  letter-spacing: -0.5rpx;
  text-shadow: 0 1rpx 4rpx rgba(0, 0, 0, 0.4);
  transition: opacity 0.3s ease;
}

.home-top-actions {
  display: flex;
  align-items: center;
  gap: 10rpx;
}

.app-toggle-pill {
  background: rgba(255, 255, 255, 0.16);
  backdrop-filter: blur(10px);
  padding: 4rpx 14rpx;
  border-radius: 999rpx;
  border: 1rpx solid rgba(255, 255, 255, 0.22);
  cursor: pointer;
  transition: all 0.2s ease;
}

.app-toggle-pill:hover {
  background: rgba(255, 255, 255, 0.28);
}

.toggle-pill-text {
  font-size: 18rpx;
  color: #fff;
  font-weight: 600;
}

.relock-btn {
  background: rgba(0, 0, 0, 0.45);
  backdrop-filter: blur(10px);
  padding: 4rpx 14rpx;
  border-radius: 999rpx;
  border: 1rpx solid rgba(255, 255, 255, 0.25);
  cursor: pointer;
  transition: all 0.2s ease;
}

.relock-btn:hover {
  background: rgba(0, 0, 0, 0.65);
}

.relock-icon {
  font-size: 18rpx;
  color: #fff;
  font-weight: 600;
}

.status-right-icons {
  display: flex;
  align-items: center;
  gap: 8rpx;
}

/* Authentic iOS Cellular Signal Bars */
.ios-signal-bars {
  display: flex;
  align-items: flex-end;
  gap: 2.5rpx;
  height: 16rpx;
  margin-right: 2rpx;
}

.signal-bar {
  width: 3.5rpx;
  background: #ffffff;
  border-radius: 1.5rpx;
  box-shadow: 0 1rpx 3rpx rgba(0, 0, 0, 0.3);
}

.signal-bar.bar-1 { height: 6rpx; }
.signal-bar.bar-2 { height: 9rpx; }
.signal-bar.bar-3 { height: 12rpx; }
.signal-bar.bar-4 { height: 16rpx; }

/* Authentic iOS 5G Typography */
.ios-5g-text {
  font-size: 16rpx;
  font-weight: 700;
  color: #ffffff;
  line-height: 1;
  letter-spacing: -0.5rpx;
  font-family: -apple-system, BlinkMacSystemFont, "SF Pro Text", "Helvetica Neue", sans-serif;
  text-shadow: 0 1rpx 3rpx rgba(0, 0, 0, 0.3);
}

/* Authentic iOS Battery Indicator */
.ios-battery-group {
  display: flex;
  align-items: center;
  gap: 4rpx;
}

.ios-battery-pct {
  font-size: 16rpx;
  font-weight: 600;
  color: #ffffff;
  line-height: 1;
  font-family: -apple-system, BlinkMacSystemFont, "SF Pro Text", "Helvetica Neue", sans-serif;
  letter-spacing: -0.5rpx;
  text-shadow: 0 1rpx 3rpx rgba(0, 0, 0, 0.3);
}

.ios-battery-pill {
  position: relative;
  width: 32rpx;
  height: 16rpx;
  border: 1.6rpx solid rgba(255, 255, 255, 0.85);
  border-radius: 4.5rpx;
  padding: 1.8rpx;
  box-sizing: border-box;
  display: flex;
  align-items: center;
  box-shadow: 0 1rpx 3rpx rgba(0, 0, 0, 0.3);
}

.ios-battery-level {
  height: 100%;
  background: #ffffff;
  border-radius: 2rpx;
}

.ios-battery-terminal {
  position: absolute;
  right: -3rpx;
  top: 50%;
  transform: translateY(-50%);
  width: 1.8rpx;
  height: 5.5rpx;
  background: rgba(255, 255, 255, 0.85);
  border-radius: 0 1.5rpx 1.5rpx 0;
}

/* Clean Mode Spacer */
.clean-springboard-spacer {
  flex: 1;
}

/* Bottom Group (Search Pill + Dock) */
.springboard-bottom-group {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16rpx;
  margin-top: auto;
  position: relative;
  z-index: 230;
  transform-style: preserve-3d;
}

/* Authentic iOS Search Pill */
.ios-search-pill {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8rpx;
  padding: 6rpx 28rpx;
  background: rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-radius: 999rpx;
  border: 1rpx solid rgba(255, 255, 255, 0.16);
  box-shadow: 0 4rpx 14rpx rgba(0, 0, 0, 0.3);
  position: relative;
  z-index: 235;
  transform: translate3d(0, 0, 15px);
}

.search-glass-icon {
  font-size: 18rpx;
  opacity: 0.85;
}

.search-text {
  font-size: 20rpx;
  color: rgba(255, 255, 255, 0.9);
  font-weight: 500;
  letter-spacing: 0.5rpx;
}

/* App Grid (4 cols x 3 rows) */
.springboard-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 22rpx 16rpx;
  margin-top: 30rpx;
  position: relative;
  z-index: 220;
  transform: translate3d(0, 0, 20px);
  transform-style: preserve-3d;
}

.app-icon-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6rpx;
  position: relative;
  z-index: 225;
  transform: translate3d(0, 0, 10px);
}

.app-icon {
  width: 76rpx;
  height: 76rpx;
  border-radius: 18rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.35);
  position: relative;
  overflow: hidden;
  z-index: 226;
}

.emoji-icon {
  font-size: 40rpx;
}

.icon-calendar {
  background: #ffffff;
  flex-direction: column;
}

.cal-day {
  font-size: 12rpx;
  color: #ef4444;
  font-weight: 700;
  line-height: 1;
}

.cal-num {
  font-size: 34rpx;
  color: #111;
  font-weight: 800;
  line-height: 1;
}

.icon-photos { background: linear-gradient(135deg, #fef08a, #f43f5e); }
.icon-camera { background: linear-gradient(135deg, #64748b, #334155); }
.icon-clock { background: #000000; border: 1rpx solid #333; }
.icon-weather { background: linear-gradient(135deg, #38bdf8, #2563eb); }
.icon-maps { background: linear-gradient(135deg, #4ade80, #059669); }
.icon-notes { background: linear-gradient(135deg, #fde047, #ca8a04); }
.icon-settings { background: linear-gradient(135deg, #94a3b8, #475569); }
.icon-appstore { background: linear-gradient(135deg, #0284c7, #0369a1); }
.appstore-letter { font-size: 44rpx; font-weight: 900; color: #fff; font-family: sans-serif; }
.icon-wechat { background: #07c160; }
.icon-music { background: linear-gradient(135deg, #fb7185, #e11d48); }
.icon-safari { background: #ffffff; }

.app-name {
  font-size: 16rpx;
  color: #ffffff;
  text-shadow: 0 2rpx 4rpx rgba(0, 0, 0, 0.8);
  font-weight: 500;
}

/* Frosted Glass Dock */
.iphone-dock-bar {
  display: flex;
  justify-content: space-around;
  align-items: center;
  width: 100%;
  box-sizing: border-box;
  background: rgba(255, 255, 255, 0.22);
  backdrop-filter: blur(28px) saturate(180%);
  -webkit-backdrop-filter: blur(28px) saturate(180%);
  border: 1rpx solid rgba(255, 255, 255, 0.3);
  padding: 12rpx 18rpx;
  border-radius: 36rpx;
  box-shadow: 0 10rpx 30rpx rgba(0, 0, 0, 0.25);
  position: relative;
  z-index: 240;
  transform: translate3d(0, 0, 20px);
}

.dock-app {
  width: 72rpx;
  height: 72rpx;
  border-radius: 18rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.dock-app.icon-phone { background: #22c55e; }
.dock-app.icon-chrome { background: #ffffff; }
.dock-app.icon-wechat { background: #07c160; }
.dock-app.icon-camera { background: #334155; }

.app-badge {
  position: absolute;
  top: -6rpx;
  right: -6rpx;
  background: #ef4444;
  color: #ffffff;
  border-radius: 999rpx;
  min-width: 24rpx;
  height: 24rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 4rpx;
  border: 2rpx solid #ffffff;
  box-shadow: 0 2rpx 6rpx rgba(0, 0, 0, 0.3);
}

.badge-num {
  font-size: 14rpx;
  font-weight: 700;
  line-height: 1;
}

.dock-app-icon {
  font-size: 38rpx;
}

.unlocked-home-bar {
  display: flex;
  justify-content: center;
  cursor: pointer;
  padding: 4rpx 0;
  position: relative;
  z-index: 250;
  transform: translate3d(0, 0, 25px);
}

/* Interactive Dock Toolbar */
.interactive-dock-toolbar {
  display: flex;
  align-items: center;
  gap: 16rpx;
  margin-top: 24rpx;
  flex-wrap: wrap;
  justify-content: center;
}

.dock-tool-btn {
  display: flex;
  align-items: center;
  gap: 8rpx;
  padding: 10rpx 22rpx;
  border-radius: 999rpx;
  background: rgba(255, 255, 255, 0.08);
  border: 1rpx solid rgba(255, 255, 255, 0.12);
  cursor: pointer;
  transition: all 0.2s ease;
}

.mockup-container.light .dock-tool-btn {
  background: rgba(0, 0, 0, 0.06);
  border: 1rpx solid rgba(0, 0, 0, 0.1);
}

.dock-tool-btn.active {
  background: #8b5cf6;
  border-color: #8b5cf6;
  box-shadow: 0 4rpx 14rpx rgba(139, 92, 246, 0.4);
}

.dock-tool-btn.active .tool-text {
  color: #ffffff;
}

.btn-action-play {
  background: linear-gradient(135deg, rgba(6, 182, 212, 0.2), rgba(59, 130, 246, 0.2));
  border-color: rgba(6, 182, 212, 0.35);
}

.btn-action-play.active {
  background: linear-gradient(135deg, #06b6d4, #3b82f6) !important;
  box-shadow: 0 4rpx 16rpx rgba(6, 182, 212, 0.45) !important;
}

.tool-icon {
  font-size: 22rpx;
}

.tool-text {
  font-size: 22rpx;
  font-weight: 600;
  color: var(--text-color, #e2e8f0);
}

.mockup-container.light .tool-text {
  color: #334155;
}

/* ======================================================== */
/* MODE 1: Standard Phone Frame */
/* ======================================================== */
.phone-frame {
  position: relative;
  width: 440rpx;
  height: 954rpx;
  aspect-ratio: 393 / 852;
  border-radius: 64rpx;
  border: 8rpx solid #2a2c38;
  overflow: hidden;
  box-shadow: 0 28rpx 90rpx rgba(0, 0, 0, 0.5);
  background: #000;
  touch-action: pan-y;
  user-select: none;
  -webkit-user-select: none;
}

.mockup-img {
  width: 100%;
  height: 100%;
  display: block;
}

.placeholder-bg {
  background: linear-gradient(135deg, #1f2937, #111827);
}

.dynamic-island {
  position: absolute;
  top: 14rpx;
  left: 50%;
  transform: translateX(-50%) translate3d(0, 0, 155px);
  width: 126rpx;
  height: 38rpx;
  background: #000000;
  border-radius: 999rpx;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding-right: 12rpx;
  z-index: 265;
  transition: opacity 0.3s ease;
}

.camera-lens {
  width: 12rpx;
  height: 12rpx;
  border-radius: 50%;
  background: #1a1a3a;
  box-shadow: inset 0 0 3rpx #000;
}

.phone-lock-header {
  position: absolute;
  top: 104rpx;
  left: 0;
  right: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  z-index: 50;
  text-shadow: 0 2rpx 10rpx rgba(0, 0, 0, 0.6);
  pointer-events: none;
}

.lock-date-row {
  display: flex;
  align-items: center;
  margin-bottom: 6rpx;
}

.lock-icon {
  font-size: 20rpx;
  margin-right: 8rpx;
  color: rgba(255, 255, 255, 0.85);
}

.lock-date-text {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.9);
  font-weight: 500;
}

.lock-time-text {
  font-size: 104rpx;
  font-weight: 700;
  color: #ffffff;
  line-height: 1;
  letter-spacing: -2rpx;
  font-family: -apple-system, BlinkMacSystemFont, "SF Pro Display", sans-serif;
  transition: color 0.3s ease;
}

.phone-lock-footer {
  position: absolute;
  bottom: 24rpx;
  left: 36rpx;
  right: 36rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  z-index: 55;
}

.footer-btn {
  width: 72rpx;
  height: 72rpx;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.45);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  display: flex;
  align-items: center;
  justify-content: center;
}

.footer-icon {
  font-size: 28rpx;
}

.home-indicator {
  width: 150rpx;
  height: 8rpx;
  background: #ffffff;
  border-radius: 999rpx;
}

/* ======================================================== */
/* MODE 2: Desktop Frame */
/* ======================================================== */
.desktop-frame {
  position: relative;
  width: 680rpx;
  height: 450rpx;
  border-radius: 24rpx;
  border: 8rpx solid #2a2c38;
  overflow: hidden;
  box-shadow: 0 28rpx 90rpx rgba(0, 0, 0, 0.5);
  background: #000;
}

.mac-menu-bar {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 38rpx;
  background: rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(16px);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 18rpx;
  z-index: 10;
}

.mac-left {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.mac-apple-logo {
  font-size: 22rpx;
  color: #ffffff;
}

.mac-menu-item {
  font-size: 18rpx;
  color: rgba(255, 255, 255, 0.85);
  font-weight: 500;
}

.mac-right {
  display: flex;
  align-items: center;
  gap: 14rpx;
}

.mac-status-icon,
.mac-status-time {
  font-size: 18rpx;
  color: rgba(255, 255, 255, 0.85);
}

.mac-wallpaper {
  width: 100%;
  height: 100%;
}

.mac-dock {
  position: absolute;
  bottom: 14rpx;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 18rpx;
  background: rgba(255, 255, 255, 0.16);
  backdrop-filter: blur(28px) saturate(200%);
  -webkit-backdrop-filter: blur(28px) saturate(200%);
  padding: 10rpx 24rpx;
  border-radius: 24rpx;
  border: 1rpx solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 8rpx 32rpx rgba(0, 0, 0, 0.3), inset 0 1rpx 2rpx rgba(255, 255, 255, 0.4);
  z-index: 10;
}

.dock-item {
  font-size: 24rpx;
}

/* ======================================================== */
/* MODE 3: Clean Frame */
/* ======================================================== */
.clean-frame {
  position: relative;
  width: 480rpx;
  height: 850rpx;
  border-radius: 36rpx;
  overflow: hidden;
  box-shadow: 0 24rpx 70rpx rgba(0, 0, 0, 0.4);
}

.clean-wallpaper {
  width: 100%;
  height: 100%;
}

.clean-hint {
  position: absolute;
  bottom: 24rpx;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(0, 0, 0, 0.65);
  padding: 10rpx 28rpx;
  border-radius: 999rpx;
  backdrop-filter: blur(12px);
}

.hint-text {
  font-size: 22rpx;
  color: #ffffff;
  font-weight: 500;
  letter-spacing: 1rpx;
}

/* Dynamic Apple CAML Sprite Layers */
.caml-sprite-layer {
  position: absolute;
  display: block;
  user-select: none;
  -webkit-user-drag: none;
  pointer-events: none;
  will-change: transform, opacity;
}

.caml-sprite-layer :deep(div),
.caml-sprite-layer :deep(img) {
  width: 100% !important;
  height: 100% !important;
  object-fit: fill !important;
  background-size: 100% 100% !important;
}

.caml-bg-layer {
  /* Situated behind the lock clock */
}

.caml-floating-layer {
  /* Situated in front of the lock clock */
}

.caml-anim-pos,
.caml-anim-rot {
  width: 100%;
  height: 100%;
  pointer-events: none;
  transform-origin: center center;
}

.caml-sprite-inner-img {
  width: 100%;
  height: 100%;
  display: block;
  pointer-events: none;
  object-fit: fill;
}
</style>
