package handlers

import (
	"iscsi-web-panel/models"
	"iscsi-web-panel/utils"

	"github.com/gin-gonic/gin"
)

// GetAPIDoc returns API documentation
func GetAPIDoc(c *gin.Context) {
	endpoints := []models.APIDocEndpoint{
		// Auth
		{Method: "POST", Path: "/api/v1/auth/login", Description: "User login, returns JWT token", Auth: false},
		{Method: "POST", Path: "/api/v1/auth/register", Description: "User registration", Auth: false},
		{Method: "GET", Path: "/api/v1/auth/me", Description: "Get current user info", Auth: true},

		// Dashboard
		{Method: "GET", Path: "/api/v1/dashboard/overview", Description: "System overview: targets, LUNs, connections, storage", Auth: true},
		{Method: "GET", Path: "/api/v1/dashboard/stats", Description: "System stats: CPU, memory, disk, network", Auth: true},

		// Targets
		{Method: "GET", Path: "/api/v1/targets", Description: "List all iSCSI targets", Auth: true},
		{Method: "POST", Path: "/api/v1/targets", Description: "Create a new iSCSI target", Auth: true},
		{Method: "GET", Path: "/api/v1/targets/:name", Description: "Get target details", Auth: true},
		{Method: "PUT", Path: "/api/v1/targets/:name", Description: "Update target", Auth: true},
		{Method: "DELETE", Path: "/api/v1/targets/:name", Description: "Delete target", Auth: true},

		// LUNs
		{Method: "GET", Path: "/api/v1/luns", Description: "List all LUNs", Auth: true},
		{Method: "POST", Path: "/api/v1/luns", Description: "Create a new LUN", Auth: true},
		{Method: "GET", Path: "/api/v1/luns/:id", Description: "Get LUN details", Auth: true},
		{Method: "PUT", Path: "/api/v1/luns/:id", Description: "Update LUN", Auth: true},
		{Method: "DELETE", Path: "/api/v1/luns/:id", Description: "Delete LUN", Auth: true},

		// Initiators
		{Method: "GET", Path: "/api/v1/initiators", Description: "List all initiators/ACLs", Auth: true},
		{Method: "POST", Path: "/api/v1/initiators", Description: "Create initiator ACL", Auth: true},
		{Method: "GET", Path: "/api/v1/initiators/:iqn", Description: "Get initiator details", Auth: true},
		{Method: "PUT", Path: "/api/v1/initiators/:iqn", Description: "Update initiator", Auth: true},
		{Method: "DELETE", Path: "/api/v1/initiators/:iqn", Description: "Delete initiator", Auth: true},

		// Network
		{Method: "GET", Path: "/api/v1/network/interfaces", Description: "List network interfaces", Auth: true},
		{Method: "GET", Path: "/api/v1/network/portals", Description: "List portals", Auth: true},
		{Method: "POST", Path: "/api/v1/network/portals", Description: "Create portal", Auth: true},
		{Method: "DELETE", Path: "/api/v1/network/portals/:addr", Description: "Delete portal", Auth: true},
		{Method: "GET", Path: "/api/v1/network/discovery", Description: "Get discovery auth settings", Auth: true},
		{Method: "PUT", Path: "/api/v1/network/discovery", Description: "Update discovery auth", Auth: true},

		// Storage
		{Method: "GET", Path: "/api/v1/storage/pools", Description: "List storage pools", Auth: true},
		{Method: "POST", Path: "/api/v1/storage/pools", Description: "Create storage pool", Auth: true},
		{Method: "GET", Path: "/api/v1/storage/pools/:name", Description: "Get storage pool details", Auth: true},
		{Method: "DELETE", Path: "/api/v1/storage/pools/:name", Description: "Delete storage pool", Auth: true},
		{Method: "GET", Path: "/api/v1/storage/devices", Description: "List available block devices", Auth: true},

		// Monitor
		{Method: "GET", Path: "/api/v1/monitor/metrics", Description: "Get current metrics", Auth: true},
		{Method: "GET", Path: "/api/v1/monitor/performance", Description: "Get performance samples", Auth: true},
		{Method: "GET", Path: "/api/v1/monitor/connections", Description: "Get active connections", Auth: true},

		// Alerts
		{Method: "GET", Path: "/api/v1/alerts", Description: "List alerts", Auth: true},
		{Method: "PUT", Path: "/api/v1/alerts/:id/ack", Description: "Acknowledge alert", Auth: true},
		{Method: "GET", Path: "/api/v1/alerts/rules", Description: "List alert rules", Auth: true},
		{Method: "POST", Path: "/api/v1/alerts/rules", Description: "Create alert rule", Auth: true},
		{Method: "DELETE", Path: "/api/v1/alerts/rules/:id", Description: "Delete alert rule", Auth: true},

		// Logs
		{Method: "GET", Path: "/api/v1/logs", Description: "List log entries", Auth: true},
		{Method: "GET", Path: "/api/v1/logs/export", Description: "Export logs as text file", Auth: true},

		// Users
		{Method: "GET", Path: "/api/v1/users", Description: "List users", Auth: true},
		{Method: "POST", Path: "/api/v1/users", Description: "Create user", Auth: true},
		{Method: "GET", Path: "/api/v1/users/:id", Description: "Get user details", Auth: true},
		{Method: "PUT", Path: "/api/v1/users/:id", Description: "Update user", Auth: true},
		{Method: "DELETE", Path: "/api/v1/users/:id", Description: "Delete user", Auth: true},

		// Settings
		{Method: "GET", Path: "/api/v1/settings", Description: "Get system settings", Auth: true},
		{Method: "PUT", Path: "/api/v1/settings", Description: "Update system settings", Auth: true},
		{Method: "GET", Path: "/api/v1/settings/about", Description: "Get system info", Auth: true},

		// Snapshots
		{Method: "GET", Path: "/api/v1/snapshots", Description: "List snapshots", Auth: true},
		{Method: "POST", Path: "/api/v1/snapshots", Description: "Create snapshot", Auth: true},
		{Method: "GET", Path: "/api/v1/snapshots/:id", Description: "Get snapshot details", Auth: true},
		{Method: "DELETE", Path: "/api/v1/snapshots/:id", Description: "Delete snapshot", Auth: true},
		{Method: "POST", Path: "/api/v1/snapshots/:id/restore", Description: "Restore from snapshot", Auth: true},
	}

	utils.Success(c, gin.H{
		"title":       "iSCSI Web Panel API",
		"version":     "0.1.0",
		"base_url":    "/api/v1",
		"endpoints":   endpoints,
		"total_count": len(endpoints),
	})
}
