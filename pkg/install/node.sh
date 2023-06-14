#!/bin/bash

if [ $# != 3 ]; then
    exit -1
fi

masterIp=$1
nodeIp=$2
seq=$3
pkgPath="/root/kubeStone/offlinePKG"
shPath="/root/kubeStone/install"
masterHost=$(ssh root@${masterIp} hostname)
log=/root/node_${nodeIp}.log
joinCmd=$(ssh root@${masterIp} "kubeadm token create --print-join-command")

token=$(grep "token:" ${masterIp}.yaml  | awk '{print $2}')
if [ $? != 0 ]; then
    exit -1
fi

setHostForM="echo $nodeIp node1"
serHostForN="echo $masterIp master"
nodeSet=(
  `ssh root@${masterIp} "$setHostForM >> /etc/hosts"`
  `ssh root@${nodeIp} "$serHostForN >> /etc/hosts"`
  `ssh root@${nodeIp} "mkdir -p /var/lib/kubeStone/offlinePKG"`
  `scp ${shPath}/node.sh root@${nodeIp}:/var/lib/kubeStone/`
  `scp ${shPath}/preSetEnv.sh root@${nodeIp}:/var/lib/kubeStone/`
  `scp -r ${pkgPath}/* root@${nodeIp}:/var/lib/kubeStone/offlinePKG/`
  `ssh root@${nodeIp} "sh /var/lib/kubeStone/preSetEnv.sh node${seq}" > $log`
)
for cmd in "${nodeSet[@]}"; do
      eval $cmd
      if [ $? != 0 ]; then
          exit -1
      fi
  done

echo "Executing: $joinCmd" >> $log
ssh root@${nodeIp}  ${joinCmd} >> $log
if [ $? != 0 ]; then
          exit -1
fi


