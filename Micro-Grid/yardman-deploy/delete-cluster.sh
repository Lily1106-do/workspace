#!/bin/bash
source ./script/config.sh

sudo umount ../data
command="./script/hlledger-delete.py ${hlledger_yardman_name} ${hlledger_cluster_name}"

data_dir="../data/${hlledger_yardman_name}"

# delete data
sudo rm -rf ${data_dir}
sudo rm -rd /opt/share/*

# delete yardman
kubectl delete -f yardman-deployment.yaml

# delete cluster
python3 ${command}
# delete namespace 
kubectl delete ns microgrid-cluster-company
kubectl delete ns microgrid-cluster-consumer
kubectl delete ns microgrid-cluster-orderer
kubectl delete ns microgrid-cluster-producer
