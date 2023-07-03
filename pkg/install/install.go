package install

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"kubeStone/pkg/config"
	"kubeStone/pkg/encrypt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

func ExecCmd(commands []string, server config.Server) error {
	logFile, err := os.OpenFile("/var/log/kubeStone_"+server.IP+".log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	log.SetOutput(logFile)

	sshCfg := &ssh.ClientConfig{
		User: server.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(encrypt.GetMD5Hash(server.Password)),
		},
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", server.IP, server.Port), sshCfg)
	if err != nil {
		return err
	}

	for _, cmd := range commands {
		session, err := client.NewSession()
		if err != nil {
			log.Printf("Failed to create session: %v\n", err)
			return err
		}
		var stdoutBuf bytes.Buffer
		session.Stdout = &stdoutBuf
		err = session.Run(cmd)
		if err != nil {
			log.Printf("Failed to run command[%s]: %s\n", cmd, stdoutBuf.String())
			return err
		} else {
			output := fmt.Sprintf("Executing [%s] on server [%s]:%s\n", cmd, server.IP, stdoutBuf.String())
			log.Print(output)
		}
	}

	return nil
}

func SetNode(node config.Server, master config.Server, seq int) error {
	node.Hostname = "node" + strconv.Itoa(seq)
	joinCmd, err := exec.Command("ssh", master.Username+"@"+master.IP, "kubeadm", "token", "create", "--print-join-command").Output()
	if err != nil {
		return err
	}
	err = exec.Command("ssh", master.Username+"@"+master.IP, "echo", node.IP+" "+node.Hostname, ">>", "/etc/hosts").Run()
	if err != nil {
		return err
	}
	pkgPath := "/var/kubeStone"
	scp := exec.Command("scp", "-r", pkgPath, node.Username+"@"+node.IP+":"+pkgPath)
	if err = scp.Run(); err != nil {
		log.Printf("Failed to copy file: %v\n", err)
		return err
	}
	err = PreSetEnv(node)
	if err != nil {
		return err
	}
	commands := []string{
		"echo \"" + master.IP + " master\" >> /etc/hosts",
		string(joinCmd),
	}
	err = ExecCmd(commands, node)
	if err != nil {
		return err
	}
	return nil
}

func SetMaster(server config.Server, master config.ClusterInfo) error {
	initYaml, err := os.ReadFile("/var/kubeStone/init.yaml")
	if err != nil {
		return err
	}

	reAdvertiseAddress := regexp.MustCompile(`(?m)(advertiseAddress:\s*)`)
	reServiceSubnet := regexp.MustCompile(`(?m)(serviceSubnet:\s*)`)
	rePodSubnet := regexp.MustCompile(`(?m)(podSubnet:\s*)`)
	reMode := regexp.MustCompile(`(?m)(mode:\s*)`)
	//	reVersion := regexp.MustCompile(`(?m)(kubernetesVersion:\s*)`)

	newYAML := reAdvertiseAddress.ReplaceAll(initYaml, []byte("advertiseAddress: "+master.MasterIp+"\n  "))
	newYAML = reServiceSubnet.ReplaceAll(newYAML, []byte("serviceSubnet: "+master.ServiceSubnet+"\n  "))
	newYAML = rePodSubnet.ReplaceAll(newYAML, []byte("podSubnet: "+master.PodSubnet+"\n"))
	newYAML = reMode.ReplaceAll(newYAML, []byte("mode: "+master.ProxyMode+"\n"))
	//	newYAML = reVersion.ReplaceAll(newYAML, []byte("mode: "+master.Version+"\n"))

	err = os.WriteFile("/var/kubeStone/"+master.MasterIp+"_init.yaml", newYAML, 0644)
	if err != nil {
		return err
	}

	calicoYaml, err := os.ReadFile("/var/kubeStone/calico-resources.yaml")
	if err != nil {
		return err
	}
	reCidr := regexp.MustCompile(`(?m)(advertiseAddress:\s*)`)
	newYAML = reCidr.ReplaceAll(calicoYaml, []byte("cidr: "+master.PodSubnet+"\n      "))
	err = os.WriteFile("/var/kubeStone/"+master.MasterIp+"_calico-resources.yaml", newYAML, 0644)
	if err != nil {
		return err
	}

	pkgPath := "/var/kubeStone"
	scp := exec.Command("scp", "-r", pkgPath, server.Username+"@"+server.IP+":"+pkgPath)
	if err = scp.Run(); err != nil {
		log.Printf("Failed to copy file: %v\n", err)
		return err
	}

	err = PreSetEnv(server)
	if err != nil {
		return err
	}
	commands := []string{
		"kubeadm init --config=/var/kubeStone/" + master.MasterIp + "_init.yaml",
		"mkdir -p $HOME/.kube",
		"cp -i /etc/kubernetes/admin.conf $HOME/.kube/config",
		"export KUBECONFIG=/etc/kubernetes/admin.conf",
		"chown $(id -u):$(id -g) $HOME/.kube/config",
		"kubectl create -f /var/kubeStone/tigera-operator.yaml",
		"kubectl create -f /var/kubeStone/" + master.MasterIp + "_calico-resources.yaml",
		"kubectl apply -f /var/kubeStone/ks-serviceaccount.yaml",
		"kubectl apply -f /var/kubeStone/ks-rolebings.yaml",
	}
	err = ExecCmd(commands, server)
	if err != nil {
		return err
	}
	token, err := exec.Command("ssh", server.Username+"@"+server.IP, "kubectl", "create", "token", "kubestone-service-account", "--duration", "87600h").Output()
	if err != nil {
		return err
	}
	err = exec.Command("kubectl", "config", "set-cluster", "cluster1", "--server=https://"+master.MasterIp+":6443", "--insecure-skip-tls-verify=true").Run()
	if err != nil {
		return err
	}
	err = exec.Command("kubectl", "config", "set-credentials", "kubestone-service-account", "--token="+string(token)).Run()
	if err != nil {
		return err
	}
	err = exec.Command("kubectl", "config", "set-context", "context1", "--cluster=cluster1", "--user=kubestone-service-account").Run()
	if err != nil {
		return err
	}
	err = exec.Command("kubectl", "config", "use-context", "context1").Run()
	if err != nil {
		return err
	}
	return nil
}

func PreSetEnv(server config.Server) error {
	var SysP string
	var ipvsModules string
	SysP = "/var/kubeStone/offlinePKG/v.1.26.3/"
	ipvsModules = "ip_vs ip_vs_lc ip_vs_wlc ip_vs_rr ip_vs_wrr ip_vs_lblc ip_vs_lblcr ip_vs_dh ip_vs_sh ip_vs_fo ip_vs_nq ip_vs_sed ip_vs_ftp nf_conntrack overlay br_netfilter vxlan"
	commands := []string{
		"hostnamectl set-hostname " + server.Hostname,
		"echo \"" + server.IP + " " + server.Hostname + "\" >> /etc/hosts",
		"swapoff -a",
		"sed -i '/swap/{s/^/#/g}' /etc/fstab",
		"setenforce 0",
		"sed -i 's/^SELINUX=enforcing$/SELINUX=permissive/' /etc/selinux/config",
		"systemctl stop firewalld",
		"systemctl disable firewalld",
		"yum install  -y --cacheonly --disablerepo=* " + SysP + "/*.rpm",
		"cat <<EOF | tee /etc/modules-load.d/k8s.conf\n   overlay\n   br_netfilter\nEOF",
		"cat <<EOF | tee /etc/sysctl.d/k8s.conf\n    net.bridge.bridge-nf-call-iptables  = 1\n    net.bridge.bridge-nf-call-ip6tables = 1\n    net.ipv4.ip_forward                 = 1\nEOF",
		"cat <<EOF | tee /etc/NetworkManager/conf.d/calico.conf\n[keyfile]\nunmanaged-devices=interface-name:cali*;interface-name:tunl*;interface-name:vxlan.calico;interface-name:vxlan-v6.calico;interface-name:wireguard.cali;interface-name:wg-v6.cali\nEOF",
		"modprobe -a " + ipvsModules,
		"sysctl -p /etc/sysctl.d/k8s.conf",
		"containerd config default > /etc/containerd/config.toml",
		"sed -i 's/SystemdCgroup = false/SystemdCgroup = true/g'  /etc/containerd/config.toml",
		"crictl config runtime-endpoint /run/containerd/containerd.sock",
		"systemctl daemon-reload",
		"systemctl enable --now containerd",
		"echo \"source <(kubectl completion bash)\" >> ~/.bashrc",
		"systemctl restart chronyd",
	}
	if err := ExecCmd(commands, server); err != nil {
		return err
	}
	return nil
}
