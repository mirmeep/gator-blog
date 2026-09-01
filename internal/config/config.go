package config

import (
	"os"
	"fmt"
	"encoding/json"
)

type Config struct {
	Url string `json:"db_url"`
	Username string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	home_dir, err := os.UserHomeDir()

	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	config_path := home_dir + "/.gatorconfig.json"

	return config_path, nil
}

func Read() (Config, error) {
	config_file, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("%w", err)
	}

	data, err := os.ReadFile(config_file)
	if err != nil {
		fmt.Errorf("%w", err)
	}

	var conf Config
	if err := json.Unmarshal(data, &conf); err != nil {
		return Config{}, err
	}

	return conf, nil
}

func (c *Config) SetUser(username string) error {
	config_file, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	c.Username = username

	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	write_err := os.WriteFile(config_file, data, 0660)
	if write_err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

