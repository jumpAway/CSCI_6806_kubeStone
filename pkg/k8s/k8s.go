package k8s

import (
	"crypto/tls"
	"kubeStone/pkg/config"
	"net/http"
)

func GetRes(cluster config.Cluster, token string, namespace string, resource string) (*http.Response, error) {
	var req *http.Request
	var err error
	var path string
	if namespace == "" {
		switch resource {
		case "pod":
			path = "/api/v1/pods"
		case "service":
			path = "/api/v1/services"
		case "deployment":
			path = "/apis/apps/v1/deployments"
		case "daemonset":
			path = "/apis/apps/v1/daemonsets"
		default:
		}
	} else {
		switch resource {
		case "pod":
			path = "/api/v1/namespaces/" + namespace + "/pods"
		case "service":
			path = "/api/v1/namespaces/" + namespace + "/services"
		case "deployment":
			path = "/apis/apps/v1/namespaces/" + namespace + "/deployments"
		case "daemonset":
			path = "/apis/apps/v1/namespaces/" + namespace + "/daemonsets"
		default:
		}
	}
	req, err = http.NewRequest("GET", "https://"+cluster.Master+":6443"+path, nil)
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
