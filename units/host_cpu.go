package units

import (
	"github.com/ultimatums/ultimate/model"
	"github.com/upmio/horus/log"
)

func init() {
	BuiltInMetrics.Put("host.cpu.idle", model.GaugeMetric("host.cpu.idle"))
	BuiltInMetrics.Put("host.cpu.user", model.GaugeMetric("host.cpu.user"))
	BuiltInMetrics.Put("host.cpu.nice", model.GaugeMetric("host.cpu.nice"))
	BuiltInMetrics.Put("host.cpu.system", model.GaugeMetric("host.cpu.system"))
	BuiltInMetrics.Put("host.cpu.iowait", model.GaugeMetric("host.cpu.iowait"))
	BuiltInMetrics.Put("host.cpu.irq", model.GaugeMetric("host.cpu.irq"))
	BuiltInMetrics.Put("host.cpu.softirq", model.GaugeMetric("host.cpu.softirq"))
	BuiltInMetrics.Put("host.cpu.steal", model.GaugeMetric("host.cpu.steal"))
	BuiltInMetrics.Put("host.cpu.guest", model.GaugeMetric("host.cpu.guest"))
	BuiltInMetrics.Put("host.cpu.guest_nice", model.GaugeMetric("host.cpu.guest_nice"))

	UnitFactories[UNIT_NAME_HOST_CPU] = new(HostCpuUnitFactory)
}

type HostCpuUnitFactory struct{}

// createUnit implements the UnitFactory interface.
func (this *HostCpuUnitFactory) createUnit() Unit {
	u := &HostCpuUnit{
		BaseUnit{
			name:      UNIT_NAME_HOST_CPU,
			fetchStop: make(chan struct{}),
		},
	}
	return u
}

type HostCpuUnit struct {
	BaseUnit
}

// EqualTo implements the Unit interface.
func (this *HostCpuUnit) EqualTo(other Unit) bool {
	otherHostCpu, ok := other.(*HostCpuUnit)
	if !ok {
		return false
	}
	return this.name == otherHostCpu.name
}

func (this *HostCpuUnit) fetchFunc() func(ch chan<- model.Metric) {
	return func(ch chan<- model.Metric) {
		log.Info("in host cpu fetch...")
	}
}
