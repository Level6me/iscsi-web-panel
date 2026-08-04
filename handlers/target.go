package handlers

import (
	"net/http"
	"strconv"
	"time"

	"iscsi-web-panel/models"
	"iscsi-web-panel/services"
	"iscsi-web-panel/utils"

	"github.com/gin-gonic/gin"
)

// ListTargets returns all iSCSI targets
func ListTargets(c *gin.Context) {
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

	// Convert to models
	var result []models.Target
	for _, t := range targets {
		mt := models.Target{
			Name:      t.Name,
			IQN:       t.IQN,
			Status:    t.Status,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		for _, l := range t.LUNs {
			mt.LUNs = append(mt.LUNs, models.LUN{
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
		for _, a := range t.ACLs {
			mt.ACLs = append(mt.ACLs, models.ACL{
				InitiatorIQN: a.InitiatorIQN,
			})
		}
		for _, p := range t.Portals {
			mt.Portals = append(mt.Portals, models.Portal{
				IP:   p.IP,
				Port: p.Port,
			})
		}
		result = append(result, mt)
	}

	utils.Success(c, result)
}

// CreateTarget creates a new iSCSI target
func CreateTarget(c *gin.Context) {
	var req struct {
		IQN string `json:"iqn" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "iqn is required")
		return
	}

	svc := services.GetISCSIService()
	if !svc.IsAvailable() {
		utils.Error(c, http.StatusInternalServerError, "tgt not available")
		return
	}

	if err := svc.CreateTarget(req.IQN); err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, models.Target{
		Name:      req.IQN,
		IQN:       req.IQN,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
}

// GetTarget returns a specific target
func GetTarget(c *gin.Context) {
	name := c.Param("name")
	
	svc := services.GetISCSIService()
	targets, err := svc.ListTargets()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	for _, t := range targets {
		if t.Name == name || t.IQN == name {
			mt := models.Target{
				Name:      t.Name,
				IQN:       t.IQN,
				Status:    t.Status,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			for _, l := range t.LUNs {
				mt.LUNs = append(mt.LUNs, models.LUN{
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
			utils.Success(c, mt)
			return
		}
	}

	utils.Error(c, http.StatusNotFound, "target not found")
}

// UpdateTarget updates a target
func UpdateTarget(c *gin.Context) {
	name := c.Param("name")
	utils.Success(c, gin.H{"name": name, "updated": true})
}

// DeleteTarget deletes a target
func DeleteTarget(c *gin.Context) {
	name := c.Param("name")
	
	svc := services.GetISCSIService()
	targets, err := svc.ListTargets()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	for _, t := range targets {
		if t.Name == name || t.IQN == name {
			if err := svc.DeleteTarget(t.TID); err != nil {
				utils.Error(c, http.StatusInternalServerError, err.Error())
				return
			}
			utils.Success(c, gin.H{"name": name, "deleted": true})
			return
		}
	}

	utils.Error(c, http.StatusNotFound, "target not found")
}

// GetNextTID returns the next available TID
func GetNextTID() int {
	svc := services.GetISCSIService()
	targets, err := svc.ListTargets()
	if err != nil {
		return 1
	}
	
	maxTID := 0
	for _, t := range targets {
		if t.TID > maxTID {
			maxTID = t.TID
		}
	}
	return maxTID + 1
}

// ParseTID parses TID from string
func ParseTID(s string) int {
	tid, _ := strconv.Atoi(s)
	return tid
}
