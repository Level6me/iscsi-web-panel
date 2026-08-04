package main

import (
	"iscsi-web-panel/config"
	"iscsi-web-panel/handlers"
	"iscsi-web-panel/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	r := gin.Default()

	// CORS
	r.Use(middleware.CORS())

	// Serve frontend static files
	r.Static("/assets", "./frontend/assets")
	r.StaticFile("/", "./frontend/index.html")
	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/index.html")
	})

	// API v1
	api := r.Group("/api/v1")
	{
		// Auth (login)
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
			auth.POST("/register", handlers.Register)
			auth.GET("/me", middleware.JWTAuth(), handlers.GetCurrentUser)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.JWTAuth())
		{
			// Dashboard
			protected.GET("/dashboard/overview", handlers.DashboardOverview)
			protected.GET("/dashboard/stats", handlers.DashboardStats)

			// Targets
			targets := protected.Group("/targets")
			{
				targets.GET("", handlers.ListTargets)
				targets.POST("", handlers.CreateTarget)
				targets.GET("/:name", handlers.GetTarget)
				targets.PUT("/:name", handlers.UpdateTarget)
				targets.DELETE("/:name", handlers.DeleteTarget)
			}

			// LUNs
			luns := protected.Group("/luns")
			{
				luns.GET("", handlers.ListLUNs)
				luns.POST("", handlers.CreateLUN)
				luns.GET("/:id", handlers.GetLUN)
				luns.PUT("/:id", handlers.UpdateLUN)
				luns.DELETE("/:id", handlers.DeleteLUN)
			}

			// Initiators
			initiators := protected.Group("/initiators")
			{
				initiators.GET("", handlers.ListInitiators)
				initiators.POST("", handlers.CreateInitiator)
				initiators.GET("/:iqn", handlers.GetInitiator)
				initiators.PUT("/:iqn", handlers.UpdateInitiator)
				initiators.DELETE("/:iqn", handlers.DeleteInitiator)
			}

			// Network
			network := protected.Group("/network")
			{
				network.GET("/interfaces", handlers.ListInterfaces)
				network.GET("/portals", handlers.ListPortals)
				network.POST("/portals", handlers.CreatePortal)
				network.DELETE("/portals/:addr", handlers.DeletePortal)
				network.GET("/discovery", handlers.DiscoveryAuth)
				network.PUT("/discovery", handlers.UpdateDiscoveryAuth)
			}

			// Storage
			storage := protected.Group("/storage")
			{
				storage.GET("/pools", handlers.ListStoragePools)
				storage.POST("/pools", handlers.CreateStoragePool)
				storage.GET("/pools/:name", handlers.GetStoragePool)
				storage.DELETE("/pools/:name", handlers.DeleteStoragePool)
				storage.GET("/devices", handlers.ListStorageDevices)
			}

			// Monitor
			monitor := protected.Group("/monitor")
			{
				monitor.GET("/metrics", handlers.GetMetrics)
				monitor.GET("/performance", handlers.GetPerformance)
				monitor.GET("/connections", handlers.GetConnections)
			}

			// Alerts
			alerts := protected.Group("/alerts")
			{
				alerts.GET("", handlers.ListAlerts)
				alerts.PUT("/:id/ack", handlers.AcknowledgeAlert)
				alerts.GET("/rules", handlers.ListAlertRules)
				alerts.POST("/rules", handlers.CreateAlertRule)
				alerts.DELETE("/rules/:id", handlers.DeleteAlertRule)
			}

			// Logs
			logs := protected.Group("/logs")
			{
				logs.GET("", handlers.ListLogs)
				logs.GET("/export", handlers.ExportLogs)
			}

			// Users
			users := protected.Group("/users")
			{
				users.GET("", handlers.ListUsers)
				users.POST("", handlers.CreateUser)
				users.GET("/:id", handlers.GetUser)
				users.PUT("/:id", handlers.UpdateUser)
				users.DELETE("/:id", handlers.DeleteUser)
			}

			// Settings
			settings := protected.Group("/settings")
			{
				settings.GET("", handlers.GetSettings)
				settings.PUT("", handlers.UpdateSettings)
				settings.GET("/about", handlers.GetAbout)
			}

			// Snapshots
			snapshots := protected.Group("/snapshots")
			{
				snapshots.GET("", handlers.ListSnapshots)
				snapshots.POST("", handlers.CreateSnapshot)
				snapshots.GET("/:id", handlers.GetSnapshot)
				snapshots.DELETE("/:id", handlers.DeleteSnapshot)
				snapshots.POST("/:id/restore", handlers.RestoreSnapshot)
			}
		}

		// API Doc (public)
		api.GET("/api-doc", handlers.GetAPIDoc)
	}

	r.Run(cfg.ListenAddr)
}
