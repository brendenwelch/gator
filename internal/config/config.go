package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func (c *Config) SetUser(user string) error {
	c.Current_user_name = user
	if err := Write(c); err != nil {
		return err
	}
	return nil
}

func Read() (Config, error) {
	var cfg Config
	path, err := getConfigPath()
	if err != nil {
		return cfg, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func Write(cfg *Config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	path, err := getConfigPath()
	if err != nil {
		return err
	}
	if err = os.WriteFile(path, data, 0777); err != nil {
		return err
	}
	return nil
}

// -- Helpers
func getConfigPath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return dir + "/.gatorconfig.json", nil
}
