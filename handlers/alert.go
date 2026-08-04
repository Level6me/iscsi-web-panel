package handlers

import (
	"net/http"
	"time"

	"iscsi-web-panel/models"
	"iscsi-web-panel/utils"

	"github.com/gin-gonic/gin"
)

// ListAlerts returns alerts
func ListAlerts(c *gin.Context) {
	mockAlerts := []models.Alert{
		{
			ID:        1,
			Level:     "warning",
			Message:   "Storage pool pool1 usage above 80%",
			Source:    "storage",
			Acked:     false,
			CreatedAt: time.Now().Add(-2 * time.Hour),
		},
		{
			ID:        2,
			Level:     "info",
			Message:   "Target target1 created",
			Source:    "target",
			Acked:     true,
			CreatedAt: time.Now().Add(-24 * time.Hour),
		},
		{
			ID:        3,
			Level:     "critical",
			Message:   "Connection lost from initiator2",
			Source:    "connection",
			Acked:     false,
			CreatedAt: time.Now().Add(-30 * time.Minute),
		},
	}
	utils.Success(c, mockAlerts)
}

// AcknowledgeAlert acknowledges an alert
func AcknowledgeAlert(c *gin.Context) {
	id := c.Param("id")
	utils.Success(c, gin.H{"id": id, "acked": true})
}

// ListAlertRules returns alert rules
func ListAlertRules(c *gin.Context) {
	mockRules := []models.AlertRule{
		{ID: 1, Name: "High Storage Usage", Metric: "storage_usage", Operator: "gt", Threshold: 80.0, Enabled: true},
		{ID: 2, Name: "High Latency", Metric: "latency_ms", Operator: "gt", Threshold: 10.0, Enabled: true},
		{ID: 3, Name: "Connection Lost", Metric: "connection_state", Operator: "eq", Threshold: 0.0, Enabled: true},
	}
	utils.Success(c, mockRules)
}

// CreateAlertRule creates an alert rule
func CreateAlertRule(c *gin.Context) {
	var req models.AlertRule
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid request")
		return
	}
	req.ID = 10
	utils.Success(c, req)
}

// DeleteAlertRule deletes an alert rule
func DeleteAlertRule(c *gin.Context) {
	id := c.Param("id")
	utils.Success(c, gin.H{"id": id, "deleted": true})
}
