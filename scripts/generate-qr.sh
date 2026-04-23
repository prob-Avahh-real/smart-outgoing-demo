#!/bin/bash

# QR Code Generator for AI Car Parking
# Generates QR codes for mobile access

set -e

OUTPUT_DIR="./public/qr-codes"
SERVER_URL="${SERVER_URL:-http://localhost:8080}"

# Create output directory
mkdir -p "$OUTPUT_DIR"

echo "=== AI Car Parking QR Code Generator ==="
echo "Server URL: $SERVER_URL"
echo ""

# Generate QR codes for different pages
echo "Generating QR codes..."

# Main parking page
qrencode -o "$OUTPUT_DIR/parking.png" -s 10 "$SERVER_URL/parking"
echo "✅ Parking page: $OUTPUT_DIR/parking.png"

# Payment page
qrencode -o "$OUTPUT_DIR/payment.png" -s 10 "$SERVER_URL/payment"
echo "✅ Payment page: $OUTPUT_DIR/payment.png"

# Main page
qrencode -o "$OUTPUT_DIR/main.png" -s 10 "$SERVER_URL/"
echo "✅ Main page: $OUTPUT_DIR/main.png"

# Create a demo poster
cat > "$OUTPUT_DIR/demo-poster.html" << 'EOF'
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>AI Car Parking Demo</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 40px 20px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
        }
        .poster {
            background: white;
            border-radius: 20px;
            padding: 40px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            text-align: center;
        }
        h1 {
            color: #333;
            margin-bottom: 10px;
            font-size: 2.5em;
        }
        .subtitle {
            color: #666;
            margin-bottom: 30px;
            font-size: 1.2em;
        }
        .qr-section {
            margin: 30px 0;
            padding: 30px;
            background: #f8f9fa;
            border-radius: 15px;
        }
        .qr-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 30px;
            margin-top: 20px;
        }
        .qr-item {
            text-align: center;
        }
        .qr-item img {
            max-width: 100%;
            height: auto;
            border: 3px solid #667eea;
            border-radius: 10px;
            padding: 10px;
            background: white;
        }
        .qr-item h3 {
            color: #333;
            margin-top: 15px;
            font-size: 1.1em;
        }
        .features {
            text-align: left;
            margin: 30px 0;
            padding: 30px;
            background: #e8f4f8;
            border-radius: 15px;
        }
        .features h2 {
            color: #667eea;
            margin-bottom: 20px;
        }
        .features ul {
            list-style: none;
            padding: 0;
        }
        .features li {
            padding: 10px 0;
            border-bottom: 1px solid #ddd;
        }
        .features li:before {
            content: "✓ ";
            color: #667eea;
            font-weight: bold;
            margin-right: 10px;
        }
        .cta {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 20px;
            border-radius: 15px;
            margin-top: 30px;
            font-size: 1.3em;
            font-weight: bold;
        }
        .footer {
            margin-top: 30px;
            color: #666;
            font-size: 0.9em;
        }
    </style>
</head>
<body>
    <div class="poster">
        <h1>🚗 AI Car Parking</h1>
        <p class="subtitle">智能停车解决方案 - 一键找车位，轻松导航</p>
        
        <div class="qr-section">
            <h2>📱 扫码体验</h2>
            <p>使用手机扫描二维码，体验完整功能</p>
            <div class="qr-grid">
                <div class="qr-item">
                    <img src="parking.png" alt="Parking QR Code">
                    <h3>停车功能</h3>
                </div>
                <div class="qr-item">
                    <img src="payment.png" alt="Payment QR Code">
                    <h3>支付功能</h3>
                </div>
                <div class="qr-item">
                    <img src="main.png" alt="Main QR Code">
                    <h3>主页功能</h3>
                </div>
            </div>
        </div>

        <div class="features">
            <h2>🎯 核心功能</h2>
            <ul>
                <li>智能停车推荐 - 基于位置和偏好的智能匹配</li>
                <li>实时地图导航 - AMap高德地图集成</li>
                <li>一键预约支付 - 微信支付、支付宝、信用卡</li>
                <li>实时状态更新 - 停车位实时可用性</li>
                <li>会话管理 - 完整的停车会话追踪</li>
            </ul>
        </div>

        <div class="cta">
            🚀 随时随地上车，一脚油门进停车场！
        </div>

        <div class="footer">
            <p>AI Car Parking - 智能停车解决方案</p>
            <p>扫码立即体验未来停车方式</p>
        </div>
    </div>
</body>
</html>
EOF

echo "✅ Demo poster: $OUTPUT_DIR/demo-poster.html"

echo ""
echo "🎉 QR codes generated successfully!"
echo ""
echo "Files created:"
echo "  - $OUTPUT_DIR/parking.png"
echo "  - $OUTPUT_DIR/payment.png"
echo "  - $OUTPUT_DIR/main.png"
echo "  - $OUTPUT_DIR/demo-poster.html"
echo ""
echo "To use:"
echo "  1. Print the demo poster: $OUTPUT_DIR/demo-poster.html"
echo "  2. Display QR codes for users to scan"
echo "  3. Or use individual QR codes for specific features"
echo ""
echo "📱 Users can scan with any QR code reader app"
echo "🌐 Server URL: $SERVER_URL"
