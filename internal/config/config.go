package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to read user's home directory: %w", err)
	}
	return fmt.Sprintf("%s/%s", home, configFileName), nil
}

func Read() (Config, error) {
	file, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("unable to get user's home directory: %w", err)
	}
	data, rErr := os.ReadFile(file)
	if rErr != nil {
		return Config{}, fmt.Errorf("unable to read json file: %w", err)
	}
	var config Config
	if jErr := json.Unmarshal(data, &config); jErr != nil {
		return Config{}, fmt.Errorf("error unmarshalling JSON: %w", err)
	}
	return config, nil
}

func (cfg Config) SetUser(username string) error {
	cfg.CurrentUserName = username
	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("unable to marshal json: %w", err)
	}
	file, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("unable to get user's home directory: %w", err)
	}
	err = os.WriteFile(file, data, 777)
	if err != nil {
		return fmt.Errorf("unable to write to file: %w", err)
	}
	return nil
}
