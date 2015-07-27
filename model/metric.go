package model

import (
	"fmt"
	"time"
)

type Metric map[string]interface{}

func (m Metric) String() string {
	return fmt.Sprintf("metric: %s,\tvalue: %v", m["metric"], m["value"].(float64))
}

func newMetric(_metric string, _type string) Metric {
	metric := Metric{
		"metric": _metric,
		"type":   _type,
		//		"value":  val,
		//		"timestamp": Time(time.Now()),
		"tags": make(map[string]interface{}),
	}
	return metric
}

func NewGauge(_metric string) Metric {
	return newMetric(_metric, "gauge")
}

func NewCounter(_metric string) Metric {
	return newMetric(_metric, "counter")
}

// Equal compares two metrics.
func (m Metric) Equal(o Metric) bool {
	if len(m) != len(o) {
		return false
	}
	for k, v := range m {
		ov, isOk := o[k]
		if !isOk {
			return false
		}
		if ov != v {
			return false
		}
	}
	return true
}

// Clone returns a copy of the Metric.
func (m Metric) Clone() Metric {
	clone := Metric{}
	for k, v := range m {
		clone[k] = v
	}
	return clone
}

func (m Metric) SetValue(value interface{}) Metric {
	m["value"] = value
	return m
}

func (m Metric) GetValue() interface{} {
	return m["value"]
}

func (m Metric) SetTimestamp(timestamp time.Time) Metric {
	m["timestamp"] = Time(timestamp)
	return m
}

func (m Metric) AddTag(tagName string, tagValue interface{}) Metric {
	m["tags"].(map[string]interface{})[tagName] = tagValue
	return m
}

func (m Metric) Collect(ch chan<- Metric) {
	ch <- m
}
