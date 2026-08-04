package handlers

import (
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"

	"iscsi-web-panel/models"
	"iscsi-web-panel/utils"

	"github.com/gin-gonic/gin"
)

// parseJSON is a helper to parse JSON
func parseJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// ListInterfaces returns real network interfaces
func ListInterfaces(c *gin.Context) {
	interfaces, err := net.Interfaces()
	if err != nil {
		utils.Error(c, 500, "failed to get interfaces")
		return
	}

	var result []models.NetworkInterface
	for _, iface := range interfaces {
		// Skip loopback and virtual interfaces
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if strings.HasPrefix(iface.Name, "docker") || strings.HasPrefix(iface.Name, "br-") || strings.HasPrefix(iface.Name, "veth") {
			continue
		}

		ni := models.NetworkInterface{
			Name:   iface.Name,
			Status: "down",
			MAC:    iface.HardwareAddr.String(),
		}

		if iface.Flags&net.FlagUp != 0 {
			ni.Status = "up"
		}

		// Get IP address
		addrs, err := iface.Addrs()
		if err == nil {
			for _, addr := range addrs {
				if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
					if ipNet.IP.To4() != nil {
						ni.IP = ipNet.IP.String()
						ni.Netmask = net.IP(ipNet.Mask).String()
						break
					}
				}
			}
		}

		// Get speed (if available) - store in a comment or skip
		// ni.Speed = getInterfaceSpeed(iface.Name)

		result = append(result, ni)
	}

	utils.Success(c, result)
}

func getInterfaceSpeed(name string) string {
	cmd := exec.Command("cat", "/sys/class/net/"+name+"/speed")
	out, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out)) + " Mbps"
}

// ListPortals returns all portals
func ListPortals(c *gin.Context) {
	// For now, return default portal
	portals := []models.Portal{
		{IP: "0.0.0.0", Port: 3260},
	}
	utils.Success(c, portals)
}

// CreatePortal creates a portal
func CreatePortal(c *gin.Context) {
	var req struct {
		IP         string `json:"ip" binding:"required"`
		Port       int    `json:"port"`
		TargetName string `json:"target_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "ip is required")
		return
	}
	if req.Port == 0 {
		req.Port = 3260
	}
	
	// TODO: Implement portal creation via tgtadm
	utils.Success(c, models.Portal{IP: req.IP, Port: req.Port})
}

// DeletePortal deletes a portal
func DeletePortal(c *gin.Context) {
	addr := c.Param("addr")
	// TODO: Implement portal deletion via tgtadm
	utils.Success(c, gin.H{"addr": addr, "deleted": true})
}

// DiscoveryAuth returns discovery auth settings
func DiscoveryAuth(c *gin.Context) {
	utils.Success(c, gin.H{
		"enabled":  false,
		"username": "",
		"password": "",
	})
}

// UpdateDiscoveryAuth updates discovery auth
func UpdateDiscoveryAuth(c *gin.Context) {
	var req struct {
		Enabled  bool   `json:"enabled"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "invalid request")
		return
	}
	// TODO: Implement discovery auth via tgtadm
	utils.Success(c, gin.H{"updated": true})
}

// ListStoragePools returns storage pools (backstores)
func ListStoragePools(c *gin.Context) {
	// For now, return empty list
	// TODO: Parse tgtadm output to get backstores
	pools := []models.StoragePool{}
	utils.Success(c, pools)
}

// CreateStoragePool creates a storage pool
func CreateStoragePool(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
		Type string `json:"type" binding:"required"`
		Path string `json:"path" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "missing required fields")
		return
	}
	// TODO: Implement storage pool creation
	utils.Success(c, models.StoragePool{
		Name:      req.Name,
		Type:      req.Type,
		Path:      req.Path,
		CreatedAt: time.Now(),
	})
}

// GetStoragePool returns a specific pool
func GetStoragePool(c *gin.Context) {
	name := c.Param("name")
	// TODO: Get specific backstore
	utils.Success(c, models.StoragePool{
		Name: name,
		Type: "fileio",
		Path: "/data/" + name,
	})
}

// DeleteStoragePool deletes a storage pool
func DeleteStoragePool(c *gin.Context) {
	name := c.Param("name")
	// TODO: Implement storage pool deletion
	utils.Success(c, gin.H{"name": name, "deleted": true})
}

// ListStorageDevices returns available block devices
func ListStorageDevices(c *gin.Context) {
	// Use lsblk to get real devices
	cmd := exec.Command("lsblk", "-J", "-o", "NAME,SIZE,TYPE,MODEL,RO,MOUNTPOINT,PATH")
	out, err := cmd.Output()
	if err != nil {
		utils.Error(c, 500, "failed to list devices")
		return
	}

	// Parse JSON output
	var result struct {
		BlockDevices []struct {
			Name       string `json:"name"`
			Size       string `json:"size"`
			Type       string `json:"type"`
			Model      string `json:"model"`
			RO         bool   `json:"ro"`
			MountPoint string `json:"mountpoint"`
			Path       string `json:"path"`
		} `json:"blockdevices"`
	}

	if err := parseJSON(out, &result); err != nil {
		utils.Error(c, 500, "failed to parse device list")
		return
	}

	var devices []models.StorageDevice
	for _, dev := range result.BlockDevices {
		if dev.Type != "disk" {
			continue
		}
		// Skip mounted devices
		if dev.MountPoint != "" {
			continue
		}
		
		devices = append(devices, models.StorageDevice{
			Name:      dev.Name,
			Path:      dev.Path,
			Type:      dev.Type,
			SizeBytes: parseSize(dev.Size),
			Model:     dev.Model,
			Vendor:    "",
		})
	}

	utils.Success(c, devices)
}

func parseSize(s string) int64 {
	// Simple parser for size strings like "100G", "1T"
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	
	var multiplier int64 = 1
	lastChar := s[len(s)-1]
	switch lastChar {
	case 'K':
		multiplier = 1024
		s = s[:len(s)-1]
	case 'M':
		multiplier = 1024 * 1024
		s = s[:len(s)-1]
	case 'G':
		multiplier = 1024 * 1024 * 1024
		s = s[:len(s)-1]
	case 'T':
		multiplier = 1024 * 1024 * 1024 * 1024
		s = s[:len(s)-1]
	}
	
	var size float64
	_, err := fmt.Sscanf(s, "%f", &size)
	if err != nil {
		return 0
	}
	
	return int64(size) * multiplier
}
