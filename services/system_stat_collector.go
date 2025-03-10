package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sanda0/vps_pilot_agent/dto"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

func StartCollectSystemStat(ctx context.Context, msgChan chan dto.Msg, interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping system stat collection")
			return
		case <-ticker.C:
			err := collectStat(msgChan)
			if err != nil {
				fmt.Println("Error collecting system stat:", err)
			}
		}
	}

}

func collectStat(msgChan chan dto.Msg) error {
	cpuUsage, err := cpu.Percent(time.Second, true)
	if err != nil {
		return err
	}
	memUsage, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	diskUsage := 0.0

	systemStat := dto.SystemStat{
		CPUUsage:  cpuUsage,
		MemUsage:  memUsage.UsedPercent,
		DiskUsage: diskUsage,
	}

	systemStatJSON, err := json.Marshal(systemStat)
	if err != nil {
		fmt.Println("Error marshalling system stat:", err)
		return err
	}

	msg := dto.Msg{
		Msg:  "sys_stat",
		Data: systemStatJSON,
	}

	msgChan <- msg
	return nil
}

func GetSystemInfo() (*dto.SystemInfo, error) {
	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}

	cpuInfo, err := cpu.Counts(true)
	if err != nil {
		return nil, err
	}

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	return &dto.SystemInfo{
		OS:              hostInfo.OS,
		Platform:        hostInfo.Platform,
		PlatformVersion: hostInfo.PlatformVersion,
		KernelVersion:   hostInfo.KernelVersion,
		CPUs:            cpuInfo,
		TotalMemory:     memInfo.Total,
	}, nil
}
