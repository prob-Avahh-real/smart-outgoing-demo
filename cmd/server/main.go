package main

import (
	"log"
	"strconv"
	"time"

	"smart-outgoing-demo/internal/config"
	"smart-outgoing-demo/internal/handlers"
	"smart-outgoing-demo/internal/store"
	"smart-outgoing-demo/internal/websocket"
	"smart-outgoing-demo/pkg/server"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	vehicleStore := store.NewVehicleStore()

	// WebSocket hub
	hub := websocket.NewHub(vehicleStore)
	go hub.Run()

	// Initialize parking handler (mock for now)
	parkingHandler := handlers.NewParkingHandler(nil, cfg)

	// Initialize payment handler (mock for now)
	paymentHandler := handlers.NewPaymentHandler(nil, cfg)

	// Initialize traffic handler
	trafficHandler := handlers.NewTrafficHandler(cfg)

	// Start scheduled cache cleanup (every 10 minutes)
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			vehicleStore.CleanupCache()
			log.Println("Scheduled cache cleanup completed")
		}
	}()

	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API routes
	api := r.Group("/api")
	{
		// Public routes (no auth required)
		api.GET("/config", handlers.GetConfig(cfg))
		api.GET("/vehicles", handlers.GetVehicles(vehicleStore))
		api.GET("/algorithm/plan", handlers.PlanRoute(cfg))
		api.GET("/cache/stats", handlers.GetCacheStats(vehicleStore))
		api.GET("/simulation/status", handlers.GetSimulationStatus(hub))
		api.GET("/simulation/results", handlers.GetSimulationResults())
		api.GET("/amap/geocode", handlers.AMapGeocode(cfg))
		api.GET("/amap/driving", handlers.AMapDriving(cfg))
		api.GET("/amap/validate", handlers.AMapValidate(cfg))

		// Parking API routes
		api.POST("/parking/find", parkingHandler.FindParking)
		api.POST("/parking/reserve", parkingHandler.ReserveSpace)
		api.POST("/parking/session/start", parkingHandler.StartParkingSession)
		api.POST("/parking/avp/start", parkingHandler.StartAVPTask)
		api.POST("/parking/avp/summon", parkingHandler.SummonAVPTask)
		api.GET("/parking/avp/tasks/:id", parkingHandler.GetAVPTask)
		api.POST("/parking/avp/tasks/:id/cancel", parkingHandler.CancelAVPTask)
		api.GET("/parking/lots", parkingHandler.GetParkingLots)
		api.GET("/parking/lots/:id", parkingHandler.GetParkingLot)
		api.GET("/parking/lots/:id/spaces", parkingHandler.GetParkingSpaces)

		// City Brain Parking API routes (龙华智行项目)
		api.POST("/parking/citybrain/entry", parkingHandler.ReportCityBrainEntry)
		api.POST("/parking/citybrain/exit", parkingHandler.ReportCityBrainExit)
		api.POST("/parking/citybrain/heartbeat", parkingHandler.SendCityBrainHeartbeat)

		// Parking Pool Management routes (三级车场池)
		api.GET("/parking/pools/stats", parkingHandler.GetParkingPoolStats)
		api.POST("/parking/pools/lot", parkingHandler.AddLotToPool)
		api.POST("/parking/pools/recommend", parkingHandler.GetPoolRecommendation)
		api.POST("/parking/pools/diversion", parkingHandler.TriggerTrafficDiversion)

		// Payment API routes
		api.POST("/payment/process", paymentHandler.ProcessPayment)
		api.GET("/payment/:payment_id/status", paymentHandler.GetPaymentStatus)
		api.POST("/payment/refund", paymentHandler.RefundPayment)
		api.GET("/payment/user/payments", paymentHandler.GetUserPayments)
		api.GET("/payment/methods", paymentHandler.GetPaymentMethods)
		api.GET("/payment/stats", paymentHandler.GetPaymentStats)

		// Traffic API routes
		api.POST("/traffic/density", trafficHandler.UpdateZoneDensity)
		api.GET("/traffic/zones/status", trafficHandler.GetZoneStatus)
		api.GET("/traffic/zone/:zone_id/analyze", trafficHandler.AnalyzeTrafficPattern)
		api.POST("/traffic/zone/:zone_id/decision", trafficHandler.MakeSchedulingDecision)
		api.GET("/traffic/decisions/history", trafficHandler.GetSchedulingHistory)
		api.GET("/traffic/density/:zone_id", trafficHandler.AnalyzeDensity)
		api.GET("/traffic/density/:zone_id/data", trafficHandler.GetDensityData)
		api.GET("/traffic/density/all", trafficHandler.GetAllDensityData)
		api.POST("/traffic/density/threshold", trafficHandler.SetAlertThreshold)
		api.GET("/traffic/density/thresholds", trafficHandler.GetAlertThresholds)
		api.POST("/traffic/density/train", trafficHandler.TrainPredictionModel)

		// C-V2X API routes
		api.POST("/cv2x/vehicle", trafficHandler.RegisterV2XVehicle)
		api.DELETE("/cv2x/vehicle/:vehicle_id", trafficHandler.UnregisterV2XVehicle)
		api.PUT("/cv2x/vehicle/:vehicle_id/position", trafficHandler.UpdateV2XVehiclePosition)
		api.GET("/cv2x/vehicles", trafficHandler.GetV2XVehicles)
		api.GET("/cv2x/rsus", trafficHandler.GetV2XRSUs)
		api.GET("/cv2x/messages", trafficHandler.GetV2XMessages)
		api.POST("/cv2x/vehicle/:vehicle_id/bsm", trafficHandler.SendV2XBSM)
		api.POST("/cv2x/hazard", trafficHandler.SimulateTrafficHazard)
		api.GET("/cv2x/statistics", trafficHandler.GetV2XStatistics)

		// Simulation API routes
		api.POST("/simulation/vehicles/generate", trafficHandler.GenerateRandomVehicles)
		api.POST("/simulation/scenario/:type", trafficHandler.GenerateTrafficScenario)
		api.DELETE("/simulation/vehicles", trafficHandler.ClearVehicles)
		api.GET("/app/url", trafficHandler.GetAppURL)

		// Protected routes (auth required)
		protected := api.Group("")
		protected.Use(handlers.RequireAuth(cfg.GetAuthConfig()))
		{
			protected.PUT("/config", handlers.UpdateConfig(cfg))
			protected.POST("/vehicles", handlers.CreateVehicle(vehicleStore, hub, cfg))
			protected.POST("/vehicles/import", handlers.ImportVehiclesFromCSV(vehicleStore, hub, cfg))
			protected.PUT("/vehicles/:id/destination", handlers.SetDestination(vehicleStore, hub, cfg))
			protected.DELETE("/vehicles/:id", handlers.DeleteVehicle(vehicleStore, hub, cfg))
			protected.POST("/algorithm/schedule", handlers.ScheduleTasks(vehicleStore))
			protected.POST("/cache/cleanup", handlers.CleanupCache(vehicleStore))
			protected.POST("/simulation/run", handlers.RunSimulation())
			protected.POST("/simulation/batch", handlers.StartSimulationBatch())
			protected.POST("/simulation/stop", handlers.StopSimulation())
		}
	}

	// WebSocket route
	r.GET("/ws", func(c *gin.Context) {
		websocket.ServeWs(hub, c.Writer, c.Request)
	})

	// HTML routes
	r.GET("/", func(c *gin.Context) {
		c.File("./public/html/index.html")
	})

	r.GET("/parking", func(c *gin.Context) {
		c.File("./public/html/parking.html")
	})

	r.GET("/payment", func(c *gin.Context) {
		c.File("./public/html/payment.html")
	})
	r.GET("/amap", func(c *gin.Context) {
		c.File("./public/html/amap.original.html")
	})
	r.GET("/dashboard", func(c *gin.Context) {
		c.File("./public/html/dashboard.html")
	})

	// Static file serving
	r.Static("/public", "./public")
	r.Static("/css", "./public/css")
	r.Static("/js", "./public/js")
	r.Static("/html", "./public/html")

	// Fallback to index.html for SPA
	r.NoRoute(func(c *gin.Context) {
		c.File("./public/html/index.html")
	})

	// Start server with HTTPS support
	serverConfig := &server.ServerConfig{
		HTTPPort:  strconv.Itoa(cfg.HTTPPort),
		HTTPSPort: strconv.Itoa(cfg.HTTPSPort),
		EnableTLS: cfg.EnableTLS,
		TLS: &server.TLSCertConfig{
			CertFile: cfg.TLSCertFile,
			KeyFile:  cfg.TLSKeyFile,
		},
	}

	log.Printf("Starting server with HTTPS support...")
	log.Printf("HTTP Port: %s, HTTPS Port: %s, TLS Enabled: %v",
		serverConfig.HTTPPort, serverConfig.HTTPSPort, serverConfig.EnableTLS)

	if err := server.RunServer(r, serverConfig); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
