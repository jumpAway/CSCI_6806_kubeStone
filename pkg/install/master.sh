#!/bin/bash

if [ $# != 4 ]; then
    exit -1
fi

token=$(kubeadm token generate)
advertiseAddress=$1
serviceSubnet=$2
podSubnet=$3
mode=$4
yamlPath="/root/kubeStone/install"
newYaml=${advertiseAddress}.yaml
log=/root/master_${advertiseAddress}.log

cd $yamlPath
cp init.yaml $newYaml

sedCMD=(
  `sed -i "s/^ *token:.*/    token: ${token}/g" $newYaml`
  `sed -i "s/^ *advertiseAddress:.*/  advertiseAddress: ${advertiseAddress}/g" $newYaml`
  `sed -i "s#^ *serviceSubnet:.*#  serviceSubnet: ${serviceSubnet}#g" $newYaml`
  `sed -i "s#^ *podSubnet:.*#  podSubnet: ${podSubnet}#g" $newYaml`
  `sed -i "s/^mode:.*/mode: ${mode}/g" $newYaml`
  `ssh root@${advertiseAddress} "mkdir -p /var/lib/kubeStone/offlinePKG"`
  `scp -r /root/kubeStone/offlinePKG/* root@${advertiseAddress}:/var/lib/kubeStone/offlinePKG/`
  `scp $newYaml root@${advertiseAddress}:/var/lib/kubeStone`
  `scp preSetEnv.sh root@${advertiseAddress}:/var/lib/kubeStone`
  `scp ks-serviceaccount.yaml root@${advertiseAddress}:/var/lib/kubeStone`
  `scp ks-rolebings.yaml root@${advertiseAddress}:/var/lib/kubeStone`
)

for cmd in "${sedCmd[@]}"; do
      eval $cmd
      if [ $? != 0 ]; then
          exit -1
      fi
  done

masterInstall=(
  `ssh root@${advertiseAddress} "sh /var/lib/kubeStone/preSetEnv.sh master" > $log`
  `ssh root@${advertiseAddress} "kubeadm init --config=/var/lib/kubeStone/${newYaml}" >> $log`
  `ssh root@${advertiseAddress} "mkdir -p $HOME/.kube" >> $log`
  `ssh root@${advertiseAddress} "cp -i /etc/kubernetes/admin.conf $HOME/.kube/config" >> $log`
  `ssh root@${advertiseAddress} "export KUBECONFIG=/etc/kubernetes/admin.conf" >> $log`
  `ssh root@${advertiseAddress} "chown $(id -u):$(id -g) $HOME/.kube/config" >> $log`
#  `ssh root@${advertiseAddress} "kubectl apply -f https://github.com/flannel-io/flannel/releases/latest/download/kube-flannel.yml" >> $log`
  `ssh root@${advertiseAddress} "kubectl create -f https://raw.githubusercontent.com/projectcalico/calico/v3.26.0/manifests/tigera-operator.yaml" >> $log`
  `ssh root@${advertiseAddress} "kubectl create -f https://raw.githubusercontent.com/projectcalico/calico/v3.26.0/manifests/custom-resources.yaml" >> $log`
#   calicoctl node status
#   calicoctl get node master -o yaml
  `ssh root@${advertiseAddress} "kubectl apply -f /var/lib/kubeStone/ks-serviceaccount.yaml" >> $log`
  `ssh root@${advertiseAddress} "kubectl apply -f /var/lib/kubeStone/ks-rolebings.yaml" >> $log`
)

for cmd in "${masterInstall[@]}"; do
      echo "Executing $cmd" >> $log
      eval $cmd
      if [ $? != 0 ]; then
          exit -1
      fi
  done

TOKEN=$(ssh root@${advertiseAddress} "kubectl create token kubestone-service-account --duration 87600h")
kubectlConfig=(
  `echo "$advertiseAddress: $TOKEN" >> /var/lib/kubeStone/TOKEN`
  "kubectl config  set-cluster cluster1 --server=https://${advertiseAddress}:6443 --insecure-skip-tls-verify=true"
  "kubectl config set-credentials kubestone-service-account --token=$TOKEN"
  "kubectl config  set-context context1 --cluster=cluster1 --user=kubestone-service-account"
  "kubectl config use-context context1"
)
for cmd in "${kubectlConfig[@]}"; do
      echo "Executing $cmd" >> $log
      eval $cmd
      if [ $? != 0 ]; then
          exit -1
      fi
done
