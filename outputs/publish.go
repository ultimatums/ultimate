package outputs

import (
	"time"

	"github.com/ultimatums/ultimate/config"
	"github.com/ultimatums/ultimate/model"
	"github.com/upmio/horus/log"
)

const (
	output_elastic = "elasticsearch"
	output_hadoop  = "hadoop"
)

var (
	Publisher = NewPublisherType()
)

type PublisherType struct {
	Outputs       map[string]Output
	Queue         chan model.Metric
	ElasticOutput *ElasticOutputType
	publistStop   chan struct{}
}

func NewPublisherType() *PublisherType {
	publisher := PublisherType{
		publistStop: make(chan struct{}),
		Outputs:     make(map[string]Output),
		Queue:       make(chan model.Metric, 200),
	}
	return &publisher
}

func (this *PublisherType) publishFromQueue() {
	for {
		select {
		case metric := <-this.Queue:
			this.publishMetric(metric)
		case <-this.publistStop:
			log.Debug("Done: This is publishFromQueue goroutine.")
			return
		}
	}
	/*
		for metric := range this.Queue {
			this.publishMetric(metric)
		}
	*/
}

func (this *PublisherType) publishMetric(metric model.Metric) {
	ts, ok := metric["timestamp"].(model.Time)
	if !ok {
		log.Error("Missing 'timestamp' field in metric.")
	}

	_, ok = metric["type"].(string)
	if !ok {
		log.Error("Missing 'type' field in metric.")
	}

	for name, output := range this.Outputs {
		// TODO Try to concurrency publish.
		err := output.PublishMetric(time.Time(ts), metric)
		if err != nil {
			log.Errorf("Fail to publish metric on output: %s, error: %s", name, err)
		}
	}
}

func (this *PublisherType) ApplyConfig(cfg *config.Config) {
	this.StopPublish()

	if cfg.OutputConfig.Elasticsearch != nil {
		if this.ElasticOutput == nil {
			this.ElasticOutput = &ElasticOutputType{}
		}
		this.ElasticOutput.Init(cfg.OutputConfig.Elasticsearch)
		this.Outputs[output_elastic] = this.ElasticOutput
	}

	this.publistStop = make(chan struct{})
	go this.publishFromQueue()
}

func (this *PublisherType) StopPublish() {
	log.Info("Stoping publish...")
	close(this.publistStop) // stop publishFromQueue gorutine
}

type Hostinfo struct {
	Hostname string
	Ips      []string
	DockerID string
}

// Output identifier
type OutputPlugin uint16

type Output interface {
	PublishMetric(ts time.Time, metric model.Metric) error
}
