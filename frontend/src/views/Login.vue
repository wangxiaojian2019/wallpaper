<template>
  <div class="auth-page">
    <div class="auth-container">
      <div class="auth-card">
        <!-- Logo -->
        <div class="auth-logo">
          <span class="logo-icon">✦</span>
          <span>WallArt</span>
        </div>

        <h2 class="auth-title">{{ t('auth.login') }}</h2>
        <p class="auth-sub">{{ locale === 'zh' ? '欢迎回来' : 'Welcome back' }}</p>

        <form @submit.prevent="handleLogin" class="auth-form">
          <div class="form-group">
            <label class="form-label">{{ t('auth.email') }}</label>
            <input
              v-model="form.email"
              type="email"
              class="form-input"
              :class="{ error: errors.email }"
              :placeholder="t('auth.email')"
              autocomplete="email"
            />
            <span v-if="errors.email" class="form-error">{{ errors.email }}</span>
          </div>

          <div class="form-group">
            <label class="form-label">{{ t('auth.password') }}</label>
            <div class="input-with-icon">
              <input
                v-model="form.password"
                :type="showPwd ? 'text' : 'password'"
                class="form-input"
                :class="{ error: errors.password }"
                :placeholder="t('auth.password')"
                autocomplete="current-password"
              />
              <button type="button" class="eye-btn" @click="showPwd = !showPwd">
                {{ showPwd ? '🙈' : '👁️' }}
              </button>
            </div>
            <span v-if="errors.password" class="form-error">{{ errors.password }}</span>
          </div>

          <div class="auth-options">
            <a href="#" class="forgot-link">{{ t('auth.forgot') }}</a>
          </div>

          <button type="submit" class="btn btn-primary btn-block btn-lg" :disabled="loading">
            <span v-if="loading" class="spinner"></span>
            {{ loading ? t('common.loading') : t('auth.login') }}
          </button>
        </form>

        <div class="auth-footer">
          {{ t('auth.no_account') }}
          <router-link to="/register">{{ t('auth.register') }}</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/store'
import { authAPI } from '@/api'
import { useToast } from '@/composables/useToast'

const { t, locale } = useI18n()
const router = useRouter()
const route = useRoute()
const auth = useAuthStore()
const { success, error } = useToast()

const loading = ref(false)
const showPwd = ref(false)
const form = reactive({ email: '', password: '' })
const errors = reactive({ email: '', password: '' })

function validate() {
  errors.email = ''
  errors.password = ''
  let ok = true
  if (!form.email) { errors.email = t('auth.email_required'); ok = false }
  if (!form.password) { errors.password = t('auth.password_required'); ok = false }
  else if (form.password.length < 6) { errors.password = t('auth.password_min'); ok = false }
  return ok
}

async function handleLogin() {
  if (!validate()) return
  loading.value = true
  try {
    const res = await authAPI.login({ email: form.email, password: form.password })
    auth.setAuth(res.user, res.token)
    success(t('auth.login_success'))
    const defaultTarget = (res.user?.role === 'admin' || window.location.hostname.includes('admin')) ? '/admin' : '/'
    const redirect = route.query.redirect || defaultTarget
    router.push(redirect)
  } catch (e) {
    error(e.message || t('common.error'))
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-page {
  min-height: calc(100vh - var(--nav-h));
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--g-hero);
  padding: 40px 20px;
}

.auth-container {
  width: 100%;
  max-width: 420px;
}

.auth-card {
  background: var(--c-surface);
  border-radius: var(--r-2xl);
  padding: 40px;
  box-shadow: var(--shadow-lg);
}

.auth-logo {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 1.3rem;
  font-weight: 800;
  letter-spacing: -0.03em;
  margin-bottom: 28px;
  justify-content: center;
}
.logo-icon { color: var(--c-accent); }

.auth-title {
  text-align: center;
  margin-bottom: 6px;
  font-size: 1.6rem;
}

.auth-sub {
  text-align: center;
  font-size: 0.9rem;
  color: var(--c-text-secondary);
  margin-bottom: 28px;
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.input-with-icon {
  position: relative;
}

.input-with-icon .form-input {
  padding-right: 44px;
}

.eye-btn {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  cursor: pointer;
  font-size: 1rem;
  padding: 0;
}

.auth-options {
  display: flex;
  justify-content: flex-end;
}

.forgot-link {
  font-size: 0.82rem;
  color: var(--c-text-tertiary);
}
.forgot-link:hover { color: var(--c-accent); }

.auth-footer {
  margin-top: 20px;
  text-align: center;
  font-size: 0.88rem;
  color: var(--c-text-secondary);
}
.auth-footer a {
  font-weight: 600;
  margin-left: 4px;
}
</style>
