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

type ClusterInfo struct {
	MasterIp      string `json:"Server,omitempty"`
	Version       string `json:"Version,omitempty"`
	ServiceSubnet string `json:"ServiceSubnet,omitempty"`
	PodSubnet     string `json:"PodSubnet,omitempty"`
	ProxyMode     string `json:"ProxyMode,omitempty"`
	NodeIp        string `json:"ServerAddress,omitempty"`
}

type Cluster struct {
	ClusterName   string `json:"ClusterName,omitempty"`
	Version       string `json:"Version,omitempty"`
	CNI           string `json:"CNI,omitempty"`
	ServiceSubnet string `json:"ServiceSubnet,omitempty"`
	PodSubnet     string `json:"PodSubnet,omitempty"`
	ProxyMode     string `json:"ProxyMode,omitempty"`
	Master        string `json:"Master,omitempty"`
	Node          string `json:"Node,omitempty"`
	Context       string `json:"Context,omitempty"`
}

type GPTResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	} `json:"choices"`
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
