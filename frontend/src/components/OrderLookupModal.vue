<template>
  <transition name="fade">
    <div v-if="modelValue" class="modal-overlay" @click.self="close">
      <div class="modal-box lookup-modal">
        <div class="modal-header">
          <div class="modal-title-wrap">
            <span class="modal-badge-icon">🔍</span>
            <h3>{{ isZh ? '订单与下载找回' : 'Retrieve Orders & Downloads' }}</h3>
          </div>
          <button class="close-btn" @click="close">✕</button>
        </div>

        <div class="modal-body">
          <p class="lookup-tip">
            {{ isZh 
              ? '如果您曾以游客身份购买或忘记登录，输入付款时填写的邮箱，即可快速找回所有已购壁纸与网盘提取码。' 
              : 'Enter your checkout email to retrieve your purchased wallpapers and cloud drive access links anytime.' 
            }}
          </p>

          <form @submit.prevent="handleSearch" class="lookup-form">
            <div class="form-group">
              <label class="form-label">
                <span>📧 {{ isZh ? '支付邮箱' : 'Checkout Email' }} *</span>
              </label>
              <input
                v-model="email"
                type="email"
                required
                class="form-input"
                :placeholder="isZh ? '例如: your_name@gmail.com' : 'e.g. your_name@gmail.com'"
              />
            </div>

            <div class="form-group">
              <label class="form-label">
                <span>🔖 {{ isZh ? '订单编号 / PayPal 账单号 (选填)' : 'Order ID / PayPal ID (Optional)' }}</span>
              </label>
              <input
                v-model="orderId"
                type="text"
                class="form-input"
                :placeholder="isZh ? '留空则查询该邮箱下的所有已购订单' : 'Leave empty to search all orders'"
              />
            </div>

            <button type="submit" class="btn btn-primary btn-block" :disabled="loading">
              <span v-if="loading" class="spinner"></span>
              <span v-else>🔍 {{ isZh ? '立即查询已购壁纸' : 'Find My Downloads' }}</span>
            </button>
          </form>

          <!-- Error Alert -->
          <div v-if="errorMsg" class="alert-error">
            ⚠️ {{ errorMsg }}
          </div>

          <!-- Results Section -->
          <div v-if="searched && orders.length > 0" class="results-section">
            <h4 class="results-title">
              🎉 {{ isZh ? `为您找到 ${orders.length} 笔已支付订单` : `Found ${orders.length} paid orders` }}
            </h4>

            <div class="order-cards-list">
              <div v-for="order in orders" :key="order.id" class="order-result-card">
                <div class="card-top">
                  <div class="order-thumb" v-if="order.cover_image">
                    <img :src="order.cover_image" :alt="order.product_title" />
                  </div>
                  <div class="order-meta">
                    <div class="order-type-badge">
                      {{ order.order_type === 'membership' ? '👑 VIP 会员' : '🖼️ 单品壁纸' }}
                    </div>
                    <div class="order-name">{{ order.product_title || '精选壁纸' }}</div>
                    <div class="order-sub">
                      <span>{{ formatDate(order.created_at) }}</span>
                      <span class="price">${{ order.amount }} USD</span>
                    </div>
                  </div>
                </div>

                <!-- Download Links Area -->
                <div v-if="order.download_links && order.download_links.length > 0" class="download-links-area">
                  <div class="links-title">
                    📥 {{ isZh ? '网盘提取链接与提取码' : 'Cloud Drive Download Links' }}
                  </div>
                  <div class="drive-btn-grid">
                    <div v-for="(dl, idx) in order.download_links" :key="idx" class="drive-item-row">
                      <div class="drive-info">
                        <span class="drive-tag" :class="dl.platform">
                          {{ getPlatformLabel(dl.platform) }}
                        </span>
                        <span v-if="dl.extract_code" class="code-tag">
                          {{ isZh ? '提取码:' : 'Code:' }} <strong>{{ dl.extract_code }}</strong>
                        </span>
                      </div>
                      <div class="drive-actions">
                        <button
                          v-if="dl.extract_code"
                          type="button"
                          class="btn btn-outline btn-xs"
                          @click="copyText(dl.extract_code)"
                        >
                          📋 {{ isZh ? '复制码' : 'Copy' }}
                        </button>
                        <a
                          :href="dl.url"
                          target="_blank"
                          rel="noopener noreferrer"
                          class="btn btn-primary btn-xs"
                        >
                          🚀 {{ isZh ? '直达网盘' : 'Open' }}
                        </a>
                      </div>
                    </div>
                  </div>
                </div>

                <div v-else-if="order.order_type === 'membership'" class="membership-active-tip">
                  ✨ {{ isZh ? '会员权益已生效，登录该邮箱即可在商城全场免费无限下载！' : 'VIP active. Log in to download any wallpaper for free!' }}
                </div>
              </div>
            </div>
          </div>

          <div v-else-if="searched && orders.length === 0" class="empty-results">
            <span class="empty-icon">📭</span>
            <p>{{ isZh ? '未找到该邮箱名下的已支付订单。' : 'No paid orders found for this email.' }}</p>
          </div>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { orderAPI } from '@/api'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue'])

const { locale } = useI18n()
const isZh = computed(() => locale.value === 'zh')

const email = ref('')
const orderId = ref('')
const loading = ref(false)
const searched = ref(false)
const orders = ref([])
const errorMsg = ref('')

function close() {
  emit('update:modelValue', false)
}

async function handleSearch() {
  if (!email.value) return
  loading.value = true
  errorMsg.value = ''
  searched.value = false
  orders.value = []

  try {
    const res = await orderAPI.lookup({
      email: email.value.trim(),
      order_id: orderId.value.trim()
    })
    orders.value = res.orders || []
    searched.value = true
  } catch (err) {
    errorMsg.value = err.message || (isZh.value ? '未查询到相关支付订单' : 'Order not found')
    searched.value = true
  } finally {
    loading.value = false
  }
}

function getPlatformLabel(platform) {
  const map = {
    baidu: isZh.value ? '百度网盘' : 'Baidu Drive',
    quark: isZh.value ? '夸克网盘' : 'Quark Drive',
    google: isZh.value ? '谷歌网盘' : 'Google Drive'
  }
  return map[platform] || platform
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  try {
    const d = new Date(dateStr)
    return d.toLocaleDateString(isZh.value ? 'zh-CN' : 'en-US', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit'
    })
  } catch (e) {
    return dateStr
  }
}

function copyText(txt) {
  if (!txt) return
  if (navigator.clipboard) {
    navigator.clipboard.writeText(txt)
  }
  alert(isZh.value ? `已复制提取码: ${txt}` : `Copied code: ${txt}`)
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.7);
  backdrop-filter: blur(8px);
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
}

.lookup-modal {
  background: var(--bg-card, #1a1b26);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.1));
  border-radius: 18px;
  width: 100%;
  max-width: 540px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.4);
}

.modal-header {
  padding: 18px 24px;
  border-bottom: 1px solid var(--border-color, rgba(255, 255, 255, 0.08));
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.modal-title-wrap {
  display: flex;
  align-items: center;
  gap: 10px;
}

.modal-title-wrap h3 {
  margin: 0;
  font-size: 1.15rem;
  font-weight: 700;
  color: #fff;
}

.close-btn {
  background: transparent;
  border: none;
  color: #888;
  font-size: 1.2rem;
  cursor: pointer;
  padding: 4px;
  transition: color 0.2s;
}

.close-btn:hover {
  color: #fff;
}

.modal-body {
  padding: 24px;
  overflow-y: auto;
  flex: 1;
}

.lookup-tip {
  font-size: 0.88rem;
  color: #9ca3af;
  margin-bottom: 20px;
  line-height: 1.5;
  background: rgba(255, 255, 255, 0.03);
  padding: 12px 16px;
  border-radius: 12px;
  border-left: 3px solid #6366f1;
}

.lookup-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-label {
  font-size: 0.84rem;
  font-weight: 600;
  color: #e5e7eb;
}

.form-input {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 10px;
  padding: 11px 14px;
  color: #fff;
  font-size: 0.95rem;
  outline: none;
  transition: border-color 0.2s;
}

.form-input:focus {
  border-color: #6366f1;
}

.alert-error {
  background: rgba(239, 68, 68, 0.15);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: #f87171;
  padding: 12px;
  border-radius: 10px;
  font-size: 0.88rem;
  margin-bottom: 16px;
}

.results-section {
  margin-top: 10px;
}

.results-title {
  font-size: 1rem;
  font-weight: 700;
  color: #10b981;
  margin-bottom: 14px;
}

.order-cards-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.order-result-card {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  padding: 14px;
}

.card-top {
  display: flex;
  gap: 12px;
  align-items: center;
  margin-bottom: 12px;
}

.order-thumb {
  width: 52px;
  height: 52px;
  border-radius: 8px;
  overflow: hidden;
  flex-shrink: 0;
  background: #222;
}

.order-thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.order-meta {
  flex: 1;
  min-width: 0;
}

.order-type-badge {
  font-size: 0.72rem;
  color: #a78bfa;
  font-weight: 600;
  margin-bottom: 2px;
}

.order-name {
  font-size: 0.95rem;
  font-weight: 600;
  color: #fff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.order-sub {
  font-size: 0.78rem;
  color: #9ca3af;
  display: flex;
  gap: 12px;
  margin-top: 2px;
}

.order-sub .price {
  color: #34d399;
  font-weight: 600;
}

.download-links-area {
  background: rgba(0, 0, 0, 0.25);
  border-radius: 10px;
  padding: 12px;
  margin-top: 10px;
}

.links-title {
  font-size: 0.8rem;
  font-weight: 600;
  color: #d1d5db;
  margin-bottom: 10px;
}

.drive-btn-grid {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.drive-item-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: rgba(255, 255, 255, 0.04);
  padding: 8px 12px;
  border-radius: 8px;
  gap: 8px;
}

.drive-info {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.drive-tag {
  font-size: 0.75rem;
  padding: 2px 8px;
  border-radius: 6px;
  font-weight: 600;
}

.drive-tag.baidu { background: rgba(33, 115, 243, 0.2); color: #60a5fa; }
.drive-tag.quark { background: rgba(245, 158, 11, 0.2); color: #fbbf24; }
.drive-tag.google { background: rgba(16, 185, 129, 0.2); color: #34d399; }

.code-tag {
  font-size: 0.78rem;
  color: #f3f4f6;
  background: rgba(255, 255, 255, 0.1);
  padding: 2px 8px;
  border-radius: 6px;
}

.drive-actions {
  display: flex;
  gap: 6px;
}

.membership-active-tip {
  font-size: 0.84rem;
  color: #fbbf24;
  background: rgba(245, 158, 11, 0.1);
  padding: 10px;
  border-radius: 8px;
  margin-top: 8px;
}

.empty-results {
  text-align: center;
  padding: 30px 10px;
  color: #9ca3af;
}

.empty-icon {
  font-size: 2.2rem;
  display: block;
  margin-bottom: 8px;
}

.btn-xs {
  padding: 4px 10px;
  font-size: 0.75rem;
  border-radius: 6px;
}
</style>
