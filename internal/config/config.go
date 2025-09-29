package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DB_URL            string `json:"db_url"`
	Current_User_Name string `json:"current_user_name"`
}

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	fileHandle, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer fileHandle.Close()
	var cfg Config
	decoder := json.NewDecoder(fileHandle)
	if err = decoder.Decode(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (cfg *Config) SetUser(user_name string) error {
	cfg.Current_User_Name = user_name
	if err := write(*cfg); err != nil {
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	home_path, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := home_path + "/" + configFileName
	return path, nil
}

func write(cfg Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	if err = encoder.Encode(cfg); err != nil {
		return err
	}
	return nil
}
