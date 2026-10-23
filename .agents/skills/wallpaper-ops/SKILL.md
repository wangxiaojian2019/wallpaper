---
name: wallpaper-ops
description: Operations, cross-compilation, deployment, and packaging workflows for the WallArt wallpaper e-commerce platform.
---

# WallArt 4K 壁纸平台运维与构建指南

本技能为 WallArt 壁纸站（Go + Gin + Uni-app + Vue3 + Docker Compose）提供标准的构建、交叉编译、打包与服务器维护操作规约。

## 1. 后端服务构建与交叉编译

因为生产服务器为 Linux amd64 容器环境，本地为 Windows PowerShell，构建时必须通过环境变量控制输出。

### Linux 交叉编译（用于生产部署）
```powershell
$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:GOCACHE = "d:\WALLPAPER\backend\.cache"
cd d:\WALLPAPER\backend
go build -o wallpaper-server-linux main.go
```

### Windows 本地调试编译
```powershell
$env:GOOS = "windows"
$env:GOARCH = "amd64"
$env:GOCACHE = "d:\WALLPAPER\backend\.cache"
cd d:\WALLPAPER\backend
go build -o wallpaper-server.exe main.go
```

## 2. 生产发布包构建

在项目根目录运行标准打包脚本：
```powershell
cd d:\WALLPAPER
powershell -ExecutionPolicy Bypass -File .\package-wallpaper.ps1
```
打包输出为 `wallpaper-docker-prod.zip`。
该脚本已自动配置保护措施：
- 排除 `backend/wallpaper.db`（防止覆盖线上已有订单、会员与商品数据）
- 排除 `backend/uploads/*` 已有图片（只保留 `.gitkeep`）
- 排除 `node_modules` 与临时构建产物

## 3. 服务器部署流程

将 `wallpaper-docker-prod.zip` 上传至香港服务器根目录后执行：
```bash
unzip -o wallpaper-docker-prod.zip
chmod +x *.sh backend/wallpaper-server-linux
./restart.sh
```

## 4. 端口异常自动纠偏

若因宝塔反向代理或历史配置导致 `wallpaper` 与 `admin` 端口错乱（5173/5174 倒置）：
```bash
./fix-ports.sh
```
标准映射：
- 5173 -> `wallpaper.icangoanyway.com` (Uni-app 买家端)
- 5174 -> `admin.icangoanyway.com` (Vue3 管理后台)
