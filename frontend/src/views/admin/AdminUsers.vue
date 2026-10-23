<template>
  <div class="admin-users">
    <!-- Toolbar -->
    <div class="toolbar">
      <div class="toolbar-left">
        <input v-model="search" type="text" class="form-input search-input" :placeholder="t('admin.search_user')" @input="handleSearch" />
        <select v-model="filterMember" class="form-input form-select" style="width:auto" @change="loadUsers">
          <option value="">{{ t('common.all') }} {{ locale === 'zh' ? '会员类型' : 'Plans' }}</option>
          <option value="none">{{ locale === 'zh' ? '普通用户' : 'Free' }}</option>
          <option value="monthly">{{ locale === 'zh' ? '月度会员' : 'Monthly' }}</option>
          <option value="quarterly">{{ locale === 'zh' ? '季度会员' : 'Quarterly' }}</option>
          <option value="yearly">{{ locale === 'zh' ? '年度会员' : 'Yearly' }}</option>
          <option value="lifetime">{{ locale === 'zh' ? '终身尊享' : 'Lifetime' }}</option>
        </select>
      </div>
      <div class="toolbar-right">
        <button class="btn btn-primary" @click="openCreateModal">
          ➕ {{ locale === 'zh' ? '添加会员/用户' : 'Add User' }}
        </button>
        <button class="btn btn-outline" @click="loadUsers">
          🔄 {{ locale === 'zh' ? '刷新' : 'Refresh' }}
        </button>
      </div>
    </div>

    <!-- Table -->
    <div class="users-table card">
      <div class="table-header">
        <span class="col-email">{{ t('auth.email') }}</span>
        <span class="col-role">{{ locale === 'zh' ? '角色' : 'Role' }}</span>
        <span class="col-member">{{ t('dashboard.membership') }}</span>
        <span class="col-expire">{{ locale === 'zh' ? '会员到期' : 'Expires' }}</span>
        <span class="col-login">{{ locale === 'zh' ? '最后登录' : 'Last Login' }}</span>
        <span class="col-reg">{{ locale === 'zh' ? '注册时间' : 'Registered' }}</span>
        <span class="col-actions">{{ locale === 'zh' ? '快捷操作' : 'Action' }}</span>
      </div>

      <div v-if="loading" class="table-loading"><span class="spinner"></span></div>
      <div v-else-if="users.length === 0" class="table-empty">{{ t('common.no_data') }}</div>

      <div v-else>
        <div v-for="user in users" :key="user.id" class="table-row">
          <!-- Email & Avatar -->
          <div class="col-email">
            <div class="user-info">
              <span class="user-avatar-mini" :class="{ 'vip-avatar': user.membership_type !== 'none' }">
                {{ user.email ? user.email[0].toUpperCase() : 'U' }}
              </span>
              <div class="user-email-col">
                <span class="email-text">{{ user.email }}</span>
                <span class="user-id-sub">ID: {{ user.id?.slice(-8) }}</span>
              </div>
            </div>
          </div>

          <!-- Role -->
          <div class="col-role">
            <span class="badge" :class="user.role === 'admin' ? 'badge-primary' : 'badge-neutral'">
              {{ user.role === 'admin' ? (locale === 'zh' ? '管理员' : 'Admin') : (locale === 'zh' ? '普通用户' : 'User') }}
            </span>
          </div>

          <!-- Membership -->
          <div class="col-member">
            <span class="badge" :class="getMemberBadge(user)">
              {{ formatMember(user.membership_type) }}
            </span>
          </div>

          <!-- Expire -->
          <div class="col-expire text-secondary">
            {{ formatExpire(user.membership_expire) }}
          </div>

          <!-- Last Login & IP -->
          <div class="col-login">
            <div v-if="user.last_login_at" class="login-meta">
              <span class="login-time">{{ formatDate(user.last_login_at) }}</span>
              <span class="login-ip font-mono text-xs">🌐 {{ user.last_login_ip || '未知IP' }}</span>
            </div>
            <span v-else class="text-tertiary text-xs">{{ locale === 'zh' ? '从未登录' : 'Never logged in' }}</span>
          </div>

          <!-- Registered -->
          <div class="col-reg text-secondary">
            {{ formatDate(user.created_at) }}
          </div>

          <!-- Actions -->
          <div class="col-actions">
            <!-- Reset Password Button -->
            <button class="btn btn-sm btn-ghost btn-icon-text" @click="openResetPwdModal(user)" :title="locale === 'zh' ? '重置登录密码' : 'Reset Password'">
              🔑 {{ locale === 'zh' ? '重置密码' : 'Password' }}
            </button>

            <!-- Force Logout Button -->
            <button
              class="btn btn-sm btn-ghost text-amber btn-icon-text"
              @click="handleForceLogout(user)"
              :title="locale === 'zh' ? '强制该用户在所有设备立即下线' : 'Force Logout'"
            >
              🚫 {{ locale === 'zh' ? '强制下线' : 'Kick' }}
            </button>

            <!-- Manage Role & Membership -->
            <button class="btn btn-sm btn-ghost btn-icon-text" @click="openEditModal(user)" :title="locale === 'zh' ? '调整会员与权限' : 'Manage'">
              ⚙️ {{ locale === 'zh' ? '权限' : 'Role' }}
            </button>

            <!-- Delete -->
            <button v-if="user.role !== 'admin'" class="btn btn-sm btn-ghost text-danger" @click="deleteUser(user)" title="删除用户">
              🗑️
            </button>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="total > pageSize" class="table-footer">
        <div class="pagination-info">
          {{ locale === 'zh' ? `共 ${total} 位注册用户` : `Total ${total} users` }}
        </div>
        <div class="pagination">
          <button class="btn btn-sm btn-ghost" :disabled="page <= 1" @click="page--; loadUsers()">
            ← {{ locale === 'zh' ? '上一页' : 'Prev' }}
          </button>
          <span class="page-num">{{ page }}</span>
          <button class="btn btn-sm btn-ghost" :disabled="page * pageSize >= total" @click="page++; loadUsers()">
            {{ locale === 'zh' ? '下一页' : 'Next' }} →
          </button>
        </div>
      </div>
    </div>

    <!-- Password Reset Modal -->
    <div v-if="resetPwdTarget" class="modal-overlay" @click.self="resetPwdTarget = null">
      <div class="modal-box" style="max-width: 440px;">
        <h3>🔑 {{ locale === 'zh' ? '管理员重置用户密码' : 'Reset User Password' }}</h3>
        <p class="user-sub">{{ resetPwdTarget.email }}</p>

        <div class="form-group">
          <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom:6px;">
            <label class="form-label font-bold" style="margin:0;">{{ locale === 'zh' ? '新登录密码' : 'New Password' }} *</label>
            <button type="button" class="btn btn-link btn-xs" @click="generateRandomPassword">
              🎲 {{ locale === 'zh' ? '生成强密码' : 'Random' }}
            </button>
          </div>
          <input
            v-model="newPassword"
            type="text"
            class="form-input font-mono font-bold"
            :placeholder="locale === 'zh' ? '输入至少 6 位的新密码' : 'Enter new password (min 6 chars)'"
          />
          <small class="form-tip">
            {{ locale === 'zh' ? '重置密码后，系统会自动更新该用户的 Token 版本，其在所有客户端的旧登录会话将被同步强制失效。' : 'Previous active login tokens will be immediately invalidated.' }}
          </small>
        </div>

        <div class="modal-actions">
          <button class="btn btn-ghost" @click="resetPwdTarget = null">{{ t('common.cancel') }}</button>
          <button class="btn btn-primary" :disabled="resettingPwd || !newPassword" @click="executeResetPassword">
            <span v-if="resettingPwd" class="spinner"></span>
            {{ locale === 'zh' ? '确认重置' : 'Confirm Reset' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Create User Modal -->
    <div v-if="showCreateModal" class="modal-overlay" @click.self="showCreateModal = false">
      <div class="modal-box">
        <h3>➕ {{ locale === 'zh' ? '添加新用户 / 手动开通会员' : 'Create New User / Grant VIP' }}</h3>
        <p class="user-sub">{{ locale === 'zh' ? '手动创建账号，支持直接开通指定天数的 VIP 会员特权' : 'Manually register user and grant membership days' }}</p>

        <div class="form-group">
          <label class="form-label">{{ locale === 'zh' ? '登录邮箱 (账号)' : 'Email' }} *</label>
          <input v-model="createForm.email" type="email" class="form-input" placeholder="user@example.com" />
        </div>

        <div class="form-group">
          <label class="form-label">{{ locale === 'zh' ? '登录密码' : 'Password' }} *</label>
          <input v-model="createForm.password" type="text" class="form-input" placeholder="至少 6 位密码" />
        </div>

        <div class="form-group">
          <label class="form-label">{{ locale === 'zh' ? '角色身份' : 'Account Role' }}</label>
          <select v-model="createForm.role" class="form-input form-select">
            <option value="user">User ({{ locale === 'zh' ? '普通用户' : 'Standard' }})</option>
            <option value="admin">Admin ({{ locale === 'zh' ? '管理员' : 'Administrator' }})</option>
          </select>
        </div>

        <div class="form-group">
          <label class="form-label">{{ locale === 'zh' ? '会员等级' : 'Membership Level' }}</label>
          <select v-model="createForm.membership_type" class="form-input form-select">
            <option value="none">{{ locale === 'zh' ? '普通用户 (无会员)' : 'None' }}</option>
            <option value="monthly">{{ locale === 'zh' ? '月度会员 (30天)' : 'Monthly (30 days)' }}</option>
            <option value="quarterly">{{ locale === 'zh' ? '季度会员 (90天)' : 'Quarterly (90 days)' }}</option>
            <option value="yearly">{{ locale === 'zh' ? '年度会员 (365天)' : 'Yearly (365 days)' }}</option>
            <option value="lifetime">{{ locale === 'zh' ? '终身尊享 (永久免单)' : 'Lifetime (Permanent)' }}</option>
          </select>
        </div>

        <div v-if="createForm.membership_type !== 'none'" class="form-group">
          <label class="form-label">{{ locale === 'zh' ? '自定义有效天数 (0 为默认)' : 'Custom Duration Days' }}</label>
          <input v-model.number="createForm.days" type="number" min="0" class="form-input" placeholder="0" />
        </div>

        <div class="modal-actions">
          <button class="btn btn-ghost" @click="showCreateModal = false">{{ t('common.cancel') }}</button>
          <button class="btn btn-primary" :disabled="creating" @click="handleCreateUser">
            {{ creating ? t('common.loading') : locale === 'zh' ? '立即创建' : 'Create' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Edit User Modal -->
    <div v-if="activeUser" class="modal-overlay" @click.self="activeUser = null">
      <div class="modal-box">
        <h3>{{ locale === 'zh' ? '调整用户权限与会员' : 'Manage User & Membership' }}</h3>
        <p class="user-sub">{{ activeUser.email }}</p>

        <div class="form-group">
          <label class="form-label">{{ locale === 'zh' ? '角色身份' : 'Account Role' }}</label>
          <select v-model="editForm.role" class="form-input form-select">
            <option value="user">User ({{ locale === 'zh' ? '普通用户' : 'Standard' }})</option>
            <option value="admin">Admin ({{ locale === 'zh' ? '管理员' : 'Administrator' }})</option>
          </select>
        </div>

        <div class="form-group">
          <label class="form-label">{{ locale === 'zh' ? '会员等级' : 'Membership Level' }}</label>
          <select v-model="editForm.membership_type" class="form-input form-select">
            <option value="none">{{ locale === 'zh' ? '无会员' : 'None' }}</option>
            <option value="monthly">{{ locale === 'zh' ? '月度会员 (30天)' : 'Monthly (30 days)' }}</option>
            <option value="quarterly">{{ locale === 'zh' ? '季度会员 (90天)' : 'Quarterly (90 days)' }}</option>
            <option value="yearly">{{ locale === 'zh' ? '年度会员 (365天)' : 'Yearly (365 days)' }}</option>
            <option value="lifetime">{{ locale === 'zh' ? '终生尊享 (永久免单)' : 'Lifetime (Permanent)' }}</option>
          </select>
        </div>

        <div class="form-group">
          <label class="form-label">{{ locale === 'zh' ? '续费/延长天数' : 'Extend Days' }}</label>
          <input v-model.number="editForm.extendDays" type="number" min="0" class="form-input" placeholder="0" />
          <span class="form-help">{{ locale === 'zh' ? '留空或 0 则按选择等级自动重设' : 'Leave 0 to set default duration' }}</span>
        </div>

        <div class="modal-actions">
          <button class="btn btn-ghost" @click="activeUser = null">{{ t('common.cancel') }}</button>
          <button class="btn btn-primary" :disabled="saving" @click="saveUser">
            {{ saving ? t('common.loading') : t('common.save') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { userAPI } from '@/api'
import { useToast } from '@/composables/useToast'
import { useDebounceFn } from '@vueuse/core'

const { t, locale } = useI18n()
const toast = useToast()

const users = ref([])
const loading = ref(false)
const search = ref('')
const filterMember = ref('')
const page = ref(1)
const pageSize = 15
const total = ref(0)

const activeUser = ref(null)
const saving = ref(false)
const editForm = ref({ role: 'user', membership_type: 'none', extendDays: 0 })

const showCreateModal = ref(false)
const creating = ref(false)
const createForm = ref({ email: '', password: '', role: 'user', membership_type: 'none', days: 0 })

// Password Reset Target & State
const resetPwdTarget = ref(null)
const newPassword = ref('')
const resettingPwd = ref(false)

function openResetPwdModal(user) {
  resetPwdTarget.value = user
  newPassword.value = ''
}

function generateRandomPassword() {
  const chars = 'abcdefghjkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789!@#$%'
  let pwd = ''
  for (let i = 0; i < 10; i++) {
    pwd += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  newPassword.value = pwd
}

async function executeResetPassword() {
  if (!newPassword.value || newPassword.value.length < 6) {
    toast.error(locale.value === 'zh' ? '新密码至少需要 6 个字符' : 'Password must be at least 6 characters')
    return
  }
  resettingPwd.value = true
  try {
    const res = await userAPI.resetPassword(resetPwdTarget.value.id, newPassword.value)
    toast.success(res.message || (locale.value === 'zh' ? '密码已成功重置' : 'Password reset successful'))
    resetPwdTarget.value = null
  } catch (err) {
    toast.error(err.message || (locale.value === 'zh' ? '重置密码失败' : 'Failed to reset password'))
  } finally {
    resettingPwd.value = false
  }
}

async function handleForceLogout(user) {
  if (!confirm(locale.value === 'zh' ? `确定强制下线用户【${user.email}】吗？该用户所有客户端将立即失效。` : `Force logout user ${user.email}?`)) {
    return
  }
  try {
    const res = await userAPI.forceLogout(user.id)
    toast.success(res.message || (locale.value === 'zh' ? '已强制该用户下线' : 'User forced to logout'))
  } catch (err) {
    toast.error(err.message || '强制下线操作失败')
  }
}

async function loadUsers() {
  loading.value = true
  try {
    const res = await userAPI.list({
      page: page.value,
      page_size: pageSize,
      membership_type: filterMember.value,
      search: search.value
    })
    users.value = res.users || []
    total.value = res.total || 0
  } catch (err) {
    toast.error(err.message || 'Failed to load users')
  } finally {
    loading.value = false
  }
}

const handleSearch = useDebounceFn(() => {
  page.value = 1
  loadUsers()
}, 350)

function openCreateModal() {
  createForm.value = { email: '', password: '', role: 'user', membership_type: 'none', days: 0 }
  showCreateModal.value = true
}

async function handleCreateUser() {
  if (!createForm.value.email || !createForm.value.password) {
    toast.error(locale.value === 'zh' ? '请填写完整的邮箱与初始密码' : 'Email and password required')
    return
  }
  creating.value = true
  try {
    await userAPI.create(createForm.value)
    toast.success(locale.value === 'zh' ? '用户创建成功！' : 'User created successfully')
    showCreateModal.value = false
    loadUsers()
  } catch (err) {
    toast.error(err.message || '创建失败')
  } finally {
    creating.value = false
  }
}

function openEditModal(user) {
  activeUser.value = user
  editForm.value = {
    role: user.role,
    membership_type: user.membership_type || 'none',
    extendDays: 0
  }
}

async function saveUser() {
  saving.value = true
  try {
    await userAPI.update(activeUser.value.id, editForm.value)
    toast.success(t('common.success'))
    activeUser.value = null
    loadUsers()
  } catch (err) {
    toast.error(err.message || t('common.error'))
  } finally {
    saving.value = false
  }
}

async function deleteUser(user) {
  if (!confirm(t('admin.confirm_delete') + `: ${user.email}?`)) return
  try {
    await userAPI.delete(user.id)
    toast.success(t('common.success'))
    loadUsers()
  } catch (err) {
    toast.error(err.message || t('common.error'))
  }
}

function formatMember(type) {
  switch (type) {
    case 'monthly': return locale.value === 'zh' ? '月度 VIP' : 'Monthly'
    case 'quarterly': return locale.value === 'zh' ? '季度 VIP' : 'Quarterly'
    case 'yearly': return locale.value === 'zh' ? '年度 VIP' : 'Yearly'
    case 'lifetime': return locale.value === 'zh' ? '终生尊享' : 'Lifetime'
    default: return locale.value === 'zh' ? '普通用户' : 'Standard'
  }
}

function getMemberBadge(user) {
  if (user.membership_type === 'lifetime') return 'badge-gold'
  if (user.membership_type === 'yearly') return 'badge-primary'
  if (user.membership_type === 'quarterly') return 'badge-info'
  if (user.membership_type === 'monthly') return 'badge-secondary'
  return 'badge-neutral'
}

function formatExpire(dateStr) {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return '-'
  return d.toISOString().split('T')[0]
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  if (isNaN(d.getTime())) return '-'
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const hh = String(d.getHours()).padStart(2, '0')
  const mm = String(d.getMinutes()).padStart(2, '0')
  return `${y}-${m}-${day} ${hh}:${mm}`
}

onMounted(loadUsers)
</script>

<style scoped>
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.toolbar-left, .toolbar-right {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.search-input { width: 220px; border-radius: var(--r-full); padding: 8px 16px; font-size: 0.88rem; }

.users-table {
  border: 1px solid var(--c-border-light);
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.table-header, .table-row {
  display: grid;
  grid-template-columns: minmax(180px, 1.2fr) 90px 100px 110px 160px 140px auto;
  gap: 12px;
  align-items: center;
  padding: 14px 18px;
  min-width: 880px;
}

.table-header {
  background: var(--c-bg);
  font-size: 0.78rem;
  font-weight: 700;
  color: var(--c-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.table-row {
  border-top: 1px solid var(--c-border-light);
  transition: background var(--t-fast);
}
.table-row:hover { background: var(--c-bg); }

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-avatar-mini {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  background: #64748b;
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 800;
  font-size: 0.85rem;
  flex-shrink: 0;
}

.user-avatar-mini.vip-avatar {
  background: linear-gradient(135deg, #f59e0b, #d97706);
  box-shadow: 0 2px 8px rgba(245, 158, 11, 0.35);
}

.user-email-col {
  display: flex;
  flex-direction: column;
}

.email-text {
  font-weight: 600;
  font-size: 0.9rem;
  color: var(--c-text);
}

.user-id-sub {
  font-size: 0.72rem;
  color: var(--c-text-tertiary);
  font-family: monospace;
}

.login-meta {
  display: flex;
  flex-direction: column;
}

.login-time {
  font-size: 0.82rem;
  color: var(--c-text);
  font-weight: 600;
}

.login-ip {
  color: var(--c-text-tertiary);
}

.col-actions {
  display: flex;
  gap: 6px;
  justify-content: flex-end;
  flex-wrap: wrap;
}

.btn-icon-text {
  font-size: 0.82rem;
  padding: 4px 8px;
}

.text-amber { color: #d97706 !important; }
.text-danger { color: #ef4444 !important; }

.badge-gold {
  background: rgba(245, 158, 11, 0.15);
  color: #d97706;
  border: 1px solid rgba(245, 158, 11, 0.35);
}
.badge-info {
  background: rgba(59, 130, 246, 0.15);
  color: #2563eb;
  border: 1px solid rgba(59, 130, 246, 0.3);
}

.table-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 18px;
  border-top: 1px solid var(--c-border-light);
  background: var(--c-bg);
}

.pagination-info { font-size: 0.85rem; color: var(--c-text-secondary); }
.pagination { display: flex; align-items: center; gap: 10px; }
.page-num { font-weight: 700; font-size: 0.9rem; }

.user-sub {
  font-size: 0.85rem;
  color: var(--c-text-secondary);
  margin-bottom: 18px;
}

.form-tip {
  font-size: 0.78rem;
  color: var(--c-text-tertiary);
  margin-top: 4px;
  display: block;
}

.btn-xs {
  font-size: 0.76rem;
  padding: 2px 6px;
  color: var(--c-accent);
  background: none;
  border: none;
  cursor: pointer;
}
.btn-xs:hover { text-decoration: underline; }
</style>
