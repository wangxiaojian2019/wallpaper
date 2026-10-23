<template>
  <div class="payment-result-page">
    <div class="result-card">
      <div class="icon-circle cancel">✕</div>
      <h2>{{ t('payment.cancel') }}</h2>
      <p class="text-secondary">
        {{ isZh ? '您已取消本次支付。您的账户未扣款，您可以随时重新发起购买。' : 'Payment was cancelled. You have not been charged and can retry anytime.' }}
      </p>

      <div class="action-buttons">
        <router-link to="/" class="btn btn-primary">
          {{ isZh ? '返回商城首页' : 'Return to Home' }}
        </router-link>
        <router-link to="/membership" class="btn btn-outline">
          {{ isZh ? '了解会员特权' : 'Explore Membership' }}
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()
const isZh = computed(() => locale.value === 'zh')
</script>

<style scoped>
.payment-result-page {
  min-height: 80vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
}

.result-card {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-xl);
  padding: 48px;
  max-width: 540px;
  width: 100%;
  text-align: center;
}

.icon-circle.cancel {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
  border: 2px solid rgba(239, 68, 68, 0.25);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  font-weight: 700;
  margin: 0 auto 24px;
}

.result-card h2 {
  font-size: 1.75rem;
  font-weight: 700;
  margin-bottom: 12px;
  color: var(--text-primary);
}

.result-card p {
  line-height: 1.6;
  margin-bottom: 32px;
}

.action-buttons {
  display: flex;
  gap: 16px;
  justify-content: center;
  flex-wrap: wrap;
}
</style>
