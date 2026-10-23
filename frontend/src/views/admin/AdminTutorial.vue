<template>
  <div class="admin-tutorial-page">
    <div class="page-header-row">
      <div>
        <h2 class="section-title">{{ locale === 'zh' ? '详情页通用教程模板' : 'Universal Tutorial Template' }}</h2>
        <p class="section-subtitle">
          {{ locale === 'zh' ? '在所有壁纸详情页（桌面端右侧、移动端底部）展示统一的设置指南。支持多版本 iOS 设置方法与视频直达超链，在此统一更新即可全站生效！' : 'Unified setup guide shown across all wallpaper detail pages with version cards and video links.' }}
        </p>
      </div>
      <div class="header-actions">
        <button class="btn btn-outline" @click="fetchTutorial" :disabled="loading">
          🔄 {{ locale === 'zh' ? '重置/刷新' : 'Refresh' }}
        </button>
        <button class="btn btn-primary" @click="saveTutorial" :disabled="saving">
          <span v-if="saving" class="spinner"></span>
          💾 {{ locale === 'zh' ? '保存全站模板' : 'Save Template' }}
        </button>
      </div>
    </div>

    <!-- Alert Banner -->
    <div v-if="saveMessage" class="save-success-banner">
      ✓ {{ saveMessage }}
    </div>

    <div v-if="loading" class="loading-state">
      <span class="spinner"></span>
      <p>{{ locale === 'zh' ? '正在加载教程模板...' : 'Loading template...' }}</p>
    </div>

    <div v-else class="tutorial-layout-grid">
      <!-- Left: Editor Column -->
      <div class="editor-col">
        <!-- Global Switch & Settings Card -->
        <div class="card settings-card">
          <div class="card-title-row">
            <span class="card-icon">⚙️</span>
            <h3>{{ locale === 'zh' ? '模块全局配置' : 'Global Configuration' }}</h3>
          </div>

          <div class="form-group toggle-row">
            <div>
              <label class="form-label font-bold">{{ locale === 'zh' ? '启用详情页通用教程模块' : 'Enable Tutorial Module' }}</label>
              <p class="form-desc">{{ locale === 'zh' ? '关闭后，前台壁纸详情页将隐藏整个教程区域' : 'When disabled, the tutorial section is hidden from all detail pages.' }}</p>
            </div>
            <label class="switch">
              <input v-model="form.enabled" type="checkbox" />
              <span class="slider round"></span>
            </label>
          </div>

          <div class="form-row-dual">
            <div class="form-group flex-1">
              <label class="form-label">{{ locale === 'zh' ? '模块主标题 (中文)' : 'Main Title (ZH)' }}</label>
              <input
                v-model="form.title_zh"
                type="text"
                class="form-input"
                placeholder="如：iOS 动态壁纸实机交互设置教程"
              />
            </div>
            <div class="form-group flex-1">
              <label class="form-label">{{ locale === 'zh' ? '模块主标题 (英文)' : 'Main Title (EN)' }}</label>
              <input
                v-model="form.title_en"
                type="text"
                class="form-input"
                placeholder="e.g. iOS Live Wallpaper Setup Guide"
              />
            </div>
          </div>

          <div class="form-group">
            <label class="form-label">{{ locale === 'zh' ? '顶部醒目提示条 (中文)' : 'Notice Banner (ZH)' }}</label>
            <input
              v-model="form.notice_zh"
              type="text"
              class="form-input"
              placeholder="💡 提示：依据当前 iPhone 系统版本选择对应指引与视频教程"
            />
          </div>

          <div class="form-group">
            <label class="form-label">{{ locale === 'zh' ? '顶部醒目提示条 (英文)' : 'Notice Banner (EN)' }}</label>
            <input
              v-model="form.notice_en"
              type="text"
              class="form-input"
              placeholder="💡 Tip: Choose the corresponding guide and video tutorial according to your iOS version."
            />
          </div>
        </div>

        <!-- Version Items List -->
        <div class="card items-card">
          <div class="card-title-row space-between">
            <div class="flex-row-center">
              <span class="card-icon">📱</span>
              <h3>{{ locale === 'zh' ? '各 iOS 系统版本教程列表' : 'Version Guides & Video Links' }}</h3>
              <span class="count-pill">{{ form.items.length }}</span>
            </div>
            <button type="button" class="btn btn-outline btn-sm" @click="addTutorialItem">
              ＋ {{ locale === 'zh' ? '添加版本教程' : 'Add Version Card' }}
            </button>
          </div>

          <div v-if="form.items.length === 0" class="empty-items-box">
            <p>{{ locale === 'zh' ? '尚未添加任何系统版本教程，请点击上方按钮添加' : 'No version items yet. Click Add above.' }}</p>
          </div>

          <!-- Item Card Loop -->
          <div v-for="(item, idx) in form.items" :key="item.id || idx" class="version-item-card">
            <div class="item-card-header">
              <div class="item-header-left">
                <span class="item-order-badge">#{{ idx + 1 }}</span>
                <span class="item-range-pill" :class="item.badge_color || 'blue'">
                  {{ item.version_range || (locale === 'zh' ? '未命名版本' : 'Untitled Version') }}
                </span>
                <span class="item-title-preview">{{ item.title_zh }}</span>
              </div>
              <div class="item-header-actions">
                <button
                  type="button"
                  class="btn-icon"
                  :disabled="idx === 0"
                  @click="moveItem(idx, -1)"
                  title="上移"
                >↑</button>
                <button
                  type="button"
                  class="btn-icon"
                  :disabled="idx === form.items.length - 1"
                  @click="moveItem(idx, 1)"
                  title="下移"
                >↓</button>
                <button
                  type="button"
                  class="btn-icon btn-danger"
                  @click="removeItem(idx)"
                  title="删除"
                >✕</button>
              </div>
            </div>

            <div class="item-card-body">
              <div class="form-row-dual">
                <div class="form-group flex-1">
                  <label class="form-label required">{{ locale === 'zh' ? '适用系统版本区间' : 'iOS Version Range' }}</label>
                  <input
                    v-model="item.version_range"
                    type="text"
                    class="form-input font-bold"
                    placeholder="如：iOS 18.0 - iOS 26.4、iOS 26.5 - iOS 27.0 Beta"
                    required
                  />
                </div>
                <div class="form-group flex-1">
                  <label class="form-label">{{ locale === 'zh' ? '徽标文字' : 'Badge Label' }}</label>
                  <input
                    v-model="item.badge_zh"
                    type="text"
                    class="form-input"
                    placeholder="如：主流推荐、Beta测试专用"
                  />
                </div>
                <div class="form-group" style="width: 130px;">
                  <label class="form-label">{{ locale === 'zh' ? '徽标主题色' : 'Badge Color' }}</label>
                  <select v-model="item.badge_color" class="form-input form-select">
                    <option value="blue">🔵 科技蓝</option>
                    <option value="purple">🟣 炫彩紫</option>
                    <option value="emerald">🟢 推荐绿</option>
                    <option value="amber">🟡 金钻金</option>
                  </select>
                </div>
              </div>

              <div class="form-row-dual">
                <div class="form-group flex-1">
                  <label class="form-label">{{ locale === 'zh' ? '教程副标题 (中文)' : 'Title (ZH)' }}</label>
                  <input
                    v-model="item.title_zh"
                    type="text"
                    class="form-input"
                    placeholder="如：经典锁屏与景深联动设置法"
                  />
                </div>
                <div class="form-group flex-1">
                  <label class="form-label">{{ locale === 'zh' ? '教程副标题 (英文)' : 'Title (EN)' }}</label>
                  <input
                    v-model="item.title_en"
                    type="text"
                    class="form-input"
                    placeholder="e.g. Classic Lock Screen Setup"
                  />
                </div>
              </div>

              <!-- Video Link -->
              <div class="video-config-row">
                <div class="form-group flex-1">
                  <label class="form-label font-bold text-accent">
                    ▶ {{ locale === 'zh' ? '视频教程超链接 (B站 / 抖音 / YouTube / MP4直链)' : 'Video Tutorial Link' }}
                  </label>
                  <div class="input-with-action">
                    <input
                      v-model="item.video_url"
                      type="url"
                      class="form-input"
                      placeholder="https://www.bilibili.com/video/BV..."
                    />
                    <a
                      v-if="item.video_url"
                      :href="item.video_url"
                      target="_blank"
                      rel="noopener noreferrer"
                      class="btn btn-outline btn-sm test-link-btn"
                    >
                      ↗ {{ locale === 'zh' ? '测试链接' : 'Open Link' }}
                    </a>
                  </div>
                </div>
                <div class="form-group" style="width: 140px;">
                  <label class="form-label">{{ locale === 'zh' ? '视频平台' : 'Platform' }}</label>
                  <select v-model="item.video_platform" class="form-input form-select">
                    <option value="bilibili">📺 哔哩哔哩 (Bilibili)</option>
                    <option value="douyin">🎵 抖音 / TikTok</option>
                    <option value="youtube">▶️ YouTube</option>
                    <option value="direct">🎬 MP4 视频直链</option>
                  </select>
                </div>
              </div>

              <!-- Steps List Builder -->
              <div class="steps-builder-section">
                <div class="steps-builder-header">
                  <label class="form-label">{{ locale === 'zh' ? '操作步骤指引' : 'Step-by-step Instructions' }}</label>
                  <button type="button" class="btn btn-ghost btn-xs" @click="addStepToItem(item)">
                    ＋ {{ locale === 'zh' ? '增加步骤' : 'Add Step' }}
                  </button>
                </div>

                <div v-for="(step, sIdx) in item.steps" :key="sIdx" class="step-input-row">
                  <span class="step-num-pill">{{ sIdx + 1 }}</span>
                  <input
                    v-model="item.steps[sIdx]"
                    type="text"
                    class="form-input flex-1"
                    :placeholder="`步骤 ${sIdx + 1} 说明...`"
                  />
                  <button
                    type="button"
                    class="btn-icon-xs"
                    @click="removeStepFromItem(item, sIdx)"
                    title="删除此步骤"
                  >✕</button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Right: Real-time Live Preview Column -->
      <div class="preview-col">
        <div class="card preview-card sticky-card">
          <div class="card-title-row">
            <span class="card-icon">👁️</span>
            <h3>{{ locale === 'zh' ? '详情页实时效果预览' : 'Live Detail Preview' }}</h3>
            <span class="preview-tag">{{ locale === 'zh' ? '移动端底部 / Web端右侧' : 'Responsive' }}</span>
          </div>

          <!-- Mock Preview Container -->
          <div class="mock-tutorial-box" :class="{ disabled: !form.enabled }">
            <div class="mock-header">
              <div class="mock-badge-tag">📖 {{ locale === 'zh' ? '实机设置教程' : 'Setup Tutorial' }}</div>
              <h4 class="mock-title">{{ form.title_zh || 'iOS 动态壁纸实机交互设置教程' }}</h4>
            </div>

            <div v-if="form.notice_zh" class="mock-notice">
              {{ form.notice_zh }}
            </div>

            <!-- Version Switch Pills -->
            <div class="mock-version-pills">
              <button
                v-for="(it, idx) in form.items"
                :key="idx"
                class="mock-ver-pill"
                :class="{ active: previewActiveIndex === idx, [it.badge_color || 'blue']: true }"
                @click="previewActiveIndex = idx"
              >
                <span>{{ it.version_range }}</span>
                <span v-if="it.badge_zh" class="pill-badge">{{ it.badge_zh }}</span>
              </button>
            </div>

            <!-- Active Version Details -->
            <div v-if="activePreviewItem" class="mock-active-guide">
              <div class="mock-guide-head">
                <span class="guide-title">{{ activePreviewItem.title_zh }}</span>
                <a
                  v-if="activePreviewItem.video_url"
                  :href="activePreviewItem.video_url"
                  target="_blank"
                  class="mock-video-jump-btn"
                >
                  ▶ {{ locale === 'zh' ? '观看视频教程' : 'Watch Video' }} ↗
                </a>
              </div>

              <div class="mock-steps-list">
                <div v-for="(st, stIdx) in activePreviewItem.steps" :key="stIdx" class="mock-step-row">
                  <span class="mock-step-num">{{ stIdx + 1 }}</span>
                  <span class="mock-step-text">{{ st }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { tutorialAPI } from '@/api'

const locale = ref(localStorage.getItem('lang') || 'zh')
const loading = ref(false)
const saving = ref(false)
const saveMessage = ref('')
const previewActiveIndex = ref(0)

const form = ref({
  enabled: true,
  title_zh: 'iOS 动态壁纸实机交互设置教程',
  title_en: 'iOS Live Wallpaper Setup Guide',
  notice_zh: '💡 提示：依据当前 iPhone 系统版本选择对应指引与视频教程，后续随 iOS 系统大版本统一在此更新维护',
  notice_en: '💡 Tip: Follow the instructions and video link according to your iOS version.',
  items: []
})

const activePreviewItem = computed(() => {
  if (form.value.items.length === 0) return null
  return form.value.items[previewActiveIndex.value] || form.value.items[0]
})

async function fetchTutorial() {
  loading.value = true
  try {
    const res = await tutorialAPI.adminGet()
    if (res) {
      form.value = {
        enabled: res.enabled !== false,
        title_zh: res.title_zh || '',
        title_en: res.title_en || '',
        notice_zh: res.notice_zh || '',
        notice_en: res.notice_en || '',
        items: res.items || []
      }
    }
  } catch (err) {
    console.error('Failed to load tutorial template:', err)
  } finally {
    loading.value = false
  }
}

async function saveTutorial() {
  saving.value = true
  saveMessage.value = ''
  try {
    await tutorialAPI.adminUpdate(form.value)
    saveMessage.value = locale.value === 'zh'
      ? '通用教程模板已保存！所有壁纸详情页已同步生效。'
      : 'Universal tutorial template saved successfully!'
    setTimeout(() => {
      saveMessage.value = ''
    }, 4000)
  } catch (err) {
    alert(err.message || '保存失败，请重试')
  } finally {
    saving.value = false
  }
}

function addTutorialItem() {
  form.value.items.push({
    id: 'tut_' + Date.now().toString(36),
    version_range: 'iOS 27.1 - 28.0',
    title_zh: '新系统景深与动作联动设置法',
    title_en: 'New System Setup',
    badge_zh: '最新测试版',
    badge_en: 'Latest',
    badge_color: 'purple',
    video_url: 'https://www.bilibili.com',
    video_platform: 'bilibili',
    steps: [
      '下载原画及 .tendies 动画文件',
      '在手机相册将照片设置为锁屏',
      '点击景深并绑定动态动作'
    ]
  })
  previewActiveIndex.value = form.value.items.length - 1
}

function removeItem(idx) {
  if (confirm(locale.value === 'zh' ? '确定删除该版本教程卡片吗？' : 'Delete this version card?')) {
    form.value.items.splice(idx, 1)
    if (previewActiveIndex.value >= form.value.items.length) {
      previewActiveIndex.value = Math.max(0, form.value.items.length - 1)
    }
  }
}

function moveItem(idx, dir) {
  const target = idx + dir
  if (target < 0 || target >= form.value.items.length) return
  const temp = form.value.items[idx]
  form.value.items[idx] = form.value.items[target]
  form.value.items[target] = temp
  previewActiveIndex.value = target
}

function addStepToItem(item) {
  if (!item.steps) item.steps = []
  item.steps.push('')
}

function removeStepFromItem(item, sIdx) {
  item.steps.splice(sIdx, 1)
}

onMounted(() => {
  fetchTutorial()
})
</script>

<style scoped>
.admin-tutorial-page {
  padding-bottom: 60px;
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

.save-success-banner {
  padding: 12px 18px;
  border-radius: 10px;
  background: rgba(16, 185, 129, 0.15);
  border: 1px solid rgba(16, 185, 129, 0.4);
  color: #10b981;
  font-weight: 600;
  font-size: 0.9rem;
  margin-bottom: 20px;
}

.loading-state {
  padding: 80px 20px;
  text-align: center;
  color: var(--c-text-secondary);
}

.tutorial-layout-grid {
  display: grid;
  grid-template-columns: 1fr 440px;
  gap: 24px;
  align-items: start;
}

@media (max-width: 1100px) {
  .tutorial-layout-grid {
    grid-template-columns: 1fr;
  }
}

.card {
  background: var(--c-card-bg, #1a1b24);
  border: 1px solid var(--c-border, rgba(255, 255, 255, 0.08));
  border-radius: 14px;
  padding: 24px;
  margin-bottom: 24px;
}

.card-title-row {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 20px;
}

.card-title-row.space-between {
  justify-content: space-between;
}

.flex-row-center {
  display: flex;
  align-items: center;
  gap: 10px;
}

.card-icon {
  font-size: 1.25rem;
}

.count-pill {
  padding: 2px 8px;
  border-radius: 999px;
  background: rgba(56, 189, 248, 0.15);
  color: #38bdf8;
  font-size: 0.8rem;
  font-weight: 700;
}

.toggle-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 18px;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.03);
  margin-bottom: 20px;
}

.form-desc {
  font-size: 0.8rem;
  color: var(--c-text-secondary);
  margin-top: 2px;
}

/* Switch */
.switch {
  position: relative;
  display: inline-block;
  width: 48px;
  height: 26px;
  flex-shrink: 0;
}

.switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(255, 255, 255, 0.2);
  transition: 0.3s;
}

.slider:before {
  position: absolute;
  content: "";
  height: 20px;
  width: 20px;
  left: 3px;
  bottom: 3px;
  background-color: white;
  transition: 0.3s;
}

input:checked + .slider {
  background-color: #38bdf8;
}

input:checked + .slider:before {
  transform: translateX(22px);
}

.slider.round {
  border-radius: 26px;
}

.slider.round:before {
  border-radius: 50%;
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

.text-accent {
  color: #38bdf8;
}

/* Version Item Card */
.version-item-card {
  border-radius: 12px;
  border: 1px solid var(--c-border, rgba(255, 255, 255, 0.08));
  background: rgba(255, 255, 255, 0.02);
  margin-bottom: 16px;
  overflow: hidden;
}

.item-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 18px;
  background: rgba(255, 255, 255, 0.03);
  border-bottom: 1px solid var(--c-border, rgba(255, 255, 255, 0.06));
}

.item-header-left {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
  overflow: hidden;
}

.item-order-badge {
  font-size: 0.78rem;
  font-weight: 700;
  color: var(--c-text-secondary);
}

.item-range-pill {
  padding: 3px 10px;
  border-radius: 999px;
  font-size: 0.82rem;
  font-weight: 700;
}

.item-range-pill.blue {
  background: rgba(56, 189, 248, 0.15);
  color: #38bdf8;
}

.item-range-pill.purple {
  background: rgba(168, 85, 247, 0.15);
  color: #c084fc;
}

.item-range-pill.emerald {
  background: rgba(16, 185, 129, 0.15);
  color: #34d399;
}

.item-range-pill.amber {
  background: rgba(245, 158, 11, 0.15);
  color: #fbbf24;
}

.item-title-preview {
  font-size: 0.85rem;
  color: var(--c-text-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.item-header-actions {
  display: flex;
  align-items: center;
  gap: 6px;
}

.btn-icon {
  width: 28px;
  height: 28px;
  border-radius: 6px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
  color: var(--c-text);
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 0.85rem;
}

.btn-icon:hover {
  background: rgba(255, 255, 255, 0.1);
}

.btn-icon.btn-danger:hover {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
}

.item-card-body {
  padding: 18px;
}

.video-config-row {
  display: flex;
  gap: 14px;
  margin-bottom: 16px;
  padding: 14px;
  border-radius: 10px;
  background: rgba(56, 189, 248, 0.04);
  border: 1px solid rgba(56, 189, 248, 0.15);
}

.input-with-action {
  display: flex;
  gap: 8px;
}

.test-link-btn {
  white-space: nowrap;
  flex-shrink: 0;
}

.steps-builder-section {
  padding: 14px;
  border-radius: 10px;
  background: rgba(0, 0, 0, 0.2);
}

.steps-builder-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

.step-input-row {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
}

.step-num-pill {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 0.72rem;
  font-weight: 700;
  color: var(--c-text-secondary);
  flex-shrink: 0;
}

.btn-icon-xs {
  width: 24px;
  height: 24px;
  border-radius: 4px;
  background: none;
  border: none;
  color: var(--c-text-secondary);
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.btn-icon-xs:hover {
  color: #ef4444;
}

/* Preview Card */
.sticky-card {
  position: sticky;
  top: 24px;
}

.preview-tag {
  font-size: 0.75rem;
  padding: 2px 8px;
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.06);
  color: var(--c-text-secondary);
}

.mock-tutorial-box {
  border-radius: 16px;
  background: linear-gradient(145deg, rgba(30, 41, 59, 0.7), rgba(15, 23, 42, 0.8));
  border: 1px solid rgba(255, 255, 255, 0.12);
  padding: 20px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.4);
}

.mock-tutorial-box.disabled {
  opacity: 0.4;
  filter: grayscale(0.8);
}

.mock-header {
  margin-bottom: 12px;
}

.mock-badge-tag {
  display: inline-block;
  font-size: 0.72rem;
  font-weight: 700;
  color: #38bdf8;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 4px;
}

.mock-title {
  font-size: 1.05rem;
  font-weight: 800;
  color: #ffffff;
}

.mock-notice {
  font-size: 0.8rem;
  line-height: 1.4;
  padding: 8px 12px;
  border-radius: 8px;
  background: rgba(245, 158, 11, 0.12);
  border: 1px solid rgba(245, 158, 11, 0.3);
  color: #fbbf24;
  margin-bottom: 16px;
}

.mock-version-pills {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 16px;
}

.mock-ver-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  border-radius: 999px;
  font-size: 0.78rem;
  font-weight: 700;
  cursor: pointer;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #94a3b8;
  transition: all 0.2s ease;
}

.mock-ver-pill.active {
  background: rgba(56, 189, 248, 0.2);
  border-color: rgba(56, 189, 248, 0.6);
  color: #38bdf8;
  box-shadow: 0 0 12px rgba(56, 189, 248, 0.3);
}

.mock-ver-pill.active.purple {
  background: rgba(168, 85, 247, 0.2);
  border-color: rgba(168, 85, 247, 0.6);
  color: #c084fc;
  box-shadow: 0 0 12px rgba(168, 85, 247, 0.3);
}

.mock-ver-pill.active.emerald {
  background: rgba(16, 185, 129, 0.2);
  border-color: rgba(16, 185, 129, 0.6);
  color: #34d399;
  box-shadow: 0 0 12px rgba(16, 185, 129, 0.3);
}

.mock-ver-pill.active.amber {
  background: rgba(245, 158, 11, 0.2);
  border-color: rgba(245, 158, 11, 0.6);
  color: #fbbf24;
  box-shadow: 0 0 12px rgba(245, 158, 11, 0.3);
}

.pill-badge {
  font-size: 0.68rem;
  padding: 1px 6px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.2);
  color: #ffffff;
}

.mock-active-guide {
  background: rgba(0, 0, 0, 0.25);
  border-radius: 12px;
  padding: 16px;
  border: 1px solid rgba(255, 255, 255, 0.06);
}

.mock-guide-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  margin-bottom: 14px;
}

.guide-title {
  font-size: 0.9rem;
  font-weight: 700;
  color: #f1f5f9;
}

.mock-video-jump-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 6px 14px;
  border-radius: 999px;
  font-size: 0.75rem;
  font-weight: 700;
  color: #ffffff;
  background: linear-gradient(135deg, #0284c7, #2563eb);
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.4);
  text-decoration: none;
}

.mock-steps-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.mock-step-row {
  display: flex;
  align-items: flex-start;
  gap: 10px;
}

.mock-step-num {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: rgba(56, 189, 248, 0.15);
  color: #38bdf8;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 0.72rem;
  font-weight: 800;
  flex-shrink: 0;
  margin-top: 2px;
}

.mock-step-text {
  font-size: 0.82rem;
  color: #cbd5e1;
  line-height: 1.45;
}
</style>
