#!/bin/bash
# ========================================================
# WallArt - 1秒一键纠正反向代理端口对调问题
# 保证:
#   * wallpaper.icangoanyway.com -> 5173 (Uni-app 客户端前台商城)
#   * admin.icangoanyway.com     -> 5174 (Vue 3 独立后台管理)
#   * ios.icangoanyway.com       -> 5173 (Uni-app 客户端前台商城)
# ========================================================

set -e
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "========================================================"
echo "  [1/2] 深度扫描并纠正宿主机反向代理 (宝塔 / Nginx)..."
echo "========================================================"

CHANGED=0

# 1. 深度检测并修复宝塔反向代理目录 (/www/server/panel/vhost/nginx/proxy/域名/*.conf)
BT_PROXY_DIR="/www/server/panel/vhost/nginx/proxy"
if [ -d "$BT_PROXY_DIR" ]; then
    # 纠正 admin 目录下的所有反代配置文件 (必须指向 5174)
    for conf in $(find "$BT_PROXY_DIR" -type f \( -name "*.conf" -o -name "proxy.conf" \) -path "*admin.icangoanyway.com*" 2>/dev/null || true); do
        if grep -q "5173" "$conf"; then
            echo "[自动修复 - 宝塔反代] 发现 $conf 错误指向 5173，正在修正为 5174..."
            sed -i -E 's/(127\.0\.0\.1|localhost):5173/\1:5174/g' "$conf"
            sed -i 's/:5173/:5174/g' "$conf"
            CHANGED=1
        fi
    done

    # 纠正 wallpaper 目录下的所有反代配置文件 (必须指向 5173)
    for conf in $(find "$BT_PROXY_DIR" -type f \( -name "*.conf" -o -name "proxy.conf" \) -path "*wallpaper.icangoanyway.com*" 2>/dev/null || true); do
        if grep -q "5174" "$conf"; then
            echo "[自动修复 - 宝塔反代] 发现 $conf 错误指向 5174，正在修正为 5173..."
            sed -i -E 's/(127\.0\.0\.1|localhost):5174/\1:5173/g' "$conf"
            sed -i 's/:5174/:5173/g' "$conf"
            CHANGED=1
        fi
    done

    # 纠正 ios 目录下的所有反代配置文件 (必须指向 5173)
    for conf in $(find "$BT_PROXY_DIR" -type f \( -name "*.conf" -o -name "proxy.conf" \) -path "*ios.icangoanyway.com*" 2>/dev/null || true); do
        if grep -q "5174" "$conf"; then
            echo "[自动修复 - 宝塔反代] 发现 $conf 错误指向 5174，正在修正为 5173..."
            sed -i -E 's/(127\.0\.0\.1|localhost):5174/\1:5173/g' "$conf"
            sed -i 's/:5174/:5173/g' "$conf"
            CHANGED=1
        fi
    done
fi

# 2. 检查并纠正宝塔数据配置文件 (/www/server/panel/data/*.json)
if [ -d "/www/server/panel/data" ]; then
    for json_file in /www/server/panel/data/proxy.json /www/server/panel/data/*.json; do
        if [ -f "$json_file" ]; then
            if grep -q "admin.icangoanyway.com" "$json_file" 2>/dev/null && grep -q "5173" "$json_file" 2>/dev/null; then
                echo "[自动修复 - 宝塔元数据] 校准 $json_file 中的 admin 代理端口为 5174..."
                sed -i -E 's/(admin\.icangoanyway\.com[^}]*):5173/\1:5174/g' "$json_file" 2>/dev/null || true
                CHANGED=1
            fi
            if grep -q "wallpaper.icangoanyway.com" "$json_file" 2>/dev/null && grep -q "5174" "$json_file" 2>/dev/null; then
                echo "[自动修复 - 宝塔元数据] 校准 $json_file 中的 wallpaper 代理端口为 5173..."
                sed -i -E 's/(wallpaper\.icangoanyway\.com[^}]*):5174/\1:5173/g' "$json_file" 2>/dev/null || true
                CHANGED=1
            fi
        fi
    done
fi

# 3. 扫描通用 Nginx 配置文件目录 (宝塔主配置、/etc/nginx、/usr/local/nginx 等)
CONFIG_DIRS=(
    "/www/server/panel/vhost/nginx"
    "/etc/nginx/conf.d"
    "/etc/nginx/sites-enabled"
    "/etc/nginx/sites-available"
    "/usr/local/nginx/conf/vhost"
)

for dir in "${CONFIG_DIRS[@]}"; do
    if [ -d "$dir" ]; then
        # 针对 admin 域名配置 (仅当该文件不含 wallpaper 时安全替换)
        for f in $(grep -rnwl "$dir" -e "admin.icangoanyway.com" 2>/dev/null || true); do
            if grep -q "5173" "$f"; then
                if ! grep -q "wallpaper.icangoanyway.com" "$f"; then
                    echo "[自动修复 - Nginx] 发现 $f (admin) 错误指向 5173，正在修正为 5174..."
                    sed -i -E 's/(127\.0\.0\.1|localhost):5173/\1:5174/g' "$f"
                    sed -i 's/:5173/:5174/g' "$f"
                    CHANGED=1
                fi
            fi
        done

        # 针对 wallpaper 域名配置
        for f in $(grep -rnwl "$dir" -e "wallpaper.icangoanyway.com" 2>/dev/null || true); do
            if grep -q "5174" "$f"; then
                if ! grep -q "admin.icangoanyway.com" "$f"; then
                    echo "[自动修复 - Nginx] 发现 $f (wallpaper) 错误指向 5174，正在修正为 5173..."
                    sed -i -E 's/(127\.0\.0\.1|localhost):5174/\1:5173/g' "$f"
                    sed -i 's/:5174/:5173/g' "$f"
                    CHANGED=1
                fi
            fi
        done
    fi
done

# 4. 重新加载 Nginx 服务
if [ $CHANGED -eq 1 ]; then
    echo "[OK] 正在重新加载宿主机 Nginx 反代配置..."
    if command -v nginx &> /dev/null; then
        nginx -t 2>/dev/null && nginx -s reload 2>/dev/null || true
    fi
    if command -v systemctl &> /dev/null; then
        systemctl reload nginx 2>/dev/null || true
    fi
    if [ -f "/etc/init.d/nginx" ]; then
        /etc/init.d/nginx reload 2>/dev/null || true
    fi
    echo "[OK] 宿主机反代端口已成功修正生效！"
else
    echo "宿主机反代配置未检测到异常（或直接使用 Docker Gateway 分流）。"
fi

echo ""
echo "========================================================"
echo "  [2/2] 热重载 Docker 网关 (wallpaper-gateway)..."
echo "========================================================"

if command -v docker &> /dev/null; then
    if docker compose version &> /dev/null; then
        docker compose restart wallpaper-gateway 2>/dev/null || true
    elif command -v docker-compose &> /dev/null; then
        docker-compose restart wallpaper-gateway 2>/dev/null || true
    fi
fi

echo ""
echo "========================================================"
echo "  🎉 纠偏校准完成！"
echo "  * 客户端前台商城 : https://wallpaper.icangoanyway.com -> 5173"
echo "  * 独立管理后台   : https://admin.icangoanyway.com     -> 5174"
echo "  * iOS 兼容域名   : https://ios.icangoanyway.com       -> 5173"
echo "========================================================"
