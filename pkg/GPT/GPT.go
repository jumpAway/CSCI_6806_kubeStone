package GPT

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var HistoryMutex = &sync.Mutex{}
var HistoryMap = make(map[string][]map[string]string)

func extractFile(answer string) ([]string, []string) {
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
		if line == "```" {
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

func ExecuteGPT(answer string) ([]string, []string, error) {
	yamlGot, cmdGot := extractFile(answer)
	if yamlGot[0] == "" {

	} else {
		cmd := exec.Command("kubectl", "config", "use-context", "context1")
		err := cmd.Run()
		if err != nil {
			return yamlGot, cmdGot, err
		}
		for _, yaml := range yamlGot {
			filePath := "/tmp/kubeStone_gpt.yaml"
			file, err := os.Create(filePath)
			if err != nil {
				return yamlGot, cmdGot, err
			}
			_, err = file.WriteString(yaml)
			if err != nil {
				return yamlGot, cmdGot, err
			}
			cmd = exec.Command("kubectl", "apply", "-f", "/tmp/kubeStone_gpt.yaml")
			err = cmd.Run()
			if err != nil {
				return yamlGot, cmdGot, err
			}
		}
	}
	return yamlGot, cmdGot, nil
}
