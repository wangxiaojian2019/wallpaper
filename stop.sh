#!/bin/bash
if docker compose version &> /dev/null; then
    docker compose down
elif command -v docker-compose &> /dev/null; then
    docker-compose down
fi
echo "WallArt 服务已停止！"
