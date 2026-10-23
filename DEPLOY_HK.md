# WallArt 壁纸商城系统 - 域名、端口与后台管理部署指南

本方案已配置为支持 **80 (HTTP) 与 443 (HTTPS) 端口** 全自动分流，且**服务器原有容器端口（5173、5174、8080）保持不变**：
- 📱 **客户端商城前台 (Uni-app)**：`http://wallpaper.icangoanyway.com` 或 `https://wallpaper.icangoanyway.com` (精准直达 5173 前台商城)
- ⚙️ **独立管理后台 (Vue 3 Admin)**：`http://admin.icangoanyway.com` 或 `https://admin.icangoanyway.com` (精准直达 5174 管理后台)
- 🔌 **备用/原端口直接访问**：`http://服务器IP:5173` 与 `:5174` 依然同时有效！

---

## 一、为什么之前提示 "显示5000报错"？

1. **原因揭秘**：这不是 HTTP 500 报错，而是之前的旧版前端在登录成功后，默认跳转到了根路径 `/`（旧前台展示页），该页面正中间有一行醒目的展示文案：
   ```
   3+ 壁纸   5000+ 用户   4K 高清
   ```
   管理员登录后看到的是前台页面上的 **`5000+ 用户`**，误以为是系统报错了 5000，并且未能直接进入后台控制台（`/admin`）。
2. **已彻底修复**：
   - 访问 `admin.icangoanyway.com` 根路径将**自动重定向至 `/admin` 管理后台**；
   - 管理员登录成功后，系统**直接跳转至 `/admin` 数据看板**，不再跳转至前台展示页；
   - 顶部导航栏右侧为管理员增加了明显的 **`⚙️ 管理后台`** 快捷直达按钮。

---

## 二、为什么之前提示 SSL 证书过期或风险？

1. **原因**：初始安装包中的自签名占位证书时间设置在了未来年份，导致当前电脑/手机访问时被浏览器判定为“证书日期无效/未生效”；并且自签名证书未经过权威机构签名，会触发浏览器的红色风险警告。
2. **已彻底修复**：
   - 安装包内的占位证书已替换为 **2020~2035 年（15年长效）** 通用多域名证书，绝不会再报证书过期；
   - 80 端口移除了强制跳转，可以直接通过 `http://wallpaper.icangoanyway.com` 和 `http://admin.icangoanyway.com` 纯净访问，无需担心证书拦截；
   - 执行 `./setup-ssl.sh` 后，即可向官方权威机构 Let's Encrypt 签发公信正式证书，点亮**绿色安全小锁**。

---

## 三、服务器一键部署与运行

将最新的 `wallpaper-docker-prod.zip` 上传到服务器并执行：

```bash
cd /data/wallpaper
unzip -o wallpaper-docker-prod.zip
chmod +x *.sh
./restart.sh
```

如需开启公网正式权威 SSL（HTTPS）安全证书：
```bash
./setup-ssl.sh
```

---

## 四、智能网关分流与端口架构

| 访问入口 (浏览器访问地址) | 端口 | 协议 | 内部代理目标 | 说明 |
| :--- | :--- | :--- | :--- | :--- |
| **`http://wallpaper.icangoanyway.com`** | 80 | HTTP | `wallpaper-app:80` (5173) | 客户端前台商城 (极速直达，无证书拦截) |
| **`https://wallpaper.icangoanyway.com`** | 443 | HTTPS | `wallpaper-app:80` (5173) | 客户端前台商城 (加密安全连接) |
| **`http://admin.icangoanyway.com`** | 80 | HTTP | `wallpaper-admin:80` (5174) | 独立管理后台 (自动跳转 /admin 登录) |
| **`https://admin.icangoanyway.com`** | 443 | HTTPS | `wallpaper-admin:80` (5174) | 独立管理后台 (加密安全连接) |
| **`http://服务器IP:5173`** | 5173 | HTTP | `wallpaper-app:80` | 原前台直接端口访问（保持不变） |
| **`http://服务器IP:5174`** | 5174 | HTTP | `wallpaper-admin:80` | 原后台直接端口访问（保持不变） |
| **`http://服务器IP:8080`** | 8080 | HTTP | `wallpaper-backend:8080` | 原后端 API 端口访问（保持不变） |

---

## 五、管理员默认登录账号与密码

- **后台入口**：`http://admin.icangoanyway.com/admin` 或 `https://admin.icangoanyway.com/admin`
- **管理员账号**：`admin@wallpaper.com`
- **管理员密码**：`admin123456`
