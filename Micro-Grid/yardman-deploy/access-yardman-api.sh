# !/bin/bash
source ./script/config.sh
namespace=${hlledger_yardman_name}
pod=$(kubectl get pods --namespace=$namespace | grep "${namespace}" | cut -d" " -f 1)
kubectl exec --namespace=$namespace ${pod} --container=yardman-api -it -- bash
