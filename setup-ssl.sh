#!/bin/bash
# ========================================================
# WallArt - Let's Encrypt SSL 免费官方安全证书一键全自动申请与配置脚本
# 覆盖域名: wallpaper.icangoanyway.com, admin.icangoanyway.com, ios.icangoanyway.com
# ========================================================

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

echo -e "${CYAN}========================================================${NC}"
echo -e "${CYAN}  WallArt 壁纸商城系统 - 一键自动化配置 HTTPS SSL 证书  ${NC}"
echo -e "${CYAN}  主申请域名: wallpaper.icangoanyway.com                 ${NC}"
echo -e "${CYAN}  后台域名:   admin.icangoanyway.com                     ${NC}"
echo -e "${CYAN}  兼容域名:   ios.icangoanyway.com                       ${NC}"
echo -e "${CYAN}========================================================${NC}\n"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# 1. 检查 Docker 环境
if ! command -v docker &> /dev/null; then
    echo -e "${RED}[错误] 未检测到 docker 命令，请先安装 Docker！${NC}"
    exit 1
fi

# 2. 确保目录结构存在
mkdir -p "$SCRIPT_DIR/certbot-www"
mkdir -p "$SCRIPT_DIR/ssl"
mkdir -p "$SCRIPT_DIR/certbot-etc"

# 3. 确保当前服务正在运行 (尤其是 80 端口网关服务用于 ACME 验证)
echo -e "${YELLOW}[1/4] 检查并确保 gateway 容器正在运行 (监听 80 端口)...${NC}"
docker compose up -d

# 等待几秒钟让 Nginx 就绪
sleep 3

# 4. 获取联系邮箱 (可选，用于证书到期前提醒)
EMAIL="$1"
if [ -z "$EMAIL" ]; then
    EMAIL="admin@icangoanyway.com"
fi

echo -e "${YELLOW}[2/4] 通过 Certbot 官方镜像申请权威免费 SSL 证书...${NC}"
echo -e "      验证域名: wallpaper.icangoanyway.com, admin.icangoanyway.com, ios.icangoanyway.com"
echo -e "      通知邮箱: $EMAIL\n"

# 5. 执行 Certbot Webroot 模式申请 (零宕机热申请)
docker run --rm \
    -v "$SCRIPT_DIR/certbot-www:/var/www/certbot" \
    -v "$SCRIPT_DIR/certbot-etc:/etc/letsencrypt" \
    certbot/certbot certonly \
    --webroot \
    -w /var/www/certbot \
    -d wallpaper.icangoanyway.com \
    -d admin.icangoanyway.com \
    -d ios.icangoanyway.com \
    --email "$EMAIL" \
    --agree-tos \
    --no-eff-email \
    --force-renewal \
    --non-interactive

# 6. 拷贝新证书到 ssl 目录供 gateway 挂载使用
echo -e "\n${YELLOW}[3/4] 正在将权威证书同步至 gateway 证书目录...${NC}"

CERT_PATH="$SCRIPT_DIR/certbot-etc/live/wallpaper.icangoanyway.com"
if [ ! -f "$CERT_PATH/fullchain.pem" ]; then
    # 智能查找 live 目录下包含 fullchain.pem 的有效证书目录
    CERT_PATH=$(find "$SCRIPT_DIR/certbot-etc/live" -name "fullchain.pem" 2>/dev/null | head -n 1 | xargs -r dirname)
fi

if [ -n "$CERT_PATH" ] && [ -f "$CERT_PATH/fullchain.pem" ]; then
    cp -L "$CERT_PATH/fullchain.pem" "$SCRIPT_DIR/ssl/fullchain.pem"
    cp -L "$CERT_PATH/privkey.pem" "$SCRIPT_DIR/ssl/privkey.pem"
    chmod 600 "$SCRIPT_DIR/ssl/privkey.pem"
    chmod 644 "$SCRIPT_DIR/ssl/fullchain.pem"
    echo -e "${GREEN}      权威证书已成功就绪！来源目录: $CERT_PATH${NC}"
else
    echo -e "${RED}[警告] 未找到生成的证书文件，请检查域名解析是否生效！${NC}"
    exit 1
fi

# 7. 热重载 Nginx 网关使 SSL 立即生效
echo -e "${YELLOW}[4/4] 正在热重载 Nginx 网关应用最新 SSL 证书...${NC}"
docker compose exec -T wallpaper-gateway nginx -s reload

echo -e "\n${GREEN}========================================================${NC}"
echo -e "${GREEN}  🎉 SSL 权威证书已成功申请并无缝激活！${NC}"
echo -e "${GREEN}  前台商城主站: https://wallpaper.icangoanyway.com${NC}"
echo -e "${GREEN}  后台管理系统: https://admin.icangoanyway.com${NC}"
echo -e "${GREEN}  前台兼容域名: https://ios.icangoanyway.com${NC}"
echo -e "${GREEN}  所有 HTTP (80) 访问将自动 301 重定向至 HTTPS (443)${NC}"
echo -e "${GREEN}========================================================${NC}\n"

echo -e "${CYAN}提示: Let's Encrypt 证书有效期为 90 天，如需全自动续签，可在服务器 crontab 添加:${NC}"
echo -e "0 3 1 * * cd $SCRIPT_DIR && ./setup-ssl.sh > /dev/null 2>&1\n"
