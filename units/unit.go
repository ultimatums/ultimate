package units

import (
	"time"

	"github.com/ultimatums/ultimate/config"
	"github.com/ultimatums/ultimate/model"
)

const (
	UNIT_NAME_HOST_CPU     = "cpu"
	UNIT_NAME_HOST_MEM     = "mem"
	UNIT_NAME_HOST_DISKIO  = "diskio"
	UNIT_NAME_HOST_NETWORK = "network"
	UNIT_NAME_CONTAINER    = "container"
)

type Unit interface {
	SetInterval(model.Duration)
	EqualTo(Unit) bool
	Start()
	Stop()
}

type UnitFactory interface {
	createUnit() Unit
}

type BaseUnit struct {
	name          string
	fetchInterval time.Duration
	fetchStop     chan struct{}
	//	metrics    []
}

func (this *BaseUnit) String() string {
	return this.name
}

func NewUnit(unitCfg *config.UnitConfig, taskCfg *config.TaskConfig) Unit {
	var fac UnitFactory
	var u Unit
	switch taskCfg.TaskName {
	case UNIT_NAME_CONTAINER:
		fac = new(ContainerUnitFactory)
	default:
		switch unitCfg.UnitName {
		case UNIT_NAME_HOST_CPU:
			fac = new(HostCpuUnitFactory)
		case UNIT_NAME_HOST_MEM:
			//		fac = new(HostMemUnitFactory)
		case UNIT_NAME_HOST_DISKIO:
		case UNIT_NAME_HOST_NETWORK:
		}
	}

	u = fac.createUnit()
	u.SetInterval(unitCfg.FetchInterval)

	return Unit(u)
}
