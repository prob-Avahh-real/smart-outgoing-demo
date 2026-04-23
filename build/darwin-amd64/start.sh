#!/bin/bash

echo "========================================="
echo "启动智慧交通演示系统"
echo "========================================="
echo ""

# 检查配置文件
if [ ! -f "config/.env" ]; then
    echo "未找到配置文件，使用默认配置..."
    export ADMIN_TOKEN=demo_admin_token
    export AMAP_JS_KEY=45109d104b3c8d03a2c84175a7749241
    export AMAP_SECURITY_CODE=c552677838e5f5e71de92ce532c936bc
else
    echo "加载配置文件..."
    export $(cat config/.env | grep -v '^#' | xargs)
fi

# 启动服务
echo "启动服务..."
./smart-traffic-macos-x64
