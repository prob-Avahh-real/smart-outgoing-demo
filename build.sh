#!/bin/bash

# 智慧交通演示系统 - 跨平台构建脚本

echo "========================================="
echo "开始构建智慧交通演示系统（跨平台）"
echo "========================================="
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo -e "${RED}错误: 未找到Go环境${NC}"
    echo "请先安装Go 1.21或更高版本"
    exit 1
fi

echo -e "${GREEN}✓ Go环境检查通过${NC}"
echo "Go版本: $(go version)"
echo ""

# 定义目标平台
PLATFORMS=(
    "linux/amd64:smart-traffic-linux"
    "linux/arm64:smart-traffic-linux-arm64"
    "darwin/amd64:smart-traffic-macos-x64"
    "darwin/arm64:smart-traffic-macos-arm64"
    "windows/amd64:smart-traffic.exe"
)

# 创建构建目录
BUILD_DIR="build"
echo "创建构建目录: $BUILD_DIR"
mkdir -p $BUILD_DIR

# 编译所有平台
echo "开始跨平台编译..."
echo ""

for platform in "${PLATFORMS[@]}"; do
    IFS=':' read -r GOARCH OUTPUT_NAME <<< "$platform"
    GOOS="${GOARCH%%/*}"
    GOARCH="${GOARCH##*/}"
    
    OUTPUT_DIR="$BUILD_DIR/$GOOS-$GOARCH"
    mkdir -p "$OUTPUT_DIR"
    mkdir -p "$OUTPUT_DIR/public"
    mkdir -p "$OUTPUT_DIR/config"
    
    echo "编译 $GOOS/$GOARCH ..."
    
    GOOS=$GOOS GOARCH=$GOARCH go build -o "$OUTPUT_DIR/$OUTPUT_NAME" cmd/server/main.go
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ $GOOS/$GOARCH 编译成功${NC}"
    else
        echo -e "${RED}✗ $GOOS/$GOARCH 编译失败${NC}"
        exit 1
    fi
done

echo ""
echo "复制前端文件到所有平台..."
for platform in "${PLATFORMS[@]}"; do
    IFS=':' read -r GOARCH OUTPUT_NAME <<< "$platform"
    GOOS="${GOARCH%%/*}"
    GOARCH="${GOARCH##*/}"
    OUTPUT_DIR="$BUILD_DIR/$GOOS-$GOARCH"
    
    cp -r public/* "$OUTPUT_DIR/public/" 2>/dev/null || true
done
echo -e "${GREEN}✓ 前端文件复制成功${NC}"

echo "复制配置文件..."
for platform in "${PLATFORMS[@]}"; do
    IFS=':' read -r GOARCH OUTPUT_NAME <<< "$platform"
    GOOS="${GOARCH%%/*}"
    GOARCH="${GOARCH##*/}"
    OUTPUT_DIR="$BUILD_DIR/$GOOS-$GOARCH"
    
    if [ -f ".env" ]; then
        cp .env "$OUTPUT_DIR/config/.env.example"
    else
        cat > "$OUTPUT_DIR/config/.env.example" << EOF
ADMIN_TOKEN=demo_admin_token
AMAP_JS_KEY=your_amap_key
AMAP_SECURITY_CODE=your_amap_security_code
EOF
    fi
done
echo -e "${GREEN}✓ 配置文件复制成功${NC}"

echo "复制文档..."
for platform in "${PLATFORMS[@]}"; do
    IFS=':' read -r GOARCH OUTPUT_NAME <<< "$platform"
    GOOS="${GOARCH%%/*}"
    GOARCH="${GOARCH##*/}"
    OUTPUT_DIR="$BUILD_DIR/$GOOS-$GOARCH"
    
    cp README.md "$OUTPUT_DIR/" 2>/dev/null || true
    cp DEMO_GUIDE.md "$OUTPUT_DIR/" 2>/dev/null || true
done
echo -e "${GREEN}✓ 文档文件复制成功${NC}"

echo "复制测试脚本..."
for platform in "${PLATFORMS[@]}"; do
    IFS=':' read -r GOARCH OUTPUT_NAME <<< "$platform"
    GOOS="${GOARCH%%/*}"
    GOARCH="${GOARCH##*/}"
    OUTPUT_DIR="$BUILD_DIR/$GOOS-$GOARCH"
    
    mkdir -p "$OUTPUT_DIR/tests"
    cp tests/run_tests.sh "$OUTPUT_DIR/tests/"
    chmod +x "$OUTPUT_DIR/tests/run_tests.sh"
done
echo -e "${GREEN}✓ 测试脚本复制成功${NC}"

# 创建平台特定的启动脚本
echo "创建平台特定的启动脚本..."

for platform in "${PLATFORMS[@]}"; do
    IFS=':' read -r GOARCH OUTPUT_NAME <<< "$platform"
    GOOS="${GOARCH%%/*}"
    GOARCH="${GOARCH##*/}"
    OUTPUT_DIR="$BUILD_DIR/$GOOS-$GOARCH"
    
    if [[ "$GOOS" == "windows" ]]; then
        # Windows启动脚本
        cat > "$OUTPUT_DIR/start.bat" << EOF
@echo off
echo =========================================
echo 启动智慧交通演示系统
echo =========================================
echo.

REM 检查配置文件
if not exist "config\.env" (
    echo 未找到配置文件，使用默认配置...
    set ADMIN_TOKEN=demo_admin_token
    set AMAP_JS_KEY=45109d104b3c8d03a2c84175a7749241
    set AMAP_SECURITY_CODE=c552677838e5f5e71de92ce532c936bc
) else (
    echo 加载配置文件...
    for /f "tokens=*" %%a in ('type config\.env') do set %%a
)

REM 启动服务
echo 启动服务...
$OUTPUT_NAME
pause
EOF
    else
        # Unix-like启动脚本
        cat > "$OUTPUT_DIR/start.sh" << EOF
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
    export \$(cat config/.env | grep -v '^#' | xargs)
fi

# 启动服务
echo "启动服务..."
./$OUTPUT_NAME
EOF
        chmod +x "$OUTPUT_DIR/start.sh"
    fi
done
echo -e "${GREEN}✓ 启动脚本创建成功${NC}"

# 创建平台特定的README
echo "创建平台特定的README..."
for platform in "${PLATFORMS[@]}"; do
    IFS=':' read -r GOARCH OUTPUT_NAME <<< "$platform"
    GOOS="${GOARCH%%/*}"
    GOARCH="${GOARCH##*/}"
    OUTPUT_DIR="$BUILD_DIR/$GOOS-$GOARCH"
    
    cat > "$OUTPUT_DIR/README.md" << EOF
# 智慧交通演示系统 - $GOOS $GOARCH

## 系统要求

- 操作系统: $GOOS ($GOARCH)
- 现代浏览器（Chrome、Edge、Safari、Firefox）
- 网络连接（需要加载高德地图API）

## 快速启动

EOF

    if [[ "$GOOS" == "windows" ]]; then
        cat >> "$OUTPUT_DIR/README.md" << EOF
双击 \`start.bat\` 或在命令行中运行：
\`\`\`cmd
start.bat
\`\`\`

或直接运行：
\`\`\`cmd
$OUTPUT_NAME
\`\`\`
EOF
    else
        cat >> "$OUTPUT_DIR/README.md" << EOF
运行启动脚本：
\`\`\`bash
./start.sh
\`\`\`

或直接运行：
\`\`\`bash
./$OUTPUT_NAME
\`\`\`
EOF
    fi

    cat >> "$OUTPUT_DIR/README.md" << EOF

## 访问系统

启动后，在浏览器中访问：\`http://localhost:8080\`

## 移动端访问

1. 确保移动设备和服务器在同一网络
2. 在电脑上访问 \`http://localhost:8080\`
3. 点击"📱 二维码"按钮
4. 用手机扫描二维码
5. 在手机上查看高德地图

## 配置

编辑 \`config/.env\` 文件配置系统参数：

\`\`\`bash
ADMIN_TOKEN=demo_admin_token
AMAP_JS_KEY=your_amap_key
AMAP_SECURITY_CODE=your_amap_security_code
\`\`\`

## 测试

运行功能测试：
\`\`\`bash
./tests/run_tests.sh
\`\`\`

## 停止服务

EOF

    if [[ "$GOOS" == "windows" ]]; then
        cat >> "$OUTPUT_DIR/README.md" << EOF
按 \`Ctrl+C\` 停止服务，或关闭命令行窗口。
EOF
    else
        cat >> "$OUTPUT_DIR/README.md" << EOF
按 \`Ctrl+C\` 停止服务，或运行：
\`\`\`bash
pkill -f $OUTPUT_NAME
\`\`\`
EOF
    fi

    cat >> "$OUTPUT_DIR/README.md" << EOF

## 详细文档

请查看 \`DEMO_GUIDE.md\` 获取详细的使用说明。

## 目录结构

\`\`\`
$GOOS-$GOARCH/
├── $OUTPUT_NAME          # 可执行文件
EOF

    if [[ "$GOOS" == "windows" ]]; then
        cat >> "$OUTPUT_DIR/README.md" << EOF
├── start.bat             # 启动脚本（Windows）
EOF
    else
        cat >> "$OUTPUT_DIR/README.md" << EOF
├── start.sh              # 启动脚本（Unix）
EOF
    fi

    cat >> "$OUTPUT_DIR/README.md" << EOF
├── public/               # 前端文件
│   ├── html/            # HTML页面
│   ├── css/             # 样式文件
│   └── js/              # JavaScript文件
├── config/              # 配置文件
│   └── .env.example     # 配置示例
├── tests/               # 测试脚本
│   └── run_tests.sh     # 功能测试
├── README.md            # 本文件
└── DEMO_GUIDE.md       # 详细指南
\`\`\`

## 移动端兼容性

本系统支持以下移动端浏览器：
- iOS Safari
- Android Chrome
- 微信内置浏览器
- 其他现代移动浏览器

## 故障排除

### 端口被占用
如果8080端口被占用，可以修改代码中的端口号或停止占用该端口的程序。

### 地图加载失败
1. 检查网络连接
2. 确认高德地图API密钥有效
3. 尝试刷新页面

### 移动端无法访问
1. 确保移动设备和服务器在同一网络
2. 检查防火墙设置
3. 使用服务器的局域网IP地址访问
EOF
done
echo -e "${GREEN}✓ README创建成功${NC}"

# 打包完成
echo ""
echo "========================================="
echo "跨平台构建完成"
echo "========================================="
echo ""
echo "构建目录: $BUILD_DIR"
echo ""
echo "已构建的平台："
for platform in "${PLATFORMS[@]}"; do
    IFS=':' read -r GOARCH OUTPUT_NAME <<< "$platform"
    GOOS="${GOARCH%%/*}"
    GOARCH="${GOARCH##*/}"
    echo "  - $GOOS/$GOARCH: $BUILD_DIR/$GOOS-$GOARCH/"
done
echo ""
echo "使用方法："
echo "  cd $BUILD_DIR/<platform>"
echo "  ./start.sh  # 或 start.bat (Windows)"
echo ""
echo -e "${GREEN}跨平台构建成功！${NC}"
