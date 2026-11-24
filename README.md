# smart-outgoing-demo（智慧交通与车路协同演示系统）

基于 Go 的智慧行车、智能泊车与 C-V2X 车路协同实时轨迹仿真演示系统（已归档于 `archive/from-home/smart-outgoing-demo`）。

## 核心特性

- **多车协同仿真**：一键生成 10–30 辆互不重叠的智慧行车与智能泊车车辆，支持自动避让与路径规划
- **C-V2X 与 RSU 看板**：实时统计路侧单元（RSU）与 V2X 通信消息（`/dashboard`）
- **停车场与支付闭环**：含车位引导（`/parking`）、计费支付（`/payment`）与高德地图联动（`/amap`）

## 目标平台（Target 判定）

- **移动端**：可选（支持手机扫码打开移动端地图视图）
- **智能硬件**：可选（可搁置）

## 快速运行与验证

```bash
# 启动 Go HTTP 服务端（默认 http://localhost:8080）
go run cmd/server/main.go

# 运行 Go 单元测试
go test ./...
```

详见 `DEMO_GUIDE.md` 与 `PARKING_FEATURE.md`。
