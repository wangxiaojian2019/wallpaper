<template>
  <div class="admin-redeem-codes">
    <!-- Header -->
    <div class="page-header">
      <div>
        <h2 class="page-title">🎫 {{ locale === 'zh' ? '口令兑换与营销卡密管理' : 'Redeem Codes & Vouchers' }}</h2>
        <p class="page-sub">
          {{ locale === 'zh'
            ? '自定义生成和管理全站口令卡密，支持免费兑换限额以下壁纸（如5元壁纸免费兑换）、开通指定VIP会员、自定义兑换次数与有效期'
            : 'Create and manage promo codes for price-limit free downloads, VIP membership grants, with custom usage limits and expiry' }}
        </p>
      </div>
      <div class="header-actions">
        <button class="btn btn-primary" @click="openCreateModal">
          ➕ {{ locale === 'zh' ? '生成新兑换口令' : 'Create Codes' }}
        </button>
      </div>
    </div>

    <!-- Main Navigation Tabs -->
    <div class="nav-tabs-bar">
      <button
        class="tab-btn"
        :class="{ active: currentTab === 'codes' }"
        @click="currentTab = 'codes'"
      >
        <span class="tab-icon">🏷️</span>
        <span>{{ locale === 'zh' ? '口令卡密列表' : 'Active Codes' }}</span>
        <span v-if="totalCodes > 0" class="tab-badge">{{ totalCodes }}</span>
      </button>

      <button
        class="tab-btn"
        :class="{ active: currentTab === 'logs' }"
        @click="switchTab('logs')"
      >
        <span class="tab-icon">📜</span>
        <span>{{ locale === 'zh' ? '兑换记录流水' : 'Redeem Logs' }}</span>
        <span v-if="totalLogs > 0" class="tab-badge">{{ totalLogs }}</span>
      </button>
    </div>

    <!-- TAB 1: CODES LIST -->
    <div v-show="currentTab === 'codes'" class="tab-content">
      <!-- Toolbar Filter -->
      <div class="card filter-bar">
        <div class="filter-inputs">
          <div class="filter-item">
            <label class="filter-label">{{ locale === 'zh' ? '搜索口令/备注' : 'Search Code/Remark' }}</label>
            <input
              v-model="search"
              type="text"
              class="form-input form-input-sm"
              :placeholder="locale === 'zh' ? '输入口令如 VIP666' : 'Enter code keyword'"
              @keyup.enter="loadCodes(true)"
            />
          </div>

          <div class="filter-item">
            <label class="filter-label">{{ locale === 'zh' ? '兑换用途' : 'Purpose' }}</label>
            <select v-model="filterPurpose" class="form-input form-input-sm form-select" @change="loadCodes(true)">
              <option value="">{{ locale === 'zh' ? '全部用途' : 'All Purposes' }}</option>
              <option value="price_limit">🏷️ {{ locale === 'zh' ? '限额免单 (如5元及以下免费)' : 'Price-Limit Free' }}</option>
              <option value="membership">👑 {{ locale === 'zh' ? '开通 VIP 会员' : 'Grant VIP' }}</option>
              <option value="product">🖼️ {{ locale === 'zh' ? '指定特定壁纸' : 'Specific Wallpaper' }}</option>
            </select>
          </div>
        </div>

        <div class="filter-actions">
          <button class="btn btn-primary btn-sm" @click="loadCodes(true)">
            🔍 {{ locale === 'zh' ? '查询' : 'Search' }}
          </button>
          <button class="btn btn-ghost btn-sm" @click="resetCodeFilter">
            ✕ {{ locale === 'zh' ? '重置' : 'Reset' }}
          </button>
        </div>
      </div>

      <!-- Codes Table -->
      <div class="card table-card">
        <div class="table-wrap">
          <table class="data-table">
            <thead>
              <tr>
                <th>{{ locale === 'zh' ? '兑换口令' : 'Code' }}</th>
                <th>{{ locale === 'zh' ? '兑换用途' : 'Purpose' }}</th>
                <th>{{ locale === 'zh' ? '口令价值 / 兑换规则' : 'Value / Rule' }}</th>
                <th>{{ locale === 'zh' ? '使用进度 (已用/上限)' : 'Uses (Used/Max)' }}</th>
                <th>{{ locale === 'zh' ? '有效期限' : 'Validity' }}</th>
                <th>{{ locale === 'zh' ? '备注说明' : 'Remark' }}</th>
                <th>{{ locale === 'zh' ? '启用状态' : 'Status' }}</th>
                <th>{{ locale === 'zh' ? '操作' : 'Actions' }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="loadingCodes">
                <td colspan="8" class="text-center py-8">
                  <span class="spinner"></span>
                  <span class="ml-2 text-secondary">{{ locale === 'zh' ? '正在加载口令列表...' : 'Loading codes...' }}</span>
                </td>
              </tr>
              <tr v-else-if="codes.length === 0">
                <td colspan="8" class="text-center py-8 text-secondary">
                  {{ locale === 'zh' ? '暂未生成任何兑换口令，点击上方按钮立即生成' : 'No redeem codes created yet' }}
                </td>
              </tr>
              <tr v-else v-for="c in codes" :key="c.id" class="table-row-hover">
                <!-- Code text with copy -->
                <td>
                  <div class="code-copy-cell">
                    <span class="code-text font-mono font-bold">{{ c.code }}</span>
                    <button class="copy-mini-btn" @click="copyCode(c.code)" title="复制口令">
                      📋
                    </button>
                  </div>
                </td>

                <!-- Purpose -->
                <td>
                  <span class="badge" :class="getPurposeBadgeClass(c.purpose)">
                    {{ formatPurpose(c.purpose) }}
                  </span>
                </td>

                <!-- Value / Rule detail -->
                <td>
                  <div class="rule-detail-cell">
                    <template v-if="c.purpose === 'price_limit'">
                      <strong class="text-emerald font-bold">
                        {{ c.currency === 'USD' ? '$' : '¥' }}{{ Number(c.max_price).toFixed(2) }}
                      </strong>
                      <span class="rule-subtext text-xs text-secondary">
                        {{ locale === 'zh' ? '及以下商品免费兑换' : 'or less free unlock' }}
                      </span>
                    </template>
                    <template v-else-if="c.purpose === 'membership'">
                      <span class="vip-tier-tag font-bold">
                        👑 {{ formatTierName(c.target_tier) }}
                      </span>
                    </template>
                    <template v-else-if="c.purpose === 'product'">
                      <span class="font-mono text-xs text-secondary">
                        商品ID: {{ c.product_id?.slice(-8) || '特定商品' }}
                      </span>
                    </template>
                  </div>
                </td>

                <!-- Uses Progress -->
                <td>
                  <div class="uses-progress-wrap">
                    <div class="uses-count-row">
                      <strong>{{ c.used_count }}</strong>
                      <span class="text-secondary">/ {{ c.max_uses > 0 ? c.max_uses : '无限' }}</span>
                    </div>
                    <div v-if="c.max_uses > 0" class="progress-bar-bg">
                      <div
                        class="progress-bar-fill"
                        :style="{ width: Math.min(100, Math.round((c.used_count / c.max_uses) * 100)) + '%' }"
                        :class="{ 'fill-full': c.used_count >= c.max_uses }"
                      ></div>
                    </div>
                  </div>
                </td>

                <!-- Expiry -->
                <td>
                  <div class="expiry-cell">
                    <span v-if="!c.expires_at" class="badge badge-success-subtle">
                      ♾️ {{ locale === 'zh' ? '永久有效' : 'Permanent' }}
                    </span>
                    <template v-else>
                      <span :class="isExpired(c.expires_at) ? 'text-danger font-semibold' : 'text-secondary'">
                        {{ formatDate(c.expires_at) }}
                      </span>
                      <small v-if="isExpired(c.expires_at)" class="text-danger-tag">
                        ({{ locale === 'zh' ? '已过期' : 'Expired' }})
                      </small>
                    </template>
                  </div>
                </td>

                <!-- Remark -->
                <td class="text-secondary text-sm" style="max-width: 180px;">
                  {{ c.remark || '-' }}
                </td>

                <!-- Active switch -->
                <td>
                  <label class="switch-control">
                    <input type="checkbox" :checked="c.is_active" @change="toggleCodeActive(c)" />
                    <span class="switch-slider"></span>
                  </label>
                </td>

                <!-- Actions -->
                <td>
                  <button class="btn btn-sm btn-ghost text-danger" @click="deleteCode(c)" title="删除口令">
                    🗑️
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination for Codes -->
        <div v-if="totalCodes > pageSize" class="table-pagination">
          <div class="pagination-info">
            {{ locale === 'zh' ? `共 ${totalCodes} 个兑换口令` : `Total ${totalCodes} codes` }}
          </div>
          <div class="pagination-btns">
            <button class="btn btn-sm btn-outline" :disabled="page <= 1" @click="page--; loadCodes()">
              ← {{ locale === 'zh' ? '上一页' : 'Prev' }}
            </button>
            <span class="page-badge">{{ page }}</span>
            <button class="btn btn-sm btn-outline" :disabled="page * pageSize >= totalCodes" @click="page++; loadCodes()">
              {{ locale === 'zh' ? '下一页' : 'Next' }} →
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- TAB 2: REDEEM LOGS -->
    <div v-show="currentTab === 'logs'" class="tab-content">
      <div class="card filter-bar">
        <div class="filter-inputs">
          <div class="filter-item">
            <label class="filter-label">{{ locale === 'zh' ? '搜索使用口令' : 'Search Code' }}</label>
            <input
              v-model="logSearchCode"
              type="text"
              class="form-input form-input-sm"
              :placeholder="locale === 'zh' ? '输入口令过滤' : 'Filter by code'"
              @keyup.enter="loadLogs(true)"
            />
          </div>
          <div class="filter-item">
            <label class="filter-label">{{ locale === 'zh' ? '搜索兑换用户邮箱' : 'User Email' }}</label>
            <input
              v-model="logSearchEmail"
              type="text"
              class="form-input form-input-sm"
              :placeholder="locale === 'zh' ? '输入用户邮箱' : 'Filter by email'"
              @keyup.enter="loadLogs(true)"
            />
          </div>
        </div>
        <div class="filter-actions">
          <button class="btn btn-primary btn-sm" @click="loadLogs(true)">
            🔍 {{ locale === 'zh' ? '查询' : 'Search' }}
          </button>
          <button class="btn btn-ghost btn-sm" @click="resetLogFilter">
            ✕ {{ locale === 'zh' ? '重置' : 'Reset' }}
          </button>
        </div>
      </div>

      <div class="card table-card">
        <div class="table-wrap">
          <table class="data-table">
            <thead>
              <tr>
                <th>{{ locale === 'zh' ? '兑换时间' : 'Time' }}</th>
                <th>{{ locale === 'zh' ? '所用口令' : 'Code' }}</th>
                <th>{{ locale === 'zh' ? '兑换用户' : 'User' }}</th>
                <th>{{ locale === 'zh' ? '用途类型' : 'Purpose' }}</th>
                <th>{{ locale === 'zh' ? '所获商品 / 会员特权' : 'Granted Asset' }}</th>
                <th>{{ locale === 'zh' ? '兑换 IP' : 'IP' }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="loadingLogs">
                <td colspan="6" class="text-center py-8">
                  <span class="spinner"></span>
                  <span class="ml-2 text-secondary">{{ locale === 'zh' ? '正在加载兑换记录...' : 'Loading logs...' }}</span>
                </td>
              </tr>
              <tr v-else-if="logs.length === 0">
                <td colspan="6" class="text-center py-8 text-secondary">
                  {{ locale === 'zh' ? '暂无兑换使用记录' : 'No redemption logs found' }}
                </td>
              </tr>
              <tr v-else v-for="l in logs" :key="l.id">
                <td class="text-secondary text-sm">{{ formatDate(l.created_at) }}</td>
                <td>
                  <span class="code-text font-mono font-bold">{{ l.code }}</span>
                </td>
                <td>
                  <div class="user-row">
                    <span>👤</span>
                    <strong class="text-sm">{{ l.user_email }}</strong>
                  </div>
                </td>
                <td>
                  <span class="badge" :class="getPurposeBadgeClass(l.purpose)">
                    {{ formatPurpose(l.purpose) }}
                  </span>
                </td>
                <td>
                  <span v-if="l.granted_tier" class="vip-tier-tag font-bold">
                    👑 {{ formatTierName(l.granted_tier) }}
                  </span>
                  <span v-else-if="l.product_title" class="font-semibold text-sm">
                    🖼️ {{ l.product_title }}
                  </span>
                  <span v-else class="text-secondary">-</span>
                </td>
                <td class="font-mono text-secondary text-xs">🌐 {{ l.ip || '-' }}</td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination for Logs -->
        <div v-if="totalLogs > logPageSize" class="table-pagination">
          <div class="pagination-info">
            {{ locale === 'zh' ? `共 ${totalLogs} 条兑换记录` : `Total ${totalLogs} records` }}
          </div>
          <div class="pagination-btns">
            <button class="btn btn-sm btn-outline" :disabled="logPage <= 1" @click="logPage--; loadLogs()">
              ← {{ locale === 'zh' ? '上一页' : 'Prev' }}
            </button>
            <span class="page-badge">{{ logPage }}</span>
            <button class="btn btn-sm btn-outline" :disabled="logPage * logPageSize >= totalLogs" @click="logPage++; loadLogs()">
              {{ locale === 'zh' ? '下一页' : 'Next' }} →
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Code Modal -->
    <div v-if="showCreateModal" class="modal-overlay" @click.self="showCreateModal = false">
      <div class="modal-box redeem-modal-box">
        <!-- Header with Close Icon -->
        <div class="modal-header-row">
          <div class="modal-header-left">
            <div class="modal-header-badge">🎫</div>
            <div>
              <h3 class="modal-title">{{ locale === 'zh' ? '生成兑换口令 / 卡密' : 'Generate Redeem Codes' }}</h3>
              <p class="modal-sub">{{ locale === 'zh' ? '可灵活配置口令用途（壁纸限额免费/会员开通/特定壁纸）、兑换次数与有效期' : 'Configure purpose, redemption limits and expiration' }}</p>
            </div>
          </div>
          <button class="modal-close-btn" @click="showCreateModal = false" title="关闭">✕</button>
        </div>

        <div class="modal-body-scroll">
          <!-- Section 1: Generation Mode & Code -->
          <div class="form-section-card">
            <div class="section-title-row">
              <span class="section-icon">🎯</span>
              <span class="section-title">{{ locale === 'zh' ? '生成模式与口令设定' : 'Generation Mode & Code' }}</span>
            </div>

            <div class="mode-cards-grid">
              <label class="mode-card" :class="{ selected: createForm.mode === 'single' }">
                <input v-model="createForm.mode" type="radio" value="single" />
                <div class="mode-card-info">
                  <div class="mode-card-title">
                    <span>✨</span>
                    <strong>{{ locale === 'zh' ? '单个自定义口令' : 'Single Custom Code' }}</strong>
                  </div>
                  <small class="mode-card-desc">{{ locale === 'zh' ? '适合自定义特征口令，如 FREE5、VIP888 等' : 'Custom recognizable codes' }}</small>
                </div>
              </label>

              <label class="mode-card" :class="{ selected: createForm.mode === 'batch' }">
                <input v-model="createForm.mode" type="radio" value="batch" />
                <div class="mode-card-info">
                  <div class="mode-card-title">
                    <span>📦</span>
                    <strong>{{ locale === 'zh' ? '批量随机卡密生成' : 'Batch Random Generation' }}</strong>
                  </div>
                  <small class="mode-card-desc">{{ locale === 'zh' ? '适合批量生成导出，支持自定义前缀' : 'Generate up to 200 codes at once' }}</small>
                </div>
              </label>
            </div>

            <!-- Single custom code input -->
            <div v-if="createForm.mode === 'single'" class="form-group mt-3">
              <div class="label-with-action">
                <label class="form-label font-bold">{{ locale === 'zh' ? '口令代码 (英数字组合)' : 'Code Text' }} *</label>
                <button type="button" class="btn btn-link btn-xs" @click="generateRandomSingleCode">
                  🎲 {{ locale === 'zh' ? '随机生成代码' : 'Generate Random' }}
                </button>
              </div>
              <input
                v-model="createForm.code"
                type="text"
                class="form-input font-mono font-bold text-lg"
                placeholder="例如: FREE5, VIP888, TIKTOK2026"
              />
            </div>

            <!-- Batch Count & Prefix -->
            <div v-else class="form-grid-2 mt-3">
              <div class="form-group">
                <label class="form-label font-bold">{{ locale === 'zh' ? '批量生成数量 (1-200)' : 'Quantity' }} *</label>
                <input v-model.number="createForm.count" type="number" min="1" max="200" class="form-input font-bold" />
              </div>
              <div class="form-group">
                <label class="form-label font-bold">{{ locale === 'zh' ? '卡密前缀 (可选)' : 'Prefix (Optional)' }}</label>
                <input v-model="createForm.prefix" type="text" class="form-input font-mono" placeholder="如 CARD, VIP" />
              </div>
            </div>
          </div>

          <!-- Section 2: Purpose & Rewards -->
          <div class="form-section-card">
            <div class="section-title-row">
              <span class="section-icon">🎁</span>
              <span class="section-title">{{ locale === 'zh' ? '兑换权益与用途设置' : 'Reward & Purpose Settings' }}</span>
            </div>

            <div class="form-group">
              <label class="form-label font-bold">{{ locale === 'zh' ? '口令用途类型' : 'Purpose' }} *</label>
              <select v-model="createForm.purpose" class="form-input form-select font-semibold">
                <option value="price_limit">🏷️ {{ locale === 'zh' ? '限额免费兑换 (例如 5 元及以下壁纸免费)' : 'Price-Limit Free Unlock' }}</option>
                <option value="membership">👑 {{ locale === 'zh' ? '开通 VIP 会员 (月度/季度/年度/终生)' : 'Grant VIP Membership' }}</option>
                <option value="product">🖼️ {{ locale === 'zh' ? '兑换指定壁纸 (指定某一款作品)' : 'Unlock Specific Wallpaper' }}</option>
              </select>
            </div>

            <!-- Condition: Price Limit -->
            <div v-if="createForm.purpose === 'price_limit'" class="form-grid-2">
              <div class="form-group">
                <label class="form-label font-bold">{{ locale === 'zh' ? '最高支持原价金额' : 'Max Price' }} *</label>
                <input
                  v-model.number="createForm.max_price"
                  type="number"
                  step="0.01"
                  min="0.1"
                  class="form-input font-bold"
                  placeholder="5.00"
                />
              </div>
              <div class="form-group">
                <label class="form-label font-bold">{{ locale === 'zh' ? '货币币种' : 'Currency' }}</label>
                <select v-model="createForm.currency" class="form-input form-select">
                  <option value="CNY">¥ 人民币 (CNY)</option>
                  <option value="USD">$ 美元 (USD)</option>
                </select>
              </div>
            </div>

            <!-- Condition: Membership Tier -->
            <div v-if="createForm.purpose === 'membership'" class="form-group">
              <label class="form-label font-bold">{{ locale === 'zh' ? '目标会员等级' : 'Target VIP Tier' }} *</label>
              <select v-model="createForm.target_tier" class="form-input form-select font-semibold">
                <option value="monthly">🌙 {{ locale === 'zh' ? '月度 VIP 会员 (30天全场免单)' : 'Monthly VIP (30 days)' }}</option>
                <option value="quarterly">⭐ {{ locale === 'zh' ? '季度 VIP 会员 (90天全场免单)' : 'Quarterly VIP (90 days)' }}</option>
                <option value="yearly">👑 {{ locale === 'zh' ? '年度 VIP 会员 (365天全场免单)' : 'Yearly VIP (365 days)' }}</option>
                <option value="lifetime">💎 {{ locale === 'zh' ? '终生尊享 VIP (永久全免买断)' : 'Lifetime VIP (Permanent)' }}</option>
              </select>
            </div>

            <!-- Condition: Specific Product ID -->
            <div v-if="createForm.purpose === 'product'" class="form-group">
              <label class="form-label font-bold">{{ locale === 'zh' ? '指定壁纸商品 ID' : 'Wallpaper Product ID' }} *</label>
              <input
                v-model="createForm.product_id"
                type="text"
                class="form-input font-mono"
                placeholder="在商品列表中复制商品 ID"
              />
            </div>
          </div>

          <!-- Section 3: Usage Limits & Expiration -->
          <div class="form-section-card">
            <div class="section-title-row">
              <span class="section-icon">⏱️</span>
              <span class="section-title">{{ locale === 'zh' ? '使用规则与时效限制' : 'Rules & Expiration' }}</span>
            </div>

            <div class="form-grid-2">
              <!-- Left: Max Uses -->
              <div class="form-group">
                <label class="form-label font-bold">{{ locale === 'zh' ? '最大兑换次数上限' : 'Max Uses' }} *</label>
                <div class="uses-input-row">
                  <input
                    v-model.number="createForm.max_uses"
                    type="number"
                    min="1"
                    class="form-input font-bold"
                    style="width: 120px;"
                  />
                  <div class="preset-pills">
                    <button type="button" class="btn btn-outline btn-xs" :class="{ 'btn-primary': createForm.max_uses === 1 }" @click="createForm.max_uses = 1">1 次 (一卡一密)</button>
                    <button type="button" class="btn btn-outline btn-xs" :class="{ 'btn-primary': createForm.max_uses === 100 }" @click="createForm.max_uses = 100">100 次</button>
                    <button type="button" class="btn btn-outline btn-xs" :class="{ 'btn-primary': createForm.max_uses === 999999 }" @click="createForm.max_uses = 999999">无限次通码</button>
                  </div>
                </div>
                <small class="form-tip">
                  {{ locale === 'zh' ? '设为 1 则单次消费即作废；设为无限次则可作为面向全体粉丝的公开口令。' : '1 for one-time code; large number for public campaign.' }}
                </small>
              </div>

              <!-- Right: Expiry -->
              <div class="form-group">
                <label class="form-label font-bold">{{ locale === 'zh' ? '有效期时效' : 'Expiration' }}</label>
                <div class="expiry-toggle-row">
                  <label class="switch-control">
                    <input v-model="createForm.is_permanent" type="checkbox" />
                    <span class="switch-slider"></span>
                  </label>
                  <span class="font-bold text-sm">
                    {{ createForm.is_permanent ? '♾️ ' + (locale === 'zh' ? '永久有效 (永不过期)' : 'Permanent') : (locale === 'zh' ? '指定截止时间' : 'Specific Expiry') }}
                  </span>
                </div>
                <input
                  v-if="!createForm.is_permanent"
                  v-model="createForm.expires_at_str"
                  type="datetime-local"
                  class="form-input mt-2"
                />
              </div>
            </div>
          </div>

          <!-- Section 4: Remark (Optional) -->
          <div class="form-group">
            <label class="form-label">{{ locale === 'zh' ? '备注说明 (可选，便于对账管理)' : 'Remark (Optional)' }}</label>
            <input
              v-model="createForm.remark"
              type="text"
              class="form-input"
              :placeholder="locale === 'zh' ? '例如：小红书新粉丝福利、抖音直播间免单券、合作达人定制' : 'Optional note for accounting'"
            />
          </div>
        </div>

        <div class="modal-actions-bar">
          <button class="btn btn-ghost" @click="showCreateModal = false">{{ t('common.cancel') }}</button>
          <button class="btn btn-primary btn-submit-code" :disabled="creating" @click="submitCreateCode">
            <span v-if="creating" class="spinner"></span>
            {{ locale === 'zh' ? '✨ 立即生成兑换口令' : '✨ Generate Codes' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { redeemAPI } from '@/api'
import { useToast } from '@/composables/useToast'

const { t, locale } = useI18n()
const toast = useToast()

const currentTab = ref('codes')

// Tab 1: Codes
const codes = ref([])
const totalCodes = ref(0)
const loadingCodes = ref(false)
const search = ref('')
const filterPurpose = ref('')
const page = ref(1)
const pageSize = 15

// Tab 2: Logs
const logs = ref([])
const totalLogs = ref(0)
const loadingLogs = ref(false)
const logSearchCode = ref('')
const logSearchEmail = ref('')
const logPage = ref(1)
const logPageSize = 20

// Modal State
const showCreateModal = ref(false)
const creating = ref(false)
const createForm = ref({
  mode: 'single',
  code: '',
  count: 10,
  prefix: 'VIP',
  purpose: 'price_limit',
  target_tier: 'monthly',
  max_price: 5.0,
  currency: 'CNY',
  product_id: '',
  max_uses: 1,
  is_permanent: true,
  expires_at_str: '',
  remark: ''
})

function switchTab(tab) {
  currentTab.value = tab
  if (tab === 'logs' && logs.value.length === 0) {
    loadLogs(true)
  }
}

async function loadCodes(reset = false) {
  if (reset) page.value = 1
  loadingCodes.value = true
  try {
    const res = await redeemAPI.list({
      page: page.value,
      page_size: pageSize,
      purpose: filterPurpose.value,
      search: search.value
    })
    codes.value = res.codes || []
    totalCodes.value = res.total || 0
  } catch (err) {
    toast.error(err.message || '加载口令失败')
  } finally {
    loadingCodes.value = false
  }
}

async function loadLogs(reset = false) {
  if (reset) logPage.value = 1
  loadingLogs.value = true
  try {
    const res = await redeemAPI.logs({
      page: logPage.value,
      page_size: logPageSize,
      code: logSearchCode.value,
      email: logSearchEmail.value
    })
    logs.value = res.logs || []
    totalLogs.value = res.total || 0
  } catch (err) {
    toast.error(err.message || '加载使用流水失败')
  } finally {
    loadingLogs.value = false
  }
}

function resetCodeFilter() {
  search.value = ''
  filterPurpose.value = ''
  loadCodes(true)
}

function resetLogFilter() {
  logSearchCode.value = ''
  logSearchEmail.value = ''
  loadLogs(true)
}

function openCreateModal() {
  const d = new Date()
  d.setDate(d.getDate() + 30)
  const pad = n => String(n).padStart(2, '0')
  const defaultDateStr = `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())}T23:59`

  createForm.value = {
    mode: 'single',
    code: '',
    count: 10,
    prefix: 'CARD',
    purpose: 'price_limit',
    target_tier: 'monthly',
    max_price: 5.0,
    currency: 'CNY',
    product_id: '',
    max_uses: 1,
    is_permanent: true,
    expires_at_str: defaultDateStr,
    remark: ''
  }
  showCreateModal.value = true
}

function generateRandomSingleCode() {
  const chars = 'ABCDEFGHJKLMNPQRSTUVWXYZ23456789'
  let res = 'FREE-'
  for (let i = 0; i < 6; i++) {
    res += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  createForm.value.code = res
}

async function submitCreateCode() {
  if (createForm.value.mode === 'single' && !createForm.value.code.trim()) {
    toast.error(locale.value === 'zh' ? '请输入口令代码' : 'Code required')
    return
  }

  creating.value = true
  try {
    let expiresAt = null
    if (!createForm.value.is_permanent && createForm.value.expires_at_str) {
      expiresAt = new Date(createForm.value.expires_at_str).toISOString()
    }

    const payload = {
      count: createForm.value.mode === 'batch' ? createForm.value.count : 1,
      code: createForm.value.code.trim(),
      prefix: createForm.value.prefix.trim(),
      purpose: createForm.value.purpose,
      target_tier: createForm.value.target_tier,
      max_price: createForm.value.max_price,
      currency: createForm.value.currency,
      product_id: createForm.value.product_id.trim(),
      max_uses: createForm.value.max_uses,
      expires_at: expiresAt,
      is_permanent: createForm.value.is_permanent,
      remark: createForm.value.remark.trim()
    }

    const res = await redeemAPI.create(payload)
    toast.success(res.message || '生成成功！')
    showCreateModal.value = false
    loadCodes(true)
  } catch (err) {
    toast.error(err.message || '生成口令失败')
  } finally {
    creating.value = false
  }
}

async function toggleCodeActive(codeItem) {
  try {
    const res = await redeemAPI.toggle(codeItem.id)
    codeItem.is_active = res.is_active
    toast.success(res.message || '状态已更新')
  } catch (err) {
    toast.error(err.message || '更新失败')
  }
}

async function deleteCode(codeItem) {
  if (!confirm(locale.value === 'zh' ? `确定删除口令【${codeItem.code}】吗？` : `Delete code ${codeItem.code}?`)) {
    return
  }
  try {
    await redeemAPI.delete(codeItem.id)
    toast.success(locale.value === 'zh' ? '口令已删除' : 'Code deleted')
    loadCodes()
  } catch (err) {
    toast.error(err.message || '删除失败')
  }
}

function copyCode(text) {
  navigator.clipboard.writeText(text).then(() => {
    toast.success(locale.value === 'zh' ? '口令已复制到剪贴板！' : 'Copied to clipboard!')
  }).catch(() => {
    toast.success(locale.value === 'zh' ? `口令：${text}` : text)
  })
}

function formatPurpose(p) {
  switch (p) {
    case 'price_limit': return locale.value === 'zh' ? '限额免单' : 'Price-Limit'
    case 'membership': return locale.value === 'zh' ? '兑换会员' : 'Membership'
    case 'product': return locale.value === 'zh' ? '指定壁纸' : 'Specific Art'
    default: return p
  }
}

function getPurposeBadgeClass(p) {
  switch (p) {
    case 'price_limit': return 'badge-success-subtle'
    case 'membership': return 'badge-gold-subtle'
    case 'product': return 'badge-info-subtle'
    default: return 'badge-neutral'
  }
}

function formatTierName(tier) {
  switch (tier) {
    case 'monthly': return '月度 VIP (30天)'
    case 'quarterly': return '季度 VIP (90天)'
    case 'yearly': return '年度 VIP (365天)'
    case 'lifetime': return '终身尊享 VIP'
    default: return 'VIP 会员'
  }
}

function isExpired(expStr) {
  if (!expStr) return false
  return new Date(expStr).getTime() < Date.now()
}

function formatDate(isoStr) {
  if (!isoStr) return '-'
  const d = new Date(isoStr)
  if (isNaN(d.getTime())) return isoStr
  const pad = n => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

onMounted(() => {
  loadCodes(true)
})
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.page-title {
  margin: 0 0 4px;
  font-size: 1.5rem;
  font-weight: 800;
}

.page-sub {
  margin: 0;
  font-size: 0.88rem;
  color: var(--c-text-secondary);
  line-height: 1.45;
}

/* Nav Tabs */
.nav-tabs-bar {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
  border-bottom: 1px solid var(--c-border-light);
  padding-bottom: 8px;
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 18px;
  border-radius: 12px;
  background: transparent;
  border: none;
  font-size: 0.92rem;
  font-weight: 700;
  color: var(--c-text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
}

.tab-btn:hover {
  background: var(--c-bg);
  color: var(--c-text);
}

.tab-btn.active {
  background: var(--c-surface);
  color: var(--c-accent);
  box-shadow: 0 4px 12px rgba(0, 113, 227, 0.12);
}

.tab-badge {
  background: var(--c-border-light);
  color: var(--c-text);
  font-size: 0.75rem;
  padding: 2px 7px;
  border-radius: 999px;
}

.tab-btn.active .tab-badge {
  background: var(--c-accent);
  color: #ffffff;
}

/* Filter Bar */
.filter-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  margin-bottom: 18px;
  gap: 14px;
  flex-wrap: wrap;
}

.filter-inputs {
  display: flex;
  align-items: center;
  gap: 14px;
  flex-wrap: wrap;
}

.filter-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.filter-label {
  font-size: 0.74rem;
  font-weight: 700;
  color: var(--c-text-tertiary);
  text-transform: uppercase;
}

.filter-actions {
  display: flex;
  align-items: flex-end;
  gap: 8px;
}

/* Table */
.table-card {
  overflow: hidden;
  border: 1px solid var(--c-border-light);
}

.table-wrap {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.code-copy-cell {
  display: flex;
  align-items: center;
  gap: 6px;
}

.code-text {
  font-size: 0.92rem;
  letter-spacing: 0.05em;
  color: var(--c-text);
}

.copy-mini-btn {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 0.85rem;
  opacity: 0.6;
  transition: opacity 0.2s;
  padding: 2px;
}
.copy-mini-btn:hover { opacity: 1; }

.badge-success-subtle {
  background: rgba(16, 185, 129, 0.12);
  color: #059669;
  border: 1px solid rgba(16, 185, 129, 0.25);
  font-weight: 700;
  border-radius: 999px;
  padding: 3px 10px;
}

.badge-gold-subtle {
  background: rgba(245, 158, 11, 0.15);
  color: #d97706;
  border: 1px solid rgba(245, 158, 11, 0.35);
  font-weight: 700;
  border-radius: 999px;
  padding: 3px 10px;
}

.badge-info-subtle {
  background: rgba(59, 130, 246, 0.12);
  color: #2563eb;
  border: 1px solid rgba(59, 130, 246, 0.25);
  font-weight: 700;
  border-radius: 999px;
  padding: 3px 10px;
}

.rule-detail-cell {
  display: flex;
  flex-direction: column;
}

.text-emerald { color: #059669; }

.vip-tier-tag {
  color: #d97706;
  font-size: 0.85rem;
}

.uses-progress-wrap {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 100px;
}

.uses-count-row {
  font-size: 0.85rem;
}

.progress-bar-bg {
  width: 100%;
  height: 5px;
  background: var(--c-border-light);
  border-radius: 999px;
  overflow: hidden;
}

.progress-bar-fill {
  height: 100%;
  background: #3b82f6;
  border-radius: 999px;
  transition: width 0.3s ease;
}

.progress-bar-fill.fill-full {
  background: #ef4444;
}

.expiry-cell {
  display: flex;
  flex-direction: column;
  font-size: 0.82rem;
}

.text-danger-tag {
  font-weight: 700;
}

.user-row {
  display: flex;
  align-items: center;
  gap: 6px;
}

/* Modals */
.radio-cards-mini {
  display: flex;
  gap: 12px;
}

.radio-card-mini {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  border: 1.5px solid var(--c-border-light);
  border-radius: 12px;
  background: var(--c-surface);
  cursor: pointer;
  font-weight: 700;
  font-size: 0.88rem;
  transition: all 0.2s ease;
}

.radio-card-mini.selected {
  border-color: var(--c-accent);
  background: rgba(0, 113, 227, 0.05);
}

.form-grid-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

@media (max-width: 640px) {
  .form-grid-2 {
    grid-template-columns: 1fr;
  }
}

.btn-xs {
  font-size: 0.76rem;
  padding: 4px 8px;
  border-radius: 6px;
}

.btn-link {
  background: none;
  border: none;
  color: var(--c-accent);
  cursor: pointer;
  font-weight: 700;
}
.btn-link:hover { text-decoration: underline; }

/* ========================================================
   Spacious Modern Redeem Modal Styling
   ======================================================== */
.redeem-modal-box {
  max-width: 760px !important;
  width: 92vw;
  max-height: 88vh;
  display: flex;
  flex-direction: column;
  background: var(--c-surface);
  border-radius: 20px;
  box-shadow: 0 25px 60px -15px rgba(0, 0, 0, 0.35);
  overflow: hidden;
  border: 1px solid var(--c-border-light);
}

.modal-header-row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding: 22px 28px 18px;
  border-bottom: 1px solid var(--c-border-light);
  background: linear-gradient(135deg, rgba(0, 113, 227, 0.04) 0%, transparent 100%);
  flex-shrink: 0;
}

.modal-header-left {
  display: flex;
  align-items: center;
  gap: 14px;
}

.modal-header-badge {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  background: rgba(0, 113, 227, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  flex-shrink: 0;
}

.modal-title {
  font-size: 1.25rem;
  font-weight: 800;
  color: var(--c-text-primary);
  margin: 0 0 4px 0;
  letter-spacing: -0.01em;
}

.modal-sub {
  font-size: 0.85rem;
  color: var(--c-text-secondary);
  margin: 0;
  line-height: 1.4;
}

.modal-close-btn {
  background: none;
  border: none;
  font-size: 1.15rem;
  width: 34px;
  height: 34px;
  border-radius: 8px;
  cursor: pointer;
  color: var(--c-text-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.modal-close-btn:hover {
  background: rgba(0, 0, 0, 0.06);
  color: var(--c-text-primary);
}

.modal-body-scroll {
  flex: 1;
  overflow-y: auto;
  padding: 22px 28px;
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.form-section-card {
  background: var(--c-bg);
  border: 1px solid var(--c-border-light);
  border-radius: 14px;
  padding: 16px 18px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.section-title-row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding-bottom: 8px;
  border-bottom: 1px dashed var(--c-border-light);
}

.section-icon {
  font-size: 1.1rem;
}

.section-title {
  font-weight: 700;
  font-size: 0.95rem;
  color: var(--c-text-primary);
}

.mode-cards-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

@media (max-width: 600px) {
  .mode-cards-grid {
    grid-template-columns: 1fr;
  }
}

.mode-card {
  position: relative;
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 12px 14px;
  border: 1.5px solid var(--c-border-light);
  border-radius: 12px;
  background: var(--c-surface);
  cursor: pointer;
  transition: all 0.2s ease;
}

.mode-card:hover {
  border-color: rgba(0, 113, 227, 0.4);
}

.mode-card.selected {
  border-color: var(--c-accent);
  background: rgba(0, 113, 227, 0.05);
  box-shadow: 0 2px 10px rgba(0, 113, 227, 0.1);
}

.mode-card input[type="radio"] {
  accent-color: var(--c-accent);
  margin-top: 3px;
  cursor: pointer;
}

.mode-card-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.mode-card-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.92rem;
  color: var(--c-text-primary);
}

.mode-card-desc {
  font-size: 0.76rem;
  color: var(--c-text-secondary);
  line-height: 1.35;
}

.label-with-action {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 6px;
}

.uses-input-row {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.preset-pills {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.form-tip {
  font-size: 0.78rem;
  color: var(--c-text-secondary);
  margin-top: 6px;
  display: block;
  line-height: 1.4;
}

.expiry-toggle-row {
  display: flex;
  align-items: center;
  gap: 10px;
  min-height: 38px;
}

.switch-control {
  position: relative !important;
  display: inline-block !important;
  width: 44px !important;
  height: 24px !important;
  flex-shrink: 0 !important;
  cursor: pointer;
  user-select: none;
  vertical-align: middle;
}

.switch-control input {
  opacity: 0;
  width: 0;
  height: 0;
  position: absolute;
}

.switch-control .switch-slider {
  position: absolute !important;
  cursor: pointer;
  inset: 0 !important;
  background-color: #cbd5e1;
  transition: .3s;
  border-radius: 24px;
}

.switch-control .switch-slider:before {
  position: absolute !important;
  content: "";
  height: 18px !important;
  width: 18px !important;
  left: 3px !important;
  bottom: 3px !important;
  background-color: white;
  transition: .3s;
  border-radius: 50%;
}

.switch-control input:checked + .switch-slider {
  background-color: #10b981;
}

.switch-control input:checked + .switch-slider:before {
  transform: translateX(20px) !important;
}

.modal-actions-bar {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 28px;
  border-top: 1px solid var(--c-border-light);
  background: var(--c-bg);
  flex-shrink: 0;
}

.btn-submit-code {
  padding: 9px 22px;
  font-weight: 700;
  border-radius: 10px;
  font-size: 0.92rem;
}

.table-pagination {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 20px;
  border-top: 1px solid var(--c-border-light);
  background: var(--c-bg);
}

.pagination-info { font-size: 0.85rem; color: var(--c-text-secondary); }
.pagination-btns { display: flex; align-items: center; gap: 8px; }
.page-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  background: var(--c-primary);
  color: white;
  border-radius: 6px;
  font-weight: 700;
  font-size: 0.85rem;
}
</style>
