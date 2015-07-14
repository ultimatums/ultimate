package units

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ultimatums/ultimate/model"
	"github.com/upmio/horus/log"
)

const (
	procStat = "/proc/stat"
)

func init() {
	UnitFactories[UNIT_NAME_HOST_CPU] = new(HostCpuUnitFactory)
}

type HostCpuUnitFactory struct{}

// createUnit implements the UnitFactory interface.
func (this *HostCpuUnitFactory) createUnit() Unit {
	u := &HostCpuUnit{
		lastCPUUsage: make([]float64, 7),
		newCPUUsage:  make([]float64, 7),

		idle:    model.NewGauge("host.cpu.idle"),
		user:    model.NewGauge("host.cpu.user"),
		nice:    model.NewGauge("host.cpu.nice"),
		system:  model.NewGauge("host.cpu.system"),
		iowait:  model.NewGauge("host.cpu.iowait"),
		irq:     model.NewGauge("host.cpu.irq"),
		softirq: model.NewGauge("host.cpu.softirq"),

		intr:          model.NewCounter("host.stat.intr"),
		ctxt:          model.NewCounter("host.stat.ctxt"),
		btime:         model.NewGauge("host.stat.btime"),
		processes:     model.NewCounter("host.stat.processes"),
		procs_running: model.NewGauge("host.stat.procs_running"),
		procs_blocked: model.NewGauge("host.stat.procs_blocked"),
	}
	u.name = UNIT_NAME_HOST_CPU
	u.fetchStop = make(chan struct{})

	go func() {
		u.updateStats()
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				u.updateStats()
			}
		}
	}()

	u.BaseUnit.unit = u
	return u
}

type HostCpuUnit struct {
	BaseUnit
	sync.RWMutex
	lastCPUUsage []float64
	newCPUUsage  []float64

	idle    model.Metric // (4) Time spent in the idle task. This value should be USER_HZ times the second entry in the /proc/uptime pseudo-file.
	user    model.Metric // (1) Time spent in user mode.
	nice    model.Metric // (2) Time spent in user mode with low priority(nice).
	system  model.Metric // (3) Time spent in system mode.
	iowait  model.Metric // (since Linux 2.5.41) (5) Time waiting for I/O to complete.
	irq     model.Metric // (since Linux 2.6.0-test4) (6) Time servicing interrupts.
	softirq model.Metric // (since Linux 2.6.0-test4) (7) Time servicing softirqs.
	// steal      model.Metric // (since Linux 2.6.11) (8) Stolen time, which is the time spent in other operating systems when running in a virtualized environment
	// guest      model.Metric // (since Linux 2.6.24) (9) Time spent running a virtual CPU for guest operating systems under the control of the Linux kernel.
	// guest_nice model.Metric // (since Linux 2.6.33) (10) Time spent running a niced guest (virtual CPU for guest operating systems under the control of the Linux kernel).

	intr          model.Metric
	ctxt          model.Metric
	btime         model.Metric
	processes     model.Metric
	procs_running model.Metric
	procs_blocked model.Metric
}

// EqualTo implements the Unit interface.
func (this *HostCpuUnit) EqualTo(other Unit) bool {
	otherHostCpu, ok := other.(*HostCpuUnit)
	if !ok {
		return false
	}
	return this.name == otherHostCpu.name
}

func (this *HostCpuUnit) Fetch(ch chan<- model.Metric) {
	log.Info("in host cpu fetch...")

}

func (this *HostCpuUnit) updateStats() error {
	file, err := os.Open(procStat)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) == 0 {
			continue
		}
		switch parts[0] {
		case "cpu":
			this.lastCPUUsage = this.newCPUUsage
			cpuFields := []string{"user", "nice", "system", "idle", "iowait", "irq", "softirq"}
			for i, _ := range cpuFields {
				value, err := strconv.ParseFloat(parts[i+1], 64)
				if err != nil {
					return err
				}
				this.newCPUUsage[i] = value
			}
		case "intr":
		}

	}

	return nil
}

type CPUUsage struct {
	idle    float64 // (4) Time spent in the idle task. This value should be USER_HZ times the second entry in the /proc/uptime pseudo-file.
	user    float64 // (1) Time spent in user mode.
	nice    float64 // (2) Time spent in user mode with low priority(nice).
	system  float64 // (3) Time spent in system mode.
	iowait  float64 // (since Linux 2.5.41) (5) Time waiting for I/O to complete.
	irq     float64 // (since Linux 2.6.0-test4) (6) Time servicing interrupts.
	softirq float64 // (since Linux 2.6.0-test4) (7) Time servicing softirqs.
	total   float64
}
