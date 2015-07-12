package units

import (
	"errors"
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

var (
	UnitFactories = make(map[string]UnitFactory)
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
	metrics       []model.Metric
}

func (this *BaseUnit) String() string {
	return this.name
}

func NewUnit(unitCfg *config.UnitConfig, taskCfg *config.TaskConfig) (Unit, error) {
	/*
		var fac UnitFactory
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
	*/
	var facKey string
	if taskCfg.TaskName == UNIT_NAME_CONTAINER {
		facKey = taskCfg.TaskName
	} else {
		facKey = unitCfg.UnitName
	}
	fac, ok := UnitFactories[facKey]
	if !ok {
		return nil, errors.New("Unrecognized unit name.")
	}

	var u Unit
	u = fac.createUnit()
	u.SetInterval(unitCfg.FetchInterval)

	return Unit(u), nil
}
