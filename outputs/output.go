package outputs

import (
	"time"

	"github.com/ultimatums/ultimate/model"
)

// Output identifier
type OutputPlugin uint16

type Output interface {
	PublishMetric(ts time.Time, metric model.Metric) error
}
