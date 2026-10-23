<template>
  <div class="admin-settings">
    <!-- Header -->
    <div class="settings-header">
      <div>
        <h2 class="settings-title">⚙️ {{ locale === 'zh' ? '商城核心系统配置' : 'Shop System Settings' }}</h2>
        <p class="settings-subtitle">
          {{ locale === 'zh'
            ? '自定义全站会员定价体系、免费壁纸下载权限规则，以及 PayPal 支付网关'
            : 'Customize VIP membership pricing tiers, free download permission rules, and PayPal gateway' }}
        </p>
      </div>
      <button class="btn btn-primary save-btn" :disabled="saving" @click="saveAllSettings">
        <span v-if="saving" class="spinner"></span>
        <span v-else>💾 {{ locale === 'zh' ? '保存所有系统配置' : 'Save Changes' }}</span>
      </button>
    </div>

    <!-- Main Navigation Tabs -->
    <div class="settings-nav-tabs">
      <button
        class="nav-tab-btn"
        :class="{ active: currentTab === 'pricing' }"
        @click="currentTab = 'pricing'"
      >
        <span class="nav-tab-icon">👑</span>
        <div class="nav-tab-info">
          <strong>{{ locale === 'zh' ? '会员定价系统' : 'Membership Pricing' }}</strong>
          <small>{{ locale === 'zh' ? '月度/季度/年度/终身定价与促销' : 'Tiers, Prices & Badges' }}</small>
        </div>
      </button>

      <button
        class="nav-tab-btn"
        :class="{ active: currentTab === 'permissions' }"
        @click="currentTab = 'permissions'"
      >
        <span class="nav-tab-icon">📥</span>
        <div class="nav-tab-info">
          <strong>{{ locale === 'zh' ? '免费下载与权限规则' : 'Free Download & Permissions' }}</strong>
          <small>{{ locale === 'zh' ? '免登录直下、每日限额、全局卡密' : 'Access rules, limits & codes' }}</small>
        </div>
      </button>

      <button
        class="nav-tab-btn"
        :class="{ active: currentTab === 'payment' }"
        @click="currentTab = 'payment'"
      >
        <span class="nav-tab-icon">💳</span>
        <div class="nav-tab-info">
          <strong>{{ locale === 'zh' ? 'PayPal 支付网关' : 'PayPal Gateway' }}</strong>
          <small>{{ locale === 'zh' ? '个人卖家收款公私钥与沙盒调试' : 'API credentials & Sandbox' }}</small>
        </div>
      </button>

      <button
        class="nav-tab-btn"
        :class="{ active: currentTab === 'channels' }"
        @click="currentTab = 'channels'"
      >
        <span class="nav-tab-icon">💱</span>
        <div class="nav-tab-info">
          <strong>{{ locale === 'zh' ? '支付渠道与汇率' : 'Channels & Currency' }}</strong>
          <small>{{ locale === 'zh' ? '支付宝/微信/PayPal范围与汇率' : 'Alipay, WeChat, PayPal & Rates' }}</small>
        </div>
      </button>
    </div>

    <!-- TAB 1: 会员定价系统 -->
    <div v-show="currentTab === 'pricing'" class="tab-pane-content">
      <div class="card settings-card">
        <div class="card-head">
          <div class="head-icon icon-gold">👑</div>
          <div class="head-title-wrap">
            <div class="head-flex">
              <div>
                <h3>{{ locale === 'zh' ? 'VIP 会员方案与动态定价' : 'VIP Membership Pricing Tiers' }}</h3>
                <p>{{ locale === 'zh' ? '后台设置的价格将实时同步到前台会员中心和收银台，支付捕获将严格按照配置金额扣款。' : 'Configured prices immediately take effect on frontend and PayPal orders.' }}</p>
              </div>
              <div class="head-actions">
                <button type="button" class="btn btn-outline btn-sm" @click="resetToDefaultPlans">
                  🔄 {{ locale === 'zh' ? '恢复默认定价' : 'Reset Defaults' }}
                </button>
                <button type="button" class="btn btn-primary btn-sm" @click="addNewPlan">
                  ➕ {{ locale === 'zh' ? '添加自定义方案' : 'Add New Tier' }}
                </button>
              </div>
            </div>
          </div>
        </div>

        <div class="card-body">
          <div class="pricing-plans-stack">
            <div
              v-for="(plan, index) in membershipPlans"
              :key="plan.type + index"
              class="plan-edit-card"
              :class="{ disabled: !plan.enabled, featured: plan.featured }"
            >
              <!-- Card Header Row -->
              <div class="plan-card-topbar">
                <div class="plan-card-left">
                  <span class="plan-icon-preview">{{ plan.icon || '👑' }}</span>
                  <div class="plan-identity">
                    <span class="plan-type-pill">{{ plan.type }}</span>
                    <strong class="plan-name-preview">{{ plan.name_zh }} ({{ plan.name_en }})</strong>
                  </div>
                </div>

                <div class="plan-card-toggles">
                  <label class="switch-control">
                    <input v-model="plan.featured" type="checkbox" />
                    <span class="switch-label">
                      {{ plan.featured ? '⭐ ' + (locale === 'zh' ? '推荐高亮' : 'Featured') : (locale === 'zh' ? '设为推荐' : 'Normal') }}
                    </span>
                  </label>

                  <label class="switch-control ml-2">
                    <input v-model="plan.enabled" type="checkbox" />
                    <span class="switch-label" :class="{ 'text-danger': !plan.enabled }">
                      {{ plan.enabled ? '🟢 ' + (locale === 'zh' ? '已上架' : 'Active') : '⚪ ' + (locale === 'zh' ? '已下架' : 'Disabled') }}
                    </span>
                  </label>

                  <button
                    v-if="membershipPlans.length > 1"
                    type="button"
                    class="btn-delete-plan"
                    title="删除此方案"
                    @click="removePlan(index)"
                  >
                    🗑️
                  </button>
                </div>
              </div>

              <!-- Form Inputs Grid -->
              <div class="plan-form-grid">
                <div class="form-group">
                  <label class="form-label">{{ locale === 'zh' ? '中文方案名' : 'Name (ZH)' }}</label>
                  <input v-model="plan.name_zh" type="text" class="form-input" placeholder="如: 月度会员" />
                </div>

                <div class="form-group">
                  <label class="form-label">{{ locale === 'zh' ? '英文方案名' : 'Name (EN)' }}</label>
                  <input v-model="plan.name_en" type="text" class="form-input" placeholder="e.g. Monthly VIP" />
                </div>

                <div class="form-group">
                  <label class="form-label">
                    {{ locale === 'zh' ? '销售价格 (USD $)' : 'Price (USD)' }}
                  </label>
                  <div class="input-with-prefix">
                    <span class="input-prefix">$</span>
                    <input
                      v-model.number="plan.price"
                      type="number"
                      step="0.01"
                      min="0.01"
                      class="form-input font-bold"
                      placeholder="9.99"
                    />
                  </div>
                </div>

                <div class="form-group">
                  <label class="form-label">{{ locale === 'zh' ? '有效天数' : 'Duration (Days)' }}</label>
                  <input
                    v-model.number="plan.duration_days"
                    type="number"
                    min="1"
                    class="form-input"
                    placeholder="30"
                  />
                  <small class="form-tip">
                    {{ plan.duration_days >= 3650 ? '永久/终身' : (plan.duration_days + ' 天') }}
                  </small>
                </div>

                <div class="form-group">
                  <label class="form-label">{{ locale === 'zh' ? '方案图标 (Emoji)' : 'Icon' }}</label>
                  <input v-model="plan.icon" type="text" class="form-input" placeholder="🌙 / ⭐ / 👑 / 💎" />
                </div>

                <div class="form-group">
                  <label class="form-label">{{ locale === 'zh' ? '促销/折扣角标' : 'Promo Badge' }}</label>
                  <input v-model="plan.badge" type="text" class="form-input" placeholder="如: 热销推荐 / 省17%" />
                </div>
              </div>
            </div>
          </div>

          <!-- Helper tip -->
          <div class="alert-box alert-info mt-4">
            <strong>💡 {{ locale === 'zh' ? '会员权益自动赋权机制' : 'VIP Automation' }}</strong>
            <p>
              {{ locale === 'zh'
                ? '当用户在前台通过 PayPal 付款或测试开通任意会员方案后，系统将自动根据该方案的“有效天数”计算并在用户账户上累加会员到期时间；同时立即解锁全场壁纸下载权限！'
                : 'When a user subscribes, their expiration is extended by the configured duration and all wallpapers unlock immediately.' }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- TAB 2: 免费下载与权限规则 -->
    <div v-show="currentTab === 'permissions'" class="tab-pane-content">
      <div class="settings-grid">
        <!-- 免费壁纸下载门槛 -->
        <div class="card settings-card">
          <div class="card-head">
            <div class="head-icon icon-green">📥</div>
            <div>
              <h3>{{ locale === 'zh' ? '免费壁纸访问门槛与下载权限' : 'Free Wallpaper Access & Permissions' }}</h3>
              <p>{{ locale === 'zh' ? '灵活控制价格为 $0 的免费壁纸是否需要注册登录后才能获取下载链接' : 'Control whether free wallpapers (price = $0) require user login' }}</p>
            </div>
          </div>

          <div class="card-body">
            <div class="form-group">
              <label class="form-label font-bold">
                {{ locale === 'zh' ? '免费壁纸下载门槛模式' : 'Free Wallpaper Download Requirement' }}
              </label>
              <div class="radio-cards">
                <label
                  class="radio-card"
                  :class="{ selected: !permissions.free_download_require_login }"
                >
                  <input
                    v-model="permissions.free_download_require_login"
                    type="radio"
                    :value="false"
                  />
                  <div class="radio-info">
                    <strong>🚀 {{ locale === 'zh' ? '免登录全开放 (访客直接下载)' : 'Open to All (No Login Required)' }}</strong>
                    <span>
                      {{ locale === 'zh'
                        ? '任意访客进入免费壁纸详情页，无需登录或注册账号，即可直接展示百度网盘、夸克网盘和谷歌网盘下载按钮及提取码。体验最好，适合快速引流！'
                        : 'Visitors can view and download free wallpapers directly without creating an account.' }}
                    </span>
                  </div>
                </label>

                <label
                  class="radio-card"
                  :class="{ selected: permissions.free_download_require_login }"
                >
                  <input
                    v-model="permissions.free_download_require_login"
                    type="radio"
                    :value="true"
                  />
                  <div class="radio-info">
                    <strong>🔒 {{ locale === 'zh' ? '强制登录后下载 (引导用户注册)' : 'Require Login to Download' }}</strong>
                    <span>
                      {{ locale === 'zh'
                        ? '访客点击免费壁纸必须先登录/注册账号，未登录访客将弹出快速登录注册弹窗。有效促进新用户注册沉淀与私域积累！'
                        : 'Visitors must log in or register before cloud download links and passcodes are revealed.' }}
                    </span>
                  </div>
                </label>
              </div>
            </div>

            <!-- 每日限额 -->
            <div class="form-group mt-4">
              <label class="form-label font-bold">
                {{ locale === 'zh' ? '普通注册会员每日免费下载限额' : 'Daily Free Download Limit for Regular Users' }}
              </label>
              <div style="max-width: 320px;">
                <div class="input-with-suffix">
                  <input
                    v-model.number="permissions.daily_free_limit"
                    type="number"
                    min="0"
                    max="999"
                    class="form-input font-bold"
                  />
                  <span class="input-suffix">{{ locale === 'zh' ? '次 / 天' : 'per day' }}</span>
                </div>
              </div>
              <small class="form-tip">
                {{ locale === 'zh'
                  ? '填写 0 表示无限制；若设置为如 5 或 10，则普通注册用户每天限免费下载对应次数，超出后提示升级 VIP 会员，VIP 尊贵会员不受任何限制。'
                  : 'Set to 0 for unlimited downloads. VIP members always have unlimited access.' }}
              </small>
            </div>
          </div>
        </div>

        <!-- 游客免登录购买 & 万能卡密 -->
        <div class="card settings-card">
          <div class="card-head">
            <div class="head-icon icon-amber">⚡</div>
            <div>
              <h3>{{ locale === 'zh' ? '游客临时购买与全局兑换口令' : 'Guest Checkout & Master Redeem Codes' }}</h3>
              <p>{{ locale === 'zh' ? '控制收费壁纸的游客免登录临时购买功能，以及作者万能卡密' : 'Configure 1-click guest purchase channel and master unlock codes' }}</p>
            </div>
          </div>

          <div class="card-body">
            <!-- 游客临时免登录购买开关 -->
            <div class="toggle-setting-row">
              <div>
                <strong>{{ locale === 'zh' ? '允许游客临时免登录购买' : 'Enable Guest Direct Checkout' }}</strong>
                <p>
                  {{ locale === 'zh'
                    ? '开启后，访客购买单张付费壁纸时可选择“⚡ 临时购买 (免登录)”，填入邮箱直接完成付款并下载；若关闭，所有访客必须注册账号后方可付款。'
                    : 'Allow guests to purchase single wallpapers without registering an account.' }}
                </p>
              </div>
              <label class="switch-control lg">
                <input v-model="permissions.guest_buy_enabled" type="checkbox" />
                <span class="switch-slider"></span>
              </label>
            </div>

            <hr class="section-divider" />

            <!-- 快速模拟测试支付通道开关 -->
            <div class="toggle-setting-row">
              <div>
                <strong>{{ locale === 'zh' ? '开启快速模拟测试支付通道' : 'Enable Fast Sandbox Test Payment' }}</strong>
                <p>
                  {{ locale === 'zh'
                    ? '方便您在未配置 PayPal 或开发调试时，无需真实扣款即可秒级模拟购买成功并解锁。上线营业时建议保持开启或按需切换。'
                    : '1-click simulated checkout for effortless testing of products and memberships.' }}
                </p>
              </div>
              <label class="switch-control lg">
                <input v-model="permissions.allow_test_pay" type="checkbox" />
                <span class="switch-slider"></span>
              </label>
            </div>

            <hr class="section-divider" />

            <!-- 全局万能口令 -->
            <div class="form-group">
              <label class="form-label font-bold">
                🔑 {{ locale === 'zh' ? '全局通用作者卡密 / 万能解锁口令' : 'Universal Author Master Unlock Codes' }}
              </label>
              <input
                v-model="permissions.master_unlock_codes"
                type="text"
                class="form-input font-mono"
                placeholder="例如: loveyou, 666, 1194, VIP888, VIP666"
              />
              <small class="form-tip">
                {{ locale === 'zh'
                  ? '多个口令用英文逗号 (,) 分隔。任何人在前台收银台【卡密兑换】中输入这些口令中的任意一个，均可立即免支付直接兑换解锁对应壁纸！'
                  : 'Comma-separated codes. Entering any of these codes instantly unlocks any wallpaper.' }}
              </small>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- TAB 3: PayPal 结算网关 -->
    <div v-show="currentTab === 'payment'" class="tab-pane-content">
      <div class="card settings-card">
        <div class="card-head">
          <div class="head-icon icon-blue">💳</div>
          <div>
            <h3>PayPal {{ locale === 'zh' ? '个人卖家结算网关配置' : 'Payment Gateway Setup' }}</h3>
            <p>{{ locale === 'zh' ? '接入 PayPal REST API v2，支持个人卖家与商业账号接入，支持单品与会员订阅收款' : 'Accept one-time wallpaper payments and membership subscriptions' }}</p>
          </div>
        </div>

        <div class="card-body">
          <div class="form-group">
            <label class="form-label font-bold">PayPal Client ID (公钥)</label>
            <input
              v-model="paymentSettings.paypal_client_id"
              type="text"
              class="form-input font-mono"
              placeholder="例如: sb 或 真实 Client ID (AXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX)"
            />
            <small class="form-tip">
              {{ locale === 'zh' ? '在 PayPal Developer 控制台 (developer.paypal.com) 创建 App 后获取。本地沙盒测试可填 sb。' : 'App Client ID from developer.paypal.com' }}
            </small>
          </div>

          <div class="form-group">
            <label class="form-label font-bold">PayPal Client Secret (私钥)</label>
            <input
              v-model="paymentSettings.paypal_client_secret"
              type="password"
              class="form-input font-mono"
              placeholder="例如: EXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
            />
            <small class="form-tip">
              {{ locale === 'zh' ? '用于服务端安全验证捕获交易并防篡改。留空或显示星号表示保持已有私钥不变。' : 'Backend secret to verify captures. Keep masked to leave unchanged.' }}
            </small>
          </div>

          <div class="form-group">
            <label class="form-label font-bold">{{ locale === 'zh' ? '运行模式 (Environment Mode)' : 'Mode' }}</label>
            <div class="radio-cards">
              <label class="radio-card" :class="{ selected: paymentSettings.paypal_mode === 'sandbox' }">
                <input v-model="paymentSettings.paypal_mode" type="radio" value="sandbox" />
                <div class="radio-info">
                  <strong>🧪 沙盒测试模式 (Sandbox)</strong>
                  <span>{{ locale === 'zh' ? '使用 PayPal 提供的虚拟测试账户与资金进行结算，适合调试全套流程' : 'Simulated transactions with sandbox accounts' }}</span>
                </div>
              </label>
              <label class="radio-card" :class="{ selected: paymentSettings.paypal_mode === 'live' }">
                <input v-model="paymentSettings.paypal_mode" type="radio" value="live" />
                <div class="radio-info">
                  <strong>🚀 正式收款模式 (Live Production)</strong>
                  <span>{{ locale === 'zh' ? '真实扣除买家信用卡/PayPal资金，并实时结算入账至您的 PayPal 个人/企业账户' : 'Process real payments into your PayPal account' }}</span>
                </div>
              </label>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- TAB 4: 支付渠道与双轨汇率管理 -->
    <div v-show="currentTab === 'channels'" class="tab-pane-content">
      <!-- 汇率设置卡片 -->
      <div class="card settings-card">
        <div class="card-head">
          <div class="head-icon icon-gold">💱</div>
          <div>
            <h3>{{ locale === 'zh' ? '双轨币种与汇率设置' : 'Dual-Currency & Exchange Rate' }}</h3>
            <p>{{ locale === 'zh' ? '用于前端针对不同地域访客自动呈现对应币种，以及发布壁纸时单币种自动折算另一币种' : 'Controls auto price conversion between USD and CNY' }}</p>
          </div>
        </div>
        <div class="card-body">
          <div class="form-group" style="max-width: 480px;">
            <label class="form-label font-bold">
              {{ locale === 'zh' ? '美元兑人民币汇率 (1 USD = ? CNY)' : 'Exchange Rate (1 USD = ? CNY)' }}
            </label>
            <div class="input-with-suffix">
              <input
                v-model.number="exchangeRate"
                type="number"
                step="0.01"
                min="0.1"
                class="form-input font-bold"
                placeholder="7.25"
              />
              <span class="input-suffix">CNY</span>
            </div>
            <small class="form-tip">
              {{ locale === 'zh'
                ? '例如填入 7.25，则 1 美元 = 7.25 元人民币。若商品只录入人民币 7.25，系统自动折算为 1.00 美元；若只录入 3.99 美元，自动折算为 28.93 元人民币；若双币种均填写，则以填写具体数字为准。'
                : 'Auto converts between USD and CNY when only one currency is filled in.' }}
            </small>
          </div>
        </div>
      </div>

      <!-- 支付渠道配置卡片 -->
      <div class="card settings-card">
        <div class="card-head">
          <div class="head-icon icon-emerald">💳</div>
          <div class="head-title-wrap">
            <div class="head-flex">
              <div>
                <h3>{{ locale === 'zh' ? '多渠道支付通道管理' : 'Payment Channels Management' }}</h3>
                <p>{{ locale === 'zh' ? '配置支持的支付渠道（如支付宝、微信支付、PayPal等），并指定其适用的地域范围（全球/仅中国/仅境外）' : 'Configure payment methods and regional availability (Global, China only, Overseas only)' }}</p>
              </div>
              <button type="button" class="btn btn-primary btn-sm" @click="openAddChannelModal">
                ➕ {{ locale === 'zh' ? '添加支付渠道' : 'Add Channel' }}
              </button>
            </div>
          </div>
        </div>

        <div class="card-body">
          <div class="channels-table-wrap">
            <table class="data-table">
              <thead>
                <tr>
                  <th>{{ locale === 'zh' ? '渠道名称' : 'Channel Name' }}</th>
                  <th>{{ locale === 'zh' ? '标识代码' : 'Code' }}</th>
                  <th>{{ locale === 'zh' ? '适用地域范围' : 'Regional Scope' }}</th>
                  <th>{{ locale === 'zh' ? '渠道描述' : 'Description' }}</th>
                  <th>{{ locale === 'zh' ? '启用状态' : 'Status' }}</th>
                  <th>{{ locale === 'zh' ? '操作' : 'Actions' }}</th>
                </tr>
              </thead>
              <tbody>
                <tr v-if="loadingChannels">
                  <td colspan="6" class="text-center py-6">
                    <span class="spinner"></span>
                    <span class="ml-2 text-secondary">{{ locale === 'zh' ? '正在加载支付渠道...' : 'Loading channels...' }}</span>
                  </td>
                </tr>
                <tr v-else-if="channels.length === 0">
                  <td colspan="6" class="text-center py-6 text-secondary">
                    {{ locale === 'zh' ? '暂未添加任何支付渠道' : 'No payment channels configured' }}
                  </td>
                </tr>
                <tr v-else v-for="ch in channels" :key="ch.id">
                  <td>
                    <div style="display:flex; align-items:center; gap:10px;">
                      <img v-if="ch.code !== 'paypal' && ch.code !== 'test' && ch.qr_code_url" :src="ch.qr_code_url" class="channel-qr-thumb" :title="locale === 'zh' ? '已配置收款码图片' : 'QR Code Configured'" />
                      <span v-else style="font-size:20px;">{{ ch.icon || (ch.code === 'paypal' ? '💳' : '⚡') }}</span>
                      <div>
                        <strong class="font-bold">{{ ch.name }}</strong>
                        <span v-if="ch.code === 'paypal'" class="paypal-gateway-badge">💳 {{ locale === 'zh' ? '官方API唤醒' : 'API Gateway' }}</span>
                        <span v-else-if="ch.code === 'test'" class="test-gateway-badge">⚡ {{ locale === 'zh' ? '沙盒免密' : 'Sandbox Demo' }}</span>
                        <span v-else-if="ch.qr_code_url" class="qr-has-badge">✓ {{ locale === 'zh' ? '已配收款码' : 'QR Ready' }}</span>
                        <span v-else class="qr-empty-badge">⚠️ {{ locale === 'zh' ? '未配收款码' : 'No QR' }}</span>
                      </div>
                    </div>
                  </td>
                  <td>
                    <span class="font-mono text-secondary">{{ ch.code }}</span>
                  </td>
                  <td>
                    <span class="scope-pill" :class="getScopeClass(ch.scope)">
                      {{ formatScope(ch.scope) }}
                    </span>
                  </td>
                  <td class="text-secondary text-sm" style="max-width: 260px;">
                    {{ ch.description || '-' }}
                  </td>
                  <td>
                    <label class="switch-control">
                      <input type="checkbox" v-model="ch.enabled" @change="toggleChannel(ch)" />
                      <span class="switch-slider"></span>
                    </label>
                  </td>
                  <td>
                    <div style="display:flex; gap:6px;">
                      <button class="btn btn-sm btn-ghost" @click="openEditChannelModal(ch)">
                        ✏️ {{ locale === 'zh' ? '编辑' : 'Edit' }}
                      </button>
                      <button class="btn btn-sm btn-ghost text-danger" @click="removeChannel(ch)">
                        🗑️
                      </button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>

    <!-- Payment Channel Modal -->
    <div v-if="showChannelModal" class="modal-overlay" @click.self="showChannelModal = false">
      <div class="modal-box">
        <h3>{{ isEditingChannel ? (locale === 'zh' ? '编辑支付渠道' : 'Edit Payment Channel') : (locale === 'zh' ? '添加新支付渠道' : 'Add Payment Channel') }}</h3>
        <p class="user-sub">{{ locale === 'zh' ? '设置支付方式名称、标识、展示图标与根据访客 IP 的适用范围' : 'Configure payment details and region targeting' }}</p>

        <div class="form-group">
          <label class="form-label font-bold">{{ locale === 'zh' ? '渠道名称' : 'Channel Name' }} *</label>
          <input v-model="channelForm.name" type="text" class="form-input" placeholder="例如: 支付宝 / 微信支付 / Stripe" />
        </div>

        <div class="form-row">
          <div class="form-group">
            <label class="form-label font-bold">{{ locale === 'zh' ? '渠道代码 (Code)' : 'Code' }} *</label>
            <select v-model="channelForm.code" class="form-input form-select">
              <option value="alipay">alipay (支付宝)</option>
              <option value="wechat">wechat (微信支付)</option>
              <option value="paypal">paypal (PayPal 国际网关)</option>
              <option value="stripe">stripe (Stripe 信用卡)</option>
              <option value="test">test (极速沙盒免密)</option>
              <option value="custom">custom (自定义收款码)</option>
            </select>
          </div>

          <div class="form-group">
            <label class="form-label font-bold">{{ locale === 'zh' ? '适用地域范围' : 'Regional Scope' }} *</label>
            <select v-model="channelForm.scope" class="form-input form-select font-bold">
              <option value="global">🌍 {{ locale === 'zh' ? '全球通用 (所有人可见)' : 'Global (All Visitors)' }}</option>
              <option value="china_only">🇨🇳 {{ locale === 'zh' ? '仅限中国 (国内/华语IP可见)' : 'China Only (Domestic IP)' }}</option>
              <option value="overseas_only">🌐 {{ locale === 'zh' ? '仅限境外 (海外/国际IP可见)' : 'Overseas Only (International IP)' }}</option>
            </select>
          </div>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label class="form-label font-bold">{{ locale === 'zh' ? '图标 Emoji' : 'Icon Emoji' }}</label>
            <input v-model="channelForm.icon" type="text" class="form-input" placeholder="例如: 🔵, 🟢, 💳, ⚡" />
          </div>

          <div class="form-group">
            <label class="form-label font-bold">{{ locale === 'zh' ? '排序权重' : 'Sort Order' }}</label>
            <input v-model.number="channelForm.sort_order" type="number" class="form-input" placeholder="0" />
          </div>
        </div>

        <div class="form-group">
          <label class="form-label font-bold">{{ locale === 'zh' ? '展示副标题/描述' : 'Description' }}</label>
          <input v-model="channelForm.description" type="text" class="form-input" placeholder="例如: 支持微信 / 支付宝扫码快捷支付" />
        </div>

        <!-- 支付模式分流：PayPal 网关提示 / 沙盒提示 / 收款码上传 -->
        <div v-if="channelForm.code === 'paypal'" class="paypal-notice-card">
          <div class="paypal-notice-icon">💳</div>
          <div class="paypal-notice-content">
            <h4>{{ locale === 'zh' ? 'PayPal 官方 REST API 在线唤醒网关' : 'PayPal REST API Gateway' }}</h4>
            <p>
              {{ locale === 'zh'
                ? 'PayPal 属于国际在线支付网关，用户点击支付按钮将直接唤醒 PayPal 官方收银台完成扣款与回调，无需也切勿上传收款码图片。'
                : 'PayPal is an official online API gateway. Users will be redirected directly to PayPal checkout. No QR code image is required.' }}
            </p>
            <div class="paypal-notice-action">
              <button type="button" class="btn btn-sm btn-primary" @click="goToPayPalConfig">
                ⚙️ {{ locale === 'zh' ? '前往配置 PayPal Client ID 与 Secret (公私钥)' : 'Configure PayPal API Credentials' }} →
              </button>
            </div>
          </div>
        </div>

        <div v-else-if="channelForm.code === 'test'" class="test-notice-card">
          <div class="test-notice-icon">⚡</div>
          <div class="test-notice-content">
            <h4>{{ locale === 'zh' ? '极速沙盒免密模拟通道' : 'Sandbox Fast Test Channel' }}</h4>
            <p>
              {{ locale === 'zh'
                ? '仅用于本地与演示环境下免扣款秒级开通测试，无需配置任何收款码或 API 参数。'
                : 'Used for zero-cost instant checkout simulation during development and testing.' }}
            </p>
          </div>
        </div>

        <!-- 微信、支付宝、自定义收款码：上传收款码图片 -->
        <div v-else class="form-group qr-upload-group">
          <label class="form-label font-bold">
            📷 {{ locale === 'zh' ? '收款二维码图片 (微信/支付宝收款码)' : 'Payment QR Code Image' }}
            <span class="label-badge">{{ locale === 'zh' ? '直接扫码收款' : 'Direct Scan' }}</span>
          </label>

          <div v-if="channelForm.qr_code_url" class="qr-preview-box">
            <img :src="channelForm.qr_code_url" class="qr-preview-img" alt="QR Code" />
            <div class="qr-preview-actions">
              <button type="button" class="btn btn-sm btn-outline" :disabled="uploadingQR" @click="triggerQRUpload">
                {{ uploadingQR ? (locale === 'zh' ? '正在上传...' : 'Uploading...') : (locale === 'zh' ? '🔄 更换收款码' : 'Change QR') }}
              </button>
              <button type="button" class="btn btn-sm btn-ghost text-danger" :disabled="uploadingQR" @click="channelForm.qr_code_url = ''">
                🗑️ {{ locale === 'zh' ? '移除' : 'Remove' }}
              </button>
              <small class="text-secondary" style="font-size: 11px; word-break: break-all;">{{ channelForm.qr_code_url }}</small>
            </div>
          </div>

          <div v-else class="qr-upload-dropzone" :class="{ 'is-uploading': uploadingQR }" @click="triggerQRUpload">
            <div class="dropzone-content">
              <span class="upload-icon">{{ uploadingQR ? '⏳' : '📥' }}</span>
              <strong class="upload-title">
                {{ uploadingQR ? (locale === 'zh' ? '正在上传图片...' : 'Uploading...') : (locale === 'zh' ? '点击直接上传微信 / 支付宝收款码' : 'Click to Upload QR Code Image') }}
              </strong>
              <small class="upload-hint">{{ locale === 'zh' ? '支持微信赞赏码/转账码、支付宝收钱码等 (PNG, JPG, WebP)' : 'Supports WeChat / Alipay QR code' }}</small>
            </div>
          </div>
          <input
            ref="qrFileInput"
            type="file"
            accept="image/png,image/jpeg,image/webp,image/jpg"
            style="display: none;"
            @change="onQRFileSelected"
          />
        </div>

        <div v-if="channelForm.code !== 'paypal' && channelForm.code !== 'test'" class="form-group">
          <label class="form-label font-bold">{{ locale === 'zh' ? '转账备注 / 补充指引说明 (选填)' : 'Payment Note / Instructions (Optional)' }}</label>
          <textarea v-model="channelForm.account_info" class="form-input" rows="2" placeholder="例如: 转账时请在备注中填写您的邮箱，以便极速秒级对账下发提取码..."></textarea>
        </div>

        <div class="form-group">
          <label class="form-label font-bold">{{ locale === 'zh' ? '是否立即启用' : 'Enabled' }}</label>
          <label class="switch-control">
            <input v-model="channelForm.enabled" type="checkbox" />
            <span class="switch-slider"></span>
          </label>
        </div>

        <div class="modal-actions">
          <button class="btn btn-ghost" @click="showChannelModal = false">{{ t('common.cancel') }}</button>
          <button class="btn btn-primary" :disabled="savingChannel" @click="saveChannel">
            {{ savingChannel ? t('common.loading') : (locale === 'zh' ? '保存渠道' : 'Save Channel') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { settingsAPI, paymentChannelAPI, productAPI } from '@/api'
import { useToast } from '@/composables/useToast'

const { locale, t } = useI18n()
const { success, error } = useToast()

const currentTab = ref('pricing')
const saving = ref(false)

// Tab 4 State: Dual-Currency & Payment Channels
const exchangeRate = ref(7.25)
const channels = ref([])
const loadingChannels = ref(false)
const showChannelModal = ref(false)
const isEditingChannel = ref(false)
const savingChannel = ref(false)
const uploadingQR = ref(false)
const qrFileInput = ref(null)

const channelForm = ref({
  id: '',
  name: '',
  code: 'alipay',
  scope: 'china_only',
  icon: '🔵',
  description: '',
  account_info: '',
  qr_code_url: '',
  sort_order: 1,
  enabled: true
})

// Tab 1: Membership Pricing Plans
const membershipPlans = ref([
  {
    type: 'monthly',
    name_zh: '月度会员',
    name_en: 'Monthly VIP',
    price: 9.99,
    period: 'month',
    duration_days: 30,
    enabled: true,
    featured: false,
    badge: '入门尝鲜',
    save_percent: 0,
    icon: '🌙'
  },
  {
    type: 'quarterly',
    name_zh: '季度会员',
    name_en: 'Quarterly VIP',
    price: 24.99,
    period: 'quarter',
    duration_days: 90,
    enabled: true,
    featured: true,
    badge: '热销推荐 / 省17%',
    save_percent: 17,
    icon: '⭐'
  },
  {
    type: 'yearly',
    name_zh: '年度会员',
    name_en: 'Yearly VIP',
    price: 79.99,
    period: 'year',
    duration_days: 365,
    enabled: true,
    featured: false,
    badge: '超值之选 / 省33%',
    save_percent: 33,
    icon: '👑'
  },
  {
    type: 'lifetime',
    name_zh: '终生尊享',
    name_en: 'Lifetime VIP',
    price: 199.99,
    period: 'lifetime',
    duration_days: 36500,
    enabled: true,
    featured: false,
    badge: '终身买断 / 限时尊享',
    save_percent: 50,
    icon: '💎'
  }
])

// Tab 2: Free Download Permissions
const permissions = ref({
  free_download_require_login: false,
  daily_free_limit: 0,
  guest_buy_enabled: true,
  master_unlock_codes: 'loveyou, 666, 1194, VIP888, VIP666',
  allow_test_pay: true
})

// Tab 3: PayPal Settings
const paymentSettings = ref({
  paypal_client_id: '',
  paypal_client_secret: '',
  paypal_mode: 'sandbox'
})

async function loadChannels() {
  loadingChannels.value = true
  try {
    const res = await paymentChannelAPI.adminList()
    channels.value = res.channels || res || []
  } catch (err) {
    console.warn('Failed to load payment channels:', err)
  } finally {
    loadingChannels.value = false
  }
}

function openAddChannelModal() {
  isEditingChannel.value = false
  channelForm.value = {
    id: '',
    name: '',
    code: 'alipay',
    scope: 'china_only',
    icon: '🔵',
    description: '',
    account_info: '',
    qr_code_url: '',
    sort_order: (channels.value.length + 1),
    enabled: true
  }
  showChannelModal.value = true
}

function openEditChannelModal(ch) {
  isEditingChannel.value = true
  channelForm.value = {
    id: ch.id,
    name: ch.name || '',
    code: ch.code || 'alipay',
    scope: ch.scope || 'global',
    icon: ch.icon || '💳',
    description: ch.description || '',
    account_info: ch.account_info || '',
    qr_code_url: ch.qr_code_url || '',
    sort_order: ch.sort_order || 0,
    enabled: ch.enabled !== false
  }
  showChannelModal.value = true
}

function triggerQRUpload() {
  if (qrFileInput.value) {
    qrFileInput.value.click()
  }
}

async function onQRFileSelected(e) {
  const file = e.target.files?.[0]
  if (!file) return

  if (file.size > 15 * 1024 * 1024) {
    error(locale.value === 'zh' ? '收款码图片大小不能超过 15MB' : 'Image exceeds 15MB')
    return
  }

  const formData = new FormData()
  formData.append('file', file)

  uploadingQR.value = true
  try {
    const res = await productAPI.uploadImage(formData)
    if (res && res.url) {
      channelForm.value.qr_code_url = res.url
      success(locale.value === 'zh' ? '收款码上传成功！' : 'QR code uploaded successfully!')
    } else {
      throw new Error(res?.message || '上传异常')
    }
  } catch (err) {
    error(err.message || (locale.value === 'zh' ? '收款码上传失败' : 'Upload failed'))
  } finally {
    uploadingQR.value = false
    if (e.target) {
      e.target.value = ''
    }
  }
}

function goToPayPalConfig() {
  showChannelModal.value = false
  currentTab.value = 'payment'
}

async function saveChannel() {
  if (!channelForm.value.name.trim()) {
    error(locale.value === 'zh' ? '请输入支付渠道名称' : 'Channel name required')
    return
  }
  // Clear QR code if PayPal or Test mode to avoid erroneous image storage
  if (channelForm.value.code === 'paypal' || channelForm.value.code === 'test') {
    channelForm.value.qr_code_url = ''
  }
  savingChannel.value = true
  try {
    if (isEditingChannel.value) {
      await paymentChannelAPI.update(channelForm.value.id, channelForm.value)
      success(locale.value === 'zh' ? '支付渠道已更新' : 'Channel updated')
    } else {
      await paymentChannelAPI.create(channelForm.value)
      success(locale.value === 'zh' ? '支付渠道已创建' : 'Channel created')
    }
    showChannelModal.value = false
    loadChannels()
  } catch (err) {
    error(err.message || '保存支付渠道失败')
  } finally {
    savingChannel.value = false
  }
}

async function toggleChannel(ch) {
  try {
    await paymentChannelAPI.update(ch.id, { ...ch, enabled: ch.enabled })
    success(locale.value === 'zh' ? `已${ch.enabled ? '启用' : '禁用'}渠道【${ch.name}】` : 'Status updated')
  } catch (err) {
    ch.enabled = !ch.enabled
    error(err.message || '更新状态失败')
  }
}

async function removeChannel(ch) {
  if (!confirm(locale.value === 'zh' ? `确定删除支付渠道【${ch.name}】吗？` : `Delete channel ${ch.name}?`)) {
    return
  }
  try {
    await paymentChannelAPI.delete(ch.id)
    success(locale.value === 'zh' ? '支付渠道已删除' : 'Channel deleted')
    loadChannels()
  } catch (err) {
    error(err.message || '删除失败')
  }
}

function formatScope(scope) {
  if (scope === 'china_only') return locale.value === 'zh' ? '🇨🇳 仅限中国 (国内IP)' : 'China Only'
  if (scope === 'overseas_only') return locale.value === 'zh' ? '🌐 仅限境外 (海外IP)' : 'Overseas Only'
  return locale.value === 'zh' ? '🌍 全球通用' : 'Global'
}

function getScopeClass(scope) {
  if (scope === 'china_only') return 'scope-china'
  if (scope === 'overseas_only') return 'scope-overseas'
  return 'scope-global'
}

async function loadAllSettings() {
  try {
    const res = await settingsAPI.getAll()
    if (res) {
      if (res.exchange_rate) {
        exchangeRate.value = Number(res.exchange_rate)
      }
      if (res.payment) {
        paymentSettings.value.paypal_client_id = res.payment.paypal_client_id || ''
        paymentSettings.value.paypal_client_secret = res.payment.paypal_client_secret || ''
        paymentSettings.value.paypal_mode = res.payment.paypal_mode || 'sandbox'
      }
      if (Array.isArray(res.membership_plans) && res.membership_plans.length > 0) {
        membershipPlans.value = res.membership_plans
      }
      if (res.permissions) {
        permissions.value.free_download_require_login = !!res.permissions.free_download_require_login
        permissions.value.daily_free_limit = Number(res.permissions.daily_free_limit) || 0
        permissions.value.guest_buy_enabled = res.permissions.guest_buy_enabled !== false
        permissions.value.master_unlock_codes = res.permissions.master_unlock_codes || 'loveyou, 666, 1194, VIP888'
        permissions.value.allow_test_pay = res.permissions.allow_test_pay !== false
      }
    }
  } catch (err) {
    console.error('Failed to load settings:', err)
  }
  loadChannels()
}

function resetToDefaultPlans() {
  membershipPlans.value = [
    {
      type: 'monthly',
      name_zh: '月度会员',
      name_en: 'Monthly VIP',
      price: 9.99,
      period: 'month',
      duration_days: 30,
      enabled: true,
      featured: false,
      badge: '入门尝鲜',
      save_percent: 0,
      icon: '🌙'
    },
    {
      type: 'quarterly',
      name_zh: '季度会员',
      name_en: 'Quarterly VIP',
      price: 24.99,
      period: 'quarter',
      duration_days: 90,
      enabled: true,
      featured: true,
      badge: '热销推荐 / 省17%',
      save_percent: 17,
      icon: '⭐'
    },
    {
      type: 'yearly',
      name_zh: '年度会员',
      name_en: 'Yearly VIP',
      price: 79.99,
      period: 'year',
      duration_days: 365,
      enabled: true,
      featured: false,
      badge: '超值之选 / 省33%',
      save_percent: 33,
      icon: '👑'
    },
    {
      type: 'lifetime',
      name_zh: '终生尊享',
      name_en: 'Lifetime VIP',
      price: 199.99,
      period: 'lifetime',
      duration_days: 36500,
      enabled: true,
      featured: false,
      badge: '终身买断 / 限时尊享',
      save_percent: 50,
      icon: '💎'
    }
  ]
  success(locale.value === 'zh' ? '已恢复默认会员定价方案！点击右上角保存即可生效。' : 'Reset to default plans.')
}

function addNewPlan() {
  membershipPlans.value.push({
    type: 'custom_' + Date.now().toString().slice(-4),
    name_zh: '自定义特惠方案',
    name_en: 'Custom Plan',
    price: 19.99,
    period: 'month',
    duration_days: 60,
    enabled: true,
    featured: false,
    badge: '限时特惠',
    save_percent: 10,
    icon: '✨'
  })
}

function removePlan(index) {
  if (confirm(locale.value === 'zh' ? '确定删除该会员方案吗？' : 'Delete this plan?')) {
    membershipPlans.value.splice(index, 1)
  }
}

async function saveAllSettings() {
  saving.value = true
  try {
    const payload = {
      exchange_rate: exchangeRate.value,
      payment: paymentSettings.value,
      membership_plans: membershipPlans.value,
      permissions: permissions.value
    }
    await settingsAPI.updateAll(payload)
    success(locale.value === 'zh' ? '🎉 会员定价、汇率与权限规则等全部配置已保存生效！' : 'All settings saved successfully!')
  } catch (err) {
    error(err.message || (locale.value === 'zh' ? '保存失败' : 'Failed to save settings'))
  } finally {
    saving.value = false
  }
}

onMounted(loadAllSettings)
</script>

<style scoped>
.admin-settings {
  max-width: 1080px;
  padding-bottom: 60px;
}

.settings-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.settings-title {
  font-size: 1.6rem;
  font-weight: 800;
  margin: 0 0 6px;
  background: linear-gradient(135deg, var(--c-text), var(--c-accent));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.settings-subtitle {
  font-size: 0.9rem;
  color: var(--c-text-secondary);
  margin: 0;
}

.save-btn {
  padding: 12px 24px;
  font-size: 0.95rem;
  font-weight: 700;
  border-radius: 14px;
  box-shadow: 0 4px 14px rgba(0, 113, 227, 0.35);
}

/* Nav Tabs */
.settings-nav-tabs {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
  margin-bottom: 24px;
}

@media (max-width: 1024px) {
  .settings-nav-tabs {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 640px) {
  .settings-nav-tabs {
    grid-template-columns: 1fr;
  }
}

.nav-tab-btn {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 16px 20px;
  background: var(--c-surface);
  border: 2px solid var(--c-border-light);
  border-radius: 18px;
  cursor: pointer;
  text-align: left;
  transition: border-color 0.2s ease, background 0.2s ease, box-shadow 0.2s ease;
}

.nav-tab-btn:hover {
  border-color: var(--c-accent);
}

.nav-tab-btn.active {
  background: var(--c-surface);
  border-color: var(--c-accent);
  box-shadow: 0 8px 24px rgba(0, 113, 227, 0.15);
}

/* Payment Channels Table & Badges */
.channels-table-wrap {
  overflow-x: auto;
  border-radius: 14px;
}

.scope-pill {
  display: inline-flex;
  align-items: center;
  padding: 3px 10px;
  border-radius: 999px;
  font-size: 0.78rem;
  font-weight: 700;
  white-space: nowrap;
}

.scope-china {
  background: rgba(239, 68, 68, 0.12);
  color: #dc2626;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.scope-overseas {
  background: rgba(59, 130, 246, 0.12);
  color: #2563eb;
  border: 1px solid rgba(59, 130, 246, 0.3);
}

.scope-global {
  background: rgba(16, 185, 129, 0.12);
  color: #059669;
  border: 1px solid rgba(16, 185, 129, 0.3);
}

.nav-tab-icon {
  font-size: 2rem;
  line-height: 1;
}

.nav-tab-info strong {
  display: block;
  font-size: 1rem;
  color: var(--c-text);
  margin-bottom: 2px;
}

.nav-tab-info small {
  font-size: 0.78rem;
  color: var(--c-text-tertiary);
  display: block;
}

/* Settings Cards */
.settings-grid {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.settings-card {
  padding: 28px 32px;
  border-radius: 24px;
  background: var(--c-surface);
  border: 1px solid var(--c-border-light);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.04);
}

.card-head {
  display: flex;
  align-items: flex-start;
  gap: 18px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--c-border-light);
  margin-bottom: 24px;
}

.head-title-wrap {
  flex: 1;
}

.head-flex {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.head-actions {
  display: flex;
  gap: 10px;
}

.head-icon {
  width: 52px;
  height: 52px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.6rem;
  flex-shrink: 0;
}

.icon-gold { background: rgba(245, 158, 11, 0.15); }
.icon-green { background: rgba(16, 185, 129, 0.15); }
.icon-amber { background: rgba(249, 115, 22, 0.15); }
.icon-blue { background: rgba(0, 113, 227, 0.12); }

.card-head h3 {
  margin: 0 0 4px;
  font-size: 1.25rem;
  font-weight: 800;
}

.card-head p {
  margin: 0;
  font-size: 0.86rem;
  color: var(--c-text-secondary);
  line-height: 1.45;
}

/* Pricing Stack */
.pricing-plans-stack {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.plan-edit-card {
  padding: 20px 24px;
  background: var(--c-bg);
  border: 2px solid var(--c-border);
  border-radius: 20px;
  transition: all 0.2s ease;
}

.plan-edit-card.featured {
  border-color: #f59e0b;
  background: rgba(245, 158, 11, 0.03);
}

.plan-edit-card.disabled {
  opacity: 0.6;
  filter: grayscale(0.5);
}

.plan-card-topbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 14px;
  margin-bottom: 16px;
  border-bottom: 1px dashed var(--c-border);
  flex-wrap: wrap;
  gap: 12px;
}

.plan-card-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.plan-icon-preview {
  font-size: 1.8rem;
}

.plan-type-pill {
  display: inline-block;
  font-size: 0.72rem;
  font-weight: 800;
  text-transform: uppercase;
  background: rgba(0, 113, 227, 0.12);
  color: var(--c-accent);
  padding: 2px 8px;
  border-radius: 6px;
  margin-right: 8px;
}

.plan-name-preview {
  font-size: 1.05rem;
  color: var(--c-text);
}

.plan-card-toggles {
  display: flex;
  align-items: center;
  gap: 12px;
}

.btn-delete-plan {
  background: none;
  border: none;
  font-size: 1.1rem;
  cursor: pointer;
  padding: 4px;
  opacity: 0.7;
  transition: opacity 0.2s;
}

.btn-delete-plan:hover {
  opacity: 1;
}

.plan-form-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 14px;
}

/* Radio Cards */
.radio-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 16px;
  margin-top: 10px;
}

.radio-card {
  display: flex;
  align-items: flex-start;
  gap: 14px;
  padding: 18px 20px;
  border: 2px solid var(--c-border);
  border-radius: 18px;
  cursor: pointer;
  transition: all 0.2s ease;
  background: var(--c-bg);
}

.radio-card:hover {
  border-color: var(--c-accent);
}

.radio-card.selected {
  border-color: var(--c-accent);
  background: rgba(0, 113, 227, 0.05);
}

.radio-card input {
  margin-top: 4px;
}

.radio-info strong {
  display: block;
  font-size: 0.96rem;
  margin-bottom: 6px;
  color: var(--c-text);
}

.radio-info span {
  font-size: 0.82rem;
  color: var(--c-text-secondary);
  line-height: 1.45;
  display: block;
}

/* Toggle row */
.toggle-setting-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 20px;
  padding: 8px 0;
}

.toggle-setting-row strong {
  display: block;
  font-size: 1rem;
  margin-bottom: 4px;
  color: var(--c-text);
}

.toggle-setting-row p {
  margin: 0;
  font-size: 0.84rem;
  color: var(--c-text-secondary);
  line-height: 1.45;
}

.section-divider {
  border: none;
  border-top: 1px solid var(--c-border-light);
  margin: 20px 0;
}

/* Inputs & Badges */
.input-with-prefix, .input-with-suffix {
  position: relative;
  display: flex;
  align-items: center;
}

.input-prefix {
  position: absolute;
  left: 14px;
  font-weight: 800;
  color: var(--c-text-secondary);
}

.input-with-prefix input {
  padding-left: 30px;
}

.input-suffix {
  position: absolute;
  right: 14px;
  font-size: 0.84rem;
  font-weight: 600;
  color: var(--c-text-secondary);
}

.input-with-suffix input {
  padding-right: 68px;
}

.font-mono {
  font-family: monospace;
  letter-spacing: 0.5px;
}

.font-bold {
  font-weight: 700;
}

.form-tip {
  display: block;
  font-size: 0.76rem;
  color: var(--c-text-tertiary);
  margin-top: 5px;
}

/* Switch */
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
  top: 0 !important;
  left: 0 !important;
  right: 0 !important;
  bottom: 0 !important;
  width: 100% !important;
  height: 100% !important;
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

.switch-control.lg {
  width: 52px !important;
  height: 28px !important;
}

.switch-control.lg .switch-slider {
  border-radius: 28px;
}

.switch-control.lg .switch-slider:before {
  height: 20px !important;
  width: 20px !important;
  left: 4px !important;
  bottom: 4px !important;
}

.switch-control.lg input:checked + .switch-slider:before {
  transform: translateX(24px) !important;
}

/* Channels Data Table */
.channels-table-wrap {
  overflow-x: auto;
  border-radius: 16px;
  border: 1px solid var(--c-border-light, #e2e8f0);
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.data-table th, .data-table td {
  padding: 14px 16px;
  border-bottom: 1px solid var(--c-border-light, #f1f5f9);
  vertical-align: middle;
}

.data-table th {
  font-size: 0.82rem;
  font-weight: 700;
  color: var(--c-text-secondary, #64748b);
  background: var(--c-bg, #f8fafc);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.data-table tbody tr:hover {
  background: rgba(0, 0, 0, 0.02);
}

/* Modal Overlay & Dialog */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.65);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 99999;
  padding: 20px;
}

.modal-box {
  background: var(--c-surface, #ffffff);
  color: var(--c-text, #0f172a);
  width: 100%;
  max-width: 540px;
  border-radius: 24px;
  padding: 28px 32px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.35);
  border: 1px solid var(--c-border, rgba(0, 0, 0, 0.1));
  max-height: 90vh;
  overflow-y: auto;
  animation: modalIn 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes modalIn {
  from { opacity: 0; transform: scale(0.95) translateY(10px); }
  to { opacity: 1; transform: scale(1) translateY(0); }
}

.modal-box h3 {
  margin: 0 0 6px;
  font-size: 1.35rem;
  font-weight: 800;
}

.user-sub {
  font-size: 0.86rem;
  color: var(--c-text-secondary, #64748b);
  margin: 0 0 20px;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid var(--c-border-light, #f1f5f9);
}

/* QR Code Uploader in Channel Modal */
.qr-upload-group {
  margin-top: 4px;
}

.label-badge {
  display: inline-block;
  font-size: 0.72rem;
  font-weight: 700;
  color: #10b981;
  background: rgba(16, 185, 129, 0.12);
  padding: 2px 8px;
  border-radius: 6px;
  margin-left: 8px;
  vertical-align: middle;
}

.qr-preview-box {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 14px 18px;
  background: var(--c-bg, #f8fafc);
  border: 1.5px dashed var(--c-border-light, #cbd5e1);
  border-radius: 16px;
}

.qr-preview-img {
  width: 100px;
  height: 100px;
  object-fit: contain;
  border-radius: 12px;
  background: #ffffff;
  padding: 4px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border: 1px solid #e2e8f0;
  flex-shrink: 0;
}

.qr-preview-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex: 1;
}

.qr-upload-dropzone {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 22px 18px;
  background: rgba(0, 113, 227, 0.03);
  border: 2px dashed rgba(0, 113, 227, 0.35);
  border-radius: 16px;
  cursor: pointer;
  transition: all 0.2s ease;
  text-align: center;
}

.qr-upload-dropzone:hover {
  background: rgba(0, 113, 227, 0.08);
  border-color: #0071e3;
  transform: translateY(-1px);
}

.qr-upload-dropzone.is-uploading {
  opacity: 0.6;
  pointer-events: none;
}

.dropzone-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
}

.upload-icon {
  font-size: 2rem;
  line-height: 1;
}

.upload-title {
  font-size: 0.94rem;
  color: var(--c-text, #0f172a);
}

.upload-hint {
  font-size: 0.78rem;
  color: var(--c-text-secondary, #64748b);
}

.channel-qr-thumb {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  object-fit: cover;
  border: 1px solid #cbd5e1;
  background: #ffffff;
  padding: 2px;
  flex-shrink: 0;
}

.qr-has-badge {
  display: inline-block;
  font-size: 0.68rem;
  color: #10b981;
  background: rgba(16, 185, 129, 0.12);
  padding: 1px 6px;
  border-radius: 4px;
  margin-left: 6px;
}

.paypal-gateway-badge {
  display: inline-block;
  font-size: 0.68rem;
  font-weight: 700;
  color: #0070ba;
  background: rgba(0, 112, 186, 0.12);
  border: 1px solid rgba(0, 112, 186, 0.25);
  padding: 1px 7px;
  border-radius: 4px;
  margin-left: 6px;
}

.test-gateway-badge {
  display: inline-block;
  font-size: 0.68rem;
  font-weight: 700;
  color: #d97706;
  background: rgba(217, 119, 6, 0.12);
  border: 1px solid rgba(217, 119, 6, 0.25);
  padding: 1px 7px;
  border-radius: 4px;
  margin-left: 6px;
}

.qr-empty-badge {
  display: inline-block;
  font-size: 0.68rem;
  color: #f59e0b;
  background: rgba(245, 158, 11, 0.12);
  padding: 1px 6px;
  border-radius: 4px;
  margin-left: 6px;
}

.paypal-notice-card, .test-notice-card {
  display: flex;
  gap: 16px;
  padding: 18px 20px;
  border-radius: 14px;
  margin-bottom: 20px;
}

.paypal-notice-card {
  background: linear-gradient(135deg, rgba(0, 112, 186, 0.08), rgba(0, 48, 135, 0.04));
  border: 1.5px solid rgba(0, 112, 186, 0.28);
}

.test-notice-card {
  background: linear-gradient(135deg, rgba(217, 119, 6, 0.08), rgba(180, 83, 9, 0.04));
  border: 1.5px solid rgba(217, 119, 6, 0.28);
}

.paypal-notice-icon, .test-notice-icon {
  font-size: 2rem;
  line-height: 1;
  flex-shrink: 0;
  padding-top: 2px;
}

.paypal-notice-content h4, .test-notice-content h4 {
  font-size: 0.96rem;
  font-weight: 800;
  margin: 0 0 6px 0;
  color: var(--c-text);
}

.paypal-notice-content p, .test-notice-content p {
  font-size: 0.82rem;
  line-height: 1.5;
  color: var(--c-text-secondary);
  margin: 0 0 12px 0;
}

.paypal-notice-action {
  display: flex;
  align-items: center;
}

/* Alerts */
.alert-box {
  padding: 16px 20px;
  border-radius: 16px;
  font-size: 0.86rem;
  line-height: 1.5;
}

.alert-info {
  background: rgba(0, 113, 227, 0.08);
  border: 1px solid rgba(0, 113, 227, 0.25);
  color: #1e40af;
}

.mt-4 { margin-top: 20px; }
.ml-2 { margin-left: 8px; }
.text-danger { color: #ef4444; }

/* Mobile Responsive Optimization for Admin Settings */
@media (max-width: 768px) {
  .settings-header {
    flex-direction: column;
    align-items: stretch;
    gap: 14px;
  }
  .save-btn {
    width: 100%;
    padding: 14px;
  }
  .settings-card {
    padding: 18px 16px !important;
    border-radius: 18px !important;
  }
  .head-flex {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  .head-actions {
    width: 100%;
    flex-direction: column;
    gap: 8px;
  }
  .head-actions button {
    width: 100%;
  }
  .plan-card-topbar {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  .plan-card-toggles {
    width: 100%;
    justify-content: space-between;
  }
  .plan-form-grid {
    grid-template-columns: 1fr 1fr;
    gap: 10px;
  }
  .radio-cards {
    grid-template-columns: 1fr;
  }
  .toggle-setting-row {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
}

@media (max-width: 480px) {
  .plan-form-grid {
    grid-template-columns: 1fr;
  }
}
</style>
