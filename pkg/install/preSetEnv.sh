#!/bin/bash

SYS_P="/var/lib/kubeStone/offlinePKG/v.1.26.3/"
#K8S_IMG_P="/var/lib/kubeStone/offlinePKG/k8s/img"
HOSTNAME=$1
IPADDRESS=$(ip a s | grep -oP 'inet \K[\d.]+' | grep -v "127.0.0.1")
pkg_containerd="containerd-1.7.1-linux-arm64.tar.gz"
pkg_servicefile="containerd.service"
pkg_cni="cni-plugins-linux-arm64-v1.3.0.tgz"
pkg_nerdctl="nerdctl-1.4.0-linux-arm64.tar.gz"
pkg_runc="runc.arm64"

BLUE='\033[0;34m'
RESET='\033[0m'
echo -e "${BLUE}Verifying the server architecture...${RESET}"
arch=$(uname -m)
expected_arch="aarch64"

if [ "$arch" != "$expected_arch" ]; then
    echo "Server architecture at this version is only supported arm64. Exiting..."
    exit -1
fi
if [ $# != 1 ]; then
    echo "Please attach HOSTNAME"
    exit -1
fi


BLUE='\033[0;34m'
RESET='\033[0m'
echo -e "${BLUE}Presetting the server environment...${RESET}"
preSetEnv=(
  "hostnamectl set-hostname $HOSTNAME"
  "echo \"$IPADDRESS $HOSTNAME\" >> /etc/hosts"
  "swapoff -a"
  "sed -i '/swap/{s/^/#/g}' /etc/fstab"
  "setenforce 0"
  "sed -i 's/^SELINUX=enforcing$/SELINUX=permissive/' /etc/selinux/config"
  "systemctl stop firewalld"
  "systemctl disable firewalld"
  "yum install  -y --cacheonly --disablerepo=* $SYS_P/*.rpm"
)
for cmd in "${preSetEnv[@]}"; do
      echo "Executing $cmd"
      eval $cmd
      if [ $? != 0 ]; then
          echo "Preset environment error"
          exit -1
      fi
  done

BLUE='\033[0;34m'
RESET='\033[0m'
echo -e "${BLUE}Installing containerd...${RESET}"
cat <<EOF | tee /etc/modules-load.d/k8s.conf
   overlay
   br_netfilter
EOF

cat <<EOF | tee /etc/sysctl.d/k8s.conf
    net.bridge.bridge-nf-call-iptables  = 1
    net.bridge.bridge-nf-call-ip6tables = 1
    net.ipv4.ip_forward                 = 1
EOF

ipvs_modules="ip_vs ip_vs_lc ip_vs_wlc ip_vs_rr ip_vs_wrr ip_vs_lblc ip_vs_lblcr ip_vs_dh ip_vs_sh ip_vs_fo ip_vs_nq ip_vs_sed ip_vs_ftp nf_conntrack"
for module in $ipvs_modules; do
    modprobe $module
    if [ $? != 0 ]; then
       echo "Install module $module error"
       exit -1
    fi
done
install_commands=(
  "modprobe overlay"
  "modprobe br_netfilter"
  "modprobe vxlan"
  "sysctl -p /etc/sysctl.d/k8s.conf"
  "containerd config default > /etc/containerd/config.toml"
  "sed -i 's/SystemdCgroup = false/SystemdCgroup = true/g'  /etc/containerd/config.toml"
  "crictl config runtime-endpoint /run/containerd/containerd.sock"
  "systemctl daemon-reload"
  "systemctl enable --now containerd"
  `echo "source <(kubectl completion bash)" >> ~/.bashrc`
)

  for cmd in "${install_commands[@]}"; do
      echo "Executing $cmd"
      eval $cmd
      if [ $? != 0 ]; then
          echo "Install containerd error"
          exit -1
      fi
  done

BLUE='\033[0;34m'
RESET='\033[0m'
echo -e "${BLUE}Verifying installation...${RESET}"
systemctl status containerd 1> /dev/null
if [ $? != 0 ]; then
      echo "Install containerd unsuccessful"
      exit -1
fi

GREEN='\033[0;32m'
RESET='\033[0m'
echo -e "${GREEN}Install containerd successful${RESET}"


cat <<EOF | tee /etc/NetworkManager/conf.d/calico.conf
[keyfile]
unmanaged-devices=interface-name:cali*;interface-name:tunl*;interface-name:vxlan.calico;interface-name:vxlan-v6.calico;interface-name:wireguard.cali;interface-name:wg-v6.cali
EOF

