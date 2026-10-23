<template>
  <div class="product-edit">
    <div class="edit-header">
      <router-link to="/admin/products" class="btn btn-ghost btn-sm">← {{ t('common.back') }}</router-link>
      <h2>{{ isNew ? t('admin.add_product') : t('admin.edit_product') }}</h2>
    </div>

    <form @submit.prevent="handleSave" class="edit-form">
      <div class="edit-grid">
        <!-- Left column: main fields -->
        <div class="edit-main">
          <!-- Title (zh/en) -->
          <div class="section-card">
            <h3 class="section-title">{{ locale === 'zh' ? '基本信息' : 'Basic Info' }}</h3>
            <div class="form-group">
              <label class="form-label">{{ locale === 'zh' ? '中文标题' : 'Chinese Title' }} *</label>
              <input v-model="form.title_zh" type="text" class="form-input" :class="{ error: errors.title_zh }" placeholder="请输入中文标题" />
              <span v-if="errors.title_zh" class="form-error">{{ errors.title_zh }}</span>
            </div>
            <div class="form-group">
              <label class="form-label">{{ locale === 'zh' ? '英文标题' : 'English Title' }}</label>
              <input v-model="form.title_en" type="text" class="form-input" placeholder="Enter English title" />
            </div>
            <div class="form-row">
              <div class="form-group">
                <label class="form-label">{{ locale === 'zh' ? '价格 (人民币 CNY - ¥)' : 'Price (CNY - ¥)' }}</label>
                <input v-model="form.price_cny" type="number" step="0.01" min="0" class="form-input" placeholder="0 为免费" @input="onCnyInput" />
              </div>
              <div class="form-group">
                <label class="form-label">{{ locale === 'zh' ? '价格 (美元 USD - $)' : 'Price (USD - $)' }}</label>
                <input v-model="form.price_usd" type="number" step="0.01" min="0" class="form-input" placeholder="0 为免费" @input="onUsdInput" />
              </div>
              <div class="form-group">
                <label class="form-label">{{ locale === 'zh' ? '🔑 提取码/解锁卡密 (选填)' : 'Unlock Verify Code' }}</label>
                <input v-model="form.verify_code" type="text" class="form-input" placeholder="如 666、loveyou" />
              </div>
              <div class="form-group">
                <label class="form-label">{{ t('admin.category') }}</label>
                <select v-model="form.category" class="form-input form-select">
                  <option value="">{{ locale === 'zh' ? '选择分类' : 'Select category' }}</option>
                  <option v-for="c in categories" :key="c.value" :value="c.value">{{ locale === 'zh' ? c.zh : c.en }}</option>
                </select>
              </div>
            </div>
            <div class="currency-hint-bar" style="margin: -6px 0 16px 0; font-size: 12px; color: #3b82f6; display: flex; align-items: center; gap: 6px;">
              <span>💱</span>
              <span>{{ locale === 'zh' ? `双轨汇率联动 (1 USD = ${exchangeRate} CNY)：只填一项时自动换算，双轨全填以具体数字为准` : `Dual currency linked (1 USD = ${exchangeRate} CNY): auto-converted if one is filled, exact numbers if both` }}</span>
            </div>
            <div class="form-group">
              <label class="form-label">{{ t('admin.tags') }}</label>
              <div class="tags-input">
                <span v-for="tag in form.tags" :key="tag" class="tag" style="cursor:default">
                  {{ tag }}
                  <button type="button" @click="removeTag(tag)" class="tag-remove">✕</button>
                </span>
                <input
                  v-model="tagInput"
                  type="text"
                  class="tag-text-input"
                  :placeholder="locale === 'zh' ? '输入标签按回车添加' : 'Type and press Enter'"
                  @keydown.enter.prevent="addTag"
                  @keydown.comma.prevent="addTag"
                />
              </div>
            </div>

            <!-- VIP Tier Download Privilege Configuration -->
            <div class="form-group vip-tier-group">
              <div class="vip-tier-header">
                <div class="vip-title-col">
                  <label class="form-label" style="margin-bottom: 2px;">
                    👑 {{ locale === 'zh' ? 'VIP 会员免单下载权限等级' : 'Allowed VIP Download Tiers' }}
                  </label>
                  <span class="vip-tier-hint">
                    {{ locale === 'zh' ? '默认全等级享受；取消勾选某些等级即可将其设为高等级专属（如仅限季度及以上）' : 'Default all VIP tiers. Uncheck tiers to restrict download rights.' }}
                  </span>
                </div>
                <!-- Quick Preset Filter Buttons -->
                <div class="vip-preset-btns">
                  <button type="button" class="preset-btn" :class="{ 'preset-active': isPresetActive('all') }" @click="setVIPPreset('all')">
                    {{ locale === 'zh' ? '全选 (默认)' : 'All' }}
                  </button>
                  <button type="button" class="preset-btn" :class="{ 'preset-active': isPresetActive('quarterly_up') }" @click="setVIPPreset('quarterly_up')">
                    {{ locale === 'zh' ? '季度及以上' : 'Quarterly+' }}
                  </button>
                  <button type="button" class="preset-btn" :class="{ 'preset-active': isPresetActive('yearly_up') }" @click="setVIPPreset('yearly_up')">
                    {{ locale === 'zh' ? '年度及以上' : 'Yearly+' }}
                  </button>
                  <button type="button" class="preset-btn" :class="{ 'preset-active': isPresetActive('lifetime_only') }" @click="setVIPPreset('lifetime_only')">
                    {{ locale === 'zh' ? '仅终生会员' : 'Lifetime' }}
                  </button>
                  <button type="button" class="preset-btn preset-danger" @click="setVIPPreset('none')">
                    {{ locale === 'zh' ? '全取消' : 'None' }}
                  </button>
                </div>
              </div>

              <!-- VIP Tier Checkbox Cards -->
              <div class="vip-tiers-grid">
                <div
                  v-for="plan in availableVIPPlans"
                  :key="plan.type"
                  class="vip-tier-card"
                  :class="{ active: isTierChecked(plan.type) }"
                  @click="toggleTier(plan.type)"
                >
                  <div class="tier-checkbox" :class="{ checked: isTierChecked(plan.type) }">
                    <span v-if="isTierChecked(plan.type)" class="tier-check-mark">✓</span>
                  </div>
                  <span class="tier-icon">{{ plan.icon || '👑' }}</span>
                  <div class="tier-content">
                    <div class="tier-name-row">
                      <span class="tier-name">{{ locale === 'zh' ? plan.name_zh : (plan.name_en || plan.name_zh) }}</span>
                      <span v-if="plan.badge" class="tier-badge">{{ plan.badge }}</span>
                    </div>
                    <span class="tier-price">${{ plan.price }} · {{ plan.duration_days || 30 }}天</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Rich Text Description -->
          <div class="section-card">
            <h3 class="section-title">{{ t('admin.product_desc') }} (TipTap 富文本)</h3>
            <div class="editor-wrapper">
              <!-- Toolbar -->
              <div class="editor-toolbar" v-if="editor">
                <button type="button" class="toolbar-btn" :class="{ active: editor.isActive('bold') }" @click="editor.chain().focus().toggleBold().run()" title="Bold">𝐁</button>
                <button type="button" class="toolbar-btn" :class="{ active: editor.isActive('italic') }" @click="editor.chain().focus().toggleItalic().run()" title="Italic">𝐼</button>
                <button type="button" class="toolbar-btn" :class="{ active: editor.isActive('strike') }" @click="editor.chain().focus().toggleStrike().run()" title="Strike">S̶</button>
                <div class="toolbar-divider"></div>
                <button type="button" class="toolbar-btn" :class="{ active: editor.isActive('heading', { level: 2 }) }" @click="editor.chain().focus().toggleHeading({ level: 2 }).run()" title="Heading 2">H2</button>
                <button type="button" class="toolbar-btn" :class="{ active: editor.isActive('heading', { level: 3 }) }" @click="editor.chain().focus().toggleHeading({ level: 3 }).run()" title="Heading 3">H3</button>
                <div class="toolbar-divider"></div>
                <button type="button" class="toolbar-btn" :class="{ active: editor.isActive('bulletList') }" @click="editor.chain().focus().toggleBulletList().run()" title="Bullet List">• List</button>
                <button type="button" class="toolbar-btn" :class="{ active: editor.isActive('orderedList') }" @click="editor.chain().focus().toggleOrderedList().run()" title="Ordered List">1. List</button>
                <div class="toolbar-divider"></div>
                <button type="button" class="toolbar-btn" :class="{ active: editor.isActive('blockquote') }" @click="editor.chain().focus().toggleBlockquote().run()" title="Blockquote">" "</button>
                <button type="button" class="toolbar-btn" @click="editor.chain().focus().setHorizontalRule().run()" title="Divider">—</button>
                <div class="toolbar-divider"></div>
                <button type="button" class="toolbar-btn" @click="setLink" title="Link">🔗</button>
                <button type="button" class="toolbar-btn" @click="editor.chain().focus().undo().run()" title="Undo">↩</button>
                <button type="button" class="toolbar-btn" @click="editor.chain().focus().redo().run()" title="Redo">↪</button>
              </div>
              <editor-content :editor="editor" class="editor-content" />
            </div>
          </div>

          <!-- Video Platform & Link Configuration -->
          <div class="section-card video-config-card">
            <div class="card-header-flex">
              <div>
                <h3 class="section-title">🎬 {{ locale === 'zh' ? '演示视频平台配置' : 'Video Platform Configuration' }}</h3>
                <p class="section-subtitle">{{ locale === 'zh' ? '支持 TikTok、抖音、小红书、B站、YouTube 等多平台跳转预览' : 'Configure video links on TikTok, Douyin, Xiaohongshu, Bilibili, YouTube' }}</p>
              </div>
              <button type="button" class="btn btn-sm btn-outline" @click="addExtraVideo">
                + {{ locale === 'zh' ? '添加多平台链接' : 'Add Another Platform' }}
              </button>
            </div>

            <!-- Primary video settings -->
            <div class="primary-video-box">
              <label class="form-label font-bold">{{ locale === 'zh' ? '主演示视频平台' : 'Primary Video Platform' }}</label>
              
              <!-- Platform Selector Buttons -->
              <div class="platform-tabs">
                <button
                  v-for="p in videoPlatforms"
                  :key="p.id"
                  type="button"
                  class="platform-tab"
                  :class="{ active: form.preview_platform === p.id }"
                  @click="form.preview_platform = p.id"
                >
                  <span class="p-icon">{{ p.icon }}</span>
                  <span class="p-name">{{ p.name }}</span>
                </button>
              </div>

              <div class="form-row mt-3">
                <div class="form-group flex-2">
                  <label class="form-label">{{ locale === 'zh' ? '视频/笔记超链接' : 'Video / Post URL' }}</label>
                  <div class="input-with-action">
                    <input
                      v-model="form.preview_video_url"
                      type="url"
                      class="form-input"
                      :placeholder="getCurrentPlatformPlaceholder(form.preview_platform)"
                      @input="autoDetectPlatform(form.preview_video_url, (detected) => { form.preview_platform = detected })"
                    />
                    <a
                      v-if="form.preview_video_url"
                      :href="form.preview_video_url"
                      target="_blank"
                      rel="noopener noreferrer"
                      class="btn btn-sm btn-ghost test-link-btn"
                      title="测试打开链接"
                    >
                      ↗ {{ locale === 'zh' ? '测试打开' : 'Test' }}
                    </a>
                  </div>
                </div>

                <div class="form-group flex-1">
                  <label class="form-label">{{ locale === 'zh' ? '按钮自定义文案（可选）' : 'Custom Button Label' }}</label>
                  <input
                    v-model="form.video_title"
                    type="text"
                    class="form-input"
                    :placeholder="locale === 'zh' ? '如: 抖音实拍演示' : 'e.g. Douyin Video'"
                  />
                </div>
              </div>
            </div>

            <!-- Additional multi-platform videos list -->
            <div v-if="form.extra_videos && form.extra_videos.length > 0" class="extra-videos-list mt-4">
              <label class="form-label font-bold">{{ locale === 'zh' ? '更多平台备选演示链接' : 'Additional Platforms' }}</label>
              <div v-for="(v, idx) in form.extra_videos" :key="idx" class="extra-video-row">
                <select v-model="v.platform" class="form-input form-select p-select">
                  <option v-for="p in videoPlatforms" :key="p.id" :value="p.id">
                    {{ p.icon }} {{ p.name }}
                  </option>
                </select>
                <input
                  v-model="v.url"
                  type="url"
                  class="form-input"
                  :placeholder="getCurrentPlatformPlaceholder(v.platform)"
                  @input="autoDetectPlatform(v.url, (detected) => { v.platform = detected })"
                />
                <input
                  v-model="v.title"
                  type="text"
                  class="form-input title-input"
                  :placeholder="locale === 'zh' ? '如: 小红书笔记' : 'Title'"
                />
                <button type="button" class="btn btn-sm btn-ghost text-danger" @click="removeExtraVideo(idx)" title="删除此项">
                  ✕
                </button>
              </div>
            </div>

            <!-- TikTok preview status badge -->
            <div v-if="adminTikTokId" class="tiktok-detected-box mt-3">
              <span class="detect-badge">🎬 TikTok 内嵌播放已激活</span>
              <span class="detect-text">检测到视频 ID: <strong>{{ adminTikTokId }}</strong>，前台详情页将自动启用 iPhone 手机画幅内嵌实机播放！</span>
            </div>
          </div>

          <!-- ⚡ iOS 17+ Interactive Wallpaper & .tendies Simulation Section -->
          <div class="section-card interactive-config-card">
            <div class="card-header-flex">
              <div>
                <h3 class="section-title">⚡ {{ locale === 'zh' ? 'iOS 交互壁纸与实机模拟配置 (.tendies)' : 'Interactive Wallpaper & .tendies Setup' }}</h3>
                <p class="section-subtitle">{{ locale === 'zh' ? '上传相关图层或 .tendies 文件后，前台将自动激活 iPhone 3D 陀螺仪景深与长按动效实机模拟' : 'Enable live 3D depth, gyroscope tilt, and long-press motion simulation on the client' }}</p>
              </div>
              <div class="toggle-switch-wrapper">
                <span class="toggle-text">{{ form.is_interactive ? (locale === 'zh' ? '已启用实机模拟' : 'Enabled') : (locale === 'zh' ? '已关闭模拟' : 'Disabled') }}</span>
                <div class="toggle-switch" :class="{ on: form.is_interactive }" @click="form.is_interactive = !form.is_interactive">
                  <div class="toggle-thumb"></div>
                </div>
              </div>
            </div>

            <!-- Configuration panel when active -->
            <div v-if="form.is_interactive" class="interactive-panel">
              <!-- Type selector -->
              <div class="form-group">
                <label class="form-label font-bold">{{ locale === 'zh' ? '交互模式规格' : 'Interactive Type' }}</label>
                <div class="interactive-type-pills">
                  <button
                    type="button"
                    class="type-pill-btn"
                    :class="{ active: form.interactive_type === 'tendies' }"
                    @click="form.interactive_type = 'tendies'"
                  >
                    🍎 {{ locale === 'zh' ? 'iOS 原生 .tendies (PosterBoard 景深+时钟遮挡)' : 'iOS .tendies Package' }}
                  </button>
                  <button
                    type="button"
                    class="type-pill-btn"
                    :class="{ active: form.interactive_type === 'gyro_depth' }"
                    @click="form.interactive_type = 'gyro_depth'"
                  >
                    🌌 {{ locale === 'zh' ? '3D 陀螺仪多层视差景深' : '3D Gyro Parallax' }}
                  </button>
                  <button
                    type="button"
                    class="type-pill-btn"
                    :class="{ active: form.interactive_type === 'live_motion' }"
                    @click="form.interactive_type = 'live_motion'"
                  >
                    🎬 {{ locale === 'zh' ? '长按触发 Live 动态流' : 'Live Motion Stream' }}
                  </button>
                </div>
              </div>

              <!-- Assets upload grid -->
              <div class="assets-upload-grid">
                <!-- 1. .tendies package file -->
                <div class="asset-slot-card">
                  <div class="slot-header">
                    <span class="slot-badge badge-purple">📦 .tendies 安装包</span>
                    <span class="slot-hint">{{ locale === 'zh' ? '原生 iOS 壁纸包' : 'iOS Bundle' }}</span>
                  </div>
                  <div class="slot-input-wrap">
                    <input
                      v-model="form.interactive_config.tendies_file"
                      type="text"
                      class="form-input form-input-sm"
                      :placeholder="locale === 'zh' ? '可直接上传或填写相对路径' : 'URL or path'"
                    />
                    <label class="btn btn-sm btn-outline file-upload-btn">
                      <span>{{ uploadingSlot === 'tendies' ? '上传中...' : (form.interactive_config.tendies_file ? '重新替换' : '上传文件') }}</span>
                      <input
                        type="file"
                        accept=".tendies,.zip"
                        class="hidden-file-input"
                        @change="(e) => handleInteractiveUpload(e, 'tendies_file', 'tendies')"
                      />
                    </label>
                    <button
                      v-if="form.interactive_config.tendies_file || form.interactive_config.bg_layers?.length"
                      type="button"
                      class="btn btn-sm btn-ghost text-danger ml-1"
                      @click="clearInteractiveConfig"
                      title="清空当前所有图层与配置，恢复空白"
                    >
                      清空
                    </button>
                  </div>
                  <small class="slot-tip">{{ locale === 'zh' ? '标记此商品支持 iOS 17+ 原生免越狱注入，大文件仍建议存入网盘交付' : 'Native package identifier' }}</small>
                  <div v-if="form.interactive_config.bg_layers?.length || form.interactive_config.floating_layers?.length" class="tendies-layers-badge mt-2">
                    <span class="badge-tag">
                      ✨ 已解析 CAML 架构: {{ form.interactive_config.bg_layers?.length || 0 }} 背景图层 + {{ form.interactive_config.floating_layers?.length || 0 }} 时钟浮动层 (共 {{ form.interactive_config.total_image_count || (form.interactive_config.bg_layers?.length + form.interactive_config.floating_layers?.length) }} 个元素，画布 {{ form.interactive_config.canvas_width || 390 }}x{{ form.interactive_config.canvas_height || 844 }})
                    </span>
                  </div>
                </div>

                <!-- 2. Background layer -->
                <div class="asset-slot-card">
                  <div class="slot-header">
                    <span class="slot-badge badge-blue">🖼️ 演示背景层 *</span>
                    <span class="slot-hint">{{ locale === 'zh' ? '底层画面' : 'Background' }}</span>
                  </div>
                  <div class="slot-input-wrap">
                    <input
                      v-model="form.interactive_config.background_url"
                      type="text"
                      class="form-input form-input-sm"
                      :placeholder="locale === 'zh' ? '底层图片/视频 URL' : 'Background URL'"
                    />
                    <label class="btn btn-sm btn-outline file-upload-btn">
                      <span>{{ uploadingSlot === 'bg' ? '上传中...' : '上传' }}</span>
                      <input
                        type="file"
                        accept="image/*,video/*"
                        class="hidden-file-input"
                        @change="(e) => handleInteractiveUpload(e, 'background_url', 'bg')"
                      />
                    </label>
                    <button
                      v-if="form.interactive_config.background_url || form.interactive_config.bg_layers?.length"
                      type="button"
                      class="btn btn-sm btn-ghost text-danger ml-1"
                      @click="clearBgSlot"
                      title="清空当前背景层"
                    >
                      ✕
                    </button>
                  </div>
                  <small class="slot-tip">{{ locale === 'zh' ? '位于时钟下方，作为沉浸式背景基底' : 'Base layer behind clock' }}</small>
                </div>

                <!-- 3. Foreground Subject Cutout layer (For iOS 17 depth clock occlusion!) -->
                <div class="asset-slot-card highlight-slot">
                  <div class="slot-header">
                    <span class="slot-badge badge-gold">✨ 演示前景主体层 (景深抠图) *</span>
                    <span class="slot-hint">{{ locale === 'zh' ? '遮挡时钟' : 'Subject Cutout' }}</span>
                  </div>
                  <div class="slot-input-wrap">
                    <input
                      v-model="form.interactive_config.foreground_url"
                      type="text"
                      class="form-input form-input-sm"
                      :placeholder="locale === 'zh' ? '透明主体 PNG/WebP' : 'Foreground cutout URL'"
                    />
                    <label class="btn btn-sm btn-outline file-upload-btn btn-gold-outline">
                      <span>{{ uploadingSlot === 'fg' ? '上传中...' : '上传抠图' }}</span>
                      <input
                        type="file"
                        accept="image/png,image/webp"
                        class="hidden-file-input"
                        @change="(e) => handleInteractiveUpload(e, 'foreground_url', 'fg')"
                      />
                    </label>
                    <button
                      v-if="form.interactive_config.foreground_url || form.interactive_config.floating_layers?.length"
                      type="button"
                      class="btn btn-sm btn-ghost text-danger ml-1"
                      @click="clearFgSlot"
                      title="清空当前前景层"
                    >
                      ✕
                    </button>
                  </div>
                  <small class="slot-tip">{{ locale === 'zh' ? '★ 核心关键：透明背景人物/主体，浮于时钟上方实现景深穿透效果' : 'Crucial: transparent cutout to overlap clock' }}</small>
                </div>

                <!-- 4. Motion Video / Particle loop -->
                <div class="asset-slot-card">
                  <div class="slot-header">
                    <span class="slot-badge badge-cyan">🎬 长按触发动态流 (可选)</span>
                    <span class="slot-hint">{{ locale === 'zh' ? 'Live Motion' : 'Touch Action' }}</span>
                  </div>
                  <div class="slot-input-wrap">
                    <input
                      v-model="form.interactive_config.motion_url"
                      type="text"
                      class="form-input form-input-sm"
                      :placeholder="locale === 'zh' ? '短视频 MP4 / WebM' : 'Motion Video URL'"
                    />
                    <label class="btn btn-sm btn-outline file-upload-btn">
                      <span>{{ uploadingSlot === 'motion' ? '上传中...' : '上传视频' }}</span>
                      <input
                        type="file"
                        accept="video/mp4,video/webm,video/quicktime"
                        class="hidden-file-input"
                        @change="(e) => handleInteractiveUpload(e, 'motion_url', 'motion')"
                      />
                    </label>
                  </div>
                  <small class="slot-tip">{{ locale === 'zh' ? '用户长按锁屏时即时浮现播放的动态视频或粒子流' : 'Plays on touch/hold' }}</small>
                </div>
              </div>

              <!-- Interactive parameters tuning -->
              <div class="interactive-controls-row mt-3">
                <div class="control-toggle-item">
                  <label class="control-label font-bold">📳 陀螺仪 / 鼠标 3D 视差感应</label>
                  <input type="checkbox" v-model="form.interactive_config.gyro_enabled" class="checkbox-input" />
                  <span class="status-badge" :class="form.interactive_config.gyro_enabled ? 'badge-on' : 'badge-off'">
                    {{ form.interactive_config.gyro_enabled ? '已开启' : '已关闭 (平面无晃动)' }}
                  </span>
                </div>
                <div class="control-toggle-item">
                  <label class="control-label font-bold">🌙 支持 AOD 熄屏常亮模拟</label>
                  <input type="checkbox" v-model="form.interactive_config.has_aod" class="checkbox-input" />
                </div>
                <div class="control-slider-item" :class="{ 'control-disabled': !form.interactive_config.gyro_enabled }">
                  <div class="slider-label-row">
                    <label class="control-label font-bold">
                      视差景深强度:
                      <span class="depth-val-display" :class="{ 'zero-val': form.interactive_config.depth_intensity === 0, 'neg-val': form.interactive_config.depth_intensity < 0 }">
                        {{ form.interactive_config.depth_intensity === 0 ? '0px (关闭 / 无视差)' : (form.interactive_config.depth_intensity > 0 ? '+' + form.interactive_config.depth_intensity + 'px (正向视差)' : form.interactive_config.depth_intensity + 'px (反向视差)') }}
                      </span>
                    </label>
                    <div class="depth-num-box">
                      <input
                        type="number"
                        min="-50"
                        max="50"
                        step="1"
                        v-model.number="form.interactive_config.depth_intensity"
                        class="form-input-sm depth-num-input"
                        title="可直接输入负数、0或正数"
                      />
                      <span class="unit-text">px</span>
                    </div>
                  </div>
                  <div class="slider-input-group">
                    <span class="range-bound">-50px</span>
                    <input
                      type="range"
                      min="-50"
                      max="50"
                      step="1"
                      v-model.number="form.interactive_config.depth_intensity"
                      class="range-slider"
                      :disabled="!form.interactive_config.gyro_enabled"
                    />
                    <span class="range-bound">+50px</span>
                  </div>
                  <!-- Quick preset buttons -->
                  <div class="depth-presets">
                    <span class="preset-label">快捷调整:</span>
                    <button type="button" class="preset-chip" :class="{ active: form.interactive_config.depth_intensity === 0 }" @click="form.interactive_config.depth_intensity = 0">0px (关闭)</button>
                    <button type="button" class="preset-chip" :class="{ active: form.interactive_config.depth_intensity === -15 }" @click="form.interactive_config.depth_intensity = -15">-15px (反向)</button>
                    <button type="button" class="preset-chip" :class="{ active: form.interactive_config.depth_intensity === -25 }" @click="form.interactive_config.depth_intensity = -25">-25px (反向强)</button>
                    <button type="button" class="preset-chip" :class="{ active: form.interactive_config.depth_intensity === 15 }" @click="form.interactive_config.depth_intensity = 15">+15px (标准)</button>
                    <button type="button" class="preset-chip" :class="{ active: form.interactive_config.depth_intensity === 25 }" @click="form.interactive_config.depth_intensity = 25">+25px (强烈)</button>
                  </div>
                </div>
              </div>

              <!-- Admin Live Simulator Mini Viewport -->
              <div class="admin-simulator-preview mt-4">
                <div class="preview-header-bar">
                  <span class="preview-tag">
                    👁️ 后台实时视差效果校验 
                    <template v-if="!form.interactive_config.gyro_enabled">
                      <span class="tag-status-pill off">(陀螺仪已关闭 - 平面无晃动)</span>
                    </template>
                    <template v-else-if="form.interactive_config.depth_intensity === 0">
                      <span class="tag-status-pill neutral">(景深已关闭 0px - 平面无视差)</span>
                    </template>
                    <template v-else-if="form.interactive_config.depth_intensity < 0">
                      <span class="tag-status-pill neg">(反向景深 {{ form.interactive_config.depth_intensity }}px - 鼠标反向偏移)</span>
                    </template>
                    <template v-else>
                      <span class="tag-status-pill on">(正向景深 +{{ form.interactive_config.depth_intensity }}px - 移动鼠标测试)</span>
                    </template>
                  </span>
                  <span class="preview-status" v-if="form.interactive_config.foreground_url">✓ 景深主体层已就绪</span>
                  <span class="preview-status warning" v-else>⚠️ 尚未配置前景抠图层，当前为平面预览</span>
                </div>
                <div
                  class="mini-iphone-mockup"
                  @mousemove="handleAdminMouseMove"
                  @mouseleave="handleAdminMouseLeave"
                >
                  <!-- 0. Dynamic Island on mini mockup -->
                  <div class="mini-dynamic-island">
                    <div class="mini-camera-dot"></div>
                  </div>

                  <!-- 1. Background layer(s) -->
                  <template v-if="form.interactive_config.bg_layers?.length">
                    <div
                      v-for="(layer, idx) in form.interactive_config.bg_layers"
                      :key="layer.id || idx"
                      class="mini-layer mini-caml-layer"
                      :style="getAdminLayerStyle(layer)"
                    >
                      <div class="mini-anim-pos" :style="getAdminPosAnimStyle(layer)">
                        <div class="mini-anim-rot" :style="getAdminRotAnimStyle(layer)">
                          <img
                            :src="getCoverFullUrl(layer.url)"
                            class="mini-sprite-inner-img"
                            :style="getAdminImgAnimStyle(layer)"
                          />
                        </div>
                      </div>
                    </div>
                  </template>
                  <img
                    v-else-if="form.interactive_config.background_url || form.cover_image"
                    :src="getCoverFullUrl(form.interactive_config.background_url || form.cover_image)"
                    class="mini-layer mini-bg"
                    :style="adminBgTransform"
                  />

                  <!-- Simulated iOS Clock (Placed between background & foreground) -->
                  <div class="mini-clock" :style="adminClockTransform">
                    <div class="mini-date-row">
                      <span class="mini-lock-icon">🔒</span>
                      <span class="mini-date-text">10月20日 星期二</span>
                    </div>
                    <span class="mini-time">09:41</span>
                  </div>

                  <!-- 2. Foreground / Floating layer(s) -->
                  <template v-if="form.interactive_config.floating_layers?.length">
                    <div
                      v-for="(layer, idx) in form.interactive_config.floating_layers"
                      :key="layer.id || idx"
                      class="mini-layer mini-caml-layer"
                      :style="getAdminLayerStyle(layer)"
                    >
                      <div class="mini-anim-pos" :style="getAdminPosAnimStyle(layer)">
                        <div class="mini-anim-rot" :style="getAdminRotAnimStyle(layer)">
                          <img
                            :src="getCoverFullUrl(layer.url)"
                            class="mini-sprite-inner-img"
                            :style="getAdminImgAnimStyle(layer)"
                          />
                        </div>
                      </div>
                    </div>
                  </template>
                  <img
                    v-else-if="form.interactive_config.foreground_url"
                    :src="getCoverFullUrl(form.interactive_config.foreground_url)"
                    class="mini-layer mini-fg"
                    :style="adminFgTransform"
                  />
                  <!-- Security Watermark overlay -->
                  <div class="mini-watermark">WallArt 4K 交互壁纸演示</div>
                  <!-- Island -->
                  <div class="mini-island"></div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Right column: media & links -->
        <div class="edit-sidebar">
          <!-- Publish, Pin & Sort Settings -->
          <div class="section-card">
            <h3 class="section-title">{{ locale === 'zh' ? '展示与排序' : 'Display & Sort' }}</h3>
            
            <div class="publish-toggle" style="margin-bottom: 14px;">
              <div>
                <p class="toggle-label">{{ form.is_published ? (locale === 'zh' ? '已上架' : 'Live') : (locale === 'zh' ? '已下架' : 'Draft') }}</p>
                <p class="toggle-desc">{{ form.is_published ? (locale === 'zh' ? '用户可见' : 'Visible to users') : (locale === 'zh' ? '仅管理员可见' : 'Admin only') }}</p>
              </div>
              <div class="toggle-switch" :class="{ on: form.is_published }" @click="form.is_published = !form.is_published">
                <div class="toggle-thumb"></div>
              </div>
            </div>

            <!-- Pin to Top Switch -->
            <div class="publish-toggle" style="margin-bottom: 14px; padding-top: 12px; border-top: 1px solid var(--c-border-light, rgba(255, 255, 255, 0.08));">
              <div>
                <p class="toggle-label" style="display: flex; align-items: center; gap: 4px;">
                  <span>📌</span>
                  <span :style="{ color: form.is_pinned ? '#ef4444' : 'inherit', fontWeight: 'bold' }">
                    {{ form.is_pinned ? (locale === 'zh' ? '最高权限 · 已置顶' : 'Pinned') : (locale === 'zh' ? '商品置顶' : 'Pin to Top') }}
                  </span>
                </p>
                <p class="toggle-desc">{{ locale === 'zh' ? '置顶商品排在最前面，最新置顶的排最前' : 'Highest authority, latest pinned sorts first' }}</p>
              </div>
              <div class="toggle-switch switch-pinned" :class="{ on: form.is_pinned }" @click="form.is_pinned = !form.is_pinned">
                <div class="toggle-thumb"></div>
              </div>
            </div>

            <!-- Sort Order Input -->
            <div class="sort-order-block" style="padding-top: 12px; border-top: 1px solid var(--c-border-light, rgba(255, 255, 255, 0.08));">
              <label class="form-label" style="font-size: 0.85rem; font-weight: 600;">
                🔢 {{ locale === 'zh' ? '排序权重 (0 为默认)' : 'Sort Order (0 default)' }}
              </label>
              <input
                v-model.number="form.sort_order"
                type="number"
                class="form-input"
                placeholder="0"
                style="margin-top: 4px;"
              />
              <p class="toggle-desc" style="margin-top: 4px;">
                {{ locale === 'zh' ? '0是默认排序，数字越大排序越靠前' : '0 is default, larger number sorts earlier' }}
              </p>
            </div>
          </div>

          <!-- Cover Image Upload & Preview -->
          <div class="section-card">
            <div class="card-header-flex">
              <h3 class="section-title">{{ t('admin.cover_image') }}</h3>
              <div class="cover-mode-toggle">
                <button
                  type="button"
                  class="mode-btn"
                  :class="{ active: coverMode === 'upload' }"
                  @click="coverMode = 'upload'"
                >
                  {{ locale === 'zh' ? '本地上传' : 'Upload' }}
                </button>
                <button
                  type="button"
                  class="mode-btn"
                  :class="{ active: coverMode === 'url' }"
                  @click="coverMode = 'url'"
                >
                  URL
                </button>
              </div>
            </div>

            <!-- Upload Mode -->
            <div v-if="coverMode === 'upload'" class="upload-zone-wrapper">
              <input
                ref="fileInputRef"
                type="file"
                accept="image/png, image/jpeg, image/jpg, image/webp, image/gif"
                class="hidden-file-input"
                @change="onFileSelected"
              />

              <div
                class="upload-dropzone"
                :class="{ 'is-dragover': isDragOver, 'has-image': !!form.cover_image }"
                @dragover.prevent="isDragOver = true"
                @dragleave.prevent="isDragOver = false"
                @drop.prevent="onFileDropped"
                @click="triggerFileInput"
              >
                <div v-if="uploading" class="upload-status">
                  <span class="spinner" style="width:32px; height:32px;"></span>
                  <p>{{ locale === 'zh' ? '正在上传图片...' : 'Uploading image...' }}</p>
                </div>
                <div v-else-if="form.cover_image" class="cover-preview-box">
                  <img :src="getCoverFullUrl(form.cover_image)" alt="Cover" />
                  <div class="cover-overlay">
                    <span>📷 {{ locale === 'zh' ? '点击更换本地图片' : 'Click to replace image' }}</span>
                  </div>
                </div>
                <div v-else class="upload-prompt">
                  <div class="upload-icon">📁</div>
                  <p class="font-bold">{{ locale === 'zh' ? '点击或拖拽壁纸截图至此处' : 'Click or drop cover image here' }}</p>
                  <small class="text-secondary">{{ locale === 'zh' ? '支持 JPG, PNG, WEBP 高清壁纸 (最大 15MB)' : 'Supports JPG, PNG, WEBP up to 15MB' }}</small>
                </div>
              </div>

              <div v-if="form.cover_image" class="cover-meta-actions">
                <span class="cover-path-preview" :title="form.cover_image">{{ form.cover_image }}</span>
                <button type="button" class="btn btn-sm btn-ghost text-danger" @click="form.cover_image = ''">
                  {{ locale === 'zh' ? '移除' : 'Remove' }}
                </button>
              </div>
            </div>

            <!-- URL Mode -->
            <div v-else class="url-input-box">
              <div class="form-group">
                <label class="form-label">{{ locale === 'zh' ? '外部图片 URL' : 'External Image URL' }}</label>
                <input v-model="form.cover_image" type="url" class="form-input" placeholder="https://..." />
              </div>
              <div v-if="form.cover_image" class="image-preview mt-2">
                <img :src="getCoverFullUrl(form.cover_image)" alt="Preview" @error="form.cover_image = ''" />
              </div>
            </div>
          </div>

          <!-- Cloud Drive Links -->
          <div class="section-card">
            <h3 class="section-title">☁️ {{ t('admin.drive_links') }}</h3>
            <p class="section-subtitle">{{ locale === 'zh' ? '只有解锁壁纸或有效会员才能看到下载提取信息' : 'Protected for unlocked users and members only' }}</p>
            
            <div v-for="drive in drives" :key="drive.key" class="drive-entry">
              <div class="drive-header">
                <span class="drive-icon">{{ drive.icon }}</span>
                <span class="drive-name">{{ drive.name }}</span>
              </div>
              <div class="form-group">
                <label class="form-label">URL</label>
                <input v-model="form.downloads[drive.key].url" type="url" class="form-input" :placeholder="drive.placeholder" />
              </div>
              <div class="form-group">
                <label class="form-label">{{ locale === 'zh' ? '提取码（可选）' : 'Extract Code (optional)' }}</label>
                <input v-model="form.downloads[drive.key].extract_code" type="text" class="form-input" :placeholder="locale === 'zh' ? '如: ab12' : 'e.g. ab12'" maxlength="16" />
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Save buttons -->
      <div class="edit-footer">
        <router-link to="/admin/products" class="btn btn-outline">{{ t('common.cancel') }}</router-link>
        <button type="submit" class="btn btn-primary btn-lg" :disabled="saving">
          <span v-if="saving" class="spinner"></span>
          {{ saving ? t('common.loading') : t('admin.save') }}
        </button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onBeforeUnmount, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Link from '@tiptap/extension-link'
import Placeholder from '@tiptap/extension-placeholder'
import { productAPI, configAPI, categoryAPI } from '@/api'
import { settingsAPI } from '@/api'
import { useToast } from '@/composables/useToast'
import { resolveAndDecryptAssetUrl } from '@/utils/assetCrypto'

const route = useRoute()
const router = useRouter()
const { t, locale } = useI18n()
const { success, error } = useToast()

// Robust computed isNew check
const isNew = computed(() => !route.params.id || route.params.id === 'new' || route.path.endsWith('/new'))
const saving = ref(false)
const uploading = ref(false)
const isDragOver = ref(false)
const coverMode = ref('upload')
const fileInputRef = ref(null)
const tagInput = ref('')
const errors = reactive({ title_zh: '' })

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
  } catch (e) {
    console.warn('Failed to load categories in product edit:', e)
  }
}

const drives = [
  { key: 'baidu', icon: '🔵', name: '百度网盘', placeholder: 'https://pan.baidu.com/s/...' },
  { key: 'quark', icon: '⚡', name: '夸克网盘', placeholder: 'https://pan.quark.cn/s/...' },
  { key: 'google', icon: '🔴', name: 'Google Drive', placeholder: 'https://drive.google.com/...' }
]

const videoPlatforms = [
  { id: 'douyin', name: '抖音', icon: '🎶', placeholder: 'https://v.douyin.com/... 或 https://www.douyin.com/video/...' },
  { id: 'xiaohongshu', name: '小红书', icon: '📕', placeholder: 'https://www.xiaohongshu.com/explore/... 或 http://xhslink.com/...' },
  { id: 'tiktok', name: 'TikTok', icon: '🎵', placeholder: 'https://www.tiktok.com/@user/video/...' },
  { id: 'bilibili', name: '哔哩哔哩', icon: '📺', placeholder: 'https://www.bilibili.com/video/BV...' },
  { id: 'youtube', name: 'YouTube', icon: '▶️', placeholder: 'https://www.youtube.com/watch?v=...' },
  { id: 'direct', name: '直链/其他', icon: '🔗', placeholder: 'https://example.com/demo.mp4' }
]

const exchangeRate = ref(7.25)
const isManualCny = ref(false)
const isManualUsd = ref(false)

function onCnyInput() {
  isManualCny.value = true
  if (!isManualUsd.value && form.price_cny !== '') {
    const cnyVal = parseFloat(form.price_cny)
    if (!isNaN(cnyVal) && cnyVal >= 0) {
      const rate = exchangeRate.value || 7.25
      form.price_usd = (cnyVal / rate).toFixed(2)
      form.price = form.price_usd
    } else {
      form.price_usd = ''
      form.price = ''
    }
  } else if (form.price_cny === '' && !isManualUsd.value) {
    form.price_usd = ''
    form.price = ''
  }
}

function onUsdInput() {
  isManualUsd.value = true
  if (!isManualCny.value && form.price_usd !== '') {
    const usdVal = parseFloat(form.price_usd)
    if (!isNaN(usdVal) && usdVal >= 0) {
      const rate = exchangeRate.value || 7.25
      form.price_cny = (usdVal * rate).toFixed(2)
      form.price = form.price_usd
    } else {
      form.price_cny = ''
      form.price = ''
    }
  } else if (form.price_usd === '' && !isManualCny.value) {
    form.price_cny = ''
    form.price = ''
  }
}

const form = reactive({
  title_zh: '',
  title_en: '',
  price: '',
  price_usd: '',
  price_cny: '',
  verify_code: '',
  category: '',
  tags: [],
  allowed_vip_levels: ['monthly', 'quarterly', 'yearly', 'lifetime'],
  cover_image: '',
  preview_platform: 'douyin',
  preview_video_url: '',
  video_title: '',
  extra_videos: [],
  description: '',
  is_published: true,
  is_pinned: false,
  sort_order: 0,
  is_interactive: false,
  interactive_type: 'tendies',
  interactive_config: {
    tendies_file: '',
    background_url: '',
    foreground_url: '',
    motion_url: '',
    gyro_enabled: true,
    depth_intensity: 15,
    has_aod: true
  },
  downloads: {
    baidu: { url: '', extract_code: '' },
    quark: { url: '', extract_code: '' },
    google: { url: '', extract_code: '' }
  }
})

const uploadingSlot = ref('')
const adminMouse = reactive({ x: 0, y: 0 })

function handleAdminMouseMove(e) {
  const rect = e.currentTarget.getBoundingClientRect()
  const normX = (e.clientX - rect.left) / rect.width - 0.5
  const normY = (e.clientY - rect.top) / rect.height - 0.5
  adminMouse.x = normX
  adminMouse.y = normY
}

function handleAdminMouseLeave() {
  adminMouse.x = 0
  adminMouse.y = 0
}

const adminBgTransform = computed(() => {
  if (!form.interactive_config.gyro_enabled) return {}
  const int = Number(form.interactive_config.depth_intensity ?? 15)
  if (int === 0) return {}
  const moveX = -adminMouse.x * (int * 0.5)
  const moveY = -adminMouse.y * (int * 0.5)
  return {
    transform: `scale(1.05) translate3d(${moveX.toFixed(1)}px, ${moveY.toFixed(1)}px, 0)`
  }
})

const adminClockTransform = computed(() => {
  if (!form.interactive_config.gyro_enabled) return {}
  const int = Number(form.interactive_config.depth_intensity ?? 15)
  if (int === 0) return {}
  const moveX = -adminMouse.x * (int * 0.2)
  const moveY = -adminMouse.y * (int * 0.2)
  return {
    transform: `translate3d(${moveX.toFixed(1)}px, ${moveY.toFixed(1)}px, 0)`
  }
})

const adminFgTransform = computed(() => {
  if (!form.interactive_config.gyro_enabled) return {}
  const int = Number(form.interactive_config.depth_intensity ?? 15)
  if (int === 0) return {}
  const moveX = adminMouse.x * (int * 1.2)
  const moveY = adminMouse.y * (int * 1.2)
  return {
    transform: `scale(1.05) translate3d(${moveX.toFixed(1)}px, ${moveY.toFixed(1)}px, 0)`
  }
})

function getAdminLayerStyle(layer) {
  const cW = form.interactive_config.canvas_width || 393
  const cH = form.interactive_config.canvas_height || 852
  const int = Number(form.interactive_config.depth_intensity ?? 15)
  let moveX = 0
  let moveY = 0
  if (form.interactive_config.gyro_enabled && int !== 0) {
    if (layer.category === 'background') {
      const factor = (layer.order + 1) * 0.1
      moveX = -adminMouse.x * (int * factor)
      moveY = -adminMouse.y * (int * factor)
    } else {
      const factor = 0.5 + (layer.order + 1) * 0.15
      moveX = adminMouse.x * (int * factor)
      moveY = adminMouse.y * (int * factor)
    }
  }

  // 1. Determine base locked state values
  const lockState = layer.states?.locked

  // Position animations in CAML specify absolute keyframes.
  // tendies_parser computed DeltaX = (rawVal - layer.pos_x) and DeltaY = (screenY - layer.pos_y).
  // Therefore, translate3d(dx, dy, 0) is relative to layer.pos_x / layer.pos_y.
  // For axes with animations, we use layer base pos to avoid double offset (which caused the ball overflow).
  // For other properties and non-animated layers, strictly respect lockState (e.g. wtksoi offscreen/hidden).
  const hasPosXAnim = Array.isArray(layer.animations) && layer.animations.some(a => a.key_path === 'position.x' || a.key_path === 'position')
  const hasPosYAnim = Array.isArray(layer.animations) && layer.animations.some(a => a.key_path === 'position.y' || a.key_path === 'position')
  const hasRotAnim = Array.isArray(layer.animations) && layer.animations.some(a => a.key_path === 'transform.rotation.z')

  const posX = (lockState && lockState.pos_x !== undefined) ? lockState.pos_x : layer.pos_x
  const posY = (lockState && lockState.pos_y !== undefined) ? lockState.pos_y : layer.pos_y
  const width = (lockState && lockState.width !== undefined && lockState.width > 0) ? lockState.width : layer.width
  const height = (lockState && lockState.height !== undefined && lockState.height > 0) ? lockState.height : layer.height
  const opacity = (lockState && lockState.opacity !== undefined) ? lockState.opacity : (layer.opacity !== undefined ? layer.opacity : 1.0)
  const rot = (!hasRotAnim && lockState && lockState.rotation !== undefined) ? lockState.rotation : (layer.rotation || 0)
  const scale = (lockState && lockState.scale !== undefined && lockState.scale > 0) ? lockState.scale : (layer.scale || 1.0)

  const isFullCanvas = (width >= cW * 0.96) && (height >= cH * 0.96) && Math.abs(posX - cW / 2) < 4 && Math.abs(posY - cH / 2) < 4 && Math.abs(rot) < 0.001 && Math.abs(scale - 1.0) < 0.001

  const leftPct = (posX / cW) * 100
  const topPct = (posY / cH) * 100
  const widthPct = (width / cW) * 100
  const heightPct = (height / cH) * 100
  const zIndex = layer.category === 'floating' ? 30 + (layer.order || 0) : 10 + (layer.order || 0)
  const rotDeg = rot * (180 / Math.PI)
  const rotStr = Math.abs(rotDeg) > 0.01 ? ` rotate(${rotDeg.toFixed(2)}deg)` : ''
  const scaleStr = Math.abs(scale - 1.0) > 0.001 ? ` scale(${scale.toFixed(3)})` : ''

  if (isFullCanvas) {
    return {
      position: 'absolute',
      left: '0',
      top: '0',
      width: '100%',
      height: '100%',
      transform: `translate3d(${moveX.toFixed(1)}px, ${moveY.toFixed(1)}px, 0)`,
      transformOrigin: '50% 50%',
      opacity: opacity,
      visibility: opacity < 0.005 ? 'hidden' : 'visible',
      zIndex: zIndex,
      pointerEvents: 'none',
      objectFit: 'cover'
    }
  }

  return {
    position: 'absolute',
    left: `${leftPct}%`,
    top: `${topPct}%`,
    width: `${widthPct}%`,
    height: `${heightPct}%`,
    transform: `translate(-50%, -50%) translate3d(${moveX.toFixed(1)}px, ${moveY.toFixed(1)}px, 0)${scaleStr}${rotStr}`,
    transformOrigin: '50% 50%',
    opacity: opacity,
    visibility: opacity < 0.005 ? 'hidden' : 'visible',
    zIndex: zIndex,
    pointerEvents: 'none',
    objectFit: 'fill'
  }
}

function getAdminPosAnimStyle(layer) {
  if (!layer.animations || !layer.animations.length) return { width: '100%', height: '100%' }
  const anim = layer.animations.find(a => a.key_path === 'position.y' || a.key_path === 'position.x' || a.key_path === 'position')
  if (!anim) return { width: '100%', height: '100%' }
  const dir = anim.autoreverses ? 'alternate' : 'normal'
  const timing = anim.key_path === 'position' && !anim.autoreverses ? 'linear' : 'ease-in-out'
  return {
    animation: `admin-pos-${layer.id} ${anim.duration || 5}s ${timing} infinite ${dir}`,
    width: '100%',
    height: '100%'
  }
}

function getAdminRotAnimStyle(layer) {
  if (!layer.animations || !layer.animations.length) return { width: '100%', height: '100%' }
  const anim = layer.animations.find(a => a.key_path === 'transform.rotation.z')
  if (!anim) return { width: '100%', height: '100%' }
  const dir = anim.autoreverses ? 'alternate' : 'normal'
  return {
    animation: `admin-rot-${layer.id} ${anim.duration || 5}s ease-in-out infinite ${dir}`,
    width: '100%',
    height: '100%'
  }
}

function getAdminImgAnimStyle(layer) {
  if (!layer.animations || !layer.animations.length) return { width: '100%', height: '100%' }
  const anim = layer.animations.find(a => a.key_path === 'opacity')
  if (!anim) return { width: '100%', height: '100%' }
  const dir = anim.autoreverses ? 'alternate' : 'normal'
  return {
    animation: `admin-op-${layer.id} ${anim.duration || 5}s ease-in-out infinite ${dir}`,
    width: '100%',
    height: '100%'
  }
}

const adminKeyframesCSS = computed(() => {
  let css = ''
  const bgL = form.interactive_config.bg_layers || []
  const flL = form.interactive_config.floating_layers || []
  const allLayers = [...bgL, ...flL]
  for (const layer of allLayers) {
    if (!layer.animations || !layer.animations.length) continue
    for (const anim of layer.animations) {
      if (anim.key_path === 'position.y' || anim.key_path === 'position.x' || anim.key_path === 'position') {
        const count = Math.max(anim.delta_x?.length || 0, anim.delta_y?.length || 0)
        if (count > 0) {
          let midX = 0
          if (anim.delta_x && anim.delta_x.length > 0) {
            const minX = Math.min(...anim.delta_x)
            const maxX = Math.max(...anim.delta_x)
            midX = (minX + maxX) / 2
          }
          let midY = 0
          if (anim.delta_y && anim.delta_y.length > 0) {
            const minY = Math.min(...anim.delta_y)
            const maxY = Math.max(...anim.delta_y)
            midY = (minY + maxY) / 2
          }
          css += `@keyframes admin-pos-${layer.id} {\n`
          for (let i = 0; i < count; i++) {
            const pct = count === 1 ? 100 : Math.round((i / (count - 1)) * 100)
            const rawDx = anim.delta_x && anim.delta_x[i] !== undefined ? anim.delta_x[i] : 0
            const rawDy = anim.delta_y && anim.delta_y[i] !== undefined ? anim.delta_y[i] : 0
            const dx = rawDx - midX
            const dy = rawDy - midY
            css += `  ${pct}% { transform: translate3d(${dx.toFixed(2)}px, ${dy.toFixed(2)}px, 0); }\n`
          }
          css += `}\n`
        }
      } else if (anim.key_path === 'transform.rotation.z') {
        const count = anim.rot_deg?.length || 0
        if (count > 0) {
          css += `@keyframes admin-rot-${layer.id} {\n`
          for (let i = 0; i < count; i++) {
            const pct = count === 1 ? 100 : Math.round((i / (count - 1)) * 100)
            const deg = anim.rot_deg[i] || 0
            css += `  ${pct}% { transform: rotate(${deg.toFixed(2)}deg); }\n`
          }
          css += `}\n`
        }
      } else if (anim.key_path === 'opacity') {
        const count = anim.opacity?.length || 0
        if (count > 0) {
          css += `@keyframes admin-op-${layer.id} {\n`
          for (let i = 0; i < count; i++) {
            const pct = count === 1 ? 100 : Math.round((i / (count - 1)) * 100)
            const op = anim.opacity[i] || 0
            css += `  ${pct}% { opacity: ${op.toFixed(3)}; }\n`
          }
          css += `}\n`
        }
      }
    }
  }
  return css
})

watch(adminKeyframesCSS, (newCSS) => {
  if (typeof document !== 'undefined') {
    let styleEl = document.getElementById('admin-caml-keyframes-style')
    if (!styleEl) {
      styleEl = document.createElement('style')
      styleEl.id = 'admin-caml-keyframes-style'
      document.head.appendChild(styleEl)
    }
    styleEl.textContent = newCSS
  }
}, { immediate: true })

onBeforeUnmount(() => {
  if (typeof document !== 'undefined') {
    const el = document.getElementById('admin-caml-keyframes-style')
    if (el) el.remove()
  }
})

function clearInteractiveConfig() {
  form.interactive_config.tendies_file = ''
  form.interactive_config.background_url = ''
  form.interactive_config.foreground_url = ''
  form.interactive_config.motion_url = ''
  form.interactive_config.bg_layers = []
  form.interactive_config.floating_layers = []
  form.interactive_config.fg_layers = []
  form.interactive_config.total_image_count = 0
  form.interactive_config.frames_dir = ''
  form.interactive_config.frame_count = 0
  success(locale.value === 'zh' ? '已清空当前交互素材配置，可重新上传全新文件' : 'Cleared interactive configuration')
}

function clearBgSlot() {
  form.interactive_config.background_url = ''
  form.interactive_config.bg_layers = []
  success(locale.value === 'zh' ? '已清空背景层' : 'Cleared background')
}

function clearFgSlot() {
  form.interactive_config.foreground_url = ''
  form.interactive_config.floating_layers = []
  success(locale.value === 'zh' ? '已清空前景主体层' : 'Cleared foreground')
}

async function handleInteractiveUpload(e, configKey, slotName) {
  const file = e.target.files?.[0]
  if (!file) return
  uploadingSlot.value = slotName
  try {
    const fd = new FormData()
    fd.append('file', file)
    const res = await productAPI.uploadInteractive(fd)
    if (res && res.url) {
      const cW = form.interactive_config.canvas_width || 390
      const cH = form.interactive_config.canvas_height || 844

      if (configKey === 'tendies_file') {
        // Force complete overwrite with the new tendies package
        form.interactive_config.tendies_file = res.url
        form.interactive_config.background_url = res.background_url || ''
        form.interactive_config.foreground_url = res.foreground_url || ''
        if (res.tendies_data) {
          form.interactive_config.canvas_width = res.tendies_data.canvas_width
          form.interactive_config.canvas_height = res.tendies_data.canvas_height
          form.interactive_config.bg_layers = res.tendies_data.bg_layers || []
          form.interactive_config.floating_layers = res.tendies_data.floating_layers || []
          form.interactive_config.fg_layers = res.tendies_data.fg_layers || []
          form.interactive_config.total_image_count = res.tendies_data.total_image_count || 0
        } else {
          form.interactive_config.bg_layers = []
          form.interactive_config.floating_layers = []
          form.interactive_config.fg_layers = []
          form.interactive_config.total_image_count = 0
        }
        form.interactive_config.frames_dir = ''
        form.interactive_config.frame_count = 0

        // Automatically default lock screen frame as cover image
        const lockCover = res.cover_url || res.tendies_data?.cover_url || res.background_url
        if (lockCover) {
          form.cover_image = lockCover
        }
      } else if (configKey === 'background_url') {
        // Explicitly replace the background material
        form.interactive_config.background_url = res.url
        form.interactive_config.bg_layers = [{
          id: 'custom_bg',
          name: 'Custom Background',
          url: res.url,
          category: 'background',
          order: 0,
          pos_x: cW / 2,
          pos_y: cH / 2,
          width: cW,
          height: cH,
          opacity: 1.0,
          rotation: 0,
          states: {
            locked: { pos_x: cW / 2, pos_y: cH / 2, width: cW, height: cH, opacity: 1.0, rotation: 0 },
            unlock: { pos_x: cW / 2, pos_y: cH / 2, width: cW, height: cH, opacity: 1.0, rotation: 0 }
          }
        }]
      } else if (configKey === 'foreground_url') {
        // Explicitly replace the foreground cutout
        form.interactive_config.foreground_url = res.url
        form.interactive_config.floating_layers = [{
          id: 'custom_fg',
          name: 'Custom Foreground Cutout',
          url: res.url,
          category: 'floating',
          order: 0,
          pos_x: cW / 2,
          pos_y: cH / 2,
          width: cW,
          height: cH,
          opacity: 1.0,
          rotation: 0,
          states: {
            locked: { pos_x: cW / 2, pos_y: cH / 2, width: cW, height: cH, opacity: 1.0, rotation: 0 },
            unlock: { pos_x: cW / 2, pos_y: cH / 2, width: cW, height: cH, opacity: 1.0, rotation: 0 }
          }
        }]
      } else {
        form.interactive_config[configKey] = res.url
      }

      const layerMsg = res.tendies_data
        ? `（已解析 ${res.tendies_data.bg_layers?.length || 0} 个背景层 + ${res.tendies_data.floating_layers?.length || 0} 个浮动前景层，共 ${res.tendies_data.total_image_count || 0} 张元素）`
        : ''
      const msg = res.foreground_url || res.tendies_data
        ? (locale.value === 'zh' ? `✨ 素材替换成功${layerMsg}！已完整同步新 CAML 坐标与图层！` : 'Tendies package replaced and auto-extracted!')
        : (locale.value === 'zh' ? '素材替换成功！' : 'Asset replaced!')
      success(msg)
    }
  } catch (err) {
    error(err.message || '上传失败')
  } finally {
    uploadingSlot.value = ''
    e.target.value = ''
  }
}

const adminTikTokId = computed(() => {
  const match = (form.preview_video_url || '').match(/(\d{15,})/)
  if (match) return match[1]
  for (const v of form.extra_videos || []) {
    const m = (v.url || '').match(/(\d{15,})/)
    if (m) return m[1]
  }
  return null
})

// TipTap Editor
const editor = useEditor({
  extensions: [
    StarterKit,
    Link.configure({ openOnClick: false }),
    Placeholder.configure({ placeholder: locale.value === 'zh' ? '在此输入商品详细描述（支持图文与多格式排版）...' : 'Enter product description here...' })
  ],
  content: '',
  onUpdate: ({ editor }) => {
    form.description = editor.getHTML()
  }
})

function getCurrentPlatformPlaceholder(platformId) {
  const p = videoPlatforms.find(item => item.id === platformId)
  return p ? p.placeholder : 'https://...'
}

function autoDetectPlatform(url, callback) {
  if (!url) return
  const lower = url.toLowerCase()
  if (lower.includes('douyin.com')) {
    callback('douyin')
  } else if (lower.includes('xiaohongshu.com') || lower.includes('xhslink.com')) {
    callback('xiaohongshu')
  } else if (lower.includes('tiktok.com')) {
    callback('tiktok')
  } else if (lower.includes('bilibili.com') || lower.includes('b23.tv')) {
    callback('bilibili')
  } else if (lower.includes('youtube.com') || lower.includes('youtu.be')) {
    callback('youtube')
  }
}

function addExtraVideo() {
  if (!form.extra_videos) form.extra_videos = []
  form.extra_videos.push({
    platform: 'xiaohongshu',
    url: '',
    title: ''
  })
}

function removeExtraVideo(idx) {
  form.extra_videos.splice(idx, 1)
}

function triggerFileInput() {
  fileInputRef.value?.click()
}

async function onFileSelected(event) {
  const file = event.target.files?.[0]
  if (file) await uploadFile(file)
}

async function onFileDropped(event) {
  isDragOver.value = false
  const file = event.dataTransfer?.files?.[0]
  if (file) await uploadFile(file)
}

async function uploadFile(file) {
  if (!file.type.startsWith('image/')) {
    error(locale.value === 'zh' ? '请选择图片文件（JPG, PNG, WEBP, GIF）' : 'Please select an image file')
    return
  }
  uploading.value = true
  try {
    const fd = new FormData()
    fd.append('file', file)
    const res = await productAPI.uploadImage(fd)
    if (res && res.url) {
      form.cover_image = res.url
      success(locale.value === 'zh' ? '封面图片上传成功！' : 'Image uploaded successfully!')
    }
  } catch (err) {
    error(err.message || (locale.value === 'zh' ? '图片上传失败' : 'Upload failed'))
  } finally {
    uploading.value = false
    if (fileInputRef.value) fileInputRef.value.value = ''
  }
}

const adminDecryptedUrls = ref({})

watch(() => form.interactive_config, (cfg) => {
  if (!cfg) return
  const allUrls = [
    ...(Array.isArray(cfg.bg_layers) ? cfg.bg_layers.map(l => l.url) : []),
    ...(Array.isArray(cfg.floating_layers) ? cfg.floating_layers.map(l => l.url) : []),
    cfg.background_url,
    cfg.foreground_url
  ].filter(Boolean)

  allUrls.forEach(url => {
    if (url.includes('.dat')) {
      if (!adminDecryptedUrls.value[url]) {
        resolveAndDecryptAssetUrl(url).then(blobUrl => {
          adminDecryptedUrls.value[url] = blobUrl
        })
      }
    }
  })
}, { deep: true, immediate: true })

function getCoverFullUrl(url) {
  if (!url) return ''
  if (adminDecryptedUrls.value[url]) return adminDecryptedUrls.value[url]
  if (url.startsWith('http://') || url.startsWith('https://')) return url
  // Relative uploaded path, prefix with backend or current origin
  return url
}

function setLink() {
  const url = prompt(locale.value === 'zh' ? '输入链接URL:' : 'Enter URL:')
  if (!url) return
  editor.value.chain().focus().setLink({ href: url }).run()
}

function addTag() {
  const tag = tagInput.value.trim().replace(/,$/, '')
  if (tag && !form.tags.includes(tag)) form.tags.push(tag)
  tagInput.value = ''
}

function removeTag(tag) {
  form.tags = form.tags.filter(t => t !== tag)
}

function validate() {
  errors.title_zh = ''
  if (!form.title_zh.trim()) {
    errors.title_zh = locale.value === 'zh' ? '请输入中文标题' : 'Chinese title is required'
    return false
  }
  return true
}

async function handleSave() {
  if (!validate()) return
  saving.value = true
  try {
    // Compile preview videos list
    const videosList = []
    if (form.preview_video_url) {
      videosList.push({
        platform: form.preview_platform || 'douyin',
        url: form.preview_video_url,
        title: form.video_title || ''
      })
    }
    if (form.extra_videos && form.extra_videos.length > 0) {
      form.extra_videos.forEach(v => {
        if (v.url) {
          videosList.push({
            platform: v.platform || 'direct',
            url: v.url,
            title: v.title || ''
          })
        }
      })
    }

    const payload = {
      title_zh: form.title_zh,
      title_en: form.title_en,
      price: parseFloat(form.price_usd) || parseFloat(form.price) || 0,
      price_usd: form.price_usd !== '' ? parseFloat(form.price_usd) : 0,
      price_cny: form.price_cny !== '' ? parseFloat(form.price_cny) : 0,
      verify_code: form.verify_code || '',
      category: form.category || 'general',
      tags: form.tags,
      allowed_vip_levels: form.allowed_vip_levels,
      cover_image: form.cover_image,
      preview_platform: form.preview_platform || 'douyin',
      preview_video_url: form.preview_video_url,
      preview_videos: videosList,
      description: form.description,
      is_published: form.is_published,
      is_pinned: !!form.is_pinned,
      sort_order: form.sort_order || 0,
      is_interactive: form.is_interactive,
      interactive_type: form.interactive_type,
      interactive_config: JSON.stringify(form.interactive_config),
      downloads: Object.entries(form.downloads)
        .filter(([, v]) => v && v.url && String(v.url).trim() !== '')
        .map(([platform, v]) => ({
          platform,
          url: String(v.url).trim(),
          extract_code: String(v.extract_code || '').trim()
        }))
    }

    if (isNew.value) {
      await productAPI.create(payload)
      success(locale.value === 'zh' ? '商品创建成功！' : 'Product created!')
    } else {
      await productAPI.update(route.params.id, payload)
      success(locale.value === 'zh' ? '商品已更新！' : 'Product updated!')
    }
    router.push('/admin/products')
  } catch (e) {
    error(e.message || t('common.error'))
  } finally {
    saving.value = false
  }
}

async function loadProduct() {
  if (isNew.value) return
  try {
    const res = await (productAPI.adminGet ? productAPI.adminGet(route.params.id) : productAPI.get(route.params.id)).catch(() => productAPI.get(route.params.id))
    const p = res.product || res
    
    // Parse tags
    let parsedTags = []
    if (Array.isArray(p.tags)) {
      parsedTags = p.tags
    } else if (p.tags) {
      try { parsedTags = JSON.parse(p.tags) } catch { parsedTags = [p.tags] }
    }

    // Parse preview videos
    let parsedVideos = []
    if (p.preview_videos) {
      try { parsedVideos = JSON.parse(p.preview_videos) } catch { /* ignore */ }
    }

    let primaryPlatform = p.preview_platform || 'douyin'
    let primaryUrl = p.preview_video_url || ''
    let primaryTitle = ''
    let extraList = []

    if (parsedVideos.length > 0) {
      primaryPlatform = parsedVideos[0].platform || primaryPlatform
      primaryUrl = parsedVideos[0].url || primaryUrl
      primaryTitle = parsedVideos[0].title || ''
      extraList = parsedVideos.slice(1)
    }

    Object.assign(form, {
      title_zh: p.title_zh || '',
      title_en: p.title_en || '',
      price: p.price !== undefined ? p.price : '',
      price_usd: p.price_usd !== undefined ? p.price_usd : (p.price || ''),
      price_cny: p.price_cny !== undefined ? p.price_cny : '',
      verify_code: p.verify_code || '',
      category: p.category || '',
      tags: parsedTags,
      allowed_vip_levels: (() => {
        if (p.allowed_vip_levels) {
          try {
            const parsed = typeof p.allowed_vip_levels === 'string' ? JSON.parse(p.allowed_vip_levels) : p.allowed_vip_levels
            if (Array.isArray(parsed)) return parsed
          } catch (e) {
            console.warn('Failed to parse allowed_vip_levels', e)
          }
        }
        return availableVIPPlans.value.map(pl => pl.type)
      })(),
      cover_image: p.cover_image || '',
      preview_platform: primaryPlatform,
      preview_video_url: primaryUrl,
      video_title: primaryTitle,
      extra_videos: extraList,
      description: p.description || '',
      is_published: p.is_published !== undefined ? p.is_published : true,
      is_pinned: !!p.is_pinned,
      sort_order: p.sort_order || 0,
      is_interactive: !!p.is_interactive,
      interactive_type: p.interactive_type || 'tendies'
    })

    if (form.price_cny !== '' && form.price_usd !== '') {
      isManualCny.value = true
      isManualUsd.value = true
    }

    if (p.interactive_config) {
      try {
        const cfg = typeof p.interactive_config === 'string' ? JSON.parse(p.interactive_config) : p.interactive_config
        Object.assign(form.interactive_config, cfg)
      } catch (e) {
        console.warn('Failed to parse interactive_config', e)
      }
    }

    // Populate downloads from any returned field
    const rawDownloads = p.downloads || res.downloads || p.download_links || res.download_links || []
    let dls = rawDownloads
    if (typeof dls === 'string') {
      try { dls = JSON.parse(dls) } catch { dls = [] }
    }
    if (Array.isArray(dls)) {
      dls.forEach(d => {
        if (d && d.platform) {
          form.downloads[d.platform] = {
            url: d.url || '',
            extract_code: d.extract_code || ''
          }
        }
      })
    }

    editor.value?.commands.setContent(p.description || '')
  } catch (err) {
    console.error(err)
  }
}

const availableVIPPlans = ref([
  { type: 'monthly', name_zh: '月度会员', name_en: 'Monthly VIP', icon: '🌙', price: 9.99, duration_days: 30, period: '月' },
  { type: 'quarterly', name_zh: '季度会员', name_en: 'Quarterly VIP', icon: '⭐', price: 24.99, duration_days: 90, period: '季', badge: '热门' },
  { type: 'yearly', name_zh: '年度会员', name_en: 'Yearly VIP', icon: '👑', price: 79.99, duration_days: 365, period: '年', badge: '省33%' },
  { type: 'lifetime', name_zh: '终生尊享', name_en: 'Lifetime VIP', icon: '💎', price: 199.99, duration_days: 36500, period: '终生', badge: '限时' }
])

function isTierChecked(type) {
  return Array.isArray(form.allowed_vip_levels) && form.allowed_vip_levels.includes(type)
}

function toggleTier(type) {
  if (!Array.isArray(form.allowed_vip_levels)) {
    form.allowed_vip_levels = []
  }
  const idx = form.allowed_vip_levels.indexOf(type)
  if (idx > -1) {
    form.allowed_vip_levels.splice(idx, 1)
  } else {
    form.allowed_vip_levels.push(type)
  }
}

function setVIPPreset(preset) {
  const all = availableVIPPlans.value.map(p => p.type)
  if (preset === 'all') {
    form.allowed_vip_levels = [...all]
  } else if (preset === 'quarterly_up') {
    form.allowed_vip_levels = all.filter(t => t !== 'monthly')
  } else if (preset === 'yearly_up') {
    form.allowed_vip_levels = all.filter(t => t !== 'monthly' && t !== 'quarterly')
  } else if (preset === 'lifetime_only') {
    form.allowed_vip_levels = ['lifetime']
  } else if (preset === 'none') {
    form.allowed_vip_levels = []
  }
}

function isPresetActive(preset) {
  const current = form.allowed_vip_levels || []
  const all = availableVIPPlans.value.map(p => p.type)
  if (preset === 'all') {
    return all.length > 0 && current.length === all.length && all.every(t => current.includes(t))
  }
  if (preset === 'quarterly_up') {
    const target = all.filter(t => t !== 'monthly')
    return current.length === target.length && target.every(t => current.includes(t))
  }
  if (preset === 'yearly_up') {
    const target = all.filter(t => t !== 'monthly' && t !== 'quarterly')
    return current.length === target.length && target.every(t => current.includes(t))
  }
  if (preset === 'lifetime_only') {
    return current.length === 1 && current[0] === 'lifetime'
  }
  return false
}

async function loadVIPPlans() {
  try {
    const cfg = await configAPI.getPublic()
    if (cfg?.membership_plans && cfg.membership_plans.length > 0) {
      availableVIPPlans.value = cfg.membership_plans.filter(p => p.enabled !== false)
      if (isNew.value) {
        form.allowed_vip_levels = availableVIPPlans.value.map(p => p.type)
      }
    }
  } catch (e) {
    console.warn('Using default VIP plans in admin:', e)
  }
}

onMounted(async () => {
  try {
    const res = await settingsAPI.getAll()
    if (res && res.exchange_rate) {
      exchangeRate.value = Number(res.exchange_rate)
    }
  } catch (e) {
    console.warn('Failed to load exchange rate in product edit:', e)
  }
  await loadVIPPlans()
  await loadCategories()
  await loadProduct()
})
onBeforeUnmount(() => { editor.value?.destroy() })
</script>

<style scoped>
.product-edit { padding-bottom: 80px; }

.edit-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
}
.edit-header h2 { font-size: 1.3rem; }

.edit-grid {
  display: grid;
  grid-template-columns: 1fr 360px;
  gap: 20px;
  align-items: start;
}

@media (max-width: 1024px) {
  .edit-grid { grid-template-columns: 1fr; }
}

.edit-main, .edit-sidebar {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.section-card {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  padding: 24px;
}

/* VIP Tier Download Privilege Configuration Styles */
.vip-tier-group {
  margin-top: 18px;
  padding: 16px;
  background: rgba(255, 215, 0, 0.03);
  border: 1px dashed rgba(255, 215, 0, 0.35);
  border-radius: var(--radius-md, 10px);
}

.vip-tier-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 14px;
}

.vip-title-col {
  display: flex;
  flex-direction: column;
}

.vip-tier-hint {
  font-size: 0.78rem;
  color: var(--text-muted, #8e8e93);
  margin-top: 2px;
}

.vip-preset-btns {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.preset-btn {
  padding: 4px 10px;
  font-size: 0.75rem;
  font-weight: 500;
  border-radius: 6px;
  background: var(--bg-hover, rgba(255, 255, 255, 0.06));
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.12));
  color: var(--text-color, #e0e0e0);
  cursor: pointer;
  transition: all 0.2s ease;
}

.preset-btn:hover {
  background: rgba(255, 215, 0, 0.12);
  border-color: rgba(255, 215, 0, 0.4);
  color: #ffd700;
}

.preset-btn.preset-active {
  background: rgba(255, 215, 0, 0.2);
  border-color: #ffd700;
  color: #ffd700;
  font-weight: 600;
}

.preset-btn.preset-danger:hover {
  background: rgba(255, 69, 58, 0.15);
  border-color: rgba(255, 69, 58, 0.5);
  color: #ff453a;
}

.vip-tiers-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 10px;
}

.vip-tier-card {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  border-radius: 8px;
  background: var(--bg-card, #1c1c1e);
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.1));
  cursor: pointer;
  user-select: none;
  transition: all 0.2s ease;
  position: relative;
}

.vip-tier-card:hover {
  border-color: rgba(255, 215, 0, 0.4);
  background: rgba(255, 215, 0, 0.04);
}

.vip-tier-card.active {
  border-color: #ffd700;
  background: rgba(255, 215, 0, 0.08);
  box-shadow: 0 0 12px rgba(255, 215, 0, 0.12);
}

.tier-checkbox {
  width: 18px;
  height: 18px;
  border-radius: 4px;
  border: 1.5px solid var(--border-color, rgba(255, 255, 255, 0.25));
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  background: rgba(0, 0, 0, 0.3);
  transition: all 0.18s ease;
}

.tier-checkbox.checked {
  background: #ffd700;
  border-color: #ffd700;
}

.tier-check-mark {
  font-size: 11px;
  font-weight: bold;
  color: #000;
  line-height: 1;
}

.tier-icon {
  font-size: 1.2rem;
  line-height: 1;
  flex-shrink: 0;
}

.tier-content {
  display: flex;
  flex-direction: column;
  min-width: 0;
  flex: 1;
}

.tier-name-row {
  display: flex;
  align-items: center;
  gap: 6px;
}

.tier-name {
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--text-color, #ffffff);
}

.tier-badge {
  font-size: 0.65rem;
  padding: 1px 5px;
  border-radius: 4px;
  background: linear-gradient(135deg, #ff7b00, #ff453a);
  color: #fff;
  font-weight: 600;
  transform: scale(0.9);
}

.tier-price {
  font-size: 0.72rem;
  color: var(--text-muted, #8e8e93);
  margin-top: 2px;
}

.section-title {
  font-size: 1.05rem;
  font-weight: 700;
  margin-bottom: 8px;
}

.section-subtitle {
  font-size: 0.8rem;
  color: var(--text-muted);
  margin-bottom: 16px;
}

.card-header-flex {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.form-group {
  margin-bottom: 16px;
}

.form-row {
  display: flex;
  gap: 16px;
}

.flex-1 { flex: 1; }
.flex-2 { flex: 2; }
.mt-2 { margin-top: 8px; }
.mt-3 { margin-top: 12px; }
.mt-4 { margin-top: 16px; }

.form-label {
  display: block;
  font-size: 0.85rem;
  font-weight: 600;
  margin-bottom: 6px;
  color: var(--text-secondary);
}

.form-error {
  color: #ef4444;
  font-size: 0.75rem;
  margin-top: 4px;
  display: block;
}

.tags-input {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  min-height: 42px;
}

.tag-remove {
  border: none;
  background: transparent;
  color: inherit;
  font-size: 0.75rem;
  margin-left: 4px;
  cursor: pointer;
}

.tag-text-input {
  border: none;
  background: transparent;
  outline: none;
  font-size: 0.85rem;
  flex: 1;
  min-width: 120px;
}

/* TipTap Editor */
.editor-wrapper {
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  overflow: hidden;
}

.editor-toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  padding: 8px 12px;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
}

.toolbar-btn {
  padding: 4px 8px;
  border-radius: var(--radius-sm);
  background: transparent;
  border: none;
  font-size: 0.85rem;
  cursor: pointer;
  color: var(--text-secondary);
  transition: all 0.15s;
}

.toolbar-btn:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.toolbar-btn.active {
  background: var(--primary);
  color: white;
}

.toolbar-divider {
  width: 1px;
  height: 20px;
  background: var(--border-color);
  margin: 0 4px;
}

.editor-content :deep(.tiptap) {
  padding: 16px;
  min-height: 180px;
  outline: none;
  font-size: 0.95rem;
  line-height: 1.6;
}

.editor-content :deep(.tiptap p.is-editor-empty:first-child::before) {
  color: var(--text-muted);
  content: attr(data-placeholder);
  float: left;
  height: 0;
  pointer-events: none;
}

/* Video platform styling */
.platform-tabs {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 8px;
}

.platform-tab {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  border-radius: 999px;
  border: 1px solid var(--border-color);
  background: var(--bg-secondary);
  font-size: 0.85rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.platform-tab:hover {
  border-color: var(--primary);
  background: var(--bg-hover);
}

.platform-tab.active {
  background: var(--text-primary);
  color: var(--bg-primary);
  border-color: var(--text-primary);
}

.input-with-action {
  display: flex;
  gap: 8px;
  align-items: center;
}

.test-link-btn {
  white-space: nowrap;
  font-size: 0.8rem;
  text-decoration: none;
  border: 1px solid var(--border-color);
}

.extra-video-row {
  display: grid;
  grid-template-columns: 140px 1fr 140px 36px;
  gap: 8px;
  align-items: center;
  margin-top: 8px;
}

/* Cover Image Upload */
.cover-mode-toggle {
  display: flex;
  background: var(--bg-secondary);
  padding: 2px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border-color);
}

.mode-btn {
  border: none;
  background: transparent;
  padding: 4px 10px;
  font-size: 0.75rem;
  border-radius: var(--radius-sm);
  cursor: pointer;
  color: var(--text-secondary);
}

.mode-btn.active {
  background: var(--bg-card);
  color: var(--text-primary);
  font-weight: 600;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

.hidden-file-input { display: none; }

.upload-dropzone {
  border: 2px dashed var(--border-color);
  border-radius: var(--radius-lg);
  padding: 24px 16px;
  text-align: center;
  cursor: pointer;
  background: var(--bg-secondary);
  transition: all 0.2s ease;
}

.upload-dropzone:hover, .upload-dropzone.is-dragover {
  border-color: var(--primary);
  background: rgba(99, 102, 241, 0.05);
}

.upload-icon {
  font-size: 2.2rem;
  margin-bottom: 8px;
}

.cover-preview-box {
  position: relative;
  border-radius: var(--radius-md);
  overflow: hidden;
  max-height: 200px;
}

.cover-preview-box img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.cover-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.2s;
  font-size: 0.85rem;
  font-weight: 600;
}

.cover-preview-box:hover .cover-overlay {
  opacity: 1;
}

.cover-meta-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 8px;
  font-size: 0.75rem;
}

.cover-path-preview {
  color: var(--text-muted);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 220px;
  font-family: monospace;
}

.drive-entry {
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
  padding: 14px;
  margin-bottom: 12px;
}

.drive-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  font-size: 0.88rem;
  margin-bottom: 10px;
}

.publish-toggle {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.toggle-label { font-weight: 600; font-size: 0.95rem; }
.toggle-desc { font-size: 0.75rem; color: var(--text-muted); }

.toggle-switch {
  width: 48px;
  height: 26px;
  background: var(--border-color);
  border-radius: 999px;
  position: relative;
  cursor: pointer;
  transition: background 0.2s;
}

.toggle-switch.on { background: #10b981; }

.toggle-thumb {
  width: 20px;
  height: 20px;
  background: white;
  border-radius: 50%;
  position: absolute;
  top: 3px;
  left: 3px;
  transition: transform 0.2s;
}

.toggle-switch.on .toggle-thumb {
  transform: translateX(22px);
}

.edit-footer {
  display: flex;
  justify-content: flex-end;
  gap: 16px;
  margin-top: 24px;
}

.text-danger { color: #ef4444; }

.tiktok-detected-box {
  display: flex;
  align-items: center;
  gap: 10px;
  background: #111;
  color: #fff;
  padding: 10px 14px;
  border-radius: 12px;
  border: 1px solid rgba(255, 45, 85, 0.4);
  font-size: 0.82rem;
}

.detect-badge {
  background: #ff2d55;
  color: #fff;
  padding: 2px 8px;
  border-radius: 999px;
  font-weight: 700;
  font-size: 0.72rem;
  white-space: nowrap;
}

.detect-text {
  color: #ddd;
}

/* Interactive Config Card & Simulation Setup */
.interactive-config-card {
  border: 1px solid rgba(139, 92, 246, 0.3);
  position: relative;
  overflow: hidden;
}

.interactive-panel {
  margin-top: 18px;
  padding-top: 16px;
  border-top: 1px dashed var(--border-color);
}

.interactive-type-pills {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 8px;
}

.type-pill-btn {
  padding: 8px 16px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 999px;
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--text-color);
  cursor: pointer;
  transition: all 0.2s ease;
}

.type-pill-btn:hover {
  border-color: #8b5cf6;
  background: rgba(139, 92, 246, 0.08);
}

.type-pill-btn.active {
  background: linear-gradient(135deg, #8b5cf6, #6366f1);
  color: #ffffff;
  border-color: transparent;
  box-shadow: 0 4px 12px rgba(139, 92, 246, 0.35);
}

.assets-upload-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 14px;
  margin-top: 14px;
}

.asset-slot-card {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  padding: 14px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.asset-slot-card.highlight-slot {
  border-color: rgba(245, 158, 11, 0.4);
  background: rgba(245, 158, 11, 0.03);
}

.slot-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.slot-badge {
  font-size: 0.78rem;
  font-weight: 700;
  padding: 3px 8px;
  border-radius: 6px;
}

.badge-purple { background: rgba(139, 92, 246, 0.15); color: #8b5cf6; }
.badge-blue { background: rgba(59, 130, 246, 0.15); color: #3b82f6; }
.badge-gold { background: rgba(245, 158, 11, 0.18); color: #f59e0b; }
.badge-cyan { background: rgba(6, 182, 212, 0.15); color: #06b6d4; }

.slot-hint {
  font-size: 0.72rem;
  color: var(--text-muted);
}

.slot-input-wrap {
  display: flex;
  gap: 8px;
  align-items: center;
}

.form-input-sm {
  padding: 6px 10px;
  font-size: 0.8rem;
  flex: 1;
}

.file-upload-btn {
  position: relative;
  overflow: hidden;
  white-space: nowrap;
  font-size: 0.75rem;
  padding: 6px 10px;
  cursor: pointer;
}

.btn-gold-outline {
  border-color: #f59e0b;
  color: #f59e0b;
}

.btn-gold-outline:hover {
  background: #f59e0b;
  color: #fff;
}

.hidden-file-input {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  opacity: 0;
  cursor: pointer;
}

.slot-tip {
  font-size: 0.72rem;
  color: var(--text-muted);
  line-height: 1.3;
}

.interactive-controls-row {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 20px;
  padding: 12px 16px;
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
}

.control-toggle-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.82rem;
}

.control-slider-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
  font-size: 0.82rem;
  flex: 1 1 360px;
  min-width: 280px;
  background: var(--bg-card);
  padding: 10px 14px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border-color);
  transition: opacity 0.2s ease;
}

.control-slider-item.control-disabled {
  opacity: 0.55;
  pointer-events: none;
}

.slider-label-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.depth-val-display {
  color: #8b5cf6;
  font-weight: 700;
  margin-left: 4px;
}

.depth-val-display.zero-val {
  color: var(--text-muted);
}

.depth-val-display.neg-val {
  color: #06b6d4;
}

.depth-num-box {
  display: flex;
  align-items: center;
  gap: 4px;
}

.depth-num-input {
  width: 58px;
  text-align: center;
  padding: 2px 4px;
  font-size: 0.82rem;
  font-weight: 600;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background: var(--bg-primary);
  color: var(--text-primary);
}

.unit-text {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.slider-input-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.range-bound {
  font-size: 0.72rem;
  color: var(--text-muted);
  font-family: monospace;
}

.depth-presets {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
  margin-top: 2px;
}

.preset-label {
  font-size: 0.72rem;
  color: var(--text-muted);
}

.preset-chip {
  border: 1px solid var(--border-color);
  background: var(--bg-secondary);
  color: var(--text-secondary);
  font-size: 0.72rem;
  padding: 2px 7px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.preset-chip:hover {
  border-color: #8b5cf6;
  color: #8b5cf6;
  background: rgba(139, 92, 246, 0.08);
}

.preset-chip.active {
  background: #8b5cf6;
  color: #fff;
  border-color: #8b5cf6;
}

.status-badge {
  font-size: 0.72rem;
  padding: 2px 6px;
  border-radius: 4px;
  font-weight: 600;
}
.status-badge.badge-on {
  background: rgba(16, 185, 129, 0.15);
  color: #10b981;
}
.status-badge.badge-off {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
}

.tag-status-pill {
  font-size: 0.75rem;
  padding: 1px 6px;
  border-radius: 4px;
  margin-left: 6px;
  font-weight: normal;
}
.tag-status-pill.on {
  background: rgba(139, 92, 246, 0.15);
  color: #8b5cf6;
}
.tag-status-pill.off {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
}
.tag-status-pill.neutral {
  background: rgba(107, 114, 128, 0.15);
  color: #9ca3af;
}
.tag-status-pill.neg {
  background: rgba(6, 182, 212, 0.15);
  color: #06b6d4;
}

.range-slider {
  flex: 1;
  accent-color: #8b5cf6;
  cursor: pointer;
}

.admin-simulator-preview {
  background: var(--bg-secondary);
  border-radius: var(--radius-md);
  padding: 16px;
  border: 1px solid var(--border-color);
}

.preview-header-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 14px;
}

.preview-tag {
  font-size: 0.82rem;
  font-weight: 700;
  color: var(--text-color);
}

.preview-status {
  font-size: 0.75rem;
  color: #10b981;
  font-weight: 600;
}

.preview-status.warning {
  color: #f59e0b;
}

.mini-iphone-mockup {
  position: relative;
  width: 216px;
  height: 468px;
  aspect-ratio: 393 / 852;
  margin: 0 auto;
  border-radius: 38px;
  background: #000;
  border: 6px solid #2d2d2f;
  box-shadow: 0 16px 36px rgba(0, 0, 0, 0.4), inset 0 0 0 1px rgba(255, 255, 255, 0.15);
  overflow: hidden;
  cursor: crosshair;
  user-select: none;
}

.mini-dynamic-island {
  position: absolute;
  top: 7px;
  left: 50%;
  transform: translateX(-50%);
  width: 68px;
  height: 20px;
  background: #000;
  border-radius: 999px;
  z-index: 60;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding-right: 6px;
  pointer-events: none;
}

.mini-camera-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #1e1b4b;
}

.mini-layer {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  pointer-events: none;
  transition: transform 0.08s ease-out;
}

.mini-caml-layer {
  max-width: none !important;
  max-height: none !important;
  object-fit: fill !important;
}

.mini-anim-pos,
.mini-anim-rot {
  width: 100%;
  height: 100%;
  pointer-events: none;
  transform-origin: center center;
}

.mini-sprite-inner-img {
  width: 100%;
  height: 100%;
  display: block;
  pointer-events: none;
  object-fit: fill;
}

.mini-bg {
  z-index: 5;
}

.mini-clock {
  position: absolute;
  top: 42px;
  left: 0;
  width: 100%;
  text-align: center;
  z-index: 25;
  pointer-events: none;
  transition: transform 0.08s ease-out;
}

.mini-date-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 3px;
  margin-bottom: 2px;
}

.mini-lock-icon {
  font-size: 10px;
  color: rgba(255, 255, 255, 0.85);
}

.mini-date-text {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.9);
  font-weight: 500;
  text-shadow: 0 1px 4px rgba(0, 0, 0, 0.4);
}

.mini-time {
  font-family: -apple-system, BlinkMacSystemFont, "SF Pro Display", sans-serif;
  font-size: 52px;
  font-weight: 700;
  line-height: 1;
  letter-spacing: -2px;
  color: rgba(255, 255, 255, 0.95);
  text-shadow: 0 2px 10px rgba(0, 0, 0, 0.35);
}

.mini-fg {
  z-index: 35;
}

.mini-island {
  position: absolute;
  top: 10px;
  left: 50%;
  transform: translateX(-50%);
  width: 72px;
  height: 20px;
  background: #000;
  border-radius: 20px;
  z-index: 10;
}

.mini-watermark {
  position: absolute;
  bottom: 30px;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(0, 0, 0, 0.45);
  backdrop-filter: blur(8px);
  color: rgba(255, 255, 255, 0.7);
  padding: 4px 10px;
  border-radius: 999px;
  font-size: 9px;
  letter-spacing: 0.5px;
  white-space: nowrap;
  z-index: 5;
}

.toggle-switch.switch-pinned.on {
  background: linear-gradient(135deg, #ef4444, #dc2626) !important;
}

.sort-order-block input {
  font-weight: 700;
}
</style>

