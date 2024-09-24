package bootstrap

import (
	"fmt"
	"time"

	"github.com/jinzhu/configor"
)

// Config ...
type Config struct {
	App  App  `yaml:"app"`
	HTTP HTTP `yaml:"http"`
}

type App struct {
	AppName string `yaml:"app_name" required:"true"`
	Region  string `yaml:"region" required:"true"`
	APIPort string `yaml:"api_port" required:"true" default:":8080"`
}

type HTTP struct {
	Timeout          time.Duration `yaml:"timeout" default:"2s"`
	ACMEClient       Client        `yaml:"acme_client"`
	PatagoniaClient  Client        `yaml:"patagonia_client"`
	PaperfliesClient Client        `yaml:"paperflies_client"`
}

type Client struct {
	EndPoint string `yaml:"endpoint" required:"true"`
}

type MergeStrategy struct {
}

func LoadConfig(path string) (*Config, error) {
	specs := Config{}
	err := configor.Load(&specs, path)
	if err != nil {
		return nil, fmt.Errorf("failed to load Config with the error : %w while loading from config file", err)
	}
	return &specs, nil
}
