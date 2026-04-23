# 智慧交通演示系统

## 重要说明

**这不是一个可以直接打开的HTML文件**，这是一个需要运行服务器的完整应用程序。

### 为什么需要运行服务器？

- 系统包含后端API（Go语言编写）
- 需要WebSocket实时通信
- 需要加载高德地图API
- 车辆数据存储在服务器端

### 如何分享给其他人？

**方法1：分享整个build目录**
1. 将对应平台的目录（如 `darwin-arm64`）打包压缩
2. 发送给对方
3. 对方解压后运行 `./start.sh`（Windows是 `start.bat`）
4. 在浏览器访问 `http://localhost:8080`

**方法2：部署到服务器**
1. 将对应平台的文件上传到服务器
2. 运行服务
3. 其他人通过服务器IP地址访问

**方法3：本地演示**
1. 在你的电脑上运行服务
2. 其他人通过你的局域网IP访问（如 `http://192.168.1.100:8080`）
3. 或使用二维码功能让手机访问

## 快速启动

### 方法1：使用启动脚本
```bash
./start.sh
```

### 方法2：直接运行
```bash
./smart-traffic
```

### 方法3：使用Go运行
```bash
go run cmd/server/main.go
```

## 访问系统

启动后，在浏览器中访问：`http://localhost:8080`

## 跨平台版本

根据你的操作系统选择对应的目录：

- **macOS Intel**: `darwin-amd64/`
- **macOS Apple Silicon**: `darwin-arm64/`
- **Linux x64**: `linux-amd64/`
- **Linux ARM64**: `linux-arm64/`
- **Windows**: `windows-amd64/`

每个平台目录都包含完整的可执行文件和前端资源。

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

```bash
./stop.sh
```

## 详细文档

请查看 `DEMO_GUIDE.md` 获取详细的使用说明。

## 目录结构

```
build/
├── smart-traffic          # 可执行文件（当前平台）
├── start.sh              # 启动脚本
├── stop.sh               # 停止脚本
├── public/               # 前端文件
│   ├── html/            # HTML页面
│   ├── css/             # 样式文件
│   └── js/              # JavaScript文件
├── config/              # 配置文件
│   └── .env.example     # 配置示例
├── tests/               # 测试脚本
│   └── run_tests.sh     # 功能测试
├── darwin-amd64/        # macOS Intel版本
├── darwin-arm64/        # macOS Apple Silicon版本
├── linux-amd64/         # Linux x64版本
├── linux-arm64/         # Linux ARM64版本
├── windows-amd64/       # Windows版本
├── README.md            # 本文件
└── DEMO_GUIDE.md       # 详细指南
```

## 演示步骤

1. **启动服务**: `./start.sh`
2. **打开浏览器**: 访问 `http://localhost:8080`
3. **生成车辆**: 点击"创建示例车（10～30）"
4. **开始演示**: 点击"开始"按钮
5. **移动端**: 点击"📱 二维码"扫码访问
