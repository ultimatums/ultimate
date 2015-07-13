package units

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

// EqualTo implements the Unit interface.
func (this *ContainerUnit) EqualTo(other Unit) bool {
	otherUnit, ok := other.(*ContainerUnit)
	if !ok {
		return false
	}
	return this.name == otherUnit.name
}
