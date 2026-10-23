<template>
  <view class="ambient-backdrop" :class="[appStore.theme]">
    <!-- Native Canvas Container with explicit ID and Ref -->
    <!-- #ifdef H5 -->
    <canvas
      id="ambient-canvas-elem"
      ref="canvasRef"
      class="ambient-canvas"
    ></canvas>
    <!-- #endif -->

    <!-- CSS Fallback for Mini-Programs -->
    <!-- #ifndef H5 -->
    <view class="mp-ambient-fx" :class="[appStore.theme]">
      <view v-if="appStore.theme === 'dark'" class="mp-stars-layer">
        <view v-for="n in 36" :key="n" class="mp-star" :style="getMpStarStyle(n)"></view>
        <view class="mp-meteor mp-meteor-1"></view>
        <view class="mp-meteor mp-meteor-2"></view>
      </view>
      <view v-else class="mp-rain-layer">
        <view v-for="n in 30" :key="n" class="mp-raindrop" :style="getMpRainStyle(n)"></view>
      </view>
    </view>
    <!-- #endif -->

  </view>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useAppStore } from '@/store/app'

const appStore = useAppStore()
const canvasRef = ref(null)

let animId = null
let ctx = null
let canvasEl = null
let width = 0
let height = 0
let dpr = 1

// Mouse coordinates
const mouse = { x: -1000, y: -1000, targetX: -1000, targetY: -1000, active: false }

/* ========================================================
   Dark Mode: Deep Cosmic Starfield & Streaming Meteors
   ======================================================== */
const stars = []
const meteors = []
let nextMeteorTime = 0

class Star {
  constructor(w, h) {
    this.reset(w, h, true)
  }

  reset(w, h, initial = false) {
    this.x = Math.random() * w
    this.y = initial ? Math.random() * h : Math.random() * h
    this.radius = 0.8 + Math.random() * 1.8
    this.baseAlpha = 0.25 + Math.random() * 0.45
    this.alpha = this.baseAlpha
    this.twinkleSpeed = 0.003 + Math.random() * 0.006
    this.phase = Math.random() * Math.PI * 2
    this.isCross = Math.random() < 0.1
    this.dx = (Math.random() - 0.5) * 0.06
    this.dy = (Math.random() - 0.5) * 0.04

    const hues = ['#ffffff', '#ffffff', '#7dd3fc', '#bae6fd', '#fef08a', '#e879f9']
    this.color = hues[Math.floor(Math.random() * hues.length)]
  }

  update(w, h) {
    this.phase += this.twinkleSpeed
    this.alpha = this.baseAlpha + Math.sin(this.phase) * 0.28
    if (this.alpha < 0.1) this.alpha = 0.1
    if (this.alpha > 0.85) this.alpha = 0.85

    this.x += this.dx
    this.y += this.dy
    if (w && h) {
      if (this.x < 0) this.x = w
      else if (this.x > w) this.x = 0
      if (this.y < 0) this.y = h
      else if (this.y > h) this.y = 0
    }
  }

  draw(c) {
    c.save()
    c.globalAlpha = this.alpha
    c.fillStyle = this.color
    c.beginPath()
    c.arc(this.x, this.y, this.radius, 0, Math.PI * 2)
    c.fill()

    if (this.radius > 1.6) {
      c.globalAlpha = this.alpha * 0.4
      c.beginPath()
      c.arc(this.x, this.y, this.radius * 3.2, 0, Math.PI * 2)
      c.fill()
    }

    if (this.isCross && this.alpha > 0.6) {
      c.strokeStyle = this.color
      c.lineWidth = 1
      c.globalAlpha = this.alpha * 0.7
      const flareLen = this.radius * 4.5
      c.beginPath()
      c.moveTo(this.x - flareLen, this.y)
      c.lineTo(this.x + flareLen, this.y)
      c.moveTo(this.x, this.y - flareLen)
      c.lineTo(this.x, this.y + flareLen)
      c.stroke()
    }
    c.restore()
  }
}

class Meteor {
  constructor(w, h) {
    this.w = w
    this.h = h
    this.angle = (32 + Math.random() * 10) * (Math.PI / 180)
    // Calm, slow, peaceful glide (taking ~2.5 - 3.5s across the screen)
    this.speed = 3.6 + Math.random() * 2.0
    this.length = 130 + Math.random() * 70
    this.x = Math.random() * (w * 1.1)
    this.y = -30 - Math.random() * 60
    this.dx = -Math.cos(this.angle) * this.speed
    this.dy = Math.sin(this.angle) * this.speed
    this.life = 0
    this.maxLife = Math.floor(130 + Math.random() * 60)
    this.opacity = 0.85
    this.dead = false
  }

  update() {
    this.x += this.dx
    this.y += this.dy
    this.life++
    if (this.life >= this.maxLife || this.y > this.h + 80 || this.x < -160) {
      this.dead = true
    }
  }

  draw(c) {
    c.save()
    const progress = this.life / this.maxLife
    let fade = 1
    if (progress < 0.18) fade = progress / 0.18
    else if (progress > 0.6) fade = (1 - progress) / 0.4

    const alpha = this.opacity * Math.max(0, fade)
    if (alpha <= 0.02) {
      c.restore()
      return
    }

    const tailX = this.x - this.dx * (this.length / this.speed)
    const tailY = this.y - this.dy * (this.length / this.speed)

    const grad = c.createLinearGradient(this.x, this.y, tailX, tailY)
    grad.addColorStop(0, `rgba(255, 255, 255, ${alpha * 0.95})`)
    grad.addColorStop(0.25, `rgba(56, 189, 248, ${alpha * 0.75})`)
    grad.addColorStop(0.7, `rgba(129, 140, 248, ${alpha * 0.3})`)
    grad.addColorStop(1, 'rgba(255, 255, 255, 0)')

    c.strokeStyle = grad
    c.lineWidth = 1.8
    c.lineCap = 'round'
    c.beginPath()
    c.moveTo(this.x, this.y)
    c.lineTo(tailX, tailY)
    c.stroke()

    // Soft glowing comet head
    c.fillStyle = `rgba(255, 255, 255, ${alpha * 0.95})`
    c.beginPath()
    c.arc(this.x, this.y, 2.2, 0, Math.PI * 2)
    c.fill()

    c.fillStyle = `rgba(56, 189, 248, ${alpha * 0.45})`
    c.beginPath()
    c.arc(this.x, this.y, 5, 0, Math.PI * 2)
    c.fill()
    c.restore()
  }
}

/* ========================================================
   Light Mode: Frosted Glass Daylight + Raindrops & Ripples
   ======================================================== */
const raindrops = []
const glassDroplets = []
const ripples = []

class Raindrop {
  constructor(w, h) {
    this.reset(w, h, true)
  }

  reset(w, h, initial = false) {
    this.x = Math.random() * (w + 80) - 40
    this.y = initial ? Math.random() * h : -20 - Math.random() * 40
    this.length = 12 + Math.random() * 14
    // Slow, serene daylight drizzle
    this.speed = 2.6 + Math.random() * 2.0
    this.wind = -0.3 + Math.random() * 0.15
    this.alpha = 0.22 + Math.random() * 0.22
    this.thickness = 1.0 + Math.random() * 0.5
  }

  update(w, h) {
    this.y += this.speed
    this.x += this.wind
    if (this.y > h + 20) {
      if (Math.random() < 0.2) {
        ripples.push(new Ripple(this.x, h - 20 + Math.random() * 20))
      }
      this.reset(w, h)
    }
  }

  draw(c) {
    c.save()
    const grad = c.createLinearGradient(this.x, this.y, this.x + this.wind * 3, this.y + this.length)
    grad.addColorStop(0, 'rgba(148, 163, 184, 0)')
    grad.addColorStop(0.3, `rgba(100, 116, 139, ${this.alpha * 0.5})`)
    grad.addColorStop(1, `rgba(71, 85, 105, ${this.alpha})`)

    c.strokeStyle = grad
    c.lineWidth = this.thickness
    c.lineCap = 'round'
    c.beginPath()
    c.moveTo(this.x, this.y)
    c.lineTo(this.x + this.wind * 3, this.y + this.length)
    c.stroke()
    c.restore()
  }
}

class GlassDroplet {
  constructor(w, h) {
    this.reset(w, h)
  }

  reset(w, h) {
    this.x = 20 + Math.random() * (w - 40)
    this.y = 30 + Math.random() * (h - 60)
    this.radius = 2.2 + Math.random() * 2.6
    this.alpha = 0.45 + Math.random() * 0.2
    this.slideSpeed = 0
    this.slideCounter = Math.floor(Math.random() * 360)
    this.shimmer = Math.random() * Math.PI * 2
  }

  update(w, h) {
    this.shimmer += 0.015
    this.slideCounter++
    // Very occasional gentle trickle
    if (this.slideCounter > 450) {
      this.slideSpeed = 0.4 + Math.random() * 0.5
      this.slideCounter = 0
    }
    if (this.slideSpeed > 0) {
      this.y += this.slideSpeed
      this.slideSpeed *= 0.98
      if (this.slideSpeed < 0.02) this.slideSpeed = 0
    }
    if (this.y > h + 10) {
      this.reset(w, h)
      this.y = -5
    }
  }

  draw(c) {
    c.save()
    const currentAlpha = this.alpha + Math.sin(this.shimmer) * 0.06

    // Soft drop shadow
    c.fillStyle = `rgba(15, 23, 42, ${currentAlpha * 0.15})`
    c.beginPath()
    c.arc(this.x + 0.8, this.y + 1, this.radius, 0, Math.PI * 2)
    c.fill()

    // Glass droplet
    c.fillStyle = `rgba(148, 163, 184, ${currentAlpha * 0.45})`
    c.beginPath()
    c.arc(this.x, this.y, this.radius, 0, Math.PI * 2)
    c.fill()

    // Specular highlight
    c.fillStyle = `rgba(255, 255, 255, ${currentAlpha * 0.85})`
    c.beginPath()
    c.arc(this.x - this.radius * 0.3, this.y - this.radius * 0.3, this.radius * 0.4, 0, Math.PI * 2)
    c.fill()
    c.restore()
  }
}

class Ripple {
  constructor(x, y) {
    this.x = x
    this.y = y
    this.radius = 1.5
    this.maxRadius = 18 + Math.random() * 14
    this.alpha = 0.4 + Math.random() * 0.2
    this.speed = 0.3 + Math.random() * 0.2
    this.dead = false
  }

  update() {
    this.radius += this.speed
    this.alpha -= 0.006
    if (this.radius >= this.maxRadius || this.alpha <= 0) {
      this.dead = true
    }
  }

  draw(c) {
    c.save()
    c.strokeStyle = `rgba(71, 85, 105, ${Math.max(0, this.alpha)})`
    c.lineWidth = 1.0
    c.beginPath()
    c.ellipse(this.x, this.y, this.radius, this.radius * 0.38, 0, 0, Math.PI * 2)
    c.stroke()
    c.restore()
  }
}

/* ========================================================
   Robust Canvas Resolution & Initialization
   ======================================================== */
function resolveCanvas() {
  if (canvasEl && ctx) return canvasEl

  // 1. Try canvasRef
  if (canvasRef.value) {
    if (typeof canvasRef.value.getContext === 'function') {
      canvasEl = canvasRef.value
    } else if (canvasRef.value.$el) {
      if (typeof canvasRef.value.$el.getContext === 'function') {
        canvasEl = canvasRef.value.$el
      } else {
        const inner = canvasRef.value.$el.querySelector('canvas')
        if (inner && typeof inner.getContext === 'function') {
          canvasEl = inner
        }
      }
    } else if (typeof canvasRef.value.querySelector === 'function') {
      const inner = canvasRef.value.querySelector('canvas')
      if (inner && typeof inner.getContext === 'function') {
        canvasEl = inner
      }
    }
  }

  // 2. Try ID lookup
  if (!canvasEl && typeof document !== 'undefined') {
    const byId = document.getElementById('ambient-canvas-elem')
    if (byId) {
      if (typeof byId.getContext === 'function') {
        canvasEl = byId
      } else {
        const inner = byId.querySelector('canvas')
        if (inner && typeof inner.getContext === 'function') {
          canvasEl = inner
        }
      }
    }
  }

  // 3. Fallback querySelector
  if (!canvasEl && typeof document !== 'undefined') {
    canvasEl = document.querySelector('.ambient-backdrop canvas') || document.querySelector('canvas')
  }

  // 4. Ultimate fallback: create native HTML5 canvas dynamically
  if (!canvasEl && typeof document !== 'undefined') {
    const parent = document.querySelector('.ambient-backdrop') || document.body
    canvasEl = document.createElement('canvas')
    canvasEl.id = 'ambient-canvas-elem'
    canvasEl.className = 'ambient-canvas'
    parent.prepend(canvasEl)
  }

  if (canvasEl && typeof canvasEl.getContext === 'function') {
    ctx = canvasEl.getContext('2d')
    canvasEl.style.zIndex = '0'
    canvasEl.style.position = 'absolute'
    canvasEl.style.top = '0px'
    canvasEl.style.left = '0px'
    canvasEl.style.pointerEvents = 'none'
  }

  return canvasEl
}

function initElements() {
  stars.length = 0
  meteors.length = 0
  raindrops.length = 0
  glassDroplets.length = 0
  ripples.length = 0

  const starCount = Math.min(140, Math.floor((width * height) / 9000))
  for (let i = 0; i < starCount; i++) {
    stars.push(new Star(width, height))
  }

  const rainCount = Math.min(80, Math.floor((width * height) / 14000))
  for (let i = 0; i < rainCount; i++) {
    raindrops.push(new Raindrop(width, height))
  }

  const dropletCount = 28
  for (let i = 0; i < dropletCount; i++) {
    glassDroplets.push(new GlassDroplet(width, height))
  }

  // Initial meteor spawns gently after user settles in
  nextMeteorTime = Date.now() + 3500
}

function handleResize() {
  const canvas = resolveCanvas()
  if (!canvas || !ctx) return

  dpr = window.devicePixelRatio || 1
  width = window.innerWidth
  height = window.innerHeight

  canvas.width = width * dpr
  canvas.height = height * dpr
  canvas.style.width = width + 'px'
  canvas.style.height = height + 'px'

  ctx.setTransform(1, 0, 0, 1, 0, 0)
  ctx.scale(dpr, dpr)
  initElements()
}

function handleMouseMove(e) {
  mouse.targetX = e.clientX
  mouse.targetY = e.clientY
  mouse.active = true

  if (appStore.theme === 'light' && Math.random() < 0.22) {
    ripples.push(new Ripple(e.clientX, e.clientY))
  }
}

function handleMouseLeave() {
  mouse.active = false
}

function render() {
  if (!ctx) {
    resolveCanvas()
    if (!ctx) {
      animId = requestAnimationFrame(render)
      return
    }
  }

  if (mouse.active) {
    mouse.x += (mouse.targetX - mouse.x) * 0.08
    mouse.y += (mouse.targetY - mouse.y) * 0.08
  }

  ctx.clearRect(0, 0, width, height)

  const isDark = appStore.theme === 'dark'

  if (isDark) {
    // ----------------------------------------------------
    // Dark: Cosmic Starfield + Ambient Glow + Meteors
    // ----------------------------------------------------
    if (mouse.active && mouse.x > 0) {
      const aura = ctx.createRadialGradient(mouse.x, mouse.y, 10, mouse.x, mouse.y, 240)
      aura.addColorStop(0, 'rgba(56, 189, 248, 0.1)')
      aura.addColorStop(0.6, 'rgba(129, 140, 248, 0.04)')
      aura.addColorStop(1, 'rgba(0, 0, 0, 0)')
      ctx.fillStyle = aura
      ctx.beginPath()
      ctx.arc(mouse.x, mouse.y, 240, 0, Math.PI * 2)
      ctx.fill()
    }

    // Stars
    for (let i = 0; i < stars.length; i++) {
      stars[i].update(width, height)
      stars[i].draw(ctx)
    }

    // Meteors: Peaceful and rare (once every 9 to 17 seconds)
    const now = Date.now()
    if (now > nextMeteorTime) {
      meteors.push(new Meteor(width, height))
      nextMeteorTime = now + 9000 + Math.random() * 8000
    }

    for (let i = meteors.length - 1; i >= 0; i--) {
      meteors[i].update()
      meteors[i].draw(ctx)
      if (meteors[i].dead) {
        meteors.splice(i, 1)
      }
    }
  } else {
    // ----------------------------------------------------
    // Light: Frosted Glass + Raindrops + Water Droplets
    // ----------------------------------------------------
    if (mouse.active && mouse.x > 0) {
      const aura = ctx.createRadialGradient(mouse.x, mouse.y, 10, mouse.x, mouse.y, 220)
      aura.addColorStop(0, 'rgba(59, 130, 246, 0.08)')
      aura.addColorStop(1, 'rgba(255, 255, 255, 0)')
      ctx.fillStyle = aura
      ctx.beginPath()
      ctx.arc(mouse.x, mouse.y, 220, 0, Math.PI * 2)
      ctx.fill()
    }

    for (let i = 0; i < raindrops.length; i++) {
      raindrops[i].update(width, height)
      raindrops[i].draw(ctx)
    }

    for (let i = ripples.length - 1; i >= 0; i--) {
      ripples[i].update()
      ripples[i].draw(ctx)
      if (ripples[i].dead) {
        ripples.splice(i, 1)
      }
    }

    for (let i = 0; i < glassDroplets.length; i++) {
      glassDroplets[i].update(width, height)
      glassDroplets[i].draw(ctx)
    }
  }

  animId = requestAnimationFrame(render)
}

function getMpStarStyle(n) {
  const left = ((n * 37) % 96) + 2
  const top = ((n * 61) % 94) + 3
  const delay = (n * 0.17) % 2.5
  const duration = 1.6 + ((n * 0.13) % 2)
  const size = 1.2 + ((n * 0.28) % 3)
  return {
    left: `${left}%`,
    top: `${top}%`,
    width: `${size}px`,
    height: `${size}px`,
    animationDelay: `${delay}s`,
    animationDuration: `${duration}s`
  }
}

function getMpRainStyle(n) {
  const left = ((n * 31) % 96) + 2
  const delay = (n * 0.15) % 1.8
  const duration = 1.0 + ((n * 0.11) % 1.0)
  return {
    left: `${left}%`,
    animationDelay: `${delay}s`,
    animationDuration: `${duration}s`
  }
}

watch(
  () => appStore.theme,
  () => {
    initElements()
  }
)

onMounted(() => {
  // #ifdef H5
  window.addEventListener('resize', handleResize)
  window.addEventListener('mousemove', handleMouseMove, { passive: true })
  document.addEventListener('mouseleave', handleMouseLeave)
  nextTick(() => {
    setTimeout(() => {
      handleResize()
      animId = requestAnimationFrame(render)
    }, 50)
  })
  // #endif
})

onUnmounted(() => {
  // #ifdef H5
  if (animId) cancelAnimationFrame(animId)
  window.removeEventListener('resize', handleResize)
  window.removeEventListener('mousemove', handleMouseMove)
  document.removeEventListener('mouseleave', handleMouseLeave)
  // #endif
})
</script>

<style scoped>
.ambient-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  pointer-events: none;
  z-index: 0;
  overflow: hidden;
  transition: background-color 0.4s ease;
}

.ambient-backdrop.dark {
  background: radial-gradient(ellipse at 50% 10%, #151828 0%, #08090f 100%);
}

.ambient-backdrop.light {
  background: radial-gradient(ellipse at 50% 0%, #f1f5f9 0%, #e2e8f0 100%);
}

.ambient-canvas {
  position: absolute !important;
  top: 0 !important;
  left: 0 !important;
  width: 100vw !important;
  height: 100vh !important;
  pointer-events: none !important;
  display: block !important;
  z-index: 0 !important;
}

/* ========================================================
   CSS Fallback Animations for Mini-Program
   ======================================================== */
.mp-ambient-fx {
  width: 100%;
  height: 100%;
  position: relative;
}

.mp-star {
  position: absolute;
  background: #ffffff;
  border-radius: 50%;
  animation: mpTwinkle ease-in-out infinite;
}

@keyframes mpTwinkle {
  0%, 100% { opacity: 0.25; transform: scale(0.8); }
  50% { opacity: 1; transform: scale(1.2); }
}

.mp-meteor {
  position: absolute;
  height: 2px;
  background: linear-gradient(90deg, rgba(255, 255, 255, 0.95), transparent);
  transform: rotate(-40deg);
  opacity: 0;
}

.mp-meteor-1 {
  width: 150px;
  top: 15%;
  right: 20%;
  animation: mpShoot 4.5s ease-in-out infinite;
}

.mp-meteor-2 {
  width: 200px;
  top: 30%;
  right: 10%;
  animation: mpShoot 6s 2s ease-in-out infinite;
}

@keyframes mpShoot {
  0% { transform: translate(120px, -120px) rotate(-40deg); opacity: 0; }
  10% { opacity: 1; }
  35% { transform: translate(-350px, 350px) rotate(-40deg); opacity: 0; }
  100% { opacity: 0; }
}

.mp-raindrop {
  position: absolute;
  top: -40px;
  width: 1.5px;
  height: 28px;
  background: linear-gradient(180deg, transparent, rgba(71, 85, 105, 0.65));
  animation: mpRain linear infinite;
}

@keyframes mpRain {
  0% { transform: translateY(-40px); opacity: 0; }
  20% { opacity: 0.7; }
  90% { opacity: 0.7; }
  100% { transform: translateY(110vh); opacity: 0; }
}
</style>
