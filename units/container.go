package units

import (
	"time"

	"github.com/ultimatums/ultimate/model"
)

type ContainerUnitFactory struct{}

// createUnit implements the UnitFactory interface.
func (this *ContainerUnitFactory) createUnit() Unit {
	u := new(ContainerUnit)
	u.fetchStop = make(chan struct{})
	return u
}

type ContainerUnit struct {
	BaseUnit
}

func (this *ContainerUnit) EqualTo(other Unit) bool {
	otherContainer := other.(*ContainerUnit)
	return this.fetchInterval == otherContainer.fetchInterval && this.name == otherContainer.name
}

func (this *ContainerUnit) SetInterval(interval model.Duration) {
	this.fetchInterval = time.Duration(interval)
}

func (this *ContainerUnit) Start() {

}
func (this *ContainerUnit) Stop() {

}
