package dto

import "encoding/json"

type SystemStat struct {
	CPUUsage  []float64 `json:"cpu_usage"`
	MemUsage  float64   `json:"mem_usage"`
	DiskUsage float64   `json:"disk_usage"`
	NetSentPS uint64    `json:"net_sent_ps"`
	NetRecvPS uint64    `json:"net_recv_ps"`
}

type Config struct {
	Host     string
	Port     int
	Interval int
}

type Msg struct {
	Msg    string
	NodeId int32
	Token  string
	Data   []byte
}

// SystemInfo represents the system information of a machine.
// It includes details about the operating system, platform, kernel version,
// number of CPUs, total memory, and disk information.
type SystemInfo struct {
	OS              string `json:"os"`               // e.g. linux, windows
	Platform        string `json:"platform"`         // e.g. ubuntu, centos
	PlatformVersion string `json:"platform_version"` // e.g. 20.04, 8
	KernelVersion   string `json:"kernel_version"`   // e.g. 5.4.0-42-generic
	CPUs            int    `json:"cpus"`             // number of CPUs
	TotalMemory     uint64 `json:"total_memory"`     // total memory in bytes
}

type Disk struct {
	Device     string `json:"device"`     // e.g. /dev/sda1
	Mountpoint string `json:"mountpoint"` // e.g. /
	Fstype     string `json:"fstype"`     // e.g. ext4
	Opts       string `json:"opts"`       // e.g. rw
	Total      uint64 `json:"total"`      // total disk space in bytes
	Used       uint64 `json:"used"`       // used disk space in bytes
}

func (s *SystemInfo) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}
