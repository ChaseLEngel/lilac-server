package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	privateKeyError = "Private key file not found"
	publicKeyError  = "Public key file not found"
	logFileError    = "Log file file not found"
	databaseError   = "Database file not found"
	portError       = "No port defined"
	userError       = "No user defined"
	passwordError   = "No password defined"
)

type Config struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
	LogFile    string `json:"log_file"`
	Database   string
	Port       string
	User       string
	Password   string
}

func (c *Config) validate() error {
	if _, err := os.Stat(c.PrivateKey); os.IsNotExist(err) {
		return errors.New(privateKeyError)
	}
	if _, err := os.Stat(c.PublicKey); os.IsNotExist(err) {
		return errors.New(publicKeyError)
	}
	if _, err := os.Stat(c.LogFile); os.IsNotExist(err) {
		return errors.New(logFileError)
	}
	if _, err := os.Stat(c.Database); os.IsNotExist(err) {
		return errors.New(databaseError)
	}
	if c.Port == "" {
		return errors.New(portError)
	}
	c.Port = fmt.Sprintf(":%v", c.Port)
	if c.User == "" {
		return errors.New(userError)
	}
	if c.Password == "" {
		return errors.New(passwordError)
	}
	return nil
}

func Parse(configPath string) (Config, error) {
	config := Config{}
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, err
	}
	json.Unmarshal(data, &config)
	if err := config.validate(); err != nil {
		return config, err
	}
	return config, nil
}
