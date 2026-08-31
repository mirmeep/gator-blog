package config

import (
	"os"
	"errors"
	"encoding/json"
)

type Config struct {
	Url string `json:"db_url"`
	Username string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	home_dir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("cannot find home directory")
	}
	config_path := home_dir + ".gatorconfig.json"
	return config_path, nil
}

func Read() (Config, error) {
	config_file, err := getConfigFilePath()
	if err != nil {
		return Config{}, errors.New("cannot find file path")
	}

	data, err := os.ReadFile(config_file)
	if err != nil {
		errors.New("cannot read file")
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
		return errors.New("cannot get config file path")
	}

	data, err := json.Marshal(c)
	if err != nil {
		return errors.New("cannot marshal json data")
	}

	write_err := os.WriteFile(config_file, data, 0660)
	if write_err != nil {
		return errors.New("cannot write to file")
	}
	return nil
}

