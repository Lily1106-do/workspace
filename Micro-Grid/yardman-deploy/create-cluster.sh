#!/bin/bash
source ./script/config.sh

# mount data 使用NFS挂载网络硬盘
bash ./script/mount-data.sh
# create yardman 创建microgrid集群
kubectl create -f yardman-deployment.yaml
