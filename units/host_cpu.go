package units

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ultimatums/ultimate/model"
	"github.com/upmio/horus/log"
)

const (
	procStat = "/proc/stat"
)

var (
	//	cpuFields = []string{"user", "nice", "system", "idle", "iowait", "irq", "softirq", "total"}
	hostname, _ = os.Hostname()
)

func init() {
	UnitFactories[UNIT_NAME_HOST_CPU] = new(HostCpuUnitFactory)
}

type HostCpuUnitFactory struct{}

// createUnit implements the UnitFactory interface.
func (this *HostCpuUnitFactory) createUnit() Unit {
	u := &HostCpuUnit{
		lastCPUUsage: make([]float64, 8),
		newCPUUsage:  make([]float64, 8),

		metrics: []model.Metric{
			model.NewGauge("host.cpu.user").AppendTag("hostname", hostname),
			model.NewGauge("host.cpu.nice").AppendTag("hostname", hostname),
			model.NewGauge("host.cpu.system").AppendTag("hostname", hostname),
			model.NewGauge("host.cpu.idle").AppendTag("hostname", hostname),
			model.NewGauge("host.cpu.iowait").AppendTag("hostname", hostname),
			model.NewGauge("host.cpu.irq").AppendTag("hostname", hostname),
			model.NewGauge("host.cpu.softirq").AppendTag("hostname", hostname),

			//			model.NewCounter("host.stat.intr").AppendTag("hostname", hostname),
			//			model.NewCounter("host.stat.ctxt").AppendTag("hostname", hostname),
			//			model.NewGauge("host.stat.btime").AppendTag("hostname", hostname),
			//			model.NewCounter("host.stat.processes").AppendTag("hostname", hostname),
			//			model.NewGauge("host.stat.procs_running").AppendTag("hostname", hostname),
			//			model.NewGauge("host.stat.procs_blocked").AppendTag("hostname", hostname),
		},
	}
	u.name = UNIT_NAME_HOST_CPU
	u.fetchStop = make(chan struct{})

	go func() {
		err := u.updateStats()
		if err != nil {
			log.Errorf("update cpu stats error: ", err)
		}
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				err = u.updateStats()
				if err != nil {
					log.Errorf("update cpu stats error: ", err)
				}
			}
		}
	}()

	u.BaseUnit.unit = u
	return u
}

type HostCpuUnit struct {
	BaseUnit
	lastCPUUsage []float64
	newCPUUsage  []float64
	/*
		user    model.Metric // (1) Time spent in user mode.
		nice    model.Metric // (2) Time spent in user mode with low priority(nice).
		system  model.Metric // (3) Time spent in system mode.
		idle    model.Metric // (4) Time spent in the idle task. This value should be USER_HZ times the second entry in the /proc/uptime pseudo-file.
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
	*/
	metrics []model.Metric
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
	log.Info("in host cpu unit fetch...")
	if this.lastCPUUsage[7] == 0 {
		return
	}
	delta_total := this.newCPUUsage[7] - this.lastCPUUsage[7]
	for i := 0; i < 7; i++ {
		this.metrics[i].
			SetValue((this.newCPUUsage[i] - this.lastCPUUsage[i]) * 100.0 / delta_total).
			SetTimestamp(time.Now()).Collect(ch)
	}
}

func (this *HostCpuUnit) updateStats() error {
	file, err := os.Open(procStat)
	if err != nil {
		return err
	}
	defer file.Close()

	//	log.Info("last--->", this.lastCPUUsage)
	//	log.Info("new --->", this.newCPUUsage)
	//	log.Infof("user=%v,nice=%v,system=%v,idle=%v,iowait=%v,irq=%v,softirq=%v", this.user, this.nice, this.system, this.idle, this.iowait, this.irq, this.softirq)

	this.Lock()
	defer this.Unlock()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) == 0 {
			continue
		}
		switch parts[0] {
		case "cpu":
			copy(this.lastCPUUsage, this.newCPUUsage)
			var cpu_total float64
			for i := 1; i < len(parts); i++ {
				value, err := strconv.ParseFloat(parts[i], 64)
				if err != nil {
					return err
				}
				if (i - 1) < 7 {
					this.newCPUUsage[i-1] = value
				}
				cpu_total += value
			}
			this.newCPUUsage[7] = cpu_total
			/*
				case "intr":
					value, err := strconv.ParseFloat(parts[1], 64)
					if err != nil {
						return err
					}
					this.metrics[8].SetValue(value)
				case "ctxt":
					value, err := strconv.ParseFloat(parts[1], 64)
					if err != nil {
						return err
					}
					this.metrics[9].SetValue(value)
				case "btime":
					value, err := strconv.ParseFloat(parts[1], 64)
					if err != nil {
						return err
					}
					this.metrics[10].SetValue(value)
				case "processes":
					value, err := strconv.ParseFloat(parts[1], 64)
					if err != nil {
						return err
					}
					this.metrics[11].SetValue(value)
				case "procs_running":
					value, err := strconv.ParseFloat(parts[1], 64)
					if err != nil {
						return err
					}
					this.metrics[12].SetValue(value)
				case "procs_blocked":
					value, err := strconv.ParseFloat(parts[1], 64)
					if err != nil {
						return err
					}
					this.metrics[13].SetValue(value)
			*/
		}
	}
	return nil
}
