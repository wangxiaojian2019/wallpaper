<template>
  <div class="admin-visitors">
    <!-- Header -->
    <div class="visitors-header">
      <div>
        <h2 class="page-title">👥 {{ locale === 'zh' ? '访客管理与访问日志' : 'Visitor Management & Logs' }}</h2>
        <p class="page-sub">
          {{ locale === 'zh' ? '实时追踪全站访客 IP、地理归属地、访问设备系统、停留时长与访问轨迹' : 'Real-time visitor logs, IP geo-location, devices and dwell duration' }}
        </p>
      </div>
      <button class="btn btn-outline" :disabled="loading" @click="loadVisitors">
        <span v-if="loading" class="spinner"></span>
        <span v-else>🔄 {{ locale === 'zh' ? '刷新访客记录' : 'Refresh Logs' }}</span>
      </button>
    </div>

    <!-- 4 Stats Cards -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon icon-blue">👀</div>
        <div class="stat-meta">
          <span class="stat-title">{{ locale === 'zh' ? '今日访问量 (PV)' : 'Today PV' }}</span>
          <span class="stat-val font-bold text-gradient">{{ stats.today_pv }}</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon icon-emerald">👤</div>
        <div class="stat-meta">
          <span class="stat-title">{{ locale === 'zh' ? '今日独立访客 (UV)' : 'Today UV' }}</span>
          <span class="stat-val font-bold">{{ stats.today_uv }}</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon icon-purple">⏱️</div>
        <div class="stat-meta">
          <span class="stat-title">{{ locale === 'zh' ? '全站平均停留时长' : 'Avg Dwell Time' }}</span>
          <span class="stat-val font-bold">{{ formatSeconds(stats.avg_duration) }}</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon icon-gold">🌐</div>
        <div class="stat-meta">
          <span class="stat-title">{{ locale === 'zh' ? '总记录访客数' : 'Total Logs' }}</span>
          <span class="stat-val font-bold">{{ total }}</span>
        </div>
      </div>
    </div>

    <!-- Filter Toolbar -->
    <div class="card filter-bar">
      <div class="filter-inputs">
        <div class="filter-item">
          <label class="filter-label">{{ locale === 'zh' ? '搜索 IP' : 'Search IP' }}</label>
          <input
            v-model="filterIP"
            type="text"
            class="form-input form-input-sm"
            :placeholder="locale === 'zh' ? '输入 IP 地址搜索' : 'Enter IP address'"
            @keyup.enter="handleSearch"
          />
        </div>

        <div class="filter-item">
          <label class="filter-label">{{ locale === 'zh' ? '访问设备' : 'Device' }}</label>
          <select v-model="filterDevice" class="form-input form-input-sm form-select" @change="handleSearch">
            <option value="all">{{ locale === 'zh' ? '全部设备类型' : 'All Devices' }}</option>
            <option value="Mobile">📱 {{ locale === 'zh' ? '手机端 (Mobile)' : 'Mobile' }}</option>
            <option value="Desktop">💻 {{ locale === 'zh' ? '电脑端 (Desktop)' : 'Desktop' }}</option>
            <option value="Tablet">📟 {{ locale === 'zh' ? '平板电脑 (Tablet)' : 'Tablet' }}</option>
          </select>
        </div>

        <div class="filter-item">
          <label class="filter-label">{{ locale === 'zh' ? '归属地' : 'Country' }}</label>
          <input
            v-model="filterCountry"
            type="text"
            class="form-input form-input-sm"
            :placeholder="locale === 'zh' ? '输入省市或国家' : 'Search Country/City'"
            @keyup.enter="handleSearch"
          />
        </div>
      </div>

      <div class="filter-actions">
        <button class="btn btn-primary btn-sm" @click="handleSearch">
          🔍 {{ locale === 'zh' ? '查询' : 'Search' }}
        </button>
        <button class="btn btn-ghost btn-sm" @click="resetFilter">
          ✕ {{ locale === 'zh' ? '重置' : 'Reset' }}
        </button>
      </div>
    </div>

    <!-- Visitor Logs Table -->
    <div class="card table-card">
      <div class="table-wrap">
        <table class="data-table">
          <thead>
            <tr>
              <th>{{ locale === 'zh' ? '访客 IP' : 'IP Address' }}</th>
              <th>{{ locale === 'zh' ? '用户身份' : 'Identity' }}</th>
              <th>{{ locale === 'zh' ? '地域归属' : 'Location' }}</th>
              <th>{{ locale === 'zh' ? '访问设备 / 系统' : 'Device / OS' }}</th>
              <th>{{ locale === 'zh' ? '浏览器' : 'Browser' }}</th>
              <th>{{ locale === 'zh' ? '访问页面' : 'Page Path' }}</th>
              <th>{{ locale === 'zh' ? '语言模式' : 'Language' }}</th>
              <th>{{ locale === 'zh' ? '停留时长' : 'Dwell Time' }}</th>
              <th>{{ locale === 'zh' ? '访问时间' : 'Visit Time' }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td colspan="9" class="text-center py-8">
                <span class="spinner"></span>
                <span class="ml-2 text-secondary">{{ locale === 'zh' ? '正在加载访客日志...' : 'Loading visitor logs...' }}</span>
              </td>
            </tr>
            <tr v-else-if="visitors.length === 0">
              <td colspan="9" class="text-center py-8 text-secondary">
                {{ locale === 'zh' ? '暂无匹配的访客访问记录' : 'No visitor records found' }}
              </td>
            </tr>
            <tr v-else v-for="log in visitors" :key="log.id" class="table-row-hover">
              <td>
                <div class="ip-cell">
                  <span class="ip-dot"></span>
                  <span class="ip-text font-mono font-bold">{{ log.ip }}</span>
                </div>
              </td>
              <!-- User Identity / Member Badge -->
              <td>
                <div v-if="log.user_email" class="user-identity-cell">
                  <div class="user-id-row">
                    <span class="user-avatar-dot">👤</span>
                    <span class="user-email-text" :title="log.user_email">{{ log.user_email }}</span>
                  </div>
                  <span v-if="log.is_member || (log.membership_type && log.membership_type !== 'none')" class="badge badge-gold-mini">
                    👑 VIP ({{ formatTierName(log.membership_type) }})
                  </span>
                  <span v-else class="badge badge-user-mini">
                    已注册用户
                  </span>
                </div>
                <div v-else class="guest-identity-cell">
                  <span class="badge badge-guest-mini">🌐 匿名访客</span>
                </div>
              </td>
              <td>
                <span class="badge" :class="isChinaIp(log.country) ? 'badge-primary' : 'badge-neutral'">
                  {{ log.country || (locale === 'zh' ? '未知地域' : 'Unknown') }}
                  <span v-if="log.city">· {{ log.city }}</span>
                </span>
              </td>
              <td>
                <div class="device-cell">
                  <span class="device-icon">{{ getDeviceIcon(log.device) }}</span>
                  <span class="device-text font-semibold">{{ log.device || 'Desktop' }}</span>
                  <span class="os-text text-secondary text-xs">({{ log.os || 'Windows' }})</span>
                </div>
              </td>
              <td>
                <span class="browser-pill">{{ log.browser || 'Browser' }}</span>
              </td>
              <td>
                <span class="path-code font-mono" :title="log.path">{{ log.path || '/' }}</span>
              </td>
              <td>
                <span class="lang-tag" :class="log.language === 'zh' ? 'lang-zh' : 'lang-en'">
                  {{ log.language === 'zh' ? '🇨🇳 中文' : '🌐 English' }}
                </span>
              </td>
              <td>
                <span class="duration-badge font-bold" :class="getDurationClass(log.duration)">
                  ⏱️ {{ formatSeconds(log.duration) }}
                </span>
              </td>
              <td class="text-secondary text-sm">
                {{ formatDate(log.created_at) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div class="table-pagination">
        <div class="pagination-info">
          {{ locale === 'zh' ? `共 ${total} 条记录 · 当前第 ${page} / ${Math.ceil(total / pageSize) || 1} 页` : `Total ${total} logs · Page ${page} of ${Math.ceil(total / pageSize) || 1}` }}
        </div>
        <div class="pagination-btns">
          <button class="btn btn-sm btn-outline" :disabled="page <= 1" @click="page--; loadVisitors()">
            ← {{ locale === 'zh' ? '上一页' : 'Prev' }}
          </button>
          <span class="page-badge">{{ page }}</span>
          <button class="btn btn-sm btn-outline" :disabled="page * pageSize >= total" @click="page++; loadVisitors()">
            {{ locale === 'zh' ? '下一页' : 'Next' }} →
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { visitorAPI } from '@/api'
import { useToast } from '@/composables/useToast'

const { locale } = useI18n()
const { error } = useToast()

const visitors = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const loading = ref(false)

const filterIP = ref('')
const filterDevice = ref('all')
const filterCountry = ref('')

const stats = ref({
  today_pv: 0,
  today_uv: 0,
  avg_duration: 0
})

function formatSeconds(sec) {
  const s = Math.round(Number(sec) || 0)
  if (s < 60) return `${s} 秒`
  const min = Math.floor(s / 60)
  const rem = s % 60
  if (min < 60) return `${min} 分 ${rem} 秒`
  const hours = Math.floor(min / 60)
  return `${hours} 小时 ${min % 60} 分`
}

function formatDate(isoStr) {
  if (!isoStr) return '-'
  const d = new Date(isoStr)
  if (isNaN(d.getTime())) return isoStr
  return d.toLocaleString(locale.value === 'zh' ? 'zh-CN' : 'en-US', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  })
}

function isChinaIp(country) {
  if (!country) return false
  const c = country.toLowerCase()
  return c.includes('中国') || c.includes('china') || c.includes('cn') || c.includes('华语') || c.includes('内网')
}

function getDeviceIcon(device) {
  if (!device) return '💻'
  const d = device.toLowerCase()
  if (d.includes('mobile') || d.includes('phone')) return '📱'
  if (d.includes('tablet') || d.includes('ipad')) return '📟'
  return '💻'
}

function getDurationClass(sec) {
  const s = Number(sec) || 0
  if (s >= 120) return 'duration-long'
  if (s >= 30) return 'duration-med'
  return 'duration-short'
}

function formatTierName(tier) {
  switch (tier) {
    case 'monthly': return locale.value === 'zh' ? '月度' : 'Month'
    case 'quarterly': return locale.value === 'zh' ? '季度' : 'Quarter'
    case 'yearly': return locale.value === 'zh' ? '年度' : 'Year'
    case 'lifetime': return locale.value === 'zh' ? '终身' : 'Lifetime'
    default: return locale.value === 'zh' ? '会员' : 'VIP'
  }
}

async function loadVisitors() {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value
    }
    if (filterIP.value.trim()) params.ip = filterIP.value.trim()
    if (filterDevice.value && filterDevice.value !== 'all') params.device = filterDevice.value
    if (filterCountry.value.trim()) params.country = filterCountry.value.trim()

    const res = await visitorAPI.list(params)
    if (res) {
      visitors.value = res.items || []
      total.value = res.total || 0
      stats.value.today_pv = res.today_pv || 0
      stats.value.today_uv = res.today_uv || 0
      stats.value.avg_duration = res.avg_duration || 0
    }
  } catch (err) {
    error(err.message || '加载访客记录失败')
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  page.value = 1
  loadVisitors()
}

function resetFilter() {
  filterIP.value = ''
  filterDevice.value = 'all'
  filterCountry.value = ''
  page.value = 1
  loadVisitors()
}

onMounted(() => {
  loadVisitors()
})
</script>

<style scoped>
.admin-visitors {
  max-width: 1280px;
  display: flex;
  flex-direction: column;
  gap: 24px;
  padding-bottom: 60px;
}

.visitors-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-title {
  font-size: 1.5rem;
  font-weight: 800;
  margin: 0 0 6px;
}

.page-sub {
  font-size: 0.9rem;
  color: var(--c-text-secondary);
  margin: 0;
}

/* Stats grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 18px 20px;
  background: var(--c-surface);
  border-radius: var(--radius-lg);
  border: 1px solid var(--c-border);
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.icon-blue { background: rgba(59, 130, 246, 0.15); }
.icon-emerald { background: rgba(16, 185, 129, 0.15); }
.icon-purple { background: rgba(168, 85, 247, 0.15); }
.icon-gold { background: rgba(245, 158, 11, 0.15); }

.stat-meta {
  display: flex;
  flex-direction: column;
}

.stat-title {
  font-size: 0.85rem;
  color: var(--c-text-secondary);
}

.stat-val {
  font-size: 1.4rem;
  color: var(--c-text);
  margin-top: 2px;
}

/* Filter bar */
.filter-bar {
  padding: 16px 20px;
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  gap: 16px;
  flex-wrap: wrap;
}

.filter-inputs {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.filter-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 160px;
}

.filter-label {
  font-size: 0.78rem;
  font-weight: 700;
  color: var(--c-text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.filter-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

/* Table card */
.table-card {
  overflow: hidden;
  padding: 0;
}

.table-wrap {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.data-table th {
  padding: 14px 18px;
  font-size: 0.8rem;
  font-weight: 700;
  color: var(--c-text-secondary);
  background: var(--c-surface-elevated, rgba(255, 255, 255, 0.03));
  border-bottom: 1px solid var(--c-border);
  white-space: nowrap;
}

.data-table td {
  padding: 14px 18px;
  border-bottom: 1px solid var(--c-border);
  vertical-align: middle;
}

.table-row-hover:hover {
  background: rgba(255, 255, 255, 0.02);
}

.ip-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.ip-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #10b981;
  box-shadow: 0 0 8px rgba(16, 185, 129, 0.6);
}

.ip-text {
  color: var(--c-text);
  font-size: 0.9rem;
}

.device-cell {
  display: flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
}

.browser-pill {
  display: inline-block;
  padding: 2px 8px;
  background: rgba(255, 255, 255, 0.06);
  border-radius: 999px;
  font-size: 0.8rem;
  color: var(--c-text-secondary);
}

.path-code {
  display: inline-block;
  max-width: 160px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 0.8rem;
  color: var(--c-primary);
  background: rgba(59, 130, 246, 0.1);
  padding: 2px 8px;
  border-radius: 4px;
}

.lang-tag {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 700;
}

.lang-zh {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
}

.lang-en {
  background: rgba(59, 130, 246, 0.15);
  color: #3b82f6;
}

.duration-badge {
  display: inline-block;
  padding: 3px 10px;
  border-radius: 999px;
  font-size: 0.8rem;
}

.duration-short {
  background: rgba(148, 163, 184, 0.15);
  color: #94a3b8;
}

.duration-med {
  background: rgba(59, 130, 246, 0.15);
  color: #38bdf8;
}

.duration-long {
  background: rgba(16, 185, 129, 0.18);
  color: #34d399;
}

/* Pagination */
.table-pagination {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  background: var(--c-surface-elevated, rgba(255, 255, 255, 0.02));
}

.pagination-info {
  font-size: 0.85rem;
  color: var(--c-text-secondary);
}

.pagination-btns {
  display: flex;
  align-items: center;
  gap: 8px;
}

.page-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: var(--c-primary);
  color: white;
  border-radius: var(--radius-sm);
  font-weight: 700;
  font-size: 0.85rem;
}

/* User Identity Cell Styles */
.user-identity-cell {
  display: flex;
  flex-direction: column;
  gap: 3px;
  max-width: 170px;
}

.user-id-row {
  display: flex;
  align-items: center;
  gap: 5px;
}

.user-avatar-dot {
  font-size: 0.85rem;
}

.user-email-text {
  font-size: 0.82rem;
  font-weight: 700;
  color: var(--c-text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.badge-gold-mini {
  background: rgba(245, 158, 11, 0.15);
  color: #d97706;
  border: 1px solid rgba(245, 158, 11, 0.35);
  font-size: 0.72rem;
  padding: 1px 6px;
  border-radius: 999px;
  width: fit-content;
  font-weight: 700;
}

.badge-user-mini {
  background: rgba(59, 130, 246, 0.12);
  color: #2563eb;
  border: 1px solid rgba(59, 130, 246, 0.25);
  font-size: 0.72rem;
  padding: 1px 6px;
  border-radius: 999px;
  width: fit-content;
  font-weight: 600;
}

.badge-guest-mini {
  background: rgba(148, 163, 184, 0.12);
  color: #64748b;
  font-size: 0.74rem;
  padding: 2px 8px;
  border-radius: 999px;
  font-weight: 500;
}
</style>
