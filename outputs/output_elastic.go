package outputs

import (
	"fmt"
	"sync"
	"time"

	"github.com/ultimatums/ultimate/config"
	"github.com/ultimatums/ultimate/model"
	"github.com/upmio/horus/log"

	"github.com/mattbaird/elastigo/lib"
)

const (
	DefaultIndex = "horus"
)

type ElasticOutputType struct {
	sync.Mutex
	Index       string
	elasticConn *elastigo.Conn
}

func (out *ElasticOutputType) Init(cfg *config.ElasticsearchConfig) {
	out.Lock()
	defer out.Unlock()

	if out.elasticConn != nil {
		out.elasticConn.Close()
	}
	out.elasticConn = elastigo.NewConn()

	out.elasticConn.Domain = cfg.Host
	out.elasticConn.Port = fmt.Sprintf("%d", cfg.Port)
	out.elasticConn.Username = cfg.Username
	out.elasticConn.Password = cfg.Password

	if cfg.Protocol != "" {
		out.elasticConn.Protocol = cfg.Protocol
	}

	if cfg.Index != "" {
		out.Index = cfg.Index
	} else {
		out.Index = DefaultIndex
	}

	log.Infof("[ElasticOutput] Using Elasticsearch %s://%s:%s", out.elasticConn.Protocol, out.elasticConn.Domain, out.elasticConn.Port)
	log.Infof("[ElasticOutput] Using index pattern [%s-]YYYY.MM.DD", out.Index)
}

func (out *ElasticOutputType) PublishHostinfo() {
	index := fmt.Sprintf(".%s", DefaultIndex)
	_, err := out.elasticConn.Index(index, "hostinfo")

}

// PublishMetric implements the Output interface.
func (out *ElasticOutputType) PublishMetric(ts time.Time, metric model.Metric) error {
	index := fmt.Sprintf("%s-%s", out.Index, ts.Format("2006.01.02"))
	_, err := out.elasticConn.Index(index, metric["type"].(string), "", nil, metric)
	return err
}
