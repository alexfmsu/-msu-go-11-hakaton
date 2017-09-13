package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Broker struct {
	ID   int32  `yaml:"id" json:"id"`
	Name string `yaml:"name" json:"name"`
}

type Ticker struct {
	ID   int32  `yaml:"id" json:"id"`
	Name string `yaml:"name" json:"name"`
}

type Config struct {
	Brokers  []Broker `yaml:"brokers"`
	Tickers  []Ticker `yaml:"tickers"`
	Port     string   `yaml:"port"`
	GrpcPort string   `yaml:"grpcPort"`
	Host     string   `yaml:"host"`
	Timeout  int      `yaml:"timeout"`
}

func ParseConfigFromFile() *Config {
	data, err := ioutil.ReadFile("src/broker/config/app.yaml")
	if err != nil {
		log.Printf("Config: Error: ParseConfigFromFile: %v", err)
		return nil
	}

	var config Config
	if err := config.Parse(data); err != nil {
		log.Printf("Config: Error: ParseConfigFromFile: %v", err)
		return nil
	}

	return &config
}

func (c *Config) Parse(data []byte) error {
	if err := yaml.Unmarshal(data, c); err != nil {
		return err
	}
	return nil
}
