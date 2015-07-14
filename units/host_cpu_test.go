package units

import (
	"fmt"
	"testing"
	"time"
)

func Test_updateStats(t *testing.T) {
	unit := HostCpuUnit{
		lastCPUUsage: make([]float64, 7),
		newCPUUsage:  make([]float64, 7),
	}
	unit.updateStats()
	fmt.Println("last : ", unit.lastCPUUsage)
	fmt.Println("new : ", unit.newCPUUsage)
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			unit.updateStats()
			fmt.Println("last : ", unit.lastCPUUsage)
			fmt.Println("new : ", unit.newCPUUsage)
		}
	}
}
