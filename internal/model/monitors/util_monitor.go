package monitors

import (
	"sync"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type UtilMonitor struct {
	mu              sync.Mutex
	TotalMemory     float64
	FreeMemory      float64
	CPUutilization1 float64
}

func (m *UtilMonitor) Update(ms *mem.VirtualMemoryStat) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.TotalMemory = float64(ms.Total)
	m.FreeMemory = float64(ms.Free)
	cpuNumber, err := cpu.Counts(false)
	if err != nil {
		return err
	}
	m.CPUutilization1 = float64(cpuNumber)

	return nil
}
