package k8s

import (
	"crypto/tls"
	"kubeStone/pkg/config"
	"net/http"
)

func GetPods(cluster config.Cluster, token string, namespace string) (*http.Response, error) {
	var req *http.Request
	var err error
	if namespace == "" {
		req, err = http.NewRequest("GET", "https://"+cluster.Master+":6443/api/v1/pods", nil)
	} else {
		req, err = http.NewRequest("GET", "https://"+cluster.Master+":6443/api/v1/namespaces/"+namespace+"/pods", nil)
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
