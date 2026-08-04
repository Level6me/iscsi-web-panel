package handlers

import (
	"time"

	"iscsi-web-panel/models"
	"iscsi-web-panel/utils"

	"github.com/gin-gonic/gin"
)

// ListLogs returns log entries
func ListLogs(c *gin.Context) {
	mockLogs := []models.LogEntry{
		{ID: 1, Level: "info", Message: "System started", Source: "system", Timestamp: time.Now().Add(-7 * 24 * time.Hour)},
		{ID: 2, Level: "info", Message: "Target target1 created", Source: "target", Timestamp: time.Now().Add(-5 * 24 * time.Hour)},
		{ID: 3, Level: "warning", Message: "High latency detected on LUN 2", Source: "monitor", Timestamp: time.Now().Add(-2 * time.Hour)},
		{ID: 4, Level: "info", Message: "Initiator initiator1 connected", Source: "connection", Timestamp: time.Now().Add(-1 * time.Hour)},
		{ID: 5, Level: "error", Message: "Failed to create portal on eth2", Source: "network", Timestamp: time.Now().Add(-30 * time.Minute)},
	}
	utils.Success(c, mockLogs)
}

// ExportLogs exports logs as text
func ExportLogs(c *gin.Context) {
	c.Header("Content-Disposition", "attachment; filename=iscsi-logs.txt")
	c.String(200, `[INFO] 2024-01-01 10:00:00 System started
[INFO] 2024-01-02 08:00:00 Target target1 created
[WARN] 2024-01-07 14:00:00 High latency detected on LUN 2
[INFO] 2024-01-07 15:00:00 Initiator initiator1 connected
[ERROR] 2024-01-07 15:30:00 Failed to create portal on eth2
`)
}
