package host

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"kubeStone/m/v2/pkg/config"
	"time"
)

// ConnectSer is a function that establishes an SSH connection to a remote server using the provided configuration.
func ConnectSer(server config.Server) error {
	sshCfg := &ssh.ClientConfig{
		User: server.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(server.Password),
		},
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	_, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", server.IP, server.Port), sshCfg)
	if err != nil {
		return err
	}
	return nil
}
