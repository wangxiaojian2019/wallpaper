<template>
  <div class="admin-categories-page">
    <div class="page-header-row">
      <div>
        <h2 class="section-title">{{ locale === 'zh' ? '分类管理' : 'Category Management' }}</h2>
        <p class="section-subtitle">
          {{ locale === 'zh' ? '自由配置前台壁纸分类、排序权重与启用状态，支持双语显示与个性化 Emoji 图标' : 'Manage wallpaper categories, sorting weights, and visibility across store.' }}
        </p>
      </div>
      <div class="header-actions">
        <button class="btn btn-outline" @click="fetchCategories" :disabled="loading">
          🔄 {{ locale === 'zh' ? '刷新' : 'Refresh' }}
        </button>
        <button class="btn btn-primary" @click="openCreateModal">
          ＋ {{ locale === 'zh' ? '新建分类' : 'New Category' }}
        </button>
      </div>
    </div>

    <!-- Category Table Card -->
    <div class="card table-card">
      <div v-if="loading" class="table-loading">
        <span class="spinner"></span>
      </div>

      <div v-else-if="categories.length === 0" class="table-empty">
        <p>📁 {{ locale === 'zh' ? '暂无分类，请点击右上角新建' : 'No categories yet. Click New Category above.' }}</p>
      </div>

      <div v-else class="table-responsive">
        <table class="modern-table">
          <thead>
            <tr>
              <th style="width: 70px;">{{ locale === 'zh' ? '图标' : 'Icon' }}</th>
              <th>{{ locale === 'zh' ? '中文名称' : 'Name (ZH)' }}</th>
              <th>{{ locale === 'zh' ? '英文名称' : 'Name (EN)' }}</th>
              <th>{{ locale === 'zh' ? '唯一编码' : 'Code' }}</th>
              <th style="width: 110px;">{{ locale === 'zh' ? '排序权重' : 'Sort Order' }}</th>
              <th style="width: 100px;">{{ locale === 'zh' ? '前台展示' : 'Status' }}</th>
              <th style="width: 160px; text-align: right;">{{ locale === 'zh' ? '操作' : 'Actions' }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="cat in categories" :key="cat.id" :class="{ 'row-disabled': !cat.is_enabled }">
              <td class="cell-icon">
                <span class="icon-bubble">{{ cat.icon || '📁' }}</span>
              </td>
              <td>
                <span class="font-bold text-main">{{ cat.name_zh }}</span>
              </td>
              <td>
                <span class="text-secondary">{{ cat.name_en || '-' }}</span>
              </td>
              <td>
                <span class="code-badge">{{ cat.code }}</span>
              </td>
              <td>
                <span class="badge badge-sort" :title="locale === 'zh' ? '数字越大排序越靠前' : 'Higher number sorts first'">
                  {{ cat.sort_order }}
                </span>
              </td>
              <td>
                <button
                  class="status-toggle-pill"
                  :class="cat.is_enabled ? 'is-active' : 'is-disabled'"
                  @click="toggleCategoryStatus(cat)"
                  :title="cat.is_enabled ? '点击隐藏' : '点击显示'"
                >
                  <span class="toggle-dot"></span>
                  <span>{{ cat.is_enabled ? (locale === 'zh' ? '已展示' : 'Enabled') : (locale === 'zh' ? '已隐藏' : 'Disabled') }}</span>
                </button>
              </td>
              <td style="text-align: right;">
                <div class="action-btn-group">
                  <button class="btn btn-ghost btn-sm" @click="openEditModal(cat)">
                    ✏️ {{ locale === 'zh' ? '编辑' : 'Edit' }}
                  </button>
                  <button class="btn btn-ghost btn-sm btn-danger-text" @click="confirmDelete(cat)">
                    🗑️
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create / Edit Modal -->
    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal-box modal-category">
        <div class="modal-header">
          <h3>{{ editingId ? (locale === 'zh' ? '编辑分类' : 'Edit Category') : (locale === 'zh' ? '新建分类' : 'New Category') }}</h3>
          <button class="btn-close" @click="showModal = false">✕</button>
        </div>

        <form @submit.prevent="submitCategory" class="modal-body">
          <div class="form-row-dual">
            <div class="form-group flex-1">
              <label class="form-label required">{{ locale === 'zh' ? '中文名称' : 'Name (ZH)' }}</label>
              <input
                v-model="form.name_zh"
                type="text"
                class="form-input"
                placeholder="如：iOS 风格、赛博科技"
                required
              />
            </div>
            <div class="form-group" style="width: 100px;">
              <label class="form-label">{{ locale === 'zh' ? '图标 Emoji' : 'Icon' }}</label>
              <input
                v-model="form.icon"
                type="text"
                class="form-input text-center"
                placeholder="📱"
                maxlength="8"
              />
            </div>
          </div>

          <div class="form-row-dual">
            <div class="form-group flex-1">
              <label class="form-label">{{ locale === 'zh' ? '英文名称' : 'Name (EN)' }}</label>
              <input
                v-model="form.name_en"
                type="text"
                class="form-input"
                placeholder="e.g. iOS Theme, Cyberpunk"
              />
            </div>
            <div class="form-group flex-1">
              <label class="form-label required">{{ locale === 'zh' ? '唯一标识编码' : 'Code' }}</label>
              <input
                v-model="form.code"
                type="text"
                class="form-input"
                placeholder="如：ios, cyberpunk"
                required
              />
            </div>
          </div>

          <div class="form-row-dual">
            <div class="form-group flex-1">
              <label class="form-label">{{ locale === 'zh' ? '排序权重' : 'Sort Order' }}</label>
              <input
                v-model.number="form.sort_order"
                type="number"
                class="form-input"
                placeholder="0"
              />
              <span class="form-hint">{{ locale === 'zh' ? '0为默认，数字越大前台排在越前面' : '0 is default, higher numbers appear first' }}</span>
            </div>
            <div class="form-group flex-1" style="display: flex; flex-direction: column; justify-content: center;">
              <label class="form-label">{{ locale === 'zh' ? '前台展示状态' : 'Visibility' }}</label>
              <label class="checkbox-label" style="margin-top: 6px; cursor: pointer;">
                <input v-model="form.is_enabled" type="checkbox" />
                <span>{{ locale === 'zh' ? '在前台画廊展示此分类' : 'Show in store gallery' }}</span>
              </label>
            </div>
          </div>

          <div v-if="modalError" class="modal-error-alert">
            ⚠️ {{ modalError }}
          </div>

          <div class="modal-footer">
            <button type="button" class="btn btn-outline" @click="showModal = false">
              {{ locale === 'zh' ? '取消' : 'Cancel' }}
            </button>
            <button type="submit" class="btn btn-primary" :disabled="submitting">
              <span v-if="submitting" class="spinner"></span>
              {{ locale === 'zh' ? '保存提交' : 'Save' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="deleteTarget" class="modal-overlay" @click.self="deleteTarget = null">
      <div class="modal-box" style="max-width: 420px; text-align: center;">
        <div style="padding: 28px 24px;">
          <p style="font-size: 2.2rem; margin-bottom: 12px;">🗑️</p>
          <h3 style="margin-bottom: 8px;">{{ locale === 'zh' ? '确认删除分类？' : 'Delete Category?' }}</h3>
          <p style="color: var(--c-text-secondary); font-size: 0.9rem; margin-bottom: 24px;">
            {{ locale === 'zh' ? `确认删除分类「${deleteTarget.name_zh}」吗？已关联该分类的壁纸不受影响。` : `Are you sure you want to delete "${deleteTarget.name_zh}"? Wallpapers with this category will not be deleted.` }}
          </p>
          <div style="display: flex; gap: 12px; justify-content: center;">
            <button class="btn btn-outline" @click="deleteTarget = null">
              {{ locale === 'zh' ? '取消' : 'Cancel' }}
            </button>
            <button class="btn btn-danger" :disabled="deleting" @click="doDelete">
              <span v-if="deleting" class="spinner"></span>
              {{ locale === 'zh' ? '确认删除' : 'Delete' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { categoryAPI } from '@/api'

const locale = ref(localStorage.getItem('lang') || 'zh')
const categories = ref([])
const loading = ref(false)
const showModal = ref(false)
const editingId = ref(null)
const submitting = ref(false)
const modalError = ref('')
const deleteTarget = ref(null)
const deleting = ref(false)

const form = ref({
  name_zh: '',
  name_en: '',
  code: '',
  icon: '📁',
  sort_order: 0,
  is_enabled: true
})

async function fetchCategories() {
  loading.value = true
  try {
    const res = await categoryAPI.adminList()
    categories.value = res.items || res.categories || res || []
  } catch (err) {
    console.error('Failed to load categories:', err)
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  editingId.value = null
  modalError.value = ''
  form.value = {
    name_zh: '',
    name_en: '',
    code: '',
    icon: '📁',
    sort_order: 0,
    is_enabled: true
  }
  showModal.value = true
}

function openEditModal(cat) {
  editingId.value = cat.id
  modalError.value = ''
  form.value = {
    name_zh: cat.name_zh || '',
    name_en: cat.name_en || '',
    code: cat.code || '',
    icon: cat.icon || '📁',
    sort_order: cat.sort_order || 0,
    is_enabled: cat.is_enabled !== false
  }
  showModal.value = true
}

async function submitCategory() {
  submitting.value = true
  modalError.value = ''
  try {
    if (editingId.value) {
      await categoryAPI.update(editingId.value, form.value)
    } else {
      await categoryAPI.create(form.value)
    }
    showModal.value = false
    await fetchCategories()
  } catch (err) {
    modalError.value = err.message || '操作失败，请重试'
  } finally {
    submitting.value = false
  }
}

async function toggleCategoryStatus(cat) {
  try {
    const updatedStatus = !cat.is_enabled
    await categoryAPI.update(cat.id, {
      ...cat,
      is_enabled: updatedStatus
    })
    cat.is_enabled = updatedStatus
  } catch (err) {
    alert(err.message || '状态切换失败')
  }
}

function confirmDelete(cat) {
  deleteTarget.value = cat
}

async function doDelete() {
  if (!deleteTarget.value) return
  deleting.value = true
  try {
    await categoryAPI.delete(deleteTarget.value.id)
    deleteTarget.value = null
    await fetchCategories()
  } catch (err) {
    alert(err.message || '删除失败')
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  fetchCategories()
})
</script>

<style scoped>
.admin-categories-page {
  padding-bottom: 40px;
}

.page-header-row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 24px;
}

.section-title {
  font-size: 1.4rem;
  font-weight: 700;
  margin-bottom: 4px;
}

.section-subtitle {
  color: var(--c-text-secondary);
  font-size: 0.88rem;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.table-card {
  padding: 0;
  overflow: hidden;
  border-radius: 14px;
  background: var(--c-card-bg, #1a1b24);
  border: 1px solid var(--c-border, rgba(255, 255, 255, 0.08));
}

.table-loading, .table-empty {
  padding: 60px 20px;
  text-align: center;
  color: var(--c-text-secondary);
}

.table-responsive {
  width: 100%;
  overflow-x: auto;
}

.modern-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 0.9rem;
}

.modern-table th {
  padding: 14px 18px;
  font-size: 0.78rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: var(--c-text-secondary);
  background: rgba(255, 255, 255, 0.02);
  border-bottom: 1px solid var(--c-border, rgba(255, 255, 255, 0.08));
}

.modern-table td {
  padding: 16px 18px;
  vertical-align: middle;
  border-bottom: 1px solid var(--c-border, rgba(255, 255, 255, 0.05));
}

.modern-table tr:hover {
  background: rgba(255, 255, 255, 0.02);
}

.row-disabled {
  opacity: 0.55;
}

.cell-icon {
  text-align: center;
}

.icon-bubble {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 38px;
  height: 38px;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.06);
  font-size: 1.25rem;
}

.text-main {
  color: var(--c-text, #ffffff);
}

.text-secondary {
  color: var(--c-text-secondary, #94a3b8);
}

.code-badge {
  display: inline-block;
  padding: 3px 8px;
  border-radius: 6px;
  background: rgba(56, 189, 248, 0.1);
  color: #38bdf8;
  font-family: monospace;
  font-size: 0.82rem;
  font-weight: 600;
}

.badge-sort {
  display: inline-block;
  padding: 3px 8px;
  border-radius: 6px;
  background: rgba(168, 85, 247, 0.12);
  color: #c084fc;
  font-weight: 700;
}

.status-toggle-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  border-radius: 999px;
  font-size: 0.78rem;
  font-weight: 600;
  cursor: pointer;
  border: 1px solid transparent;
  transition: all 0.2s ease;
}

.status-toggle-pill.is-active {
  background: rgba(16, 185, 129, 0.12);
  color: #10b981;
  border-color: rgba(16, 185, 129, 0.3);
}

.status-toggle-pill.is-disabled {
  background: rgba(100, 116, 139, 0.15);
  color: #94a3b8;
  border-color: rgba(100, 116, 139, 0.3);
}

.toggle-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
}

.action-btn-group {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.btn-danger-text:hover {
  color: #ef4444;
}

/* Modal */
.modal-category {
  max-width: 520px;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid var(--c-border, rgba(255, 255, 255, 0.08));
}

.modal-header h3 {
  font-size: 1.15rem;
  font-weight: 700;
}

.btn-close {
  background: none;
  border: none;
  color: var(--c-text-secondary);
  font-size: 1.1rem;
  cursor: pointer;
}

.modal-body {
  padding: 24px;
}

.form-row-dual {
  display: flex;
  gap: 14px;
  margin-bottom: 16px;
}

.flex-1 {
  flex: 1;
}

.form-label.required::after {
  content: ' *';
  color: #ef4444;
}

.form-hint {
  display: block;
  font-size: 0.75rem;
  color: var(--c-text-secondary);
  margin-top: 4px;
}

.text-center {
  text-align: center;
}

.checkbox-label {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-size: 0.88rem;
}

.modal-error-alert {
  padding: 10px 14px;
  border-radius: 8px;
  background: rgba(239, 68, 68, 0.12);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: #f87171;
  font-size: 0.84rem;
  margin-bottom: 16px;
}

.modal-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}
</style>
