package handlers

import (
	"net/http"
	"os/exec"
	"strings"

	"iscsi-web-panel/models"
	"iscsi-web-panel/utils"

	"github.com/gin-gonic/gin"
)

// GetSettings returns system settings
func GetSettings(c *gin.Context) {
	// For now, return default settings
	// TODO: Load from config file or database
	utils.Success(c, models.Settings{
		SystemName:       "iSCSI Web Panel",
		DefaultCHAPUser:  "iscsi_user",
		AutoSnapshot:     false,
		SnapshotInterval: 60,
		LogRetentionDays: 30,
	})
}

// UpdateSettings updates system settings
func UpdateSettings(c *gin.Context) {
	var req models.Settings
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid request")
		return
	}
	
	// TODO: Save to config file or database
	utils.Success(c, req)
}

// GetAbout returns system info
func GetAbout(c *gin.Context) {
	kernelVer := getKernelVersion()
	tgtVer := getTgtVersion()
	
	utils.Success(c, gin.H{
		"version":           "0.2.0",
		"go_version":        "1.21",
		"kernel_version":    kernelVer,
		"tgt_version":       tgtVer,
		"platform":          "linux/arm64",
		"iscsi_available":   isTgtAvailable(),
	})
}

func getKernelVersion() string {
	cmd := exec.Command("uname", "-r")
	out, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

func getTgtVersion() string {
	cmd := exec.Command("tgtadm", "--version")
	out, err := cmd.Output()
	if err != nil {
		return "not installed"
	}
	return strings.TrimSpace(string(out))
}

func isTgtAvailable() bool {
	_, err := exec.LookPath("tgtadm")
	return err == nil
}
