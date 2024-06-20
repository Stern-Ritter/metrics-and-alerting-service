package monitors

import (
	"sync"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// UtilMonitor holds system utilization metrics statistics.
type UtilMonitor struct {
	mu              sync.Mutex
	TotalMemory     float64
	FreeMemory      float64
	CPUutilization1 float64
}

// Update updates the UtilMonitor with the latest statistics.
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
