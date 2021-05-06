#!/bin/bash
namespace=mychannel-cli
pod=$(kubectl get pods --namespace=$namespace | grep "mychannel-cli" | cut -d" " -f 1)
kubectl exec --namespace=${namespace} ${pod} --container=cli -it -- bash
