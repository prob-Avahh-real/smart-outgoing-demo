#!/bin/bash

# 简单的测试脚本 - 运行所有功能测试

echo "========================================="
echo "开始运行功能测试"
echo "========================================="
echo ""

BASE_URL="http://localhost:8080"
FAILED=0
PASSED=0

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试函数
test_endpoint() {
    local name=$1
    local method=$2
    local endpoint=$3
    local data=$4

    echo -n "测试: $name ... "

    if [ -z "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data")
    fi

    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')

    if [ "$http_code" -ge 200 ] && [ "$http_code" -lt 300 ]; then
        echo -e "${GREEN}通过${NC} (HTTP $http_code)"
        ((PASSED++))
    else
        echo -e "${RED}失败${NC} (HTTP $http_code)"
        echo "  响应: $body"
        ((FAILED++))
    fi
}

echo "1. 基础API测试"
echo "----------------------------------------"
test_endpoint "获取应用URL" "GET" "/api/app/url"
test_endpoint "获取车辆列表" "GET" "/api/vehicles"
test_endpoint "获取C-V2X车辆" "GET" "/api/cv2x/vehicles"
test_endpoint "获取RSU列表" "GET" "/api/cv2x/rsus"
test_endpoint "获取C-V2X统计" "GET" "/api/cv2x/statistics"
test_endpoint "获取C-V2X消息" "GET" "/api/cv2x/messages"
echo ""

echo "2. 模拟功能测试"
echo "----------------------------------------"
test_endpoint "生成随机车辆（10～50）" "POST" "/api/simulation/vehicles/generate" '{"count":25}'
test_endpoint "生成拥堵场景" "POST" "/api/simulation/scenario/congestion"
test_endpoint "生成正常场景" "POST" "/api/simulation/scenario/normal"
test_endpoint "生成紧急场景" "POST" "/api/simulation/scenario/emergency"
test_endpoint "清空所有车辆" "DELETE" "/api/simulation/vehicles"
echo ""

echo "3. 页面访问测试"
echo "----------------------------------------"
test_endpoint "首页" "GET" "/"
test_endpoint "Dashboard页面" "GET" "/dashboard"
test_endpoint "Parking页面" "GET" "/parking"
test_endpoint "Payment页面" "GET" "/payment"
test_endpoint "AMap页面" "GET" "/amap"
echo ""

echo "========================================="
echo "测试完成"
echo "========================================="
echo -e "通过: ${GREEN}$PASSED${NC}"
echo -e "失败: ${RED}$FAILED${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}所有测试通过！${NC}"
    exit 0
else
    echo -e "${RED}有 $FAILED 个测试失败${NC}"
    exit 1
fi
