<template>
  <div class="admin-products">
    <!-- Toolbar -->
    <div class="toolbar">
      <div class="toolbar-left">
        <input v-model="search" type="text" class="form-input search-input" :placeholder="t('nav.search')" @input="handleSearch" />
        <select v-model="filterCat" class="form-input form-select" style="width:auto" @change="loadProducts(true)">
          <option value="">{{ t('common.all') }}</option>
          <option v-for="c in categories" :key="c.value" :value="c.value">{{ locale === 'zh' ? c.zh : c.en }}</option>
        </select>
      </div>
      <div class="toolbar-right">
        <router-link to="/admin/products/new" class="btn btn-primary">
          + {{ t('admin.add_product') }}
        </router-link>
      </div>
    </div>

    <!-- Batch Action Floating Bar -->
    <transition name="slide-down">
      <div v-if="selectedIds.length > 0" class="batch-action-bar card">
        <div class="batch-left">
          <span class="batch-count-badge">✓ {{ selectedIds.length }}</span>
          <span class="batch-tip">
            {{ locale === 'zh' ? `已选中 ${selectedIds.length} 款壁纸商品` : `${selectedIds.length} wallpapers selected` }}
          </span>
        </div>
        <div class="batch-actions">
          <button class="btn btn-sm btn-accent" :disabled="batchProcessing" @click="handleBatchAction('pin')">
            📌 {{ locale === 'zh' ? '批量置顶' : 'Batch Pin' }}
          </button>
          <button class="btn btn-sm btn-outline" :disabled="batchProcessing" @click="handleBatchAction('unpin')">
            📍 {{ locale === 'zh' ? '取消置顶' : 'Batch Unpin' }}
          </button>
          <button class="btn btn-sm btn-primary" :disabled="batchProcessing" @click="handleBatchAction('publish')">
            🚀 {{ locale === 'zh' ? '批量上架' : 'Batch Publish' }}
          </button>
          <button class="btn btn-sm btn-outline" :disabled="batchProcessing" @click="handleBatchAction('unpublish')">
            ⏸️ {{ locale === 'zh' ? '批量下架' : 'Batch Unpublish' }}
          </button>
          <button class="btn btn-sm btn-danger" :disabled="batchProcessing" @click="confirmBatchDelete">
            🗑️ {{ locale === 'zh' ? '批量删除' : 'Batch Delete' }}
          </button>
          <button class="btn btn-sm btn-ghost" @click="selectedIds = []">
            ✕ {{ locale === 'zh' ? '取消全选' : 'Deselect All' }}
          </button>
        </div>
      </div>
    </transition>

    <!-- Table -->
    <div class="products-table card">
      <div class="table-header">
        <span class="col-checkbox">
          <input
            type="checkbox"
            :checked="isAllSelected"
            :indeterminate="isIndeterminate"
            @change="toggleSelectAll"
          />
        </span>
        <span class="col-thumb">{{ locale === 'zh' ? '图片' : 'Image' }}</span>
        <span class="col-title">{{ locale === 'zh' ? '标题 / 排序' : 'Title & Sort' }}</span>
        <span class="col-cat">{{ t('admin.category') }}</span>
        <span class="col-price">{{ locale === 'zh' ? '双轨价格' : 'Price' }}</span>
        <span class="col-status">{{ locale === 'zh' ? '状态' : 'Status' }}</span>
        <span class="col-actions">{{ locale === 'zh' ? '操作' : 'Actions' }}</span>
      </div>

      <div v-if="loading" class="table-loading"><span class="spinner"></span></div>

      <div v-else-if="products.length === 0" class="table-empty">{{ t('common.no_data') }}</div>

      <div v-else>
        <div
          v-for="product in products"
          :key="product.id"
          class="table-row"
          :class="{ 'row-selected': selectedIds.includes(product.id), 'row-pinned': product.is_pinned }"
        >
          <div class="col-checkbox">
            <input
              type="checkbox"
              :checked="selectedIds.includes(product.id)"
              @change="toggleSelectRow(product.id)"
            />
          </div>
          <div class="col-thumb">
            <div class="product-thumb">
              <img v-if="product.cover_image" :src="product.cover_image" />
              <div v-else class="thumb-ph" :style="{ background: product._color || '#667eea' }">🖼️</div>
            </div>
          </div>
          <div class="col-title">
            <div class="title-top-meta">
              <span v-if="product.is_pinned" class="badge-pinned-tag" title="当前已置顶，最新置顶排在最前">
                📌 {{ locale === 'zh' ? '已置顶' : 'Pinned' }}
              </span>
              <span v-if="product.sort_order" class="badge-sort-tag" :title="`排序权重: ${product.sort_order} (数字越大越靠前)`">
                权重: {{ product.sort_order }}
              </span>
            </div>
            <p class="product-name">{{ locale === 'zh' ? product.title_zh : product.title_en }}</p>
            <p class="product-id">ID: {{ product.id?.slice(-8) }}</p>
          </div>
          <div class="col-cat">
            <span class="badge badge-neutral">{{ getCat(product.category) }}</span>
          </div>
          <div class="col-price">
            <strong>$ {{ product.price }}</strong>
            <small v-if="product.price_cny" class="cny-hint"> (¥{{ product.price_cny }})</small>
          </div>
          <div class="col-status">
            <span class="badge" :class="product.is_published ? 'badge-success' : 'badge-neutral'">
              {{ product.is_published ? (locale === 'zh' ? '已上架' : 'Live') : (locale === 'zh' ? '已下架' : 'Draft') }}
            </span>
          </div>
          <div class="col-actions">
            <button
              class="btn btn-sm"
              :class="product.is_pinned ? 'btn-pinned-active' : 'btn-outline'"
              @click="togglePin(product)"
              :title="product.is_pinned ? '已置顶，点击取消置顶' : '点击将此商品置顶至最前'"
            >
              📌 {{ product.is_pinned ? (locale === 'zh' ? '取消置顶' : 'Unpin') : (locale === 'zh' ? '置顶' : 'Pin') }}
            </button>
            <button
              class="btn btn-ghost btn-sm"
              @click="editSortOrder(product)"
              :title="locale === 'zh' ? '修改排序权重 (数字越大排序越靠前)' : 'Edit sort order'"
            >
              🔢 {{ product.sort_order || 0 }}
            </button>
            <router-link :to="`/admin/products/${product.id}/edit`" class="btn btn-ghost btn-sm">
              ✏️ {{ t('common.edit') }}
            </router-link>
            <button class="btn btn-sm" :class="product.is_published ? 'btn-outline' : 'btn-primary'"
              @click="togglePublish(product)">
              {{ product.is_published ? t('admin.unpublish') : t('admin.publish') }}
            </button>
            <button class="btn btn-sm btn-danger" @click="confirmDelete(product)">
              🗑️
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Pagination -->
    <div class="pagination">
      <button class="btn btn-ghost btn-sm" :disabled="page <= 1" @click="prevPage">← {{ locale === 'zh' ? '上一页' : 'Prev' }}</button>
      <span class="page-info">{{ locale === 'zh' ? `第 ${page} 页` : `Page ${page}` }}</span>
      <button class="btn btn-ghost btn-sm" :disabled="!hasMore" @click="nextPage">{{ locale === 'zh' ? '下一页' : 'Next' }} →</button>
    </div>

    <!-- Single Delete confirm modal -->
    <div v-if="deleteTarget" class="modal-overlay" @click.self="deleteTarget = null">
      <div class="modal-box" style="max-width:400px">
        <div class="modal-body" style="padding:28px;text-align:center">
          <p style="font-size:2rem;margin-bottom:12px">🗑️</p>
          <h3 style="margin-bottom:8px">{{ t('admin.confirm_delete') }}</h3>
          <p style="color:var(--c-text-secondary);font-size:0.9rem;margin-bottom:24px">
            {{ locale === 'zh' ? deleteTarget.title_zh : deleteTarget.title_en }}
          </p>
          <div style="display:flex;gap:10px;justify-content:center">
            <button class="btn btn-outline" @click="deleteTarget = null">{{ t('common.cancel') }}</button>
            <button class="btn btn-danger" :disabled="deleting" @click="doDelete">
              <span v-if="deleting" class="spinner"></span>
              {{ t('common.delete') }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Batch Delete confirm modal -->
    <div v-if="showBatchDeleteModal" class="modal-overlay" @click.self="showBatchDeleteModal = false">
      <div class="modal-box" style="max-width:420px">
        <div class="modal-body" style="padding:28px;text-align:center">
          <p style="font-size:2.4rem;margin-bottom:12px">⚠️</p>
          <h3 style="margin-bottom:8px">{{ locale === 'zh' ? '批量删除商品确认' : 'Confirm Batch Delete' }}</h3>
          <p style="color:var(--c-text-secondary);font-size:0.9rem;margin-bottom:24px">
            {{ locale === 'zh' ? `您已选择 ${selectedIds.length} 款壁纸，删除后不可恢复，是否确定全部删除？` : `Are you sure you want to delete ${selectedIds.length} products permanently?` }}
          </p>
          <div style="display:flex;gap:10px;justify-content:center">
            <button class="btn btn-outline" @click="showBatchDeleteModal = false">{{ t('common.cancel') }}</button>
            <button class="btn btn-danger" :disabled="batchProcessing" @click="executeBatchDelete">
              <span v-if="batchProcessing" class="spinner"></span>
              {{ locale === 'zh' ? `确认删除 (${selectedIds.length})` : 'Confirm Delete' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { productAPI, categoryAPI } from '@/api'
import { useToast } from '@/composables/useToast'
import { useDebounceFn } from '@vueuse/core'

const { t, locale } = useI18n()
const { success, error } = useToast()

const products = ref([])
const loading = ref(true)
const page = ref(1)
const hasMore = ref(true)
const search = ref('')
const filterCat = ref('')
const deleteTarget = ref(null)
const deleting = ref(false)

// Batch Operation State
const selectedIds = ref([])
const batchProcessing = ref(false)
const showBatchDeleteModal = ref(false)

const isAllSelected = computed(() => {
  return products.value.length > 0 && products.value.every(p => selectedIds.value.includes(p.id))
})

const isIndeterminate = computed(() => {
  const count = selectedIds.value.length
  return count > 0 && count < products.value.length
})

function toggleSelectAll(e) {
  if (e.target.checked) {
    selectedIds.value = products.value.map(p => p.id)
  } else {
    selectedIds.value = []
  }
}

function toggleSelectRow(id) {
  const idx = selectedIds.value.indexOf(id)
  if (idx > -1) {
    selectedIds.value.splice(idx, 1)
  } else {
    selectedIds.value.push(id)
  }
}

async function handleBatchAction(action) {
  if (selectedIds.value.length === 0) return
  batchProcessing.value = true
  try {
    const res = await productAPI.batchAction(action, selectedIds.value)
    success(res.message || (locale.value === 'zh' ? '批量操作已完成' : 'Batch action completed'))
    selectedIds.value = []
    loadProducts()
  } catch (err) {
    error(err.message || (locale.value === 'zh' ? '批量操作失败' : 'Batch action failed'))
  } finally {
    batchProcessing.value = false
  }
}

function confirmBatchDelete() {
  if (selectedIds.value.length === 0) return
  showBatchDeleteModal.value = true
}

async function executeBatchDelete() {
  batchProcessing.value = true
  try {
    const res = await productAPI.batchAction('delete', selectedIds.value)
    success(res.message || (locale.value === 'zh' ? '批量删除成功' : 'Batch delete successful'))
    showBatchDeleteModal.value = false
    selectedIds.value = []
    loadProducts(true)
  } catch (err) {
    error(err.message || '删除失败')
  } finally {
    batchProcessing.value = false
  }
}

const categories = ref([
  { value: 'ios', zh: '📱 iOS 风格', en: '📱 iOS Theme' },
  { value: 'anime', zh: '✨ 动漫二次元', en: '✨ Anime' },
  { value: 'cyberpunk', zh: '🏙️ 赛博科技', en: '🏙️ Cyberpunk' },
  { value: 'nature', zh: '🌿 治愈风景', en: '🌿 Nature' },
  { value: 'minimal', zh: '🌸 极简冷淡', en: '🌸 Minimal' }
])

async function loadCategories() {
  try {
    const res = await categoryAPI.adminList()
    const items = res.items || res.categories || res || []
    if (items.length > 0) {
      categories.value = items.map(c => ({
        value: c.code,
        zh: `${c.icon || ''} ${c.name_zh}`.trim(),
        en: `${c.icon || ''} ${c.name_en || c.name_zh}`.trim()
      }))
    }
  } catch (e) {}
}

const catMap = computed(() => Object.fromEntries(categories.value.map(c => [c.value, c])))
function getCat(cat) {
  if (!cat) return '-'
  const item = catMap.value[cat]
  return item ? (locale.value === 'zh' ? item.zh : item.en) : cat
}

async function loadProducts(reset = false) {
  if (reset) { page.value = 1; products.value = [] }
  loading.value = true
  try {
    const res = await productAPI.list({
      page: page.value,
      page_size: 15,
      limit: 15,
      search: search.value,
      category: filterCat.value
    })
    const items = res.items || res.products || (Array.isArray(res) ? res : [])
    products.value = items
    const total = res.total || 0
    hasMore.value = total > 0 ? (page.value * 15 < total) : (items.length === 15)
  } catch (err) {
    console.error('Failed to load products:', err)
    products.value = []
    hasMore.value = false
  } finally { loading.value = false }
}

const handleSearch = useDebounceFn(() => loadProducts(true), 400)

function prevPage() { if (page.value > 1) { page.value--; loadProducts() } }
function nextPage() { if (hasMore.value) { page.value++; loadProducts() } }

async function togglePin(product) {
  try {
    const res = await productAPI.togglePin(product.id)
    product.is_pinned = res.is_pinned
    product.pinned_at = res.pinned_at
    success(res.message || (product.is_pinned ? '商品已置顶' : '已取消置顶'))
    loadProducts()
  } catch (err) {
    error(err.message || '操作失败')
  }
}

async function editSortOrder(product) {
  const current = product.sort_order || 0
  const val = prompt(
    locale.value === 'zh'
      ? `请输入「${product.title_zh}」的排序权重\n(0为默认，数字越大排序越靠前)：`
      : `Enter sort order for "${product.title_en || product.title_zh}"\n(0 is default, higher numbers sort first):`,
    String(current)
  )
  if (val === null) return
  const num = parseInt(val, 10)
  if (isNaN(num)) {
    alert(locale.value === 'zh' ? '请输入合法整数数字' : 'Please enter a valid integer')
    return
  }
  try {
    await productAPI.updateSort(product.id, num)
    product.sort_order = num
    success(locale.value === 'zh' ? '排序权重更新成功' : 'Sort order updated')
    loadProducts()
  } catch (err) {
    error(err.message || '更新失败')
  }
}

async function togglePublish(product) {
  try {
    await productAPI.update(product.id, { is_published: !product.is_published })
    product.is_published = !product.is_published
    success(t('common.success'))
  } catch { error(t('common.error')) }
}

function confirmDelete(product) { deleteTarget.value = product }

async function doDelete() {
  deleting.value = true
  try {
    await productAPI.delete(deleteTarget.value.id)
    products.value = products.value.filter(p => p.id !== deleteTarget.value.id)
    deleteTarget.value = null
    success(t('common.success'))
  } catch { error(t('common.error')) }
  finally { deleting.value = false }
}

onMounted(() => {
  loadCategories()
  loadProducts(true)
})
</script>

<style scoped>
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}
.toolbar-left { display: flex; gap: 8px; flex-wrap: wrap; }
.search-input { width: 220px; border-radius: var(--r-full); padding: 8px 16px; font-size: 0.88rem; }

/* Batch Action Bar */
.batch-action-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 20px;
  margin-bottom: 16px;
  background: var(--c-surface);
  border: 2px solid var(--c-accent);
  border-radius: 16px;
  box-shadow: 0 6px 20px rgba(0, 113, 227, 0.12);
  gap: 16px;
  flex-wrap: wrap;
}

.batch-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.batch-count-badge {
  background: var(--c-accent);
  color: #ffffff;
  padding: 3px 10px;
  border-radius: 999px;
  font-size: 0.82rem;
  font-weight: 700;
}

.batch-tip {
  font-weight: 700;
  font-size: 0.92rem;
  color: var(--c-text);
}

.batch-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.products-table {
  border: 1px solid var(--c-border-light);
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.table-header, .table-row {
  display: grid;
  grid-template-columns: 36px 60px minmax(180px, 1fr) 110px 110px 90px auto;
  gap: 12px;
  align-items: center;
  padding: 12px 16px;
  min-width: 760px;
}

.col-checkbox {
  display: flex;
  align-items: center;
  justify-content: center;
}

.col-checkbox input[type="checkbox"] {
  width: 17px;
  height: 17px;
  cursor: pointer;
  accent-color: var(--c-accent);
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
.table-row.row-selected {
  background: rgba(0, 113, 227, 0.05);
}

.product-thumb {
  width: 40px;
  height: 54px;
  border-radius: var(--r-xs);
  overflow: hidden;
}
.product-thumb img { width: 100%; height: 100%; object-fit: cover; }
.thumb-ph { width: 100%; height: 100%; display: flex; align-items: center; justify-content: center; font-size: 0.9rem; }

.product-name { font-size: 0.88rem; font-weight: 600; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.product-id { font-size: 0.72rem; color: var(--c-text-tertiary); font-family: monospace; }
.cny-hint { font-size: 0.78rem; color: var(--c-text-tertiary); font-weight: normal; }

.col-actions { display: flex; gap: 6px; justify-content: flex-end; }

.table-loading, .table-empty {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 40px;
  color: var(--c-text-tertiary);
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  margin-top: 20px;
}
.page-info { font-size: 0.88rem; color: var(--c-text-secondary); }

.slide-down-enter-active, .slide-down-leave-active {
  transition: all 0.25s ease;
}
.slide-down-enter-from, .slide-down-leave-to {
  transform: translateY(-10px);
  opacity: 0;
}

.table-row.row-pinned {
  background: rgba(239, 68, 68, 0.04);
  border-left: 3px solid #ef4444;
}

.title-top-meta {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 3px;
}

.badge-pinned-tag {
  display: inline-block;
  padding: 1px 7px;
  border-radius: 4px;
  font-size: 0.72rem;
  font-weight: 700;
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.badge-sort-tag {
  display: inline-block;
  padding: 1px 6px;
  border-radius: 4px;
  font-size: 0.7rem;
  font-weight: 600;
  background: rgba(168, 85, 247, 0.12);
  color: #c084fc;
}

.btn-pinned-active {
  background: rgba(239, 68, 68, 0.15) !important;
  color: #ef4444 !important;
  border: 1px solid rgba(239, 68, 68, 0.4) !important;
  font-weight: 700;
}

.btn-accent {
  background: linear-gradient(135deg, #ef4444, #dc2626);
  color: #ffffff;
  border: none;
}
.btn-accent:hover {
  opacity: 0.9;
}
</style>
