package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get user home directory: %w", err)
	}
	return filepath.Join(homeDir, configFileName), nil
}

func Read() (Config, error) {
	cfg := Config{}
	filePath, err := getConfigFilePath()
	if err != nil {
		return cfg, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		// If file doesn't exist, return the zero-value Config without error
		// as the file might be created later by a write operation.
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("could not open config file %s: %w", filePath, err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return cfg, fmt.Errorf("could not stat config file %s: %w", filePath, err)
	}

	// If file is empty, return zero-value Config, nothing to decode
	if fileInfo.Size() == 0 {
		return cfg, nil
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		// Check if it's just an empty JSON object or unexpected EOF
		// which decoder might return error for depending on exact content
		if err.Error() == "EOF" {
			// Consider an empty file or empty JSON object as valid empty config
			return Config{}, nil
		}
		return cfg, fmt.Errorf("could not decode config file %s: %w", filePath, err)
	}

	return cfg, nil
}

func write(cfg Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("could not create/open config file %s for writing: %w", filePath, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	// Ensure pretty printing matching the assignment's example structure
	encoder.SetIndent("", "  ")
	err = encoder.Encode(cfg)
	if err != nil {
		return fmt.Errorf("could not encode config to file %s: %w", filePath, err)
	}

	return nil
}

func (c *Config) SetUser(name string) error {
	c.CurrentUserName = name
	return write(*c)
}
