#!/bin/bash
set -e

chmod +x *.sh 2>/dev/null || true

echo "========================================================"
echo "  WallArt 壁纸商城系统 - Docker 生产服务一键启动"
echo "========================================================"

if command -v docker &> /dev/null; then
    if docker compose version &> /dev/null; then
        docker compose up -d --build
        echo ""
        docker compose ps
    elif command -v docker-compose &> /dev/null; then
        docker-compose up -d --build
        echo ""
        docker-compose ps
    else
        echo "[错误] 未找到 docker compose 命令，请先安装 Docker！"
        exit 1
    fi
else
    echo "[错误] 系统未检测到 Docker 环境，请先安装 Docker！"
    exit 1
fi

echo ""
echo "========================================================"
echo "服务已成功启动！"
echo "  * 客户端主站 (H5/Web) : https://wallpaper.icangoanyway.com (80/443端口直达)"
echo "  * 独立管理后台 (Admin): https://admin.icangoanyway.com (80/443端口直达)"
echo "  * 备用直接端口访问    : http://服务器IP:5173 与 :5174"
echo "  * 一键部署权威SSL证书 : ./setup-ssl.sh"
echo "========================================================"
