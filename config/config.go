package config

import (
	"fmt"
	"os"
	"strconv"
)

const (
	CONFIG_MQ = "MQ"
	CONFIG_DB = "DB"
)

type Config struct {
	Host     string
	Name     string
	Username string
	Password string
	Port     uint
}

// since both the mq and postgres have similar env we are using the same load func for both
// this needs is a type in the form of a string "DB" or "MQ"
func LoadConfigFromEnv(t string) (*Config, error) {
	var cfg Config

	if host, ok := os.LookupEnv(fmt.Sprintf("%s_HOST", t)); ok {
		cfg.Host = host
	}

	if name, ok := os.LookupEnv(fmt.Sprintf("%s_NAME", t)); ok {
		cfg.Name = name
	}

	if username, ok := os.LookupEnv(fmt.Sprintf("%s_USER", t)); ok {
		cfg.Username = username
	}

	if password, ok := os.LookupEnv(fmt.Sprintf("%s_PASS", t)); ok {
		cfg.Password = password
	}

	if port, ok := os.LookupEnv(fmt.Sprintf("%s_PORT", t)); ok {
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