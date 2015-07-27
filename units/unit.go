package units

import (
	"errors"
	"reflect"
	"sync"
	"time"

	"github.com/ultimatums/ultimate/config"
	"github.com/ultimatums/ultimate/model"
	"github.com/upmio/horus/log"
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
	Start(chan<- model.Metric)
	Stop()
}

type UnitFactory interface {
	createUnit() Unit
}

type BaseUnit struct {
	unit Unit
	sync.RWMutex
	name          string
	fetchInterval time.Duration
	fetchStop     chan struct{}
	metrics       []model.Metric
}

func (this *BaseUnit) String() string {
	return this.name
}

// SetInterval implements the Unit interface.
func (this *BaseUnit) SetInterval(interval model.Duration) {
	this.fetchInterval = time.Duration(interval)
}

// Start implements the Unit interface.
func (this *BaseUnit) Start(ch chan<- model.Metric) {
	this.RLock()
	fetchInterval := this.fetchInterval
	this.RUnlock()

	log.Debugf("Starting fetch for unit %v...", this)

	ticker := time.NewTicker(fetchInterval)
	defer ticker.Stop()

	c := reflect.ValueOf(this.unit)
	methodFetch := c.MethodByName("Fetch")
	methodFetch.Call([]reflect.Value{reflect.ValueOf(ch)})

	for {
		select {
		case <-this.fetchStop:
			return
		case <-ticker.C:
			methodFetch.Call([]reflect.Value{reflect.ValueOf(ch)})
		}
	}
}

// Stop implements the Unit interface.
func (this *BaseUnit) Stop() {
	log.Debugf("Stopping fetch for unit %v...", this)
	close(this.fetchStop)
	log.Debugf("Fetch for unit %v stopped.", this)
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
	log.Info("factory key = ", facKey)
	fac, ok := UnitFactories[facKey]
	if !ok {
		return nil, errors.New("Unrecognized unit name: " + facKey)
	}

	var u Unit
	u = fac.createUnit()
	u.SetInterval(unitCfg.FetchInterval)

	return Unit(u), nil
}
