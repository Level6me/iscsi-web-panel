package handlers

import (
	"net/http"
	"time"

	"iscsi-web-panel/models"
	"iscsi-web-panel/utils"

	"github.com/gin-gonic/gin"
)

// ListSnapshots returns all snapshots
func ListSnapshots(c *gin.Context) {
	mockSnapshots := []models.Snapshot{
		{
			ID:          1,
			Name:        "snap-target1-lun0-20240107",
			LUNID:       1,
			TargetName:  "iqn.2024-01.com.example:target1",
			SizeBytes:   10737418240,
			Description: "Daily backup",
			CreatedAt:   time.Now().Add(-24 * time.Hour),
		},
		{
			ID:          2,
			Name:        "snap-target2-lun0-20240106",
			LUNID:       2,
			TargetName:  "iqn.2024-01.com.example:target2",
			SizeBytes:   21474836480,
			Description: "Pre-update snapshot",
			CreatedAt:   time.Now().Add(-48 * time.Hour),
		},
	}
	utils.Success(c, mockSnapshots)
}

// CreateSnapshot creates a new snapshot
func CreateSnapshot(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		LUNID       int    `json:"lun_id" binding:"required"`
		TargetName  string `json:"target_name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "missing required fields")
		return
	}
	utils.Success(c, models.Snapshot{
		ID:          10,
		Name:        req.Name,
		LUNID:       req.LUNID,
		TargetName:  req.TargetName,
		SizeBytes:   10737418240,
		Description: req.Description,
		CreatedAt:   time.Now(),
	})
}

// GetSnapshot returns a specific snapshot
func GetSnapshot(c *gin.Context) {
	id := c.Param("id")
	utils.Success(c, models.Snapshot{
		ID:          1,
		Name:        "snap-target1-lun0-20240107",
		LUNID:       1,
		TargetName:  "iqn.2024-01.com.example:target1",
		SizeBytes:   10737418240,
		Description: "Daily backup",
		CreatedAt:   time.Now().Add(-24 * time.Hour),
	})
	_ = id
}

// DeleteSnapshot deletes a snapshot
func DeleteSnapshot(c *gin.Context) {
	id := c.Param("id")
	utils.Success(c, gin.H{"id": id, "deleted": true})
}

// RestoreSnapshot restores from a snapshot
func RestoreSnapshot(c *gin.Context) {
	id := c.Param("id")
	utils.Success(c, gin.H{"id": id, "restored": true})
}
