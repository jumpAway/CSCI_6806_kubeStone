package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const ScriptPath = "/root/kubeStone/install/"

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

/*
ClusterInfo is a struct that represents information about a cluster.

The struct fields represent various details about the cluster, and are tagged with JSON keys,

which will be used when the struct is serialized to JSON or deserialized from JSON.
*/
type ClusterInfo struct {
	MasterIp      string `json:"Server,omitempty"`
	Version       string `json:"Version,omitempty"`
	ServiceSubnet string `json:"ServiceSubnet,omitempty"`
	PodSubnet     string `json:"PodSubnet,omitempty"`
	ProxyMode     string `json:"ProxyMode,omitempty"`
	NodeIp        string `json:"ServerAddress,omitempty"`
}

/*
	InitConfig is a function that reads a configuration file named "config.json",

deserializes the JSON data into a Config struct, and returns it.
*/
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
