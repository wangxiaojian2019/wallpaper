<template>
  <view v-if="visible" class="auth-overlay" @tap.self="close">
    <view class="auth-sheet" :class="[appStore.theme]">
      <!-- Drag handle -->
      <view class="drag-bar"></view>

      <!-- Close icon -->
      <view class="auth-close-btn" @tap="close">
        <text class="close-icon">✕</text>
      </view>

      <!-- Top Icon & Header -->
      <view class="auth-header">
        <view class="header-icon-badge">
          <text class="header-badge-icon">👑</text>
        </view>
        <text class="auth-title">{{ isRegisterMode ? (appStore.lang === 'zh' ? '免费注册会员' : 'Create Account') : (appStore.lang === 'zh' ? '会员登录' : 'Sign In') }}</text>
        <text class="auth-subtitle">
          {{ appStore.lang === 'zh' ? '云端永久保存已购订单 · 换设备不丢失 · 极速下载' : 'Save all purchases permanently to cloud. Access anywhere.' }}
        </text>
      </view>

      <!-- Member Benefit Callout (Encourages Guest Conversion) -->
      <view class="benefit-callout">
        <view class="benefit-row">
          <text class="benefit-icon">🛡️</text>
          <text class="benefit-text">
            {{ appStore.lang === 'zh' ? '强烈建议注册账号：游客购买若误清浏览器缓存易丢失提取链接，注册后永久保存在个人中心！' : 'Registered members never lose download links when clearing browser cache.' }}
          </text>
        </view>
      </view>

      <!-- Mode Switcher Tabs -->
      <view class="auth-tabs">
        <view
          class="auth-tab"
          :class="{ active: !isRegisterMode }"
          @tap="isRegisterMode = false"
        >
          <text class="tab-text">{{ appStore.lang === 'zh' ? '登录' : 'Sign In' }}</text>
        </view>
        <view
          class="auth-tab"
          :class="{ active: isRegisterMode }"
          @tap="isRegisterMode = true"
        >
          <text class="tab-text">{{ appStore.lang === 'zh' ? '免费注册' : 'Register' }}</text>
          <text class="tab-badge-new">Free</text>
        </view>
      </view>

      <!-- Form Inputs -->
      <view class="form-body">
        <view class="input-group">
          <text class="input-label">{{ appStore.lang === 'zh' ? '电子邮箱' : 'Email Address' }}</text>
          <view class="input-wrap">
            <text class="input-icon">✉️</text>
            <input
              v-model="email"
              type="text"
              :placeholder="appStore.lang === 'zh' ? '请输入常用邮箱 (如 name@domain.com)' : 'name@example.com'"
              class="auth-input"
            />
          </view>
        </view>

        <view class="input-group">
          <text class="input-label">{{ appStore.lang === 'zh' ? '登录密码' : 'Password' }}</text>
          <view class="input-wrap">
            <text class="input-icon">🔒</text>
            <input
              v-model="password"
              :password="true"
              :placeholder="appStore.lang === 'zh' ? '请输入密码 (不少于6位)' : 'At least 6 characters'"
              class="auth-input"
            />
          </view>
        </view>

        <view v-if="isRegisterMode" class="input-group">
          <text class="input-label">{{ appStore.lang === 'zh' ? '确认密码' : 'Confirm Password' }}</text>
          <view class="input-wrap">
            <text class="input-icon">🔐</text>
            <input
              v-model="confirmPassword"
              :password="true"
              :placeholder="appStore.lang === 'zh' ? '请再次输入密码确认' : 'Re-enter password'"
              class="auth-input"
            />
          </view>
        </view>

        <!-- Error Tip -->
        <view v-if="errorMsg" class="error-box">
          <text class="error-icon">⚠️</text>
          <text class="error-text">{{ errorMsg }}</text>
        </view>

        <!-- Submit Button -->
        <button
          class="btn-submit"
          :loading="loading"
          :disabled="loading"
          @tap="handleSubmit"
        >
          <text class="btn-text">
            {{ isRegisterMode ? (appStore.lang === 'zh' ? '立即免费注册' : 'Create Free Account') : (appStore.lang === 'zh' ? '立即登录' : 'Sign In') }}
          </text>
        </button>

        <!-- Toggle Mode Link -->
        <view class="toggle-footer" @tap="isRegisterMode = !isRegisterMode">
          <text class="toggle-text">
            {{ isRegisterMode ? (appStore.lang === 'zh' ? '已有账号？点击直接登录' : 'Already have an account? Sign In') : (appStore.lang === 'zh' ? '还没有账号？3秒免费注册' : 'No account? Create one for free') }}
          </text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useUserStore } from '@/store/user'
import { useAppStore } from '@/store/app'

const props = defineProps({
  visible: Boolean,
  initialMode: {
    type: String,
    default: 'login' // 'login' | 'register'
  },
  initialEmail: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['close', 'success'])

const userStore = useUserStore()
const appStore = useAppStore()

const isRegisterMode = ref(props.initialMode === 'register')
const email = ref(props.initialEmail || '')
const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const errorMsg = ref('')

watch(() => props.visible, (val) => {
  if (val) {
    errorMsg.value = ''
    password.value = ''
    confirmPassword.value = ''
    if (props.initialEmail) email.value = props.initialEmail
    isRegisterMode.value = props.initialMode === 'register'
  }
})

function close() {
  emit('close')
}

async function handleSubmit() {
  errorMsg.value = ''
  const em = email.value.trim()
  const pwd = password.value.trim()

  if (!em || !em.includes('@')) {
    errorMsg.value = appStore.lang === 'zh' ? '请输入有效的邮箱地址' : 'Please enter a valid email address'
    return
  }

  if (!pwd || pwd.length < 6) {
    errorMsg.value = appStore.lang === 'zh' ? '密码长度不能少于 6 位' : 'Password must be at least 6 characters'
    return
  }

  if (isRegisterMode.value && pwd !== confirmPassword.value.trim()) {
    errorMsg.value = appStore.lang === 'zh' ? '两次输入的密码不一致' : 'Passwords do not match'
    return
  }

  loading.value = true
  try {
    if (isRegisterMode.value) {
      await userStore.register({ email: em, password: pwd })
      uni.showToast({
        title: appStore.lang === 'zh' ? '注册并登录成功！' : 'Registered successfully!',
        icon: 'success'
      })
    } else {
      await userStore.login({ email: em, password: pwd })
      uni.showToast({
        title: appStore.lang === 'zh' ? '登录成功' : 'Signed in successfully',
        icon: 'success'
      })
    }
    emit('success', userStore.userInfo)
    close()
  } catch (err) {
    errorMsg.value = err.message || (appStore.lang === 'zh' ? '操作失败，请重试' : 'Operation failed, please try again')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.72);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  z-index: 9999;
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.auth-sheet {
  width: 100%;
  max-width: 480px;
  max-height: 90vh;
  overflow-y: auto;
  border-radius: 28px 28px 0 0;
  padding: 16px 24px 32px;
  box-sizing: border-box;
  position: relative;
  animation: slideUp 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes slideUp {
  from {
    transform: translateY(100%);
  }
  to {
    transform: translateY(0);
  }
}

.auth-sheet.dark {
  background: #151722;
  border-top: 1px solid rgba(255, 255, 255, 0.12);
  color: #f1f5f9;
  box-shadow: 0 -12px 48px rgba(0, 0, 0, 0.65);
}

.auth-sheet.light {
  background: #ffffff;
  border-top: 1px solid rgba(0, 0, 0, 0.08);
  color: #0f172a;
  box-shadow: 0 -8px 36px rgba(0, 0, 0, 0.16);
}

.drag-bar {
  width: 40px;
  height: 4px;
  background: rgba(148, 163, 184, 0.4);
  border-radius: 999px;
  margin: 0 auto 14px;
}

.auth-close-btn {
  position: absolute;
  top: 16px;
  right: 18px;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: rgba(148, 163, 184, 0.14);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.close-icon {
  font-size: 14px;
  color: #94a3b8;
}

.auth-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  margin-bottom: 16px;
}

.header-icon-badge {
  width: 52px;
  height: 52px;
  border-radius: 50%;
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.2), rgba(217, 119, 6, 0.2));
  border: 1px solid rgba(245, 158, 11, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 10px;
}

.header-badge-icon {
  font-size: 24px;
}

.auth-title {
  font-size: 20px;
  font-weight: 800;
  letter-spacing: 0.3px;
  margin-bottom: 4px;
}

.auth-subtitle {
  font-size: 12px;
  opacity: 0.65;
  line-height: 1.4;
}

/* Benefit Callout Card */
.benefit-callout {
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.12), rgba(217, 119, 6, 0.08));
  border: 1px solid rgba(245, 158, 11, 0.3);
  border-radius: 14px;
  padding: 10px 14px;
  margin-bottom: 16px;
}

.benefit-row {
  display: flex;
  align-items: flex-start;
  gap: 8px;
}

.benefit-icon {
  font-size: 14px;
  flex-shrink: 0;
  margin-top: 1px;
}

.benefit-text {
  font-size: 12px;
  line-height: 1.45;
  color: #fbbf24;
  font-weight: 600;
}

/* Tabs */
.auth-tabs {
  display: flex;
  background: rgba(148, 163, 184, 0.12);
  border-radius: 12px;
  padding: 3px;
  margin-bottom: 18px;
}

.auth-tab {
  flex: 1;
  text-align: center;
  padding: 8px 0;
  border-radius: 10px;
  cursor: pointer;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  transition: all 0.2s ease;
}

.auth-tab.active {
  background: #2563eb;
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.3);
}

.auth-tab.active .tab-text {
  color: #ffffff;
  font-weight: 700;
}

.tab-text {
  font-size: 13px;
  font-weight: 600;
  opacity: 0.75;
}

.auth-tab.active .tab-text {
  opacity: 1;
}

.tab-badge-new {
  font-size: 9px;
  font-weight: 800;
  background: #10b981;
  color: #ffffff;
  padding: 1px 5px;
  border-radius: 999px;
}

/* Form */
.form-body {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.input-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.input-label {
  font-size: 12px;
  font-weight: 600;
  opacity: 0.7;
}

.input-wrap {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 46px;
  border-radius: 14px;
  padding: 0 14px;
  border: 1px solid rgba(148, 163, 184, 0.25);
  background: rgba(148, 163, 184, 0.08);
}

.input-icon {
  font-size: 14px;
  opacity: 0.6;
}

.auth-input {
  flex: 1;
  height: 100%;
  font-size: 14px;
  background: transparent;
  border: none;
  outline: none;
}

.error-box {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  border-radius: 10px;
  background: rgba(239, 68, 68, 0.14);
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.error-icon {
  font-size: 13px;
}

.error-text {
  font-size: 12px;
  color: #f87171;
  font-weight: 600;
}

.btn-submit {
  height: 48px;
  border-radius: 14px;
  background: linear-gradient(135deg, #2563eb, #3b82f6);
  color: #ffffff;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 6px 20px rgba(37, 99, 235, 0.35);
  margin-top: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-submit:active {
  transform: scale(0.98);
}

.btn-text {
  font-size: 15px;
  font-weight: 700;
  color: #ffffff;
}

.toggle-footer {
  text-align: center;
  padding: 6px 0;
  cursor: pointer;
}

.toggle-text {
  font-size: 12px;
  color: #38bdf8;
  font-weight: 600;
}
</style>
