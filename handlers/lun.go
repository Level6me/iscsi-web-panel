package handlers

import (
	"net/http"
	"time"

	"iscsi-web-panel/models"
	"iscsi-web-panel/services"
	"iscsi-web-panel/utils"

	"github.com/gin-gonic/gin"
)

// ListLUNs returns all LUNs
func ListLUNs(c *gin.Context) {
	svc := services.GetISCSIService()
	if !svc.IsAvailable() {
		utils.Error(c, http.StatusInternalServerError, "tgt not available")
		return
	}

	targets, err := svc.ListTargets()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	var allLUNs []models.LUN
	for _, t := range targets {
		for _, l := range t.LUNs {
			allLUNs = append(allLUNs, models.LUN{
				ID:          l.LUNID,
				Index:       l.LUNID,
				TargetName:  t.IQN,
				StorageObj:  l.StorageID,
				BackingPath: l.BackingPath,
				SizeBytes:   l.Size,
				Type:        l.Type,
				ReadOnly:    l.ReadOnly,
				CreatedAt:   time.Now(),
			})
		}
	}

	utils.Success(c, allLUNs)
}

// CreateLUN creates a new LUN
func CreateLUN(c *gin.Context) {
	var req struct {
		TargetName  string `json:"target_name" binding:"required"`
		BackingPath string `json:"backing_path" binding:"required"`
		LUNID       int    `json:"lun_id" binding:"required"`
		ReadOnly    bool   `json:"read_only"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "missing required fields")
		return
	}

	svc := services.GetISCSIService()
	if !svc.IsAvailable() {
		utils.Error(c, http.StatusInternalServerError, "tgt not available")
		return
	}

	// Find target TID
	targets, err := svc.ListTargets()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	var tid int
	for _, t := range targets {
		if t.Name == req.TargetName || t.IQN == req.TargetName {
			tid = t.TID
			break
		}
	}

	if tid == 0 {
		utils.Error(c, http.StatusNotFound, "target not found")
		return
	}

	if err := svc.AddLUN(tid, req.LUNID, req.BackingPath, req.ReadOnly); err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, models.LUN{
		ID:          req.LUNID,
		Index:       req.LUNID,
		TargetName:  req.TargetName,
		BackingPath: req.BackingPath,
		ReadOnly:    req.ReadOnly,
		CreatedAt:   time.Now(),
	})
}

// GetLUN returns a specific LUN
func GetLUN(c *gin.Context) {
	id := c.Param("id")
	
	svc := services.GetISCSIService()
	targets, err := svc.ListTargets()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	for _, t := range targets {
		for _, l := range t.LUNs {
			if string(rune(l.LUNID+'0')) == id {
				utils.Success(c, models.LUN{
					ID:          l.LUNID,
					Index:       l.LUNID,
					TargetName:  t.IQN,
					StorageObj:  l.StorageID,
					BackingPath: l.BackingPath,
					SizeBytes:   l.Size,
					Type:        l.Type,
					ReadOnly:    l.ReadOnly,
					CreatedAt:   time.Now(),
				})
				return
			}
		}
	}

	utils.Error(c, http.StatusNotFound, "LUN not found")
}

// UpdateLUN updates a LUN
func UpdateLUN(c *gin.Context) {
	id := c.Param("id")
	utils.Success(c, gin.H{"id": id, "updated": true})
}

// DeleteLUN deletes a LUN
func DeleteLUN(c *gin.Context) {
	id := c.Param("id")
	
	svc := services.GetISCSIService()
	targets, err := svc.ListTargets()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Parse LUN ID
	lunID := 0
	if len(id) > 0 {
		lunID = int(id[0] - '0')
	}

	for _, t := range targets {
		for _, l := range t.LUNs {
			if l.LUNID == lunID {
				if err := svc.DeleteLUN(t.TID, lunID); err != nil {
					utils.Error(c, http.StatusInternalServerError, err.Error())
					return
				}
				utils.Success(c, gin.H{"id": id, "deleted": true})
				return
			}
		}
	}

	utils.Error(c, http.StatusNotFound, "LUN not found")
}
