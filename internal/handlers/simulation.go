package handlers

import (
	"net/http"

	"smart-outgoing-demo/internal/simulation"
	"smart-outgoing-demo/internal/websocket"

	"github.com/gin-gonic/gin"
)

// RunSimulation 运行交通仿真
func RunSimulation() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req simulation.SimulationRequest

		// 解析JSON请求
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 验证请求
		bridge := simulation.NewSimulationBridge()
		if err := bridge.ValidateRequest(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 运行仿真
		result, err := bridge.RunSimulation(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Simulation failed: " + err.Error(),
				"success": false,
			})
			return
		}

		// 返回结果
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    result,
		})
	}
}

// GetSimulationStatus 获取仿真状态
func GetSimulationStatus(hub *websocket.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 返回当前WebSocket连接数和车辆数
		connectionCount := hub.GetClientCount()
		vehicleCount := hub.GetVehicleCount()

		status := gin.H{
			"websocket_connections": connectionCount,
			"active_vehicles":       vehicleCount,
			"simulation_ready":      true,
		}

		c.JSON(http.StatusOK, status)
	}
}

// StartSimulationBatch 批量启动仿真
func StartSimulationBatch() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 预定义场景
		scenarios := []simulation.ScenarioConfig{
			{Start: simulation.Point{X: 0, Y: 0}, Goal: simulation.Point{X: 100, Y: 100}},
			{Start: simulation.Point{X: 100, Y: 0}, Goal: simulation.Point{X: 0, Y: 100}},
			{Start: simulation.Point{X: 0, Y: 100}, Goal: simulation.Point{X: 100, Y: 100}},
			{Start: simulation.Point{X: 100, Y: 100}, Goal: simulation.Point{X: 0, Y: 0}},
		}

		// 预定义车辆
		vehicles := []simulation.VehicleConfig{
			{ID: "veh_0", Start: simulation.Point{X: 0, Y: 0}, Goal: simulation.Point{X: 100, Y: 100}, Speed: 8.0},
			{ID: "veh_1", Start: simulation.Point{X: 100, Y: 0}, Goal: simulation.Point{X: 0, Y: 100}, Speed: 10.0},
		}

		req := simulation.SimulationRequest{
			Scenarios: scenarios,
			Vehicles:  vehicles,
			MaxSteps:  2000,
			EnableGIF: true,
		}

		// 运行仿真
		bridge := simulation.NewSimulationBridge()
		result, err := bridge.RunSimulation(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    result,
		})
	}
}

// StopSimulation 停止仿真
func StopSimulation() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Simulation stopped",
			"success": true,
		})
	}
}

// GetSimulationResults 获取仿真结果历史
func GetSimulationResults() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这里可以从数据库或文件系统读取历史结果
		// 暂时返回模拟数据
		results := []gin.H{
			{
				"id":            "sim_001",
				"timestamp":     "2026-04-20T15:00:00Z",
				"vehicles":      2,
				"arrived_count": 2,
				"total_steps":   1850,
				"success_rate":  "100%",
			},
			{
				"id":            "sim_002",
				"timestamp":     "2026-04-20T14:30:00Z",
				"vehicles":      4,
				"arrived_count": 3,
				"total_steps":   2100,
				"success_rate":  "75%",
			},
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"count":   len(results),
			"results": results,
		})
	}
}
