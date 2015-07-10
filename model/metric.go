package model

import (
	"packetbeat/common"
	"time"
)

type Metric map[string]interface{}

func newMetric(_metric string, _type string, val interface{}) Metric {
	metric := Metric{
		"metric":    _metric,
		"type":      _type,
		"value":     val,
		"timestamp": common.Time(time.Now()),
		"tags":      make(map[string]interface{}),
	}
	return metric
}

func GaugeMetric(_metric string, val interface{}) Metric {
	return newMetric(_metric, "gauge", val)
}

func CounterMetric(_metric string, val interface{}) Metric {
	return newMetric(_metric, "counter", val)
}

func (m *Metric) Collect(ch chan<- Metric) {

}
