package fetch

import (
	"time"

	"github.com/ultimatums/ultimate/config"
)

type Unit struct {
	fetchInterval time.Duration
	fetchStop     chan struct{}

	identity string
}

func NewUnit(unitCfg *config.UnitConfig) *Unit {
	u := &Unit{
		fetchStop: make(chan struct{}),
	}
	u.fetchInterval = time.Duration(unitCfg.FetchInterval)
	u.identity = unitCfg.Identity
	return u
}
