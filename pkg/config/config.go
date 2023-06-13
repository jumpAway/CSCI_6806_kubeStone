package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Database struct {
		Host         string `json:"host"`
		Port         int    `json:"port"`
		Username     string `json:"username"`
		Password     string `json:"password"`
		DatabaseName string `json:"database_name"`
		Volume       string `json:"volume"`
	} `json:"database"`
	Server struct {
		Port int `json:"port"`
	} `json:"server"`
}
type Server struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ipaddress"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func InitConfig() (config Config, err error) {
	data, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return config, err
	}
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Println("Error parsing config file:", err)
		return config, err
	}
	return config, err
}
