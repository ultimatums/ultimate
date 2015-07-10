package units

import (
	"time"

	"github.com/prometheus/log"
	"github.com/ultimatums/ultimate/model"
)

func init() {

}

type HostCpuUnitFactory struct{}

// createUnit implements the UnitFactory interface.
func (this *HostCpuUnitFactory) createUnit() Unit {
	u := new(HostCpuUnit)
	u.fetchStop = make(chan struct{})
	return u
}

type HostCpuUnit struct {
	BaseUnit
}

// SetInterval implements the Unit interface.
func (this *HostCpuUnit) SetInterval(interval model.Duration) {
	this.fetchInterval = time.Duration(interval)
}

// EqualTo implements the Unit interface.
func (this *HostCpuUnit) EqualTo(other Unit) bool {
	otherHostCpu := other.(*HostCpuUnit)
	return this.fetchInterval == otherHostCpu.fetchInterval && this.name == otherHostCpu.name
}

func (this *HostCpuUnit) Start() {

}
func (this *HostCpuUnit) Stop() {
	log.Debugf("Stopping fetch for unit %v...", this)
	close(this.fetchStop)
	log.Debugf("Fetch for unit %v stopped.", this)
}
