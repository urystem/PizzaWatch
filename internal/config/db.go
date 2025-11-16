package config

// import (
// 	"fmt"
// 	"os"

// 	"go.yaml.in/yaml/v4"
// )

// type cfgDB struct {
// 	Database struct {
// 		Host     string `yaml:"host"`
// 		Port     int    `yaml:"port"`
// 		User     string `yaml:"user"`
// 		Password string `yaml:"password"`
// 		Database string `yaml:"database"`
// 	} `yaml:"database"`
// }

type CfgDBInter interface {
	GetHostName() string
	GetDBPort() int
	GetUser() string
	GetPassword() string
	GetDBName() string
}

// func GetDBConfig() (CfgDBInter, error) {
// 	f, err := os.Open("config.yaml")
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to open config file: %w", err)
// 	}
// 	defer f.Close()

// 	var cfg cfgDB
// 	decoder := yaml.NewDecoder(f)
// 	if err := decoder.Decode(&cfg); err != nil {
// 		return nil, fmt.Errorf("failed to decode YAML: %w", err)
// 	}
// 	return &cfg, nil
// }

// func (c *cfgDB) GetHostName() string { return c.Database.Host }

// func (c *cfgDB) GetDBPort() int { return c.Database.Port }

// func (c *cfgDB) GetUser() string { return c.Database.User }

// func (c *cfgDB) GetPassword() string { return c.Database.Password }

// func (c *cfgDB) GetDBName() string { return c.Database.Database }
