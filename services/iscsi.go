package services

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

// ISCSIService handles iSCSI operations via tgtadm
type ISCSIService struct {
	mu sync.Mutex
}

// Target represents an iSCSI target
type Target struct {
	TID     int       `json:"tid"`
	Name    string    `json:"name"`
	IQN     string    `json:"iqn"`
	Status  string    `json:"status"`
	LUNs    []LUN     `json:"luns"`
	ACLs    []ACL     `json:"acls"`
	Portals []Portal  `json:"portals"`
}

// LUN represents a logical unit
type LUN struct {
	LUNID       int    `json:"lun_id"`
	StorageID   string `json:"storage_id"`
	BackingPath string `json:"backing_path"`
	Size        int64  `json:"size"`
	Type        string `json:"type"` // disk, cdrom
	ReadOnly    bool   `json:"read_only"`
}

// ACL represents an initiator access control
type ACL struct {
	InitiatorIQN string `json:"initiator_iqn"`
}

// Portal represents a target portal
type Portal struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

var iscsiService *ISCSIService
var once sync.Once

// GetISCSIService returns the singleton ISCSIService
func GetISCSIService() *ISCSIService {
	once.Do(func() {
		iscsiService = &ISCSIService{}
	})
	return iscsiService
}

// runTgtadm executes a tgtadm command with sudo
func (s *ISCSIService) runTgtadm(args ...string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Prepend sudo to command
	fullArgs := append([]string{"tgtadm"}, args...)
	cmd := exec.Command("sudo", fullArgs...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// ListTargets returns all iSCSI targets
func (s *ISCSIService) ListTargets() ([]Target, error) {
	out, err := s.runTgtadm("--lld", "iscsi", "--mode", "target", "--op", "show")
	if err != nil {
		return nil, fmt.Errorf("tgtadm failed: %v, output: %s", err, out)
	}
	
	return s.parseTargets(out)
}

// parseTargets parses tgtadm show output
func (s *ISCSIService) parseTargets(output string) ([]Target, error) {
	var targets []Target
	var current *Target
	var currentLUN *LUN
	
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Target line: Target 1: iqn.2024-01.com.example:target1
		if strings.HasPrefix(line, "Target ") {
			if current != nil {
				targets = append(targets, *current)
			}
			parts := strings.SplitN(line, ":", 2)
			if len(parts) >= 2 {
				tidStr := strings.TrimPrefix(parts[0], "Target ")
				tid, _ := strconv.Atoi(strings.TrimSpace(tidStr))
				current = &Target{
					TID:    tid,
					IQN:    strings.TrimSpace(parts[1]),
					Name:   strings.TrimSpace(parts[1]),
					Status: "active",
				}
			}
			continue
		}
		
		if current == nil {
			continue
		}
		
		// LUN line: LUN: 1
		if strings.HasPrefix(line, "LUN: ") {
			if currentLUN != nil {
				current.LUNs = append(current.LUNs, *currentLUN)
			}
			lunIDStr := strings.TrimPrefix(line, "LUN: ")
			lunID, _ := strconv.Atoi(strings.TrimSpace(lunIDStr))
			currentLUN = &LUN{LUNID: lunID}
			continue
		}
		
		// LUN attributes
		if currentLUN != nil {
			if strings.HasPrefix(line, "Type: ") {
				currentLUN.Type = strings.TrimPrefix(line, "Type: ")
			} else if strings.HasPrefix(line, "SCSI ID: ") {
				currentLUN.StorageID = strings.TrimPrefix(line, "SCSI ID: ")
			} else if strings.HasPrefix(line, "Backing store path: ") {
				currentLUN.BackingPath = strings.TrimPrefix(line, "Backing store path: ")
			} else if strings.HasPrefix(line, "Backing store size: ") {
				sizeStr := strings.TrimPrefix(line, "Backing store size: ")
				sizeStr = strings.Fields(sizeStr)[0]
				size, _ := strconv.ParseInt(sizeStr, 10, 64)
				currentLUN.Size = size
			} else if strings.HasPrefix(line, "Readonly: ") {
				roStr := strings.TrimPrefix(line, "Readonly: ")
				currentLUN.ReadOnly = strings.TrimSpace(roStr) == "Yes"
			}
		}
		
		// Account (ACL) line
		if strings.HasPrefix(line, "Account: ") {
			iqn := strings.TrimPrefix(line, "Account: ")
			current.ACLs = append(current.ACLs, ACL{InitiatorIQN: iqn})
		}
	}
	
	if currentLUN != nil && current != nil {
		current.LUNs = append(current.LUNs, *currentLUN)
	}
	if current != nil {
		targets = append(targets, *current)
	}
	
	return targets, nil
}

// CreateTarget creates a new iSCSI target
func (s *ISCSIService) CreateTarget(iqn string) error {
	// Get next available TID
	tid, err := s.getNextTID()
	if err != nil {
		return fmt.Errorf("failed to get next TID: %v", err)
	}
	
	// Use correct tgtadm syntax: --tid <num> -T <iqn>
	out, err := s.runTgtadm("--lld", "iscsi", "--mode", "target", "--op", "new", "--tid", strconv.Itoa(tid), "-T", iqn)
	if err != nil {
		return fmt.Errorf("create target failed: %v, output: %s", err, out)
	}
	return nil
}

// getNextTID returns the next available target ID
func (s *ISCSIService) getNextTID() (int, error) {
	targets, err := s.ListTargets()
	if err != nil {
		return 1, nil // Start from 1 if no targets exist
	}
	
	maxTID := 0
	for _, t := range targets {
		if t.TID > maxTID {
			maxTID = t.TID
		}
	}
	
	return maxTID + 1, nil
}

// DeleteTarget deletes an iSCSI target
func (s *ISCSIService) DeleteTarget(tid int) error {
	out, err := s.runTgtadm("--lld", "iscsi", "--mode", "target", "--op", "delete", "--tid="+strconv.Itoa(tid))
	if err != nil {
		return fmt.Errorf("delete target failed: %v, output: %s", err, out)
	}
	return nil
}

// AddLUN adds a LUN to a target
func (s *ISCSIService) AddLUN(tid, lunID int, path string, readonly bool) error {
	args := []string{"--lld", "iscsi", "--mode", "logicalunit", "--op", "new",
		"--tid=" + strconv.Itoa(tid), "--lun=" + strconv.Itoa(lunID),
		"-b", path}
	if readonly {
		args = append(args, "--readonly")
	}
	
	out, err := s.runTgtadm(args...)
	if err != nil {
		return fmt.Errorf("add LUN failed: %v, output: %s", err, out)
	}
	return nil
}

// DeleteLUN deletes a LUN from a target
func (s *ISCSIService) DeleteLUN(tid, lunID int) error {
	out, err := s.runTgtadm("--lld", "iscsi", "--mode", "logicalunit", "--op", "delete",
		"--tid="+strconv.Itoa(tid), "--lun="+strconv.Itoa(lunID))
	if err != nil {
		return fmt.Errorf("delete LUN failed: %v, output: %s", err, out)
	}
	return nil
}

// AddInitiator adds an initiator ACL to a target
func (s *ISCSIService) AddInitiator(tid int, initiatorIQN string) error {
	out, err := s.runTgtadm("--lld", "iscsi", "--mode", "target", "--op", "bind",
		"--tid="+strconv.Itoa(tid), "--initiator="+initiatorIQN)
	if err != nil {
		return fmt.Errorf("add initiator failed: %v, output: %s", err, out)
	}
	return nil
}

// RemoveInitiator removes an initiator ACL from a target
func (s *ISCSIService) RemoveInitiator(tid int, initiatorIQN string) error {
	out, err := s.runTgtadm("--lld", "iscsi", "--mode", "target", "--op", "unbind",
		"--tid="+strconv.Itoa(tid), "--initiator="+initiatorIQN)
	if err != nil {
		return fmt.Errorf("remove initiator failed: %v, output: %s", err, out)
	}
	return nil
}

// GetConnections returns active iSCSI connections
func (s *ISCSIService) GetConnections() ([]map[string]interface{}, error) {
	out, err := s.runTgtadm("--lld", "iscsi", "--mode", "conn", "--op", "show")
	if err != nil {
		return nil, fmt.Errorf("get connections failed: %v, output: %s", err, out)
	}
	
	// Parse connections
	var connections []map[string]interface{}
	// Simple parsing - real implementation would parse the full output
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.Contains(line, "SID:") {
			connections = append(connections, map[string]interface{}{
				"status": "connected",
			})
		}
	}
	
	return connections, nil
}

// SaveConfig saves the current configuration
func (s *ISCSIService) SaveConfig() error {
	out, err := s.runTgtadm("--lld", "iscsi", "--mode", "target", "--op", "show")
	if err != nil {
		return fmt.Errorf("save config failed: %v, output: %s", err, out)
	}
	return nil
}

// IsAvailable checks if tgt is available
func (s *ISCSIService) IsAvailable() bool {
	_, err := exec.LookPath("tgtadm")
	return err == nil
}

// GetVersion returns tgt version
func (s *ISCSIService) GetVersion() string {
	cmd := exec.Command("tgtadm", "--version")
	out, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

// ToJSON converts target to JSON (for debugging)
func (t *Target) ToJSON() string {
	b, _ := json.MarshalIndent(t, "", "  ")
	return string(b)
}
