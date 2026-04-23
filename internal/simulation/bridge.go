package simulation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

// SimulationRequest 仿真请求
type SimulationRequest struct {
	Vehicles  []VehicleConfig  `json:"vehicles"`
	Scenarios []ScenarioConfig `json:"scenarios"`
	MaxSteps  int              `json:"max_steps"`
	EnableGIF bool             `json:"enable_gif"`
}

// VehicleConfig 车辆配置
type VehicleConfig struct {
	ID    string  `json:"id"`
	Start Point   `json:"start"`
	Goal  Point   `json:"goal"`
	Speed float64 `json:"speed"`
}

// ScenarioConfig 场景配置
type ScenarioConfig struct {
	Start Point `json:"start"`
	Goal  Point `json:"goal"`
}

// Point 坐标点
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// SimulationResult 仿真结果
type SimulationResult struct {
	Success      bool                `json:"success"`
	TotalSteps   int                 `json:"total_steps"`
	ArrivedCount int                 `json:"arrived_count"`
	Trajectories []VehicleTrajectory `json:"trajectories"`
	GIFPath      string              `json:"gif_path,omitempty"`
	Error        string              `json:"error,omitempty"`
}

// VehicleTrajectory 车辆轨迹
type VehicleTrajectory struct {
	ID          string  `json:"id"`
	Trajectory  []Point `json:"trajectory"`
	Arrived     bool    `json:"arrived"`
	ArrivalStep int     `json:"arrival_step"`
}

// SimulationBridge Python仿真桥接服务
type SimulationBridge struct {
	pythonScript string
	timeout      time.Duration
}

// NewSimulationBridge 创建仿真桥接服务
func NewSimulationBridge() *SimulationBridge {
	// 检查Python脚本是否存在
	scriptPath := "simulation/traffic_simulator.py"
	return &SimulationBridge{
		pythonScript: scriptPath,
		timeout:      5 * time.Minute,
	}
}

// RunSimulation 运行仿真
func (s *SimulationBridge) RunSimulation(req SimulationRequest) (*SimulationResult, error) {
	// 准备Python脚本参数
	configData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 执行Python仿真脚本
	cmd := exec.Command("python3", s.pythonScript, string(configData))

	// 设置输入
	cmd.Stdin = bytes.NewReader([]byte(configData))

	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute simulation: %w", err)
	}

	// 解析结果
	var result SimulationResult
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return nil, fmt.Errorf("failed to parse simulation result: %w", err)
	}

	return &result, nil
}

// ValidateRequest 验证仿真请求
func (s *SimulationBridge) ValidateRequest(req SimulationRequest) error {
	if len(req.Vehicles) == 0 {
		return fmt.Errorf("at least one vehicle required")
	}

	if len(req.Scenarios) == 0 {
		return fmt.Errorf("at least one scenario required")
	}

	if req.MaxSteps <= 0 || req.MaxSteps > 10000 {
		return fmt.Errorf("max_steps must be between 1 and 10000")
	}

	// 验证坐标范围
	for _, vehicle := range req.Vehicles {
		if vehicle.Start.X < -1000 || vehicle.Start.X > 1000 ||
			vehicle.Start.Y < -1000 || vehicle.Start.Y > 1000 {
			return fmt.Errorf("vehicle start coordinates out of range: %v", vehicle.Start)
		}

		if vehicle.Goal.X < -1000 || vehicle.Goal.X > 1000 ||
			vehicle.Goal.Y < -1000 || vehicle.Goal.Y > 1000 {
			return fmt.Errorf("vehicle goal coordinates out of range: %v", vehicle.Goal)
		}
	}

	return nil
}
