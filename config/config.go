package config

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/ultimatums/ultimate/model"

	"gopkg.in/yaml.v2"
)

var (
	DefaultConfig = Config{
		GlobalConfig: DefaultGlobalConfig,
	}

	DefaultGlobalConfig = GlobalConfig{
		FetchInterval: model.Duration(10 * time.Second),
	}

	DefaultTaskConfig = TaskConfig{}
)

func LoadConfig(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = yaml.Unmarshal(content, cfg)
	if err != nil {
		return nil, err
	}
	cfg.raw = string(content)
	return cfg, err
}

// Config is the top_level configuration.
type Config struct {
	GlobalConfig GlobalConfig  `yaml:"global"`
	FetchConfigs []*TaskConfig `yaml:"fetch_configs,omitempty"`

	//raw is the orginal content from the configuration file.
	raw string
}

func (c Config) String() string {
	if c.raw != "" {
		return c.raw
	}
	b, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

// GlobalConfig configures global environment.
type GlobalConfig struct {
	FetchInterval model.Duration `yaml:"fetch_interval,omitempty"`
}

// TaskConfig configures a fetching task.
type TaskConfig struct {
	TaskName      string         `yaml:"task_name"`
	FetchInterval model.Duration `yaml:"fetch_interval,omitempty"`
	TaskTags      model.TagMap   `yaml:"task_tags,omitempty"`
	UnitSets      []*UnitSet     `yaml:"unit_sets,omitempty"`
	//	ContainerSets []*ContainerSet `yaml:"container_sets,omitempty"`
}

type UnitSet struct {
	UnitTags      model.TagMap   `yaml:"unit_tags,omitempty"`
	FetchInterval model.Duration `yaml:"fetch_interval,omitempty"`
	Identity      string
}

/*
type ContainerSet struct {
	DockerEndpoint string            `yaml:"docker_endpoint"`
	Containers     []ContainerConfig `yaml:"containers,omitempty"`
}

type ContainerConfig struct {
	ID            string         `yaml:"id"`
	FetchInterval model.Duration `yaml:"fetch_interval,omitempty"`
}
*/
