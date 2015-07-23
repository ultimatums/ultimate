package outputs

import (
	"fmt"
	"time"

	"github.com/ultimatums/ultimate/config"
	"github.com/ultimatums/ultimate/model"
	"github.com/upmio/horus/log"

	"github.com/mattbaird/elastigo/lib"
)

type ElasticOutputType struct {
	Index string
}

var (
	elasticConn *elastigo.Conn
)

func (out *ElasticOutputType) Init(cfg *config.ElasticsearchConfig) {

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

}

// PublishMetric implements the Output interface.
func (out *ElasticOutputType) PublishMetric(ts time.Time, metric model.Metric) error {
	index := fmt.Sprintf("%s-%s", out.Index, ts.Format("2006.01.02"))
	fmt.Println(metric)
	//	return nil
	_, err := elasticConn.Index(index, metric["type"].(string), "", nil, metric)
	return err
}
