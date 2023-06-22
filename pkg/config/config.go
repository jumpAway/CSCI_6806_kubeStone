package config

import (
	"encoding/json"
	"fmt"
	"os"
)

/*
The Config struct represents a configuration object.
It is used to hold configurations for a software application.
*/
type Config struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Database struct {
		//The hostname or IP address of the database server
		Host string `json:"host"`
		//The port number on which the database server is listening
		Port         int    `json:"port"`
		Username     string `json:"username"`
		Password     string `json:"password"`
		DatabaseName string `json:"database_name"`
		Volume       string `json:"volume"`
	} `json:"database"`
	Server struct {
		//The port number on which the application server should listen
		Port int `json:"port"`
	} `json:"server"`
}

// The Server struct encapsulates information about a server.
type Server struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ipaddress"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

/* InitConfig is a function that reads a configuration file named "config.json",
deserializes the JSON data into a Config struct, and returns it.*/
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
