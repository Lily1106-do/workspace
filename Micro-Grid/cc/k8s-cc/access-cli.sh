#!/bin/bash
namespace=conandpro-cli
pod=$(kubectl get pods --namespace=$namespace | grep "conandpro-cli" | cut -d" " -f 1)
kubectl exec --namespace=${namespace} ${pod} --container=cli -it -- bash
