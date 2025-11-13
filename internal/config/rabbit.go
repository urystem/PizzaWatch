package config

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"
)

type cfgRabbitMQ struct {
	RabbitMQ struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"rabbitmq"`
}

type CfgRabbitInter interface {
	GetHostName() string
	GetDBPort() int
	GetUser() string
	GetPassword() string
}

func GetRabbitMQConfig() (CfgRabbitInter, error) {
	f, err := os.Open("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer f.Close()

	var cfg cfgRabbitMQ
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode YAML: %w", err)
	}
	return &cfg, nil
}

func (c *cfgRabbitMQ) GetHostName() string { return c.RabbitMQ.Host }

func (c *cfgRabbitMQ) GetDBPort() int { return c.RabbitMQ.Port }

func (c *cfgRabbitMQ) GetUser() string { return c.RabbitMQ.User }

func (c *cfgRabbitMQ) GetPassword() string { return c.RabbitMQ.Password }

