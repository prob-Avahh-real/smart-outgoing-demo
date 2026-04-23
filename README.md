# 高德多智能体互动演示

这是一个用于现场演示的最小骨架：

- 前端使用高德 JavaScript API v2 展示地图、车辆、终点和路线
- 后端提供车辆实体管理接口（已落地到本地数据库文件）
- 后端提供 WebSocket 实时推送车辆快照
- 后端提供高德 Web 服务代理接口，避免前端直接调用 Web Service
- 算法逻辑暂时留空，先用占位路线把前后端链路接通，并预留独立算法接入层

## 启动

```bash
npm install
npm start
```

默认访问：

```text
http://localhost:3000
```

## 环境变量

启动前请显式设置（不再内置默认 key）：

```bash
export AMAP_WEB_KEY=你的Web服务Key
export AMAP_JS_KEY=你的JSAPIKey
export AMAP_SECURITY_CODE=你的安全密钥
export DEMO_ADMIN_TOKEN=你的写接口令牌
npm start
```

可选限流参数：

```bash
export WRITE_RATE_LIMIT_WINDOW_MS=60000
export WRITE_RATE_LIMIT_MAX=30
```

## 当前接口

- `GET /api/health`：健康检查
- `GET /api/config`：返回前端初始化地图所需配置
- `GET /api/vehicles`：获取全部车辆
- `POST /api/vehicles`：创建车辆实体（需令牌）
- `PUT /api/vehicles/:id/destination`：更新车辆目的地，并尝试调用高德驾车路线（需令牌）
- `POST /api/algorithm/plan`：算法接入层占位接口，可预览路径规划结果（需令牌）
- `POST /api/amap/geocode`：地址转坐标（需令牌）
- `POST /api/amap/route/driving`：驾车路线代理（需令牌）

写接口鉴权头二选一：

- `x-admin-token: <DEMO_ADMIN_TOKEN>`
- `Authorization: Bearer <DEMO_ADMIN_TOKEN>`

## 实时同步

- WebSocket 地址：`/ws`
- 服务端在创建车辆、修改终点后会主动广播完整车辆快照
- 前端优先使用 WebSocket，同步断开时自动回退到轮询

## 数据存储

- 当前采用 `NeDB` 文件数据库，数据文件在 `data/vehicles.db`
- 服务重启后车辆不会丢失
- 后续如果要上真正多机部署，建议将 `services/vehicle-store.js` 替换为 Redis 或云数据库实现

## 后续接算法的建议

- 保留前端轮询和地图渲染逻辑不变
- 将 [`algorithm-adapter.js`](file:///Users/skat/smart-outgoing-demo/services/algorithm-adapter.js) 里的 `planRoute()` 替换为真实调度/路径规划结果
- 如果后续要做多机联动，可以把内存态 `vehicles` 换成 Redis 或数据库
