package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) GetDBURL() string {
	return c.DBURL
}

func (c *Config) GetUser() string {
	return c.CurrentUserName
}

func (c *Config) SetUser(user string) error {
	c.CurrentUserName = user
	return c.write()
}

func (c *Config) String() string {
	return fmt.Sprintf("db_url: %s, current_user_name: %s", c.DBURL, c.CurrentUserName)
}
func (c *Config) write() error {
	configPath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("error while calling getConfigFilePath: %s", err)
	}

	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("could not create %s: %s", configPath, err)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	if err := enc.Encode(c); err != nil {
		return fmt.Errorf("could not encode the config to JSON: %s", err)
	}

	return nil
}

func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not fetch user home directory: %s", err)
	}
	configPath := path.Join(homePath, configFileName)
	return configPath, nil
}

func Read() (Config, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("error while calling getConfigFilePath: %s", err)
	}

	file, err := os.Open(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("could not open %s: %s", configPath, err)
	}
	defer file.Close()

	var config Config
	dec := json.NewDecoder(file)
	if err := dec.Decode(&config); err != nil {
		return Config{}, fmt.Errorf("could not decode the content of %s: %s", configPath, err)
	}

	return config, nil
}
