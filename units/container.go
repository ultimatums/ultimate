package units

import (
	"github.com/ultimatums/ultimate/model"
	"github.com/upmio/horus/log"
)

type ContainerUnitFactory struct{}

// createUnit implements the UnitFactory interface.
func (this *ContainerUnitFactory) createUnit() Unit {
	u := &ContainerUnit{
		BaseUnit{
			name:      UNIT_NAME_CONTAINER,
			fetchStop: make(chan struct{}),
		},
	}
	u.BaseUnit.unit = u
	return u
}

type ContainerUnit struct {
	BaseUnit
}

// EqualTo implements the Unit interface.
func (this *ContainerUnit) EqualTo(other Unit) bool {
	otherUnit, ok := other.(*ContainerUnit)
	if !ok {
		return false
	}
	return this.name == otherUnit.name
}

func (this *ContainerUnit) Fetch(ch chan<- model.Metric) {
	log.Info("in container fetch...")
}
