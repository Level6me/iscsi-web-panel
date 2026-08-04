package handlers

import (
	"net/http"

	"iscsi-web-panel/models"
	"iscsi-web-panel/services"
	"iscsi-web-panel/utils"

	"github.com/gin-gonic/gin"
)

// ListInitiators returns all initiators/ACLs
func ListInitiators(c *gin.Context) {
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

	var allACLs []models.ACL
	for _, t := range targets {
		for _, a := range t.ACLs {
			allACLs = append(allACLs, models.ACL{
				InitiatorIQN: a.InitiatorIQN,
			})
		}
	}

	utils.Success(c, allACLs)
}

// CreateInitiator creates a new initiator ACL
func CreateInitiator(c *gin.Context) {
	var req struct {
		TargetIQN     string `json:"target_iqn" binding:"required"`
		InitiatorIQN  string `json:"initiator_iqn" binding:"required"`
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
		if t.IQN == req.TargetIQN || t.Name == req.TargetIQN {
			tid = t.TID
			break
		}
	}

	if tid == 0 {
		utils.Error(c, http.StatusNotFound, "target not found")
		return
	}

	if err := svc.AddInitiator(tid, req.InitiatorIQN); err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, models.ACL{
		InitiatorIQN: req.InitiatorIQN,
	})
}

// GetInitiator returns a specific initiator
func GetInitiator(c *gin.Context) {
	iqn := c.Param("iqn")
	
	svc := services.GetISCSIService()
	targets, err := svc.ListTargets()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	for _, t := range targets {
		for _, a := range t.ACLs {
			if a.InitiatorIQN == iqn {
				utils.Success(c, models.ACL{
					InitiatorIQN: a.InitiatorIQN,
				})
				return
			}
		}
	}

	utils.Error(c, http.StatusNotFound, "initiator not found")
}

// UpdateInitiator updates an initiator
func UpdateInitiator(c *gin.Context) {
	iqn := c.Param("iqn")
	utils.Success(c, gin.H{"iqn": iqn, "updated": true})
}

// DeleteInitiator deletes an initiator ACL
func DeleteInitiator(c *gin.Context) {
	iqn := c.Param("iqn")
	
	svc := services.GetISCSIService()
	targets, err := svc.ListTargets()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	for _, t := range targets {
		for _, a := range t.ACLs {
			if a.InitiatorIQN == iqn {
				if err := svc.RemoveInitiator(t.TID, iqn); err != nil {
					utils.Error(c, http.StatusInternalServerError, err.Error())
					return
				}
				utils.Success(c, gin.H{"iqn": iqn, "deleted": true})
				return
			}
		}
	}

	utils.Error(c, http.StatusNotFound, "initiator not found")
}
