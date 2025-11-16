package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type rabbitmq struct {
	host     string
	port     int
	user     string
	password string
}

func (r *rabbitmq) GetHostName() string {
	return r.host
}

func (r *rabbitmq) GetDBPort() int {
	return r.port
}

func (r *rabbitmq) GetUser() string {
	return r.user
}

func (r *rabbitmq) GetPassword() string {
	return r.password
}

func GetRabbitMQConfig() (CfgRabbitInter, error) {
	f, err := os.Open("config.yaml")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := &rabbitmq{}
	scanner := bufio.NewScanner(f)

	var section string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// пропускаем пустые строки и комментарии
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// rabbitmq:
		if strings.HasSuffix(line, ":") {
			section = strings.TrimSuffix(line, ":")
			continue
		}

		if section == "rabbitmq" {
			keyVal := strings.SplitN(line, ":", 2)
			if len(keyVal) != 2 {
				continue
			}

			key := strings.TrimSpace(keyVal[0])
			val := strings.TrimSpace(keyVal[1])
			val = strings.Trim(val, `"'`)

			switch key {
			case "host":
				cfg.host = val
			case "port":
				n, err := strconv.Atoi(val)
				if err != nil {
					return nil, err
				}
				cfg.port = n
			case "user":
				cfg.user = val
			case "password":
				cfg.password = val
			}
		}
	}

	return cfg, scanner.Err()
}
