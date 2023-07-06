package GPT

import (
	"bufio"
	yaml2 "github.com/ghodss/yaml"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var HistoryMutex = &sync.Mutex{}
var HistoryMap = make(map[string][]map[string]string)

func ExtractFile(answer string) ([]string, []string) {
	var yamlGot []string
	var cmdGot []string
	yaml := ""
	cmd := ""
	scanner := bufio.NewScanner(strings.NewReader(answer))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "```yaml" {
			for scanner.Scan() {
				line := scanner.Text()
				if line == "```" {
					yamlGot = append(yamlGot, yaml)
					yaml = ""
					break
				}
				yaml = yaml + line + "\n"
			}
		}
		if line == "```" || line == "```shell" {
			for scanner.Scan() {
				line := scanner.Text()
				if line == "```" {
					cmdGot = append(cmdGot, cmd)
					cmd = ""
					break
				}
				cmd = cmd + line + "\n"
			}
		}
	}
	return yamlGot, cmdGot
}

func ExecuteYaml(json string) error {
	yaml, err := yaml2.JSONToYAML([]byte(json))
	if err != nil {
		return err
	}
	cmd := exec.Command("kubectl", "config", "use-context", "context1")
	err = cmd.Run()
	if err != nil {
		return err
	}
	filePath := "/tmp/kubeStone_gpt.yaml"
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	_, err = file.WriteString(string(yaml))
	if err != nil {
		return err
	}
	cmd = exec.Command("kubectl", "apply", "-f", "/tmp/kubeStone_gpt.yaml")
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func ExecuteKubectl(kubeCmd string) error {
	cmd := exec.Command("kubectl", "config", "use-context", "context1")
	args := strings.Fields(kubeCmd)
	err := cmd.Run()
	if err != nil {
		return err
	}
	cmd = exec.Command("kubectl", args[1:]...)
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
