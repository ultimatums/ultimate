package units

import "github.com/ultimatums/ultimate/model"

const (
	procMeminfo = "/proc/meminfo"
)

func init() {
	UnitFactories[UNIT_NAME_HOST_MEM] = new(HostMemoryUnitFactory)
}

type HostMemoryUnitFactory struct{}

func (this *HostMemoryUnitFactory) createUnit() Unit {
	u := &HostMemoryUnit{
		metrics: []model.Metric{
			model.NewGauge("host.mem.total").AddTag("hostname", hostname),
			model.NewGauge("host.mem.free").AddTag("hostname", hostname),
			model.NewGauge("host.mem.available").AddTag("hostname", hostname),
			model.NewGauge("host.mem.app").AddTag("hostname", hostname),
			model.NewGauge("host.mem.os").AddTag("hostname", hostname),
			model.NewGauge("host.mem.buffers").AddTag("hostname", hostname),
			model.NewGauge("host.mem.cached").AddTag("hostname", hostname),
			model.NewGauge("host.mem.used").AddTag("hostname", hostname),
			model.NewGauge("host.mem.swaptotal").AddTag("hostname", hostname),
			model.NewGauge("host.mem.swapfree").AddTag("hostname", hostname),
			model.NewGauge("host.mem.swapused").AddTag("hostname", hostname),
			model.NewGauge("host.mem.free.percent").AddTag("hostname", hostname),
			model.NewGauge("host.mem.used.percent").AddTag("hostname", hostname),
			model.NewGauge("host.mem.swapfree.percent").AddTag("hostname", hostname),
			model.NewGauge("host.mem.swapused.percent").AddTag("hostname", hostname),
		},
	}
	u.name = UNIT_NAME_HOST_MEM
	u.fetchStop = make(chan struct{})
	u.BaseUnit.unit = u
	return u
}

type HostMemoryUnit struct {
	BaseUnit
	metrics []model.Metric
}

// EqualTo implements the Unit interface.
func (this *HostMemoryUnit) EqualTo(other Unit) bool {
	otherHostCpu, ok := other.(*HostMemoryUnit)
	if !ok {
		return false
	}
	return this.name == otherHostCpu.name
}
