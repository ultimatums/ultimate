package model

type Metric map[string]interface{}

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

func GaugeMetric(_metric string) Metric {
	return newMetric(_metric, "gauge")
}

func CounterMetric(_metric string) Metric {
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

func (m Metric) Collect(ch chan<- Metric) {
	ch <- m
}
