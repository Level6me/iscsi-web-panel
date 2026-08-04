package handlers

import (
	"os/exec"
	"strconv"
	"strings"
	"time"

	"iscsi-web-panel/models"
	"iscsi-web-panel/services"
	"iscsi-web-panel/utils"

	"github.com/gin-gonic/gin"
)

// GetMetrics returns current metrics
func GetMetrics(c *gin.Context) {
	svc := services.GetISCSIService()
	
	// Get target count
	targets, _ := svc.ListTargets()
	targetCount := len(targets)
	
	// Get LUN count
	lunCount := 0
	for _, t := range targets {
		lunCount += len(t.LUNs)
	}
	
	// Get connection count (simplified)
	connCount := 0
	
	// Get disk I/O stats
	txBytes, rxBytes := getDiskIOStats()
	
	utils.Success(c, models.Metrics{
		TxBytes:     txBytes,
		RxBytes:     rxBytes,
		IOPS:        0, // TODO: Calculate from /proc/diskstats
		LatencyMs:   0, // TODO: Calculate
		Connections: connCount,
		Targets:     targetCount,
		LUNs:        lunCount,
		Timestamp:   time.Now(),
	})
}

func getDiskIOStats() (txBytes, rxBytes int64) {
	// Read from /proc/diskstats
	cmd := exec.Command("cat", "/proc/diskstats")
	out, err := cmd.Output()
	if err != nil {
		return 0, 0
	}
	
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 14 {
			continue
		}
		
		// Skip non-disk devices
		devName := fields[2]
		if !strings.HasPrefix(devName, "sd") && !strings.HasPrefix(devName, "nvme") {
			continue
		}
		
		// Parse read/write sectors (assuming 512 bytes per sector)
		readSectors, _ := strconv.ParseInt(fields[5], 10, 64)
		writeSectors, _ := strconv.ParseInt(fields[9], 10, 64)
		
		rxBytes += readSectors * 512
		txBytes += writeSectors * 512
	}
	
	return txBytes, rxBytes
}

// GetPerformance returns performance samples
func GetPerformance(c *gin.Context) {
	// For now, return empty samples
	// TODO: Implement performance sampling with historical data
	samples := []models.PerformanceSample{}
	utils.Success(c, samples)
}

// GetConnections returns active connections
func GetConnections(c *gin.Context) {
	svc := services.GetISCSIService()
	
	// Get connections from tgt
	conns, err := svc.GetConnections()
	if err != nil {
		// Return empty list on error
		utils.Success(c, []models.ConnectionInfo{})
		return
	}
	
	var result []models.ConnectionInfo
	for _, conn := range conns {
		result = append(result, models.ConnectionInfo{
			InitiatorIQN: getString(conn, "initiator_iqn"),
			TargetIQN:    getString(conn, "target_iqn"),
			SessionID:    getString(conn, "session_id"),
			State:        getString(conn, "status"),
			ConnectedAt:  time.Now(),
			TxBytes:      getInt64(conn, "tx_bytes"),
			RxBytes:      getInt64(conn, "rx_bytes"),
		})
	}
	
	utils.Success(c, result)
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getInt64(m map[string]interface{}, key string) int64 {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case int64:
			return val
		case float64:
			return int64(val)
		case int:
			return int64(val)
		}
	}
	return 0
}
