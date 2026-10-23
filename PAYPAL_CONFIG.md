# PayPal 个人卖家 / 商业收款接入指南

本项目已经完整实现了 **PayPal Orders v2 REST API** 与 **PayPal Smart Checkout** 流程。

---

## 1. 为什么个人卖家可以用 PayPal 收款？

PayPal 支持个人卖家（Personal/Premier 账号或升级为免费的 Business 账号）直接收款：
1. 前往 [PayPal Developer 开发者后台](https://developer.paypal.com/dashboard/) 登录你的 PayPal 账号。
2. 点击 **Apps & Credentials**。
3. 在 **REST API apps** 下点击 **Create App**。
4. 获得你的 **Client ID** 和 **Secret Key**。

---

## 2. 环境变量配置

在 `backend` 目录下创建或编辑 `.env` 文件（或直接配置系统环境变量）：

```env
# 服务端口
PORT=8080

# 数据库路径（SQLite 零配置开箱即用）
DB_PATH=wallpaper.db

# JWT 鉴权密钥
JWT_SECRET=your-random-jwt-secret-key-2026

# 前端访问地址
APP_BASE_URL=http://localhost:5173

# PayPal 配置
# 沙盒测试模式填 sandbox，正式收款模式填 live
PAYPAL_MODE=sandbox

# 你的 PayPal App Client ID
PAYPAL_CLIENT_ID=your_paypal_client_id_here

# 你的 PayPal App Secret Key
PAYPAL_CLIENT_SECRET=your_paypal_client_secret_here
```

> **提示**：未配置真实密钥时，系统内置了**沙盒模拟测试通道**，你可以直接在本地体验完整的支付下单、解锁网盘、开通会员流程。配置真实密钥后自动无缝切换到 PayPal 官方线上网关！

---

## 3. 两种售卖模式如何运作？

### 模式 A：单品壁纸售卖
- 用户在壁纸详情页点击「立即购买 $X.XX USD」。
- 发起 PayPal 支付订单。
- 支付完成后，该商品立即解锁，详情页与个人中心均显示**百度网盘、夸克网盘、Google Drive**的下载链接与提取码。

### 模式 B：月 / 季 / 年会员订阅模式
- 用户在「会员」页面选择方案：
  - **月度会员**：30 天无限畅享所有壁纸下载
  - **季度会员**：90 天无限畅享
  - **年度会员**：365 天无限畅享
- 支付完成后，用户角色升级为对应会员等级，全站所有已上架壁纸全部自动解锁，无需再逐个购买。

---

## 4. 商品富文本编辑与演示视频

在管理后台（`/admin/products/new` 或编辑页面）：
- **封面图片**：输入高清截图或大图 URL，前台自适应展示。
- **演示视频**：输入 Bilibili、YouTube 或直链 MP4 视频链接，详情页以超链/播放器呈现。
- **富文本描述**：基于 TipTap 编辑器，支持标题、粗体、斜体、列表、引用块、链接等丰富排版。
- **网盘下载地址**：可同时配置百度云盘、夸克网盘、Google 云盘的下载 URL 与提取码。只有已付费或会员才能查看真实下载信息，防盗链泄露。

---

## 5. 管理员账号与默认凭证

- **后台地址**：`http://localhost:5173/admin`
- **默认管理员邮箱**：`admin@wallpaper.com`
- **默认管理员密码**：`admin123456`
