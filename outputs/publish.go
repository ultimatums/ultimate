package outputs

import (
	"errors"
	"fmt"
	"time"

	"github.com/ultimatums/ultimate/model"
	"github.com/upmio/horus/log"
)

var (
	Publisher = NewPublisherType()
)

type PublisherType struct {
	Outputs []Output
	Queue   chan model.Metric
	//	ElasticOutput ElasticOutputType
}

func NewPublisherType() *PublisherType {
	return &PublisherType{}
}

func (this *PublisherType) publishFromQueue() {
	for metric := range this.Queue {
		err := this.publishMetric(metric)
		if err != nil {
			log.Error(err)
		}
	}
}

func (this *PublisherType) publishMetric(metric model.Metric) error {
	ts, ok := metric["timestamp"].(model.Time)
	if !ok {
		return errors.New("Missing 'timestamp' field from metric.")
	}

	_, ok = metric["input"].(string)
	if !ok {
		return errors.New("Missing 'type' field from metric.")
	}

	//	metric["host"] = os.Hostname()
	has_error := false
	for _, output := range this.Outputs {
		//TODO Try to concurrency publish.
		err := output.PublishMetric(time.Time(ts), metric)
		if err != nil {
			errors.New(fmt.Sprintf("Fail to publish metric on output: %v, error: %s", output, err))
			has_error = true
		}
	}

	if has_error {
		return errors.New("Fail to publish metrics")
	}

	return nil
}

func (this *PublisherType) Init() {
	//	output, isExists := config.AppConfig().Output["elasticsearch"]
	//	log.Debug("is Exists = ", isExists)
	elasticOutput := &ElasticOutputType{}
	//	if isExists {
	//		elasticOutput.Init(output)
	//	}

	this.Outputs = append(this.Outputs, Output(elasticOutput))

	// TODO other outputs initalize...

	this.Queue = make(chan model.Metric, 200)
	go this.publishFromQueue()
}

// Output identifier
type OutputPlugin uint16

type Output interface {
	PublishMetric(ts time.Time, metric model.Metric) error
}
