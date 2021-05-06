#!/bin/bash
source ./script/config.sh

# mount data
bash ./script/mount-data.sh
# create yardman
kubectl create -f yardman-deployment.yaml
