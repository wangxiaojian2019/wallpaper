<template>
  <div class="admin-orders">
    <!-- Header Summary Stats -->
    <div class="orders-summary-grid">
      <div class="summary-stat-card">
        <div class="stat-icon-wrap bg-blue-subtle">📋</div>
        <div class="stat-info">
          <span class="stat-label">{{ locale === 'zh' ? '订单总数' : 'Total Orders' }}</span>
          <strong class="stat-value">{{ total }}</strong>
        </div>
      </div>

      <div class="summary-stat-card clickable" :class="{ active: filterStatus === 'pending_review' }" @click="filterStatus = (filterStatus === 'pending_review' ? '' : 'pending_review'); loadOrders()">
        <div class="stat-icon-wrap bg-amber-subtle">⏳</div>
        <div class="stat-info">
          <div class="stat-label-row">
            <span class="stat-label">{{ locale === 'zh' ? '待审核对账' : 'Pending Review' }}</span>
            <span v-if="pendingReviewTotal > 0" class="stat-badge-pulse">{{ pendingReviewTotal }}</span>
          </div>
          <strong class="stat-value text-amber">{{ pendingReviewTotal }}</strong>
        </div>
      </div>

      <div class="summary-stat-card clickable" :class="{ active: filterStatus === 'paid' }" @click="filterStatus = (filterStatus === 'paid' ? '' : 'paid'); loadOrders()">
        <div class="stat-icon-wrap bg-emerald-subtle">✓</div>
        <div class="stat-info">
          <span class="stat-label">{{ locale === 'zh' ? '已支付完成' : 'Paid Completed' }}</span>
          <strong class="stat-value text-emerald">{{ paidTotal }}</strong>
        </div>
      </div>

      <div class="summary-stat-card">
        <div class="stat-icon-wrap bg-purple-subtle">💰</div>
        <div class="stat-info">
          <span class="stat-label">{{ locale === 'zh' ? '当前页成交额' : 'Page Revenue' }}</span>
          <strong class="stat-value text-base">{{ pageRevenueText }}</strong>
        </div>
      </div>
    </div>

    <!-- Toolbar -->
    <div class="toolbar card">
      <div class="toolbar-left">
        <div class="search-input-wrap">
          <span class="search-icon">🔍</span>
          <input
            v-model="search"
            type="text"
            class="form-input search-input"
            :placeholder="locale === 'zh' ? '搜索订单编号 / 买家邮箱 / 流水号...' : 'Search Order ID / Email / Ref...'"
            @input="handleSearch"
          />
        </div>

        <select v-model="filterStatus" class="form-input form-select" style="width: 150px;" @change="loadOrders">
          <option value="">{{ t('common.all') }} {{ locale === 'zh' ? '状态' : 'Status' }}</option>
          <option value="pending_review">⏳ {{ locale === 'zh' ? '待审核' : 'Pending Review' }}</option>
          <option value="paid">✓ {{ t('admin.paid') }}</option>
          <option value="pending">⚪ {{ locale === 'zh' ? '待付款' : 'Pending Pay' }}</option>
          <option value="cancelled">✕ {{ t('admin.cancelled') }}</option>
        </select>

        <select v-model="filterType" class="form-input form-select" style="width: 140px;" @change="loadOrders">
          <option value="">{{ t('common.all') }} {{ locale === 'zh' ? '类型' : 'Types' }}</option>
          <option value="single">🖼️ {{ locale === 'zh' ? '单购壁纸' : 'Single Wallpaper' }}</option>
          <option value="membership">👑 {{ locale === 'zh' ? '会员订阅' : 'Membership' }}</option>
        </select>
      </div>

      <div class="toolbar-right">
        <button class="btn btn-outline btn-refresh" :disabled="loading" @click="loadOrders">
          <span :class="{ spinning: loading }">🔄</span> {{ locale === 'zh' ? '刷新' : 'Refresh' }}
        </button>
      </div>
    </div>

    <!-- Table -->
    <div class="orders-table-wrapper card">
      <div class="table-scroll-container">
        <div class="table-header">
          <span class="col-order-id">{{ locale === 'zh' ? '订单编号' : 'Order ID' }}</span>
          <span class="col-email">{{ locale === 'zh' ? '买家邮箱' : 'Buyer Email' }}</span>
          <span class="col-product">{{ locale === 'zh' ? '购买内容 / 类型' : 'Item / Type' }}</span>
          <span class="col-amount">{{ t('payment.amount') }}</span>
          <span class="col-channel">{{ locale === 'zh' ? '支付方式 / 流水' : 'Channel / Ref' }}</span>
          <span class="col-status">{{ t('admin.order_status') }}</span>
          <span class="col-time">{{ locale === 'zh' ? '创建时间' : 'Created' }}</span>
          <span class="col-actions text-right">{{ locale === 'zh' ? '操作' : 'Action' }}</span>
        </div>

        <div v-if="loading" class="table-loading">
          <span class="spinner"></span>
          <span class="ml-2 text-secondary">{{ locale === 'zh' ? '正在加载订单列表...' : 'Loading orders...' }}</span>
        </div>
        <div v-else-if="orders.length === 0" class="table-empty">
          <div class="empty-icon-wrap">📂</div>
          <p>{{ t('common.no_data') }}</p>
          <small class="text-secondary">{{ locale === 'zh' ? '未找到符合当前筛选条件的订单' : 'No orders matched your filter' }}</small>
        </div>

        <div v-else class="table-body">
          <div v-for="order in orders" :key="order.id" class="table-row" :class="{ 'row-pending-review': order.status === 'pending_review' }">
            <!-- Col 1: Order ID -->
            <div class="col-order-id">
              <div class="id-badge-wrap" @click="copyText(order.id)">
                <span class="mono-id" :title="order.id">{{ order.id?.slice(0, 10) }}...</span>
                <span class="copy-hint" title="复制完整订单号">📋</span>
              </div>
            </div>

            <!-- Col 2: Buyer Email -->
            <div class="col-email">
              <div class="email-cell" :title="order.user_email || '游客'">
                <span class="email-icon">✉️</span>
                <span class="email-text">{{ order.user_email || '-' }}</span>
              </div>
            </div>

            <!-- Col 3: Product / Membership -->
            <div class="col-product">
              <div class="product-cell-wrap">
                <span class="type-pill" :class="order.order_type === 'membership' ? 'pill-membership' : 'pill-single'">
                  {{ order.order_type === 'membership' ? (locale === 'zh' ? '👑 会员' : '👑 Member') : (locale === 'zh' ? '🖼️ 单购' : '🖼️ Single') }}
                </span>
                <span class="product-name font-bold" :title="order.product_title || order.membership_type">
                  {{ order.product_title || getMembershipPlanTitle(order.membership_type) }}
                </span>
              </div>
            </div>

            <!-- Col 4: Amount -->
            <div class="col-amount">
              <span class="amount-val">{{ getCurrencySymbol(order.currency) }}{{ Number(order.amount || 0).toFixed(2) }}</span>
              <span class="amount-cur" :class="(order.currency || '').toUpperCase() === 'CNY' ? 'badge-cny' : 'badge-usd'">
                {{ order.currency || 'USD' }}
              </span>
            </div>

            <!-- Col 5: Payment Channel & Ref -->
            <div class="col-channel">
              <div class="channel-cell">
                <span class="channel-chip" :class="getChannelMeta(order.paypal_order_id).chipClass">
                  {{ getChannelMeta(order.paypal_order_id).icon }} {{ getChannelMeta(order.paypal_order_id).name }}
                </span>
                <span class="ref-id font-mono text-xs" :title="order.paypal_order_id">
                  {{ order.paypal_order_id || '-' }}
                </span>
              </div>
            </div>

            <!-- Col 6: Order Status -->
            <div class="col-status">
              <span class="status-pill" :class="getStatusBadgeClass(order.status)">
                <span class="status-dot"></span>
                {{ getStatusLabel(order.status) }}
              </span>
            </div>

            <!-- Col 7: Created Time -->
            <div class="col-time">
              <span class="time-text">{{ formatDate(order.created_at) }}</span>
            </div>

            <!-- Col 8: Action Buttons -->
            <div class="col-actions">
              <!-- When order is pending review -->
              <template v-if="order.status === 'pending_review'">
                <button class="btn btn-xs btn-approve" title="核验到账并为用户下发网盘直链" @click="updateStatus(order.id, 'paid')">
                  ✓ {{ locale === 'zh' ? '通过' : 'Approve' }}
                </button>
                <button class="btn btn-xs btn-reject" title="驳回未付款记录" @click="updateStatus(order.id, 'cancelled')">
                  ✕ {{ locale === 'zh' ? '驳回' : 'Reject' }}
                </button>
              </template>

              <!-- When order is pending pay -->
              <template v-else-if="order.status === 'pending'">
                <button class="btn btn-xs btn-mark-paid" title="管理员手动补单" @click="updateStatus(order.id, 'paid')">
                  {{ locale === 'zh' ? '标已付' : 'Mark' }}
                </button>
                <button class="btn btn-xs btn-cancel-order" @click="updateStatus(order.id, 'cancelled')">
                  {{ locale === 'zh' ? '取消' : 'Cancel' }}
                </button>
              </template>

              <!-- When order is already paid -->
              <div v-else-if="order.status === 'paid'" class="fulfilled-tag">
                <span class="fulfilled-check">✓</span>
                <span>{{ locale === 'zh' ? '已下发' : 'Fulfilled' }}</span>
              </div>

              <!-- When cancelled -->
              <span v-else class="text-xs text-muted">{{ locale === 'zh' ? '已关闭' : 'Closed' }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="total > pageSize" class="table-footer">
        <div class="pagination-info">
          {{ locale === 'zh' ? `共 ${total} 笔订单 · 每页 ${pageSize} 笔` : `Total ${total} orders · ${pageSize}/page` }}
        </div>
        <div class="pagination">
          <button class="btn btn-sm btn-ghost" :disabled="page <= 1" @click="page--; loadOrders()">
            ← {{ locale === 'zh' ? '上一页' : 'Prev' }}
          </button>
          <span class="page-num">{{ page }}</span>
          <button class="btn btn-sm btn-ghost" :disabled="page * pageSize >= total" @click="page++; loadOrders()">
            {{ locale === 'zh' ? '下一页' : 'Next' }} →
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { orderAPI } from '@/api'
import { useToast } from '@/composables/useToast'

const { t, locale } = useI18n()
const toast = useToast()

const orders = ref([])
const loading = ref(false)
const search = ref('')
const filterStatus = ref('')
const filterType = ref('')
const page = ref(1)
const pageSize = 15
const total = ref(0)

// Quick stats
const pendingReviewTotal = computed(() => {
  return orders.value.filter(o => o.status === 'pending_review').length
})

const paidTotal = computed(() => {
  return orders.value.filter(o => o.status === 'paid').length
})

const pageRevenueText = computed(() => {
  const paidOrders = orders.value.filter(o => o.status === 'paid')
  const cnySum = paidOrders
    .filter(o => (o.currency || '').toUpperCase() === 'CNY')
    .reduce((acc, cur) => acc + Number(cur.amount || 0), 0)
  const usdSum = paidOrders
    .filter(o => (o.currency || '').toUpperCase() !== 'CNY')
    .reduce((acc, cur) => acc + Number(cur.amount || 0), 0)

  const parts = []
  if (cnySum > 0) parts.push(`¥${cnySum.toFixed(2)} CNY`)
  if (usdSum > 0 || parts.length === 0) parts.push(`$${usdSum.toFixed(2)} USD`)
  return parts.join(' + ')
})

function getCurrencySymbol(curr) {
  const c = (curr || '').toUpperCase()
  if (c === 'CNY' || c === 'RMB') return '¥'
  if (c === 'USD') return '$'
  if (c === 'EUR') return '€'
  if (c === 'GBP') return '£'
  return c ? c + ' ' : '$'
}

let searchTimer = null
function handleSearch() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    page.value = 1
    loadOrders()
  }, 300)
}

async function loadOrders() {
  loading.value = true
  try {
    const res = await orderAPI.adminList({
      page: page.value,
      page_size: pageSize,
      status: filterStatus.value || undefined,
      order_type: filterType.value || undefined,
      search: search.value || undefined
    })
    orders.value = res.orders || res.items || res || []
    total.value = res.total || orders.value.length
  } catch (err) {
    console.error('Failed to load orders:', err)
  } finally {
    loading.value = false
  }
}

async function updateStatus(id, newStatus) {
  try {
    await orderAPI.updateStatus(id, newStatus)
    toast.success(locale.value === 'zh' ? '订单状态已成功更新' : 'Order status updated')
    loadOrders()
  } catch (err) {
    toast.error(err.message || t('common.error'))
  }
}

function getStatusBadgeClass(status) {
  if (status === 'paid') return 'status-paid'
  if (status === 'pending_review') return 'status-review'
  if (status === 'pending') return 'status-pending'
  return 'status-cancelled'
}

function getStatusLabel(status) {
  if (status === 'paid') return locale.value === 'zh' ? '已支付' : 'Paid'
  if (status === 'pending_review') return locale.value === 'zh' ? '待审核' : 'Review'
  if (status === 'pending') return locale.value === 'zh' ? '待付款' : 'Unpaid'
  if (status === 'cancelled') return locale.value === 'zh' ? '已取消' : 'Cancelled'
  return status
}

function getChannelMeta(refId) {
  if (!refId) return { name: locale.value === 'zh' ? '通用渠道' : 'General', icon: '💳', chipClass: 'chip-default' }
  const upper = refId.toUpperCase()
  if (upper.startsWith('WECHAT-')) return { name: locale.value === 'zh' ? '微信支付' : 'WeChat Pay', icon: '🟢', chipClass: 'chip-wechat' }
  if (upper.startsWith('ALIPAY-')) return { name: locale.value === 'zh' ? '支付宝' : 'Alipay', icon: '🔵', chipClass: 'chip-alipay' }
  if (upper.startsWith('PAYPAL-')) return { name: 'PayPal 国际', icon: '💳', chipClass: 'chip-paypal' }
  if (upper.startsWith('TEST-')) return { name: locale.value === 'zh' ? '沙盒免密' : 'Sandbox Test', icon: '⚡', chipClass: 'chip-test' }
  return { name: refId.slice(0, 8), icon: '💳', chipClass: 'chip-default' }
}

function getMembershipPlanTitle(type) {
  switch (type) {
    case 'monthly': return locale.value === 'zh' ? '月度 VIP 会员' : 'Monthly VIP'
    case 'quarterly': return locale.value === 'zh' ? '季度 VIP 会员' : 'Quarterly VIP'
    case 'yearly': return locale.value === 'zh' ? '年度 VIP 会员' : 'Yearly VIP'
    case 'lifetime': return locale.value === 'zh' ? '终生尊享 VIP' : 'Lifetime VIP'
    default: return locale.value === 'zh' ? 'VIP 会员方案' : 'VIP Plan'
  }
}

function copyText(txt) {
  if (!txt) return
  if (navigator?.clipboard?.writeText) {
    navigator.clipboard.writeText(txt).then(() => {
      toast.success(locale.value === 'zh' ? '完整订单号已复制' : 'Order ID copied')
    })
  } else {
    toast.info(txt)
  }
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString(locale.value === 'zh' ? 'zh-CN' : 'en-US', {
    month: 'numeric',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  })
}

onMounted(() => {
  loadOrders()
})
</script>

<style scoped>
.admin-orders {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

/* ========================================================
   Summary Stat Cards
   ======================================================== */
.orders-summary-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 14px;
}

@media (max-width: 960px) {
  .orders-summary-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

.summary-stat-card {
  background: var(--c-surface, #ffffff);
  border: 1px solid var(--c-border-light, #e2e8f0);
  border-radius: 14px;
  padding: 16px 18px;
  display: flex;
  align-items: center;
  gap: 14px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  transition: all 0.2s ease;
}

.summary-stat-card.clickable {
  cursor: pointer;
}

.summary-stat-card.clickable:hover {
  transform: translateY(-1px);
  border-color: rgba(0, 113, 227, 0.4);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.06);
}

.summary-stat-card.active {
  border-color: #0071e3;
  background: rgba(0, 113, 227, 0.03);
}

.stat-icon-wrap {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  flex-shrink: 0;
}

.bg-blue-subtle { background: rgba(59, 130, 246, 0.12); color: #2563eb; }
.bg-amber-subtle { background: rgba(245, 158, 11, 0.15); color: #d97706; }
.bg-emerald-subtle { background: rgba(16, 185, 129, 0.15); color: #059669; }
.bg-purple-subtle { background: rgba(139, 92, 246, 0.12); color: #7c3aed; }

.stat-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
  flex: 1;
}

.stat-label-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.stat-label {
  font-size: 0.82rem;
  color: var(--c-text-secondary, #64748b);
  font-weight: 600;
}

.stat-badge-pulse {
  background: #f59e0b;
  color: #ffffff;
  font-size: 10px;
  font-weight: 800;
  padding: 1px 6px;
  border-radius: 999px;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.75; transform: scale(1.08); }
}

.stat-value {
  font-size: 1.35rem;
  font-weight: 800;
  color: var(--c-text-primary, #0f172a);
}

.text-amber { color: #d97706 !important; }
.text-emerald { color: #059669 !important; }

/* ========================================================
   Toolbar
   ======================================================== */
.toolbar {
  padding: 14px 18px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
  background: var(--c-surface, #ffffff);
  border: 1px solid var(--c-border-light, #e2e8f0);
  border-radius: 14px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
  flex: 1;
}

.search-input-wrap {
  position: relative;
  min-width: 280px;
  flex: 1;
  max-width: 360px;
}

.search-icon {
  position: absolute;
  left: 10px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 14px;
  opacity: 0.6;
  pointer-events: none;
}

.search-input {
  padding-left: 32px;
  width: 100%;
  height: 38px;
  border-radius: 9px;
}

.btn-refresh {
  display: flex;
  align-items: center;
  gap: 6px;
  height: 38px;
}

.spinning {
  display: inline-block;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  100% { transform: rotate(360deg); }
}

/* ========================================================
   Orders Table Layout (Spacious, Non-overlapping)
   ======================================================== */
.orders-table-wrapper {
  background: var(--c-surface, #ffffff);
  border: 1px solid var(--c-border-light, #e2e8f0);
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
}

.table-scroll-container {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.table-header, .table-row {
  display: grid;
  /* 8 distinct, spacious columns totaling 1180px min-width */
  grid-template-columns: 140px 180px 1.5fr 110px 160px 120px 120px 150px;
  align-items: center;
  padding: 13px 20px;
  gap: 16px;
  min-width: 1100px;
  box-sizing: border-box;
}

.table-header {
  background: var(--c-bg, #f8fafc);
  font-size: 0.8rem;
  font-weight: 700;
  color: var(--c-text-secondary, #64748b);
  border-bottom: 1px solid var(--c-border-light, #e2e8f0);
  text-transform: uppercase;
  letter-spacing: 0.03em;
}

.table-row {
  border-bottom: 1px solid var(--c-border-light, #f1f5f9);
  font-size: 0.88rem;
  transition: background 0.15s ease;
}

.table-row:hover {
  background: rgba(0, 113, 227, 0.02);
}

.table-row.row-pending-review {
  background: rgba(245, 158, 11, 0.04);
}

.table-row.row-pending-review:hover {
  background: rgba(245, 158, 11, 0.08);
}

/* Column 1: Order ID */
.col-order-id {
  display: flex;
  align-items: center;
  overflow: hidden;
}

.id-badge-wrap {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px 8px;
  background: rgba(0, 0, 0, 0.04);
  border-radius: 6px;
  cursor: pointer;
  max-width: 100%;
}

.id-badge-wrap:hover {
  background: rgba(0, 113, 227, 0.1);
  color: #0071e3;
}

.mono-id {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.82rem;
  font-weight: 600;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.copy-hint {
  font-size: 11px;
  opacity: 0.7;
}

/* Column 2: Buyer Email */
.col-email {
  overflow: hidden;
}

.email-cell {
  display: flex;
  align-items: center;
  gap: 6px;
  overflow: hidden;
}

.email-icon {
  font-size: 13px;
  flex-shrink: 0;
  opacity: 0.7;
}

.email-text {
  font-size: 0.84rem;
  color: var(--c-text-primary, #0f172a);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* Column 3: Product / Membership */
.col-product {
  overflow: hidden;
}

.product-cell-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
  overflow: hidden;
}

.type-pill {
  font-size: 0.72rem;
  font-weight: 700;
  padding: 2px 7px;
  border-radius: 6px;
  flex-shrink: 0;
  white-space: nowrap;
}

.pill-membership {
  background: rgba(245, 158, 11, 0.15);
  color: #d97706;
  border: 1px solid rgba(245, 158, 11, 0.3);
}

.pill-single {
  background: rgba(59, 130, 246, 0.1);
  color: #2563eb;
  border: 1px solid rgba(59, 130, 246, 0.25);
}

.product-name {
  font-size: 0.88rem;
  color: var(--c-text-primary, #0f172a);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* Column 4: Amount */
.col-amount {
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.amount-val {
  font-size: 0.95rem;
  font-weight: 800;
  color: #0f172a;
}

.amount-cur {
  font-size: 0.72rem;
  color: var(--c-text-secondary, #64748b);
  font-weight: 600;
}

.amount-cur.badge-cny {
  background: rgba(245, 158, 11, 0.14);
  color: #d97706;
  border: 1px solid rgba(245, 158, 11, 0.28);
  padding: 1px 5px;
  border-radius: 4px;
  font-weight: 700;
  font-size: 10px;
}

.amount-cur.badge-usd {
  background: rgba(16, 185, 129, 0.12);
  color: #059669;
  border: 1px solid rgba(16, 185, 129, 0.28);
  padding: 1px 5px;
  border-radius: 4px;
  font-weight: 700;
  font-size: 10px;
}

/* Column 5: Payment Channel & Ref */
.col-channel {
  overflow: hidden;
}

.channel-cell {
  display: flex;
  flex-direction: column;
  gap: 2px;
  overflow: hidden;
}

.channel-chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 0.76rem;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 6px;
  width: fit-content;
  white-space: nowrap;
}

.chip-wechat { background: rgba(16, 185, 129, 0.12); color: #059669; }
.chip-alipay { background: rgba(59, 130, 246, 0.12); color: #2563eb; }
.chip-paypal { background: rgba(0, 113, 227, 0.12); color: #0071e3; }
.chip-test { background: rgba(245, 158, 11, 0.12); color: #d97706; }
.chip-default { background: rgba(0, 0, 0, 0.05); color: #64748b; }

.ref-id {
  color: var(--c-text-secondary, #64748b);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 11px;
}

/* Column 6: Status */
.col-status {
  display: flex;
  align-items: center;
}

.status-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 0.78rem;
  font-weight: 700;
  padding: 3px 10px;
  border-radius: 999px;
  white-space: nowrap;
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
}

.status-paid {
  background: rgba(16, 185, 129, 0.12);
  color: #059669;
  border: 1px solid rgba(16, 185, 129, 0.25);
}

.status-review {
  background: rgba(245, 158, 11, 0.15);
  color: #d97706;
  border: 1px solid rgba(245, 158, 11, 0.35);
  box-shadow: 0 0 8px rgba(245, 158, 11, 0.2);
}

.status-pending {
  background: rgba(100, 116, 139, 0.1);
  color: #64748b;
  border: 1px solid rgba(100, 116, 139, 0.2);
}

.status-cancelled {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.2);
}

/* Column 7: Time */
.col-time {
  font-size: 0.82rem;
  color: var(--c-text-secondary, #64748b);
  white-space: nowrap;
}

/* Column 8: Actions */
.col-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 6px;
}

.btn-xs {
  font-size: 0.74rem;
  font-weight: 700;
  padding: 3px 8px;
  border-radius: 6px;
  border: none;
  cursor: pointer;
  transition: all 0.15s ease;
  white-space: nowrap;
}

.btn-approve {
  background: #0071e3;
  color: #ffffff;
}
.btn-approve:hover {
  background: #005bb5;
  box-shadow: 0 2px 6px rgba(0, 113, 227, 0.3);
}

.btn-reject {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.2);
}
.btn-reject:hover {
  background: #ef4444;
  color: #ffffff;
}

.btn-mark-paid {
  background: rgba(0, 113, 227, 0.1);
  color: #0071e3;
}
.btn-mark-paid:hover {
  background: #0071e3;
  color: #ffffff;
}

.btn-cancel-order {
  background: none;
  color: #64748b;
}
.btn-cancel-order:hover {
  color: #ef4444;
}

.fulfilled-tag {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  font-size: 0.75rem;
  font-weight: 700;
  color: #059669;
}

.fulfilled-check {
  font-size: 11px;
}

/* Empty / Loading States */
.table-loading, .table-empty {
  padding: 56px 20px;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.empty-icon-wrap {
  font-size: 40px;
  opacity: 0.5;
}

/* Pagination */
.table-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 20px;
  background: var(--c-bg, #f8fafc);
  border-top: 1px solid var(--c-border-light, #e2e8f0);
}

.pagination-info {
  font-size: 0.84rem;
  color: var(--c-text-secondary, #64748b);
}

.pagination {
  display: flex;
  align-items: center;
  gap: 6px;
}

.page-num {
  font-weight: 700;
  font-size: 0.85rem;
  min-width: 28px;
  height: 28px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  background: #0071e3;
  color: #ffffff;
}
</style>
