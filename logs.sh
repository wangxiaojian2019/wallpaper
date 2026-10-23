#!/bin/bash
if docker compose version &> /dev/null; then
    docker compose logs -f "$@"
elif command -v docker-compose &> /dev/null; then
    docker-compose logs -f "$@"
fi
