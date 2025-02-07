package services

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/sanda0/vps_pilot_agent/dto"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

func CollectSystemStat(msgChan chan dto.Msg, interval int, wg *sync.WaitGroup) {
	defer wg.Done()

	n := 0
	for {
		cpuUsage, err := cpu.Percent(time.Second, true)
		if err != nil {
			continue
		}
		memUsage, err := mem.VirtualMemory()
		if err != nil {
			continue
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
			continue
		}

		msg := dto.Msg{
			Msg:  "sys_stat",
			Data: systemStatJSON,
		}

		msgChan <- msg

		n++
		time.Sleep(time.Duration(interval) * time.Second)
	}

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
