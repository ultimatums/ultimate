package units

import (
	"sync"

	"github.com/ultimatums/ultimate/model"
)

var (
	BuiltInMetrics = MetricMap{
		Metrics: make(map[string]model.Metric),
	}
)

type MetricMap struct {
	sync.RWMutex
	Metrics map[string]model.Metric
}

func (this *MetricMap) Put(key string, elem model.Metric) model.Metric {
	this.Lock()
	defer this.Unlock()
	oldElem := this.Metrics[key]
	this.Metrics[key] = elem
	return oldElem
}

func (this *MetricMap) Get(key string) (model.Metric, bool) {
	this.RLock()
	defer this.RUnlock()
	elem, ok := this.Metrics[key]
	return elem, ok
}

func (this *MetricMap) Remove(key string) (model.Metric, bool) {
	this.Lock()
	defer this.Unlock()
	elem, ok := this.Metrics[key]
	if !ok {
		return nil, ok
	}
	delete(this.Metrics, key)
	return elem, true
}
