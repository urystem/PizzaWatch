package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type database struct {
	host     string
	port     int
	user     string
	password string
	database string
}

func (d *database) GetHostName() string {
	return d.host
}

func (d *database) GetDBPort() int {
	return d.port
}

func (d *database) GetUser() string {
	return d.user
}

func (d *database) GetPassword() string {
	return d.password
}

func (d *database) GetDBName() string {
	return d.database
}

func GetDBConfig() (CfgDBInter, error) {
	f, err := os.Open("config.yaml")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := &database{}
	scanner := bufio.NewScanner(f)

	var section string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// пустые строки и комментарии игнорируем
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// секция: database:
		if strings.HasSuffix(line, ":") {
			section = strings.TrimSuffix(line, ":")
			continue
		}

		// ключ: значение
		if section == "database" {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				continue
			}

			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])

			// убираем кавычки, если есть
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
			case "database":
				cfg.database = val
			}
		}
	}

	return cfg, scanner.Err()
}
