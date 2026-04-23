# 智慧交通演示系统 - darwin arm64

## 系统要求

- 操作系统: darwin (arm64)
- 现代浏览器（Chrome、Edge、Safari、Firefox）
- 网络连接（需要加载高德地图API）

## 快速启动

运行启动脚本：
```bash
./start.sh
```

或直接运行：
```bash
./smart-traffic-macos-arm64
```

## 访问系统

启动后，在浏览器中访问：`http://localhost:8080`

## 移动端访问

1. 确保移动设备和服务器在同一网络
2. 在电脑上访问 `http://localhost:8080`
3. 点击"📱 二维码"按钮
4. 用手机扫描二维码
5. 在手机上查看高德地图

## 配置

编辑 `config/.env` 文件配置系统参数：

```bash
ADMIN_TOKEN=demo_admin_token
AMAP_JS_KEY=your_amap_key
AMAP_SECURITY_CODE=your_amap_security_code
```

## 测试

运行功能测试：
```bash
./tests/run_tests.sh
```

## 停止服务

按 `Ctrl+C` 停止服务，或运行：
```bash
pkill -f smart-traffic-macos-arm64
```

## 详细文档

请查看 `DEMO_GUIDE.md` 获取详细的使用说明。

## 目录结构

```
darwin-arm64/
├── smart-traffic-macos-arm64          # 可执行文件
├── start.sh              # 启动脚本（Unix）
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
```

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
