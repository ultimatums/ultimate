package outputs

import (
	"fmt"
	"time"

	"github.com/mattbaird/elastigo/lib"
	"github.com/ultimatums/ultimate/model"
)

type ElasticOutputType struct {
	Index string
}

var (
	elasticConn *elastigo.Conn
)

func (out *ElasticOutputType) Init( /*cfg *config.OutputConfig*/ ) {
	/*
		elasticConn = elastigo.NewConn()
		elasticConn.Domain = cfg.Host
		elasticConn.Port = fmt.Sprintf("%d", cfg.Port)
		elasticConn.Username = cfg.Username
		elasticConn.Password = cfg.Password

		if cfg.Protocol != "" {
			elasticConn.Protocol = cfg.Protocol
		}

		if cfg.Index != "" {
			out.Index = cfg.Index
		} else {
			out.Index = "horus"
		}

		log.Infof("[ElasticOutput] Using Elasticsearch %s://%s:%s", elasticConn.Protocol, elasticConn.Domain, elasticConn.Port)
		log.Infof("[ElasticOutput] Using index pattern [%s-]YYYY.MM.DD", out.Index)
	*/
}

// PublishMetric implements the Output interface.
func (out *ElasticOutputType) PublishMetric(ts time.Time, metrics model.Metric) error {
	index := fmt.Sprintf("%s-%s", out.Index, ts.Format("2006.01.02"))
	_, err := elasticConn.Index(index, metrics["input"].(string), "", nil, metrics)
	return err
}
