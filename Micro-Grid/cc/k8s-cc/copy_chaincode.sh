# !bin/bash
set -x
namespace=conandpro-cli
pod=$(kubectl get pods --namespace=$namespace | grep "conandpro-cli" | cut -d" " -f 1)
kubectl cp ../chaincode -n conandpro-cli ${pod}:/opt/gopath/src/github.com/hyperledger/fabric/chaincode/
kubectl cp ../scripts/excute_chaincode.sh -n conandpro-cli ${pod}:/opt/gopath/src/github.com/hyperledger/fabric/peer/scripts/
