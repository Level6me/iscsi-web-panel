package models

import "time"

// iSCSI Target
type Target struct {
	Name       string    `json:"name"`
	IQN        string    `json:"iqn"`
	Status     string    `json:"status"`
	LUNs       []LUN     `json:"luns"`
	ACLs       []ACL     `json:"acls"`
	Portals    []Portal  `json:"portals"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// LUN
type LUN struct {
	ID         int       `json:"id"`
	Index      int       `json:"index"`
	TargetName string    `json:"target_name"`
	StorageObj string    `json:"storage_obj"`
	BackingPath string   `json:"backing_path"`
	SizeBytes  int64     `json:"size_bytes"`
	Type       string    `json:"type"` // blockio, fileio, pscsi, iblock
	ReadOnly   bool      `json:"read_only"`
	CreatedAt  time.Time `json:"created_at"`
}

// ACL (Initiator access)
type ACL struct {
	InitiatorIQN string     `json:"initiator_iqn"`
	CHAP         *CHAPAuth  `json:"chap,omitempty"`
	LUNMapping   []LUNMap   `json:"lun_mapping"`
}

type LUNMap struct {
	LUNIndex   int    `json:"lun_index"`
	MappedLUN  int    `json:"mapped_lun"`
	TargetName string `json:"target_name"`
}

// CHAP Authentication
type CHAPAuth struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	MutualUser string `json:"mutual_username,omitempty"`
	MutualPass string `json:"mutual_password,omitempty"`
}

// Portal
type Portal struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

// Network Interface
type NetworkInterface struct {
	Name     string   `json:"name"`
	IP       string   `json:"ip"`
	Netmask  string   `json:"netmask"`
	Status   string   `json:"status"`
	MAC      string   `json:"mac"`
	IsPortal bool     `json:"is_portal"`
}

// Storage Pool
type StoragePool struct {
	Name       string    `json:"name"`
	Type       string    `json:"type"` // fileio, blockio, iblock
	Path       string    `json:"path"`
	TotalBytes int64     `json:"total_bytes"`
	UsedBytes  int64     `json:"used_bytes"`
	FreeBytes  int64     `json:"free_bytes"`
	CreatedAt  time.Time `json:"created_at"`
}

// Storage Device
type StorageDevice struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	Type       string `json:"type"` // disk, partition
	SizeBytes  int64  `json:"size_bytes"`
	Model      string `json:"model"`
	Vendor     string `json:"vendor"`
}

// Monitor
type Metrics struct {
	TxBytes     int64   `json:"tx_bytes"`
	RxBytes     int64   `json:"rx_bytes"`
	IOPS        float64 `json:"iops"`
	LatencyMs   float64 `json:"latency_ms"`
	Connections int     `json:"connections"`
	Targets     int     `json:"targets"`
	LUNs        int     `json:"luns"`
	Timestamp   time.Time `json:"timestamp"`
}

type PerformanceSample struct {
	Timestamp time.Time `json:"timestamp"`
	IOPS      float64   `json:"iops"`
	ThroughputMBps float64 `json:"throughput_mbps"`
	LatencyMs float64   `json:"latency_ms"`
}

type ConnectionInfo struct {
	InitiatorIQN string    `json:"initiator_iqn"`
	TargetIQN    string    `json:"target_iqn"`
	SessionID    string    `json:"session_id"`
	State        string    `json:"state"`
	ConnectedAt  time.Time `json:"connected_at"`
	TxBytes      int64     `json:"tx_bytes"`
	RxBytes      int64     `json:"rx_bytes"`
}

// Alert
type Alert struct {
	ID        int       `json:"id"`
	Level     string    `json:"level"` // info, warning, critical
	Message   string    `json:"message"`
	Source    string    `json:"source"`
	Acked     bool      `json:"acked"`
	CreatedAt time.Time `json:"created_at"`
}

type AlertRule struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Metric    string `json:"metric"`
	Operator  string `json:"operator"` // gt, lt, eq
	Threshold float64 `json:"threshold"`
	Enabled   bool   `json:"enabled"`
}

// Log Entry
type LogEntry struct {
	ID        int       `json:"id"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Source    string    `json:"source"`
	Timestamp time.Time `json:"timestamp"`
}

// User
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"` // admin, operator, viewer
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// Settings
type Settings struct {
	SystemName       string `json:"system_name"`
	DefaultCHAPUser  string `json:"default_chap_user"`
	AutoSnapshot     bool   `json:"auto_snapshot"`
	SnapshotInterval int    `json:"snapshot_interval_minutes"`
	LogRetentionDays int    `json:"log_retention_days"`
}

// Snapshot
type Snapshot struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	LUNID       int       `json:"lun_id"`
	TargetName  string    `json:"target_name"`
	SizeBytes   int64     `json:"size_bytes"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// Dashboard
type DashboardOverview struct {
	TotalTargets     int     `json:"total_targets"`
	ActiveTargets    int     `json:"active_targets"`
	TotalLUNs        int     `json:"total_luns"`
	TotalConnections int     `json:"total_connections"`
	TotalStorageGB   float64 `json:"total_storage_gb"`
	UsedStorageGB    float64 `json:"used_storage_gb"`
	AlertCount       int     `json:"alert_count"`
	Uptime           string  `json:"uptime"`
}

type DashboardStats struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemUsage    float64 `json:"mem_usage"`
	DiskUsage   float64 `json:"disk_usage"`
	NetworkIn   int64   `json:"network_in"`
	NetworkOut  int64   `json:"network_out"`
	KernelVer   string  `json:"kernel_version"`
	TargetCliVer string `json:"targetcli_version"`
}

// API Response
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
	User      User   `json:"user"`
}

// API Doc
type APIDocEndpoint struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	Description string `json:"description"`
	Auth        bool   `json:"auth"`
}
