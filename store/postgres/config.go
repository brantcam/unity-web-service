package postgres

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Host     string
	Name     string
	Username string
	Password string
	Port     uint
}

func LoadConfigFromEnv() (*Config, error) {
	var cfg Config

	if host, ok := os.LookupEnv("DB_HOST"); ok {
		cfg.Host = host
	}

	if name, ok := os.LookupEnv("DB_NAME"); ok {
		cfg.Name = name
	}

	if username, ok := os.LookupEnv("DB_USER"); ok {
		cfg.Username = username
	}

	if password, ok := os.LookupEnv("DB_PASS"); ok {
		cfg.Password = password
	}

	if port, ok := os.LookupEnv("DB_PORT"); ok {
		p, err := strconv.ParseUint(port, 10, 16)
		if err != nil {
			return nil, fmt.Errorf("invalid port value: %v", err)
		}
		cfg.Port = uint(p)
	}

	return &cfg, validateConfig(cfg)
}

func validateConfig(c Config) error {
	var missing []string

	if c.Host == "" {
		missing = append(missing, "Host")
	}
	if c.Name == "" {
		missing = append(missing, "Name")
	}
	if c.Username == "" {
		missing = append(missing, "Username")
	}
	if c.Password == "" {
		missing = append(missing, "Password")
	}
	if c.Port == 0 {
		missing = append(missing, "Port")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required fields: %s", missing)
	}

	return nil
}
