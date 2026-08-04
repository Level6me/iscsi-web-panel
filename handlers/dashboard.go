package handlers

import (
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"iscsi-web-panel/services"
	"iscsi-web-panel/utils"

	"github.com/gin-gonic/gin"
)

// DashboardOverview returns system overview
func DashboardOverview(c *gin.Context) {
	svc := services.GetISCSIService()
	
	// Get target count
	targets, _ := svc.ListTargets()
	targetCount := len(targets)
	activeTargets := 0
	for _, t := range targets {
		if t.Status == "active" {
			activeTargets++
		}
	}
	
	// Get LUN count
	lunCount := 0
	for _, t := range targets {
		lunCount += len(t.LUNs)
	}
	
	// Get storage info
	totalStorage, usedStorage := getStorageInfo()
	
	utils.Success(c, gin.H{
		"targets":        targetCount,
		"active_targets": activeTargets,
		"luns":           lunCount,
		"connections":    0, // TODO: Get from tgt
		"storage_total":  totalStorage,
		"storage_used":   usedStorage,
	})
}

func getStorageInfo() (total, used int64) {
	// Get disk usage for root filesystem
	var stat syscall.Statfs_t
	if err := syscall.Statfs("/", &stat); err == nil {
		total = int64(stat.Blocks) * int64(stat.Bsize)
		used = int64(stat.Blocks-stat.Bfree) * int64(stat.Bsize)
	}
	return total, used
}

// DashboardStats returns system stats
func DashboardStats(c *gin.Context) {
	cpuUsage := getCPUUsage()
	memTotal, memUsed := getMemoryInfo()
	diskTotal, diskUsed := getStorageInfo()
	
	utils.Success(c, gin.H{
		"cpu_usage":  cpuUsage,
		"mem_total":  memTotal,
		"mem_used":   memUsed,
		"disk_total": diskTotal,
		"disk_used":  diskUsed,
		"uptime":     getUptime(),
	})
}

func getCPUUsage() float64 {
	// Read from /proc/stat
	cmd := exec.Command("cat", "/proc/stat")
	out, err := cmd.Output()
	if err != nil {
		return 0
	}
	
	lines := strings.Split(string(out), "\n")
	if len(lines) == 0 {
		return 0
	}
	
	fields := strings.Fields(lines[0])
	if len(fields) < 5 {
		return 0
	}
	
	// Parse CPU times
	var user, nice, system, idle int64
	user, _ = strconv.ParseInt(fields[1], 10, 64)
	nice, _ = strconv.ParseInt(fields[2], 10, 64)
	system, _ = strconv.ParseInt(fields[3], 10, 64)
	idle, _ = strconv.ParseInt(fields[4], 10, 64)
	
	total := user + nice + system + idle
	if total == 0 {
		return 0
	}
	
	// Calculate usage percentage (simplified)
	usage := float64(user+nice+system) / float64(total) * 100
	return usage
}

func getMemoryInfo() (total, used int64) {
	cmd := exec.Command("cat", "/proc/meminfo")
	out, err := cmd.Output()
	if err != nil {
		return 0, 0
	}
	
	lines := strings.Split(string(out), "\n")
	var memTotal, memAvailable int64
	
	for _, line := range lines {
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				memTotal, _ = strconv.ParseInt(fields[1], 10, 64)
			}
		} else if strings.HasPrefix(line, "MemAvailable:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				memAvailable, _ = strconv.ParseInt(fields[1], 10, 64)
			}
		}
	}
	
	// Convert from KB to bytes
	total = memTotal * 1024
	used = (memTotal - memAvailable) * 1024
	
	return total, used
}

func getUptime() int64 {
	cmd := exec.Command("cat", "/proc/uptime")
	out, err := cmd.Output()
	if err != nil {
		return 0
	}
	
	fields := strings.Fields(string(out))
	if len(fields) < 1 {
		return 0
	}
	
	uptime, _ := strconv.ParseFloat(fields[0], 64)
	return int64(uptime)
}
