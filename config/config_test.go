package config

import (
	"reflect"
	"testing"
	"time"

	"github.com/ultimatums/ultimate/model"

	"gopkg.in/yaml.v2"
)

var expConf = &Config{
	GlobalConfig: &GlobalConfig{
		FetchInterval: model.Duration(15 * time.Second),
	},
	OutputConfig: &OutputConfig{
		Elasticsearch: &ElasticsearchConfig{
			Host: "192.168.2.78",
			Port: 9200,
		},
	},
	TaskConfigs: []*TaskConfig{
		{
			TaskName:      "host",
			FetchInterval: model.Duration(5 * time.Second),
			UnitConfigs: []*UnitConfig{
				{
					UnitName: "cpu",
				},
				{
					UnitName: "mem",
				},
				{
					UnitName: "diskio",
				},
				{
					UnitName: "network",
				},
			},
		},
		{
			TaskName:      "container",
			FetchInterval: model.Duration(10 * time.Second),
			TaskTags: model.TagMap{
				"docker_endpoint": "unix:///var/run/docker.sock",
				"key1":            "value1",
			},
			UnitConfigs: []*UnitConfig{
				{
					UnitName:      "02e1f960f516",
					FetchInterval: model.Duration(5 * time.Second),
				},
				{
					UnitName:      "78b0817479ce",
					FetchInterval: model.Duration(6 * time.Second),
				},
			},
		},
	},
}

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig("testdata/test.yml")
	if err != nil {
		t.Fatalf("Error parsing %s: %s", "testdata/test.yml", err)
	}

	bgot, err := yaml.Marshal(cfg)
	if err != nil {
		t.Fatalf("%s", err)
	}

	bexp, err := yaml.Marshal(expConf)
	if err != nil {
		t.Fatalf("%s", err)
	}
	expConf.raw = cfg.raw

	if !reflect.DeepEqual(cfg, expConf) {
		t.Fatalf("%s: expected: \n\n%s\n but got: \n\n%s\n", "testdata/test.yml", bexp, bgot)
	}
}
