package monitors

import (
	"runtime"
	"sync"
)

type Monitor struct {
	mu            sync.Mutex
	Alloc         float64
	BuckHashSys   float64
	Frees         float64
	GCCPUFraction float64
	GCSys         float64
	HeapAlloc     float64
	HeapIdle      float64
	HeapInuse     float64
	HeapObjects   float64
	HeapReleased  float64
	HeapSys       float64
	LastGC        float64
	Lookups       float64
	MCacheInuse   float64
	MCacheSys     float64
	MSpanInuse    float64
	MSpanSys      float64
	Mallocs       float64
	NextGC        float64
	NumForcedGC   float64
	NumGC         float64
	OtherSys      float64
	PauseTotalNs  float64
	StackInuse    float64
	StackSys      float64
	Sys           float64
	TotalAlloc    float64
}

func (m *Monitor) Update(ms *runtime.MemStats) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Alloc = float64(ms.Alloc)
	m.BuckHashSys = float64(ms.BuckHashSys)
	m.Frees = float64(ms.Frees)
	m.GCCPUFraction = float64(ms.GCCPUFraction)
	m.GCSys = float64(ms.GCSys)
	m.HeapAlloc = float64(ms.HeapAlloc)
	m.HeapIdle = float64(ms.HeapIdle)
	m.HeapInuse = float64(ms.HeapInuse)
	m.HeapObjects = float64(ms.HeapObjects)
	m.HeapReleased = float64(ms.HeapReleased)
	m.HeapSys = float64(ms.HeapSys)
	m.LastGC = float64(ms.LastGC)
	m.Lookups = float64(ms.Lookups)
	m.MCacheInuse = float64(ms.MCacheInuse)
	m.MCacheSys = float64(ms.MCacheSys)
	m.MSpanInuse = float64(ms.MSpanInuse)
	m.MSpanSys = float64(ms.MSpanSys)
	m.Mallocs = float64(ms.Mallocs)
	m.NextGC = float64(ms.NextGC)
	m.NumForcedGC = float64(ms.NumForcedGC)
	m.NumGC = float64(ms.NumGC)
	m.OtherSys = float64(ms.OtherSys)
	m.PauseTotalNs = float64(ms.PauseTotalNs)
	m.StackInuse = float64(ms.StackInuse)
	m.StackSys = float64(ms.StackSys)
	m.Sys = float64(ms.Sys)
	m.TotalAlloc = float64(ms.TotalAlloc)
}

func (m *Monitor) Lock() {
	m.mu.Lock()
}

func (m *Monitor) Unlock() {
	m.mu.Unlock()
}
