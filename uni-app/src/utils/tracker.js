// Visitor traffic & dwell duration tracker
import request from './request'

let currentVisitId = ''
let startTime = Date.now()
let heartbeatTimer = null

function getDeviceInfo() {
  let device = 'Desktop'
  let os = 'Unknown'
  let browser = 'Browser'
  let ua = ''

  if (typeof navigator !== 'undefined') {
    ua = navigator.userAgent || ''
    const uaLower = ua.toLowerCase()

    if (uaLower.includes('ipad') || uaLower.includes('tablet')) {
      device = 'Tablet'
    } else if (uaLower.includes('mobile') || uaLower.includes('android') || uaLower.includes('iphone')) {
      device = 'Mobile'
    } else {
      device = 'Desktop'
    }

    if (uaLower.includes('windows')) os = 'Windows'
    else if (uaLower.includes('iphone') || uaLower.includes('ipad') || uaLower.includes('ios')) os = 'iOS'
    else if (uaLower.includes('android')) os = 'Android'
    else if (uaLower.includes('macintosh') || uaLower.includes('mac os')) os = 'macOS'
    else if (uaLower.includes('linux')) os = 'Linux'

    if (uaLower.includes('micromessenger')) browser = 'WeChat'
    else if (uaLower.includes('edg')) browser = 'Edge'
    else if (uaLower.includes('chrome')) browser = 'Chrome'
    else if (uaLower.includes('safari')) browser = 'Safari'
    else if (uaLower.includes('firefox')) browser = 'Firefox'
  }

  return { device, os, browser, user_agent: ua }
}

export function startTracking(path = window.location?.hash || '/') {
  // If already tracking, report previous page duration
  stopHeartbeat()
  startTime = Date.now()

  const info = getDeviceInfo()
  const lang = uni.getStorageSync('app_lang') || ''

  let userEmail = ''
  let userId = ''
  try {
    const raw = uni.getStorageSync('user_info')
    if (raw) {
      const u = typeof raw === 'string' ? JSON.parse(raw) : raw
      userEmail = u.email || ''
      userId = u.id || ''
    }
  } catch (e) {}

  request.post('/stats/track', {
    path,
    language: lang,
    device: info.device,
    os: info.os,
    browser: info.browser,
    user_agent: info.user_agent,
    user_email: userEmail,
    user_id: userId
  }).then(res => {
    if (res && res.visit_id) {
      currentVisitId = res.visit_id
      startHeartbeat()
    }
  }).catch(() => {})
}

function startHeartbeat() {
  if (heartbeatTimer) clearInterval(heartbeatTimer)
  heartbeatTimer = setInterval(() => {
    reportDuration()
  }, 15000)
}

function stopHeartbeat() {
  if (heartbeatTimer) {
    clearInterval(heartbeatTimer)
    heartbeatTimer = null
  }
}

function reportDuration() {
  if (!currentVisitId) return
  const seconds = Math.floor((Date.now() - startTime) / 1000)
  if (seconds <= 0) return

  // In browser, try sendBeacon on unload if needed, otherwise normal request
  request.post('/stats/heartbeat', {
    visit_id: currentVisitId,
    duration: seconds
  }).catch(() => {})
}

// Lifecycle listeners for web
if (typeof window !== 'undefined') {
  window.addEventListener('beforeunload', () => {
    reportDuration()
  })

  document.addEventListener('visibilitychange', () => {
    if (document.visibilityState === 'hidden') {
      reportDuration()
    }
  })
}

export default {
  startTracking,
  reportDuration
}
