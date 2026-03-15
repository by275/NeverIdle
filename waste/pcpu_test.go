package waste

import (
	"sync"
	"testing"
	"time"
)

func TestMachineControlConcurrentAccess(t *testing.T) {
	m := &machine{
		runtimePeriod:   time.Second,
		maxControlValue: 100000,
		idleTime:        time.Second,
	}

	var wg sync.WaitGroup
	for i := range 8 {
		wg.Add(1)
		go func(step float64) {
			defer wg.Done()
			for range 1000 {
				m.Control(step)
				busyTime, idleTime := m.currentTimings()
				if busyTime < 0 {
					t.Errorf("busyTime must be non-negative, got %d", busyTime)
				}
				if idleTime < 0 {
					t.Errorf("idleTime must be non-negative, got %s", idleTime)
				}
				if busyTime > m.runtimePeriod.Nanoseconds() {
					t.Errorf("busyTime must not exceed runtime period, got %d", busyTime)
				}
				if idleTime > m.runtimePeriod {
					t.Errorf("idleTime must not exceed runtime period, got %s", idleTime)
				}
			}
		}(float64(i+1) * 10)
	}

	wg.Wait()
}
